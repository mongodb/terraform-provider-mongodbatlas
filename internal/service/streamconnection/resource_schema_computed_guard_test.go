package streamconnection_test

import (
	"context"
	"slices"
	"strings"
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
			if a.IsComputed() && !slices.Contains(computedAttributesHandledByAliasOnlyTransition, path) && !keepsStateWhenUnknown(ctx, t, a) {
				t.Errorf("computed attribute %q has no UseStateForUnknown and is not handled by IsAliasOnlyTransition: "+
					"see the comment on computedAttributesHandledByAliasOnlyTransition for what to do", path)
			}
			if nested, ok := nestedAttributes(a); ok {
				walk(path+".", nested)
			}
		}
	}
	walk("", resp.Schema.Attributes)

	// Keep the allowlist honest: every entry must still exist as a Computed attribute without UseStateForUnknown.
	for _, path := range computedAttributesHandledByAliasOnlyTransition {
		if !isUnstableComputed(ctx, t, resp.Schema.Attributes, path) {
			t.Errorf("%q is in computedAttributesHandledByAliasOnlyTransition but is no longer a Computed attribute "+
				"without UseStateForUnknown: remove it from the list and its handling from IsAliasOnlyTransition", path)
		}
	}
}

// keepsStateWhenUnknown runs the attribute's plan modifiers with an unknown plan value and a known state value.
// It fails the test for an attribute kind it does not model, so a new kind never passes the guard unchecked.
func keepsStateWhenUnknown(ctx context.Context, t *testing.T, a schema.Attribute) bool {
	t.Helper()
	state, plan := nonNullState(), nonNullPlan()
	switch attribute := a.(type) {
	case schema.StringAttribute:
		stateValue := types.StringValue("state")
		return anyModifierKeepsState(attribute.PlanModifiers,
			planmodifier.StringRequest{State: state, Plan: plan, StateValue: stateValue, PlanValue: types.StringUnknown(), ConfigValue: types.StringNull()},
			func(r planmodifier.StringRequest) *planmodifier.StringResponse {
				return &planmodifier.StringResponse{PlanValue: r.PlanValue}
			},
			func(m planmodifier.String, r planmodifier.StringRequest, p *planmodifier.StringResponse) {
				m.PlanModifyString(ctx, r, p)
			},
			func(p *planmodifier.StringResponse) attr.Value { return p.PlanValue }, stateValue)
	case schema.BoolAttribute:
		stateValue := types.BoolValue(true)
		return anyModifierKeepsState(attribute.PlanModifiers,
			planmodifier.BoolRequest{State: state, Plan: plan, StateValue: stateValue, PlanValue: types.BoolUnknown(), ConfigValue: types.BoolNull()},
			func(r planmodifier.BoolRequest) *planmodifier.BoolResponse {
				return &planmodifier.BoolResponse{PlanValue: r.PlanValue}
			},
			func(m planmodifier.Bool, r planmodifier.BoolRequest, p *planmodifier.BoolResponse) {
				m.PlanModifyBool(ctx, r, p)
			},
			func(p *planmodifier.BoolResponse) attr.Value { return p.PlanValue }, stateValue)
	case schema.Int64Attribute:
		stateValue := types.Int64Value(1)
		return anyModifierKeepsState(attribute.PlanModifiers,
			planmodifier.Int64Request{State: state, Plan: plan, StateValue: stateValue, PlanValue: types.Int64Unknown(), ConfigValue: types.Int64Null()},
			func(r planmodifier.Int64Request) *planmodifier.Int64Response {
				return &planmodifier.Int64Response{PlanValue: r.PlanValue}
			},
			func(m planmodifier.Int64, r planmodifier.Int64Request, p *planmodifier.Int64Response) {
				m.PlanModifyInt64(ctx, r, p)
			},
			func(p *planmodifier.Int64Response) attr.Value { return p.PlanValue }, stateValue)
	case schema.ListAttribute:
		stateValue := types.ListValueMust(types.StringType, []attr.Value{})
		return anyModifierKeepsState(attribute.PlanModifiers,
			planmodifier.ListRequest{State: state, Plan: plan, StateValue: stateValue, PlanValue: types.ListUnknown(types.StringType), ConfigValue: types.ListNull(types.StringType)},
			func(r planmodifier.ListRequest) *planmodifier.ListResponse {
				return &planmodifier.ListResponse{PlanValue: r.PlanValue}
			},
			func(m planmodifier.List, r planmodifier.ListRequest, p *planmodifier.ListResponse) {
				m.PlanModifyList(ctx, r, p)
			},
			func(p *planmodifier.ListResponse) attr.Value { return p.PlanValue }, stateValue)
	case schema.MapAttribute:
		stateValue := types.MapValueMust(types.StringType, map[string]attr.Value{})
		return anyModifierKeepsState(attribute.PlanModifiers,
			planmodifier.MapRequest{State: state, Plan: plan, StateValue: stateValue, PlanValue: types.MapUnknown(types.StringType), ConfigValue: types.MapNull(types.StringType)},
			func(r planmodifier.MapRequest) *planmodifier.MapResponse {
				return &planmodifier.MapResponse{PlanValue: r.PlanValue}
			},
			func(m planmodifier.Map, r planmodifier.MapRequest, p *planmodifier.MapResponse) {
				m.PlanModifyMap(ctx, r, p)
			},
			func(p *planmodifier.MapResponse) attr.Value { return p.PlanValue }, stateValue)
	case schema.SingleNestedAttribute:
		stateValue := types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{})
		return anyModifierKeepsState(attribute.PlanModifiers,
			planmodifier.ObjectRequest{State: state, Plan: plan, StateValue: stateValue, PlanValue: types.ObjectUnknown(map[string]attr.Type{}), ConfigValue: types.ObjectNull(map[string]attr.Type{})},
			func(r planmodifier.ObjectRequest) *planmodifier.ObjectResponse {
				return &planmodifier.ObjectResponse{PlanValue: r.PlanValue}
			},
			func(m planmodifier.Object, r planmodifier.ObjectRequest, p *planmodifier.ObjectResponse) {
				m.PlanModifyObject(ctx, r, p)
			},
			func(p *planmodifier.ObjectResponse) attr.Value { return p.PlanValue }, stateValue)
	default:
		t.Fatalf("attribute kind %T is not modeled by keepsStateWhenUnknown: add a case for it", a)
		return false
	}
}

