package advancedcluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
)

const (
	errDefaultMaxTimeMS    = "`advanced_configuration.default_max_time_ms` can only be configured if the mongo_db_major_version is 8.0 or higher"
	errChangeStreamOptions = "`advanced_configuration.change_stream_options_pre_and_post_images_expire_after_seconds` can only be configured if the mongo_db_major_version is 7.0 or higher"
)

func TestPlanMustUseMongoDBVersion(t *testing.T) {
	versionSchema := schema.Schema{Attributes: map[string]schema.Attribute{
		"mongo_db_major_version": schema.StringAttribute{Optional: true, Computed: true},
	}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"mongo_db_major_version": tftypes.String}}
	rawVersion := func(v any) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{"mongo_db_major_version": tftypes.NewValue(tftypes.String, v)})
	}
	unknown := tftypes.UnknownValue
	var null any // nil -> null value

	// planVersion mirrors what Terraform plans for an Optional+Computed attribute: the config value when known,
	// unknown when the config value is unknown, and unknown or the prior state when it is not configured.
	tests := map[string]struct {
		configVersion any
		planVersion   any
		stateVersion  any
		operator      advancedcluster.MajorVersionOperator
		threshold     float64
		expectError   bool
	}{
		// Version cannot be resolved (create without mongo_db_major_version): check is skipped, whatever the threshold.
		"null config, null state, >= 9.0 is skipped": {configVersion: null, planVersion: unknown, stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 9.0},
		"null config, null state, < 9.0 is skipped":  {configVersion: null, planVersion: unknown, stateVersion: null, operator: advancedcluster.Lower, threshold: 9.0},
		"null config, null state, < 7.0 is skipped":  {configVersion: null, planVersion: unknown, stateVersion: null, operator: advancedcluster.Lower, threshold: 7.0},
		// Set through a reference unknown at plan time: state may be stale, so the check is skipped.
		"unknown config, state 7.0, >= 8.0 is skipped": {configVersion: unknown, planVersion: unknown, stateVersion: "7.0", operator: advancedcluster.EqualOrHigher, threshold: 8.0},
		"unknown config, state 8.0, < 8.0 is skipped":  {configVersion: unknown, planVersion: unknown, stateVersion: "8.0", operator: advancedcluster.Lower, threshold: 8.0},
		"unknown config, null state is skipped":        {configVersion: unknown, planVersion: unknown, stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 9.0},
		// Version resolved from config or (when not configured) state is still validated.
		"config version fails >= 9.0":         {configVersion: "8.0", planVersion: "8.0", stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 9.0, expectError: true},
		"config version passes >= 8.0":        {configVersion: "8.0", planVersion: "8.0", stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 8.0},
		"config takes precedence over state":  {configVersion: "8.0", planVersion: "8.0", stateVersion: "7.0", operator: advancedcluster.EqualOrHigher, threshold: 8.0},
		"state used when config null, fails":  {configVersion: null, planVersion: unknown, stateVersion: "7.0", operator: advancedcluster.EqualOrHigher, threshold: 8.0, expectError: true},
		"state used when config null, passes": {configVersion: null, planVersion: "8.0", stateVersion: "8.0", operator: advancedcluster.EqualOrHigher, threshold: 8.0},
		// Plan and config differ only to assert which one is read: config is the source, plan is ignored.
		"config is read, not plan": {configVersion: unknown, planVersion: "7.0", stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 8.0},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := planmodifier.Int64Request{
				Path:        path.Root("advanced_configuration").AtName("default_max_time_ms"),
				ConfigValue: types.Int64Value(100),
				PlanValue:   types.Int64Value(100),
				Config:      tfsdk.Config{Schema: versionSchema, Raw: rawVersion(tc.configVersion)},
				Plan:        tfsdk.Plan{Schema: versionSchema, Raw: rawVersion(tc.planVersion)},
				State:       tfsdk.State{Schema: versionSchema, Raw: rawVersion(tc.stateVersion)},
			}
			resp := &planmodifier.Int64Response{}
			advancedcluster.PlanMustUseMongoDBVersion(tc.threshold, tc.operator).PlanModifyInt64(context.Background(), req, resp)
			assert.Equal(t, tc.expectError, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics)
		})
	}
}

