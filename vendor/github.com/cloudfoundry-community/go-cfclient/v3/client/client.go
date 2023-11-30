package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	http2 "net/http"
	"net/url"
	"os"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/check"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/ios"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// Client used to communicate with Cloud Foundry
type Client struct {
	Admin                     *AdminClient
	Applications              *AppClient
	AppFeatures               *AppFeatureClient
	AppUsageEvents            *AppUsageClient
	AuditEvents               *AuditEventClient
	Buildpacks                *BuildpackClient
	Builds                    *BuildClient
	Deployments               *DeploymentClient
	Domains                   *DomainClient
	Droplets                  *DropletClient
	EnvVarGroups              *EnvVarGroupClient
	FeatureFlags              *FeatureFlagClient
	IsolationSegments         *IsolationSegmentClient
	Jobs                      *JobClient
	Manifests                 *ManifestClient
	Organizations             *OrganizationClient
	OrganizationQuotas        *OrganizationQuotaClient
	Packages                  *PackageClient
	Processes                 *ProcessClient
	Revisions                 *RevisionClient
	ResourceMatches           *ResourceMatchClient
	Roles                     *RoleClient
	Root                      *RootClient
	Routes                    *RouteClient
	SecurityGroups            *SecurityGroupClient
	ServiceBrokers            *ServiceBrokerClient
	ServiceCredentialBindings *ServiceCredentialBindingClient
	ServiceInstances          *ServiceInstanceClient
	ServiceOfferings          *ServiceOfferingClient
	ServicePlans              *ServicePlanClient
	ServicePlansVisibility    *ServicePlanVisibilityClient
	ServiceRouteBindings      *ServiceRouteBindingClient
	ServiceUsageEvents        *ServiceUsageClient
	Sidecars                  *SidecarClient
	Spaces                    *SpaceClient
	SpaceFeatures             *SpaceFeatureClient
	SpaceQuotas               *SpaceQuotaClient
	Stacks                    *StackClient
	Tasks                     *TaskClient
	Users                     *UserClient

	common commonClient // Reuse a single struct instead of allocating one for each commonClient on the heap.
	config *config.Config

	unauthenticatedClientProvider *http.UnauthenticatedClientProvider
	unauthenticatedHTTPExecutor   *http.Executor
	authenticatedHTTPExecutor     *http.Executor
	authenticatedClientProvider   *http.OAuthSessionManager
}

type commonClient struct {
	client *Client
}

