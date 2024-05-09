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
		resourceName     = "mongodbatlas_federated_settings_org_role_mapping.test"
		dataSourceName   = "data.mongodbatlas_federated_settings_org_role_mapping.test"
		dataSourcePlural = "data.mongodbatlas_federated_settings_org_role_mappings.test"

		extGroupName1        = "newtestgroup"
		extGroupName2        = "newupdatedgroup"
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID                = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
		groupID              = os.Getenv("MONGODB_ATLAS_FEDERATED_GROUP_ID")
		mapAttrs             = map[string]string{
			"federation_settings_id": federationSettingsID,
			"org_id":                 orgID,
			"external_group_name":    extGroupName1,
		}
		sliceAttrs       = []string{"role_assignments.#"}
		sliceAttrsPlural = []string{"results.#", "federation_settings_id"}
	)
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks = acc.AddAttrChecks(resourceName, checks, mapAttrs)
	checks = acc.AddAttrChecks(dataSourceName, checks, mapAttrs)
	checks = acc.AddAttrSetChecks(resourceName, checks, sliceAttrs...)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, sliceAttrs...)
	checks = acc.AddAttrSetChecks(dataSourcePlural, checks, sliceAttrsPlural...)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(federationSettingsID, orgID, groupID, extGroupName1),

				Check: resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: configBasic(federationSettingsID, orgID, groupID, extGroupName2),

				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(resourceName, "external_group_name", extGroupName2)),
			},
			{
				Config:            configBasic(federationSettingsID, orgID, groupID, extGroupName2),
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
		_, _, err := acc.ConnV2().FederatedAuthenticationApi.GetRoleMapping(context.Background(),
			rs.Primary.Attributes["federation_settings_id"],
			rs.Primary.Attributes["role_mapping_id"],
			rs.Primary.Attributes["org_id"],
		).Execute()
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
		roleMapping, _, err := acc.ConnV2().FederatedAuthenticationApi.GetRoleMapping(context.Background(), ids["federation_settings_id"], ids["role_mapping_id"], ids["org_id"]).Execute()
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

func configBasic(federationSettingsID, orgID, groupID, externalGroupName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_org_role_mapping" "test" {
		federation_settings_id = "%[1]s"
		org_id                 = "%[2]s"
		external_group_name    = %[4]q
		role_assignments {
			org_id = "%[2]s"
			roles  = ["ORG_MEMBER","ORG_GROUP_CREATOR"]
		}
		
		  role_assignments {
			group_id = "%[3]s"
			roles    = ["GROUP_OWNER","GROUP_DATA_ACCESS_ADMIN","GROUP_SEARCH_INDEX_EDITOR","GROUP_DATA_ACCESS_READ_ONLY"]
		}

	}
	data "mongodbatlas_federated_settings_org_role_mapping" "test" {
		federation_settings_id = "%[1]s"
		org_id                 = "%[2]s"
		role_mapping_id        = mongodbatlas_federated_settings_org_role_mapping.test.role_mapping_id
	}
	data "mongodbatlas_federated_settings_org_role_mappings" "test" {
		depends_on 			   = [mongodbatlas_federated_settings_org_role_mapping.test]

		federation_settings_id = "%[1]s"
		org_id                 = "%[2]s"
		page_num = 1
		items_per_page = 100
	}
	  `, federationSettingsID, orgID, groupID, externalGroupName)
}
