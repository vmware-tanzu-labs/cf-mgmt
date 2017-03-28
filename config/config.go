package config

import (
	"fmt"
	"os"

	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/utils"
)

//Manager -
type Manager interface {
	AddOrgToConfig(orgConfig *OrgConfig) (err error)
	AddSpaceToConfig(spaceConfig *SpaceConfig) (err error)
}

//DefaultManager -
type DefaultManager struct {
	Config string
}

//OrgConfig Describes attributes for an org
type OrgConfig struct {
	OrgName          string
	OrgBillingMgrGrp string
	OrgMgrGrp        string
	OrgAuditorGrp    string
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
func NewManager(config string) Manager {
	return &DefaultManager{
		Config: config,
	}
}

//AddOrgToConfig -
func (m *DefaultManager) AddOrgToConfig(orgConfig *OrgConfig) (err error) {
	orgList := &organization.InputOrgs{}
	orgFileName := fmt.Sprintf("%s/orgs.yml", m.Config)
	orgName := orgConfig.OrgName

	if err = utils.NewDefaultManager().LoadFile(orgFileName, orgList); err == nil {
		if orgList.Contains(orgName) {
			fmt.Println(orgName, "already added to config")
		} else {
			fmt.Println("Adding org", orgName)
			orgList.Orgs = append(orgList.Orgs, orgName)
			if err = utils.NewDefaultManager().WriteFile(orgFileName, orgList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s", m.Config, orgName), 0755); err == nil {
					orgConfigYml := &organization.InputUpdateOrgs{
						Org:                     orgName,
						BillingManager:          organization.UserMgmt{LdapGroup: orgConfig.OrgBillingMgrGrp},
						Manager:                 organization.UserMgmt{LdapGroup: orgConfig.OrgMgrGrp},
						Auditor:                 organization.UserMgmt{LdapGroup: orgConfig.OrgAuditorGrp},
						EnableOrgQuota:          false,
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
						RemoveUsers:             true,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/orgConfig.yml", m.Config, orgName), orgConfigYml)
					spaces := &space.InputCreateSpaces{
						Org: orgName,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/spaces.yml", m.Config, orgName), spaces)
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
	spaceFileName := fmt.Sprintf("%s/%s/spaces.yml", m.Config, orgName)
	if err = utils.NewDefaultManager().LoadFile(spaceFileName, spaceList); err == nil {
		if spaceList.Contains(spaceName) {
			fmt.Println(spaceName, "already added to config")
		} else {
			fmt.Println("Adding space", spaceName)
			spaceList.Spaces = append(spaceList.Spaces, spaceName)
			if err = utils.NewDefaultManager().WriteFile(spaceFileName, spaceList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s/%s", m.Config, orgName, spaceName), 0755); err == nil {
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
						RemoveUsers:             true,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/%s/spaceConfig.yml", m.Config, orgName, spaceName), spaceConfigYml)
					utils.NewDefaultManager().WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.Config, orgName, spaceName), []byte("[]"))
				}
			}
		}
	}
	return
}
