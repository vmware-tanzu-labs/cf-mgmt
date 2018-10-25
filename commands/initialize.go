package commands

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/configcommands"
	"github.com/pivotalservices/cf-mgmt/isosegment"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/privatedomain"
	"github.com/pivotalservices/cf-mgmt/quota"
	"github.com/pivotalservices/cf-mgmt/securitygroup"
	"github.com/pivotalservices/cf-mgmt/serviceaccess"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pivotalservices/cf-mgmt/user"
	"github.com/xchapter7x/lo"
)

type CFMgmt struct {
	UAAManager              uaa.Manager
	OrgManager              organization.Manager
	SpaceManager            space.Manager
	UserManager             user.Manager
	QuotaManager            quota.Manager
	PrivateDomainManager    privatedomain.Manager
	ConfigManager           config.Updater
	ConfigDirectory         string
	SystemDomain            string
	SecurityGroupManager    securitygroup.Manager
	IsolationSegmentManager isosegment.Manager
	ServiceAccessManager    *serviceaccess.Manager
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
	var err error
	cfMgmt := &CFMgmt{}
	cfMgmt.ConfigDirectory = baseCommand.ConfigDirectory
	cfMgmt.SystemDomain = baseCommand.SystemDomain
	cfMgmt.ConfigManager = config.NewManager(cfMgmt.ConfigDirectory)

	uaaMgr, err := uaa.NewDefaultUAAManager(cfMgmt.SystemDomain, baseCommand.UserID, baseCommand.ClientSecret, peek)
	if err != nil {
		return nil, err
	}
	cfMgmt.UAAManager = uaaMgr

	var c *cfclient.Config
	if baseCommand.Password != "" {
		lo.G.Warning("Password parameter is deprecated, create uaa client and client-secret instead")
		c = &cfclient.Config{
			ApiAddress:        fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain),
			SkipSslValidation: true,
			Username:          baseCommand.UserID,
			Password:          baseCommand.Password,
			UserAgent:         fmt.Sprintf("cf-mgmt/%s", configcommands.VERSION),
		}
	} else {
		c = &cfclient.Config{
			ApiAddress:        fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain),
			SkipSslValidation: true,
			ClientID:          baseCommand.UserID,
			ClientSecret:      baseCommand.ClientSecret,
			UserAgent:         fmt.Sprintf("cf-mgmt/%s", configcommands.VERSION),
		}
	}
	client, err := cfclient.NewClient(c)
	if err != nil {
		return nil, err
	}
	cfMgmt.OrgManager = organization.NewManager(client, cfg, peek)
	cfMgmt.SpaceManager = space.NewManager(client, cfMgmt.UAAManager, cfMgmt.OrgManager, cfg, peek)
	cfMgmt.UserManager = user.NewManager(client, cfg, cfMgmt.SpaceManager, cfMgmt.OrgManager, cfMgmt.UAAManager, peek)
	cfMgmt.SecurityGroupManager = securitygroup.NewManager(client, cfMgmt.SpaceManager, cfg, peek)
	cfMgmt.QuotaManager = quota.NewManager(client, cfMgmt.SpaceManager, cfMgmt.OrgManager, cfg, peek)
	cfMgmt.PrivateDomainManager = privatedomain.NewManager(client, cfMgmt.OrgManager, cfg, peek)
	if isoSegmentManager, err := isosegment.NewManager(client, cfg, cfMgmt.OrgManager, cfMgmt.SpaceManager, peek); err == nil {
		cfMgmt.IsolationSegmentManager = isoSegmentManager
	} else {
		return nil, err
	}
	cfMgmt.ServiceAccessManager = serviceaccess.NewManager(client, cfMgmt.OrgManager, cfg, peek)
	return cfMgmt, nil
}
