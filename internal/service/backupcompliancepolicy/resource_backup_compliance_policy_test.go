package backupcompliancepolicy_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName       = "mongodbatlas_backup_compliance_policy.backup_policy_res"
	dataSourceName     = "data.mongodbatlas_backup_compliance_policy.backup_policy"
	projectIDTerraform = "mongodbatlas_project.test.id"
)

func TestAccBackupCompliancePolicy_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t, true))
}

func TestAccBackupCompliancePolicy_update(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid conflicts with backup compliance policy
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithoutOptionals(projectName, orgID, projectOwnerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_last_name", "Last"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "false"),
				),
			},
			{
				Config: configBasic(projectName, orgID, projectOwnerID, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_last_name", "Last"),
				),
			},
		},
	})
}

func TestAccBackupCompliancePolicy_overwriteBackupPolicies(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid conflicts with backup compliance policy
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		req            = acc.ClusterRequest{
			AdvancedConfiguration: map[string]any{
				acc.ClusterAdvConfigOplogMinRetentionHours: 8,
			},
			ProjectID:            projectIDTerraform,
			MongoDBMajorVersion:  "6.0",
			CloudBackup:          true,
			DiskSizeGb:           12,
			RetainBackupsEnabled: true,
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{EbsVolumeType: "STANDARD", AutoScalingDiskGbEnabled: true, NodeCount: 3},
			},
		}
		clusterInfo = acc.GetClusterInfo(t, &req)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configClusterWithBackupSchedule(projectName, orgID, projectOwnerID, &clusterInfo),
			},
			{
				Config:      configOverwriteIncompatibleBackupPoliciesError(projectName, orgID, projectOwnerID, &clusterInfo),
				ExpectError: regexp.MustCompile(`BACKUP_POLICIES_NOT_MEETING_BACKUP_COMPLIANCE_POLICY_REQUIREMENTS`),
			},
		},
	})
}

func TestAccBackupCompliancePolicy_withoutRestoreWindowDaysAndOnDemand(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid conflicts with backup compliance policy
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithoutRestoreDaysAndOnDemand(projectName, orgID, projectOwnerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_policy_item.#", "0"),
				),
			},
		},
	})
}

func TestAccBackupCompliancePolicy_UpdateSetsAllAttributes(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid conflicts with backup compliance policy
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasicWithOptionalAttributesWithNonDefaultValues(projectName, orgID, projectOwnerID, "7"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_last_name", "Last"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "on_demand_policy_item.0.retention_value", "3"),
				),
			},
			{
				Config: configBasicWithOptionalAttributesWithNonDefaultValues(projectName, orgID, projectOwnerID, "8"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_first_name", "First"),
					resource.TestCheckResourceAttr(resourceName, "authorized_user_last_name", "Last"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption_at_rest_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "copy_protection_enabled", "true"),
				),
			},
		},
	})
}

func basicTestCase(tb testing.TB, useYearly bool) *resource.TestCase {
	tb.Helper()

	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid conflicts with backup compliance policy
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, projectOwnerID, useYearly),
				Check:  resource.ComposeAggregateTestCheckFunc(basicChecks()...),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"state"},
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
		ids := conversion.DecodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		policy, _, err := acc.ConnV2().CloudBackupsApi.GetDataProtectionSettings(context.Background(), projectID).Execute()
		if err != nil || policy == nil {
			return fmt.Errorf("backup compliance policy (%s) does not exist: %s", rs.Primary.ID, err)
		}
		time.Sleep(30 * time.Second) // Wait for the bcp to be fully applied, see more details in CLOUDP-324378.
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

