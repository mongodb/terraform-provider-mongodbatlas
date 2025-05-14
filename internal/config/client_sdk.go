package config

// This package contains functions got from the Atlas Go SDK to support UntypedAPICall.

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

var (
	jsonCheck       = regexp.MustCompile(`(?i:(?:application|text)/(?:vnd\.[^;]+\+)?json)`)
	xmlCheck        = regexp.MustCompile(`(?i:(?:application|text)/xml)`)
	queryParamSplit = regexp.MustCompile(`(^|&)([^&]+)`)
	queryDescape    = strings.NewReplacer("%5B", "[", "%5D", "]")
)

type formFile struct {
	fileName     string
	formFileName string
	fileBytes    []byte
}

// prepareRequest build the request
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

// callAPI do the request.
func callAPI(c *admin.APIClient, request *http.Request) (*http.Response, error) {
	if c.GetConfig().Debug {
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		log.Printf("\n%s\n", string(dump))
	}

	resp, err := c.GetConfig().HTTPClient.Do(request)
	if err != nil {
		return resp, err
	}

	if c.GetConfig().Debug {
		dump, err1 := httputil.DumpResponse(resp, true)
		if err1 != nil {
			return resp, err
		}
		log.Printf("\n%s\n", string(dump))
	}
	return resp, err
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

// setBody sets the request body from an any
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

// detectContentType is used to figure out `Request.Body` content type for request header
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

// addFile adds a file to the multipart request
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