// New returns a new CF client
func New(config *config.Config) (*Client, error) {
	// construct an unauthenticated root client
	unauthenticatedClientProvider := http.NewUnauthenticatedClientProvider(config.HTTPClient())
	unauthenticatedHTTPExecutor := http.NewExecutor(unauthenticatedClientProvider, config.APIEndpointURL, config.UserAgent)
	rootClient := NewRootClient(unauthenticatedHTTPExecutor)
	err := authServiceDiscovery(context.Background(), config, rootClient)
	if err != nil {
		return nil, err
	}

	// create the client instance
	authenticatedClientProvider := http.NewOAuthSessionManager(config)
	authenticatedHTTPExecutor := http.NewExecutor(authenticatedClientProvider, config.APIEndpointURL, config.UserAgent)
	client := &Client{
		config:                        config,
		unauthenticatedHTTPExecutor:   unauthenticatedHTTPExecutor,
		unauthenticatedClientProvider: unauthenticatedClientProvider,
		authenticatedHTTPExecutor:     authenticatedHTTPExecutor,
		authenticatedClientProvider:   authenticatedClientProvider,
	}

	// populate sub-clients
	client.common.client = client
	client.Root = rootClient
	client.Admin = (*AdminClient)(&client.common)
	client.Applications = (*AppClient)(&client.common)
	client.AppFeatures = (*AppFeatureClient)(&client.common)
	client.AppUsageEvents = (*AppUsageClient)(&client.common)
	client.AuditEvents = (*AuditEventClient)(&client.common)
	client.Buildpacks = (*BuildpackClient)(&client.common)
	client.Builds = (*BuildClient)(&client.common)
	client.Deployments = (*DeploymentClient)(&client.common)
	client.Domains = (*DomainClient)(&client.common)
	client.Droplets = (*DropletClient)(&client.common)
	client.EnvVarGroups = (*EnvVarGroupClient)(&client.common)
	client.FeatureFlags = (*FeatureFlagClient)(&client.common)
	client.IsolationSegments = (*IsolationSegmentClient)(&client.common)
	client.Jobs = (*JobClient)(&client.common)
	client.Manifests = (*ManifestClient)(&client.common)
	client.Organizations = (*OrganizationClient)(&client.common)
	client.OrganizationQuotas = (*OrganizationQuotaClient)(&client.common)
	client.Packages = (*PackageClient)(&client.common)
	client.Processes = (*ProcessClient)(&client.common)
	client.Revisions = (*RevisionClient)(&client.common)
	client.ResourceMatches = (*ResourceMatchClient)(&client.common)
	client.Roles = (*RoleClient)(&client.common)
	client.Routes = (*RouteClient)(&client.common)
	client.SecurityGroups = (*SecurityGroupClient)(&client.common)
	client.ServiceBrokers = (*ServiceBrokerClient)(&client.common)
	client.ServiceCredentialBindings = (*ServiceCredentialBindingClient)(&client.common)
	client.ServiceInstances = (*ServiceInstanceClient)(&client.common)
	client.ServiceOfferings = (*ServiceOfferingClient)(&client.common)
	client.ServicePlans = (*ServicePlanClient)(&client.common)
	client.ServicePlansVisibility = (*ServicePlanVisibilityClient)(&client.common)
	client.ServiceRouteBindings = (*ServiceRouteBindingClient)(&client.common)
	client.ServiceUsageEvents = (*ServiceUsageClient)(&client.common)
	client.Sidecars = (*SidecarClient)(&client.common)
	client.Spaces = (*SpaceClient)(&client.common)
	client.SpaceQuotas = (*SpaceQuotaClient)(&client.common)
	client.SpaceFeatures = (*SpaceFeatureClient)(&client.common)
	client.Stacks = (*StackClient)(&client.common)
	client.Tasks = (*TaskClient)(&client.common)
	client.Users = (*UserClient)(&client.common)
	return client, nil
}

// AccessToken returns the raw encoded OAuth access token without the bearer prefix
func (c *Client) AccessToken(ctx context.Context) (string, error) {
	token, err := c.authenticatedClientProvider.AccessToken(ctx)
	if err != nil {
		return "", err
	}
	return token, nil
}

// SSHCode generates an SSH code that can be used by generic SSH clients to SSH into app instances
func (c *Client) SSHCode(ctx context.Context) (string, error) {
	// need this to grab the SSH client id, should probably be cached in config
	r, err := c.Root.Get(ctx)
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", r.Links.AppSSH.Meta.OauthClient) // client_id，used by cf server

	token, err := c.authenticatedClientProvider.AccessToken(ctx)
	if err != nil {
		return "", err
	}

	req := http.NewRequest(ctx, http2.MethodGet, path.Format("/oauth/authorize?%s", values)).
		WithHeader("Authorization", fmt.Sprintf("bearer %s", token)).
		WithFollowRedirects(false)

	uaaHTTPExecutor := http.NewExecutor(c.unauthenticatedClientProvider, c.config.UAAEndpointURL, c.config.UserAgent)
	resp, err := uaaHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to get one-time code: %w", err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)
	if resp.StatusCode != http2.StatusFound {
		return "", fmt.Errorf(
			"expected UAA to return a 302 location that contains the code, but instead got a %d", resp.StatusCode)
	}

	loc, err := resp.Location()
	if err != nil {
		return "", fmt.Errorf("error getting the redirected location: %w", err)
	}
	codes := loc.Query()["code"]
	if len(codes) != 1 {
		return "", errors.New("unable to acquire one time code from authorization response")
	}

	return codes[0], nil
}

