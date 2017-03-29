package config

import (
	"fmt"
	"os"

	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/utils"
)

//Manager -
type Manager interface {
	AddOrgToConfig(orgConfig *OrgConfig) (err error)
	AddSpaceToConfig(spaceConfig *SpaceConfig) (err error)
	CreateConfigIfNotExists() error
}

//DefaultManager -
type DefaultManager struct {
	ConfigDir string
}

//OrgConfig Describes attributes for an org
type OrgConfig struct {
	OrgName               string
	OrgBillingMgrLDAPGrp  string
	OrgMgrLDAPGrp         string
	OrgAuditorLDAPGrp     string
	OrgBillingMgrUAAUsers []string
	OrgMgrUAAUsers        []string
	OrgAuditorUAAUsers    []string
}

//SpaceConfig Describes attributes for a space
type SpaceConfig struct {
	OrgName         string
	SpaceName       string
	SpaceDevGrp     string
	SpaceMgrGrp     string
	SpaceAuditorGrp string
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

	if err = utils.NewDefaultManager().LoadFile(orgFileName, orgList); err == nil {
		if orgList.Contains(orgName) {
			fmt.Println(orgName, "already added to config")
		} else {
			fmt.Println("Adding org", orgName)
			orgList.Orgs = append(orgList.Orgs, orgName)
			if err = utils.NewDefaultManager().WriteFile(orgFileName, orgList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s", m.ConfigDir, orgName), 0755); err == nil {
					orgConfigYml := &organization.InputUpdateOrgs{
						Org:                     orgName,
						BillingManager:          organization.UserMgmt{LdapGroup: orgConfig.OrgBillingMgrLDAPGrp, Users: orgConfig.OrgBillingMgrUAAUsers},
						Manager:                 organization.UserMgmt{LdapGroup: orgConfig.OrgMgrLDAPGrp, Users: orgConfig.OrgMgrUAAUsers},
						Auditor:                 organization.UserMgmt{LdapGroup: orgConfig.OrgAuditorLDAPGrp, Users: orgConfig.OrgAuditorUAAUsers},
						EnableOrgQuota:          false,
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
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
	if err = utils.NewDefaultManager().LoadFile(spaceFileName, spaceList); err == nil {
		if spaceList.Contains(spaceName) {
			fmt.Println(spaceName, "already added to config")
		} else {
			fmt.Println("Adding space", spaceName)
			spaceList.Spaces = append(spaceList.Spaces, spaceName)
			if err = utils.NewDefaultManager().WriteFile(spaceFileName, spaceList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s/%s", m.ConfigDir, orgName, spaceName), 0755); err == nil {
					spaceConfigYml := &space.InputUpdateSpaces{
						Org:                     orgName,
						Space:                   spaceName,
						Developer:               space.UserMgmt{LdapGroup: spaceConfig.SpaceDevGrp},
						Manager:                 space.UserMgmt{LdapGroup: spaceConfig.SpaceMgrGrp},
						Auditor:                 space.UserMgmt{LdapGroup: spaceConfig.SpaceAuditorGrp},
						EnableSpaceQuota:        false,
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
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
func (m *DefaultManager) CreateConfigIfNotExists() error {
	var err error
	if !utils.NewDefaultManager().DoesFileOrDirectoryExists(m.ConfigDir) {
		if err = os.MkdirAll(m.ConfigDir, 0755); err == nil {
			utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/ldap.yml", m.ConfigDir), &ldap.Config{TLS: false, Origin: "ldap"})
			utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &organization.InputOrgs{})
			utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), &space.ConfigSpaceDefaults{})
		}
	} else {
		fmt.Println(m.ConfigDir, "already exists, skipping creation")
	}
	return err
}
