package federatedsettingsorgconfig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsOrgDSPlural_basic(t *testing.T) {
	var (
		resourceName        = "data.mongodbatlas_federated_settings_org_configs.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasicPluralDS(federatedSettingsID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.identity_provider_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.data_access_identity_provider_ids.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.org_id"),
				),
			},
		},
	})
}

func configBasicPluralDS(federatedSettingsID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_org_configs" "test" {
			federation_settings_id = "%[1]s"
			page_num = 1
			items_per_page = 100
		}
`, federatedSettingsID)
}
