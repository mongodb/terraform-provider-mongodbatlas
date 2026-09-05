package clusteradaptivesettings_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/clusteradaptivesettings"
)

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
