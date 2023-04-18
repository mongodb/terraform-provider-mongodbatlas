package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBackupDSBackupCompliancePolicy_basic(t *testing.T) {
	projectName := fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectOwnerID := os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists("data.mongodbatlas_backup_compliance_policy.backup_policy"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "encryption_at_rest_enabled", "false"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name                                             = "%s"
			org_id                                           = "%s"
			project_owner_id                                 = "%s"
			with_default_alerts_settings                     = false
			is_collect_database_specifics_statistics_enabled = false
			is_data_explorer_enabled                         = false
			is_performance_advisor_enabled                   = false
			is_realtime_performance_panel_enabled            = false
			is_schema_advisor_enabled                        = false
		  }
		  
		  data "mongodbatlas_backup_compliance_policy" "backup_policy" {
			project_id = mongodbatlas_backup_compliance_policy.backup_policy_res.project_id
		  }
		  
		  resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = mongodbatlas_project.test.id
			authorized_email           = "martin@akava.io"
			copy_protection_enabled    = false
			pit_enabled                = false
			encryption_at_rest_enabled = false
		  
			restore_window_days = 7
		  
			on_demand_policy_item {
		  
			  frequency_interval = 0
			  frequency_type     = "ondemand"
			  retention_unit     = "days"
			  retention_value    = 3
			}
			
			scheduled_policy_items {
			  
				frequency_interval = 6
				frequency_type     = "hourly"
				retention_unit     = "days"
				retention_value    = 7
			  }
		  
		  scheduled_policy_items {
			  
				frequency_interval = 0
				frequency_type     = "daily"
				retention_unit     = "days"
				retention_value    = 7
			  }
		  
		  scheduled_policy_items {
			  
				frequency_interval = 0
				frequency_type     = "weekly"
				retention_unit     = "weeks"
				retention_value    = 4
			  }
		  
		  scheduled_policy_items {
			  
				frequency_interval = 0
				frequency_type     = "monthly"
				retention_unit     = "months"
				retention_value    = 12
			  }
		  
		  }
`, projectName, orgID, projectOwnerID)
}
