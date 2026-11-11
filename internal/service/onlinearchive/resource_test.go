package onlinearchive_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/onlinearchive"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clusterRequest() *acc.ClusterRequest {
	return &acc.ClusterRequest{
		ReplicationSpecs: []acc.ReplicationSpecRequest{
			// Must use US_EAST_1 in dev for online_archive to work
			{AutoScalingDiskGbEnabled: true, Region: "US_EAST_1"},
		},
		Tags: map[string]string{
			"ArchiveTest": "true", "Owner": "test",
		},
	}
}
func TestAccBackupRSOnlineArchive(t *testing.T) {
	var (
		onlineArchiveResourceName    = "mongodbatlas_online_archive.users_archive"
		onlineArchiveDataSourceName  = "data.mongodbatlas_online_archive.read_archive"
		onlineArchivesDataSourceName = "data.mongodbatlas_online_archives.all"
		clusterInfo                  = acc.GetClusterInfo(t, clusterRequest())
		clusterName                  = clusterInfo.Name
		projectID                    = clusterInfo.ProjectID
		clusterTerraformStr          = clusterInfo.TerraformStr
		clusterResourceName          = clusterInfo.ResourceName
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: clusterTerraformStr,
				Check:  acc.PopulateWithSampleDataTestCheck(projectID, clusterName),
			},
			{
				Config: configWithDailySchedule(clusterTerraformStr, clusterResourceName, 1, 7, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_minute"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_minute"),
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_expiration_rule.0.expire_after_days", "7"),
					resource.TestCheckResourceAttr(onlineArchiveDataSourceName, "data_expiration_rule.0.expire_after_days", "7"),
					resource.TestCheckResourceAttr(onlineArchivesDataSourceName, "results.0.data_expiration_rule.0.expire_after_days", "7"),
				),
			},
			{
				Config: configWithDailySchedule(clusterTerraformStr, clusterResourceName, 2, 8, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_minute"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_minute"),
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_expiration_rule.0.expire_after_days", "8"),
					resource.TestCheckResourceAttr(onlineArchiveDataSourceName, "data_expiration_rule.0.expire_after_days", "8"),
					resource.TestCheckResourceAttr(onlineArchivesDataSourceName, "results.0.data_expiration_rule.0.expire_after_days", "8"),
				),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithWeeklySchedule(clusterTerraformStr, clusterResourceName, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_minute"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_minute"),
				),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithMonthlySchedule(clusterTerraformStr, clusterResourceName, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_minute"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_minute"),
				),
			},
			{
				Config: configWithoutSchedule(clusterTerraformStr, clusterResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
					resource.TestCheckNoResourceAttr(onlineArchiveResourceName, "schedule.#"),
				),
			},
			{
				Config: configWithoutSchedule(clusterTerraformStr, clusterResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "partition_fields.0.field_name", "last_review"),
				),
			},
		},
	})
}

func TestAccBackupRSOnlineArchiveBasic(t *testing.T) {
	var (
		clusterInfo               = acc.GetClusterInfo(t, clusterRequest())
		clusterResourceName       = clusterInfo.ResourceName
		clusterName               = clusterInfo.Name
		projectID                 = clusterInfo.ProjectID
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		clusterTerraformStr       = clusterInfo.TerraformStr
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: clusterTerraformStr,
				Check:  acc.PopulateWithSampleDataTestCheck(projectID, clusterName),
			},
			{
				Config: configWithoutSchedule(clusterTerraformStr, clusterResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
				),
			},
			{
				// The archive's backing federated instance is only listed with type=ONLINE_ARCHIVE.
				// No exact count: the project is shared with parallel archive tests.
				Config: configWithoutSchedule(clusterTerraformStr, clusterResourceName) + configFederatedArchiveInstances(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mongodbatlas_federated_database_instances.archive", "results.0.name"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_federated_database_instances.archive", "results.0.hostnames.0"),
				),
			},
			{
				Config: configWithDailySchedule(clusterTerraformStr, clusterResourceName, 1, 1, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.type"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.end_minute"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_hour"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "schedule.0.start_minute"),
				),
			},
			{
				ResourceName:            onlineArchiveResourceName,
				ImportStateIdFunc:       importStateIDFunc(onlineArchiveResourceName, "/"),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sync_creation", "delete_on_create_timeout", "state"},
			},
			{
				ResourceName:            onlineArchiveResourceName,
				ImportStateIdFunc:       importStateIDFunc(onlineArchiveResourceName, "-"),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"sync_creation", "delete_on_create_timeout", "state"},
			},
		},
	})
}