// delete does an HTTP DELETE to the specified endpoint and returns the job ID if any
//
// This function takes the relative API resource path. If the resource returns an async job ID
// then the function returns the job GUID which the caller can reference via the job endpoint.
func (c *Client) delete(ctx context.Context, path string) (string, error) {
	req := http.NewRequest(ctx, http2.MethodDelete, path)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error deleting %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// some endpoints return accepted and others return no content
	if resp.StatusCode != http2.StatusAccepted && resp.StatusCode != http2.StatusNoContent {
		return "", c.decodeError(resp)
	}
	return c.decodeJobIDOrBody(resp, nil)
}

func (c *Client) list(ctx context.Context, urlPathFormat string, queryStrFunc func() (url.Values, error), result any) error {
	params, err := queryStrFunc()
	if err != nil {
		return fmt.Errorf("error while generate query params: %w", err)
	}
	if len(params) > 0 {
		urlPathFormat = strings.TrimSuffix(urlPathFormat+"?"+params.Encode(), "?")
	}
	return c.get(ctx, urlPathFormat, result)
}

// get does an HTTP GET to the specified endpoint and automatically handles unmarshalling
// the result JSON body
func (c *Client) get(ctx context.Context, path string, result any) error {
	if !check.IsNil(result) && !check.IsPointer(result) {
		return errors.New("expected result to be nil or a pointer type")
	}

	req := http.NewRequest(ctx, http2.MethodGet, path)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return fmt.Errorf("error getting %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http2.StatusOK {
		return c.decodeError(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, resp.Body)
		return fmt.Errorf("error decoding %s get response JSON before '%s': %w", path, buf.String(), err)
	}
	return nil
}

// patch does an HTTP PATCH to the specified endpoint and automatically handles the result
// whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to PATCH and an optional
// struct to unmarshall the result body. If the resource returns an async job ID instead of a
// response body, then the body won't be unmarshalled and the function returns the job GUID
// which the caller can reference via the job endpoint.
func (c *Client) patch(ctx context.Context, path string, params any, result any) (string, error) {
	if !check.IsNil(result) && !check.IsPointer(result) {
		return "", errors.New("expected result to be nil or a pointer type")
	}

	req := http.NewRequest(ctx, http2.MethodPatch, path).WithObject(params)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error updating %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http2.StatusOK && resp.StatusCode != http2.StatusAccepted && resp.StatusCode != http2.StatusNoContent {
		return "", c.decodeError(resp)
	}
	return c.decodeJobIDOrBody(resp, result)
}

// post does an HTTP POST to the specified endpoint and automatically handles the result
// whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to POST and an optional
// struct to unmarshall the result body. If the resource returns an async job ID in the Location
// header then the job GUID is returned which the caller can reference via the job endpoint.
func (c *Client) post(ctx context.Context, path string, params, result any) (string, error) {
	if !check.IsNil(result) && !check.IsPointer(result) {
		return "", errors.New("expected result to be a pointer type, or nil")
	}

	req := http.NewRequest(ctx, http2.MethodPost, path).WithObject(params)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error creating %s: %w", path, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// Endpoints return different status codes for posts
	if resp.StatusCode != http2.StatusCreated && resp.StatusCode != http2.StatusOK && resp.StatusCode != http2.StatusAccepted {
		return "", c.decodeError(resp)
	}
	return c.decodeJobIDOrBody(resp, result)
}

