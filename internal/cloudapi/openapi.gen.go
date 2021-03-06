// Package cloudapi provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package cloudapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// AWSUploadRequestOptions defines model for AWSUploadRequestOptions.
type AWSUploadRequestOptions struct {
	Ec2    AWSUploadRequestOptionsEc2 `json:"ec2"`
	Region string                     `json:"region"`
	S3     AWSUploadRequestOptionsS3  `json:"s3"`
}

// AWSUploadRequestOptionsEc2 defines model for AWSUploadRequestOptionsEc2.
type AWSUploadRequestOptionsEc2 struct {
	AccessKeyId       string    `json:"access_key_id"`
	SecretAccessKey   string    `json:"secret_access_key"`
	ShareWithAccounts *[]string `json:"share_with_accounts,omitempty"`
	SnapshotName      *string   `json:"snapshot_name,omitempty"`
}

// AWSUploadRequestOptionsS3 defines model for AWSUploadRequestOptionsS3.
type AWSUploadRequestOptionsS3 struct {
	AccessKeyId     string `json:"access_key_id"`
	Bucket          string `json:"bucket"`
	SecretAccessKey string `json:"secret_access_key"`
}

// AWSUploadStatus defines model for AWSUploadStatus.
type AWSUploadStatus struct {
	Ami    string `json:"ami"`
	Region string `json:"region"`
}

// AzureUploadRequestOptions defines model for AzureUploadRequestOptions.
type AzureUploadRequestOptions struct {

	// Name of the uploaded image. It must be unique in the given resource group.
	// If name is omitted from the request, a random one based on a UUID is
	// generated.
	ImageName *string `json:"image_name,omitempty"`

	// Location where the image should be uploaded and registered. This link explain
	// how to list all locations:
	// https://docs.microsoft.com/en-us/cli/azure/account?view=azure-cli-latest#az_account_list_locations'
	Location string `json:"location"`

	// Name of the resource group where the image should be uploaded.
	ResourceGroup string `json:"resource_group"`

	// ID of subscription where the image should be uploaded.
	SubscriptionId string `json:"subscription_id"`

	// ID of the tenant where the image should be uploaded. This link explains how
	// to find it in the Azure Portal:
	// https://docs.microsoft.com/en-us/azure/active-directory/fundamentals/active-directory-how-to-find-tenant
	TenantId string `json:"tenant_id"`
}

// AzureUploadStatus defines model for AzureUploadStatus.
type AzureUploadStatus struct {
	ImageName string `json:"image_name"`
}

// ComposeRequest defines model for ComposeRequest.
type ComposeRequest struct {
	Customizations *Customizations `json:"customizations,omitempty"`
	Distribution   string          `json:"distribution"`
	ImageRequests  []ImageRequest  `json:"image_requests"`
}

// ComposeResult defines model for ComposeResult.
type ComposeResult struct {
	Id string `json:"id"`
}

// ComposeStatus defines model for ComposeStatus.
type ComposeStatus struct {
	ImageStatus ImageStatus `json:"image_status"`
}

// Customizations defines model for Customizations.
type Customizations struct {
	Packages     *[]string     `json:"packages,omitempty"`
	Subscription *Subscription `json:"subscription,omitempty"`
}

// GCPUploadRequestOptions defines model for GCPUploadRequestOptions.
type GCPUploadRequestOptions struct {

	// Name of an existing STANDARD Storage class Bucket.
	Bucket string `json:"bucket"`

	// The name to use for the imported and shared Compute Node image.
	// The image name must be unique within the GCP project, which is used
	// for the OS image upload and import. If not specified a random
	// 'composer-api-<uuid>' string is used as the image name.
	ImageName *string `json:"image_name,omitempty"`

	// The GCP region where the OS image will be imported to and shared from.
	// The value must be a valid GCP location. See https://cloud.google.com/storage/docs/locations.
	// If not specified, the multi-region location closest to the source
	// (source Storage Bucket location) is chosen automatically.
	Region *string `json:"region,omitempty"`

	// List of valid Google accounts to share the imported Compute Node image with.
	// Each string must contain a specifier of the account type. Valid formats are:
	//   - 'user:{emailid}': An email address that represents a specific
	//     Google account. For example, 'alice@example.com'.
	//   - 'serviceAccount:{emailid}': An email address that represents a
	//     service account. For example, 'my-other-app@appspot.gserviceaccount.com'.
	//   - 'group:{emailid}': An email address that represents a Google group.
	//     For example, 'admins@example.com'.
	//   - 'domain:{domain}': The G Suite domain (primary) that represents all
	//     the users of that domain. For example, 'google.com' or 'example.com'.
	// If not specified, the imported Compute Node image is not shared with any
	// account.
	ShareWithAccounts *[]string `json:"share_with_accounts,omitempty"`
}

// GCPUploadStatus defines model for GCPUploadStatus.
type GCPUploadStatus struct {
	ImageName string `json:"image_name"`
	ProjectId string `json:"project_id"`
}

