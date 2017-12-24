package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/xchapter7x/lo"
)

// yamlManager is the default implementation of Manager.
// It is backed by a directory of YAML files.
type yamlManager struct {
	ConfigDir string
}

// Orgs reads the config for all orgs.
func (m *yamlManager) Orgs() (Orgs, error) {
	configFile := filepath.Join(m.ConfigDir, "orgs.yml")
	lo.G.Debug("Processing org file", configFile)
	input := Orgs{}
	if err := LoadFile(configFile, &input); err != nil {
		return Orgs{}, err
	}
	return input, nil
}

// GetASGConfigs reads all ASGs from the cf-mgmt configuration.
func (m *yamlManager) GetASGConfigs() ([]ASGConfig, error) {

	files, err := FindFiles(path.Join(m.ConfigDir, "asgs"), ".json")
	if err != nil {
		return nil, err
	}
	var result []ASGConfig
	for _, securityGroupFile := range files {
		lo.G.Debug("Loading security group contents", securityGroupFile)
		bytes, err := ioutil.ReadFile(securityGroupFile)
		if err != nil {
			return nil, err
		}
		asgConfig := ASGConfig{}
		lo.G.Debug("setting security group contents", string(bytes))
		asgConfig.Rules = string(bytes)
		asgConfig.Name = strings.Replace(filepath.Base(securityGroupFile), ".json", "", 1)
		result = append(result, asgConfig)
	}
	return result, nil
}

// GetIsolationSegmentConfig reads isolation segment config
func (m *yamlManager) GetGlobalConfig() (GlobalConfig, error) {
	globalConfig := &GlobalConfig{}
	LoadFile(path.Join(m.ConfigDir, "cf-mgmt.yml"), globalConfig)
	return *globalConfig, nil
}

// GetOrgConfigs reads all orgs from the cf-mgmt configuration.
func (m *yamlManager) GetOrgConfigs() ([]OrgConfig, error) {
	files, err := FindFiles(m.ConfigDir, "orgConfig.yml")
	if err != nil {
		return nil, err
	}
	result := make([]OrgConfig, len(files))
	for i, f := range files {
		result[i].AppInstanceLimit = -1
		result[i].TotalReservedRoutePorts = 0
		result[i].TotalPrivateDomains = -1
		result[i].TotalServiceKeys = -1

		if err = LoadFile(f, &result[i]); err != nil {
			lo.G.Error(err)
			return nil, err
		}
	}
	return result, nil
}

func (m *yamlManager) spaceList(path string) ([]Spaces, error) {
	files, err := FindFiles(path, "spaces.yml")
	if err != nil {
		return nil, err
	}

	spaceList := make([]Spaces, len(files))
	for i, f := range files {
		lo.G.Debug("Processing space file", f)

		if err = LoadFile(f, &spaceList[i]); err != nil {
			lo.G.Errorf("reading config for space %s: %v", f, err)
			return nil, err
		}
	}
	return spaceList, nil
}
func (m *yamlManager) Spaces() ([]Spaces, error) {
	return m.spaceList(m.ConfigDir)
}
func (m *yamlManager) OrgSpaces(orgName string) (*Spaces, error) {
	spaceList, err := m.spaceList(path.Join(m.ConfigDir, orgName))
	if err != nil {
		return nil, err
	}
	if len(spaceList) == 1 {
		return &spaceList[0], nil
	}
	return nil, fmt.Errorf("No spaces found for org [%s]", orgName)
}

