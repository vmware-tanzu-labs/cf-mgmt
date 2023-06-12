package commands

import (
	"fmt"
)

type ApplyCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
	BaseLDAPCommand
}

// Execute - applies all the config in order
func (c *ApplyCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err != nil {
		return err
	}
	if err = cfMgmt.UserManager.InitializeLdap(c.LdapUser, c.LdapPassword, c.LdapServer); err != nil {
		return err
	}
	defer cfMgmt.UserManager.DeinitializeLdap()
	fmt.Println("*********  Creating Orgs")
	if err = cfMgmt.OrgManager.CreateOrgs(); err != nil {
		return err
	}

	fmt.Println("*********  Update Orgs Metadata")
	if err = cfMgmt.OrgManager.UpdateOrgsMetadata(); err != nil {
		return err
	}

	fmt.Println("*********  Delete Orgs")
	if err = cfMgmt.OrgManager.DeleteOrgs(); err != nil {
		return err
	}

	fmt.Println("*********  Update Org Users")
	if err = cfMgmt.UserManager.UpdateOrgUsers(); err != nil {
		return err
	}

	fmt.Println("*********  Create Global Security Groups")
	if err = cfMgmt.SecurityGroupManager.CreateGlobalSecurityGroups(); err != nil {
		return err
	}

	fmt.Println("*********  Assign Default Security Groups")
	if err = cfMgmt.SecurityGroupManager.AssignDefaultSecurityGroups(); err != nil {
		return err
	}

	fmt.Println("*********  Create Private Domains")
	if err = cfMgmt.PrivateDomainManager.CreatePrivateDomains(); err != nil {
		return err
	}

	fmt.Println("*********  Share Private Domains")
	if err = cfMgmt.PrivateDomainManager.SharePrivateDomains(); err != nil {
		return err
	}

	fmt.Println("*********  Create Org Quotas")
	if err = cfMgmt.QuotaManager.CreateOrgQuotas(); err != nil {
		return err
	}

	fmt.Println("*********  Create Spaces")
	if err = cfMgmt.SpaceManager.CreateSpaces(); err != nil {
		return err
	}

	fmt.Println("*********  Delete Spaces")
	if err = cfMgmt.SpaceManager.DeleteSpaces(); err != nil {
		return err
	}

	fmt.Println("*********  Update Spaces")
	if err = cfMgmt.SpaceManager.UpdateSpaces(); err != nil {
		return err
	}

	fmt.Println("*********  Update Spaces Metadata")
	if err = cfMgmt.SpaceManager.UpdateSpacesMetadata(); err != nil {
		return err
	}

	fmt.Println("*********  Update Space Users")
	if err = cfMgmt.UserManager.UpdateSpaceUsers(); err != nil {
		return err
	}

	fmt.Println("*********  Create Space Quotas")
	if err = cfMgmt.QuotaManager.CreateSpaceQuotas(); err != nil {
		return err
	}

	fmt.Println("*********  Create Application Security Groups")
	if err = cfMgmt.SecurityGroupManager.CreateApplicationSecurityGroups(); err != nil {
		return err
	}

	fmt.Println("*********  Isolation Segments")
	if err = cfMgmt.IsolationSegmentManager.Apply(); err != nil {
		return err
	}

	fmt.Println("*********  Service Access")
	if err = cfMgmt.ServiceAccessManager.Apply(); err != nil {
		return err
	}

	fmt.Println("*********  Cleanup Org Users")
	if err = cfMgmt.UserManager.CleanupOrgUsers(); err != nil {
		return err
	}

	fmt.Println("*********  Shared Domains")
	if err = cfMgmt.SharedDomainManager.Apply(); err != nil {
		return err
	}

	return nil
}
