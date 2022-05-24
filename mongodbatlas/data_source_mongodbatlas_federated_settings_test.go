package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceMongoDBAtlasFederatedSettings_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_cloud_federated_settings.config"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		name           = "Terraform Official Testing for Federation"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasFederatedSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasFederatedSettingsConfig(orgID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "identity_provider_id"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasFederatedSettingsConfig(orgID, name string) string {
	return fmt.Sprintf(`
	data "mongodbatlas_cloud_federated_settings" "federated_settings" {
		org_id = "%s"
		name = "%s"
	  }

	`, orgID, name)
}

func testAccCheckMongoDBAtlasFederatedSettingsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_federated_settings" {
			continue
		}

		// Try to find the cluster
		globalConfig, _, err := conn.FederatedSettings.Get(context.Background(), rs.Primary.Attributes["org_id"])
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("No federated settings identity provider %s exists in org %s", rs.Primary.Attributes["identity_provider_id"], rs.Primary.Attributes["org_id"])) {
				return nil
			}

			return err
		}

		if len(globalConfig.IdentityProviderID) > 0 || len(globalConfig.IdentityProviderStatus) > 0 {
			return fmt.Errorf("Federated settings identity provider(%s) still exists", rs.Primary.Attributes["identity_provider_id"])
		}
	}

	return nil
}
