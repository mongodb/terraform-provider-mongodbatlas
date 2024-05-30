package federatedsettingsidentityprovider_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsIdentityProvidersDS_basic(t *testing.T) {
	var (
		dataSourceName      = "data.mongodbatlas_federated_settings_identity_providers.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configPluralDS(federatedSettingsID, nil, []string{oidcProtocol, samlProtocol}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "2"),
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, nil, []string{samlProtocol}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.display_name", "SAML-test"),
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, nil, []string{oidcProtocol}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.display_name", "OIDC-test"),
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, nil, []string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.display_name", "SAML-test"), // if no protocol is specified, it defaults to SAML
				),
			},
			{
				Config: configPluralDS(federatedSettingsID, conversion.StringPtr("WORKLOAD"), []string{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "federation_settings_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "0"),
				),
			},
		},
	})
}

func configPluralDS(federatedSettingsID string, idpType *string, protocols []string) string {
	var protocolString string
	if len(protocols) > 1 {
		protocolString = fmt.Sprintf(`protocols = [%[1]q, %[2]q]`, protocols[0], protocols[1])
	} else if len(protocols) > 0 {
		protocolString = fmt.Sprintf(`protocols = [%[1]q]`, protocols[0])
	}
	var idpTypeString string
	if idpType != nil {
		idpTypeString = fmt.Sprintf(`idp_types = [%[1]q]`, *idpType)
	}

	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_identity_providers" "test" {
			federation_settings_id = "%[1]s"
			page_num = 1
			items_per_page = 100
			%[2]s
			%[3]s
		}
`, federatedSettingsID, protocolString, idpTypeString)
}