func (m *yamlManager) GetSpaceConfigs() ([]SpaceConfig, error) {

	spaceDefaults := SpaceConfig{}
	LoadFile(filepath.Join(m.ConfigDir, "spaceDefaults.yml"), &spaceDefaults)

	files, err := FindFiles(m.ConfigDir, "spaceConfig.yml")
	if err != nil {
		return nil, err
	}
	result := make([]SpaceConfig, len(files))
	for i, f := range files {
		result[i].AppInstanceLimit = -1
		result[i].TotalReservedRoutePorts = 0
		result[i].TotalPrivateDomains = -1
		result[i].TotalServiceKeys = -1

		if err = LoadFile(f, &result[i]); err != nil {
			return nil, err
		}

		result[i].Developer.LDAPUsers = append(result[i].Developer.LDAPUsers, spaceDefaults.Developer.LDAPUsers...)
		result[i].Developer.Users = append(result[i].Developer.Users, spaceDefaults.Developer.Users...)
		result[i].Developer.SamlUsers = append(result[i].Developer.SamlUsers, spaceDefaults.Developer.SamlUsers...)

		result[i].Auditor.LDAPUsers = append(result[i].Auditor.LDAPUsers, spaceDefaults.Auditor.LDAPUsers...)
		result[i].Auditor.Users = append(result[i].Auditor.Users, spaceDefaults.Auditor.Users...)
		result[i].Auditor.SamlUsers = append(result[i].Auditor.SamlUsers, spaceDefaults.Auditor.SamlUsers...)

		result[i].Manager.LDAPUsers = append(result[i].Manager.LDAPUsers, spaceDefaults.Manager.LDAPUsers...)
		result[i].Manager.Users = append(result[i].Manager.Users, spaceDefaults.Manager.Users...)
		result[i].Manager.SamlUsers = append(result[i].Manager.SamlUsers, spaceDefaults.Manager.SamlUsers...)

		result[i].Developer.LDAPGroups = append(result[i].GetDeveloperGroups(), spaceDefaults.GetDeveloperGroups()...)
		result[i].Auditor.LDAPGroups = append(result[i].GetAuditorGroups(), spaceDefaults.GetAuditorGroups()...)
		result[i].Manager.LDAPGroups = append(result[i].GetManagerGroups(), spaceDefaults.GetManagerGroups()...)

		if result[i].EnableSecurityGroup {
			securityGroupFile := strings.Replace(f, "spaceConfig.yml", "security-group.json", -1)
			lo.G.Debug("Loading security group contents", securityGroupFile)
			bytes, err := ioutil.ReadFile(securityGroupFile)
			if err != nil {
				return nil, err
			}
			lo.G.Debug("setting security group contents", string(bytes))
			result[i].SecurityGroupContents = string(bytes)
		}
	}
	return result, nil
}

func (m *yamlManager) GetOrgConfig(orgName string) (*OrgConfig, error) {
	configs, err := m.GetOrgConfigs()
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		if config.Org == orgName {
			return &config, nil
		}
	}
	return nil, fmt.Errorf("Org [%s] not found in config", orgName)
}

func (m *yamlManager) SaveOrgConfig(orgConfig *OrgConfig) error {
	directory := fmt.Sprintf("%s/%s", m.ConfigDir, orgConfig.Org)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return err
		}
	}

	return WriteFile(filepath.Join(m.ConfigDir, orgConfig.Org, "orgConfig.yml"), orgConfig)
}

func (m *yamlManager) GetSpaceConfig(orgName, spaceName string) (*SpaceConfig, error) {
	configs, err := m.GetSpaceConfigs()
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		if config.Org == orgName && config.Space == spaceName {
			return &config, nil
		}
	}
	return nil, fmt.Errorf("Space [%s] not found in org [%s] config", spaceName, orgName)
}

func (m *yamlManager) SaveSpaceConfig(spaceConfig *SpaceConfig) error {
	if err := os.MkdirAll(fmt.Sprintf("%s/%s/%s", m.ConfigDir, spaceConfig.Org, spaceConfig.Space), 0755); err != nil {
		return err
	}
	return WriteFile(fmt.Sprintf("%s/%s/%s/spaceConfig.yml", m.ConfigDir, spaceConfig.Org, spaceConfig.Space), spaceConfig)
}

func (m *yamlManager) DeleteOrgConfig(orgName string) error {
	orgs, err := m.Orgs()
	if err != nil {
		return err
	}
	if orgs.Contains(orgName) {
		var orgList []string
		for _, org := range orgs.Orgs {
			if org != orgName {
				orgList = append(orgList, org)
			}
		}
		orgs.Orgs = orgList
		if err := m.saveOrgList(orgs); err != nil {
			return err
		}
		os.RemoveAll(path.Join(m.ConfigDir, orgName))
	}
	return nil
}

func (m *yamlManager) DeleteSpaceConfig(orgName, spaceName string) error {
	spaces, err := m.OrgSpaces(orgName)
	if err != nil {
		return err
	}
	if spaces.Contains(spaceName) {
		var spaceList []string
		for _, space := range spaces.Spaces {
			if space != spaceName {
				spaceList = append(spaceList, space)
			}
		}
		spaces.Spaces = spaceList
		if err := m.saveSpaceList(*spaces); err != nil {
			return err
		}
		os.RemoveAll(path.Join(m.ConfigDir, orgName, spaceName))
	}
	return nil
}

// GetSpaceDefaults returns the default space configuration, if one was provided.
// If no space defaults were configured, a nil config and a nil error are returned.
func (m *yamlManager) GetSpaceDefaults() (*SpaceConfig, error) {
	fp := filepath.Join(m.ConfigDir, "spaceDefaults.yml")

	if !FileOrDirectoryExists(fp) {
		return nil, nil
	}
	result := SpaceConfig{}
	err := LoadFile(fp, &result)
	return &result, err
}

func (m *yamlManager) saveOrgList(orgs Orgs) error {
	if err := WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), orgs); err != nil {
		return err
	}
	return nil
}

