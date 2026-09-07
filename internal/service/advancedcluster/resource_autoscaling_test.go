package advancedcluster_test

import (
	"encoding/json"
	"errors"
	"maps"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
	"go.mongodb.org/atlas-sdk/v20250312024/mockadmin"
)

func TestUpdateRemovesShardSizeLimit(t *testing.T) {
	compute := map[string]any{
		"compute_enabled": true, "compute_scale_down_enabled": true,
		"compute_min_instance_size": "M10", "compute_max_instance_size": "M30",
	}
	hardware := map[string]any{"instance_size": "M10", "node_count": int64(2)}
	zeroNodes := map[string]any{"instance_size": "M10", "node_count": int64(0)}
	electableOnly := map[string]any{"electable_specs": hardware}
	withZeroNodes := func(autoScaling any) map[string]any {
		return map[string]any{
			"electable_specs": hardware, "read_only_specs": zeroNodes, "analytics_specs": zeroNodes,
			"auto_scaling": autoScaling,
		}
	}
	testCases := map[string]struct {
		expectedRegions string
		configRegions   []any
		planRegions     []any
		withTags        bool
	}{
		"preserves active planned hardware and analytics scaling when omitted from configuration": {
			configRegions: []any{map[string]any{}},
			planRegions: []any{map[string]any{
				"electable_specs": hardware, "read_only_specs": hardware, "analytics_specs": hardware,
				"analytics_auto_scaling": compute,
			}},
			expectedRegions: `[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"analyticsAutoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M30","minInstanceSize":"M10","scaleDownEnabled":true}},"electableSpecs":{"instanceSize":"M10","nodeCount":2},"readOnlySpecs":{"instanceSize":"M10","nodeCount":2},"analyticsSpecs":{"instanceSize":"M10","nodeCount":2}}]`,
		},
		"preserves planned compute when auto scaling is removed from configuration": {
			configRegions:   []any{electableOnly},
			planRegions:     []any{withZeroNodes(compute)},
			withTags:        true,
			expectedRegions: `[{"autoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M30","minInstanceSize":"M10","scaleDownEnabled":true},"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":2}}]`,
		},
		"preserves planned compute false without copying zero-node hardware": {
			configRegions:   []any{electableOnly},
			planRegions:     []any{withZeroNodes(map[string]any{"compute_enabled": false})},
			expectedRegions: `[{"autoScaling":{"compute":{"enabled":false},"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":2}}]`,
		},
		"omits empty auto scaling and preserves configured compute in every region": {
			configRegions: []any{
				map[string]any{"electable_specs": hardware, "auto_scaling": map[string]any{}},
				map[string]any{"electable_specs": hardware, "auto_scaling": compute},
			},
			expectedRegions: `[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":2}},{"autoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M30","minInstanceSize":"M10","scaleDownEnabled":true},"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":2}}]`,
		},
		"omits computed zero-node specs and preserves unrelated patch fields": {
			configRegions:   []any{electableOnly},
			planRegions:     []any{withZeroNodes(nil)},
			withTags:        true,
			expectedRegions: `[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":2}}]`,
		},
		"preserves explicitly configured zero-node hardware": {
			configRegions:   []any{withZeroNodes(nil)},
			expectedRegions: `[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":2},"readOnlySpecs":{"instanceSize":"M10","nodeCount":0},"analyticsSpecs":{"instanceSize":"M10","nodeCount":0}}]`,
		},
		"preserves zero-node hardware resolved from unknown configuration": {
			configRegions: []any{map[string]any{
				"electable_specs": hardware, "read_only_specs": tftypes.UnknownValue, "analytics_specs": tftypes.UnknownValue,
			}},
			planRegions:     []any{withZeroNodes(nil)},
			expectedRegions: `[{"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}},"electableSpecs":{"instanceSize":"M10","nodeCount":2},"readOnlySpecs":{"instanceSize":"M10","nodeCount":0},"analyticsSpecs":{"instanceSize":"M10","nodeCount":0}}]`,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := t.Context()
			r := advancedcluster.Resource()
			var schemaResp resource.SchemaResponse
			r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
			typ := schemaResp.Schema.Type().TerraformType(ctx)
			model := func(regions []any, withTags bool) tfsdk.Plan {
				attributes := map[string]any{
					"name": "example", "project_id": dummyProjectID, "cluster_type": "REPLICASET", "database_edition": "INFINITE",
					"replication_specs": []any{map[string]any{"region_configs": regions}},
				}
				if withTags {
					attributes["tags"] = map[string]tftypes.Value{"environment": tftypes.NewValue(tftypes.String, "test")}
				}
				return tfsdk.Plan{Schema: schemaResp.Schema, Raw: planTestValue(typ, attributes)}
			}
			planRegions := tc.planRegions
			if planRegions == nil {
				planRegions = tc.configRegions
			}
			stateRegions := make([]any, len(planRegions))
			for i, raw := range planRegions {
				region := maps.Clone(raw.(map[string]any))
				autoScaling, _ := region["auto_scaling"].(map[string]any)
				autoScaling = maps.Clone(autoScaling)
				if autoScaling == nil {
					autoScaling = map[string]any{}
				}
				autoScaling["storage_config"] = map[string]any{"shard_size_limit_gb": int64(1024)}
				region["auto_scaling"] = autoScaling
				stateRegions[i] = region
			}
			api := mockadmin.NewClustersAPI(t)
			r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
			api.On("UpdateCluster", mock.Anything, dummyProjectID, "example", mock.Anything).Run(func(args mock.Arguments) {
				payload := args[3].(*admin.ClusterDescription20240805)
				require.Len(t, payload.GetReplicationSpecs(), 1)
				encoded, err := json.Marshal(payload.GetReplicationSpecs()[0].GetRegionConfigs())
				require.NoError(t, err)
				require.JSONEq(t, tc.expectedRegions, string(encoded))
				if tc.withTags {
					require.Equal(t, []admin.ResourceTag{{Key: "environment", Value: "test"}}, payload.GetTags())
				}
			}).Return(admin.UpdateClusterApiRequest{ApiService: api}).Once()
			// Stop after inspecting the request; polling and Atlas behavior have separate lifecycle coverage.
			apiError := errors.New("request inspected")
			api.EXPECT().UpdateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
			prior := model(stateRegions, false)
			configuration := model(tc.configRegions, tc.withTags)
			var resp resource.UpdateResponse
			r.Update(ctx, resource.UpdateRequest{
				Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: configuration.Raw},
				Plan:   model(planRegions, tc.withTags),
				State:  tfsdk.State{Schema: schemaResp.Schema, Raw: prior.Raw},
			}, &resp)
			require.Len(t, resp.Diagnostics.Errors(), 1)
			require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), apiError.Error())
		})
	}
}

