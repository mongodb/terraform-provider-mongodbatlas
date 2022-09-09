package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceMongoDBAtlasPrivateEndpointRegionalMode_basic(t *testing.T) {
	resourceName := "mongodbatlas_private_endpoint_regional_mode.test"
	projectID := os.Getenv("MONGODB_ATLAS_NETWORK_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeDataSourceConfig(projectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					testAccMongoDBAtlasPrivateEndpointRegionalModeUnmanagedResource(resourceName, projectID),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeUnmanagedResource(resourceName, projectID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		setting, _, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), projectID)

		if err != nil || setting == nil {
			return fmt.Errorf("Could not get regionalized private endpoint setting for project_id (%s)", projectID)
		}

		return resource.TestCheckResourceAttr(resourceName, "enabled", strconv.FormatBool(setting.Enabled))(s)
	}
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeDataSourceConfig(projectID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id = data.mongodbatlas_private_endpoint_regional_mode.test.project_id
		}

		data "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id  = %q
		}
	`, projectID)
}
