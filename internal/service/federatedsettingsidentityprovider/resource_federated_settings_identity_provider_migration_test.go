package federatedsettingsidentityprovider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func TestAccMigration_FederatedSettingsIdentityProvider(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		federatedSettingsIdentityProvider admin.FederationIdentityProvider
		resourceName                      = "mongodbatlas_federated_settings_identity_provider.test"
		federationSettingsID              = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		idpID                             = os.Getenv("MONGODB_ATLAS_FEDERATED_OKTA_IDP_ID")
		ssoURL                            = os.Getenv("MONGODB_ATLAS_FEDERATED_SSO_URL")
		issuerURI                         = os.Getenv("MONGODB_ATLAS_FEDERATED_ISSUER_URI")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { acc.PreCheckFederatedSettings(t) },
		// ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:            testAccMongoDBAtlasFederatedSettingsIdentityProviderConfig(federationSettingsID, ssoURL, issuerURI),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderImportStateIDFunc(resourceName, federationSettingsID, idpID),
				ImportState:       true,
				ImportStateVerify: true,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderExists(resourceName, &federatedSettingsIdentityProvider, idpID),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "name", "mongodb_federation_test"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   testAccMongoDBAtlasFederatedSettingsIdentityProviderConfig(federationSettingsID, ssoURL, issuerURI),
				// ImportStateIdFunc:        testAccCheckMongoDBAtlasFederatedSettingsIdentityProviderImportStateIDFunc(resourceName, federationSettingsID, idpID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
