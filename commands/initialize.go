package commands

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	routing_api "code.cloudfoundry.org/routing-api"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	v3cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	v3config "github.com/cloudfoundry-community/go-cfclient/v3/config"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/configcommands"
	"github.com/vmwarepivotallabs/cf-mgmt/isosegment"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/organization"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/privatedomain"
	"github.com/vmwarepivotallabs/cf-mgmt/quota"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/securitygroup"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"
	"github.com/vmwarepivotallabs/cf-mgmt/shareddomain"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
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
	RoleManager             role.Manager
}

type Initialize struct {
	ConfigDir, SystemDomain, UserID, Password, ClientSecret, LdapPwd string
	Peek                                                             bool
}

func InitializeManagers(baseCommand BaseCFConfigCommand) (*CFMgmt, error) {
	return InitializePeekManagers(baseCommand, false, nil)
}

func InitializeLdapManager(baseCommand BaseCFConfigCommand, ldapCommand BaseLDAPCommand) (*ldap.Manager, error) {
	cfg := config.NewManager(baseCommand.ConfigDirectory)
	ldapConfig, err := cfg.LdapConfig(ldapCommand.LdapUser, ldapCommand.LdapPassword, ldapCommand.LdapServer)
	if err != nil {
		return nil, err
	}
	if ldapConfig.Enabled {
		ldapMgr, err := ldap.NewManager(ldapConfig)
		if err != nil {
			return nil, err
		}
		return ldapMgr, nil
	}
	return nil, nil
}

func InitializePeekManagers(baseCommand BaseCFConfigCommand, peek bool, ldapMgr *ldap.Manager) (*CFMgmt, error) {
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

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "trace") {
		loggingTranport := util.NewLoggingTransport(httpClient.Transport)
		httpClient = &http.Client{
			Transport: loggingTranport,
		}
	}

	userAgent := fmt.Sprintf("cf-mgmt/%s", configcommands.VERSION)
	uaaMgr, err := uaa.NewDefaultUAAManager(cfMgmt.SystemDomain, baseCommand.UserID, baseCommand.ClientSecret, userAgent, httpClient, peek)
	if err != nil {
		return nil, err
	}
	cfMgmt.UAAManager = uaaMgr

	var c *cfclient.Config
	var cv3 *v3config.Config
	if baseCommand.Password != "" {
		lo.G.Warning("Password parameter is deprecated, create uaa client and client-secret instead")
		c = &cfclient.Config{
			ApiAddress: fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain),
			Username:   baseCommand.UserID,
			Password:   baseCommand.Password,
			UserAgent:  userAgent,
		}
		cv3, err = v3config.NewUserPassword(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain),
			baseCommand.UserID,
			baseCommand.Password)
		if err != nil {
			return nil, err
		}
	} else {
		c = &cfclient.Config{
			ApiAddress:   fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain),
			ClientID:     baseCommand.UserID,
			ClientSecret: baseCommand.ClientSecret,
			UserAgent:    userAgent,
		}
		cv3, err = v3config.NewClientSecret(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain),
			baseCommand.UserID,
			baseCommand.ClientSecret)
		if err != nil {
			return nil, err
		}
	}
	c.HttpClient = httpClient
	cv3.UserAgent = userAgent
	cv3.WithHTTPClient(httpClient)

	client, err := cfclient.NewClient(c)
	if err != nil {
		return nil, err
	}
	v3client, err := v3cfclient.New(cv3)

	if err != nil {
		return nil, err
	}

	cfMgmt.OrgReader = organizationreader.NewReader(client, v3client.Organizations, cfg, peek)
	cfMgmt.SpaceManager = space.NewManager(v3client.Spaces, v3client.SpaceFeatures, cfMgmt.UAAManager, cfMgmt.OrgReader, cfg, peek)
	cfMgmt.OrgManager = organization.NewManager(v3client.Organizations, cfMgmt.OrgReader, cfg, peek)
	cfMgmt.RoleManager = role.New(v3client.Roles, v3client.Users, v3client.Jobs, uaaMgr, peek)

	userManager, err := user.NewManager(cfg, cfMgmt.SpaceManager, cfMgmt.OrgReader, cfMgmt.UAAManager, cfMgmt.RoleManager, ldapMgr, peek)
	if err != nil {
		return nil, err
	}
	cfMgmt.UserManager = userManager
	cfMgmt.SecurityGroupManager = securitygroup.NewManager(v3client.SecurityGroups, cfMgmt.SpaceManager, cfg, peek)
	cfMgmt.QuotaManager = quota.NewManager(v3client.SpaceQuotas, v3client.OrganizationQuotas, cfMgmt.SpaceManager, cfMgmt.OrgReader, cfg, peek)
	cfMgmt.PrivateDomainManager = privatedomain.NewManager(client, cfMgmt.OrgReader, cfg, peek)
	if isoSegmentManager, err := isosegment.NewManager(client, cfg, cfMgmt.OrgReader, cfMgmt.SpaceManager, peek); err == nil {
		cfMgmt.IsolationSegmentManager = isoSegmentManager
	} else {
		return nil, err
	}
	cfMgmt.ServiceAccessManager = serviceaccess.NewManager(client, cfMgmt.OrgReader, cfg, peek)
	token, err := client.GetToken()
	if err != nil {
		return nil, err
	}
	//needs to not include bearer prefix
	token = strings.Replace(token, "bearer ", "", 1)
	routingAPIClient := routing_api.NewClient(c.ApiAddress, true)
	routingAPIClient.SetToken(token)
	cfMgmt.SharedDomainManager = shareddomain.NewManager(client, routingAPIClient, cfg, peek)
	return cfMgmt, nil
}