func TestUpdateRemovesShardSizeLimitPreservesPlannedValues(t *testing.T) {
	for _, omittedField := range []string{"node_count", "zone_name"} {
		for _, operation := range []string{"keep storage", "clear storage"} {
			t.Run(omittedField+"/"+operation, func(t *testing.T) {
				ctx := t.Context()
				r := advancedcluster.Resource()
				var schemaResp resource.SchemaResponse
				r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
				typ := schemaResp.Schema.Type().TerraformType(ctx)
				model := func(withStorage, configOnly, updated bool) tfsdk.Plan {
					scaling := map[string]any{"compute_enabled": true, "compute_max_instance_size": "M20", "compute_scale_down_enabled": false}
					if withStorage {
						scaling["storage_config"] = map[string]any{"shard_size_limit_gb": int64(1024)}
					}
					region := map[string]any{
						"provider_name": "AWS", "region_name": "US_EAST_1", "priority": int64(7), "auto_scaling": scaling,
					}
					for role, nodeCount := range map[string]int64{"electable_specs": 2, "read_only_specs": 1, "analytics_specs": 1} {
						hardware := map[string]any{"instance_size": "M10", "node_count": nodeCount}
						if configOnly && omittedField == "node_count" {
							delete(hardware, "node_count")
						}
						region[role] = hardware
					}
					spec := map[string]any{"zone_name": "Existing zone", "region_configs": []any{region}}
					if !configOnly {
						spec["external_id"] = "650000000000000000000001"
						spec["zone_id"] = "650000000000000000000002"
					} else if omittedField == "zone_name" {
						delete(spec, "zone_name")
					}
					return tfsdk.Plan{Schema: schemaResp.Schema, Raw: planTestValue(typ, map[string]any{
						"name": "example", "project_id": dummyProjectID, "cluster_type": "REPLICASET", "database_edition": "INFINITE",
						"redact_client_log_data": updated, "replication_specs": []any{spec},
					})}
				}
				clearStorage := operation == "clear storage"
				prior := model(true, false, false)
				plan := model(!clearStorage, false, true)
				configuration := model(!clearStorage, true, true)
				dynamic := func(value tfsdk.Plan) *tfprotov6.DynamicValue {
					result, err := tfprotov6.NewDynamicValue(typ, value.Raw)
					require.NoError(t, err)
					return &result
				}
				server, err := acc.TestAccProviderV6Factories["mongodbatlas"]()
				require.NoError(t, err)
				planned, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
					TypeName: "mongodbatlas_advanced_cluster", PriorState: dynamic(prior),
					ProposedNewState: dynamic(plan), Config: dynamic(configuration),
				})
				require.NoError(t, err)
				require.Empty(t, planned.Diagnostics)
				plan.Raw, err = planned.PlannedState.Unmarshal(typ)
				require.NoError(t, err)

				api := mockadmin.NewClustersAPI(t)
				r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
				api.On("UpdateCluster", mock.Anything, dummyProjectID, "example", mock.Anything).Run(func(args mock.Arguments) {
					encoded, err := json.Marshal(args[3].(*admin.ClusterDescription20240805))
					require.NoError(t, err)
					expected := `{"redactClientLogData":true}`
					if clearStorage {
						expected = `{"redactClientLogData":true,"replicationSpecs":[{
							"id":"650000000000000000000001","zoneId":"650000000000000000000002","zoneName":"Existing zone",
							"regionConfigs":[{
								"providerName":"AWS","regionName":"US_EAST_1","priority":7,
								"electableSpecs":{"instanceSize":"M10","nodeCount":2},
								"readOnlySpecs":{"instanceSize":"M10","nodeCount":1},
								"analyticsSpecs":{"instanceSize":"M10","nodeCount":1},
								"autoScaling":{"compute":{"enabled":true,"maxInstanceSize":"M20","scaleDownEnabled":false},"storageConfig":{"shardSizeLimitGB":null}}
							}]
						}]}`
					}
					require.JSONEq(t, expected, string(encoded))
				}).Return(admin.UpdateClusterApiRequest{ApiService: api}).Once()
				apiError := errors.New("request inspected")
				api.EXPECT().UpdateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
				var resp resource.UpdateResponse
				r.Update(ctx, resource.UpdateRequest{
					Plan: plan, State: tfsdk.State{Schema: schemaResp.Schema, Raw: prior.Raw},
					Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: configuration.Raw},
				}, &resp)
				require.Len(t, resp.Diagnostics.Errors(), 1)
				require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), apiError.Error())
			})
		}
	}
}

