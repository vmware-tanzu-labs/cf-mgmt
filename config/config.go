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
	AddOrgToConfig(inputOrg string) (err error)
	AddSpaceToConfig(inputOrg, inputSpace string) (err error)
}

//DefaultManager -
type DefaultManager struct {
	Config string
}

//NewManager -
func NewManager(config string) Manager {
	return &DefaultManager{
		Config: config,
	}
}

//AddOrgToConfig -
func (m *DefaultManager) AddOrgToConfig(inputOrg string) (err error) {
	orgList := &organization.InputOrgs{}
	orgFileName := fmt.Sprintf("%s/orgs.yml", m.Config)
	if err = utils.NewDefaultManager().LoadFile(orgFileName, orgList); err == nil {
		if orgList.Contains(inputOrg) {
			fmt.Println(inputOrg, "already added to config")
		} else {
			fmt.Println("Adding org", inputOrg)
			orgList.Orgs = append(orgList.Orgs, inputOrg)
			if err = utils.NewDefaultManager().WriteFile(orgFileName, orgList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s", m.Config, inputOrg), 0755); err == nil {
					orgConfig := &organization.InputUpdateOrgs{
						Org:                     inputOrg,
						EnableOrgQuota:          false,
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/orgConfig.yml", m.Config, inputOrg), orgConfig)
					spaces := &space.InputCreateSpaces{
						Org: inputOrg,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/spaces.yml", m.Config, inputOrg), spaces)
				}
			}
		}
	}
	return
}

//AddSpaceToConfig -
func (m *DefaultManager) AddSpaceToConfig(inputOrg, inputSpace string) (err error) {
	spaceList := &space.InputCreateSpaces{}
	spaceFileName := fmt.Sprintf("%s/%s/spaces.yml", m.Config, inputOrg)
	if err = utils.NewDefaultManager().LoadFile(spaceFileName, spaceList); err == nil {
		if spaceList.Contains(inputSpace) {
			fmt.Println(inputSpace, "already added to config")
		} else {
			fmt.Println("Adding space", inputSpace)
			spaceList.Spaces = append(spaceList.Spaces, inputSpace)
			if err = utils.NewDefaultManager().WriteFile(spaceFileName, spaceList); err == nil {
				if err = os.MkdirAll(fmt.Sprintf("%s/%s/%s", m.Config, inputOrg, inputSpace), 0755); err == nil {
					spaceConfig := &space.InputUpdateSpaces{
						Org:                     inputOrg,
						Space:                   inputSpace,
						EnableSpaceQuota:        false,
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
					}
					utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/%s/%s/spaceConfig.yml", m.Config, inputOrg, inputSpace), spaceConfig)
					utils.NewDefaultManager().WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.Config, inputOrg, inputSpace), []byte("[]"))
				}
			}
		}
	}
	return
}
