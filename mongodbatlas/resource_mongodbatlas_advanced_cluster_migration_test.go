package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccClusterRSAdvancedClusterMigrateState_empty_advancedConfig(t *testing.T) {
	v0State := map[string]interface{}{
		"project_id":   "test-id",
		"name":         "test-cluster",
		"cluster_type": "REPLICASET",
		"replication_specs": []interface{}{
			map[string]interface{}{
				"region_configs": []interface{}{
					map[string]interface{}{
						"electable_specs": []interface{}{
							map[string]interface{}{
								"instance_size": "M30",
								"node_count":    3,
							},
						},
						"provider_name": "AWS",
						"region_name":   "US_EAST_1",
						"priority":      7,
					},
				},
			},
		},
		"bi_connector": []interface{}{
			map[string]interface{}{
				"enabled":         1,
				"read_preference": "secondary",
			},
		},
	}

	v0Config := terraform.NewResourceConfigRaw(v0State)
	diags := resourceMongoDBAtlasAdvancedClusterResourceV0().Validate(v0Config)

	if len(diags) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State := migrateBIConnectorConfig(v0State)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = resourceMongoDBAtlasAdvancedCluster().Validate(v1Config)
	if len(diags) > 0 {
		fmt.Println(diags)
		t.Error("migrated cluster advanced config is invalid")

		return
	}
}
