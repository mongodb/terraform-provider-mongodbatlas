package advancedcluster_test

// func TestMigAdvancedCluster_empty_advancedConfig(t *testing.T) {
// 	acc.SkipIfAdvancedClusterV2Schema(t) // This test is specific to the legacy schema
// 	acc.SkipInUnitTest(t)                // needed because TF test infra is not used
// 	v0State := map[string]any{
// 		"project_id":   "test-id",
// 		"name":         "test-cluster",
// 		"cluster_type": "REPLICASET",
// 		"replication_specs": []any{
// 			map[string]any{
// 				"region_configs": []any{
// 					map[string]any{
// 						"electable_specs": []any{
// 							map[string]any{
// 								"instance_size": "M30",
// 								"node_count":    3,
// 							},
// 						},
// 						"provider_name": "AWS",
// 						"region_name":   "US_WEST_2",
// 						"priority":      7,
// 					},
// 				},
// 			},
// 		},
// 		"bi_connector": []any{
// 			map[string]any{
// 				"enabled":         1,
// 				"read_preference": "secondary",
// 			},
// 		},
// 	}

// 	v0Config := terraform.NewResourceConfigRaw(v0State)
// 	diags := advancedcluster.ResourceV0().Validate(v0Config)

// 	if len(diags) > 0 {
// 		t.Error("test precondition failed - invalid mongodb cluster v0 config")

// 		return
// 	}

// 	// test migrate function
// 	v1State := advancedcluster.MigrateBIConnectorConfig(v0State)

// 	v1Config := terraform.NewResourceConfigRaw(v1State)
// 	diags = advancedcluster.Resource().Validate(v1Config)
// 	if len(diags) > 0 {
// 		fmt.Println(diags)
// 		t.Error("migrated cluster advanced config is invalid")

// 		return
// 	}
// }

// func TestMigAdvancedCluster_v0StateUpgrade_ReplicationSpecs(t *testing.T) {
// 	acc.SkipIfAdvancedClusterV2Schema(t) // This test is specific to the legacy schema
// 	acc.SkipInUnitTest(t)                // needed because TF test infra is not used
// 	v0State := map[string]any{
// 		"project_id":     "test-id",
// 		"name":           "test-cluster",
// 		"cluster_type":   "REPLICASET",
// 		"backup_enabled": true,
// 		"disk_size_gb":   256,
// 		"replication_specs": []any{
// 			map[string]any{
// 				"zone_name": "Test Zone",
// 				"region_configs": []any{
// 					map[string]any{
// 						"priority":      7,
// 						"provider_name": "AWS",
// 						"region_name":   "US_WEST_2",
// 						"electable_specs": []any{
// 							map[string]any{
// 								"instance_size": "M30",
// 								"node_count":    3,
// 							},
// 						},
// 						"read_only_specs": []any{
// 							map[string]any{
// 								"disk_iops":     0,
// 								"instance_size": "M30",
// 								"node_count":    0,
// 							},
// 						},
// 						"auto_scaling": []any{
// 							map[string]any{
// 								"compute_enabled":            true,
// 								"compute_max_instance_size":  "M60",
// 								"compute_min_instance_size":  "M30",
// 								"compute_scale_down_enabled": true,
// 								"disk_gb_enabled":            false,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	v0Config := terraform.NewResourceConfigRaw(v0State)
// 	diags := advancedcluster.ResourceV0().Validate(v0Config)

// 	if diags.HasError() {
// 		fmt.Println(diags)
// 		t.Error("test precondition failed - invalid mongodb cluster v0 config")

// 		return
// 	}

// 	// test migrate function
// 	v1State := advancedcluster.MigrateBIConnectorConfig(v0State)

// 	v1Config := terraform.NewResourceConfigRaw(v1State)
// 	diags = advancedcluster.Resource().Validate(v1Config)
// 	if diags.HasError() {
// 		fmt.Println(diags)
// 		t.Error("migrated advanced cluster replication_specs invalid")

// 		return
// 	}

// 	if len(v1State["replication_specs"].([]any)) != len(v0State["replication_specs"].([]any)) {
// 		t.Error("migrated replication specs did not contain the same number of elements")

// 		return
// 	}
// }
