package commands

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/isosegment"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/securitygroup"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pivotalservices/cf-mgmt/uaac"
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
	UAACManager             uaac.Manager
	CloudController         cloudcontroller.Manager
	SecurityGroupManager    securitygroup.Manager
	IsolationSegmentUpdater *isosegment.Updater
}

type Initialize struct {
	ConfigDir, SystemDomain, UserID, Password, ClientSecret, LdapPwd string
	Peek                                                             bool
}

func InitializeManagers(baseCommand BaseCFConfigCommand) (*CFMgmt, error) {

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
	cfMgmt.UAAManager = uaa.NewDefaultUAAManager(cfMgmt.SystemDomain, baseCommand.UserID)

	if uaacToken, err = cfMgmt.UAAManager.GetUAACToken(baseCommand.ClientSecret); err != nil {
		return nil, err
	}
	cfMgmt.UaacToken = uaacToken
	cfMgmt.UAACManager = uaac.NewManager(cfMgmt.SystemDomain, uaacToken)

	if baseCommand.Password != "" {
		lo.G.Warning("Password parameter is deprecated, create uaa client and client-secret instead")
		if cfToken, err = cfMgmt.UAAManager.GetCFToken(baseCommand.Password); err != nil {
			return nil, err
		}
		cfMgmt.CloudController = cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain), cfToken)
	} else {
		cfToken = uaacToken
		cfMgmt.CloudController = cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain), uaacToken)
	}
	cfMgmt.OrgManager = organization.NewManager(cfMgmt.SystemDomain, cfToken, uaacToken, cfg)
	cfMgmt.SpaceManager = space.NewManager(cfMgmt.SystemDomain, cfToken, uaacToken, cfg)
	cfMgmt.SecurityGroupManager = securitygroup.NewManager(cfMgmt.SystemDomain, cfToken, cfg)

	cfMgmt.ConfigManager = config.NewManager(cfMgmt.ConfigDirectory)
	if isoSegmentUpdater, err := isosegment.NewUpdater(VERSION, cfMgmt.SystemDomain, cfToken, cfg); err == nil {
		cfMgmt.IsolationSegmentUpdater = isoSegmentUpdater
	} else {
		return nil, err
	}
	return cfMgmt, nil
}
