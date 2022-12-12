package mongodbatlas

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccClusterRSAdvancedClusterV0StateUpgrade_ReplicationSpecs(t *testing.T) {
	v0State := map[string]interface{}{
		"project_id":     "test-id",
		"name":           "test-cluster",
		"cluster_type":   "REPLICASET",
		"backup_enabled": true,
		"disk_size_gb":   256,
		"replication_specs": []interface{}{
			map[string]interface{}{
				"zone_name": "Test Zone",
				"region_configs": []interface{}{
					map[string]interface{}{
						"priority":      7,
						"provider_name": "AWS",
						"region_name":   "US_EAST_1",
						"electable_specs": []interface{}{
							map[string]interface{}{
								"instance_size": "M30",
								"node_count":    3,
							},
						},
						"read_only_specs": []interface{}{
							map[string]interface{}{
								"disk_iops":     0,
								"instance_size": "M30",
								"node_count":    0,
							},
						},
						"auto_scaling": []interface{}{
							map[string]interface{}{
								"compute_enabled":            true,
								"compute_max_instance_size":  "M60",
								"compute_min_instance_size":  "M30",
								"compute_scale_down_enabled": true,
								"disk_gb_enabled":            false,
							},
						},
					},
				},
			},
		},
	}

	v0Config := terraform.NewResourceConfigRaw(v0State)
	diags := resourceMongoDBAtlasAdvancedClusterV0().Validate(v0Config)

	if len(diags) > 0 {
		fmt.Println(diags)
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State, _ := resourceMongoDBAtlasAdvancedClusterStateUpgradeV0(context.Background(), v0State, nil)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = resourceMongoDBAtlasAdvancedCluster().Validate(v1Config)
	if len(diags) > 0 {
		fmt.Println(diags)
		t.Error("migrated advanced cluster replication_specs invalid")

		return
	}

	if len(v1State["replication_specs"].([]interface{})) != len(v0State["replication_specs"].([]interface{})) {
		t.Error("migrated replication specs did not contain the same number of elements")

		return
	}
}
