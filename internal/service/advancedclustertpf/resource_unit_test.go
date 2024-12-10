package advancedclustertpf_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/tc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var (
	mockConfig = &unit.MockHTTPDataConfig{AllowMissingRequests: true, SideEffect: shortenRetries, IsDiffMustSubstrings: []string{"/clusters"}}
)

func shortenRetries() error {
	advancedclustertpf.RetryMinTimeout = 100 * time.Millisecond
	advancedclustertpf.RetryDelay = 100 * time.Millisecond
	advancedclustertpf.RetryPollInterval = 100 * time.Millisecond
	return nil
}

func TestMockAdvancedCluster_replicasetAdvConfigUpdate(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
		fullUpdate  = `
	backup_enabled = true
	bi_connector_config = {
		enabled = true
	}
	# config_server_management_mode = "ATLAS_MANAGED" UNSTABLE: After applying this test step, the non-refresh plan was not empty
	labels = [{
		key   = "env"
		value = "test"
	}]
	tags = [{
		key   = "env"
		value = "test"
	}]
	mongo_db_major_version = "8.0"
	pit_enabled = true
	redact_client_log_data = true
	replica_set_scaling_strategy = "NODE_TYPE"
	# retain_backups_enabled = true # only set on delete
	root_cert_type = "ISRGROOTX1"
	# termination_protection_enabled = true # must be reset to false to enable delete
	version_release_system = "CONTINUOUS"
	
	advanced_configuration = {
		change_stream_options_pre_and_post_images_expire_after_seconds = 100
		default_read_concern                                           = "available"
		default_write_concern                                          = "majority"
		javascript_enabled                                             = true
		minimum_enabled_tls_protocol                                   = "TLS1_0"
		no_table_scan                                                  = true
		sample_refresh_interval_bi_connector                           = 310
		sample_size_bi_connector                                       = 110
		transaction_lifetime_limit_seconds                             = 300
	}
`
		// # oplog_min_retention_hours                                      = 5.5
		// # oplog_size_mb                                                  = 1000
		// # fail_index_key_too_long 								        = true # only valid for MongoDB version 4.4 and earlier
	)
	testCase := resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "2000s"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.AWS:US_EAST_1", "67345bd9905b8c30c54fd220"),
				),
			},
			{
				Config: configBasic(projectID, clusterName, fullUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mongo_db_major_version", "8.0"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.change_stream_options_pre_and_post_images_expire_after_seconds", "100"),
				),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
	unit.MockTestCaseAndRun(t, mockConfig, &testCase)
}

func TestMockAdvancedCluster_configSharded(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	testCase := resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configSharded(projectID, clusterName, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
			},
			{
				Config: configSharded(projectID, clusterName, true),
				Check:  resource.TestCheckResourceAttr(resourceName, "name", clusterName),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
	unit.MockTestCaseAndRun(t, mockConfig, &testCase)
}

func configSharded(projectID, clusterName string, withUpdate bool) string {
	var autoScaling, analyticsSpecs string
	if withUpdate {
		autoScaling = `
			auto_scaling = {
				disk_gb_enabled = true
			}`
		analyticsSpecs = `
			analytics_specs = {
				instance_size   = "M30"
				node_count      = 1
				ebs_volume_type = "PROVISIONED"
				disk_iops       = 2000
			}`
	}
	// SDK v2 Implementation receives many warnings, one of them: `.replication_specs[1].region_configs[0].analytics_specs[0].disk_iops: was cty.NumberIntVal(2000), but now cty.NumberIntVal(1000)`
	// Therefore, in TPF we are forced to set the value that will be returned by the API (1000)
	// The rule is: For any replication spec, the `(analytics|electable|read_only)_spec.disk_iops` must be the same across all region_configs
	// The API raises no errors, but the response reflects this rule
	analyticsSpecsForSpec2 := strings.ReplaceAll(analyticsSpecs, "2000", "1000")
	return fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[2]q
		cluster_type = "SHARDED"

		replication_specs = [
			{ # shard 1
			region_configs = [{
				electable_specs = {
					instance_size   = "M30"
					disk_iops       = 2000
					node_count      = 3
					ebs_volume_type = "PROVISIONED"
				}
				%[3]s
				%[4]s
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
				}]
				},
				{ # shard 2
			region_configs = [{
				electable_specs = {
					instance_size   = "M30"
					ebs_volume_type = "PROVISIONED"
					disk_iops       = 1000
					node_count      = 3
				}
				%[3]s
				%[5]s
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
			}]
		}]
	}
	`, projectID, clusterName, autoScaling, analyticsSpecs, analyticsSpecsForSpec2)
}

func TestMockClusterAdvancedCluster_basicTenant(t *testing.T) {
	testCase := tc.BasicTenantTestCase(t)
	unit.MockTestCaseAndRun(t, mockConfig, testCase)
}

func TestMockClusterAdvancedClusterConfig_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
	testCase := tc.SymmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t)
	unit.MockTestCaseAndRun(t, mockConfig, testCase)
}

func TestMockClusterAdvancedClusterConfig_symmetricShardedOldSchema(t *testing.T) {
	testCase := tc.SymmetricShardedOldSchema(t)
	unit.MockTestCaseAndRun(t, mockConfig, testCase)
}

func TestMockClusterAdvancedCluster_tenantUpgrade(t *testing.T) {
	testCase := tc.TenantUpgrade(t)
	unit.MockTestCaseAndRun(t, mockConfig, testCase)
}
