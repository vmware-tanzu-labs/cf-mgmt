// Package config provides utilities for reading and writing cf-mgmt's configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

// Manager can read and write the cf-mgmt configuration.
type Manager interface {
	Updater
	Reader
}

// Updater is used to update the cf-mgmt configuration.
type Updater interface {
	AddOrgToConfig(orgConfig *OrgConfig) error
	AddSpaceToConfig(spaceConfig *SpaceConfig) error
	CreateConfigIfNotExists(uaaOrigin string) error
	DeleteConfigIfExists() error
}

// Reader is used to read the cf-mgmt configuration.
type Reader interface {
	Orgs() (Orgs, error)
	GetOrgConfigs() ([]OrgConfig, error)
}

// yamlManager is the default implementation of Manager.
// It is backed by a directory of YAML files.
type yamlManager struct {
	ConfigDir string
}

// NewManager creates a Manager that is backed by a set of YAML
// files in the specified configuration directory.
func NewManager(configDir string) Manager {
	return &yamlManager{
		ConfigDir: configDir,
	}
}

// Orgs reads the config for all orgs.
func (m *yamlManager) Orgs() (Orgs, error) {
	configFile := filepath.Join(m.ConfigDir, "orgs.yml")
	lo.G.Info("Processing org file", configFile)
	input := Orgs{}
	if err := utils.NewDefaultManager().LoadFile(configFile, &input); err != nil {
		return Orgs{}, err
	}
	return input, nil
}

// GetOrgConfigs reads all orgs from the cf-mgmt configuration.
func (m *yamlManager) GetOrgConfigs() ([]OrgConfig, error) {
	fs := utils.NewDefaultManager()
	files, err := fs.FindFiles(m.ConfigDir, "orgConfig.yml")
	if err != nil {
		return nil, err
	}
	result := make([]OrgConfig, len(files))
	for i, f := range files {
		result[i].AppInstanceLimit = -1
		result[i].TotalReservedRoutePorts = 0
		result[i].TotalPrivateDomains = -1
		result[i].TotalServiceKeys = -1

		if err = fs.LoadFile(f, &result[i]); err != nil {
			lo.G.Error(err)
			return nil, err
		}
	}
	return result, nil
}

// AddOrgToConfig adds an organization to the cf-mgmt configuration.
func (m *yamlManager) AddOrgToConfig(orgConfig *OrgConfig) error {
	orgFileName := filepath.Join(m.ConfigDir, "orgs.yml")
	orgName := orgConfig.Org
	if orgName == "" {
		return errors.New("cannot have an empty org name")
	}

	mgr := utils.NewDefaultManager()
	orgList := &Orgs{}
	err := mgr.LoadFile(orgFileName, orgList)
	if err != nil {
		return err
	}

	if orgList.Contains(orgName) {
		lo.G.Infof("%s already added to config", orgName)
		return nil
	}
	lo.G.Infof("Adding org: %s ", orgName)
	orgList.Orgs = append(orgList.Orgs, orgName)
	if err = mgr.WriteFile(orgFileName, orgList); err != nil {
		return err
	}

	if err = os.MkdirAll(fmt.Sprintf("%s/%s", m.ConfigDir, orgName), 0755); err != nil {
		return err
	}
	orgConfig.RemoveUsers = true
	orgConfig.RemovePrivateDomains = true
	mgr.WriteFile(filepath.Join(m.ConfigDir, orgName, "orgConfig.yml"), orgConfig)
	return mgr.WriteFile(filepath.Join(m.ConfigDir, orgName, "spaces.yml"), &Spaces{
		Org:                orgName,
		EnableDeleteSpaces: true,
	})
}