func TestUpdateRemovesShardSizeLimitAcrossReplicationSpecs(t *testing.T) {
	ctx := t.Context()
	r := advancedcluster.Resource()
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	typ := schemaResp.Schema.Type().TerraformType(ctx)
	// This exercises request construction; Atlas does not yet support Infinite geosharded clusters.
	model := func(withStorage, configOnly bool) tfsdk.Plan {
		var specs []any
		for _, spec := range []struct {
			id, zoneID, zoneName, instanceSize string
			regions                            []string
			nodeCount                          int64
		}{
			{"650000000000000000000001", "650000000000000000000002", "US zone", "M10", []string{"US_EAST_1", "US_WEST_2"}, 2},
			{"650000000000000000000003", "650000000000000000000004", "EU zone", "M20", []string{"EU_WEST_1"}, 3},
		} {
			var regions []any
			for i, regionName := range spec.regions {
				hardware := map[string]any{"instance_size": spec.instanceSize}
				region := map[string]any{
					"provider_name": "AWS", "region_name": regionName, "priority": int64(7 - i), "electable_specs": hardware,
				}
				if !configOnly {
					hardware["node_count"] = spec.nodeCount
					region["read_only_specs"] = map[string]any{"instance_size": spec.instanceSize, "node_count": int64(0)}
					region["analytics_specs"] = map[string]any{"instance_size": spec.instanceSize, "node_count": int64(0)}
				}
				if withStorage {
					region["auto_scaling"] = map[string]any{"storage_config": map[string]any{"shard_size_limit_gb": int64(1024)}}
				}
				regions = append(regions, region)
			}
			replicationSpec := map[string]any{"region_configs": regions}
			if !configOnly {
				replicationSpec["external_id"] = spec.id
				replicationSpec["zone_id"] = spec.zoneID
				replicationSpec["zone_name"] = spec.zoneName
			}
			specs = append(specs, replicationSpec)
		}
		return tfsdk.Plan{Schema: schemaResp.Schema, Raw: planTestValue(typ, map[string]any{
			"name": "example", "project_id": dummyProjectID, "cluster_type": "GEOSHARDED", "database_edition": "INFINITE",
			"replication_specs": specs,
		})}
	}
	api := mockadmin.NewClustersAPI(t)
	r.(config.ImplementedResource).SetClient(&config.MongoDBClient{AtlasV2: &admin.APIClient{ClustersAPI: api}})
	api.On("UpdateCluster", mock.Anything, dummyProjectID, "example", mock.Anything).Run(func(args mock.Arguments) {
		encoded, err := json.Marshal(args[3].(*admin.ClusterDescription20240805))
		require.NoError(t, err)
		require.JSONEq(t, `{"replicationSpecs":[{
			"id":"650000000000000000000001","zoneId":"650000000000000000000002","zoneName":"US zone",
			"regionConfigs":[{
				"providerName":"AWS","regionName":"US_EAST_1","priority":7,
				"electableSpecs":{"instanceSize":"M10","nodeCount":2},
				"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}
			},{
				"providerName":"AWS","regionName":"US_WEST_2","priority":6,
				"electableSpecs":{"instanceSize":"M10","nodeCount":2},
				"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}
			}]
		},{
			"id":"650000000000000000000003","zoneId":"650000000000000000000004","zoneName":"EU zone",
			"regionConfigs":[{
				"providerName":"AWS","regionName":"EU_WEST_1","priority":7,
				"electableSpecs":{"instanceSize":"M20","nodeCount":3},
				"autoScaling":{"storageConfig":{"shardSizeLimitGB":null}}
			}]
		}]}`, string(encoded))
	}).Return(admin.UpdateClusterApiRequest{ApiService: api}).Once()
	apiError := errors.New("request inspected")
	api.EXPECT().UpdateClusterExecute(mock.Anything).Return(nil, nil, apiError).Once()
	prior, plan, configuration := model(true, false), model(false, false), model(false, true)
	var resp resource.UpdateResponse
	r.Update(ctx, resource.UpdateRequest{
		Plan: plan, State: tfsdk.State{Schema: schemaResp.Schema, Raw: prior.Raw},
		Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: configuration.Raw},
	}, &resp)
	require.Len(t, resp.Diagnostics.Errors(), 1)
	require.Contains(t, resp.Diagnostics.Errors()[0].Detail(), apiError.Error())
}

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
