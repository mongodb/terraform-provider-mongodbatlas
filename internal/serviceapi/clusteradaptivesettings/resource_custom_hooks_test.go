package clusteradaptivesettings_test

import (
	"errors"
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

// TestAdaptiveSettingsPostReadHook verifies that the PostRead hook normalizes
// the state when Atlas omits adaptiveSettingsOverrides after a whole-map reset.
func TestAdaptiveSettingsPostReadHook(t *testing.T) {
	t.Parallel()

	r, _ := configuredResource(t, http.StatusOK, `{"effectiveAdaptiveSettings":{"LOAD_SHEDDING":true}}`)
	hook, ok := r.(autogen.PostReadAPICallHook)
	require.True(t, ok)

	state := &clusteradaptivesettings.TFModel{
		AdaptiveSettingsOverrides: jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`),
	}
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(`{}`)})
	require.NoError(t, result.Err)
	require.Equal(t, jsontypes.NewNormalizedNull(), state.AdaptiveSettingsOverrides)
}

// TestAdaptiveSettingsPostReadHookPreservesError verifies that the hook does
// not clear the state when the read failed.
func TestAdaptiveSettingsPostReadHookPreservesError(t *testing.T) {
	t.Parallel()

	r, _ := configuredResource(t, http.StatusOK, `{"effectiveAdaptiveSettings":{"LOAD_SHEDDING":true}}`)
	hook, ok := r.(autogen.PostReadAPICallHook)
	require.True(t, ok)

	overrides := jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`)
	state := &clusteradaptivesettings.TFModel{
		AdaptiveSettingsOverrides: overrides,
	}
	testErr := errors.New("test error")
	result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Err: testErr})
	require.Equal(t, testErr, result.Err)
	require.Equal(t, overrides, state.AdaptiveSettingsOverrides)
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
		require.Equal(t, apiVersionHeader, req.Header.Get("Accept"))
		require.Equal(t, apiVersionHeader, req.Header.Get("Content-Type"))
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
