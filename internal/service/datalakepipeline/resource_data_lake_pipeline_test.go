package datalakepipeline_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccDataLakePipeline_basic(t *testing.T) {
	var (
		resourceName         = "mongodbatlas_data_lake_pipeline.test"
		dataSourceName       = "data.mongodbatlas_data_lake_pipeline.testDataSource"
		pluralDataSourceName = "data.mongodbatlas_data_lake_pipelines.testDataSource"
		orgID                = os.Getenv("MONGODB_ATLAS_ORG_ID")
		firstClusterName     = acc.RandomClusterName()
		secondClusterName    = acc.RandomClusterName()
		firstPipelineName    = acc.RandomName()
		secondPipelineName   = acc.RandomName()
		projectName          = acc.RandomProjectName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicWithPluralDS(orgID, projectName, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", firstPipelineName),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),

					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", firstPipelineName),
					resource.TestCheckResourceAttr(dataSourceName, "state", "ACTIVE"),
					resource.TestCheckResourceAttrSet(pluralDataSourceName, "results.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["name"]), nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_data_lake_pipeline" {
			continue
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().DataLakePipelinesApi.GetPipeline(context.Background(), ids["project_id"], ids["name"]).Execute()
		if err == nil {
			return fmt.Errorf("datalake (%s) still exists", ids["project_id"])
		}
	}
	return nil
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().DataLakePipelinesApi.GetPipeline(context.Background(), ids["project_id"], ids["name"]).Execute()
		if err == nil {
			return nil
		}
		return fmt.Errorf("DataLake pipeline (%s) does not exist", ids["name"])
	}
}

func configBasic(orgID, projectName, clusterName, pipelineName string) string {
	return fmt.Sprintf(`

		resource "mongodbatlas_project" "project" {
			org_id = %[1]q
			name   = %[2]q
		}
	
		resource "mongodbatlas_advanced_cluster" "aws_conf" {
			project_id   = mongodbatlas_project.project.id
			name         = %[3]q
			cluster_type = "REPLICASET"
		  
			replication_specs {
			  region_configs {
				electable_specs {
				  instance_size = "M10"
				  node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
			  }
			}
			backup_enabled               = true
		  }

		resource "mongodbatlas_data_lake_pipeline" "test" {
			project_id       = mongodbatlas_project.project.id
			name			 = %[4]q
			sink {
				type = "DLS"
				partition_fields {
						field_name = "access"
						order = 0
				}
			}	

			source {
				type = "ON_DEMAND_CPS"
				cluster_name = mongodbatlas_advanced_cluster.aws_conf.name
				database_name = "sample_airbnb"
				collection_name = "listingsAndReviews"
			}

			transformations {
				field = "test"
				type =  "EXCLUDE"
			}
		}
	`, orgID, projectName, clusterName, pipelineName)
}

func configBasicWithPluralDS(orgID, projectName, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName string) string {
	config := configBasic(orgID, projectName, firstClusterName, firstPipelineName)
	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_advanced_cluster" "azure_conf2" {
			project_id   = mongodbatlas_project.project.id
			name         = %[2]q
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

		resource "mongodbatlas_data_lake_pipeline" "test2" {
			project_id       =  mongodbatlas_project.project.id
			name			 = 	%[3]q
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

		data "mongodbatlas_data_lake_pipeline" "testDataSource" {
			project_id       = mongodbatlas_data_lake_pipeline.test.project_id
			name			 = mongodbatlas_data_lake_pipeline.test.name
		}

		data "mongodbatlas_data_lake_pipelines" "testDataSource" {
			project_id       = mongodbatlas_data_lake_pipeline.test.project_id
		}
	`, config, secondClusterName, secondPipelineName)
}
