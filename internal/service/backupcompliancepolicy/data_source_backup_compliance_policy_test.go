package backupcompliancepolicy_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccGenericBackupDSBackupCompliancePolicy_basic(t *testing.T) {
	projectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectOwnerID := os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists("data.mongodbatlas_backup_compliance_policy.backup_policy"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "authorized_email", "test@example.com"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "authorized_user_last_name", "Last"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID string) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) + `

		data "mongodbatlas_backup_compliance_policy" "backup_policy" {
			project_id = mongodbatlas_backup_compliance_policy.backup_policy_res.project_id
		}
		  
		resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = mongodbatlas_project.test.id
			authorized_email           = "test@example.com"
			authorized_user_first_name = "First"
			authorized_user_last_name  = "Last"
			copy_protection_enabled    = false
			pit_enabled                = false
			encryption_at_rest_enabled = false
		  
			restore_window_days = 7
		  
			on_demand_policy_item {
			  frequency_interval = 0
			  retention_unit     = "days"
			  retention_value    = 3
			}
			
			policy_item_hourly {
				frequency_interval = 6
				retention_unit     = "days"
				retention_value    = 7
			}
		  
			policy_item_daily {
				frequency_interval = 0
				retention_unit     = "days"
				retention_value    = 7
			}
		  
			policy_item_weekly {
				frequency_interval = 0
				retention_unit     = "weeks"
				retention_value    = 4
			}
		  
			policy_item_monthly {
				frequency_interval = 0
				retention_unit     = "months"
				retention_value    = 12
			}
		}
	`
}
