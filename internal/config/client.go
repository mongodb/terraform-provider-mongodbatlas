package config

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	adminpreview "github.com/mongodb/atlas-sdk-go/admin"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	admin20240805 "go.mongodb.org/atlas-sdk/v20240805005/admin"
	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
	matlasClient "go.mongodb.org/atlas/mongodbatlas"
	realmAuth "go.mongodb.org/realm/auth"
	"go.mongodb.org/realm/realm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mongodb-forks/digest"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/version"
)

const (
	toolName              = "terraform-provider-mongodbatlas"
	terraformPlatformName = "Terraform"
)

var (
	jsonCheck       = regexp.MustCompile(`(?i:(?:application|text)/(?:vnd\.[^;]+\+)?json)`)
	xmlCheck        = regexp.MustCompile(`(?i:(?:application|text)/xml)`)
	queryParamSplit = regexp.MustCompile(`(^|&)([^&]+)`)
	queryDescape    = strings.NewReplacer("%5B", "[", "%5D", "]")
)

// MongoDBClient contains the mongodbatlas clients and configurations
type MongoDBClient struct {
	Atlas           *matlasClient.Client
	AtlasV2         *admin.APIClient
	AtlasPreview    *adminpreview.APIClient
	AtlasV220240805 *admin20240805.APIClient // used in advanced_cluster to avoid adopting 2024-10-23 release with ISS autoscaling
	AtlasV220240530 *admin20240530.APIClient // used in advanced_cluster and cloud_backup_schedule for avoiding breaking changes (supporting deprecated replication_specs.id)
	AtlasV220241113 *admin20241113.APIClient // used in teams and atlas_users to avoiding breaking changes
	Config          *Config
}

// Config contains the configurations needed to use SDKs
type Config struct {
	AssumeRole       *AssumeRole
	PublicKey        string
	PrivateKey       string
	BaseURL          string
	RealmBaseURL     string
	TerraformVersion string
}

type AssumeRole struct {
	Tags              map[string]string
	RoleARN           string
	ExternalID        string
	Policy            string
	SessionName       string
	SourceIdentity    string
	PolicyARNs        []string
	TransitiveTagKeys []string
	Duration          time.Duration
}

type SecretData struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type PlatformVersion struct {
	Name    string
	Version string
}

// NewClient func...
func (c *Config) NewClient(ctx context.Context) (any, error) {
	// setup a transport to handle digest
	transport := digest.NewTransport(cast.ToString(c.PublicKey), cast.ToString(c.PrivateKey))

	// initialize the client
	client, err := transport.Client()
	if err != nil {
		return nil, err
	}

	client.Transport = logging.NewTransport("MongoDB Atlas", transport)

	optsAtlas := []matlasClient.ClientOpt{matlasClient.SetUserAgent(userAgent(c))}
	if c.BaseURL != "" {
		optsAtlas = append(optsAtlas, matlasClient.SetBaseURL(c.BaseURL))
	}

	// Initialize the MongoDB Atlas API Client.
	atlasClient, err := matlasClient.New(client, optsAtlas...)
	if err != nil {
		return nil, err
	}

	sdkV2Client, err := c.newSDKV2Client(client)
	if err != nil {
		return nil, err
	}

	sdkPreviewClient, err := c.newSDKPreviewClient(client)
	if err != nil {
		return nil, err
	}

	sdkV220240530Client, err := c.newSDKV220240530Client(client)
	if err != nil {
		return nil, err
	}

	sdkV220240805Client, err := c.newSDKV220240805Client(client)
	if err != nil {
		return nil, err
	}

	sdkV220241113Client, err := c.newSDKV220241113Client(client)
	if err != nil {
		return nil, err
	}

	clients := &MongoDBClient{
		Atlas:           atlasClient,
		AtlasV2:         sdkV2Client,
		AtlasPreview:    sdkPreviewClient,
		AtlasV220240530: sdkV220240530Client,
		AtlasV220240805: sdkV220240805Client,
		AtlasV220241113: sdkV220241113Client,
		Config:          c,
	}
	return clients, nil
}

