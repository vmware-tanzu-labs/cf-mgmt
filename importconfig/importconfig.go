package importconfig

import (
	"fmt"
	"os"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
)

func NewManager(
	configDir string,
	uaacMgr uaa.Manager,
	orgMgr organization.Manager,
	spaceMgr space.Manager,
	cloudController cloudcontroller.Manager) Manager {
	return &DefaultImportManager{
		ConfigDir:       configDir,
		UAACMgr:         uaacMgr,
		OrgMgr:          orgMgr,
		SpaceMgr:        spaceMgr,
		CloudController: cloudController,
	}
}

func (im *DefaultImportManager) ImportConfig(excludedOrgs map[string]string) error {
	var err error
	var orgs []*cloudcontroller.Org
	var configMgr config.Manager
	orgs, err = im.CloudController.ListOrgs()
	if err != nil && len(orgs) > 0 {
		configMgr = config.NewManager(im.ConfigDir)
		err = configMgr.CreateConfigIfNotExists()
		if err != nil {
			for _, org := range orgs {
				if _, ok := excludedOrgs[org.Entity.Name]; !ok {
					var orgUsers map[string]string
					orgUsers, err = im.CloudController.GetCFUsers(org.MetaData.GUID, organization.ORGS, organization.ROLE_ORG_AUDITORS)
					if err != nil && len(orgUsers) > 0 {

					} else {
						fmt.Fprintf(os.Stdout, "No org auditors found for org : %s", org.Entity.Name)
					}
					orgConfig := &config.OrgConfig{OrgName: org.Entity.Name}
					configMgr.AddOrgToConfig(orgConfig)
				}
			}
		} else {
			fmt.Fprintf(os.Stdout, "Unable to create config directory : %s", im.ConfigDir)
			return err
		}

	}

	return nil
}
