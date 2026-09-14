package clusteradaptivesettings_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/clusteradaptivesettings"
)

// TestAdaptiveSettingsPostReadHook verifies the hook normalizes state when
// Atlas omits adaptiveSettingsOverrides after a whole-map reset, and preserves
// state on read errors.
func TestAdaptiveSettingsPostReadHook(t *testing.T) {
	t.Parallel()

	t.Run("clears stale overrides on success", func(t *testing.T) {
		t.Parallel()
		r, _ := configuredResource(t, http.StatusOK, `{}`)
		hook := r.(autogen.PostReadAPICallHook)
		state := &clusteradaptivesettings.TFModel{
			AdaptiveSettingsOverrides: jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`),
		}
		result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Body: []byte(`{}`)})
		require.NoError(t, result.Err)
		require.Equal(t, jsontypes.NewNormalizedNull(), state.AdaptiveSettingsOverrides)
	})

	t.Run("preserves state on error", func(t *testing.T) {
		t.Parallel()
		r, _ := configuredResource(t, http.StatusOK, `{}`)
		hook := r.(autogen.PostReadAPICallHook)
		overrides := jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`)
		state := &clusteradaptivesettings.TFModel{
			AdaptiveSettingsOverrides: overrides,
		}
		testErr := errors.New("test error")
		result := hook.PostReadAPICall(autogen.HandleReadReq{State: state}, autogen.APICallResult{Err: testErr})
		require.Equal(t, testErr, result.Err)
		require.Equal(t, overrides, state.AdaptiveSettingsOverrides)
	})
}

// TestAdaptiveSettingsReadOverrides verifies read behavior across the
// replacement PATCH contract: omitted field, explicit null, empty map, and
// populated overrides.
func TestAdaptiveSettingsReadOverrides(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		response string
		expected jsontypes.Normalized
	}{
		"externally reset": {
			response: `{"effectiveAdaptiveSettings":{"LOAD_SHEDDING":false}}`,
			expected: jsontypes.NewNormalizedNull(),
		},
		"explicit null": {
			response: `{"adaptiveSettingsOverrides":null,"effectiveAdaptiveSettings":{"LOAD_SHEDDING":false}}`,
			expected: jsontypes.NewNormalizedNull(),
		},
		"empty overrides": {
			response: `{"adaptiveSettingsOverrides":{},"effectiveAdaptiveSettings":{"LOAD_SHEDDING":false}}`,
			expected: jsontypes.NewNormalizedValue(`{}`),
		},
		"changed overrides": {
			response: `{"adaptiveSettingsOverrides":{"LOAD_SHEDDING":false},"effectiveAdaptiveSettings":{"LOAD_SHEDDING":false}}`,
			expected: jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":false}`),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, calls := configuredResource(t, http.StatusOK, test.response)
			state := adaptiveSettingsState(t, r, jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`))
			resp := resource.ReadResponse{State: state}
			r.Read(t.Context(), resource.ReadRequest{State: state}, &resp)
			require.False(t, resp.Diagnostics.HasError(), resp.Diagnostics)
			var actual clusteradaptivesettings.TFModel
			diags := resp.State.Get(t.Context(), &actual)
			require.False(t, diags.HasError(), diags)
			require.Equal(t, test.expected, actual.AdaptiveSettingsOverrides)
			require.JSONEq(t, `{"LOAD_SHEDDING":false}`, actual.EffectiveAdaptiveSettings.ValueString())
			require.EqualValues(t, 1, calls.Load())
		})
	}
}

// TestAdaptiveSettingsValidateConfig verifies the JSON-object and null-entry
// validation rules for adaptive_settings_overrides.
func TestAdaptiveSettingsValidateConfig(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value   jsontypes.Normalized
		invalid bool
	}{
		"omitted":        {value: jsontypes.NewNormalizedNull()},
		"unknown":        {value: jsontypes.NewNormalizedUnknown()},
		"empty object":   {value: jsontypes.NewNormalizedValue(`{}`)},
		"configured":     {value: jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":false}`)},
		"future setting": {value: jsontypes.NewNormalizedValue(`{"future":{"enabled":true,"limit":9007199254740993}}`)},
		"null entry":     {value: jsontypes.NewNormalizedValue(`{"SEARCH_LOAD_SHEDDING": null}`), invalid: true},
		"JSON null":      {value: jsontypes.NewNormalizedValue(`null`), invalid: true},
		"array":          {value: jsontypes.NewNormalizedValue(`[]`), invalid: true},
		"scalar":         {value: jsontypes.NewNormalizedValue(`false`), invalid: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := config.AnalyticsResourceFunc(clusteradaptivesettings.Resource())()
			state := adaptiveSettingsState(t, r, test.value)
			attribute := state.Schema.(schema.Schema).Attributes["adaptive_settings_overrides"].(schema.StringAttribute)
			require.NotEmpty(t, attribute.Validators)
			resp := validator.StringResponse{}
			for _, validation := range attribute.Validators {
				validation.ValidateString(t.Context(), validator.StringRequest{
					Path:        path.Root("adaptive_settings_overrides"),
					ConfigValue: test.value.StringValue,
					Config:      tfsdk.Config(state),
				}, &resp)
			}
			require.Equal(t, test.invalid, resp.Diagnostics.HasError(), resp.Diagnostics)
			if test.invalid {
				require.Contains(t, resp.Diagnostics[0].Summary(), "Invalid Adaptive Settings")
			}
		})
	}
}

func adaptiveSettingsState(t *testing.T, r resource.Resource, overrides jsontypes.Normalized) tfsdk.State {
	t.Helper()
	var schemaResp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &schemaResp)
	require.False(t, schemaResp.Diagnostics.HasError(), schemaResp.Diagnostics)
	state := tfsdk.State{Schema: schemaResp.Schema}
	diags := state.Set(t.Context(), &clusteradaptivesettings.TFModel{
		ProjectId:                 types.StringValue("projectID"),
		ClusterName:               types.StringValue("clusterName"),
		AdaptiveSettingsOverrides: overrides,
		EffectiveAdaptiveSettings: jsontypes.NewNormalizedValue(`{"LOAD_SHEDDING":true}`),
	})
	require.False(t, diags.HasError(), diags)
	return state
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
