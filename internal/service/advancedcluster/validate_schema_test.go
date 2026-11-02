package advancedcluster_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/assert"
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

	tests := map[string]struct {
		planVersion  any
		stateVersion any
		operator     advancedcluster.MajorVersionOperator
		threshold    float64
		expectError  bool
	}{
		// Version cannot be resolved (create without mongo_db_major_version): check is skipped, whatever the threshold.
		"unknown in plan, null state, >= 9.0 is skipped": {planVersion: unknown, stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 9.0},
		"unknown in plan, null state, < 9.0 is skipped":  {planVersion: unknown, stateVersion: null, operator: advancedcluster.Lower, threshold: 9.0},
		"null in plan, null state, >= 9.0 is skipped":    {planVersion: null, stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 9.0},
		"null in plan, null state, < 7.0 is skipped":     {planVersion: null, stateVersion: null, operator: advancedcluster.Lower, threshold: 7.0},
		// Version resolved from plan or state is still validated.
		"plan version fails >= 9.0":         {planVersion: "8.0", stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 9.0, expectError: true},
		"plan version passes >= 8.0":        {planVersion: "8.0", stateVersion: null, operator: advancedcluster.EqualOrHigher, threshold: 8.0},
		"state version used when plan null": {planVersion: null, stateVersion: "7.0", operator: advancedcluster.EqualOrHigher, threshold: 8.0, expectError: true},
		"state version used when unknown":   {planVersion: unknown, stateVersion: "8.0", operator: advancedcluster.EqualOrHigher, threshold: 8.0},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := planmodifier.Int64Request{
				Path:        path.Root("advanced_configuration").AtName("default_max_time_ms"),
				ConfigValue: types.Int64Value(100),
				Plan:        tfsdk.Plan{Schema: versionSchema, Raw: rawVersion(tc.planVersion)},
				State:       tfsdk.State{Schema: versionSchema, Raw: rawVersion(tc.stateVersion)},
			}
			resp := &planmodifier.Int64Response{}
			advancedcluster.PlanMustUseMongoDBVersion(tc.threshold, tc.operator).PlanModifyInt64(context.Background(), req, resp)
			assert.Equal(t, tc.expectError, resp.Diagnostics.HasError(), "diagnostics: %v", resp.Diagnostics)
		})
	}
}
