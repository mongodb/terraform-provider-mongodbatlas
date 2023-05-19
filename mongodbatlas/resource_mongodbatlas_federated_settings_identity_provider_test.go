package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccFedRSFederatedSettingsIdentityProvider_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		federatedSettingsIdentityProvider matlas.FederatedSettingsIdentityProvider
		resourceName                      = "mongodbatlas_federated_settings_identity_provider.test"
		federationSettingsID              = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		idpID                             = os.Getenv("MONGODB_ATLAS_FEDERATED_OKTA_IDP_ID")
		ssoURL                            = os.Getenv("MONGODB_ATLAS_FEDERATED_SSO_URL")
		issuerURI                         = os.Getenv("MONGODB_ATLAS_FEDERATED_ISSUER_URI")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testCheckFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:            testAccMongoDBAtlasFederatedSettingsIdentityProviderConfig(federationSettingsID, ssoURL, issuerURI),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderImportStateIDFunc(resourceName, federationSettingsID, idpID),
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config:            testAccMongoDBAtlasFederatedSettingsIdentityProviderConfig(federationSettingsID, ssoURL, issuerURI),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderImportStateIDFunc(resourceName, federationSettingsID, idpID),

				ImportState: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderExists(resourceName, &federatedSettingsIdentityProvider, idpID),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "name", "mongodb_federation_test"),
				),
			},
		},
	})
}

func TestAccFedRSFederatedSettingsIdentityProvider_importBasic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName         = "mongodbatlas_federated_settings_identity_provider.test"
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		idpID                = os.Getenv("MONGODB_ATLAS_FEDERATED_OKTA_IDP_ID")
		ssoURL               = os.Getenv("MONGODB_ATLAS_FEDERATED_SSO_URL")
		issuerURI            = os.Getenv("MONGODB_ATLAS_FEDERATED_ISSUER_URI")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testCheckFederatedSettings(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config:            testAccMongoDBAtlasFederatedSettingsIdentityProviderConfig(federationSettingsID, ssoURL, issuerURI),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderImportStateIDFunc(resourceName, federationSettingsID, idpID),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderExists(resourceName string,
	federatedSettingsIdentityProvider *matlas.FederatedSettingsIdentityProvider, idpID string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		response, _, err := conn.FederatedSettings.GetIdentityProvider(context.Background(),
			rs.Primary.Attributes["federation_settings_id"],
			idpID)
		if err == nil {
			*federatedSettingsIdentityProvider = *response
			return nil
		}

		return fmt.Errorf("identity provider (%s) does not exist", idpID)
	}
}

func testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderImportStateIDFunc(resourceName, federationSettingsID, idpID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ID := encodeStateID(map[string]string{
			"federation_settings_id": federationSettingsID,
			"okta_idp_id":            idpID,
		})

		ids := decodeStateID(ID)
		return fmt.Sprintf("%s-%s", ids["federation_settings_id"], ids["okta_idp_id"]), nil
	}
}

func testAccMongoDBAtlasFederatedSettingsIdentityProviderConfig(federationSettingsID, ssoURL, issuerURI string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_identity_provider" "test" {
		federation_settings_id = "%[1]s"
		name = "mongodb_federation_test"
        associated_domains           = ["reorganizeyourworld.com"]
        sso_debug_enabled = true
        status = "ACTIVE"
        sso_url = "%[2]s"
        issuer_uri = "%[3]s"
        request_binding = "HTTP-POST"
        response_signature_algorithm = "SHA-256"
	  }`, federationSettingsID, ssoURL, issuerURI)
}