// AddSpaceToConfig adds a space to the cf-mgmt configuration, so long as a
// space with the specified name doesn't already exist.
func (m *yamlManager) AddSpaceToConfig(spaceConfig *SpaceConfig) error {
	orgName := spaceConfig.Org
	spaceFileName := filepath.Join(m.ConfigDir, orgName, "spaces.yml")
	spaceList := &Spaces{}
	spaceName := spaceConfig.Space
	mgr := utils.NewDefaultManager()

	if err := mgr.LoadFile(spaceFileName, spaceList); err != nil {
		return err
	}
	if spaceList.Contains(spaceName) {
		lo.G.Infof("%s already added to config", spaceName)
		return nil
	}
	lo.G.Infof("Adding space: %s ", spaceName)
	spaceList.Spaces = append(spaceList.Spaces, spaceName)
	if err := mgr.WriteFile(spaceFileName, spaceList); err != nil {
		return err
	}
	if err := os.MkdirAll(fmt.Sprintf("%s/%s/%s", m.ConfigDir, orgName, spaceName), 0755); err != nil {
		return err
	}
	spaceConfig.RemoveUsers = true

	mgr.WriteFile(fmt.Sprintf("%s/%s/%s/spaceConfig.yml", m.ConfigDir, orgName, spaceName), spaceConfig)
	mgr.WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.ConfigDir, orgName, spaceName), []byte("[]"))
	return nil
}

// CreateConfigIfNotExists initializes a new configuration directory.
// If the specified configuration directory already exists, it is left unmodified.
func (m *yamlManager) CreateConfigIfNotExists(uaaOrigin string) error {
	mgr := utils.NewDefaultManager()
	if mgr.FileOrDirectoryExists(m.ConfigDir) {
		lo.G.Infof("Config directory %s already exists, skipping creation", m.ConfigDir)
		return nil
	}
	if err := os.MkdirAll(m.ConfigDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", m.ConfigDir, err)
		return fmt.Errorf("cannot create directory %s: %v", m.ConfigDir, err)
	}
	lo.G.Infof("Config directory %s created", m.ConfigDir)
	mgr.WriteFile(fmt.Sprintf("%s/ldap.yml", m.ConfigDir), &ldap.Config{TLS: false, Origin: uaaOrigin})
	mgr.WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &Orgs{
		EnableDeleteOrgs: true,
		ProtectedOrgs:    []string{"system"},
	})
	mgr.WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), struct {
		Developer UserMgmt `yaml:"space-developer"`
		Manager   UserMgmt `yaml:"space-manager"`
		Auditor   UserMgmt `yaml:"space-auditor"`
	}{})
	return nil
}

// DeleteConfigIfExists deletes config directory if it exists.
func (m *yamlManager) DeleteConfigIfExists() error {
	utilsManager := utils.NewDefaultManager()
	if !utilsManager.FileOrDirectoryExists(m.ConfigDir) {
		lo.G.Infof("%s doesn't exists, nothing to delete", m.ConfigDir)
		return nil
	}
	if err := os.RemoveAll(m.ConfigDir); err != nil {
		lo.G.Errorf("Error deleting config folder. Error: %s", err)
		return fmt.Errorf("cannot delete %s: %v", m.ConfigDir, err)
	}
	lo.G.Info("Config directory deleted")
	return nil
}

// UserMgmt specifies users and groups that can be associated to a particular org or space.
type UserMgmt struct {
	LDAPUsers  []string `yaml:"ldap_users"`
	Users      []string `yaml:"users"`
	LDAPGroup  string   `yaml:"ldap_group"`
	LDAPGroups []string `yaml:"ldap_groups"`
}

func (u *UserMgmt) groups(groupName string) []string {
	groupMap := make(map[string]string)
	for _, group := range u.LDAPGroups {
		groupMap[group] = group
	}
	if u.LDAPGroup != "" {
		groupMap[u.LDAPGroup] = u.LDAPGroup
	}
	if groupName != "" {
		groupMap[groupName] = groupName
	}

	result := make([]string, 0, len(groupMap))
	for k := range groupMap {
		result = append(result, k)
	}
	return result
}