// importStateIDFunc builds the import ID {project_id}{sep}{cluster_name}{sep}{archive_id}. Test cluster names contain hyphens.
func importStateIDFunc(resourceName, sep string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		attrs := rs.Primary.Attributes
		return strings.Join([]string{attrs["project_id"], attrs["cluster_name"], attrs["archive_id"]}, sep), nil
	}
}

func TestAccBackupRSOnlineArchiveWithProcessRegion(t *testing.T) {
	var (
		onlineArchiveResourceName   = "mongodbatlas_online_archive.users_archive"
		onlineArchiveDataSourceName = "data.mongodbatlas_online_archive.read_archive"
		clusterInfo                 = acc.GetClusterInfo(t, clusterRequest())
		clusterResourceName         = clusterInfo.ResourceName
		clusterName                 = clusterInfo.Name
		projectID                   = clusterInfo.ProjectID
		clusterTerraformStr         = clusterInfo.TerraformStr
		cloudProvider               = "AWS"
		processRegion               = "US_EAST_1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: clusterTerraformStr,
				Check:  acc.PopulateWithSampleDataTestCheck(projectID, clusterName),
			},
			{
				Config: configWithDataProcessRegion(clusterTerraformStr, clusterResourceName, cloudProvider, processRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.region", processRegion),
					resource.TestCheckResourceAttr(onlineArchiveDataSourceName, "data_process_region.0.cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(onlineArchiveDataSourceName, "data_process_region.0.region", processRegion),
				),
			},
			{
				Config:      configWithDataProcessRegion(clusterTerraformStr, clusterResourceName, cloudProvider, "AP_SOUTH_1"),
				ExpectError: regexp.MustCompile("data_process_region can't be modified"),
			},
			{
				Config: configWithoutSchedule(clusterTerraformStr, clusterResourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.region", processRegion),
				),
			},
		},
	})
}

func TestAccBackupRSOnlineArchive_ErrorMessages(t *testing.T) {
	var (
		clusterInfo         = acc.GetClusterInfo(t, clusterRequest())
		clusterTerraformStr = clusterInfo.TerraformStr
		cloudProvider       = "AWS"
		clusterResourceName = clusterInfo.ResourceName
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configWithInvalidDateFormat(clusterTerraformStr, clusterResourceName),
				ExpectError: regexp.MustCompile("An invalid enumeration value INVALID_FORMAT was specified"),
			},
			{
				Config:      configWithDataProcessRegion(clusterTerraformStr, clusterResourceName, cloudProvider, "UNKNOWN"),
				ExpectError: regexp.MustCompile("INVALID_ATTRIBUTE"),
			},
		},
	})
}