// postFileUpload does an HTTP POST to the specified endpoint and automatically handles uploading the specified file
// and handling the result whether that's a JSON body or job ID.
//
// This function takes the relative API resource path, any parameters to POST and an optional
// struct to unmarshall the result body. If the resource returns an async job ID in the Location
// header then the job GUID is returned which the caller can reference via the job endpoint.
func (c *Client) postFileUpload(ctx context.Context, path, fieldName, fileName string, fileToUpload io.Reader, result any) (string, error) {
	if !check.IsNil(result) && !check.IsPointer(result) {
		return "", errors.New("expected result to be a pointer type, or nil")
	}

	requestFile, err := os.CreateTemp("", "upload-*.tmp")
	if err != nil {
		return "", fmt.Errorf("could not create temp file for %s upload form: %w", path, err)
	}
	defer ios.CleanupTempFile(requestFile)

	formWriter := multipart.NewWriter(requestFile)
	part, err := formWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s: %w", path, err)
	}
	_, err = io.Copy(part, fileToUpload)
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s, failed on copy: %w", path, err)
	}
	err = formWriter.Close()
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s, failed to close multipart form writer: %w", path, err)
	}
	_, err = requestFile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s, failed to seek beginning of temp file: %w", path, err)
	}
	fileStats, err := requestFile.Stat()
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s, failed to stat temp file: %w", path, err)
	}

	req := http.NewRequest(ctx, http2.MethodPost, path).
		WithContentType(formWriter.FormDataContentType()).
		WithContentLength(fileStats.Size()).
		WithBody(requestFile)
	resp, err := c.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error uploading file to %s: %w", path, err)
	}
	defer ios.CloseReaderIgnoreError(resp.Body)

	if resp.StatusCode != http2.StatusOK && resp.StatusCode != http2.StatusAccepted {
		return "", c.decodeError(resp)
	}
	return c.decodeJobIDAndBody(resp, result)
}

// decodeError attempts to unmarshall the response body as a CF error
func (c *Client) decodeError(resp *http2.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	// Unmarshal v3 error response
	var errs resource.CloudFoundryErrors
	if err := json.Unmarshal(body, &errs); err != nil {
		return CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	// ensure we got an error back
	if len(errs.Errors) == 0 {
		return CloudFoundryHTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	// TODO handle 2+ errors
	return errs.Errors[0]
}

// decodeJobIDAndBody returns the jobGUID if specified in the Location response header and
// unmarshalls the JSON response body to result if available
func (c *Client) decodeJobIDAndBody(resp *http2.Response, result any) (string, error) {
	jobGUID := c.decodeJobID(resp)
	err := c.decodeBody(resp, result)
	return jobGUID, err
}

// decodeJobIDOrBody returns the jobGUID if specified in the Location response header or
// unmarshalls the JSON response body if no job ID and result is non nil
func (c *Client) decodeJobIDOrBody(resp *http2.Response, result any) (string, error) {
	jobGUID := c.decodeJobID(resp)
	if jobGUID != "" {
		return jobGUID, nil
	}
	return "", c.decodeBody(resp, result)
}

// decodeJobID returns the jobGUID if specified in the Location response header
func (c *Client) decodeJobID(resp *http2.Response) string {
	location, err := resp.Location()
	if err == nil && strings.Contains(location.Path, "jobs") {
		p := strings.Split(location.Path, "/")
		return p[len(p)-1]
	}
	return ""
}

// decodeBody unmarshalls the JSON response body if the result is non nil
func (c *Client) decodeBody(resp *http2.Response, result any) error {
	if result != nil && resp.StatusCode != http2.StatusNoContent {
		err := json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return fmt.Errorf("error decoding response JSON: %w", err)
		}
	}
	return nil
}

// authServiceDiscovery sets the UAA and Login endpoint if the user didn't configure these manually
func authServiceDiscovery(ctx context.Context, config *config.Config, rootClient *RootClient) error {
	if config.UAAEndpointURL != "" && config.LoginEndpointURL != "" {
		return nil
	}
	root, err := rootClient.Get(ctx)
	if err != nil {
		return err
	}
	config.UAAEndpointURL = root.Links.Uaa.Href
	config.LoginEndpointURL = root.Links.Login.Href
	return nil
}
