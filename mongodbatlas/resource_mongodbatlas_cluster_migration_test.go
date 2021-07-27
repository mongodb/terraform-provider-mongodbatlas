package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMongoDBAtlasClusterMigrateState_empty_advancedConfig(t *testing.T) {
	v0State := map[string]interface{}{
		"project_id":                  "test-id",
		"name":                        "test-cluster",
		"provider_instance_size_name": "M10",
		"provider_name":               "AWS",
		"replication_specs": []interface{}{
			map[string]interface{}{
				"num_shards": 1,
			},
		},
		"advanced_configuration": map[string]interface{}{},
	}

	v0Config := terraform.NewResourceConfigRaw(v0State)
	diags := resourceMongoDBAtlasClusterResourceV0().Validate(v0Config)

	if len(diags) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State := migrateAdvancedConfiguration(v0State)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = resourceMongoDBAtlasCluster().Validate(v1Config)
	if len(diags) > 0 {
		fmt.Println(diags)
		t.Error("migrated cluster advanced config is invalid")

		return
	}
}

func TestAccMongoDBAtlasClusterMigrateState_with_advancedConfig(t *testing.T) {
	v0State := map[string]interface{}{
		"project_id":                  "test-id",
		"name":                        "test-cluster",
		"provider_instance_size_name": "M10",
		"provider_name":               "AWS",
		"replication_specs": []interface{}{
			map[string]interface{}{
				"num_shards": 1,
			},
		},
		"advanced_configuration": map[string]interface{}{
			"fail_index_key_too_long":              "true",
			"javascript_enabled":                   "true",
			"minimum_enabled_tls_protocol":         "TLS1_2",
			"no_table_scan":                        "false",
			"oplog_size_mb":                        "1000",
			"sample_refresh_interval_bi_connector": "310",
			"sample_size_bi_connector":             "110",
		},
	}

	v0Config := terraform.NewResourceConfigRaw(v0State)
	diags := resourceMongoDBAtlasClusterResourceV0().Validate(v0Config)
	if len(diags) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State := migrateAdvancedConfiguration(v0State)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = resourceMongoDBAtlasCluster().Validate(v1Config)
	if len(diags) > 0 {
		fmt.Println(diags)
		t.Error("migrated cluster advanced config is invalid")

		return
	}
}

func TestAccMongoDBAtlasClusterMigrateState_with_defaultAdvancedConfig_v0_5_1(t *testing.T) {
	v0State := map[string]interface{}{
		"project_id":                  "test-id",
		"name":                        "test-cluster",
		"provider_instance_size_name": "M10",
		"provider_name":               "AWS",
		"replication_specs": []interface{}{
			map[string]interface{}{
				"num_shards": 1,
			},
		},
		"advanced_configuration": map[string]interface{}{
			"fail_index_key_too_long":              "true",
			"javascript_enabled":                   "true",
			"minimum_enabled_tls_protocol":         "TLS1_2",
			"no_table_scan":                        "false",
			"oplog_size_mb":                        "",
			"sample_refresh_interval_bi_connector": "",
			"sample_size_bi_connector":             "",
		},
	}

	v0Config := terraform.NewResourceConfigRaw(v0State)
	diags := resourceMongoDBAtlasClusterResourceV0().Validate(v0Config)
	if len(diags) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State := migrateAdvancedConfiguration(v0State)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = resourceMongoDBAtlasCluster().Validate(v1Config)
	if len(diags) > 0 {
		fmt.Println(diags)
		t.Error("migrated cluster advanced config is invalid")

		return
	}
}
