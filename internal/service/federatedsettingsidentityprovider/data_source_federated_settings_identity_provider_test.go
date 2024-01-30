package federatedsettingsidentityprovider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFedDSFederatedSettingsIdentityProvider_samlBasic(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		resourceName        = "data.mongodbatlas_federated_settings_identity_provider.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		idpID               = os.Getenv("MONGODB_ATLAS_FEDERATED_IDP_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceFederatedSettingsIdentityProviderConfig(federatedSettingsID, idpID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsIdentityProvidersExists(resourceName),

					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_orgs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "acs_url"),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "TestConfig"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "SAML"),
					resource.TestCheckResourceAttr(resourceName, "okta_idp_id", "0oafbloyfixJjK4VI357"),
					resource.TestCheckResourceAttr(resourceName, "idp_id", idpID),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federatedSettingsID),
				),
			},
		},
	})
}

func TestAccFedDSFederatedSettingsIdentityProvider_oidcBasic(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		resourceName        = "data.mongodbatlas_federated_settings_identity_provider.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		idpID               = os.Getenv("MONGODB_ATLAS_FEDERATED_OIDC_IDP_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceFederatedSettingsIdentityProviderConfig(federatedSettingsID, idpID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsIdentityProvidersExists(resourceName),

					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_orgs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "audience_claim.#"),
					resource.TestCheckResourceAttrSet(resourceName, "client_id"),
					resource.TestCheckResourceAttrSet(resourceName, "groups_claim"),
					resource.TestCheckResourceAttrSet(resourceName, "requested_scopes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "user_claim"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "OIDC"),
					resource.TestCheckResourceAttr(resourceName, "okta_idp_id", ""),
					resource.TestCheckResourceAttr(resourceName, "idp_id", idpID),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federatedSettingsID),
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
