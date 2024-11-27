package acc_test

import (
	"testing"

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
							instance_size = "M10"
							node_count    = 3
							disk_size_gb  = 10
						}
						analytics_specs {
							instance_size = "M10"
							node_count    = 1
							disk_size_gb  = 10
						}
						provider_name = "AWS"
						priority      = 7
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
							instance_size = "M10"
							node_count    = 3
							disk_size_gb  = 10
						}
						analytics_specs {
							instance_size = "M10"
							node_count    = 1
							disk_size_gb  = 10
						}
						provider_name = "AWS"
						priority      = 7
						region_name   = "EU_WEST_1"
					}
				}
			}	
 		`

		expected = `
			resource "mongodbatlas_advanced_cluster" "cluster2" {
				project_id   = "MY-PROJECT-ID"
				name         = "cluster2"
				cluster_type = "SHARDED"

				replication_specs = [
					{
						region_configs = [{
							electable_specs = {
								instance_size = "M10"
								node_count    = 3
								disk_size_gb  = 10
							}
							analytics_specs = {
								instance_size = "M10"
								node_count    = 1
								disk_size_gb  = 10
							}
							provider_name = "AWS"
							priority      = 7
							region_name   = "EU_WEST_1"
						},  {
							electable_specs = {
								instance_size = "M30"
								node_count    = 2
							}
							provider_name = "AZURE"
							priority      = 6
							region_name   = "US_EAST_2"
						}] }, {
						region_configs = [{
							electable_specs = {
								instance_size = "M10"
								node_count    = 3
								disk_size_gb  = 10
							}
							analytics_specs = {
								instance_size = "M10"
								node_count    = 1
								disk_size_gb  = 10
							}
							provider_name = "AWS"
							priority      = 7
							region_name   = "EU_WEST_1"
						}]
					}
				]
			}
 		`
	)
	actual := acc.ConvertAdvancedClusterToTPF(t, input)
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
