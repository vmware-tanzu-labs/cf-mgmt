// Package config provides utilities for reading and writing cf-mgmt's configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

// Manager is used to update the cf-mgmt configuration.
type Manager interface {
	AddOrgToConfig(orgConfig *OrgConfig) error
	AddSpaceToConfig(spaceConfig *SpaceConfig) error
	CreateConfigIfNotExists(uaaOrigin string) error
	DeleteConfigIfExists() error
}

// yamlManager is the default implementation of Manager.
// It is backed by a directory of YAML files.
type yamlManager struct {
	ConfigDir string
}

// OrgConfig describes attributes for an org.
type OrgConfig struct {
	OrgName                string
	OrgBillingMgrLDAPGrp   string
	OrgMgrLDAPGrp          string
	OrgAuditorLDAPGrp      string
	OrgBillingMgrUAAUsers  []string
	OrgMgrUAAUsers         []string
	OrgAuditorUAAUsers     []string
	OrgBillingMgrLDAPUsers []string
	OrgMgrLDAPUsers        []string
	OrgAuditorLDAPUsers    []string
	IsoSegments            []string
	DefaultIsoSegment      string
	OrgQuota               cloudcontroller.QuotaEntity
}

// SpaceConfig describes attributes for a space.
type SpaceConfig struct {
	OrgName               string
	SpaceName             string
	SpaceDevLDAPGrp       string
	SpaceMgrLDAPGrp       string
	SpaceAuditorLDAPGrp   string
	SpaceDevUAAUsers      []string
	SpaceMgrUAAUsers      []string
	SpaceAuditorUAAUsers  []string
	SpaceDevLDAPUsers     []string
	SpaceMgrLDAPUsers     []string
	SpaceAuditorLDAPUsers []string
	SpaceQuota            cloudcontroller.QuotaEntity
	AllowSSH              bool
}

// NewManager creates a Manager that is backed by a set of YAML
// files in the specified configuration directory.
func NewManager(configDir string) Manager {
	return &yamlManager{
		ConfigDir: configDir,
	}
}

// AddOrgToConfig adds an organization to the cf-mgmt configuration.
func (m *yamlManager) AddOrgToConfig(orgConfig *OrgConfig) error {
	orgFileName := filepath.Join(m.ConfigDir, "orgs.yml")
	orgName := orgConfig.OrgName
	if orgName == "" {
		return errors.New("cannot have an empty org name")
	}

	mgr := utils.NewDefaultManager()
	orgList := &organization.InputOrgs{}
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
	orgConfigYml := &organization.InputUpdateOrgs{
		Org:                     orgName,
		BillingManager:          newUserMgmt(orgConfig.OrgBillingMgrLDAPGrp, orgConfig.OrgBillingMgrUAAUsers, orgConfig.OrgBillingMgrLDAPUsers),
		Manager:                 newUserMgmt(orgConfig.OrgMgrLDAPGrp, orgConfig.OrgMgrUAAUsers, orgConfig.OrgMgrLDAPUsers),
		Auditor:                 newUserMgmt(orgConfig.OrgAuditorLDAPGrp, orgConfig.OrgAuditorUAAUsers, orgConfig.OrgAuditorLDAPUsers),
		EnableOrgQuota:          orgConfig.OrgQuota.IsQuotaEnabled(),
		MemoryLimit:             orgConfig.OrgQuota.GetMemoryLimit(),
		InstanceMemoryLimit:     orgConfig.OrgQuota.GetInstanceMemoryLimit(),
		TotalRoutes:             orgConfig.OrgQuota.GetTotalRoutes(),
		TotalServices:           orgConfig.OrgQuota.GetTotalServices(),
		PaidServicePlansAllowed: orgConfig.OrgQuota.IsPaidServicesAllowed(),
		RemoveUsers:             true,
		RemovePrivateDomains:    true,
		IsoSegments:             orgConfig.IsoSegments,
	}
	mgr.WriteFile(filepath.Join(m.ConfigDir, orgName, "orgConfig.yml"), orgConfigYml)
	return mgr.WriteFile(filepath.Join(m.ConfigDir, orgName, "spaces.yml"), &space.InputSpaces{
		Org:                orgName,
		EnableDeleteSpaces: true,
	})
}

func newUserMgmt(ldapGroup string, users, ldapUsers []string) organization.UserMgmt {
	return organization.UserMgmt{
		LdapGroup: ldapGroup,
		Users:     users,
		LdapUsers: ldapUsers,
	}
}

// AddSpaceToConfig adds a space to the cf-mgmt configuration, so long as a
// space with the specified name doesn't already exist.
func (m *yamlManager) AddSpaceToConfig(spaceConfig *SpaceConfig) error {
	orgName := spaceConfig.OrgName
	spaceFileName := filepath.Join(m.ConfigDir, orgName, "spaces.yml")
	spaceList := &space.InputSpaces{}
	spaceName := spaceConfig.SpaceName
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
	spaceConfigYml := &space.InputSpaceConfig{
		Org:                     orgName,
		Space:                   spaceName,
		Developer:               space.UserMgmt{LdapGroup: spaceConfig.SpaceDevLDAPGrp, Users: spaceConfig.SpaceDevUAAUsers, LdapUsers: spaceConfig.SpaceDevLDAPUsers},
		Manager:                 space.UserMgmt{LdapGroup: spaceConfig.SpaceMgrLDAPGrp, Users: spaceConfig.SpaceMgrUAAUsers, LdapUsers: spaceConfig.SpaceMgrLDAPUsers},
		Auditor:                 space.UserMgmt{LdapGroup: spaceConfig.SpaceAuditorLDAPGrp, Users: spaceConfig.SpaceAuditorUAAUsers, LdapUsers: spaceConfig.SpaceAuditorLDAPUsers},
		EnableSpaceQuota:        spaceConfig.SpaceQuota.IsQuotaEnabled(),
		MemoryLimit:             spaceConfig.SpaceQuota.GetMemoryLimit(),
		InstanceMemoryLimit:     spaceConfig.SpaceQuota.GetInstanceMemoryLimit(),
		TotalRoutes:             spaceConfig.SpaceQuota.GetTotalRoutes(),
		TotalServices:           spaceConfig.SpaceQuota.GetTotalServices(),
		PaidServicePlansAllowed: spaceConfig.SpaceQuota.IsPaidServicesAllowed(),
		RemoveUsers:             true,
		AllowSSH:                spaceConfig.AllowSSH,
	}
	mgr.WriteFile(fmt.Sprintf("%s/%s/%s/spaceConfig.yml", m.ConfigDir, orgName, spaceName), spaceConfigYml)
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
	mgr.WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &organization.InputOrgs{
		EnableDeleteOrgs: true,
		ProtectedOrgs:    []string{"system"},
	})
	mgr.WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), &space.ConfigSpaceDefaults{})
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
