package datalakepipeline_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataLakePipeline_basic(t *testing.T) {
	var (
		pipeline     matlas.DataLakePipeline
		resourceName = "mongodbatlas_data_lake_pipeline.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-index")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		name         = acctest.RandomWithPrefix("test-acc-index")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, clusterName, name),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, &pipeline),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
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
		// Try to find the data lake pipeline
		_, _, err := acc.Conn().DataLakePipeline.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			return fmt.Errorf("datalake (%s) still exists", ids["project_id"])
		}
	}
	return nil
}

func checkExists(resourceName string, pipeline *matlas.DataLakePipeline) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		response, _, err := acc.Conn().DataLakePipeline.Get(context.Background(), ids["project_id"], ids["name"])
		if err == nil {
			*pipeline = *response
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
				region_name   = "EU_WEST_1"
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
