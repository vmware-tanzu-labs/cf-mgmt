package importconfig

import (
	"fmt"
	"os"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaac"
)

const LDAP string = "ldap"
const SAML string = "saml"

func NewManager(
	configDir string,
	uaacMgr uaac.Manager,
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
	var userIDToUserMap map[string]uaac.User

	//Get all the users from the foundation
	userIDToUserMap, err = im.UAACMgr.UsersByID()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to retrieve users. Error : %s", err)
		return err
	}
	//Get all the orgs
	orgs, err = im.CloudController.ListOrgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to retrieve orgs. Error : %s", err)
		return err
	}

	configMgr = config.NewManager(im.ConfigDir)

	//Delete existing config directory
	configMgr.DeleteConfigIfExists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to delete config directory %s . Error:  %s", im.ConfigDir, err)
		return err
	}

	//Create a brand new directory
	err = configMgr.CreateConfigIfNotExists("ldap")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create config directory %s. Error : %s", im.ConfigDir, err)
		return err
	}

	for _, org := range orgs {
		if _, ok := excludedOrgs[org.Entity.Name]; !ok {
			orgConfig := &config.OrgConfig{OrgName: org.Entity.Name}
			//Get Org manager users for this org
			orgMgrs := getCFUsers(im.CloudController, org.MetaData.GUID, organization.ORGS, organization.ROLE_ORG_MANAGERS)
			for _, orgMgrUser := range orgMgrs {
				if usr, ok := userIDToUserMap[orgMgrUser]; ok {
					if usr.Origin == LDAP || usr.Origin == SAML {
						orgConfig.OrgMgrLDAPUsers = append(orgConfig.OrgMgrLDAPUsers, usr.UserName)
					} else {
						orgConfig.OrgMgrUAAUsers = append(orgConfig.OrgMgrUAAUsers, usr.UserName)
					}
				}
			}

			orgBillingMgrs := getCFUsers(im.CloudController, org.MetaData.GUID, organization.ORGS, organization.ROLE_ORG_BILLING_MANAGERS)
			for _, orgBillingMgrUser := range orgBillingMgrs {
				if usr, ok := userIDToUserMap[orgBillingMgrUser]; ok {
					if usr.Origin == LDAP || usr.Origin == SAML {
						orgConfig.OrgBillingMgrLDAPUsers = append(orgConfig.OrgBillingMgrLDAPUsers, usr.UserName)
					} else {
						orgConfig.OrgBillingMgrUAAUsers = append(orgConfig.OrgBillingMgrUAAUsers, usr.UserName)
					}
				}
			}

			orgAuditors := getCFUsers(im.CloudController, org.MetaData.GUID, organization.ORGS, organization.ROLE_ORG_AUDITORS)
			for _, orgAuditorUser := range orgAuditors {
				if usr, ok := userIDToUserMap[orgAuditorUser]; ok {
					if usr.Origin == LDAP || usr.Origin == SAML {
						orgConfig.OrgAuditorLDAPUsers = append(orgConfig.OrgAuditorLDAPUsers, usr.UserName)
					} else {
						orgConfig.OrgAuditorUAAUsers = append(orgConfig.OrgAuditorUAAUsers, usr.UserName)
					}
				}
			}
			configMgr.AddOrgToConfig(orgConfig)
		} else {
			fmt.Fprintf(os.Stdout, "Skipping org : %s as it is ignored from import", org.Entity.Name)
		}
	}
	return nil
}

func getCFUsers(cc cloudcontroller.Manager, entityGUID, entityType, role string) []string {
	userIDMap, err := cc.GetCFUsers(entityGUID, entityType, role)
	if err != nil && len(userIDMap) > 0 {
		orgUsers := make([]string, len(userIDMap))
		for userID := range userIDMap {
			orgUsers = append(orgUsers, userID)
		}
		return orgUsers
	}
	return make([]string, 0)
}
