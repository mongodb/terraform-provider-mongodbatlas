package streamworkspace_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamworkspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStreamWorkspaceUpdateReq(t *testing.T) {
	regionAttrTypes := map[string]attr.Type{
		"cloud_provider": types.StringType,
		"region":         types.StringType,
	}
	regionObjType := types.ObjectType{AttrTypes: regionAttrTypes}

	dataProcessRegionObj, _ := types.ObjectValue(regionAttrTypes, map[string]attr.Value{
		"cloud_provider": types.StringValue("AWS"),
		"region":         types.StringValue("VIRGINIA_USA"),
	})
	failoverRegionObj, _ := types.ObjectValue(regionAttrTypes, map[string]attr.Value{
		"cloud_provider": types.StringValue("AWS"),
		"region":         types.StringValue("DUBLIN_IRL"),
	})
	failoverList, _ := types.ListValue(regionObjType, []attr.Value{failoverRegionObj})

	testCases := map[string]struct {
		plan                      streamworkspace.TFModel
		state                     streamworkspace.TFModel
		expectCloudProvider       string
		expectRegion              string
		expectFailoverRegionCount int
	}{
		"no_failover_sends_data_process_region": {
			plan: streamworkspace.TFModel{
				DataProcessRegion: dataProcessRegionObj,
				FailoverRegions:   types.ListNull(regionObjType),
			},
			state: streamworkspace.TFModel{
				DataProcessRegion: dataProcessRegionObj,
				FailoverRegions:   types.ListNull(regionObjType),
			},
			expectCloudProvider:       "AWS",
			expectRegion:              "VIRGINIA_USA",
			expectFailoverRegionCount: 0,
		},
		"with_failover_sends_only_failover_regions": {
			plan: streamworkspace.TFModel{
				DataProcessRegion: dataProcessRegionObj,
				FailoverRegions:   failoverList,
			},
			state: streamworkspace.TFModel{
				DataProcessRegion: dataProcessRegionObj,
				FailoverRegions:   types.ListNull(regionObjType),
			},
			expectCloudProvider:       "",
			expectRegion:              "",
			expectFailoverRegionCount: 1,
		},
		"failover_unchanged_in_state_sends_data_process_region": {
			plan: streamworkspace.TFModel{
				DataProcessRegion: dataProcessRegionObj,
				FailoverRegions:   failoverList,
			},
			state: streamworkspace.TFModel{
				DataProcessRegion: dataProcessRegionObj,
				FailoverRegions:   failoverList,
			},
			expectCloudProvider:       "AWS",
			expectRegion:              "VIRGINIA_USA",
			expectFailoverRegionCount: 0,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			req, diags := streamworkspace.NewStreamWorkspaceUpdateReq(context.Background(), &tc.plan, &tc.state)
			require.False(t, diags.HasError())
			require.NotNil(t, req)

			if tc.expectFailoverRegionCount > 0 {
				require.NotNil(t, req.FailoverRegions)
				assert.Len(t, *req.FailoverRegions, tc.expectFailoverRegionCount)
				assert.Equal(t, "DUBLIN_IRL", (*req.FailoverRegions)[0].Region)
				assert.Nil(t, req.CloudProvider)
				assert.Nil(t, req.Region)
			} else {
				assert.Nil(t, req.FailoverRegions)
				require.NotNil(t, req.CloudProvider)
				require.NotNil(t, req.Region)
				assert.Equal(t, tc.expectCloudProvider, *req.CloudProvider)
				assert.Equal(t, tc.expectRegion, *req.Region)
			}
		})
	}
}