func (c *Config) newSDKV2Client(client *http.Client) (*admin.APIClient, error) {
	opts := []admin.ClientModifier{
		admin.UseHTTPClient(client),
		admin.UseUserAgent(userAgent(c)),
		admin.UseBaseURL(c.BaseURL),
		admin.UseDebug(false)}

	sdk, err := admin.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *Config) newSDKPreviewClient(client *http.Client) (*adminpreview.APIClient, error) {
	opts := []adminpreview.ClientModifier{
		adminpreview.UseHTTPClient(client),
		adminpreview.UseUserAgent(userAgent(c)),
		adminpreview.UseBaseURL(c.BaseURL),
		adminpreview.UseDebug(false)}

	sdk, err := adminpreview.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *Config) newSDKV220240530Client(client *http.Client) (*admin20240530.APIClient, error) {
	opts := []admin20240530.ClientModifier{
		admin20240530.UseHTTPClient(client),
		admin20240530.UseUserAgent(userAgent(c)),
		admin20240530.UseBaseURL(c.BaseURL),
		admin20240530.UseDebug(false)}

	sdk, err := admin20240530.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *Config) newSDKV220240805Client(client *http.Client) (*admin20240805.APIClient, error) {
	opts := []admin20240805.ClientModifier{
		admin20240805.UseHTTPClient(client),
		admin20240805.UseUserAgent(userAgent(c)),
		admin20240805.UseBaseURL(c.BaseURL),
		admin20240805.UseDebug(false)}

	sdk, err := admin20240805.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *Config) newSDKV220241113Client(client *http.Client) (*admin20241113.APIClient, error) {
	opts := []admin20241113.ClientModifier{
		admin20241113.UseHTTPClient(client),
		admin20241113.UseUserAgent(userAgent(c)),
		admin20241113.UseBaseURL(c.BaseURL),
		admin20241113.UseDebug(false)}

	sdk, err := admin20241113.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (c *MongoDBClient) GetRealmClient(ctx context.Context) (*realm.Client, error) {
	// Realm
	if c.Config.PublicKey == "" && c.Config.PrivateKey == "" {
		return nil, errors.New("please set `public_key` and `private_key` in order to use the realm client")
	}

	optsRealm := []realm.ClientOpt{realm.SetUserAgent(userAgent(c.Config))}

	authConfig := realmAuth.NewConfig(nil)
	if c.Config.BaseURL != "" && c.Config.RealmBaseURL != "" {
		adminURL := c.Config.RealmBaseURL + "api/admin/v3.0/"
		optsRealm = append(optsRealm, realm.SetBaseURL(adminURL))
		authConfig.AuthURL, _ = url.Parse(adminURL + "auth/providers/mongodb-cloud/login")
	}

	token, err := authConfig.NewTokenFromCredentials(ctx, c.Config.PublicKey, c.Config.PrivateKey)
	if err != nil {
		return nil, err
	}

	clientRealm := realmAuth.NewClient(realmAuth.BasicTokenSource(token))
	clientRealm.Transport = logging.NewTransport("MongoDB Realm", clientRealm.Transport)

	// Initialize the MongoDB Realm API Client.
	realmClient, err := realm.New(clientRealm, optsRealm...)
	if err != nil {
		return nil, err
	}

	return realmClient, nil
}

type APICallParams struct {
	VersionHeader string
	RelativePath  string
	PathParams    map[string]string
	Method        string
}

func (c *MongoDBClient) UntypedAPICall(ctx context.Context, params *APICallParams, bodyReq []byte) (*http.Response, error) {
	localBasePath, _ := c.AtlasV2.GetConfig().ServerURLWithContext(ctx, "")
	localVarPath := localBasePath + params.RelativePath

	for key, value := range params.PathParams {
		localVarPath = strings.ReplaceAll(localVarPath, "{"+key+"}", url.PathEscape(value))
	}

	headerParams := make(map[string]string)
	headerParams["Content-Type"] = params.VersionHeader
	headerParams["Accept"] = params.VersionHeader

	var bodyPost any
	if bodyReq != nil { // if nil slice is sent with application/json content type SDK method returns an error
		bodyPost = bodyReq
	}
	// TODO DELETE apiReq, err := c.AtlasV2.PrepareRequest(ctx, localVarPath, params.Method, bodyPost, headerParams, nil, nil, nil)
	apiReq, err := prepareRequest(ctx, c.AtlasV2, localVarPath, params.Method, bodyPost, headerParams, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	apiResp, err := c.AtlasV2.CallAPI(apiReq)

	if apiResp.StatusCode >= 300 {
		newErr := makeAPIError(apiResp, params.Method, localVarPath)
		return apiResp, newErr
	}

	return apiResp, err
}

type formFile struct {
	fileName     string
	formFileName string
	fileBytes    []byte
}

// PrepareRequest build the request
func prepareRequest(
	ctx context.Context,
	c *admin.APIClient,
	path string, method string,
	postBody any,
	headerParams map[string]string,
	queryParams url.Values,
	formParams url.Values,
	formFiles []formFile) (localVarRequest *http.Request, err error) {
	var body *bytes.Buffer

	// Detect postBody type and post.
	if postBody != nil {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}

		body, err = setBody(postBody, contentType)
		if err != nil {
			return nil, err
		}
	}

	// add form parameters and file if available.
	if strings.HasPrefix(headerParams["Content-Type"], "multipart/form-data") && len(formParams) > 0 || (len(formFiles) > 0) {
		if body != nil {
			return nil, errors.New("cannot specify postBody and multipart form at the same time")
		}
		body = &bytes.Buffer{}
		w := multipart.NewWriter(body)

		for k, v := range formParams {
			for _, iv := range v {
				if strings.HasPrefix(k, "@") { // file
					err = addFile(w, k[1:], iv)
					if err != nil {
						return nil, err
					}
				} else { // form value
					err = w.WriteField(k, iv)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		for _, formFile := range formFiles {
			if !(len(formFile.fileBytes) > 0 && formFile.fileName != "") {
				continue
			}
			w.Boundary()
			part, err1 := w.CreateFormFile(formFile.formFileName, filepath.Base(formFile.fileName))
			if err1 != nil {
				return nil, err
			}
			_, err1 = part.Write(formFile.fileBytes)
			if err1 != nil {
				return nil, err
			}
		}

		// Set the Boundary in the Content-Type
		headerParams["Content-Type"] = w.FormDataContentType()

		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
		w.Close()
	}

	if strings.HasPrefix(headerParams["Content-Type"], "application/x-www-form-urlencoded") && len(formParams) > 0 {
		if body != nil {
			return nil, errors.New("cannot specify postBody and x-www-form-urlencoded form at the same time")
		}
		body = &bytes.Buffer{}
		body.WriteString(formParams.Encode())
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
	}

	// Setup path and query parameters
	urlData, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Override request host, if applicable
	if c.GetConfig().Host != "" {
		urlData.Host = c.GetConfig().Host
	}

	// Override request scheme, if applicable
	if c.GetConfig().Scheme != "" {
		urlData.Scheme = c.GetConfig().Scheme
	}

	// Adding Query Param
	query := urlData.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	// Encode the parameters.
	urlData.RawQuery = queryParamSplit.ReplaceAllStringFunc(query.Encode(), func(s string) string {
		pieces := strings.Split(s, "=")
		pieces[0] = queryDescape.Replace(pieces[0])
		return strings.Join(pieces, "=")
	})

	// Generate a new request
	if body != nil {
		localVarRequest, err = http.NewRequest(method, urlData.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, urlData.String(), http.NoBody)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers[h] = []string{v}
		}
		localVarRequest.Header = headers
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.GetConfig().UserAgent)

	if ctx != nil {
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)
	}

	for header, value := range c.GetConfig().DefaultHeader {
		localVarRequest.Header.Add(header, value)
	}
	return localVarRequest, nil
}

func makeAPIError(res *http.Response, httpMethod, httpPath string) error {
	defer res.Body.Close()

	newErr := new(admin.GenericOpenAPIError)
	newErr.SetError(res.Status)

	localVarBody, err := io.ReadAll(res.Body)
	if err != nil {
		newErr.SetError(fmt.Sprintf("(%s) failed to read response body: %s", res.Status, err.Error()))
		return newErr
	}
	// newErr.body = localVarBody // TODO: body not set

	var v admin.ApiError
	err = decode(&v, io.NopCloser(bytes.NewBuffer(localVarBody)), res.Header.Get("Content-Type"))
	if err != nil {
		newErr.SetError(fmt.Sprintf("(%s) failed to decode response body: %s", res.Status, err.Error()))
		return newErr
	}
	newErr.SetError(admin.FormatErrorMessageWithDetails(res.Status, httpMethod, httpPath, v))
	newErr.SetModel(v)
	return newErr
}

func decode(v any, b io.ReadCloser, contentType string) (err error) {
	switch r := v.(type) {
	case *string:
		buf, err := io.ReadAll(b)
		_ = b.Close()
		if err != nil {
			return err
		}
		*r = string(buf)
		return nil
	case *io.ReadCloser:
		*r = b
		return nil
	case **io.ReadCloser:
		*r = &b
		return nil
	default:
		buf, err := io.ReadAll(b)
		_ = b.Close()
		if err != nil {
			return err
		}
		if len(buf) == 0 {
			return nil
		}
		if xmlCheck.MatchString(contentType) {
			return xml.Unmarshal(buf, v)
		}
		if jsonCheck.MatchString(contentType) {
			if actualObj, ok := v.(interface{ GetActualInstance() any }); ok { // oneOf, anyOf schemas
				if unmarshalObj, ok := actualObj.(interface{ UnmarshalJSON([]byte) error }); ok { // make sure it has UnmarshalJSON defined
					if err := unmarshalObj.UnmarshalJSON(buf); err != nil {
						return err
					}
				} else {
					return errors.New("unknown type with GetActualInstance but no unmarshalObj.UnmarshalJSON defined")
				}
			} else if err := json.Unmarshal(buf, v); err != nil { // simple model
				return err
			}
			return nil
		}
		return errors.New("undefined response type")
	}
}

func userAgent(c *Config) string {
	platformVersions := []PlatformVersion{
		{toolName, version.ProviderVersion},
		{terraformPlatformName, c.TerraformVersion},
	}

	var parts []string
	for _, info := range platformVersions {
		part := fmt.Sprintf("%s/%s", info.Name, info.Version)
		parts = append(parts, part)
	}

	return strings.Join(parts, " ")
}

// Set request body from an any
func setBody(body any, contentType string) (bodyBuf *bytes.Buffer, err error) {
	bodyBuf = &bytes.Buffer{}
	switch v := body.(type) {
	case io.Reader:
		_, err = bodyBuf.ReadFrom(v)
	case *io.ReadCloser:
		_, err = bodyBuf.ReadFrom(*v)
	case []byte:
		_, err = bodyBuf.Write(v)
	case string:
		_, err = bodyBuf.WriteString(v)
	case *string:
		_, err = bodyBuf.WriteString(*v)
	default:
		if jsonCheck.MatchString(contentType) {
			err = json.NewEncoder(bodyBuf).Encode(body)
		} else if xmlCheck.MatchString(contentType) {
			err = xml.NewEncoder(bodyBuf).Encode(body)
		}
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("invalid body type %s", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

// detectContentType method is used to figure out `Request.Body` content type for request header
func detectContentType(body any) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	case reflect.Slice:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else {
			contentType = "application/json; charset=utf-8"
		}
	case reflect.Invalid, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array, reflect.Chan,
		reflect.Func, reflect.Interface, reflect.UnsafePointer:
		contentType = "text/plain; charset=utf-8"
	}

	return contentType
}

// Add a file to the multipart request
func addFile(w *multipart.Writer, fieldName, path string) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	part, err := w.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	return err
}
