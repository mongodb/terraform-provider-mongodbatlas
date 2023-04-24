package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBackupRSBackupCompliancePolicy_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_backup_compliance_policy.backup_policy_res"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupCompliancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
				),
			},
			{
				Config: testAccMongoDBAtlasBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists(resourceName),
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists("data.mongodbatlas_backup_compliance_policy.backup_policy"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
				),
			},
		},
	})
}

func TestAccBackupRSBackupCompliancePolicy_withoutRestoreWindowDays(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_backup_compliance_policy.backup_policy_res"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupCompliancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupCompliancePolicyConfigWithoutRestoreDays(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists(resourceName),
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists("data.mongodbatlas_backup_compliance_policy.backup_policy"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "encryption_at_rest_enabled", "false"),
				),
			},
			{
				Config: testAccMongoDBAtlasBackupCompliancePolicyConfigWithoutRestoreDays(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists(resourceName),
					testAccCheckMongoDBAtlasBackupCompliancePolicyExists("data.mongodbatlas_backup_compliance_policy.backup_policy"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr("mongodbatlas_backup_compliance_policy.backup_policy_res", "encryption_at_rest_enabled", "false"),
				),
			},
		},
	})
}

func TestAccBackupRSBackupCompliancePolicy_importBasic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_backup_compliance_policy.backup_policy_res"
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupCompliancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasBackupCompliancePolicyImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func testAccCheckMongoDBAtlasBackupCompliancePolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("N o ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]

		schedule, _, err := conn.BackupCompliancePolicy.Get(context.Background(), projectID)
		if err != nil || schedule == nil {
			return fmt.Errorf("Backup Compliance Policy (%s) does not exist: %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasBackupCompliancePolicyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_backup_compliance_policy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]

		compliancePolicy, _, err := conn.BackupCompliancePolicy.Get(context.Background(), projectID)
		if compliancePolicy != nil || err == nil {
			return fmt.Errorf("Backup Compliance Policy (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasBackupCompliancePolicyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return ids["project_id"], nil
	}
}

func testAccMongoDBAtlasBackupCompliancePolicyConfig(projectName, orgID, projectOwnerID string) string {
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
		authorized_email           = "test@example.com"
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
	`, projectName, orgID, projectOwnerID)
}

func testAccMongoDBAtlasBackupCompliancePolicyConfigWithoutRestoreDays(projectName, orgID, projectOwnerID string) string {
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
		authorized_email           = "test@example.com"
		copy_protection_enabled    = false
		pit_enabled                = false
		encryption_at_rest_enabled = false
	  
		//restore_window_days = 7
	  
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
	`, projectName, orgID, projectOwnerID)
}
