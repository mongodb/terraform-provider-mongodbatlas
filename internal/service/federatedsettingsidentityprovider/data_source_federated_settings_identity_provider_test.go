package federatedsettingsidentityprovider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsIdentityProviderDS_samlBasic(t *testing.T) {
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
				Config: configBasicDS(federatedSettingsID, idpID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "associated_orgs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "acs_url"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "SAML-test"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "SAML"),
					resource.TestCheckResourceAttrSet(resourceName, "okta_idp_id"),
					resource.TestCheckResourceAttr(resourceName, "idp_id", idpID),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federatedSettingsID),
				),
			},
		},
	})
}

func configBasicDS(federatedSettingsID, idpID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_identity_provider" "test" {
			federation_settings_id = "%[1]s"
            identity_provider_id   = "%[2]s"
		}
`, federatedSettingsID, idpID)
}
