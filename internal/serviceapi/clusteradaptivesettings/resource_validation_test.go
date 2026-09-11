package clusteradaptivesettings_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/serviceapi/clusteradaptivesettings"
)

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
		"null entry":     {value: jsontypes.NewNormalizedValue(`{"SEARCH_OVERLOAD_PROTECTION": null}`), invalid: true},
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
