package clusteradaptivesettings_test

import (
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/clusteradaptivesettings"
)

func TestAdaptiveSettingsRequestHooks(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request      string
		expected     string
		expectedGETs int32
	}{
		"omitted overrides reset the whole map": {
			request:  `{}`,
			expected: `{"adaptiveSettingsOverrides":null}`,
		},
		"empty overrides reset every effective key": {
			request:      `{"adaptiveSettingsOverrides":{}}`,
			expected:     `{"adaptiveSettingsOverrides":{"key1":null,"key2":null}}`,
			expectedGETs: 2,
		},
		"planned overrides replace effective candidates": {
			request:      `{"adaptiveSettingsOverrides":{"key2":"val22"}}`,
			expected:     `{"adaptiveSettingsOverrides":{"key1":null,"key2":"val22"}}`,
			expectedGETs: 2,
		},
		"new planned keys are preserved": {
			request:      `{"adaptiveSettingsOverrides":{"key3":{"enabled":true}}}`,
			expected:     `{"adaptiveSettingsOverrides":{"key1":null,"key2":null,"key3":{"enabled":true}}}`,
			expectedGETs: 2,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resourceInstance, getCalls := configuredResource(t, http.StatusOK, `{"effectiveAdaptiveSettings":{"key1":true,"key2":false}}`)
			createHook, ok := resourceInstance.(autogen.PreCreateAPICallHook)
			require.True(t, ok)
			updateHook, ok := resourceInstance.(autogen.PreUpdateAPICallHook)
			require.True(t, ok)

			callParams := adaptiveSettingsCallParams()
			createParams, createBody, err := createHook.PreCreateAPICall(t.Context(), callParams, []byte(testCase.request))
			require.NoError(t, err)
			require.Equal(t, callParams, createParams)
			require.JSONEq(t, testCase.expected, string(createBody))

			updateParams, updateBody, err := updateHook.PreUpdateAPICall(t.Context(), callParams, []byte(testCase.request))
			require.NoError(t, err)
			require.Equal(t, callParams, updateParams)
			require.JSONEq(t, testCase.expected, string(updateBody))
			require.Equal(t, testCase.expectedGETs, getCalls.Load())
		})
	}
}

func TestRemovedAdaptiveSettingsOverridesResetTheWholeMap(t *testing.T) {
	t.Parallel()

	plan := clusteradaptivesettings.TFModel{
		AdaptiveSettingsOverrides: jsontypes.NewNormalizedNull(),
	}
	bodyReq, err := autogen.Marshal(&plan, true)
	require.NoError(t, err)
	require.JSONEq(t, `{"adaptiveSettingsOverrides":null}`, string(bodyReq))

	resourceInstance, getCalls := configuredResource(t, http.StatusOK, `{"effectiveAdaptiveSettings":{"key1":true,"key2":false}}`)
	hook, ok := resourceInstance.(autogen.PreUpdateAPICallHook)
	require.True(t, ok)
	_, updatedBody, err := hook.PreUpdateAPICall(t.Context(), adaptiveSettingsCallParams(), bodyReq)
	require.NoError(t, err)
	require.JSONEq(t, `{"adaptiveSettingsOverrides":null}`, string(updatedBody))
	require.Zero(t, getCalls.Load())
}

func TestAdaptiveSettingsRequestHookErrors(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		responseBody string
		errorText    string
		statusCode   int
	}{
		"GET error": {
			statusCode:   http.StatusInternalServerError,
			responseBody: `{"error":500}`,
			errorText:    "get Adaptive Settings before PATCH",
		},
		"missing effective settings": {
			statusCode:   http.StatusOK,
			responseBody: `{}`,
			errorText:    "missing effectiveAdaptiveSettings",
		},
		"null effective settings": {
			statusCode:   http.StatusOK,
			responseBody: `{"effectiveAdaptiveSettings":null}`,
			errorText:    "null effectiveAdaptiveSettings",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resourceInstance, _ := configuredResource(t, testCase.statusCode, testCase.responseBody)
			hook, ok := resourceInstance.(autogen.PreUpdateAPICallHook)
			require.True(t, ok)
			_, _, err := hook.PreUpdateAPICall(t.Context(), adaptiveSettingsCallParams(), []byte(`{"adaptiveSettingsOverrides":{"key1":true}}`))
			require.ErrorContains(t, err, testCase.errorText)
		})
	}
}

func configuredResource(t *testing.T, statusCode int, responseBody string) (resource.Resource, *atomic.Int32) {
	t.Helper()
	client, err := config.NewClient(&config.Credentials{BaseURL: "http://atlas.example.test"}, "")
	require.NoError(t, err)
	var getCalls atomic.Int32
	client.AtlasV2.GetConfig().HTTPClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		getCalls.Add(1)
		if req.Method != http.MethodGet || req.URL.Path != "/api/atlas/v2/groups/projectID/clusters/clusterName/adaptiveSettings" {
			t.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
		return &http.Response{
			StatusCode: statusCode,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(responseBody)),
			Request:    req,
		}, nil
	})
	resourceInstance := clusteradaptivesettings.Resource()
	clientSetter, ok := resourceInstance.(interface{ SetClient(*config.MongoDBClient) })
	require.True(t, ok)
	clientSetter.SetClient(client)
	return resourceInstance, &getCalls
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func adaptiveSettingsCallParams() config.APICallParams {
	return config.APICallParams{
		VersionHeader: "application/vnd.atlas.preview+json",
		RelativePath:  "/api/atlas/v2/groups/{projectId}/clusters/{clusterName}/adaptiveSettings",
		PathParams: map[string]string{
			"projectId":   "projectID",
			"clusterName": "clusterName",
		},
		Method: http.MethodPatch,
	}
}