// anyModifierKeepsState reports whether any modifier turns the unknown plan value in req into stateValue.
func anyModifierKeepsState[M, Req, Resp any](mods []M, req Req, newResp func(Req) *Resp, modify func(M, Req, *Resp), planValue func(*Resp) attr.Value, stateValue attr.Value) bool {
	for _, m := range mods {
		resp := newResp(req)
		modify(m, req, resp)
		if planValue(resp).Equal(stateValue) {
			return true
		}
	}
	return false
}

// nestedAttributes returns the attributes nested in a, whatever the nesting kind. Paths through list, set and map
// nesting have no index or key segment, e.g. "parent.child".
func nestedAttributes(a schema.Attribute) (map[string]schema.Attribute, bool) {
	switch attribute := a.(type) {
	case schema.SingleNestedAttribute:
		return attribute.Attributes, true
	case schema.ListNestedAttribute:
		return attribute.NestedObject.Attributes, true
	case schema.SetNestedAttribute:
		return attribute.NestedObject.Attributes, true
	case schema.MapNestedAttribute:
		return attribute.NestedObject.Attributes, true
	}
	return nil, false
}

func isUnstableComputed(ctx context.Context, t *testing.T, attrs map[string]schema.Attribute, path string) bool {
	t.Helper()
	for name, a := range attrs {
		if path == name {
			return a.IsComputed() && !keepsStateWhenUnknown(ctx, t, a)
		}
		if nested, ok := nestedAttributes(a); ok && strings.HasPrefix(path, name+".") {
			return isUnstableComputed(ctx, t, nested, strings.TrimPrefix(path, name+"."))
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
