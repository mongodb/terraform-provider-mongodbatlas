package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoDBAtlasFederatedSettingsIdentityProvider_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName        = "data.mongodbatlas_federated_settings_identity_provider.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		idpID               = os.Getenv("MONGODB_ATLAS_FEDERATED_IDP_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { checkFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceFederatedSettingsIdentityProviderConfig(federatedSettingsID, idpID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsIdentityProvidersExists(resourceName),

					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_orgs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "acs_url"),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "mongodb_federation_test"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceFederatedSettingsIdentityProviderConfig(federatedSettingsID, idpID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_identity_provider" "test" {
			federation_settings_id = "%[1]s"
            identity_provider_id   = "%[2]s"
		}
`, federatedSettingsID, idpID)
}
