package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

const defaultAPIResourcePath = "https://api.example.org/v3/somepagedresource"

type JSONResource struct {
	GUID   string
	Name   string
	JSON   string
	Params map[string]string
}

type ResourceResult struct {
	Resource string

	// extra included resources
	// https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#resources-with-includes
	Apps             []string
	Spaces           []string
	Organizations    []string
	Domains          []string
	Users            []string
	ServiceOfferings []string
	ServiceInstances []string
	Routes           []string
}

type PagedResult struct {
	Resources []string

	// extra included resources
	// https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#resources-with-includes
	Apps             []string
	Spaces           []string
	Organizations    []string
	Domains          []string
	Users            []string
	ServiceOfferings []string
	ServiceInstances []string
	Routes           []string
}

type resultTemplate struct {
	TotalResults int
	TotalPages   int
	FirstPage    string
	LastPage     string
	NextPage     string
	PreviousPage string

	Resources        string
	Apps             string
	Spaces           string
	Organizations    string
	Domains          string
	Users            string
	ServiceOfferings string
	ServiceInstances string
	Routes           string
}

type ObjectJSONGenerator struct {
}

func NewObjectJSONGenerator(seed int) *ObjectJSONGenerator {
	rand.Seed(int64(seed)) // stable random
	return &ObjectJSONGenerator{}
}

func (o ObjectJSONGenerator) Application() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "app.json")
}

func (o ObjectJSONGenerator) AppFeature() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "app_feature.json")
}

func (o ObjectJSONGenerator) AppUsage() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "app_usage.json")
}

func (o ObjectJSONGenerator) AuditEvent() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "audit_event.json")
}

func (o ObjectJSONGenerator) AppUpdateEnvVars() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "app_update_envvar.json")
}

func (o ObjectJSONGenerator) AppEnvironment() *JSONResource {
	r := &JSONResource{
		Name: RandomName(),
	}
	return o.renderTemplate(r, "app_environment.json")
}

func (o ObjectJSONGenerator) AppEnvVar() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "app_envvar.json")
}

func (o ObjectJSONGenerator) AppSSH() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "app_ssh.json")
}

func (o ObjectJSONGenerator) AppPermission() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "app_permissions.json")
}

func (o ObjectJSONGenerator) Build(state string) *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Params: map[string]string{
			"state": state,
		},
	}
	return o.renderTemplate(r, "build.json")
}

func (o ObjectJSONGenerator) Buildpack() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "buildpack.json")
}

func (o ObjectJSONGenerator) Droplet() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "droplet.json")
}

func (o ObjectJSONGenerator) DropletAssociation() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "droplet_association.json")
}

func (o ObjectJSONGenerator) Deployment() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "deployment.json")
}

func (o ObjectJSONGenerator) Domain() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "domain.json")
}

func (o ObjectJSONGenerator) DomainShared() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "domain_shared.json")
}

func (o ObjectJSONGenerator) EnvVarGroup() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "environment_variable_group.json")
}

func (o ObjectJSONGenerator) FeatureFlag() *JSONResource {
	r := &JSONResource{
		Name: RandomName(),
	}
	return o.renderTemplate(r, "feature_flag.json")
}

func (o ObjectJSONGenerator) IsolationSegment() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "isolation_segment.json")
}

func (o ObjectJSONGenerator) IsolationSegmentRelationships() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "isolation_segment_relationships.json")
}

func (o ObjectJSONGenerator) Job(state string) *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Params: map[string]string{
			"state": state,
		},
	}
	return o.renderTemplate(r, "job.json")
}

func (o ObjectJSONGenerator) Manifest() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "manifest.yml")
}

func (o ObjectJSONGenerator) ManifestDiff() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "manifest_diff.yml")
}

func (o ObjectJSONGenerator) Organization() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "org.json")
}

func (o ObjectJSONGenerator) OrganizationUsageSummary() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "org_usage_summary.json")
}

func (o ObjectJSONGenerator) OrganizationQuota() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "org_quota.json")
}

func (o ObjectJSONGenerator) Package(state string) *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Params: map[string]string{
			"state": state,
		},
	}
	return o.renderTemplate(r, "package.json")
}

func (o ObjectJSONGenerator) PackageDocker() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "package_docker.json")
}

func (o ObjectJSONGenerator) Process() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "process.json")
}

func (o ObjectJSONGenerator) ProcessStats() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "process_stats.json")
}

func (o ObjectJSONGenerator) ResourceMatch() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "resource_match.json")
}

