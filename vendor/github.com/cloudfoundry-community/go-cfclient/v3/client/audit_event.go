package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type AuditEventClient commonClient

// AuditEventListOptions list filters
type AuditEventListOptions struct {
	*ListOptions

	Types             Filter          `qs:"types"`        //  list of event types to filter by
	TargetGUIDs       ExclusionFilter `qs:"target_guids"` // list of target guids to filter by
	OrganizationGUIDs Filter          `qs:"organization_guids"`
	SpaceGUIDs        Filter          `qs:"space_guids"`
}

// NewAuditEventListOptions creates new options to pass to list
func NewAuditEventListOptions() *AuditEventListOptions {
	return &AuditEventListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o AuditEventListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// First returns the first audit event matching the options or an error when less than 1 match
func (c *AuditEventClient) First(ctx context.Context, opts *AuditEventListOptions) (*resource.AuditEvent, error) {
	return First[*AuditEventListOptions, *resource.AuditEvent](opts, func(opts *AuditEventListOptions) ([]*resource.AuditEvent, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get retrieves the specified audit event
func (c *AuditEventClient) Get(ctx context.Context, guid string) (*resource.AuditEvent, error) {
	var a resource.AuditEvent
	err := c.client.get(ctx, path.Format("/v3/audit_events/%s", guid), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// List pages all audit events the user has access to
func (c *AuditEventClient) List(ctx context.Context, opts *AuditEventListOptions) ([]*resource.AuditEvent, *Pager, error) {
	if opts == nil {
		opts = NewAuditEventListOptions()
	}
	var res resource.AuditEventList
	err := c.client.list(ctx, "/v3/audit_events", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all audit events the user has access to
func (c *AuditEventClient) ListAll(ctx context.Context, opts *AuditEventListOptions) ([]*resource.AuditEvent, error) {
	if opts == nil {
		opts = NewAuditEventListOptions()
	}

	var all []*resource.AuditEvent
	for {
		page, pager, err := c.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, nil
}

// Single returns a single audit event matching the options or an error if not exactly 1 match
func (c *AuditEventClient) Single(ctx context.Context, opts *AuditEventListOptions) (*resource.AuditEvent, error) {
	return Single[*AuditEventListOptions, *resource.AuditEvent](opts, func(opts *AuditEventListOptions) ([]*resource.AuditEvent, *Pager, error) {
		return c.List(ctx, opts)
	})
}
