package deletenotfound_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	adminOlderSDK "go.mongodb.org/atlas-sdk/v20241113005/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/deletenotfound"
)

func response(status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"detail":"not found","errorCode":"USER_NOT_FOUND"}`)),
		Request:    &http.Request{},
	}
}

func TestIsNotFound(t *testing.T) {
	testCases := map[string]struct {
		resp     *http.Response
		err      error
		expected bool
	}{
		"404 with error is tolerated":        {response(http.StatusNotFound), errors.New("404 USER_NOT_FOUND"), true},
		"204 without error is not tolerated": {response(http.StatusNoContent), nil, false},
		"400 is not tolerated":               {response(http.StatusBadRequest), errors.New("400 BAD_REQUEST"), false},
		"500 is not tolerated":               {response(http.StatusInternalServerError), errors.New("500 INTERNAL"), false},
		"nil response with error":            {nil, errors.New("transport failure"), false},
		"nil response without error":         {nil, nil, false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, deletenotfound.IsNotFound(context.Background(), tc.resp, tc.err))
		})
	}
}

// deleteCall returns the response and error of an Atlas API delete call.
type deleteCall func(t *testing.T) (*http.Response, error)

// errorServer returns the URL of a server that answers every request with the given status and Atlas error code.
func errorServer(t *testing.T, status int, errorCode string) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = fmt.Fprintf(w, `{"error":%d,"errorCode":%q}`, status, errorCode)
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

// sdkDelete makes a real SDK call so the error has the same type and body as in production.
func sdkDelete(status int, errorCode string) deleteCall {
	return func(t *testing.T) (*http.Response, error) {
		t.Helper()
		client, err := admin.NewClient(admin.UseBaseURL(errorServer(t, status, errorCode)))
		require.NoError(t, err)
		return client.ProgrammaticAPIKeysAPI.DeleteOrgApiKey(t.Context(), "orgID", "apiKeyID").Execute()
	}
}

// olderSDKDelete is sdkDelete with an older SDK version, whose error type is different.
func olderSDKDelete(status int, errorCode string) deleteCall {
	return func(t *testing.T) (*http.Response, error) {
		t.Helper()
		client, err := adminOlderSDK.NewClient(adminOlderSDK.UseBaseURL(errorServer(t, status, errorCode)))
		require.NoError(t, err)
		_, resp, err := client.ProgrammaticAPIKeysApi.DeleteApiKey(t.Context(), "orgID", "apiKeyID").Execute()
		return resp, err
	}
}

func wrapped(call deleteCall) deleteCall {
	return func(t *testing.T) (*http.Response, error) {
		t.Helper()
		resp, err := call(t)
		return resp, fmt.Errorf("error deleting: %w", err)
	}
}

func withoutResponse(call deleteCall) deleteCall {
	return func(t *testing.T) (*http.Response, error) {
		t.Helper()
		_, err := call(t)
		return nil, err
	}
}

func plain(resp *http.Response, err error) deleteCall {
	return func(*testing.T) (*http.Response, error) {
		return resp, err
	}
}

func TestIsNotFoundOrErrorCode(t *testing.T) {
	const errorCode = "API_KEY_NOT_FOUND"
	testCases := map[string]struct {
		call     deleteCall
		expected bool
	}{
		"404 is tolerated":                                    {sdkDelete(http.StatusNotFound, "RESOURCE_NOT_FOUND"), true},
		"400 with matching error code is tolerated":           {sdkDelete(http.StatusBadRequest, errorCode), true},
		"matching error code from an older SDK is tolerated":  {olderSDKDelete(http.StatusBadRequest, errorCode), true},
		"wrapped error with matching error code is tolerated": {wrapped(sdkDelete(http.StatusBadRequest, errorCode)), true},
		"nil response with matching error code is tolerated":  {withoutResponse(sdkDelete(http.StatusBadRequest, errorCode)), true},
		"400 with other error code is not tolerated":          {sdkDelete(http.StatusBadRequest, "INVALID_ATTRIBUTE"), false},
		"older SDK with other error code is not tolerated":    {olderSDKDelete(http.StatusBadRequest, "INVALID_ATTRIBUTE"), false},
		"500 is not tolerated":                                {sdkDelete(http.StatusInternalServerError, "UNEXPECTED_ERROR"), false},
		"400 with non-API error is not tolerated":             {plain(response(http.StatusBadRequest), errors.New("400 BAD_REQUEST")), false},
		"204 without error is not tolerated":                  {plain(response(http.StatusNoContent), nil), false},
		"nil response with transport error is not tolerated":  {plain(nil, errors.New("transport failure")), false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			resp, err := tc.call(t)
			assert.Equal(t, tc.expected, deletenotfound.IsNotFoundOrErrorCode(t.Context(), resp, err, errorCode))
		})
	}
}

func TestIsNotFoundLegacy(t *testing.T) {
	testCases := map[string]struct {
		resp     *mongodbatlas.Response
		err      error
		expected bool
	}{
		"404 with error is tolerated":   {&mongodbatlas.Response{Response: response(http.StatusNotFound)}, errors.New("404 USER_NOT_FOUND"), true},
		"400 is not tolerated":          {&mongodbatlas.Response{Response: response(http.StatusBadRequest)}, errors.New("400 BAD_REQUEST"), false},
		"nil wrapper with error":        {nil, errors.New("transport failure"), false},
		"nil wrapper without error":     {nil, nil, false},
		"nil inner response with error": {&mongodbatlas.Response{}, errors.New("error without response"), false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, deletenotfound.IsNotFoundLegacy(context.Background(), tc.resp, tc.err))
		})
	}
}
