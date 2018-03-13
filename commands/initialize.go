package commands

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/configcommands"
	"github.com/pivotalservices/cf-mgmt/isosegment"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/securitygroup"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/spacequota"
	"github.com/pivotalservices/cf-mgmt/spaceusers"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

type CFMgmt struct {
	UAAManager              uaa.Manager
	OrgManager              organization.Manager
	SpaceManager            space.Manager
	SpaceUserManager        spaceusers.Manager
	ConfigManager           config.Updater
	ConfigDirectory         string
	UaacToken               string
	SystemDomain            string
	SecurityGroupManager    securitygroup.Manager
	IsolationSegmentManager isosegment.Manager
	SpaceQuotaManager       spacequota.Manager
}

type Initialize struct {
	ConfigDir, SystemDomain, UserID, Password, ClientSecret, LdapPwd string
	Peek                                                             bool
}

func InitializeManagers(baseCommand BaseCFConfigCommand) (*CFMgmt, error) {
	return InitializePeekManagers(baseCommand, false)
}

func InitializePeekManagers(baseCommand BaseCFConfigCommand, peek bool) (*CFMgmt, error) {
	if baseCommand.SystemDomain == "" ||
		baseCommand.UserID == "" ||
		baseCommand.ClientSecret == "" {
		return nil, fmt.Errorf("must set system-domain, user-id, client-secret properties")
	}

	cfg := config.NewManager(baseCommand.ConfigDirectory)
	var cfToken, uaacToken string
	var err error
	cfMgmt := &CFMgmt{}
	cfMgmt.ConfigDirectory = baseCommand.ConfigDirectory
	cfMgmt.SystemDomain = baseCommand.SystemDomain
	cfMgmt.ConfigManager = config.NewManager(cfMgmt.ConfigDirectory)

	uaaHost := fmt.Sprintf("https://uaa.%s", cfMgmt.SystemDomain)
	if uaacToken, err = uaa.GetUAACToken(uaaHost, baseCommand.UserID, baseCommand.ClientSecret); err != nil {
		return nil, err
	}
	cfMgmt.UaacToken = uaacToken
	cfMgmt.UAAManager = uaa.NewDefaultUAAManager(cfMgmt.SystemDomain, uaacToken, peek)

	if baseCommand.Password != "" {
		lo.G.Warning("Password parameter is deprecated, create uaa client and client-secret instead")
		if cfToken, err = uaa.GetCFToken(uaaHost, baseCommand.UserID, baseCommand.Password); err != nil {
			return nil, err
		}
	} else {
		cfToken = uaacToken
	}

	c := &cfclient.Config{
		ApiAddress:        fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain),
		SkipSslValidation: true,
		Token:             cfToken,
		UserAgent:         fmt.Sprintf("cf-mgmt/%s", configcommands.VERSION),
	}

	client, err := cfclient.NewClient(c)
	if err != nil {
		return nil, err
	}
	ldapMgr := ldap.NewManager()
	cfMgmt.OrgManager = organization.NewManager(client, cfMgmt.UAAManager, cfg, peek)
	cfMgmt.SpaceManager = space.NewManager(client, cfMgmt.UAAManager, cfMgmt.OrgManager, cfg, peek)
	cfMgmt.SpaceUserManager = spaceusers.NewManager(client, cfg, cfMgmt.SpaceManager, ldapMgr, cfMgmt.UAAManager, peek)
	cfMgmt.SecurityGroupManager = securitygroup.NewManager(client, cfMgmt.SpaceManager, cfg, peek)
	cfMgmt.SpaceQuotaManager = spacequota.NewManager(client, cfMgmt.SpaceManager, cfg, peek)
	if isoSegmentManager, err := isosegment.NewManager(client, cfg, peek); err == nil {
		cfMgmt.IsolationSegmentManager = isoSegmentManager
	} else {
		return nil, err
	}
	return cfMgmt, nil
}