func configWithDailySchedule(clusterTerraformStr, clusterResourceName string, startHour, deleteExpirationDays int, deleteOnTimeout bool) string {
	var dataExpirationRuleBlock string
	if deleteExpirationDays > 0 {
		dataExpirationRuleBlock = fmt.Sprintf(`
		data_expiration_rule {
			expire_after_days = %d
		}
		`, deleteExpirationDays)
	}
	deleteOnCreateTimeoutStr := ""
	if deleteOnTimeout {
		deleteOnCreateTimeoutStr = `
		delete_on_create_timeout = true
		timeouts {
			create = "1s"
		}
		`
	}

	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = %[4]s.project_id
		cluster_name = %[4]s.name
		coll_name = "listingsAndReviews"
		collection_type = "STANDARD"
		db_name = "sample_airbnb"
	
		criteria {
			type = "DATE"
			date_field = "last_review"
			date_format = "ISODATE"
			expire_after_days = 2
		}
		
		%[3]s

		schedule {
			type = "DAILY"
			end_hour = 1
			end_minute = 1
			start_hour = %[2]d
			start_minute = 1
		}

		partition_fields {
			field_name = "last_review"
			order = 0
		}
	
		partition_fields {
			field_name = "maximum_nights"
			order = 1
		}
	
		partition_fields {
			field_name = "name"
			order = 2
		}

		sync_creation = true

		%[5]s
	}
	
	data "mongodbatlas_online_archive" "read_archive" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
		archive_id = mongodbatlas_online_archive.users_archive.archive_id
	}
	
	data "mongodbatlas_online_archives" "all" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
	}
	`, clusterTerraformStr, startHour, dataExpirationRuleBlock, clusterResourceName, deleteOnCreateTimeoutStr)
}

func configWithoutSchedule(clusterTerraformStr, clusterResourceName string) string {
	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = %[2]s.project_id
		cluster_name = %[2]s.name
		coll_name = "listingsAndReviews"
		collection_type = "STANDARD"
		db_name = "sample_airbnb"
	
		criteria {
			type = "DATE"
			date_field = "last_review"
			date_format = "ISODATE"
			expire_after_days = 2
		}

		partition_fields {
			field_name = "last_review"
			order = 0
		}

		partition_fields {
			field_name = "maximum_nights"
			order = 1
		}
	
		partition_fields {
			field_name = "name"
			order = 2
		}

		sync_creation = true
	}
	
	data "mongodbatlas_online_archive" "read_archive" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
		archive_id = mongodbatlas_online_archive.users_archive.archive_id
	}
	
	data "mongodbatlas_online_archives" "all" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
	}
	`, clusterTerraformStr, clusterResourceName)
}

func configFederatedArchiveInstances() string {
	return `
	data "mongodbatlas_federated_database_instances" "archive" {
		project_id = mongodbatlas_online_archive.users_archive.project_id
		type       = "ONLINE_ARCHIVE"
	}
	`
}

func configWithDataProcessRegion(clusterTerraformStr, clusterResourceName, cloudProvider, region string) string {
	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = %[4]s.project_id
		cluster_name = %[4]s.name
		coll_name = "listingsAndReviews"
		collection_type = "STANDARD"
		db_name = "sample_airbnb"

		criteria {
			type = "DATE"
			date_field = "last_review"
			date_format = "ISODATE"
			expire_after_days = 2
		}

		partition_fields {
			field_name = "last_review"
			order = 0
		}

		partition_fields {
			field_name = "maximum_nights"
			order = 1
		}

		partition_fields {
			field_name = "name"
			order = 2
		}

		data_process_region {
			cloud_provider = %[2]q
			region = %[3]q
		}

		sync_creation = true
	}

	data "mongodbatlas_online_archive" "read_archive" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
		archive_id = mongodbatlas_online_archive.users_archive.archive_id
	}

	data "mongodbatlas_online_archives" "all" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
	}
	`, clusterTerraformStr, cloudProvider, region, clusterResourceName)
}

func testAccBackupRSOnlineArchiveConfigWithWeeklySchedule(clusterTerraformStr, clusterResourceName string, startHour int) string {
	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = %[3]s.project_id
		cluster_name = %[3]s.name
		coll_name = "listingsAndReviews"
		collection_type = "STANDARD"
		db_name = "sample_airbnb"
	
		criteria {
			type = "DATE"
			date_field = "last_review"
			date_format = "ISODATE"
			expire_after_days = 2
		}

		schedule {
			type = "WEEKLY"
			day_of_week = 1
			end_hour = 1
			end_minute = 1
			start_hour = %[2]d
			start_minute = 1
		}

		partition_fields {
			field_name = "last_review"
			order = 0
		}
	
		partition_fields {
			field_name = "maximum_nights"
			order = 1
		}
	
		partition_fields {
			field_name = "name"
			order = 2
		}

		sync_creation = true
	}
	
	data "mongodbatlas_online_archive" "read_archive" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
		archive_id = mongodbatlas_online_archive.users_archive.archive_id
	}
	
	data "mongodbatlas_online_archives" "all" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
	}
	`, clusterTerraformStr, startHour, clusterResourceName)
}

