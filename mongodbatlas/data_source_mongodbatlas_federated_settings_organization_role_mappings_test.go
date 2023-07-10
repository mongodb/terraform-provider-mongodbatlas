package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFedDSFederatedSettingsOrganizationRoleMappings_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName        = "data.mongodbatlas_federated_settings_org_role_mappings.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID               = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testCheckFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceFederatedSettingsOrganizationRoleMappingsConfig(federatedSettingsID, orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingsExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.external_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.role_assignments.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceFederatedSettingsOrganizationRoleMappingsConfig(federatedSettingsID, orgID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_org_role_mappings" "test" {
			federation_settings_id = "%[1]s"
			org_id                 = "%[2]s"
			page_num = 1
			items_per_page = 100
		}
`, federatedSettingsID, orgID)
}

func testAccCheckMongoDBAtlasFederatedSettingsOrganizationRoleMappingsExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		_, _, err := conn.FederatedSettings.ListRoleMappings(context.Background(), rs.Primary.Attributes["federation_settings_id"], rs.Primary.Attributes["org_id"], nil)
		if err != nil {
			return fmt.Errorf("FederatedSettingsOrganizationRoleMappings (%s) does not exist", rs.Primary.ID)
		}

		return nil
	}
}
