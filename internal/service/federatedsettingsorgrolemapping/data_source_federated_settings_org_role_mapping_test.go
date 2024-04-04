package federatedsettingsorgrolemapping_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsOrgRoleMappingDS_basic(t *testing.T) {
	acc.SkipTestForCI(t) // affects the org

	var (
		resourceName        = "data.mongodbatlas_federated_settings_org_role_mapping.test"
		federatedSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID               = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
		roleMappingID       = os.Getenv("MONGODB_ATLAS_FEDERATED_ROLE_MAPPING_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasicDS(federatedSettingsID, orgID, roleMappingID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "federation_settings_id"),
					resource.TestCheckResourceAttrSet(resourceName, "external_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "role_assignments.#"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "external_group_name", "group2"),
				),
			},
		},
	})
}

func configBasicDS(federatedSettingsID, orgID, roleMappingID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_federated_settings_org_role_mapping" "test" {
			federation_settings_id = "%[1]s"
			org_id                 = "%[2]s"
			role_mapping_id        = "%[3]s"
		}
`, federatedSettingsID, orgID, roleMappingID)
}