func (o ObjectJSONGenerator) Revision() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "revision.json")
}

func (o ObjectJSONGenerator) Role() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "role.json")
}

func (o ObjectJSONGenerator) Route() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "route.json")
}

func (o ObjectJSONGenerator) RouteSpaceRelationships() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "route_space_relationships.json")
}

func (o ObjectJSONGenerator) RouteDestinations() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "route_destinations.json")
}

func (o ObjectJSONGenerator) RouteDestinationWithLinks() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "route_destination_with_links.json")
}

func (o ObjectJSONGenerator) ServiceBroker() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "service_broker.json")
}

func (o ObjectJSONGenerator) SecurityGroup() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "security_group.json")
}

func (o ObjectJSONGenerator) ServiceCredentialBinding() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "service_credential_binding.json")
}

func (o ObjectJSONGenerator) ServiceCredentialBindingDetails() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "service_credential_binding_detail.json")
}

func (o ObjectJSONGenerator) ServiceInstance() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "service_instance.json")
}

func (o ObjectJSONGenerator) ServiceInstanceUsageSummary() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "service_instance_usage_summary.json")
}

func (o ObjectJSONGenerator) ServiceInstanceSpaceRelationships() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "service_instance_space_relationships.json")
}

func (o ObjectJSONGenerator) ServiceOffering() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "service_offering.json")
}

func (o ObjectJSONGenerator) ServicePlan() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "service_plan.json")
}

func (o ObjectJSONGenerator) ServicePlanVisibility() *JSONResource {
	r := &JSONResource{}
	return o.renderTemplate(r, "service_plan_visibility.json")
}

func (o ObjectJSONGenerator) ServiceRouteBinding() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "service_route_binding.json")
}

func (o ObjectJSONGenerator) ServiceUsage() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "service_usage.json")
}

func (o ObjectJSONGenerator) Sidecar() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "sidecar.json")
}

func (o ObjectJSONGenerator) Space() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "space.json")
}

func (o ObjectJSONGenerator) SpaceQuota() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "space_quota.json")
}

func (o ObjectJSONGenerator) Stack() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "stack.json")
}

func (o ObjectJSONGenerator) Task() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
	}
	return o.renderTemplate(r, "task.json")
}

func (o ObjectJSONGenerator) User() *JSONResource {
	r := &JSONResource{
		GUID: RandomGUID(),
		Name: RandomName(),
	}
	return o.renderTemplate(r, "user.json")
}

// ResourceWithInclude merges the included resources under the primary resource's included key
func (o ObjectJSONGenerator) ResourceWithInclude(rr ResourceResult) []string {
	j := map[string]any{}
	err := json.Unmarshal([]byte(rr.Resource), &j)
	if err != nil {
		panic(err)
	}

	t, err := template.New("res").Parse(singleTemplate)
	if err != nil {
		panic(err)
	}
	p := resultTemplate{
		Apps:             strings.Join(rr.Apps, ","),
		Spaces:           strings.Join(rr.Spaces, ","),
		Organizations:    strings.Join(rr.Organizations, ","),
		Domains:          strings.Join(rr.Domains, ","),
		Users:            strings.Join(rr.Users, ","),
		Routes:           strings.Join(rr.Routes, ","),
		ServiceOfferings: strings.Join(rr.ServiceOfferings, ","),
		ServiceInstances: strings.Join(rr.ServiceInstances, ","),
	}

	var h bytes.Buffer
	err = t.Execute(&h, p)
	if err != nil {
		panic(err)
	}
	s := h.String()
	j["included"] = json.RawMessage(s)

	b, err := json.Marshal(&j)
	if err != nil {
		panic(err)
	}
	s = string(b)
	return []string{s}
}

