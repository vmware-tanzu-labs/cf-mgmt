package cloudcontroller

import "github.com/pivotalservices/cf-mgmt/http"

type Manager interface {
	CreateSpace(spaceName, orgGUID string) error
	DeleteSpace(spaceGUID string) error
	ListSpaces(orgGUID string) ([]*Space, error)
	ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error)
	AddUserToSpaceRole(userName, role, spaceGUID string) error
	UpdateSpaceSSH(sshAllowed bool, spaceGUID string) error

	AssignRunningSecurityGroup(sgGUID string) error
	AssignStagingSecurityGroup(sgGUID string) error
	UnassignRunningSecurityGroup(sgGUID string) error
	UnassignStagingSecurityGroup(sgGUID string) error

	AssignSecurityGroupToSpace(spaceGUID, sgGUID string) error
	ListNonDefaultSecurityGroups() (map[string]SecurityGroupInfo, error)
	ListDefaultSecurityGroups() (map[string]SecurityGroupInfo, error)
	ListSecurityGroups() (map[string]SecurityGroupInfo, error)
	CreateSecurityGroup(sgName, contents string) (string, error)
	UpdateSecurityGroup(sgGUID, sgName, contents string) error
	GetSecurityGroupRules(sgGUID string) ([]byte, error)

	CreateSpaceQuota(quota SpaceQuotaEntity) (string, error)
	UpdateSpaceQuota(quotaGUID string, quota SpaceQuotaEntity) error
	ListAllSpaceQuotasForOrg(orgGUID string) (map[string]string, error)
	AssignQuotaToSpace(spaceGUID, quotaGUID string) error

	CreateOrg(orgName string) error
	DeleteOrg(orgGUID string) error
	DeleteOrgByName(orgName string) error
	ListOrgs() ([]*Org, error)
	AddUserToOrgRole(userName, role, orgGUID string) error
	AddUserToOrg(userName, orgGUID string) error

	ListAllOrgQuotas() (quotas map[string]string, err error)
	CreateQuota(quota QuotaEntity) (string, error)
	UpdateQuota(quotaGUID string, quota QuotaEntity) error

	AssignQuotaToOrg(orgGUID, quotaGUID string) error

	GetCFUsers(entityGUID, entityType, role string) (map[string]string, error)

	RemoveCFUser(entityGUID, entityType, userGUID, role string) error
	//Returns a specific quota definition for either an org or space
	QuotaDef(quotaDefGUID string, entityType string) (*Quota, error)

	ListAllPrivateDomains() (map[string]PrivateDomainInfo, error)
	ListOrgOwnedPrivateDomains(orgGUID string) (map[string]string, error)
	ListOrgSharedPrivateDomains(orgGUID string) (map[string]string, error)
	DeletePrivateDomain(guid string) error
	CreatePrivateDomain(orgGUID, privateDomain string) (string, error)
	SharePrivateDomain(sharedOrgGUID, privateDomainGUID string) error
	RemoveSharedPrivateDomain(sharedOrgGUID, privateDomainGUID string) error
}

type DefaultManager struct {
	Host  string
	Token string
	HTTP  http.Manager
	Peek  bool
}

//PrivateDomainResources -
type PrivateDomainResources struct {
	PrivateDomains []*PrivateDomain `json:"resources"`
	NextURL        string           `json:"next_url"`
}

//PrivateDomain -
type PrivateDomain struct {
	MetaData PrivateDomainMetaData `json:"metadata"`
	Entity   PrivateDomainEntity   `json:"entity"`
}

type PrivateDomainInfo struct {
	OrgGUID           string
	PrivateDomainGUID string
}

//PrivateDomainMetaData -
type PrivateDomainMetaData struct {
	GUID string `json:"guid"`
}

//PrivateDomainEntity -
type PrivateDomainEntity struct {
	Name    string `json:"name"`
	OrgGUID string `json:"owning_organization_guid"`
}

type Pagination interface {
	GetNextURL() string
	AddInstances(Pagination)
}

//SpaceResources -
type SpaceResources struct {
	Spaces  []*Space `json:"resources"`
	NextURL string   `json:"next_url"`
}

type Space struct {
	MetaData SpaceMetaData `json:"metadata"`
	Entity   SpaceEntity   `json:"entity"`
}

//SpaceMetaData -
type SpaceMetaData struct {
	GUID string `json:"guid"`
}

//SpaceEntity -
type SpaceEntity struct {
	Name                string `json:"name"`
	AllowSSH            bool   `json:"allow_ssh"`
	OrgGUID             string `json:"organization_guid"`
	QuotaDefinitionGUID string `json:"space_quota_definition_guid"`
}

//Orgs -
type Orgs struct {
	NextURL string `json:"next_url"`
	Orgs    []*Org `json:"resources"`
}

//Org -
type Org struct {
	Entity   OrgEntity   `json:"entity"`
	MetaData OrgMetaData `json:"metadata"`
}

//OrgEntity -
type OrgEntity struct {
	Name                string `json:"name"`
	QuotaDefinitionGUID string `json:"quota_definition_guid"`
}

//OrgMetaData -
type OrgMetaData struct {
	GUID string `json:"guid"`
}

//SecurityGroupResources -
type SecurityGroupResources struct {
	NextURL        string          `json:"next_url"`
	SecurityGroups []SecurityGroup `json:"resources"`
}

type SecurityGroup struct {
	MetaData SecurityGroupMetaData `json:"metadata"`
	Entity   SecurityGroupEntity   `json:"entity"`
}

type SecurityGroupRule struct {
	Entity SecurityGroupRuleEntity `json:"entity"`
}

//SecurityGroupMetaData -
type SecurityGroupMetaData struct {
	GUID string `json:"guid"`
}

//SecurityGroupEntity -
type SecurityGroupEntity struct {
	Name           string      `json:"name"`
	Rules          interface{} `json:"rules"`
	DefaultStaging bool        `json:"staging_default"`
	DefaultRunning bool        `json:"running_default"`
}

//SecurityGroupInfo -
type SecurityGroupInfo struct {
	GUID           string
	Rules          string
	DefaultStaging bool
	DefaultRunning bool
}

//SecurityGroupRuleEntity -
type SecurityGroupRuleEntity struct {
	Name  string `json:"name"`
	Rules []Rule `json:"rules"`
}

//Rule -
type Rule struct {
	Protocol    string `json:"protocol,omitempty"`
	Ports       string `json:"ports,omitempty"`
	Destination string `json:"destination,omitempty"`
	Type        int    `json:"type,omitempty"`
	Code        int    `json:"code,omitempty"`
	Log         bool   `json:"log,omitempty"`
	Description string `json:"description,omitempty"`
}

//Quotas -
type Quotas struct {
	NextURL string  `json:"next_url"`
	Quotas  []Quota `json:"resources"`
}

type Quota struct {
	MetaData QuotaMetaData `json:"metadata"`
	Entity   QuotaEntity   `json:"entity"`
}

//QuotaMetaData -
type QuotaMetaData struct {
	GUID string `json:"guid"`
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

//OrgSpaceUsers -
type OrgSpaceUsers struct {
	NextURL string          `json:"next_url"`
	Users   []*OrgSpaceUser `json:"resources"`
}

//OrgSpaceUser -
type OrgSpaceUser struct {
	Entity   UserEntity   `json:"entity"`
	MetaData UserMetaData `json:"metadata"`
}

//UserEntity -
type UserEntity struct {
	UserName string `json:"username"`
}

//UserMetaData -
type UserMetaData struct {
	GUID string `json:"guid"`
}
