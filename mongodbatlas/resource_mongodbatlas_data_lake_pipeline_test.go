package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccClusterRSDataLakePipeline_basic(t *testing.T) {
	var (
		pipeline     matlas.DataLakePipeline
		resourceName = "mongodbatlas_data_lake_pipeline.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-index")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name         = acctest.RandomWithPrefix("test-acc-index")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakePipelineConfig(projectID, clusterName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDataLakePipelineExists(resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasDataLakePipelineImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasDataLakePipelineImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s", ids["project_id"], ids["name"]), nil
	}
}

func testAccCheckMongoDBAtlasDataLakePipelineExists(resourceName string, pipeline *matlas.DataLakePipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		response, _, err := conn.DataLakePipeline.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			*pipeline = *response
			return nil
		}

		return fmt.Errorf("DataLake pipeline (%s) does not exist", ids["name"])
	}
}

func testAccMongoDBAtlasDataLakePipelineConfig(projectID, clusterName, pipelineName string) string {
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

		resource "mongodbatlas_data_lake_pipeline" "test" {
			project_id       = "%[1]s"
			name			 = "%[3]s"
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
				database_name = "sample_airbnb"
				collection_name = "listingsAndReviews"
			}

			transformations {
				field = "test"
				type =  "EXCLUDE"
			}
		}
	`, projectID, clusterName, pipelineName)
}
