package streamconnection_test

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamconnection"
)

// computedAttributesHandledByAliasOnlyTransition lists the Computed attributes that have no UseStateForUnknown
// plan modifier and are handled explicitly by IsAliasOnlyTransition.
//
// TestComputedAttributesKeepStateWhenUnknownGuard fails when a Computed attribute is added to the schema
// without UseStateForUnknown and without being listed here. If this guard fails, do ONE of the following:
//
//  1. Preferred: add UseStateForUnknown to the new attribute when its value cannot change on update.
//  2. Otherwise: the attribute is unknown in the plan while state holds the API value, so
//     IsAliasOnlyTransition's reflect.DeepEqual never matches and an alias-only migration
//     (instance_name -> workspace_name) sends a PATCH. Atlas rejects every update on a connection used by
//     a stream processor (STREAM_CONNECTION_HAS_STREAM_PROCESSORS, CLOUDP-446873). Handle the attribute in
//     IsAliasOnlyTransition (see withStateForUnknownAzureRegion), add a case to
//     TestIsAliasOnlyTransition with the attribute unknown in the plan, then list its path here.
var computedAttributesHandledByAliasOnlyTransition = []string{
	"azure.region",
}

func TestComputedAttributesKeepStateWhenUnknownGuard(t *testing.T) {
	ctx := t.Context()
	var resp resource.SchemaResponse
	streamconnection.Resource().Schema(ctx, resource.SchemaRequest{}, &resp)

	var walk func(prefix string, attrs map[string]schema.Attribute)
	walk = func(prefix string, attrs map[string]schema.Attribute) {
		for name, a := range attrs {
			path := prefix + name
			if a.IsComputed() && !slices.Contains(computedAttributesHandledByAliasOnlyTransition, path) && !keepsStateWhenUnknown(ctx, a) {
				t.Errorf("computed attribute %q has no UseStateForUnknown and is not handled by IsAliasOnlyTransition: "+
					"see the comment on computedAttributesHandledByAliasOnlyTransition for what to do", path)
			}
			if nested, ok := a.(schema.SingleNestedAttribute); ok {
				walk(path+".", nested.Attributes)
			}
		}
	}
	walk("", resp.Schema.Attributes)

	// Keep the allowlist honest: every entry must still exist as a Computed attribute without UseStateForUnknown.
	for _, path := range computedAttributesHandledByAliasOnlyTransition {
		if !isUnstableComputed(ctx, resp.Schema.Attributes, path) {
			t.Errorf("%q is in computedAttributesHandledByAliasOnlyTransition but is no longer a Computed attribute "+
				"without UseStateForUnknown: remove it from the list and its handling from IsAliasOnlyTransition", path)
		}
	}
}

// keepsStateWhenUnknown runs the attribute's plan modifiers with an unknown plan value and a known state value.
func keepsStateWhenUnknown(ctx context.Context, a schema.Attribute) bool {
	switch attribute := a.(type) {
	case schema.StringAttribute:
		for _, m := range attribute.PlanModifiers {
			req := planmodifier.StringRequest{State: nonNullState(), Plan: nonNullPlan(), StateValue: types.StringValue("state"), PlanValue: types.StringUnknown(), ConfigValue: types.StringNull()}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
			m.PlanModifyString(ctx, req, resp)
			if resp.PlanValue.Equal(req.StateValue) {
				return true
			}
		}
	case schema.SingleNestedAttribute:
		for _, m := range attribute.PlanModifiers {
			state := types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{})
			req := planmodifier.ObjectRequest{State: nonNullState(), Plan: nonNullPlan(), StateValue: state, PlanValue: types.ObjectUnknown(map[string]attr.Type{}), ConfigValue: types.ObjectNull(map[string]attr.Type{})}
			resp := &planmodifier.ObjectResponse{PlanValue: req.PlanValue}
			m.PlanModifyObject(ctx, req, resp)
			if resp.PlanValue.Equal(req.StateValue) {
				return true
			}
		}
	}
	return false
}

func isUnstableComputed(ctx context.Context, attrs map[string]schema.Attribute, path string) bool {
	for name, a := range attrs {
		if path == name {
			return a.IsComputed() && !keepsStateWhenUnknown(ctx, a)
		}
		if nested, ok := a.(schema.SingleNestedAttribute); ok && len(path) > len(name)+1 && path[:len(name)+1] == name+"." {
			return isUnstableComputed(ctx, nested.Attributes, path[len(name)+1:])
		}
	}
	return false
}

// UseStateForUnknown does nothing when the resource has no state or plan, so the requests need non-null raw values.
func nonNullState() tfsdk.State {
	return tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})}
}

func nonNullPlan() tfsdk.Plan {
	return tfsdk.Plan{Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{})}
}
