package mongodbatlas

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoDBAtlasPrivateEndpointRegionalMode_basic(t *testing.T) {
	datasourceName := "data.mongodbatlas_private_endpoint_regional_mode.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	enabled := true

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeDataSourceConfig(projectID, enabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceADLExists(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "enabled", strconv.FormatBool(enabled)),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeDataSourceConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id = "%[1]s"
			enabled    = "%[2]b"
		}

		data "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id      = mongodbatlas_private_endpoint_regional_mode.test.project_id
			endpoint_id = mongodbatlas_private_endpoint_regional_mode.test.enabled
		}
	`, projectID, enabled)
}
