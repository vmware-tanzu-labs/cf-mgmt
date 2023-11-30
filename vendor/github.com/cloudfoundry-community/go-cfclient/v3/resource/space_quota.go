package resource

import "time"

type SpaceQuota struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 	Name of the quota
	Name string `json:"name"`

	// quotas that affect apps
	Apps SpaceQuotaApps `json:"apps"`

	// quotas that affect services
	Services SpaceQuotaServices `json:"services"`

	// quotas that affect routes
	Routes SpaceQuotaRoutes `json:"routes"`

	// relationships to the organization and spaces where the quota belongs
	Relationships SpaceQuotaRelationships `json:"relationships"`

	Links map[string]Link `json:"links"`
}

type SpaceQuotaList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*SpaceQuota `json:"resources"`
}

type SpaceQuotaCreateOrUpdate struct {
	// 	Name of the quota
	Name *string `json:"name,omitempty"`

	// relationships to the organization and spaces where the quota belongs
	Relationships *SpaceQuotaRelationships `json:"relationships,omitempty"`

	// quotas that affect apps
	Apps *SpaceQuotaApps `json:"apps,omitempty"`

	// quotas that affect services
	Services *SpaceQuotaServices `json:"services,omitempty"`

	// quotas that affect routes
	Routes *SpaceQuotaRoutes `json:"routes,omitempty"`
}

type SpaceQuotaApps struct {
	// Total memory allowed for all the started processes and running tasks in a space
	TotalMemoryInMB *int `json:"total_memory_in_mb"`

	// Maximum memory for a single process or task
	PerProcessMemoryInMB *int `json:"per_process_memory_in_mb"`

	// Total log rate limit allowed for all the started processes and running tasks in an organization
	LogRateLimitInBytesPerSecond *int `json:"log_rate_limit_in_bytes_per_second"`

	// Total instances of all the started processes allowed in a space
	TotalInstances *int `json:"total_instances"`

	// Maximum number of running tasks in a space
	PerAppTasks *int `json:"per_app_tasks"`
}

type SpaceQuotaServices struct {
	// Specifies whether instances of paid service plans can be created
	PaidServicesAllowed *bool `json:"paid_services_allowed,omitempty"`

	// Total number of service instances allowed in a space
	TotalServiceInstances *int `json:"total_service_instances"`

	// Total number of service keys allowed in a space
	TotalServiceKeys *int `json:"total_service_keys"`
}

type SpaceQuotaRoutes struct {
	// Total number of routes allowed in a space
	TotalRoutes *int `json:"total_routes"`

	// Total number of ports that are reservable by routes in a space
	TotalReservedPorts *int `json:"total_reserved_ports"`
}

type SpaceQuotaRelationships struct {
	Organization *ToOneRelationship   `json:"organization,omitempty"`
	Spaces       *ToManyRelationships `json:"spaces,omitempty"`
}

func NewSpaceQuotaCreate(name, orgGUID string) *SpaceQuotaCreateOrUpdate {
	return &SpaceQuotaCreateOrUpdate{
		Name: &name,
		Relationships: &SpaceQuotaRelationships{
			Organization: &ToOneRelationship{
				Data: &Relationship{
					GUID: orgGUID,
				},
			},
		},
	}
}

func NewSpaceQuotaUpdate() *SpaceQuotaCreateOrUpdate {
	return &SpaceQuotaCreateOrUpdate{}
}

func (s *SpaceQuotaCreateOrUpdate) WithName(name string) *SpaceQuotaCreateOrUpdate {
	s.Name = &name
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalMemoryInMB(mb int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &SpaceQuotaApps{}
	}
	s.Apps.TotalMemoryInMB = &mb
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithPerProcessMemoryInMB(mb int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &SpaceQuotaApps{}
	}
	s.Apps.PerProcessMemoryInMB = &mb
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithLogRateLimitInBytesPerSecond(mbps int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &SpaceQuotaApps{}
	}
	s.Apps.LogRateLimitInBytesPerSecond = &mbps
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalInstances(count int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &SpaceQuotaApps{}
	}
	s.Apps.TotalInstances = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithPerAppTasks(count int) *SpaceQuotaCreateOrUpdate {
	if s.Apps == nil {
		s.Apps = &SpaceQuotaApps{}
	}
	s.Apps.PerAppTasks = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithPaidServicesAllowed(allowed bool) *SpaceQuotaCreateOrUpdate {
	if s.Services == nil {
		s.Services = &SpaceQuotaServices{}
	}
	s.Services.PaidServicesAllowed = &allowed
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalServiceInstances(count int) *SpaceQuotaCreateOrUpdate {
	if s.Services == nil {
		s.Services = &SpaceQuotaServices{}
	}
	s.Services.TotalServiceInstances = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalServiceKeys(count int) *SpaceQuotaCreateOrUpdate {
	if s.Services == nil {
		s.Services = &SpaceQuotaServices{}
	}
	s.Services.TotalServiceKeys = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalRoutes(count int) *SpaceQuotaCreateOrUpdate {
	if s.Routes == nil {
		s.Routes = &SpaceQuotaRoutes{}
	}
	s.Routes.TotalRoutes = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithTotalReservedPorts(count int) *SpaceQuotaCreateOrUpdate {
	if s.Routes == nil {
		s.Routes = &SpaceQuotaRoutes{}
	}
	s.Routes.TotalReservedPorts = &count
	return s
}

func (s *SpaceQuotaCreateOrUpdate) WithSpaces(spaceGUIDs ...string) *SpaceQuotaCreateOrUpdate {
	if s.Relationships == nil {
		s.Relationships = &SpaceQuotaRelationships{
			Spaces: &ToManyRelationships{},
		}
	}
	for _, g := range spaceGUIDs {
		r := Relationship{
			GUID: g,
		}
		s.Relationships.Spaces.Data = append(s.Relationships.Spaces.Data, r)
	}
	return s
}
