package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccFedDSFederatedSettings_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		federatedSettings matlas.FederatedSettings
		resourceName      = "data.mongodbatlas_federated_settings.test"
		orgID             = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testCheckFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceFederatedSettingsConfig(orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsExists(resourceName, &federatedSettings),

					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "identity_provider_id"),
					resource.TestCheckResourceAttrSet(resourceName, "identity_provider_status"),
					resource.TestCheckResourceAttrSet(resourceName, "has_role_mappings"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceFederatedSettingsConfig(orgID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings" "test" {
			org_id = "%[1]s"
		}
`, orgID)
}

func testAccCheckMongoDBAtlasFederatedSettingsExists(resourceName string, federatedSettings *matlas.FederatedSettings) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		federatedSettingsRes, _, err := conn.FederatedSettings.Get(context.Background(), rs.Primary.Attributes["org_id"])
		if err != nil {
			return fmt.Errorf("FederatedSettings (%s) does not exist", rs.Primary.ID)
		}

		federatedSettings = federatedSettingsRes

		return nil
	}
}
