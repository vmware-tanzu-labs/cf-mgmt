package shareddomain

import (
	"strings"

	"code.cloudfoundry.org/routing-api/models"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type Manager struct {
	CFClient      CFClient
	RoutingClient RoutingClient
	Cfg           config.Reader
	Peek          bool
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_cf_client.go . CFClient
type CFClient interface {
	ListSharedDomains() ([]cfclient.SharedDomain, error)
	DeleteSharedDomain(guid string, async bool) error
	CreateSharedDomain(name string, internal bool, router_group_guid string) (*cfclient.SharedDomain, error)
}

//counterfeiter:generate -o fakes/fake_routing_client.go . RoutingClient
type RoutingClient interface {
	RouterGroupWithName(string) (models.RouterGroup, error)
	RouterGroups() ([]models.RouterGroup, error)
}

func NewManager(cfclient CFClient, routingClient RoutingClient, cfg config.Reader, peek bool) *Manager {
	return &Manager{
		CFClient:      cfclient,
		RoutingClient: routingClient,
		Cfg:           cfg,
		Peek:          peek,
	}
}

func (m *Manager) Apply() error {
	global, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}

	currentDomains, err := m.CFClient.ListSharedDomains()
	if err != nil {
		return err
	}
	domainMap := make(map[string]string)
	for _, domain := range currentDomains {
		domainMap[strings.ToLower(domain.Name)] = domain.Guid
	}
	for expectedDomain, sharedDomainConfig := range global.SharedDomains {
		if _, ok := domainMap[strings.ToLower(expectedDomain)]; !ok {
			if m.Peek {
				lo.G.Infof("[dry-run]: create shared domain %s as internal [%t] for router group [%s]", expectedDomain, sharedDomainConfig.Internal, sharedDomainConfig.RouterGroup)
				continue
			}

			routerGroupGUID := ""
			if sharedDomainConfig.RouterGroup != "" {
				routingGroup, err := m.RoutingClient.RouterGroupWithName(sharedDomainConfig.RouterGroup)
				if err != nil {
					return err
				}
				routerGroupGUID = routingGroup.Guid
			}
			lo.G.Infof("create shared domain %s as internal [%t] for router group [%s]", expectedDomain, sharedDomainConfig.Internal, sharedDomainConfig.RouterGroup)
			_, err := m.CFClient.CreateSharedDomain(expectedDomain, sharedDomainConfig.Internal, routerGroupGUID)
			if err != nil {
				return err
			}
		} else {
			delete(domainMap, strings.ToLower(expectedDomain))
		}
	}
	if global.EnableDeleteSharedDomains {
		for domain, domainGUID := range domainMap {
			if m.Peek {
				lo.G.Infof("[dry-run]: deleting shared domain %s", domain)
				continue
			}
			lo.G.Infof("deleting shared domain %s", domain)
			err := m.CFClient.DeleteSharedDomain(domainGUID, false)
			if err != nil {
				return err
			}
		}
	} else {
		lo.G.Infof("Shared Domains will not be removed, must set enable-remove-shared-domains: true in cf-mgmt.yml")
	}
	return nil
}
