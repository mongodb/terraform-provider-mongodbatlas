package acc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestConvertAdvancedClusterToTPF(t *testing.T) {
	var (
		input = `
			resource "mongodbatlas_advanced_cluster" "test" {
				project_id   = "66d979971ec97b7de1ef8777"
				name         = "test-acc-tf-c-2683795087811441116"
				cluster_type = "REPLICASET"
				replication_specs {
								region_configs {
												electable_specs {
																instance_size = "M5"
												}
												provider_name         = "TENANT"
												backing_provider_name = "AWS"
												region_name           = "US_EAST_1"
												priority              = 7
								}
				}
			}

			data "mongodbatlas_advanced_cluster" "test" {
				project_id = mongodbatlas_advanced_cluster.test.project_id
				name         = mongodbatlas_advanced_cluster.test.name
			}

			data "mongodbatlas_advanced_clusters" "test" {
				project_id = mongodbatlas_advanced_cluster.test.project_id
			}
 		`

		expected = `
			resource "mongodbatlas_advanced_cluster" "test" {
				project_id  = "66d979971ec97b7de1ef8777"
				name         = "test-acc-tf-c-2683795087811441116"
				cluster_type = "REPLICASET"
				replication_specs {
								region_configs {
												electable_specs {
																instance_size = "M5"
												}
												provider_name         = "TENANT"
												backing_provider_name = "AWS"
												region_name           = "US_EAST_1"
												priority              = 7
								}
				}
			}

			data "mongodbatlas_advanced_cluster" "test" {
				project_id = mongodbatlas_advanced_cluster.test.project_id
				name         = mongodbatlas_advanced_cluster.test.name
			}

			data "mongodbatlas_advanced_clusters" "test" {
				project_id = mongodbatlas_advanced_cluster.test.project_id
			}
 		`
	)
	actual := input
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