// ImageRequest defines model for ImageRequest.
type ImageRequest struct {
	Architecture  string        `json:"architecture"`
	ImageType     string        `json:"image_type"`
	Repositories  []Repository  `json:"repositories"`
	UploadRequest UploadRequest `json:"upload_request"`
}

// ImageStatus defines model for ImageStatus.
type ImageStatus struct {
	Status       ImageStatusValue `json:"status"`
	UploadStatus *UploadStatus    `json:"upload_status,omitempty"`
}

// ImageStatusValue defines model for ImageStatusValue.
type ImageStatusValue string

// List of ImageStatusValue
const (
	ImageStatusValue_building    ImageStatusValue = "building"
	ImageStatusValue_failure     ImageStatusValue = "failure"
	ImageStatusValue_pending     ImageStatusValue = "pending"
	ImageStatusValue_registering ImageStatusValue = "registering"
	ImageStatusValue_success     ImageStatusValue = "success"
	ImageStatusValue_uploading   ImageStatusValue = "uploading"
)

// Repository defines model for Repository.
type Repository struct {
	Baseurl    *string `json:"baseurl,omitempty"`
	Metalink   *string `json:"metalink,omitempty"`
	Mirrorlist *string `json:"mirrorlist,omitempty"`
	Rhsm       bool    `json:"rhsm"`
}

// Subscription defines model for Subscription.
type Subscription struct {
	ActivationKey string `json:"activation-key"`
	BaseUrl       string `json:"base-url"`
	Insights      bool   `json:"insights"`
	Organization  int    `json:"organization"`
	ServerUrl     string `json:"server-url"`
}

// UploadRequest defines model for UploadRequest.
type UploadRequest struct {
	Options interface{} `json:"options"`
	Type    UploadTypes `json:"type"`
}

// UploadStatus defines model for UploadStatus.
type UploadStatus struct {
	Options interface{} `json:"options"`
	Status  string      `json:"status"`
	Type    UploadTypes `json:"type"`
}

// UploadTypes defines model for UploadTypes.
type UploadTypes string

// List of UploadTypes
const (
	UploadTypes_aws   UploadTypes = "aws"
	UploadTypes_azure UploadTypes = "azure"
	UploadTypes_gcp   UploadTypes = "gcp"
)

// Version defines model for Version.
type Version struct {
	Version string `json:"version"`
}

// ComposeJSONBody defines parameters for Compose.
type ComposeJSONBody ComposeRequest

// ComposeRequestBody defines body for Compose for application/json ContentType.
type ComposeJSONRequestBody ComposeJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = http.DefaultClient
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditor = fn
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Compose request  with any body
	ComposeWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error)

	Compose(ctx context.Context, body ComposeJSONRequestBody) (*http.Response, error)

	// ComposeStatus request
	ComposeStatus(ctx context.Context, id string) (*http.Response, error)

	// GetOpenapiJson request
	GetOpenapiJson(ctx context.Context) (*http.Response, error)

	// GetVersion request
	GetVersion(ctx context.Context) (*http.Response, error)
}