func configBasic(projectName, orgID, projectOwnerID string, useYearly bool) string {
	var strYearly string
	if useYearly {
		strYearly = `
			policy_item_yearly {
				frequency_interval = 1
				retention_unit     = "years"
				retention_value    = 1
			}
		`
	}

	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) +
		fmt.Sprintf(`	  
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

			%s
	  }
		
		data "mongodbatlas_backup_compliance_policy" "backup_policy" {
			project_id = mongodbatlas_backup_compliance_policy.backup_policy_res.project_id
		}
	`, strYearly)
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

			policy_item_yearly {
				frequency_interval = 1
				retention_unit     = "years"
				retention_value    = 1
			}
	  }
	`
}

func configWithoutRestoreDaysAndOnDemand(projectName, orgID, projectOwnerID string) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) + `	  
	  resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
			project_id                 = mongodbatlas_project.test.id
			authorized_email           = "test@example.com"
			authorized_user_first_name = "First"
			authorized_user_last_name  = "Last"
			copy_protection_enabled    = false
			pit_enabled                = false
			encryption_at_rest_enabled = false
			
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

			policy_item_yearly {
				frequency_interval = 1
				retention_unit     = "years"
				retention_value    = 1
			}
	  }
	`
}

func configOverwriteIncompatibleBackupPoliciesError(projectName, orgID, projectOwnerID string, info *acc.ClusterInfo) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) + fmt.Sprintf(`	  
	  %[1]s
	  resource "mongodbatlas_cloud_backup_schedule" "test" {
		cluster_name 			   = %[2]s.name
		project_id                 = mongodbatlas_project.test.id
	  
		reference_hour_of_day    = 3
		reference_minute_of_hour = 45
		restore_window_days      = 2
	  
		copy_settings {
		  cloud_provider      = "AWS"
		  frequencies         = ["DAILY"]
		  region_name         = "US_WEST_1"
		  replication_spec_id = one(%[2]s.replication_specs).id
		  should_copy_oplogs  = false
		}
	  }

	  resource "mongodbatlas_backup_compliance_policy" "test" {
		project_id                 = mongodbatlas_project.test.id
		authorized_email           = "test@example.com"
		  authorized_user_first_name = "First"
		  authorized_user_last_name  = "Last"
		copy_protection_enabled    = true
		pit_enabled                = false
		encryption_at_rest_enabled = false
	  
		on_demand_policy_item {
		  frequency_interval = 1
		  retention_unit     = "days"
		  retention_value    = 1
		}
	  
		policy_item_daily {
		  frequency_interval = 1
		  retention_unit     = "days"
		  retention_value    = 1
		}
	  }
	`, info.TerraformStr, info.ResourceName)
}

func configClusterWithBackupSchedule(projectName, orgID, projectOwnerID string, info *acc.ClusterInfo) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) + fmt.Sprintf(`	  
	  %[1]s
	  resource "mongodbatlas_cloud_backup_schedule" "test" {
		cluster_name 			  = %[2]s.name
		project_id                 = mongodbatlas_project.test.id
	  
		reference_hour_of_day    = 3
		reference_minute_of_hour = 45
		restore_window_days      = 2
	  
		copy_settings {
		  cloud_provider      = "AWS"
		  frequencies         = ["DAILY"]
		  region_name         = "US_WEST_1"
		  replication_spec_id = one(%[2]s.replication_specs).id
		  should_copy_oplogs  = false
		}
	  }
	`, info.TerraformStr, info.ResourceName)
}

func basicChecks() []resource.TestCheckFunc {
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

func configBasicWithOptionalAttributesWithNonDefaultValues(projectName, orgID, projectOwnerID, restreWindowDays string) string {
	return acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false) +
		fmt.Sprintf(`resource "mongodbatlas_backup_compliance_policy" "backup_policy_res" {
		project_id                 = mongodbatlas_project.test.id
		authorized_email           = "test@example.com"
		authorized_user_first_name = "First"
		authorized_user_last_name  = "Last"
		copy_protection_enabled    = true
		pit_enabled                = false
		encryption_at_rest_enabled = false
		
		restore_window_days = %[1]s
		
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
  }`, restreWindowDays)
}
