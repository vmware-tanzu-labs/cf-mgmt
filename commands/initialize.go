package commands

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/configcommands"
	"github.com/pivotalservices/cf-mgmt/isosegment"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/securitygroup"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

type CFMgmt struct {
	UAAManager              uaa.Manager
	OrgManager              organization.Manager
	SpaceManager            space.Manager
	ConfigManager           config.Updater
	ConfigDirectory         string
	UaacToken               string
	SystemDomain            string
	CloudController         cloudcontroller.Manager
	SecurityGroupManager    securitygroup.Manager
	IsolationSegmentUpdater *isosegment.Updater
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

	if uaacToken, err = uaa.GetUAACToken(cfMgmt.SystemDomain, baseCommand.UserID, baseCommand.ClientSecret); err != nil {
		return nil, err
	}
	cfMgmt.UaacToken = uaacToken
	cfMgmt.UAAManager = uaa.NewDefaultUAAManager(cfMgmt.SystemDomain, uaacToken, peek)

	if baseCommand.Password != "" {
		lo.G.Warning("Password parameter is deprecated, create uaa client and client-secret instead")
		if cfToken, err = uaa.GetCFToken(cfMgmt.SystemDomain, baseCommand.UserID, baseCommand.Password); err != nil {
			return nil, err
		}
		cfMgmt.CloudController = cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain), cfToken, peek)
	} else {
		cfToken = uaacToken
		cfMgmt.CloudController = cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain), uaacToken, peek)
	}
	cfMgmt.OrgManager = organization.NewManager(cfMgmt.CloudController, cfMgmt.UAAManager, cfg)
	cfMgmt.SpaceManager = space.NewManager(cfMgmt.CloudController, cfMgmt.UAAManager, cfMgmt.OrgManager, cfg)
	cfMgmt.SecurityGroupManager = securitygroup.NewManager(cfMgmt.CloudController, cfg)
	if isoSegmentUpdater, err := isosegment.NewUpdater(configcommands.VERSION, cfMgmt.SystemDomain, cfToken, baseCommand.UserID, baseCommand.ClientSecret, cfg); err == nil {
		cfMgmt.IsolationSegmentUpdater = isoSegmentUpdater
	} else {
		return nil, err
	}
	return cfMgmt, nil
}