// PagedWithInclude takes the list of resources and inserts them into a paged API response
func (o ObjectJSONGenerator) PagedWithInclude(pagesOfResourcesJSON ...PagedResult) []string {
	totalPages := len(pagesOfResourcesJSON)
	totalResults := 0
	for _, pageOfResourcesJSON := range pagesOfResourcesJSON {
		totalResults += len(pageOfResourcesJSON.Resources)
	}

	// iterate through each page of resources and build a list of paginated responses
	var resultPages []string
	for i, pageOfResourcesJSON := range pagesOfResourcesJSON {
		pageIndex := i + 1
		resourcesPerPage := len(pageOfResourcesJSON.Resources)

		p := resultTemplate{
			TotalResults:     totalResults,
			TotalPages:       totalPages,
			FirstPage:        fmt.Sprintf("%s?page=1&per_page=%d", defaultAPIResourcePath, resourcesPerPage),
			LastPage:         fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, totalPages, resourcesPerPage),
			Resources:        strings.Join(pageOfResourcesJSON.Resources, ","),
			Apps:             strings.Join(pageOfResourcesJSON.Apps, ","),
			Spaces:           strings.Join(pageOfResourcesJSON.Spaces, ","),
			Organizations:    strings.Join(pageOfResourcesJSON.Organizations, ","),
			Domains:          strings.Join(pageOfResourcesJSON.Domains, ","),
			Users:            strings.Join(pageOfResourcesJSON.Users, ","),
			Routes:           strings.Join(pageOfResourcesJSON.Routes, ","),
			ServiceOfferings: strings.Join(pageOfResourcesJSON.ServiceOfferings, ","),
			ServiceInstances: strings.Join(pageOfResourcesJSON.ServiceInstances, ","),
		}
		if pageIndex < totalPages {
			p.NextPage = fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, pageIndex+1, resourcesPerPage)
		}
		if pageIndex > 1 {
			p.PreviousPage = fmt.Sprintf("%s?page=%d&per_page=%d", defaultAPIResourcePath, pageIndex-1, resourcesPerPage)
		}

		t, err := template.New("page").Parse(listTemplate)
		if err != nil {
			panic(err)
		}
		var h bytes.Buffer
		err = t.Execute(&h, p)
		if err != nil {
			panic(err)
		}
		s := h.String()
		resultPages = append(resultPages, s)

	}
	return resultPages
}

func (o ObjectJSONGenerator) Single(resourceJSON string) []string {
	return []string{resourceJSON}
}

func (o ObjectJSONGenerator) SinglePaged(resourceJSON string) []string {
	return o.Paged([]string{resourceJSON})
}

func (o ObjectJSONGenerator) Paged(pagesOfResourcesJSON ...[]string) []string {
	var pagedResults []PagedResult
	for _, pageOfResourcesJSON := range pagesOfResourcesJSON {
		p := PagedResult{
			Resources: pageOfResourcesJSON,
		}
		pagedResults = append(pagedResults, p)
	}
	return o.PagedWithInclude(pagedResults...)
}

func (o ObjectJSONGenerator) Array(resourcesJSON ...string) string {
	return "[" + strings.Join(resourcesJSON, ",") + "]"
}

func (o ObjectJSONGenerator) renderTemplate(rt *JSONResource, fileName string) *JSONResource {
	_, filename, _, _ := runtime.Caller(1)
	p := path.Join(path.Dir(filename), "template", fileName)
	f, err := os.ReadFile(p)
	if err != nil {
		panic(err)
	}

	t, err := template.New("resource").Parse(string(f))
	if err != nil {
		panic(err)
	}
	var b bytes.Buffer
	err = t.Execute(&b, rt)
	if err != nil {
		panic(err)
	}
	rt.JSON = b.String()
	return rt
}

const listTemplate = `
{
  "pagination": {
    "total_results": {{.TotalResults}},
    "total_pages": {{.TotalPages}},
    "first": { "href": "{{.FirstPage}}" },
    "last": { "href": "{{.LastPage}}" },
    {{if .NextPage}}"next": { "href": "{{.NextPage}}" },{{else}}"next": null,{{end}}
    {{if .PreviousPage}}"previous": { "href": "{{.PreviousPage}}" }{{else}}"previous": null{{end}}
  },
  "resources": [
    {{.Resources}}
  ],
  "included": {
    "apps": [
      {{.Apps}}
    ],
    "spaces": [
      {{.Spaces}}
    ],
    "domains": [
      {{.Domains}}
    ],
    "users": [
      {{.Users}}
    ],
    "routes": [
      {{.Routes}}
    ],
    "service_offerings": [
      {{.ServiceOfferings}}
    ],
    "service_instances": [
      {{.ServiceInstances}}
    ],
    "organizations": [
      {{.Organizations}}
    ]
  }
}
`

const singleTemplate = `
{
    "apps": [
      {{.Apps}}
    ],
    "spaces": [
      {{.Spaces}}
    ],
    "domains": [
      {{.Domains}}
    ],
    "users": [
      {{.Users}}
    ],
    "routes": [
      {{.Routes}}
    ],
    "service_offerings": [
      {{.ServiceOfferings}}
    ],
    "service_instances": [
      {{.ServiceInstances}}
    ],
    "organizations": [
      {{.Organizations}}
    ]
}
`
