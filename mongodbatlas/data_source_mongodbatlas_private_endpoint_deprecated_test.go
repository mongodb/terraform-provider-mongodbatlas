package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasPrivateEndpointDeprecated_basic(t *testing.T) {
	SkipTestExtCred(t)
	resourceName := "data.mongodbatlas_private_endpoint_deprecated.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); checkPeeringEnvAWS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointDeprecatedDataSourceConfig(projectID, providerName, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointDeprecatedExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateEndpointDeprecatedDataSourceConfig(projectID, providerName, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_deprecated" "test" {
			project_id    = "%s"
			provider_name = "%s"
			region        = "%s"
		}

		data "mongodbatlas_private_endpoint_deprecated" "test" {
			project_id      = "${mongodbatlas_private_endpoint_deprecated.test.project_id}"
			private_link_id = "${mongodbatlas_private_endpoint_deprecated.test.private_link_id}"
		}
	`, projectID, providerName, region)
}
