package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceClusterRSDataLakePipeline_basic(t *testing.T) {
	var (
		pipeline     matlas.DataLakePipeline
		resourceName = "mongodbatlas_data_lake_pipeline.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-index")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name         = acctest.RandomWithPrefix("test-acc-index")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDataLakeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasDataLakePipelineConfig(projectID, clusterName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDataLakePipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasDataLakeDestroy(s *terraform.State) error {
	conn := testAccProviderSdkV2.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_data_lake" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		// Try to find the database user
		_, _, err := conn.DataLakes.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			return fmt.Errorf("datalake (%s) still exists", ids["project_id"])
		}
	}

	return nil
}

func testAccDataSourceMongoDBAtlasDataLakePipelineConfig(projectID, clusterName, pipelineName string) string {
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

		resource "mongodbatlas_data_lake_pipeline" "test" {
			project_id       = "%[1]s"
			name			 = "%[3]s"
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

		data "mongodbatlas_data_lake_pipeline" "testDataSource" {
			project_id       = mongodbatlas_data_lake_pipeline.test.project_id
			name			 = mongodbatlas_data_lake_pipeline.test.name	
		}
	`, projectID, clusterName, pipelineName)
}
