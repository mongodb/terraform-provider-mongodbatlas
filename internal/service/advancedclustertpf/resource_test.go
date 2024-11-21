package advancedclustertpf_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

const (
	resourceName        = "mongodbatlas_advanced_cluster.test"
	processResponseOnly = "processResponseOnly"
	projectID      = "111111111111111111111111"
	clusterName    = "test"
)

func TestAdvancedCluster_replicaset(t *testing.T) {
	var (
		oneNewVariable = "backup_enabled = false"
		fullUpdate     = `
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
		paused = true
		pit_enabled = true
		redact_client_log_data = true
		replica_set_scaling_strategy = "NODE_TYPE"
		# retain_backups_enabled = true # only set on delete
		root_cert_type = "ISRGROOTX1"
		# termination_protection_enabled = true # must be reset to false to enable delete
		version_release_system = "CONTINUOUS"
		`
		// # oplog_min_retention_hours                                      = 5.5
		// # oplog_size_mb                                                  = 1000
		// # fail_index_key_too_long 								        = true # only valid for MongoDB version 4.4 and earlier
		advClusterConfig = `
		advanced_configuration = {
			change_stream_options_pre_and_post_images_expire_after_seconds = 100
			default_read_concern                                           = "available"
			default_write_concern                                          = "majority"
			javascript_enabled                                             = false
			minimum_enabled_tls_protocol                                   = "TLS1_0"
			no_table_scan                                                  = true
			sample_refresh_interval_bi_connector                           = 310
			sample_size_bi_connector                                       = 110
			transaction_lifetime_limit_seconds                             = 300
		}
		`
		fullUpdateResumed = strings.Replace(fullUpdate, "paused = true", "paused = false", 1)
		vars              = map[string]string{
			"groupId":     projectID,
			"clusterName": clusterName,
		}
	)
	advancedclustertpf.RetryMinTimeout = 1 * time.Second
	advancedclustertpf.RetryDelay = 1 * time.Second
	advancedclustertpf.RetryPollInterval = 100 * time.Millisecond
	mockTransport, checkFunc := unit.MockRoundTripper(t, vars, &unit.MockHTTPDataConfig{AllowMissingRequests: true, AllowReReadGet: true})

	resource.Test(t, resource.TestCase{ // Sequential as it is using global variables
		ProtoV6ProviderFactories: acc.TestAccProviderV6FactoriesWithMock(mockTransport),
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
					checkFunc,
				),
			},
			{
				Config: configBasic(projectID, clusterName, oneNewVariable),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "backup_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
					checkFunc,
				),
			},
			{
				Config: configBasic(projectID, clusterName, fullUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
					resource.TestCheckResourceAttr(resourceName, "mongo_db_major_version", "8.0"),
					resource.TestCheckResourceAttr(resourceName, "backup_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
					checkFunc,
				),
			},
			{
				Config: configBasic(projectID, clusterName, fullUpdateResumed),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
					resource.TestCheckResourceAttr(resourceName, "backup_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
					checkFunc,
				),
			},
			{
				Config: configBasic(projectID, clusterName, fullUpdateResumed+advClusterConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mongo_db_major_version", "8.0"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.change_stream_options_pre_and_post_images_expire_after_seconds", "100"),
					checkFunc,
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func TestAdvancedCluster_configSharded(t *testing.T) {
	var (
		projectID   = "111111111111111111111111"
		clusterName = "sharded-multi-replication"
		vars        = map[string]string{
			"groupId":     projectID,
			"clusterName": clusterName,
		}
	)
	advancedclustertpf.RetryMinTimeout = 1 * time.Second
	advancedclustertpf.RetryDelay = 1 * time.Second
	advancedclustertpf.RetryPollInterval = 100 * time.Millisecond
	mockTransport, checkFunc := unit.MockRoundTripper(t, vars, &unit.MockHTTPDataConfig{AllowMissingRequests: true, AllowReReadGet: true})
	resource.Test(t, resource.TestCase{ // Sequential as it is using global variables
		ProtoV6ProviderFactories: acc.TestAccProviderV6FactoriesWithMock(mockTransport),
		Steps: []resource.TestStep{
			{
				Config: configSharded(projectID, clusterName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					checkFunc,
				),
			},
			{
				Config: configSharded(projectID, clusterName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkFunc,
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func configBasic(projectID, clusterName, extra string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			timeouts = {
				create = "20s"
			}
			project_id = %[1]q
			name = %[2]q
			cluster_type = "REPLICASET"
			replication_specs = [{
				region_configs = [{
					priority        = 7
					provider_name = "AWS"
					region_name     = "US_EAST_1"
					auto_scaling = {
						compute_scale_down_enabled = false # necessary to have similar SDKv2 request
						compute_enabled = false # necessary to have similar SDKv2 request
						disk_gb_enabled = true
					}
					electable_specs = {
						node_count = 3
						instance_size = "M10"
						disk_size_gb = 10
					}
				}]
			}]
			%[3]s
		}
	`, projectID, clusterName, extra)
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
					%[5]s
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}]
			}]
			}
		

	`, projectID, clusterName, autoScaling, analyticsSpecs, analyticsSpecsForSpec2)
}
func TestClusterAdvancedCluster_basicTenant(t *testing.T) {
	var (
		clusterName = "test-acc-tf-c-8049930413007488732"
		clusterNameUpdated = "test-acc-tf-c-91771214182147246"
		vars = map[string]string{
			"groupId":   projectID,
			"clusterName": clusterName,
			"clusterName2": clusterNameUpdated,
		}
	)
	advancedclustertpf.RetryMinTimeout = 1 * time.Second
	advancedclustertpf.RetryDelay = 1 * time.Second
	advancedclustertpf.RetryPollInterval = 100 * time.Millisecond
	mockTransport, checkFunc := unit.MockRoundTripper(t, vars, &unit.MockHTTPDataConfig{AllowMissingRequests: true, AllowReReadGet: true})

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6FactoriesWithMock(mockTransport),
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configTenant(projectID, clusterName),
				Check:  resource.ComposeAggregateTestCheckFunc(checkTenant(projectID, clusterName), checkFunc),
			},
			{
				Config: configTenant(projectID, clusterNameUpdated),
				Check:  resource.ComposeAggregateTestCheckFunc(checkTenant(projectID, clusterNameUpdated), checkFunc),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}


func configTenant(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M5"
					}
					provider_name         = "TENANT"
					backing_provider_name = "AWS"
					region_name           = "US_EAST_1"
					priority              = 7
				}]
			}]
		}
	`, projectID, name)
}

func checkTenant(projectID, name string) resource.TestCheckFunc {
	attrsSet := []string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"}
	attrsMap := 		map[string]string{
		"project_id":                           projectID,
		"name":                                 name,
		"termination_protection_enabled":       "false",
		"global_cluster_self_managed_sharding": "false",
		"labels.#":                             "0",
	}
	checks := acc.AddAttrSetChecks(resourceName, nil, attrsSet...)
	checks = acc.AddAttrChecks(resourceName, checks, attrsMap)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}