package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceClusterRSDataLakePipelines_basic(t *testing.T) {
	var (
		pipeline           matlas.DataLakePipeline
		resourceName       = "mongodbatlas_data_lake_pipeline.test"
		dataSourceName     = "data.mongodbatlas_data_lake_pipelines.testDataSource"
		firstClusterName   = acctest.RandomWithPrefix("test-acc-index")
		secondClusterName  = acctest.RandomWithPrefix("test-acc-index")
		projectID          = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		firstPipelineName  = acctest.RandomWithPrefix("test-acc-index")
		secondPipelineName = acctest.RandomWithPrefix("test-acc-index")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasDataLakeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasDataLakePipelinesConfig(projectID, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDataLakePipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.1.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.1.state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.1.project_id"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasDataLakePipelinesConfig(projectID, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "aws_conf" {
			project_id   = %[1]q
			name         = %[2]q
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

		resource "mongodbatlas_advanced_cluster" "aws_conf2" {
			project_id   = %[1]q
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
			project_id       =  "%[1]s"
			name = "%[4]s"
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

		resource "mongodbatlas_data_lake_pipeline" "test2" {
			project_id       =  "%[1]s"
			name			 = 	"%[5]s"
			sink {
				type = "DLS"
				partition_fields {
						field_name = "access"
						order = 0
				}
			}	
	
			source {
				type = "ON_DEMAND_CPS"
				cluster_name = mongodbatlas_advanced_cluster.aws_conf2.name
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
	`, projectID, firstClusterName, secondClusterName, firstPipelineName, secondPipelineName)
}
