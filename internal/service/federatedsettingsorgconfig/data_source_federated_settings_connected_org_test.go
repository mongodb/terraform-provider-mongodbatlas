package federatedsettingsorgconfig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsOrgDS_basic(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		resourceName        = "data.mongodbatlas_federated_settings_org_config.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID               = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasicDS(federatedSettingsID, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_mappings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "identity_provider_id"),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttr(resourceName, "identity_provider_id", "0oad4fas87jL5Xnk1297"),
				),
			},
		},
	})
}

func configBasicDS(federatedSettingsID, orgID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_org_config" "test" {
			federation_settings_id = "%[1]s"
			org_id = "%[2]s"

		}
`, federatedSettingsID, orgID)
}