func TestValidateMongoDBVersionAtApply(t *testing.T) {
	advConfigAttrTypes := map[string]tftypes.Type{
		"change_stream_options_pre_and_post_images_expire_after_seconds": tftypes.Number,
		"default_max_time_ms": tftypes.Number,
	}
	configSchema := schema.Schema{Attributes: map[string]schema.Attribute{
		"advanced_configuration": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
			"change_stream_options_pre_and_post_images_expire_after_seconds": schema.Int64Attribute{Optional: true},
			"default_max_time_ms": schema.Int64Attribute{Optional: true},
		}},
	}}
	advConfigType := tftypes.Object{AttributeTypes: advConfigAttrTypes}
	rootType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"advanced_configuration": advConfigType}}
	// rawConfig builds the config; a nil map means advanced_configuration is not configured.
	rawConfig := func(advConfig map[string]any) tftypes.Value {
		advConfigValue := tftypes.NewValue(advConfigType, nil)
		if advConfig != nil {
			values := map[string]tftypes.Value{}
			for name := range advConfigAttrTypes {
				values[name] = tftypes.NewValue(tftypes.Number, advConfig[name])
			}
			advConfigValue = tftypes.NewValue(advConfigType, values)
		}
		return tftypes.NewValue(rootType, map[string]tftypes.Value{"advanced_configuration": advConfigValue})
	}
	defaultMaxTimeMS := map[string]any{"default_max_time_ms": 100}
	changeStreamOptions := map[string]any{"change_stream_options_pre_and_post_images_expire_after_seconds": 100}
	bothGated := map[string]any{"default_max_time_ms": 100, "change_stream_options_pre_and_post_images_expire_after_seconds": 100}

	tests := map[string]struct {
		advConfig      map[string]any
		plannedVersion types.String
		stateVersion   types.String
		expectErrors   []string
	}{
		// Version resolved at apply, e.g. from a reference that was unknown at plan time.
		"resolved version fails gate":                  {advConfig: defaultMaxTimeMS, plannedVersion: types.StringValue("7.0"), stateVersion: types.StringNull(), expectErrors: []string{errDefaultMaxTimeMS}},
		"resolved version passes gate":                 {advConfig: defaultMaxTimeMS, plannedVersion: types.StringValue("8.0"), stateVersion: types.StringNull()},
		"resolved version fails every gate":            {advConfig: bothGated, plannedVersion: types.StringValue("6.0"), stateVersion: types.StringNull(), expectErrors: []string{errChangeStreamOptions, errDefaultMaxTimeMS}},
		"only configured gated attributes are checked": {advConfig: changeStreamOptions, plannedVersion: types.StringValue("7.0"), stateVersion: types.StringNull()},
		"planned version takes precedence over state":  {advConfig: defaultMaxTimeMS, plannedVersion: types.StringValue("8.0"), stateVersion: types.StringValue("7.0")},
		// Update without mongo_db_major_version in config: the planned value can be unknown, state is used.
		"unknown planned version uses state, fails":  {advConfig: defaultMaxTimeMS, plannedVersion: types.StringUnknown(), stateVersion: types.StringValue("7.0"), expectErrors: []string{errDefaultMaxTimeMS}},
		"unknown planned version uses state, passes": {advConfig: defaultMaxTimeMS, plannedVersion: types.StringUnknown(), stateVersion: types.StringValue("8.0")},
		// Create without mongo_db_major_version: Atlas picks the default, so the check is skipped.
		"no version is skipped":                 {advConfig: bothGated, plannedVersion: types.StringUnknown(), stateVersion: types.StringNull()},
		"null version is skipped":               {advConfig: bothGated, plannedVersion: types.StringNull(), stateVersion: types.StringNull()},
		"advanced_configuration not configured": {advConfig: nil, plannedVersion: types.StringValue("6.0"), stateVersion: types.StringNull()},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			config := tfsdk.Config{Schema: configSchema, Raw: rawConfig(tc.advConfig)}
			diags := &diag.Diagnostics{}
			advancedcluster.ValidateMongoDBVersionAtApply(context.Background(), diags, &config, tc.plannedVersion, tc.stateVersion)
			var gotErrors []string
			for _, d := range diags.Errors() {
				gotErrors = append(gotErrors, d.Summary())
			}
			assert.Equal(t, tc.expectErrors, gotErrors)
		})
	}
}
