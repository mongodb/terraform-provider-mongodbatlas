package federatedsettingsorgconfig_test

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

func TestAccFederatedSettingsOrg_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // affects the org

	var (
		resourceName         = "mongodbatlas_federated_settings_org_config.test"
		federationSettingsID = os.Getenv("MONGODB_ATLAS_FEDERATION_SETTINGS_ID")
		orgID                = os.Getenv("MONGODB_ATLAS_FEDERATED_ORG_ID")
		idpID                = os.Getenv("MONGODB_ATLAS_FEDERATED_IDP_ID")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckFederatedSettings(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:            configBasic(federationSettingsID, orgID, idpID),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName, federationSettingsID, orgID),
				ImportState:       true,
				ImportStateVerify: false,
			},
			{
				Config:            configBasic(federationSettingsID, orgID, idpID),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName, federationSettingsID, orgID),

				ImportState: true,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "federation_settings_id", federationSettingsID),
					resource.TestCheckResourceAttr(resourceName, "name", "mongodb_federation_test"),
				),
			},
			{
				Config:            configBasic(federationSettingsID, orgID, idpID),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName, federationSettingsID, orgID),
				ImportState:       true,
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
		_, _, err := acc.Conn().FederatedSettings.GetConnectedOrg(context.Background(),
			rs.Primary.Attributes["federation_settings_id"],
			rs.Primary.Attributes["org_id"])
		if err == nil {
			return nil
		}
		return fmt.Errorf("connected org  (%s) does not exist", rs.Primary.Attributes["org_id"])
	}
}

func importStateIDFunc(resourceName, federationSettingsID, orgID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ID := conversion.EncodeStateID(map[string]string{
			"federation_settings_id": federationSettingsID,
			"org_id":                 orgID,
		})

		ids := conversion.DecodeStateID(ID)
		return fmt.Sprintf("%s-%s", ids["federation_settings_id"], ids["org_id"]), nil
	}
}

func configBasic(federationSettingsID, orgID, identityProviderID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_federated_settings_org_config" "test" {
		federation_settings_id = "%[1]s"
		org_id                 = "%[2]s"
		domain_restriction_enabled = false
		domain_allow_list          = ["reorganizeyourworld.com"]
		identity_provider_id = "%[3]s"
	  }`, federationSettingsID, orgID, identityProviderID)
}
