package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//Manager -
type Manager interface {
	AddOrgToConfig(orgConfig *OrgConfig) (err error)
	AddSpaceToConfig(spaceConfig *SpaceConfig) (err error)
	CreateConfigIfNotExists(uaaOrigin string) error
	DeleteConfigIfExists() error
}

//DefaultManager -
type DefaultManager struct {
	ConfigDir string
}

//OrgConfig Describes attributes for an org
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
	OrgQuota               cloudcontroller.QuotaEntity
}

//SpaceConfig Describes attributes for a space
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

//NewManager -
func NewManager(configDir string) Manager {
	return &DefaultManager{
		ConfigDir: configDir,
	}
}

//AddOrgToConfig -
func (m *DefaultManager) AddOrgToConfig(orgConfig *OrgConfig) error {
	orgList := &organization.InputOrgs{}
	orgFileName := fmt.Sprintf("%s/orgs.yml", m.ConfigDir) // TODO filepath.Join
	orgName := orgConfig.OrgName
	if orgName == "" {
		return errors.New("cannot have an empty org name")
	}

	mgr := utils.NewDefaultManager()
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
	}
	mgr.WriteFile(fmt.Sprintf("%s/%s/orgConfig.yml", m.ConfigDir, orgName), orgConfigYml) // TODO: filepath.Join
	spaces := &space.InputCreateSpaces{
		Org: orgName,
	}
	mgr.WriteFile(fmt.Sprintf("%s/%s/spaces.yml", m.ConfigDir, orgName), spaces) // TODO: filepath.Join
	return nil
}

func newUserMgmt(ldapGroup string, users, ldapUsers []string) organization.UserMgmt {
	return organization.UserMgmt{
		LdapGroup: ldapGroup,
		Users:     users,
		LdapUsers: ldapUsers,
	}
}

//AddSpaceToConfig -
func (m *DefaultManager) AddSpaceToConfig(spaceConfig *SpaceConfig) error {
	orgName := spaceConfig.OrgName
	spaceFileName := fmt.Sprintf("%s/%s/spaces.yml", m.ConfigDir, orgName) // TODO: filepath.Join
	spaceList := &space.InputCreateSpaces{}
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
	spaceConfigYml := &space.InputUpdateSpaces{
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

//CreateConfigIfNotExists Create org and space config directory. If directory already exists, it is left as is
func (m *DefaultManager) CreateConfigIfNotExists(uaaOrigin string) error {
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
	mgr.WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &organization.InputOrgs{})
	mgr.WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), &space.ConfigSpaceDefaults{})
	return nil
}

//DeleteConfigIfExists Deletes config directory if it exists
func (m *DefaultManager) DeleteConfigIfExists() error {
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
