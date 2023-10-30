package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type RouteClient commonClient

// RouteListOptions list filters
type RouteListOptions struct {
	*ListOptions

	AppGUIDs             Filter `qs:"app_guids"`
	SpaceGUIDs           Filter `qs:"space_guids"`
	DomainGUIDs          Filter `qs:"domain_guids"`
	OrganizationGUIDs    Filter `qs:"organization_guids"`
	ServiceInstanceGUIDs Filter `qs:"service_instance_guids"`

	Hosts Filter `qs:"hosts"`
	Paths Filter `qs:"paths"`
	Ports Filter `qs:"ports"`

	Include resource.RouteIncludeType `qs:"include"`
}

// NewRouteListOptions creates new options to pass to list
func NewRouteListOptions() *RouteListOptions {
	return &RouteListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o RouteListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// RouteReservationListOptions list filters
type RouteReservationListOptions struct {
	*ListOptions

	Hosts string `qs:"host"`
	Paths string `qs:"path"`
	Ports int    `qs:"port"`
}

// NewRouteReservationListOptions creates new options to pass to IsRouteReserved
func NewRouteReservationListOptions() *RouteReservationListOptions {
	return &RouteReservationListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o RouteReservationListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new route
func (c *RouteClient) Create(ctx context.Context, r *resource.RouteCreate) (*resource.Route, error) {
	var Route resource.Route
	_, err := c.client.post(ctx, "/v3/routes", r, &Route)
	if err != nil {
		return nil, err
	}
	return &Route, nil
}

// Delete the specified route asynchronously and return a jobGUID
func (c *RouteClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/routes/%s", guid))
}

// DeleteUnmappedRoutesForSpace deletes all routes in a space that are not mapped to any applications and not
// bound to any service instances and returns the async JobGUID
func (c *RouteClient) DeleteUnmappedRoutesForSpace(ctx context.Context, spaceGUID string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/spaces/%s/routes?unmapped=true", spaceGUID))
}

// First returns the first route matching the options or an error when less than 1 match
func (c *RouteClient) First(ctx context.Context, opts *RouteListOptions) (*resource.Route, error) {
	return First[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// FirstForApp returns the first route matching the options and app or an error when less than 1 match
func (c *RouteClient) FirstForApp(ctx context.Context, appGUID string, opts *RouteListOptions) (*resource.Route, error) {
	return First[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Get the specified route
func (c *RouteClient) Get(ctx context.Context, guid string) (*resource.Route, error) {
	var r resource.Route
	err := c.client.get(ctx, path.Format("/v3/routes/%s", guid), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// GetIncludeDomain allows callers to fetch a route and include the parent domain
func (c *RouteClient) GetIncludeDomain(ctx context.Context, guid string) (*resource.Route, *resource.Domain, error) {
	var r resource.RouteWithIncluded
	err := c.client.get(ctx, path.Format("/v3/routes/%s?include=%s", guid, resource.RouteIncludeDomain), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.Route, r.Included.Domains[0], nil
}

// GetIncludeSpace allows callers to fetch a route and include the parent space
func (c *RouteClient) GetIncludeSpace(ctx context.Context, guid string) (*resource.Route, *resource.Space, error) {
	var r resource.RouteWithIncluded
	err := c.client.get(ctx, path.Format("/v3/routes/%s?include=%s", guid, resource.RouteIncludeSpaceOrganization), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.Route, r.Included.Spaces[0], nil
}

// GetIncludeSpaceAndOrganization allows callers to fetch a route and include the parent space and organization
func (c *RouteClient) GetIncludeSpaceAndOrganization(ctx context.Context, guid string) (*resource.Route, *resource.Space, *resource.Organization, error) {
	var r resource.RouteWithIncluded
	err := c.client.get(ctx, path.Format("/v3/routes/%s?include=%s", guid, resource.RouteIncludeSpaceOrganization), &r)
	if err != nil {
		return nil, nil, nil, err
	}
	return &r.Route, r.Included.Spaces[0], r.Included.Organizations[0], nil
}

// GetSharedSpacesRelationships retrieves the spaces that the route has been shared to
func (c *RouteClient) GetSharedSpacesRelationships(ctx context.Context, guid string) (*resource.RouteSharedSpaceRelationships, error) {
	var r resource.RouteSharedSpaceRelationships
	err := c.client.get(ctx, path.Format("/v3/routes/%s/relationships/shared_spaces", guid), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// GetDestinations retrieves all destinations associated with a route
func (c *RouteClient) GetDestinations(ctx context.Context, guid string) (*resource.RouteDestinations, error) {
	var r resource.RouteDestinations
	err := c.client.get(ctx, path.Format("/v3/routes/%s/destinations", guid), &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// InsertDestinations add one or more destinations to a route, preserving any existing destinations
//
// Note that weighted destinations cannot be added with this endpoint. To add weighted destinations, replace
// all destinations for a route at once using the replace destinations endpoint.
func (c *RouteClient) InsertDestinations(ctx context.Context, guid string, dest []*resource.RouteDestinationInsertOrReplace) (*resource.RouteDestinations, error) {
	destinations := &resource.RouteDestinationsInsertOrReplace{
		Destinations: dest,
	}
	var r resource.RouteDestinations
	_, err := c.client.post(ctx, path.Format("/v3/routes/%s/destinations", guid), destinations, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// IsRouteReserved checks if a specific route for a domain exists, regardless of the user’s visibility for the
// route in case the route belongs to a space the user does not belong to
func (c *RouteClient) IsRouteReserved(ctx context.Context, domainGUID string, opts *RouteReservationListOptions) (bool, error) {
	if opts == nil {
		opts = NewRouteReservationListOptions()
	}
	var match map[string]bool
	err := c.client.list(ctx, "/v3/domains/"+domainGUID+"/route_reservations", opts.ToQueryString, &match)
	if err != nil {
		return false, err
	}
	return match["matching_route"], nil
}

// List pages routes the user has access to
func (c *RouteClient) List(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeNone

	var res resource.RouteList
	err := c.client.list(ctx, "/v3/routes", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all routes the user has access to
func (c *RouteClient) ListAll(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	return AutoPage[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForApp pages routes for the specified app the user has access to
func (c *RouteClient) ListForApp(ctx context.Context, appGUID string, opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeNone

	var res resource.RouteList
	err := c.client.list(ctx, "/v3/apps/"+appGUID+"/routes", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all routes for the specified app the user has access to
func (c *RouteClient) ListForAppAll(ctx context.Context, appGUID string, opts *RouteListOptions) ([]*resource.Route, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	return AutoPage[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// ListIncludeDomains page all routes the user has access to and include the parent domains
func (c *RouteClient) ListIncludeDomains(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, []*resource.Domain, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeDomain

	var res resource.RouteList
	err := c.client.list(ctx, "/v3/routes", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Domains, pager, nil
}

// ListIncludeDomainsAll retrieves all routes the user has access to and includes the parent domains
func (c *RouteClient) ListIncludeDomainsAll(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, []*resource.Domain, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}

	var all []*resource.Route
	var allDomains []*resource.Domain
	for {
		page, domains, pager, err := c.ListIncludeDomains(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allDomains = append(allDomains, domains...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allDomains, nil
}

// ListIncludeSpaces page all routes the user has access to and include the parent spaces
func (c *RouteClient) ListIncludeSpaces(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, []*resource.Space, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeSpace

	var res resource.RouteList
	err := c.client.list(ctx, "/v3/routes", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, pager, nil
}

// ListIncludeSpacesAll retrieves all routes the user has access to and includes the parent spaces
func (c *RouteClient) ListIncludeSpacesAll(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, []*resource.Space, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}

	var all []*resource.Route
	var allSpaces []*resource.Space
	for {
		page, spaces, pager, err := c.ListIncludeSpaces(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allSpaces = append(allSpaces, spaces...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSpaces, nil
}

// ListIncludeSpacesAndOrganizations page all routes the user has access to and include the parent spaces and organizations
func (c *RouteClient) ListIncludeSpacesAndOrganizations(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, []*resource.Space, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}
	opts.Include = resource.RouteIncludeSpaceOrganization

	var res resource.RouteList
	err := c.client.list(ctx, "/v3/routes", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, res.Included.Organizations, pager, nil
}

// ListIncludeSpacesAndOrganizationsAll retrieves all routes the user has access to and includes the parent spaces and organization
func (c *RouteClient) ListIncludeSpacesAndOrganizationsAll(ctx context.Context, opts *RouteListOptions) ([]*resource.Route, []*resource.Space, []*resource.Organization, error) {
	if opts == nil {
		opts = NewRouteListOptions()
	}

	var all []*resource.Route
	var allSpaces []*resource.Space
	var allOrgs []*resource.Organization
	for {
		page, spaces, orgs, pager, err := c.ListIncludeSpacesAndOrganizations(ctx, opts)
		if err != nil {
			return nil, nil, nil, err
		}
		all = append(all, page...)
		allSpaces = append(allSpaces, spaces...)
		allOrgs = append(allOrgs, orgs...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSpaces, allOrgs, nil
}

// RemoveDestination removes a destination from a route
func (c *RouteClient) RemoveDestination(ctx context.Context, guid, destinationGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/routes/%s/destinations/%s", guid, destinationGUID))
	return err
}

// ReplaceDestinations replaces all destinations for a route, removing any destinations not included in the provided list
//
// If using weighted destinations, all destinations provided here must have a weight specified, and all weights for
// this route must sum to 100. If not, all provided destinations must not have a weight. Mixing weighted and unweighted
// destinations for a route is not allowed.
func (c *RouteClient) ReplaceDestinations(ctx context.Context, guid string, dest []*resource.RouteDestinationInsertOrReplace) (*resource.RouteDestinations, error) {
	destinations := &resource.RouteDestinationsInsertOrReplace{
		Destinations: dest,
	}
	var r resource.RouteDestinations
	_, err := c.client.patch(ctx, path.Format("/v3/routes/%s/destinations", guid), destinations, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// ShareWithSpace shares the route with the specified space
//
// In order to share into a space the requesting user must be a space developer in the target space
func (c *RouteClient) ShareWithSpace(ctx context.Context, guid string, spaceGUID string) (*resource.RouteSharedSpaceRelationships, error) {
	return c.ShareWithSpaces(ctx, guid, []string{spaceGUID})
}

// ShareWithSpaces shares the route with the specified spaces
//
// In order to share into a space the requesting user must be a space developer in the target space
func (c *RouteClient) ShareWithSpaces(ctx context.Context, guid string, spaceGUIDs []string) (*resource.RouteSharedSpaceRelationships, error) {
	req := resource.NewToManyRelationships(spaceGUIDs)
	var relationships resource.RouteSharedSpaceRelationships
	_, err := c.client.post(ctx, path.Format("/v3/routes/%s/relationships/shared_spaces", guid), req, &relationships)
	if err != nil {
		return nil, err
	}
	return &relationships, nil
}

// Single returns a single route matching the options or an error if not exactly 1 match
func (c *RouteClient) Single(ctx context.Context, opts *RouteListOptions) (*resource.Route, error) {
	return Single[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// SingleForApp returns a single route matching the options and app or an error if not exactly 1 match
func (c *RouteClient) SingleForApp(ctx context.Context, appGUID string, opts *RouteListOptions) (*resource.Route, error) {
	return Single[*RouteListOptions, *resource.Route](opts, func(opts *RouteListOptions) ([]*resource.Route, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// TransferOwnership transfers the ownership of a route to another space
//
// Users must have write access for both spaces to perform this action. The original owning space will still
// retain access to the route as a shared space. To completely remove a space from a route, users will have
// to un-share the route.
func (c *RouteClient) TransferOwnership(ctx context.Context, guid string, spaceGUID string) error {
	req := resource.ToOneRelationship{
		Data: &resource.Relationship{
			GUID: spaceGUID,
		},
	}
	_, err := c.client.patch(ctx, path.Format("/v3/routes/%s/relationships/space", guid), req, nil)
	return err
}

// UnShareWithSpace un-shares the route with the specified space
//
// This will automatically unbind any applications bound to this route in the specified space
// Un-sharing a route from a space will not delete any service keys
func (c *RouteClient) UnShareWithSpace(ctx context.Context, guid string, spaceGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/routes/%s/relationships/shared_spaces/%s", guid, spaceGUID))
	return err
}

// UnShareWithSpaces un-shares the route with the specified spaces
//
// This will automatically unbind any applications bound to this route in the specified space
// Un-sharing a route from a space will not delete any service keys
func (c *RouteClient) UnShareWithSpaces(ctx context.Context, guid string, spaceGUIDs []string) error {
	for _, s := range spaceGUIDs {
		err := c.UnShareWithSpace(ctx, guid, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update the specified attributes of the app
func (c *RouteClient) Update(ctx context.Context, guid string, r *resource.RouteUpdate) (*resource.Route, error) {
	var res resource.Route
	_, err := c.client.patch(ctx, path.Format("/v3/routes/%s", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateDestinationProtocol updates the protocol of a route destination (app, port and weight cannot be updated)
//
// Protocol the destination will use. Valid protocols are http1 or http2 if route protocol is http, tcp if route
// protocol is tcp. An empty string will set it to either http1 or tcp based on the route protocol
func (c *RouteClient) UpdateDestinationProtocol(ctx context.Context, guid, destinationGUID, protocol string) (*resource.RouteDestinationWithLinks, error) {
	// use nil/null for empty string
	var p *string
	if protocol != "" {
		p = &protocol
	}
	u := &resource.RouteDestinationProtocolUpdate{
		Protocol: p,
	}
	var r resource.RouteDestinationWithLinks
	_, err := c.client.patch(ctx, path.Format("/v3/routes/%s/destinations/%s", guid, destinationGUID), u, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
