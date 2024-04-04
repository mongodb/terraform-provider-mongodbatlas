package federatedsettingsorgrolemapping_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccFederatedSettingsOrgRoleMapping_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // affects the org

	var (
		resourceName         = "mongodbatlas_federated_settings_org_role_mapping.test"
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID                = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
		groupID              = os.Getenv("MONGODB_ATLAS_FEDERATED_GROUP_ID")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(federationSettingsID, orgID, groupID),

				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "external_group_name", "newtestgroup"),
				),
			},
			{
				Config:            configBasic(federationSettingsID, orgID, groupID),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       false,
				ImportStateVerify: false,
			},
		},
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		_, _, err := acc.Conn().FederatedSettings.GetRoleMapping(context.Background(),
			rs.Primary.Attributes["federation_settings_id"],
			rs.Primary.Attributes["org_id"],
			rs.Primary.Attributes["role_mapping_id"])
		if err == nil {
			return nil
		}
		return fmt.Errorf("role mapping (%s) does not exist", rs.Primary.Attributes["role_mapping_id"])
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_federated_settings_org_role_mapping" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		roleMapping, _, err := acc.Conn().FederatedSettings.GetRoleMapping(context.Background(), ids["federation_settings_id"], ids["org_id"], ids["role_mapping_id"])
		if err == nil && roleMapping != nil {
			return fmt.Errorf("role mapping (%s) still exists", ids["okta_idp_id"])
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["federation_settings_id"], ids["org_id"], ids["role_mapping_id"]), nil
	}
}

func configBasic(federationSettingsID, orgID, groupID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_org_role_mapping" "test" {
		federation_settings_id = "%[1]s"
		org_id                 = "%[2]s"
		external_group_name    = "newtestgroup"
		role_assignments {
			org_id = "%[2]s"
			roles  = ["ORG_MEMBER","ORG_GROUP_CREATOR"]
		}
		
		  role_assignments {
			group_id = "%[3]s"
			roles    = ["GROUP_OWNER","GROUP_DATA_ACCESS_ADMIN","GROUP_SEARCH_INDEX_EDITOR","GROUP_DATA_ACCESS_READ_ONLY"]
		}

	  }`, federationSettingsID, orgID, groupID)
}
