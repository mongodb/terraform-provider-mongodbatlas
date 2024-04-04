package federatedsettingsorgrolemapping_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsOrgRoleMappingDSPlural_basic(t *testing.T) {
	var (
		resourceName        = "data.mongodbatlas_federated_settings_org_role_mappings.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID               = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasicPluralDS(federatedSettingsID, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.external_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.role_assignments.#"),
				),
			},
		},
	})
}

func configBasicPluralDS(federatedSettingsID, orgID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_org_role_mappings" "test" {
			federation_settings_id = "%[1]s"
			org_id                 = "%[2]s"
			page_num = 1
			items_per_page = 100
		}
`, federatedSettingsID, orgID)
}
