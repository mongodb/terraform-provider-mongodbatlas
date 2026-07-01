package streamworkspace_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamworkspace"
	"github.com/stretchr/testify/assert"
)

func TestFailoverRegionsWriteOnce(t *testing.T) {
	regionObjType := types.ObjectType{AttrTypes: map[string]attr.Type{
		"cloud_provider": types.StringType,
		"region":         types.StringType,
	}}

	regionObj, _ := types.ObjectValue(regionObjType.AttrTypes, map[string]attr.Value{
		"cloud_provider": types.StringValue("AWS"),
		"region":         types.StringValue("VIRGINIA_USA"),
	})
	listWithRegion, _ := types.ListValue(regionObjType, []attr.Value{regionObj})
	emptyList, _ := types.ListValue(regionObjType, []attr.Value{})

	testCases := map[string]struct {
		stateValue            types.List
		planValue             types.List
		expectRequiresReplace bool
	}{
		"null_state_allows_first_time_set": {
			stateValue:            types.ListNull(regionObjType),
			planValue:             listWithRegion,
			expectRequiresReplace: false,
		},
		"empty_state_allows_first_time_set": {
			stateValue:            emptyList,
			planValue:             listWithRegion,
			expectRequiresReplace: false,
		},
		"no_change_no_replace": {
			stateValue:            listWithRegion,
			planValue:             listWithRegion,
			expectRequiresReplace: false,
		},
		"configured_regions_removed_requires_replace": {
			stateValue:            listWithRegion,
			planValue:             types.ListNull(regionObjType),
			expectRequiresReplace: true,
		},
		"configured_regions_changed_requires_replace": {
			stateValue:            listWithRegion,
			planValue:             emptyList,
			expectRequiresReplace: true,
		},
		"unknown_plan_value_does_not_replace": {
			stateValue:            listWithRegion,
			planValue:             types.ListUnknown(regionObjType),
			expectRequiresReplace: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			modifier := streamworkspace.FailoverRegionsWriteOnce{}
			req := planmodifier.ListRequest{
				StateValue: tc.stateValue,
				PlanValue:  tc.planValue,
			}
			resp := &planmodifier.ListResponse{PlanValue: tc.planValue}
			modifier.PlanModifyList(context.Background(), req, resp)
			assert.Equal(t, tc.expectRequiresReplace, resp.RequiresReplace)
		})
	}
}
