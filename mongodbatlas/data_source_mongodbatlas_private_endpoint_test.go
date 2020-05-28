package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasPrivateEndpoint_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_private_endpoint.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); checkPeeringEnvAWS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointDataSourceConfig(projectID, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
				),
			},
		},
	})

}

func testAccMongoDBAtlasPrivateEndpointDataSourceConfig(projectID, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint" "test" {
			project_id    = "%s"
			provider_name = "%s"
			region        = "%s"
		}

		data "mongodbatlas_private_endpoint" "test" {
			project_id      = "${mongodbatlas_private_endpoint.test.project_id}"
			private_link_id = "${mongodbatlas_private_endpoint.test.private_link_id}"
		}
	`, projectID, providerName, region)
}
