package orgmaintenancesettings_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_org_maintenance_settings.test"
	dataSourceName = "data.mongodbatlas_org_maintenance_settings.test"
)

func TestAccOrgMaintenanceSettings_basic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithMode(orgID, "MANUAL"),
				Check:  checkWithMode(orgID, "MANUAL"),
			},
			{
				Config: configWithMode(orgID, "ENV_TAG_MAPPING"),
				Check:  checkWithMode(orgID, "ENV_TAG_MAPPING"),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    importStateIDFunc(resourceName),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "org_id",
			},
			{
				Config: configEmpty(orgID),
				Check:  checkEmpty(orgID),
			},
			{
				Config:   configEmpty(orgID),
				PlanOnly: true, // Check that wave_assignment_mode was correctly unset in API.
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

		data "mongodbatlas_org_maintenance_settings" "test" {
			org_id = mongodbatlas_org_maintenance_settings.test.org_id
		}
	`, orgID, waveAssignmentMode)
}

func configEmpty(orgID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_org_maintenance_settings" "test" {
			org_id = %[1]q
		}

		data "mongodbatlas_org_maintenance_settings" "test" {
			org_id = mongodbatlas_org_maintenance_settings.test.org_id
		}
	`, orgID)
}

func checkWithMode(orgID, waveAssignmentMode string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
		resource.TestCheckResourceAttr(resourceName, "wave_assignment_mode", waveAssignmentMode),
		resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
		resource.TestCheckResourceAttr(dataSourceName, "wave_assignment_mode", waveAssignmentMode),
		resource.TestCheckResourceAttrSet(dataSourceName, "effective_wave_assignment_mode"),
	)
}

func checkEmpty(orgID string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
		resource.TestCheckNoResourceAttr(resourceName, "wave_assignment_mode"),
		resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
		resource.TestCheckNoResourceAttr(dataSourceName, "wave_assignment_mode"),
		resource.TestCheckResourceAttrSet(dataSourceName, "effective_wave_assignment_mode"),
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
