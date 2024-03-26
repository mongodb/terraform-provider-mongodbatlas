package backupcompliancepolicy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_backup_compliance_policy.backup_policy_res"
	dataSourceName = "data.mongodbatlas_backup_compliance_policy.backup_policy"
)

func TestAccBackupCompliancePolicy_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID),
				Check:  resource.ComposeTestCheckFunc(checks()...),
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

func TestAccBackupCompliancePolicy_withoutOptionals(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithoutOptionals(projectID),
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

func TestAccBackupCompliancePolicy_withoutRestoreWindowDays(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithoutRestoreDays(projectID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
				),
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
		policy, _, err := acc.ConnV2().CloudBackupsApi.GetDataProtectionSettings(context.Background(), projectID).Execute()
		if err != nil || policy == nil {
			return fmt.Errorf("backup compliance policy (%s) does not exist: %s", rs.Primary.ID, err)
		}
		return nil
	}
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

func configBasic(projectID string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = %q
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

		data "mongodbatlas_backup_compliance_policy" "backup_policy" {
			project_id = mongodbatlas_backup_compliance_policy.backup_policy_res.project_id
		}
	`, projectID)
}

func configWithoutOptionals(projectID string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = %q
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
		`, projectID)
}

func configWithoutRestoreDays(projectID string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = %q
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
	`, projectID)
}

func checks() []resource.TestCheckFunc {
	commonChecks := map[string]string{
		"copy_protection_enabled":    "false",
		"encryption_at_rest_enabled": "false",
		"authorized_user_first_name": "First",
		"authorized_user_last_name":  "Last",
		"authorized_email":           "test@example.com",
		"restore_window_days":        "7",
	}
	checks := acc.AddAttrChecks(resourceName, nil, commonChecks)
	checks = acc.AddAttrChecks(dataSourceName, checks, commonChecks)
	checks = append(checks, checkExists(resourceName), checkExists(dataSourceName))
	return checks
}
