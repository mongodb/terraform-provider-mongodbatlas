package cluster_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	deprecatedResourceDiagSummary = "Deprecated Resource"
)

func TestAccClusterRSClusterMigrateState_empty_advancedConfig(t *testing.T) {
	acc.SkipInUnitTest(t)
	v0State := map[string]any{
		"project_id":                  "test-id",
		"name":                        "test-cluster",
		"provider_instance_size_name": "M10",
		"provider_name":               "AWS",
		"replication_specs": []any{
			map[string]any{
				"num_shards": 1,
			},
		},
		"advanced_configuration": map[string]any{},
	}

	v0Config := terraform.NewResourceConfigRaw(v0State)
	diags := cluster.ResourceClusterResourceV0().Validate(v0Config)

	if len(diags) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State := cluster.MigrateAdvancedConfiguration(v0State)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = cluster.Resource().Validate(v1Config)
	if isErrorDiags(diags) {
		fmt.Println(diags)
		t.Error("migrated cluster advanced config is invalid")

		return
	}
}

func TestAccClusterRSClusterMigrateState_with_advancedConfig(t *testing.T) {
	acc.SkipInUnitTest(t)
	v0State := map[string]any{
		"project_id":                  "test-id",
		"name":                        "test-cluster",
		"provider_instance_size_name": "M10",
		"provider_name":               "AWS",
		"replication_specs": []any{
			map[string]any{
				"num_shards": 1,
			},
		},
		"advanced_configuration": map[string]any{
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
	diags := cluster.ResourceClusterResourceV0().Validate(v0Config)
	if len(diags) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State := cluster.MigrateAdvancedConfiguration(v0State)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = cluster.Resource().Validate(v1Config)
	if isErrorDiags(diags) {
		fmt.Println(diags)
		t.Error("migrated cluster advanced config is invalid")

		return
	}
}

func TestAccClusterRSClusterMigrateState_with_defaultAdvancedConfig_v0_5_1(t *testing.T) {
	acc.SkipInUnitTest(t)
	v0State := map[string]any{
		"project_id":                  "test-id",
		"name":                        "test-cluster",
		"provider_instance_size_name": "M10",
		"provider_name":               "AWS",
		"replication_specs": []any{
			map[string]any{
				"num_shards": 1,
			},
		},
		"advanced_configuration": map[string]any{
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
	diags := cluster.ResourceClusterResourceV0().Validate(v0Config)
	if len(diags) > 0 {
		t.Error("test precondition failed - invalid mongodb cluster v0 config")

		return
	}

	// test migrate function
	v1State := cluster.MigrateAdvancedConfiguration(v0State)

	v1Config := terraform.NewResourceConfigRaw(v1State)
	diags = cluster.Resource().Validate(v1Config)
	if isErrorDiags(diags) {
		fmt.Println(diags)
		t.Error("migrated cluster advanced config is invalid")

		return
	}
}

func isErrorDiags(diags diag.Diagnostics) bool {
	return len(diags) > 0 && !strings.Contains(diags[0].Summary, deprecatedResourceDiagSummary)
}
