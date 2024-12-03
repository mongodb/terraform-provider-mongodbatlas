package acc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestConvertAdvancedClusterToTPF(t *testing.T) {
	var (
		input = `
			resource "mongodbatlas_advanced_cluster" "cluster2" {
				project_id   = "MY-PROJECT-ID"
				name         = "cluster2"
				cluster_type = "SHARDED"

				replication_specs {
					region_configs {
						electable_specs {
							disk_size_gb  = 10
							instance_size = "M10"
							node_count    = 3
						}
						analytics_specs {
							disk_size_gb  = 10
							instance_size = "M10"
							node_count    = 1
						}
						priority      = 7
						provider_name = "AWS"
						region_name   = "EU_WEST_1"
					}
					region_configs {
						electable_specs {
							instance_size = "M30"
							node_count    = 2
						}
						provider_name = "AZURE"
						priority      = 6
						region_name   = "US_EAST_2"
					}
				}

				replication_specs {
					region_configs {
						electable_specs {
							disk_size_gb  = 10
							instance_size = "M10"
							node_count    = 3
						}
						analytics_specs {
							disk_size_gb  = 10
							instance_size = "M10"
							node_count    = 1
						}
						priority      = 7
						provider_name = "AWS"
						region_name   = "EU_WEST_1"
					}
				}

 				tags {
					key   = "Key Tag 2"
					value = "Value Tag 2"
  			}

 				labels {
					key   = "Key Label 1"
					value = "Value Label 1"
  			}

				tags {
					key   = "Key Tag 1"
					value = "Value Tag 1"
			  }

 				labels {
					key   = "Key Label 2"
					value = "Value Label 2"
  			}

 				labels {
					key   = "Key Label 3"
					value = "Value Label 3"
  			}

				advanced_configuration  {
					fail_index_key_too_long              = false
					javascript_enabled                   = true
					minimum_enabled_tls_protocol         = "TLS1_1"
					no_table_scan                        = false
					oplog_size_mb                        = 1000
					sample_size_bi_connector			 = 110
					sample_refresh_interval_bi_connector = 310
			    transaction_lifetime_limit_seconds   = 300  
			    change_stream_options_pre_and_post_images_expire_after_seconds = 100
				}

				bi_connector_config {
 					enabled         = true
  				read_preference = "secondary"
				}
			}	
 		`
		// expected has the attributes sorted alphabetically to match the output of ConvertAdvancedClusterToTPF
		expected = `
			resource "mongodbatlas_advanced_cluster" "cluster2" {
				project_id   = "MY-PROJECT-ID"
				name         = "cluster2"
				cluster_type = "SHARDED"








				
				labels = [{
					key   = "Key Label 2"
					value = "Value Label 1"
  			}, {
					key   = "Key Label 2"
					value = "Value Label 2"
  			}, {
					key   = "Key Label 3"
					value = "Value Label 3"
  			}]
				tags = [{
					key   = "Key Tag 2"
					value = "Value Tag 2"
  			}, {
					key   = "Key Tag 1"
					value = "Value Tag 1"
  			}]
				replication_specs = [{
						region_configs = [{
							analytics_specs = {
								disk_size_gb  = 10
								instance_size = "M10"
								node_count    = 1
							}
							electable_specs = {
								disk_size_gb  = 10
								instance_size = "M10"
								node_count    = 3
							}
							priority      = 7
							provider_name = "AWS"
							region_name   = "EU_WEST_1"
						},  {
							electable_specs = {
								instance_size = "M30"
								node_count    = 2
							}
							priority      = 6
							provider_name = "AZURE"
							region_name   = "US_EAST_2"
						}] 
						}, {
						region_configs = [{
							analytics_specs = {
								disk_size_gb  = 10
								instance_size = "M10"
								node_count    = 1
							}
							electable_specs = {
								disk_size_gb  = 10
								instance_size = "M10"
								node_count    = 3
							}
							priority      = 7
							provider_name = "AWS"
							region_name   = "EU_WEST_1"
						}]
					}]
				advanced_configuration = {
			    change_stream_options_pre_and_post_images_expire_after_seconds = 100
					fail_index_key_too_long              = false
					javascript_enabled                   = true
					minimum_enabled_tls_protocol         = "TLS1_1"
					no_table_scan                        = false
					oplog_size_mb                        = 1000
					sample_refresh_interval_bi_connector = 310
					sample_size_bi_connector			 = 110
			    transaction_lifetime_limit_seconds   = 300  
				}
				bi_connector_config = {
 					enabled         = true
  				read_preference = "secondary"
				}
			}
 		`
	)
	t.Setenv(config.LatestAdvancedClusterEnabledEnvVar, "true")
	if config.LatestAdvancedClusterEnabled() {
		actual := acc.ConvertAdvancedClusterToTPF(t, input)
		acc.AssertEqualHCL(t, expected, actual, "convert in new format")
	}

	t.Setenv(config.LatestAdvancedClusterEnabledEnvVar, "false")
	actual := acc.ConvertAdvancedClusterToTPF(t, input)
	acc.AssertEqualHCL(t, input, actual, "convert in old format")
}

func TestAssertEqualHCL(t *testing.T) {
	var (
		val1 = `
			resource "type1" "name1" {
				attr1 = "val1"
				block1 {
					attr2 = "val2"
				}
			}
 		`
		val2 = `
			resource "type1"      "name1" {
				attr1 =        "val1"
				block1    {
					attr2="val2"
				      }			
		}
 		`
	)
	acc.AssertEqualHCL(t, val1, val2)
}
