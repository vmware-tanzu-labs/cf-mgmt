package shareddomain

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"strings"

	"code.cloudfoundry.org/routing-api/models"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type Manager struct {
	DomainClient  CFDomainClient
	JobClient     CFJobClient
	RoutingClient RoutingClient
	Cfg           config.Reader
	Peek          bool
}

//go:generate counterfeiter -o fakes/fake_domain_client.go . CFDomainClient
type CFDomainClient interface {
	ListAll(ctx context.Context, opts *cfclient.DomainListOptions) ([]*resource.Domain, error)
	Create(ctx context.Context, r *resource.DomainCreate) (*resource.Domain, error)
	Delete(ctx context.Context, guid string) (string, error)
}

//go:generate counterfeiter -o fakes/fake_job_client.go . CFJobClient
type CFJobClient interface {
	PollComplete(ctx context.Context, jobGUID string, opts *cfclient.PollingOptions) error
}

//go:generate counterfeiter -o fakes/fake_routing_client.go . RoutingClient
type RoutingClient interface {
	RouterGroupWithName(string) (models.RouterGroup, error)
	RouterGroups() ([]models.RouterGroup, error)
}

func NewManager(cfclient CFDomainClient, jobClient CFJobClient, routingClient RoutingClient, cfg config.Reader, peek bool) *Manager {
	return &Manager{
		DomainClient:  cfclient,
		JobClient:     jobClient,
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

	currentDomains, err := m.DomainClient.ListAll(context.Background(), nil)
	if err != nil {
		return err
	}
	domainMap := make(map[string]string)
	for _, domain := range currentDomains {
		domainMap[strings.ToLower(domain.Name)] = domain.GUID
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
			r := resource.NewDomainCreate(expectedDomain)
			r.Internal = &sharedDomainConfig.Internal
			r.RouterGroup = &resource.Relationship{
				GUID: routerGroupGUID,
			}
			_, err := m.DomainClient.Create(context.Background(), r)
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
			jobGUID, err := m.DomainClient.Delete(context.Background(), domainGUID)
			if err != nil {
				return err
			}
			err = m.JobClient.PollComplete(context.Background(), jobGUID, nil)
			if err != nil {
				return err
			}
		}
	} else {
		lo.G.Infof("Shared Domains will not be removed, must set enable-remove-shared-domains: true in cf-mgmt.yml")
	}
	return nil
}
