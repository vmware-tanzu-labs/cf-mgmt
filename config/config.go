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
func (m *DefaultManager) AddOrgToConfig(orgConfig *OrgConfig) (err error) {
	orgList := &organization.InputOrgs{}
	orgFileName := fmt.Sprintf("%s/orgs.yml", m.ConfigDir)
	orgName := orgConfig.OrgName
	if orgName == "" {
		err = errors.New("Cannot have an empty org name")
		return
	}

	orgQuota := orgConfig.OrgQuota
	if err = utils.NewDefaultManager().LoadFile(orgFileName, orgList); err == nil {
		if orgList.Contains(orgName) {
			lo.G.Infof("%s already added to config", orgName)
		} else {
			lo.G.Infof("Adding org: %s ", orgName)
			orgList.Orgs = append(orgList.Orgs, orgName)
			if err = utils.NewDefaultManager().WriteFile(orgFileName, orgList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s", m.ConfigDir, orgName), 0755); err == nil {
					orgConfigYml := &organization.InputUpdateOrgs{
						Org:                     orgName,
						BillingManager:          organization.UserMgmt{LdapGroup: orgConfig.OrgBillingMgrLDAPGrp, Users: orgConfig.OrgBillingMgrUAAUsers, LdapUsers: orgConfig.OrgBillingMgrLDAPUsers},
						Manager:                 organization.UserMgmt{LdapGroup: orgConfig.OrgMgrLDAPGrp, Users: orgConfig.OrgMgrUAAUsers, LdapUsers: orgConfig.OrgMgrLDAPUsers},
						Auditor:                 organization.UserMgmt{LdapGroup: orgConfig.OrgAuditorLDAPGrp, Users: orgConfig.OrgAuditorUAAUsers, LdapUsers: orgConfig.OrgAuditorLDAPUsers},
						EnableOrgQuota:          orgQuota.IsQuotaEnabled(),
						MemoryLimit:             orgQuota.GetMemoryLimit(),
						InstanceMemoryLimit:     orgQuota.GetInstanceMemoryLimit(),
						TotalRoutes:             orgQuota.GetTotalRoutes(),
						TotalServices:           orgQuota.GetTotalServices(),
						PaidServicePlansAllowed: orgQuota.IsPaidServicesAllowed(),
						RemoveUsers:             true,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/orgConfig.yml", m.ConfigDir, orgName), orgConfigYml)
					spaces := &space.InputCreateSpaces{
						Org: orgName,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/spaces.yml", m.ConfigDir, orgName), spaces)
				}
			}
		}
	}
	return
}

//AddSpaceToConfig -
func (m *DefaultManager) AddSpaceToConfig(spaceConfig *SpaceConfig) (err error) {
	spaceList := &space.InputCreateSpaces{}
	spaceName := spaceConfig.SpaceName
	orgName := spaceConfig.OrgName
	spaceFileName := fmt.Sprintf("%s/%s/spaces.yml", m.ConfigDir, orgName)
	spaceQuota := spaceConfig.SpaceQuota
	if err = utils.NewDefaultManager().LoadFile(spaceFileName, spaceList); err == nil {
		if spaceList.Contains(spaceName) {
			lo.G.Infof("%s already added to config", spaceName)
		} else {
			lo.G.Infof("Adding space: %s ", spaceName)
			spaceList.Spaces = append(spaceList.Spaces, spaceName)
			if err = utils.NewDefaultManager().WriteFile(spaceFileName, spaceList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s/%s", m.ConfigDir, orgName, spaceName), 0755); err == nil {
					spaceConfigYml := &space.InputUpdateSpaces{
						Org:                     orgName,
						Space:                   spaceName,
						Developer:               space.UserMgmt{LdapGroup: spaceConfig.SpaceDevLDAPGrp, Users: spaceConfig.SpaceDevUAAUsers, LdapUsers: spaceConfig.SpaceDevLDAPUsers},
						Manager:                 space.UserMgmt{LdapGroup: spaceConfig.SpaceMgrLDAPGrp, Users: spaceConfig.SpaceMgrUAAUsers, LdapUsers: spaceConfig.SpaceMgrLDAPUsers},
						Auditor:                 space.UserMgmt{LdapGroup: spaceConfig.SpaceAuditorLDAPGrp, Users: spaceConfig.SpaceAuditorUAAUsers, LdapUsers: spaceConfig.SpaceAuditorLDAPUsers},
						EnableSpaceQuota:        spaceQuota.IsQuotaEnabled(),
						MemoryLimit:             spaceQuota.GetMemoryLimit(),
						InstanceMemoryLimit:     spaceQuota.GetInstanceMemoryLimit(),
						TotalRoutes:             spaceQuota.GetTotalRoutes(),
						TotalServices:           spaceQuota.GetTotalServices(),
						PaidServicePlansAllowed: spaceQuota.IsPaidServicesAllowed(),
						RemoveUsers:             true,
						AllowSSH:                spaceConfig.AllowSSH,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/%s/spaceConfig.yml", m.ConfigDir, orgName, spaceName), spaceConfigYml)
					utils.NewDefaultManager().WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.ConfigDir, orgName, spaceName), []byte("[]"))
				}
			}
		}
	}
	return
}

//CreateConfigIfNotExists Create org and space config directory. If directory already exists, it is left as is
func (m *DefaultManager) CreateConfigIfNotExists(uaaOrigin string) error {
	var err error
	utilsManager := utils.NewDefaultManager()
	if !utilsManager.FileOrDirectoryExists(m.ConfigDir) {
		if err = os.MkdirAll(m.ConfigDir, 0755); err == nil {
			lo.G.Infof("Config directory %s created", m.ConfigDir)
			utilsManager.WriteFile(fmt.Sprintf("%s/ldap.yml", m.ConfigDir), &ldap.Config{TLS: false, Origin: uaaOrigin})
			utilsManager.WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &organization.InputOrgs{})
			utilsManager.WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), &space.ConfigSpaceDefaults{})
		} else {
			lo.G.Errorf("Error creating config directory %s. Error : %s", m.ConfigDir, err)
		}
	} else {
		lo.G.Infof("Config directory %s already exists, skipping creation", m.ConfigDir)
	}
	return err
}

//DeleteConfigIfExists Deletes config directory if it exists
func (m *DefaultManager) DeleteConfigIfExists() error {
	var err error
	utilsManager := utils.NewDefaultManager()
	if utilsManager.FileOrDirectoryExists(m.ConfigDir) {
		err = os.RemoveAll(m.ConfigDir)
		if err != nil {
			lo.G.Errorf("Error deleting config folder. Error : %s", err)
			return err
		}
		lo.G.Info("Config directory deleted")
	} else {
		lo.G.Infof("%s doesn't exists, nothing to delete", m.ConfigDir)
	}
	return err
}
