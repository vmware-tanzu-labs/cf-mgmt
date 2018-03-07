package cloudcontroller

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type Manager interface {
	ListIsolationSegments() ([]cfclient.IsolationSegment, error)
	OrgQuotaByName(name string) (cfclient.OrgQuota, error)
	SpaceQuotaByName(name string) (cfclient.SpaceQuota, error)
}

type DefaultManager struct {
	Client cfclient.Client
	Peek   bool
}

//QuotaEntity -
type QuotaEntity struct {
	Name                    string `json:"name"`
	MemoryLimit             int    `json:"memory_limit"`
	InstanceMemoryLimit     int    `json:"instance_memory_limit"`
	TotalRoutes             int    `json:"total_routes"`
	TotalServices           int    `json:"total_services"`
	PaidServicePlansAllowed bool   `json:"non_basic_services_allowed"`
	TotalPrivateDomains     int    `json:"total_private_domains"`
	TotalReservedRoutePorts int    `json:"total_reserved_route_ports"`
	TotalServiceKeys        int    `json:"total_service_keys"`
	AppInstanceLimit        int    `json:"app_instance_limit"`
}

//SpaceQuotaEntity -
type SpaceQuotaEntity struct {
	QuotaEntity
	OrgGUID string `json:"organization_guid"`
}

//GetName --
func (qe *QuotaEntity) GetName() string {
	return qe.Name
}

//IsQuotaEnabled --
func (qe *QuotaEntity) IsQuotaEnabled() bool {
	return qe.Name != ""
}

//GetMemoryLimit --
func (qe *QuotaEntity) GetMemoryLimit() int {
	if qe.MemoryLimit == 0 {
		return 10240
	}
	return qe.MemoryLimit
}

//GetInstanceMemoryLimit --
func (qe *QuotaEntity) GetInstanceMemoryLimit() int {
	if qe.InstanceMemoryLimit == 0 {
		return -1
	}
	return qe.InstanceMemoryLimit
}

//GetTotalServices --
func (qe *QuotaEntity) GetTotalServices() int {
	if qe.TotalServices == 0 {
		return -1
	}
	return qe.TotalServices
}

//GetTotalRoutes --
func (qe *QuotaEntity) GetTotalRoutes() int {
	if qe.TotalRoutes == 0 {
		return 1000
	}
	return qe.TotalRoutes
}

//IsPaidServicesAllowed  --
func (qe *QuotaEntity) IsPaidServicesAllowed() bool {
	return qe.PaidServicePlansAllowed
}
