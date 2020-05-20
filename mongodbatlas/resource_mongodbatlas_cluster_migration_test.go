package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
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

	rawConfigV0, err := config.NewRawConfig(v0State)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	v0Config := terraform.NewResourceConfig(rawConfigV0)
	warns, errs := resourceMongoDBAtlasClusterResourceV0().Validate(v0Config)
	if len(warns) > 0 || len(errs) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")
		return
	}

	//test migrate function
	v1State := migrateAdvancedConfiguration(v0State)

	rawConfigV1, err := config.NewRawConfig(v1State)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	v1Config := terraform.NewResourceConfig(rawConfigV1)
	warns, errs = resourceMongoDBAtlasCluster().Validate(v1Config)
	if len(warns) > 0 || len(errs) > 0 {
		fmt.Println(warns, errs)
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

	rawConfigV0, err := config.NewRawConfig(v0State)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	v0Config := terraform.NewResourceConfig(rawConfigV0)
	warns, errs := resourceMongoDBAtlasClusterResourceV0().Validate(v0Config)
	if len(warns) > 0 || len(errs) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")
		return
	}

	//test migrate function
	v1State := migrateAdvancedConfiguration(v0State)

	rawConfigV1, err := config.NewRawConfig(v1State)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	v1Config := terraform.NewResourceConfig(rawConfigV1)
	warns, errs = resourceMongoDBAtlasCluster().Validate(v1Config)
	if len(warns) > 0 || len(errs) > 0 {
		fmt.Println(warns, errs)
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

	rawConfigV0, err := config.NewRawConfig(v0State)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	v0Config := terraform.NewResourceConfig(rawConfigV0)
	warns, errs := resourceMongoDBAtlasClusterResourceV0().Validate(v0Config)
	if len(warns) > 0 || len(errs) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")
		return
	}

	//test migrate function
	v1State := migrateAdvancedConfiguration(v0State)

	rawConfigV1, err := config.NewRawConfig(v1State)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	v1Config := terraform.NewResourceConfig(rawConfigV1)
	warns, errs = resourceMongoDBAtlasCluster().Validate(v1Config)
	if len(warns) > 0 || len(errs) > 0 {
		fmt.Println(warns, errs)
		t.Error("migrated cluster advanced config is invalid")
		return
	}
}
