package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceMongoDBAtlasPrivateEndpointRegionalMode_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_private_endpoint_regional_mode.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(projectID, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateEndpointRegionalModeConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id   = %[1]q
			enabled      = %[2]t
		}
	`, projectID, enabled)
}

func testAccCheckMongoDBAtlasPrivateEndpointRegionalModeExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		projectID := rs.Primary.ID

		_, _, err := conn.PrivateEndpoints.GetRegionalizedPrivateEndpointSetting(context.Background(), projectID)

		if err == nil {
			return nil
		}

		return fmt.Errorf("regional mode for project_id (%s) does not exist", projectID)
	}
}
