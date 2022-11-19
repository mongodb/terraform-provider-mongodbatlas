package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoDBAtlasPrivateLinkEndpointsServiceServerless_basic(t *testing.T) {
	datasourceName := "data.mongodbatlas_privatelink_endpoints_service_serverless.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	instanceID := "serverless"
	comments := "Test Comments"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointsServerlessDataSourceConfig(projectID, instanceID, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "instance_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateLinkEndpointsServerlessDataSourceConfig(projectID, instanceID, comments string) string {
	return fmt.Sprintf(`
	data "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
		project_id   = "%[1]s"
		instance_name = mongodbatlas_serverless_instance.test.name
		endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
	  }

	data "mongodbatlas_privatelink_endpoints_service_serverless" "test" {
	  project_id   = "%[1]s"
	  instance_name = mongodbatlas_serverless_instance.test.name
	}

	resource "mongodbatlas_privatelink_endpoint_serverless" "test" {
		project_id   = "%[1]s"
		instance_name = mongodbatlas_serverless_instance.test.name
		provider_name = "AWS"
	  }
	  
	  
	  resource "mongodbatlas_privatelink_endpoint_service_serverless" "test" {
		project_id   = "%[1]s"
		instance_name = "%[2]s"
		endpoint_id = mongodbatlas_privatelink_endpoint_serverless.test.endpoint_id
		provider_name = "AWS"
		comment = "%[3]s"
	  }

	resource "mongodbatlas_serverless_instance" "test" {
		project_id   = "%[1]s"
		name         = "%[2]s"
		provider_settings_backing_provider_name = "AWS"
		provider_settings_provider_name = "SERVERLESS"
		provider_settings_region_name = "US_EAST_1"
		continuous_backup_enabled = true
	}
	`, projectID, instanceID, comments)
}