// AddOrgToConfig adds an organization to the cf-mgmt configuration.
func (m *yamlManager) AddOrgToConfig(orgConfig *OrgConfig) error {
	orgList, err := m.Orgs()
	if err != nil {
		return err
	}
	orgName := orgConfig.Org
	if orgName == "" {
		return errors.New("cannot have an empty org name")
	}

	if orgList.Contains(orgName) {
		lo.G.Infof("%s already added to config", orgName)
		return nil
	}
	lo.G.Infof("Adding org: %s ", orgName)
	orgList.Orgs = append(orgList.Orgs, orgName)
	if err = m.saveOrgList(orgList); err != nil {
		return err
	}
	m.SaveOrgConfig(orgConfig)
	return m.saveSpaceList(Spaces{
		Org:                orgName,
		EnableDeleteSpaces: true,
	})
}

func (m *yamlManager) saveSpaceList(spaces Spaces) error {
	return WriteFile(filepath.Join(m.ConfigDir, spaces.Org, "spaces.yml"), spaces)
}

// AddSpaceToConfig adds a space to the cf-mgmt configuration, so long as a
// space with the specified name doesn't already exist.
func (m *yamlManager) AddSpaceToConfig(spaceConfig *SpaceConfig) error {
	orgName := spaceConfig.Org
	spaceFileName := filepath.Join(m.ConfigDir, orgName, "spaces.yml")
	spaceList := &Spaces{}
	spaceName := spaceConfig.Space

	if err := LoadFile(spaceFileName, spaceList); err != nil {
		return err
	}
	if spaceList.Contains(spaceName) {
		lo.G.Infof("%s already added to config", spaceName)
		return nil
	}
	lo.G.Infof("Adding space: %s ", spaceName)
	spaceList.Spaces = append(spaceList.Spaces, spaceName)
	if err := WriteFile(spaceFileName, spaceList); err != nil {
		return err
	}

	if err := m.SaveSpaceConfig(spaceConfig); err != nil {
		return err
	}
	if err := WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.ConfigDir, orgName, spaceName), []byte("[]")); err != nil {
		return err
	}
	return nil
}

//AddSecurityGroupToSpace - adds security group json to org/space location
func (m *yamlManager) AddSecurityGroupToSpace(orgName, spaceName string, securityGroupDefinition []byte) error {
	return WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.ConfigDir, orgName, spaceName), securityGroupDefinition)
}

//AddSecurityGroupToSpace - adds security group json to org/space location
func (m *yamlManager) AddSecurityGroup(securityGroupName string, securityGroupDefinition []byte) error {
	lo.G.Infof("Writing out bytes for security group %s", securityGroupName)
	return WriteFileBytes(fmt.Sprintf("%s/asgs/%s.json", m.ConfigDir, securityGroupName), securityGroupDefinition)
}

// CreateConfigIfNotExists initializes a new configuration directory.
// If the specified configuration directory already exists, it is left unmodified.
func (m *yamlManager) CreateConfigIfNotExists(uaaOrigin string) error {
	if FileOrDirectoryExists(m.ConfigDir) {
		lo.G.Infof("Config directory %s already exists, skipping creation", m.ConfigDir)
		return nil
	}
	if err := os.MkdirAll(m.ConfigDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", m.ConfigDir, err)
		return fmt.Errorf("cannot create directory %s: %v", m.ConfigDir, err)
	}
	lo.G.Infof("Config directory %s created", m.ConfigDir)

	asgDir := path.Join(m.ConfigDir, "asgs")
	if err := os.MkdirAll(asgDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", asgDir, err)
		return fmt.Errorf("cannot create directory %s: %v", asgDir, err)
	}
	lo.G.Infof("ASG directory %s created", asgDir)

	if err := WriteFile(fmt.Sprintf("%s/cf-mgmt.yml", m.ConfigDir), &GlobalConfig{}); err != nil {
		return err
	}
	if err := WriteFile(fmt.Sprintf("%s/ldap.yml", m.ConfigDir), &ldap.Config{TLS: false, Origin: uaaOrigin}); err != nil {
		return err
	}

	if err := WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &Orgs{
		EnableDeleteOrgs: true,
		ProtectedOrgs:    DefaultProtectedOrgs,
	}); err != nil {
		return err
	}
	if err := WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), struct {
		Developer UserMgmt `yaml:"space-developer"`
		Manager   UserMgmt `yaml:"space-manager"`
		Auditor   UserMgmt `yaml:"space-auditor"`
	}{}); err != nil {
		return err
	}
	return nil
}

// DeleteConfigIfExists deletes config directory if it exists.
func (m *yamlManager) DeleteConfigIfExists() error {
	if !FileOrDirectoryExists(m.ConfigDir) {
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
