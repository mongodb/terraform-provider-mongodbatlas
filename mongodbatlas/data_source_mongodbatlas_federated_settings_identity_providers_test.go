package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoDBAtlasFederatedSettingsIdentityProvider_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName        = "data.mongodbatlas_cloud_federated_settings_identity_providers.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { checkFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		//CheckDestroy:      testAccCheckMongoDBAtlasFederatedSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceFederatedSettingsIdentityProviderConfig(federatedSettingsID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderExists(resourceName),

					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.acs_url"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.display_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceFederatedSettingsIdentityProviderConfig(federatedSettingsID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_cloud_federated_settings_identity_providers" "test" {
			federation_settings_id = "%[1]s"
			page_num = 1
			items_per_page = 100
		}
`, federatedSettingsID)
}

func testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		_, _, err := conn.FederatedSettings.ListIdentityProviders(context.Background(), rs.Primary.Attributes["federation_settings_id"], nil)
		if err != nil {
			return fmt.Errorf("FederatedSettingsIdentityProviders (%s) does not exist", rs.Primary.ID)
		}

		return nil
	}
}
