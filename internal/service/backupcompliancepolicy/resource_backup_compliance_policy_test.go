package backupcompliancepolicy_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const resourceName = "mongodbatlas_backup_compliance_policy.backup_policy_res"

func TestAccGenericBackupRSBackupCompliancePolicy_basic(t *testing.T) {
	var (
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_last_name", "Last"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
				),
			},
		},
	})
}

func TestAccGenericBackupRSBackupCompliancePolicy_withFirstLastName(t *testing.T) {
	var (
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_last_name", "Last"),
				),
			},
		},
	})
}

func TestAccGenericBackupRSBackupCompliancePolicy_withoutOptionals(t *testing.T) {
	var (
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithoutOptionals(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_last_name", "Last"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "false"),
				),
			},
		},
	})
}

func TestAccGenericBackupRSBackupCompliancePolicy_withoutRestoreWindowDays(t *testing.T) {
	var (
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithoutRestoreDays(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
				),
			},
		},
	})
}

func TestAccGenericBackupRSBackupCompliancePolicy_importBasic(t *testing.T) {
	var (
		projectName    = fmt.Sprintf("testacc-project-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, projectOwnerID),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"state"},
			},
		},
	})
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
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		schedule, _, err := acc.Conn().BackupCompliancePolicy.Get(context.Background(), projectID)
		if err != nil || schedule == nil {
			return fmt.Errorf("backup compliance policy (%s) does not exist: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_backup_compliance_policy" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		compliancePolicy, _, err := acc.Conn().BackupCompliancePolicy.Get(context.Background(), projectID)
		if compliancePolicy != nil || err == nil {
			return fmt.Errorf("Backup Compliance Policy (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return ids["project_id"], nil
	}
}

func configBasic(projectName, orgID, projectOwnerID string) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) + `	  
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

func configWithoutOptionals(projectName, orgID, projectOwnerID string) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) + `	  
	  resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = mongodbatlas_project.test.id
			authorized_email           = "test@example.com"
			authorized_user_first_name = "First"
			authorized_user_last_name  = "Last"
			
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

func configWithoutRestoreDays(projectName, orgID, projectOwnerID string) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) + `	  
	  resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = mongodbatlas_project.test.id
			authorized_email           = "test@example.com"
			authorized_user_first_name = "First"
			authorized_user_last_name  = "Last"
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
	`
}
