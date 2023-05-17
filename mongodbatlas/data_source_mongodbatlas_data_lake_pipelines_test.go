package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceClusterRSDataLakePipelines_basic(t *testing.T) {
	var (
		pipeline       matlas.DataLakePipeline
		resourceName   = "mongodbatlas_data_lake_pipeline.test"
		dataSourceName = "data.mongodbatlas_data_lake_pipelines.testDataSource"
		clusterName    = acctest.RandomWithPrefix("test-acc-index")
		projectID      = "63f4d4a47baeac59406dc131"
		firstPipelineName           = acctest.RandomWithPrefix("test-acc-index")
		secondPipelineName           = acctest.RandomWithPrefix("test-acc-index")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasDataLakeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasDataLakePipelinesConfig(projectID, clusterName, firstPipelineName, secondPipelineName),
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

func testAccDataSourceMongoDBAtlasDataLakePipelinesConfig(projectID, clusterName, firstPipelineName, secondPipelineName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "aws_conf" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 10
		
			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_EAST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}
			backup_enabled               = false
			auto_scaling_disk_gb_enabled = false
		
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
		}

		resource "mongodbatlas_cluster" "aws_conf2" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 10
		
			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_EAST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}
			backup_enabled               = false
			auto_scaling_disk_gb_enabled = false
		
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
		}

		resource "mongodbatlas_data_lake_pipeline" "test" {
			project_id       =  "%[1]s"
			name = "%[3]s"
			sink {
				type = "DLS"
				partition_fields {
						name = "access"
						order = 0
				}
			}	
	
			source {
				type = "ON_DEMAND_CPS"
				cluster_name = mongodbatlas_cluster.aws_conf.name
				cluster_name = "Cluster0"
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
			name			 = " "%[4]s"
			sink {
				type = "DLS"
				partition_fields {
						name = "access"
						order = 0
				}
			}	
	
			source {
				type = "ON_DEMAND_CPS"
				name			 = cluster_name = mongodbatlas_cluster.aws_conf2.name
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
	`, projectID, clusterName, firstPipelineName, secondPipelineName)
}
