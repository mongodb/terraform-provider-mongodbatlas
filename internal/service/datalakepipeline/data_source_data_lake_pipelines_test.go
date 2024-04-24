package datalakepipeline_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccDataLakeDSPlural_basic(t *testing.T) {
	var (
		resourceName       = "mongodbatlas_data_lake_pipeline.test"
		dataSourceName     = "data.mongodbatlas_data_lake_pipelines.testDataSource"
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		firstClusterName   = acc.RandomClusterName()
		secondClusterName  = acc.RandomClusterName()
		firstPipelineName  = acc.RandomName()
		secondPipelineName = acc.RandomName()
		projectName        = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDSPlural(orgID, projectName, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func configDSPlural(orgID, projectName, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName string) string {
	return fmt.Sprintf(`

		resource "mongodbatlas_project" "project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "azure_conf" {
			project_id   = mongodbatlas_project.project.id
			name         = %[3]q
			cluster_type = "REPLICASET"
		
			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					provider_name = "AZURE"
					priority      = 7
					region_name   = "US_EAST_2"
				}
			}
			backup_enabled               = true
		}

		resource "mongodbatlas_advanced_cluster" "azure_conf2" {
			project_id   = mongodbatlas_project.project.id
			name         = %[4]q
			cluster_type = "REPLICASET"
		
			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					provider_name = "AZURE"
					priority      = 7
					region_name   = "US_EAST_2"
				}
			}
			backup_enabled               = true
		}

		resource "mongodbatlas_data_lake_pipeline" "test" {
			project_id 	=  mongodbatlas_project.project.id
			name 		= %[5]q
			sink {
				type = "DLS"
				partition_fields {
						field_name = "access"
						order = 0
				}
			}	
	
			source {
				type = "ON_DEMAND_CPS"
				cluster_name = mongodbatlas_advanced_cluster.azure_conf.name
				database_name = "sample_airbnb"
				collection_name = "listingsAndReviews"
			}

			transformations {
				field = "test"
				type =  "EXCLUDE"
			}
		}

		resource "mongodbatlas_data_lake_pipeline" "test2" {
			project_id       =  mongodbatlas_project.project.id
			name			 = 	%[6]q
			sink {
				type = "DLS"
				partition_fields {
						field_name = "access"
						order = 0
				}
			}	
	
			source {
				type = "ON_DEMAND_CPS"
				cluster_name = mongodbatlas_advanced_cluster.azure_conf2.name
				database_name = "sample_airbnb"
				collection_name = "listingsAndReviews"
			}

			transformations {
				field = "test"
				type =  "EXCLUDE"
			}
		}

		data "mongodbatlas_data_lake_pipelines" "testDataSource" {
			project_id       = mongodbatlas_data_lake_pipeline.test.project_id
		}
	`, orgID, projectName, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName)
}
