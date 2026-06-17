package orgmaintenancesettings_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_org_maintenance_settings.test"

func TestAccOrgMaintenanceSettings_basic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithMode(orgID, "MANUAL"),
				Check:  checkWithMode("MANUAL"),
			},
			{
				Config: configWithMode(orgID, "ENV_TAG_MAPPING"),
				Check:  checkWithMode("ENV_TAG_MAPPING"),
			},
			{
				Config: configEmpty(orgID),
				Check:  checkEmpty(),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "org_id",
			},
		},
	})
}

func configWithMode(orgID, waveAssignmentMode string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_maintenance_settings" "test" {
			org_id               = %[1]q
			wave_assignment_mode = %[2]q
		}
	`, orgID, waveAssignmentMode)
}

func configEmpty(orgID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_maintenance_settings" "test" {
			org_id = %[1]q
		}
	`, orgID)
}

func checkWithMode(waveAssignmentMode string) resource.TestCheckFunc {
	checks := acc.AddAttrChecks(resourceName, nil, map[string]string{
		"wave_assignment_mode": waveAssignmentMode,
	})
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkEmpty() resource.TestCheckFunc {
	// When wave_assignment_mode is removed from config, the provider sends null to the API which resets it to the default (MANUAL).
	// The field may be omitted from GET responses when maintenance sequencing is disabled for the org.
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "org_id"),
	)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		orgID := rs.Primary.Attributes["org_id"]
		if orgID == "" {
			return "", fmt.Errorf("import, org_id not found for: %s", resourceName)
		}
		return orgID, nil
	}
}
