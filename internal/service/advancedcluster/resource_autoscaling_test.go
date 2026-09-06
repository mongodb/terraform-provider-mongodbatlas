package advancedcluster_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
	"go.mongodb.org/atlas-sdk/v20250312024/mockadmin"
)

func TestAutoScalingRequestWithoutStorage(t *testing.T) {
	testCases := map[string]struct {
		edition  string
		settings map[string]any
		expected string
	}{
		"Infinite compute only":              {"INFINITE", map[string]any{"compute_enabled": true, "compute_max_instance_size": "M20"}, `{"compute":{"enabled":true,"maxInstanceSize":"M20"}}`},
		"Infinite disabled compute":          {"INFINITE", map[string]any{"compute_enabled": false}, `{"compute":{"enabled":false}}`},
		"Infinite empty auto scaling":        {"INFINITE", map[string]any{}, `{}`},
		"Infinite omitted auto scaling":      {"INFINITE", nil, `null`},
		"Infinite disk true reaches server":  {"INFINITE", map[string]any{"disk_gb_enabled": true}, `{"diskGB":{"enabled":true}}`},
		"Infinite disk false reaches server": {"INFINITE", map[string]any{"disk_gb_enabled": false}, `{"diskGB":{"enabled":false}}`},
		"CORE keeps empty disk":              {"CORE", map[string]any{"compute_enabled": true, "compute_max_instance_size": "M20"}, `{"compute":{"enabled":true,"maxInstanceSize":"M20"},"diskGB":{}}`},
		"default edition keeps empty disk":   {"", map[string]any{"compute_enabled": true, "compute_max_instance_size": "M20"}, `{"compute":{"enabled":true,"maxInstanceSize":"M20"},"diskGB":{}}`},
	}
	for name, tc := range testCases {
		for _, autoScalingKey := range []string{"auto_scaling", "analytics_auto_scaling"} {
			for _, operation := range []string{"create", "update"} {
				t.Run(name+"/"+autoScalingKey+"/"+operation, func(t *testing.T) {
					ctx := t.Context()
					r := advancedcluster.Resource()
					var schemaResp resource.SchemaResponse
					r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
					typ := schemaResp.Schema.Type().TerraformType(ctx)
					model := func(settings map[string]any, instanceSize string) tfsdk.Plan {
						var autoScaling, edition any
						if settings != nil {
							autoScaling = settings
						}
						if tc.edition != "" {
							edition = tc.edition
						}
						regions := []any{}
						for _, region := range []string{"US_EAST_1", "US_WEST_2"} {
							regions = append(regions, map[string]any{
								"provider_name": "AWS", "region_name": region, "priority": int64(7),
								"electable_specs": map[string]any{"instance_size": instanceSize, "node_count": int64(2)},
								autoScalingKey:    autoScaling,
							})
						}
						return tfsdk.Plan{Schema: schemaResp.Schema, Raw: planTestValue(typ, map[string]any{
							"name": "example", "project_id": dummyProjectID, "cluster_type": "REPLICASET", "database_edition": edition,
							"replication_specs": []any{map[string]any{"region_configs": regions}},
						})}
					}
					plan := model(tc.settings, "M10")
					api := mockadmin.NewClustersAPI(t)
					r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
					checkPayload := func(args mock.Arguments) {
						payload := args[len(args)-1].(*admin.ClusterDescription20240805)
						regions := payload.GetReplicationSpecs()[0].GetRegionConfigs()
						require.Len(t, regions, 2)
						for _, region := range regions {
							autoScaling := region.AutoScaling
							if autoScalingKey == "analytics_auto_scaling" {
								autoScaling = region.AnalyticsAutoScaling
							}
							encoded, err := json.Marshal(autoScaling)
							require.NoError(t, err)
							require.JSONEq(t, tc.expected, string(encoded))
						}
					}
					// Stop at the API boundary after inspecting the outgoing request, without polling or provisioning.
					apiError := errors.New("request inspected")
					if operation == "create" {
						api.On("CreateCluster", mock.Anything, dummyProjectID, mock.Anything).Run(checkPayload).
							Return(admin.CreateClusterApiRequest{ApiService: api}).Once()
						api.EXPECT().CreateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
						var resp resource.CreateResponse
						r.Create(ctx, resource.CreateRequest{Plan: plan}, &resp)
						require.Len(t, resp.Diagnostics.Errors(), 1)
						require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), apiError.Error())
					} else {
						api.On("UpdateCluster", mock.Anything, dummyProjectID, "example", mock.Anything).Run(checkPayload).
							Return(admin.UpdateClusterApiRequest{ApiService: api}).Once()
						api.EXPECT().UpdateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
						prior := model(tc.settings, "M20")
						var resp resource.UpdateResponse
						r.Update(ctx, resource.UpdateRequest{Plan: plan, State: tfsdk.State{Schema: schemaResp.Schema, Raw: prior.Raw}}, &resp)
						require.Len(t, resp.Diagnostics.Errors(), 1)
						require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), apiError.Error())
					}
				})
			}
		}
	}
}