func (c *Client) ComposeWithBody(ctx context.Context, contentType string, body io.Reader) (*http.Response, error) {
	req, err := NewComposeRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) Compose(ctx context.Context, body ComposeJSONRequestBody) (*http.Response, error) {
	req, err := NewComposeRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) ComposeStatus(ctx context.Context, id string) (*http.Response, error) {
	req, err := NewComposeStatusRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetOpenapiJson(ctx context.Context) (*http.Response, error) {
	req, err := NewGetOpenapiJsonRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

func (c *Client) GetVersion(ctx context.Context) (*http.Response, error) {
	req, err := NewGetVersionRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(ctx, req)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewComposeRequest calls the generic Compose builder with application/json body
func NewComposeRequest(server string, body ComposeJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewComposeRequestWithBody(server, "application/json", bodyReader)
}

// NewComposeRequestWithBody generates requests for Compose with any type of body
func NewComposeRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/compose")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryUrl.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// NewComposeStatusRequest generates requests for ComposeStatus
func NewComposeStatusRequest(server string, id string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParam("simple", false, "id", id)
	if err != nil {
		return nil, err
	}

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/compose/%s", pathParam0)
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetOpenapiJsonRequest generates requests for GetOpenapiJson
func NewGetOpenapiJsonRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/openapi.json")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetVersionRequest generates requests for GetVersion
func NewGetVersionRequest(server string) (*http.Request, error) {
	var err error

	queryUrl, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	basePath := fmt.Sprintf("/version")
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	queryUrl, err = queryUrl.Parse(basePath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// Compose request  with any body
	ComposeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*ComposeResponse, error)

	ComposeWithResponse(ctx context.Context, body ComposeJSONRequestBody) (*ComposeResponse, error)

	// ComposeStatus request
	ComposeStatusWithResponse(ctx context.Context, id string) (*ComposeStatusResponse, error)

	// GetOpenapiJson request
	GetOpenapiJsonWithResponse(ctx context.Context) (*GetOpenapiJsonResponse, error)

	// GetVersion request
	GetVersionWithResponse(ctx context.Context) (*GetVersionResponse, error)
}

type ComposeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *ComposeResult
}

// Status returns HTTPResponse.Status
func (r ComposeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ComposeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ComposeStatusResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ComposeStatus
}

// Status returns HTTPResponse.Status
func (r ComposeStatusResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ComposeStatusResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetOpenapiJsonResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetOpenapiJsonResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetOpenapiJsonResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetVersionResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Version
}

// Status returns HTTPResponse.Status
func (r GetVersionResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetVersionResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ComposeWithBodyWithResponse request with arbitrary body returning *ComposeResponse
func (c *ClientWithResponses) ComposeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader) (*ComposeResponse, error) {
	rsp, err := c.ComposeWithBody(ctx, contentType, body)
	if err != nil {
		return nil, err
	}
	return ParseComposeResponse(rsp)
}

func (c *ClientWithResponses) ComposeWithResponse(ctx context.Context, body ComposeJSONRequestBody) (*ComposeResponse, error) {
	rsp, err := c.Compose(ctx, body)
	if err != nil {
		return nil, err
	}
	return ParseComposeResponse(rsp)
}

// ComposeStatusWithResponse request returning *ComposeStatusResponse
func (c *ClientWithResponses) ComposeStatusWithResponse(ctx context.Context, id string) (*ComposeStatusResponse, error) {
	rsp, err := c.ComposeStatus(ctx, id)
	if err != nil {
		return nil, err
	}
	return ParseComposeStatusResponse(rsp)
}

// GetOpenapiJsonWithResponse request returning *GetOpenapiJsonResponse
func (c *ClientWithResponses) GetOpenapiJsonWithResponse(ctx context.Context) (*GetOpenapiJsonResponse, error) {
	rsp, err := c.GetOpenapiJson(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetOpenapiJsonResponse(rsp)
}

// GetVersionWithResponse request returning *GetVersionResponse
func (c *ClientWithResponses) GetVersionWithResponse(ctx context.Context) (*GetVersionResponse, error) {
	rsp, err := c.GetVersion(ctx)
	if err != nil {
		return nil, err
	}
	return ParseGetVersionResponse(rsp)
}

// ParseComposeResponse parses an HTTP response from a ComposeWithResponse call
func ParseComposeResponse(rsp *http.Response) (*ComposeResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ComposeResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest ComposeResult
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	}

	return response, nil
}

// ParseComposeStatusResponse parses an HTTP response from a ComposeStatusWithResponse call
func ParseComposeStatusResponse(rsp *http.Response) (*ComposeStatusResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &ComposeStatusResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ComposeStatus
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetOpenapiJsonResponse parses an HTTP response from a GetOpenapiJsonWithResponse call
func ParseGetOpenapiJsonResponse(rsp *http.Response) (*GetOpenapiJsonResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetOpenapiJsonResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	}

	return response, nil
}

// ParseGetVersionResponse parses an HTTP response from a GetVersionWithResponse call
func ParseGetVersionResponse(rsp *http.Response) (*GetVersionResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &GetVersionResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Version
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create compose
	// (POST /compose)
	Compose(w http.ResponseWriter, r *http.Request)
	// The status of a compose
	// (GET /compose/{id})
	ComposeStatus(w http.ResponseWriter, r *http.Request, id string)
	// get the openapi json specification
	// (GET /openapi.json)
	GetOpenapiJson(w http.ResponseWriter, r *http.Request)
	// get the service version
	// (GET /version)
	GetVersion(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Compose operation middleware
func (siw *ServerInterfaceWrapper) Compose(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.Compose(w, r.WithContext(ctx))
}

// ComposeStatus operation middleware
func (siw *ServerInterfaceWrapper) ComposeStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid format for parameter id: %s", err), http.StatusBadRequest)
		return
	}

	siw.Handler.ComposeStatus(w, r.WithContext(ctx), id)
}

// GetOpenapiJson operation middleware
func (siw *ServerInterfaceWrapper) GetOpenapiJson(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.GetOpenapiJson(w, r.WithContext(ctx))
}

// GetVersion operation middleware
func (siw *ServerInterfaceWrapper) GetVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	siw.Handler.GetVersion(w, r.WithContext(ctx))
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerFromMux(si, chi.NewRouter())
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	r.Group(func(r chi.Router) {
		r.Post("/compose", wrapper.Compose)
	})
	r.Group(func(r chi.Router) {
		r.Get("/compose/{id}", wrapper.ComposeStatus)
	})
	r.Group(func(r chi.Router) {
		r.Get("/openapi.json", wrapper.GetOpenapiJson)
	})
	r.Group(func(r chi.Router) {
		r.Get("/version", wrapper.GetVersion)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file.
func GetSwagger() (*openapi3.Swagger, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error loading Swagger: %s", err)
	}
	return swagger, nil
}
