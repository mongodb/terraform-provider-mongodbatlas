package acc_test

import (
	"sort"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
)

func TestConvertToPreviewProviderV2AttrsMapAndAttrsSet(t *testing.T) {
	if !config.PreviewProviderV2AdvancedCluster() {
		t.Skip("Skipping test as not in PreviewProviderV2AdvancedCluster")
	}
	attrsMap := map[string]string{
		"attr":                            "val1",
		"electable_specs.0":               "val2",
		"prefixbi_connector_config.0":     "val3",
		"advanced_configuration.0postfix": "val4",
		"electable_specs.0advanced_configuration.0bi_connector_config.0": "val5",
		"connection_strings.0.standard":                                  "val6",
		"connection_strings.0.standard_srv":                              "val6",
	}
	expectedMap := map[string]string{
		"attr":                          "val1",
		"electable_specs":               "val2",
		"prefixbi_connector_config":     "val3",
		"advanced_configurationpostfix": "val4",
		"electable_specsadvanced_configurationbi_connector_config": "val5",
		"connection_strings.standard":                              "val6",
		"connection_strings.standard_srv":                          "val6",
	}
	actualMap := acc.ConvertToPreviewProviderV2AttrsMap(true, attrsMap)
	assert.Equal(t, expectedMap, actualMap)

	attrsSet := make([]string, 0, len(attrsMap))
	for name := range attrsMap {
		attrsSet = append(attrsSet, name)
	}
	expectedSet := make([]string, 0, len(attrsMap))
	for name := range expectedMap {
		expectedSet = append(expectedSet, name)
	}
	actualSet := acc.ConvertToPreviewProviderV2AttrsSet(true, attrsSet)
	sort.Strings(expectedSet)
	sort.Strings(actualSet)
	assert.Equal(t, expectedSet, actualSet)
}

func TestConvertAdvancedClusterToPreviewProviderV2(t *testing.T) {
	if !config.PreviewProviderV2AdvancedCluster() {
		t.Skip("Skipping test as not in PreviewProviderV2AdvancedCluster")
	}
	var (
		input = `
			resource "mongodbatlas_advanced_cluster" "cluster2" {
				project_id   = "MY-PROJECT-ID"
				name         = "cluster2"
				cluster_type = "SHARDED"

				replication_specs {
					zone_name = "zone1"
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

				labels {
					key   = "label"
					value = "labelvalue"
  			}

				tags {
					key   = "tag"
					value = "tagvalue"
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
				timeouts {
					create = "5m"
				}
			}	
 		`
		// expected has the attributes sorted alphabetically to match the output of ConvertAdvancedClusterToPreviewProviderV2
		expected = `
			resource "mongodbatlas_advanced_cluster" "cluster2" {
				project_id   = "MY-PROJECT-ID"
				name         = "cluster2"
				cluster_type = "SHARDED"







				



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
						zone_name = "zone1" 
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
				timeouts = {
					create = "5m"
				}
 				labels = {
					"Key Label 1" = "Value Label 1"
					"Key Label 2" = "Value Label 2"
					"Key Label 3" = "Value Label 3"
					label = "labelvalue"
  			}
				tags = {
					"Key Tag 1" = "Value Tag 1"
					"Key Tag 2" = "Value Tag 2"
					tag = "tagvalue"
  			}
			}
 		`
	)
	actual := acc.ConvertAdvancedClusterToTPF(t, true, input)
	acc.AssertEqualHCL(t, expected, actual)
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
