package privateendpointregionalmode_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccNetworkDSPrivateEndpointRegionalMode_basic(t *testing.T) {
	resourceName := "mongodbatlas_private_endpoint_regional_mode.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
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
		conn := acc.TestAccProviderSdkV2.Meta().(*config.MongoDBClient).Atlas

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
