package onlinearchive_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	clusterRequest = acc.ClusterRequest{
		ReplicationSpecs: []acc.ReplicationSpecRequest{
			// Must use US_EAST_1 in dev for online_archive to work
			{AutoScalingDiskGbEnabled: true, Region: "US_EAST_1"},
		},
		Tags: map[string]string{
			"ArchiveTest": "true", "Owner": "test",
		},
	}
)

func TestAccBackupRSOnlineArchive(t *testing.T) {
	var (
		onlineArchiveResourceName    = "mongodbatlas_online_archive.users_archive"
		onlineArchiveDataSourceName  = "data.mongodbatlas_online_archive.read_archive"
		onlineArchivesDataSourceName = "data.mongodbatlas_online_archives.all"
		clusterInfo                  = acc.GetClusterInfo(t, &clusterRequest)
		clusterName                  = clusterInfo.ClusterName
		projectID                    = clusterInfo.ProjectID
		clusterTerraformStr          = clusterInfo.ClusterTerraformStr
		clusterResourceName          = clusterInfo.ClusterResourceName
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: clusterTerraformStr,
				Check: resource.ComposeAggregateTestCheckFunc(
					populateWithSampleData(clusterResourceName, projectID, clusterName),
				),
			},
			{
				Config: configWithDailySchedule(clusterTerraformStr, clusterResourceName, 1, 7),
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
				Config: configWithDailySchedule(clusterTerraformStr, clusterResourceName, 2, 8),
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
		clusterInfo               = acc.GetClusterInfo(t, &clusterRequest)
		clusterResourceName       = clusterInfo.ClusterResourceName
		clusterName               = clusterInfo.ClusterName
		projectID                 = clusterInfo.ProjectID
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		clusterTerraformStr       = clusterInfo.ClusterTerraformStr
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: clusterTerraformStr,
				Check: resource.ComposeAggregateTestCheckFunc(
					populateWithSampleData(clusterResourceName, projectID, clusterName),
				),
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
				Config: configWithDailySchedule(clusterTerraformStr, clusterResourceName, 1, 1),
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
		},
	})
}

func TestAccBackupRSOnlineArchiveWithProcessRegion(t *testing.T) {
	var (
		onlineArchiveResourceName   = "mongodbatlas_online_archive.users_archive"
		onlineArchiveDataSourceName = "data.mongodbatlas_online_archive.read_archive"
		clusterInfo                 = acc.GetClusterInfo(t, &clusterRequest)
		clusterResourceName         = clusterInfo.ClusterResourceName
		clusterName                 = clusterInfo.ClusterName
		projectID                   = clusterInfo.ProjectID
		clusterTerraformStr         = clusterInfo.ClusterTerraformStr
		cloudProvider               = "AWS"
		processRegion               = "US_EAST_1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: clusterTerraformStr,
				Check: resource.ComposeAggregateTestCheckFunc(
					populateWithSampleData(clusterResourceName, projectID, clusterName),
				),
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

func TestAccBackupRSOnlineArchiveInvalidProcessRegion(t *testing.T) {
	var (
		clusterInfo         = acc.GetClusterInfo(t, &clusterRequest)
		clusterTerraformStr = clusterInfo.ClusterTerraformStr
		cloudProvider       = "AWS"
		clusterResourceName = clusterInfo.ClusterResourceName
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configWithDataProcessRegion(clusterTerraformStr, clusterResourceName, cloudProvider, "UNKNOWN"),
				ExpectError: regexp.MustCompile("INVALID_ATTRIBUTE"),
			},
		},
	})
}

// populateWithSampleData adds Sample Data to the cluster otherwise online archive won't work
func populateWithSampleData(resourceName, projectID, clusterName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		conn := acc.ConnV2()
		ctx := context.Background()
		_, _, err := conn.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
		if err != nil {
			return fmt.Errorf("cluster(%s:%s) does not exist %s", projectID, clusterName, err)
		}
		job, _, err := conn.ClustersApi.LoadSampleDataset(context.Background(), projectID, clusterName).Execute()
		if err != nil {
			return fmt.Errorf("cluster(%s:%s) loading sample data set error %s", projectID, clusterName, err)
		}
		if job == nil {
			return fmt.Errorf("cluster(%s:%s) loading sample data set error, no job found", projectID, clusterName)
		}
		ticker := time.NewTicker(30 * time.Second)

	JOB:
		for {
			select {
			case <-time.After(20 * time.Second):
				log.Println("timeout elapsed ....")
			case <-ticker.C:
				job, _, err = conn.ClustersApi.GetSampleDatasetLoadStatus(ctx, projectID, job.GetId()).Execute()
				fmt.Println("querying for job ")
				if err != nil {
					return fmt.Errorf("cluster(%s:%s) failed to query for job, %s", projectID, clusterName, err)
				}
				if job == nil {
					return fmt.Errorf("cluster(%s:%s) failed to query for job, no job found", projectID, clusterName)
				}
				if job.GetState() != "WORKING" {
					break JOB
				}
			}
		}

		if job.GetState() != "COMPLETED" {
			return fmt.Errorf("cluster(%s:%s) working sample data set error %s", projectID, job.GetId(), job.GetState())
		}
		return nil
	}
}

func configWithDailySchedule(clusterTerraformStr, clusterResourceName string, startHour, deleteExpirationDays int) string {
	var dataExpirationRuleBlock string
	if deleteExpirationDays > 0 {
		dataExpirationRuleBlock = fmt.Sprintf(`
		data_expiration_rule {
			expire_after_days = %d
		}
		`, deleteExpirationDays)
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
	`, clusterTerraformStr, startHour, dataExpirationRuleBlock, clusterResourceName)
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
