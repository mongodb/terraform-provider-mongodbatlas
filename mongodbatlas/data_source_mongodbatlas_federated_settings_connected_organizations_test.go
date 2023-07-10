package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFedDSFederatedSettingsOrganizationConfigs_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName        = "data.mongodbatlas_federated_settings_org_configs.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testCheckFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceFederatedSettingsOrganizationConfigsConfig(federatedSettingsID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsOrganizationConfigsExists(resourceName),

					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.identity_provider_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.org_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceFederatedSettingsOrganizationConfigsConfig(federatedSettingsID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_org_configs" "test" {
			federation_settings_id = "%[1]s"
			page_num = 1
			items_per_page = 100
		}
`, federatedSettingsID)
}

func testAccCheckMongoDBAtlasFederatedSettingsOrganizationConfigsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		_, _, err := conn.FederatedSettings.ListConnectedOrgs(context.Background(), rs.Primary.Attributes["federation_settings_id"], nil)
		if err != nil {
			return fmt.Errorf("FederatedSettingsConnectedOrganization (%s) does not exist", rs.Primary.ID)
		}

		return nil
	}
}
