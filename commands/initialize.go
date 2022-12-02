package commands

import (
	routingapi "code.cloudfoundry.org/routing-api"
	"context"
	"fmt"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	cfconfig "github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/configcommands"
	"github.com/vmwarepivotallabs/cf-mgmt/isosegment"
	"github.com/vmwarepivotallabs/cf-mgmt/organization"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/privatedomain"
	"github.com/vmwarepivotallabs/cf-mgmt/quota"
	"github.com/vmwarepivotallabs/cf-mgmt/securitygroup"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"
	"github.com/vmwarepivotallabs/cf-mgmt/shareddomain"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/xchapter7x/lo"
)

type CFMgmt struct {
	UAAManager              uaa.Manager
	OrgReader               organizationreader.Reader
	OrgManager              organization.Manager
	SpaceManager            space.Manager
	UserManager             user.Manager
	QuotaManager            *quota.Manager
	PrivateDomainManager    privatedomain.Manager
	ConfigManager           config.Updater
	ConfigDirectory         string
	SystemDomain            string
	SecurityGroupManager    securitygroup.Manager
	IsolationSegmentManager isosegment.Manager
	ServiceAccessManager    *serviceaccess.Manager
	SharedDomainManager     *shareddomain.Manager
}

type Initialize struct {
	ConfigDir, SystemDomain, UserID, Password, ClientSecret, LdapPwd string
	Peek                                                             bool
}

func InitializeManagers(baseCommand BaseCFConfigCommand) (*CFMgmt, error) {
	return InitializePeekManagers(baseCommand, false)
}

func InitializePeekManagers(baseCommand BaseCFConfigCommand, peek bool) (*CFMgmt, error) {
	lo.G.Debugf("Using %s of cf-mgmt", configcommands.GetFormattedVersion())
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

	userAgent := fmt.Sprintf("cf-mgmt/%s", configcommands.VERSION)
	uaaMgr, err := uaa.NewDefaultUAAManager(cfMgmt.SystemDomain, baseCommand.UserID, baseCommand.ClientSecret, userAgent, peek)
	if err != nil {
		return nil, err
	}
	cfMgmt.UAAManager = uaaMgr

	u := fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain)
	var c *cfconfig.Config
	if baseCommand.Password != "" {
		lo.G.Warning("Password parameter is deprecated, create uaa client and client-secret instead")
		c, err = cfconfig.NewUserPassword(u, baseCommand.UserID, baseCommand.Password)
		if err != nil {
			return nil, err
		}
	} else {
		c, err = cfconfig.NewClientSecret(u, baseCommand.UserID, baseCommand.ClientSecret)
		if err != nil {
			return nil, err
		}
	}
	c.WithSkipTLSValidation(true)
	c.UserAgent = userAgent

	// if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
	// 	c.Debug = true
	// }
	client, err := cfclient.New(c)
	if err != nil {
		return nil, err
	}
	cfMgmt.OrgReader = organizationreader.NewReader(client.Organizations, cfg, peek)
	cfMgmt.SpaceManager = space.NewManager(client.Spaces, client.SpaceFeatures, client.Organizations, client.Jobs, cfMgmt.UAAManager, cfMgmt.OrgReader, cfg, peek)
	cfMgmt.OrgManager = organization.NewManager(client.Organizations, client.Domains, client.Jobs, cfMgmt.OrgReader, cfMgmt.SpaceManager, cfg, peek)
	userManager, err := user.NewManager(client.Roles, client.Users, client.Spaces, client.Jobs, cfg, cfMgmt.SpaceManager, cfMgmt.OrgReader, cfMgmt.UAAManager, peek)
	if err != nil {
		return nil, err
	}
	cfMgmt.UserManager = userManager
	cfMgmt.SecurityGroupManager = securitygroup.NewManager(client.SecurityGroups, cfMgmt.SpaceManager, cfg, peek)
	cfMgmt.QuotaManager = quota.NewManager(client.SpaceQuotas, client.OrganizationQuotas, cfMgmt.SpaceManager, cfMgmt.OrgReader, cfMgmt.OrgManager, cfg, peek)
	cfMgmt.PrivateDomainManager = privatedomain.NewManager(client.Domains, client.Jobs, cfMgmt.OrgReader, cfg, peek)
	if isoSegmentManager, err := isosegment.NewManager(client.IsolationSegments, client.Organizations, client.Spaces, cfg, cfMgmt.OrgReader, cfMgmt.SpaceManager, peek); err == nil {
		cfMgmt.IsolationSegmentManager = isoSegmentManager
	} else {
		return nil, err
	}
	cfMgmt.ServiceAccessManager = serviceaccess.NewManager(client.ServicePlans, client.ServicePlansVisibility, client.ServiceOfferings, client.ServiceBrokers, cfMgmt.OrgReader, cfg, peek)

	token, err := client.AccessToken(context.Background())
	if err != nil {
		return nil, err
	}
	routingAPIClient := routingapi.NewClient(c.APIEndpointURL, true)
	routingAPIClient.SetToken(token)
	cfMgmt.SharedDomainManager = shareddomain.NewManager(client.Domains, client.Jobs, routingAPIClient, cfg, peek)

	return cfMgmt, nil
}
