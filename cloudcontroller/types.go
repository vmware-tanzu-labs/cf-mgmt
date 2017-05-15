package cloudcontroller

import "github.com/pivotalservices/cf-mgmt/http"

type Manager interface {
	CreateSpace(spaceName, orgGUID string) error
	ListSpaces(orgGUID string) ([]*Space, error)
	AddUserToSpaceRole(userName, role, spaceGUID string) error
	UpdateSpaceSSH(sshAllowed bool, spaceGUID string) error

	AssignSecurityGroupToSpace(spaceGUID, sgGUID string) error
	ListSecurityGroups() (map[string]string, error)
	CreateSecurityGroup(sgName, contents string) (string, error)
	UpdateSecurityGroup(sgGUID, sgName, contents string) error

	CreateSpaceQuota(orgGUID, quotaName string,
		memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
		paidServicePlansAllowed bool) (string, error)
	UpdateSpaceQuota(orgGUID, quotaGUID, quotaName string,
		memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
		paidServicePlansAllowed bool) error
	ListAllSpaceQuotasForOrg(orgGUID string) (map[string]string, error)
	AssignQuotaToSpace(spaceGUID, quotaGUID string) error

	CreateOrg(orgName string) error
	DeleteOrg(orgName string) error
	ListOrgs() ([]*Org, error)
	AddUserToOrgRole(userName, role, orgGUID string) error
	AddUserToOrg(userName, orgGUID string) error

	ListAllOrgQuotas() (quotas map[string]string, err error)
	CreateQuota(quotaName string,
		memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
		paidServicePlansAllowed bool) (string, error)
	UpdateQuota(quotaGUID, quotaName string,
		memoryLimit, instanceMemoryLimit, totalRoutes, totalServices int,
		paidServicePlansAllowed bool) error

	AssignQuotaToOrg(orgGUID, quotaGUID string) error

	GetCFUsers(entityGUID, entityType, role string) (map[string]string, error)

	RemoveCFUser(entityGUID, entityType, userGUID, role string) error
	//Returns a specific quota definition for either an org or space
	QuotaDef(quotaDefGUID string, entityType string) (*Quota, error)
}

type DefaultManager struct {
	Host  string
	Token string
	HTTP  http.Manager
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

//SecurityGroupMetaData -
type SecurityGroupMetaData struct {
	GUID string `json:"guid"`
}

//SecurityGroupEntity -
type SecurityGroupEntity struct {
	Name string `json:"name"`
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
	PaidServicePlansAllowed bool   `json:"paid_service_plans_allowed"`
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
