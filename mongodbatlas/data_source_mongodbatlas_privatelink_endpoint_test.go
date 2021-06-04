package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasPrivateLinkEndpoint_basic(t *testing.T) {
	SkipTestExtCred(t)
	resourceName := "data.mongodbatlas_privatelink_endpoint.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointDataSourceConfig(projectID, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateLinkEndpointDataSourceConfig(projectID, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%s"
			provider_name = "%s"
			region        = "%s"
		}

		data "mongodbatlas_privatelink_endpoint" "test" {
			project_id      = mongodbatlas_privatelink_endpoint.test.project_id
			private_link_id = mongodbatlas_privatelink_endpoint.test.id
			provider_name = "%[2]s"
		}
	`, projectID, providerName, region)
}