func testAccBackupRSOnlineArchiveConfigWithMonthlySchedule(clusterTerraformStr, clusterResourceName string, startHour int) string {
	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = %[3]s.project_id
		cluster_name = %[3]s.name
		coll_name = "listingsAndReviews"
		collection_type = "STANDARD"
		db_name = "sample_airbnb"
	
		criteria {
			type = "DATE"
			date_field = "last_review"
			date_format = "ISODATE"
			expire_after_days = 2
		}

		schedule {
			type = "MONTHLY"
			day_of_month = 1
			end_hour = 1
			end_minute = 1
			start_hour = %[2]d
			start_minute = 1
		}

		partition_fields {
			field_name = "last_review"
			order = 0
		}
	
		partition_fields {
			field_name = "maximum_nights"
			order = 1
		}


		partition_fields {
			field_name = "name"
			order = 2
		}


		sync_creation = true
	}
	
	data "mongodbatlas_online_archive" "read_archive" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
		archive_id = mongodbatlas_online_archive.users_archive.archive_id
	}
	
	data "mongodbatlas_online_archives" "all" {
		project_id =  mongodbatlas_online_archive.users_archive.project_id
		cluster_name = mongodbatlas_online_archive.users_archive.cluster_name
	}
	`, clusterTerraformStr, startHour, clusterResourceName)
}

func TestAccOnlineArchive_deleteOnCreateTimeout(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, clusterRequest())
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: clusterInfo.TerraformStr,
				Check:  acc.PopulateWithSampleDataTestCheck(clusterInfo.ProjectID, clusterInfo.Name),
			},
			{
				Config:      configWithDailySchedule(clusterInfo.TerraformStr, clusterInfo.ResourceName, 1, 7, true),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	})
}

func configWithInvalidDateFormat(clusterTerraformStr, clusterResourceName string) string {
	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = %[2]s.project_id
		cluster_name = %[2]s.name
		coll_name = "listingsAndReviews"
		collection_type = "STANDARD"
		db_name = "sample_airbnb"
	
		criteria {
			type = "DATE"
			date_field = "last_review"
			date_format = "INVALID_FORMAT"
			expire_after_days = 2
		}

		partition_fields {
			field_name = "last_review"
			order = 0
		}
	}
	`, clusterTerraformStr, clusterResourceName)
}

func TestSplitOnlineArchiveImportID(t *testing.T) {
	const errorString = onlinearchive.ErrorImportFormat
	tests := map[string]struct {
		importID    string
		projectID   string
		clusterName string
		archiveID   string
		errorString string
	}{
		"slash format": {
			importID:    "664619d870c247237f4b86a6/my-cluster/664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my-cluster",
			archiveID:   "664619d870c247237f4b86a7",
		},
		"slash format with more parts": {
			importID:    "part1/part2/part3/part4",
			errorString: errorString,
		},
		"cluster name without hyphens": {
			importID:    "664619d870c247237f4b86a6-Cluster0-664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "Cluster0",
			archiveID:   "664619d870c247237f4b86a7",
		},
		"cluster name with one hyphen": {
			importID:    "664619d870c247237f4b86a6-my-cluster-664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my-cluster",
			archiveID:   "664619d870c247237f4b86a7",
		},
		"cluster name with several hyphens": {
			importID:    "664619d870c247237f4b86a6-my-prod-cluster-01-664619d870c247237f4b86a7",
			projectID:   "664619d870c247237f4b86a6",
			clusterName: "my-prod-cluster-01",
			archiveID:   "664619d870c247237f4b86a7",
		},
		"missing archive ID": {
			importID:    "664619d870c247237f4b86a6-Cluster0",
			errorString: errorString,
		},
		"empty input": {
			importID:    "",
			errorString: errorString,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, clusterName, archiveID, err := onlinearchive.SplitOnlineArchiveImportID(tc.importID)
			if tc.errorString != "" {
				require.EqualError(t, err, tc.errorString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.clusterName, clusterName)
			assert.Equal(t, tc.archiveID, archiveID)
		})
	}
}
