package mongodbatlas_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccBackupRSOnlineArchive(t *testing.T) {
	var (
		cluster                      matlas.Cluster
		resourceName                 = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName    = "mongodbatlas_online_archive.users_archive"
		onlineArchiveDataSourceName  = "data.mongodbatlas_online_archive.read_archive"
		onlineArchivesDataSourceName = "data.mongodbatlas_online_archives.all"
		orgID                        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName                  = acctest.RandomWithPrefix("test-acc")
		name                         = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				// We need this step to pupulate the cluster with Sample Data
				// The online archive won't work if the cluster does not have data
				Config: testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					populateWithSampleData(resourceName, &cluster),
				),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithDailySchedule(orgID, projectName, name, 1, 7),
				Check: resource.ComposeTestCheckFunc(
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
				Config: testAccBackupRSOnlineArchiveConfigWithDailySchedule(orgID, projectName, name, 2, 8),
				Check: resource.ComposeTestCheckFunc(
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
				Config: testAccBackupRSOnlineArchiveConfigWithWeeklySchedule(orgID, projectName, name, 2),
				Check: resource.ComposeTestCheckFunc(
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
				Config: testAccBackupRSOnlineArchiveConfigWithMonthlySchedule(orgID, projectName, name, 2),
				Check: resource.ComposeTestCheckFunc(
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
				Config: testAccBackupRSOnlineArchiveConfigWithoutSchedule(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
					resource.TestCheckNoResourceAttr(onlineArchiveResourceName, "schedule.#"),
				),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithoutSchedule(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "partition_fields.0.field_name", "last_review"),
				),
			},
		},
	})
}

func TestAccBackupRSOnlineArchiveBasic(t *testing.T) {
	var (
		cluster                   matlas.Cluster
		resourceName              = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acctest.RandomWithPrefix("test-acc")
		name                      = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				// We need this step to pupulate the cluster with Sample Data
				// The online archive won't work if the cluster does not have data
				Config: testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					populateWithSampleData(resourceName, &cluster),
				),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithoutSchedule(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
				),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithDailySchedule(orgID, projectName, name, 1, 1),
				Check: resource.ComposeTestCheckFunc(
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
		cluster                     matlas.Cluster
		resourceName                = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName   = "mongodbatlas_online_archive.users_archive"
		onlineArchiveDataSourceName = "data.mongodbatlas_online_archive.read_archive"
		orgID                       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName                 = acctest.RandomWithPrefix("test-acc")
		name                        = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		cloudProvider               = "AWS"
		processRegion               = "US_EAST_1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				// We need this step to pupulate the cluster with Sample Data
				// The online archive won't work if the cluster does not have data
				Config: testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					populateWithSampleData(resourceName, &cluster),
				),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithProcessRegion(orgID, projectName, name, cloudProvider, processRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.region", processRegion),
					resource.TestCheckResourceAttr(onlineArchiveDataSourceName, "data_process_region.0.cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(onlineArchiveDataSourceName, "data_process_region.0.region", processRegion),
				),
			},
			{
				Config:      testAccBackupRSOnlineArchiveConfigWithProcessRegion(orgID, projectName, name, cloudProvider, "AP_SOUTH_1"),
				ExpectError: regexp.MustCompile("ONLINE_ARCHIVE_CANNOT_MODIFY_FIELD"),
			},
			{
				Config: testAccBackupRSOnlineArchiveConfigWithoutSchedule(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.cloud_provider", cloudProvider),
					resource.TestCheckResourceAttr(onlineArchiveResourceName, "data_process_region.0.region", processRegion),
				),
			},
		},
	})
}

func TestAccBackupRSOnlineArchiveInvalidProcessRegion(t *testing.T) {
	var (
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName   = acctest.RandomWithPrefix("test-acc")
		name          = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		cloudProvider = "AWS"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBackupRSOnlineArchiveConfigWithProcessRegion(orgID, projectName, name, cloudProvider, "UNKNOWN"),
				ExpectError: regexp.MustCompile("INVALID_ATTRIBUTE"),
			},
		},
	})
}

func populateWithSampleData(resourceName string, cluster *matlas.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acc.TestMongoDBClient.(*config.MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		log.Printf("[DEBUG] projectID: %s, name %s", ids["project_id"], ids["cluster_name"])

		clusterResp, _, err := conn.Clusters.Get(context.Background(), ids["project_id"], ids["cluster_name"])

		if err != nil {
			return fmt.Errorf("cluster(%s:%s) does not exist %s", rs.Primary.Attributes["project_id"], rs.Primary.ID, err)
		}

		*cluster = *clusterResp

		job, _, err := conn.Clusters.LoadSampleDataset(context.Background(), ids["project_id"], ids["cluster_name"])

		if err != nil {
			return fmt.Errorf("cluster(%s:%s) loading sample data set error %s", rs.Primary.Attributes["project_id"], rs.Primary.ID, err)
		}

		ticker := time.NewTicker(30 * time.Second)

	JOB:
		for {
			select {
			case <-time.After(20 * time.Second):
				log.Println("timeout elapsed ....")
			case <-ticker.C:
				job, _, err = conn.Clusters.GetSampleDatasetStatus(context.Background(), ids["project_id"], job.ID)
				fmt.Println("querying for job ")
				if job.State != "WORKING" {
					break JOB
				}
			}
		}

		if err != nil {
			return fmt.Errorf("cluster(%s:%s) loading sample data set error %s", rs.Primary.Attributes["project_id"], rs.Primary.ID, err)
		}

		if job.State != "COMPLETED" {
			return fmt.Errorf("cluster(%s:%s) working sample data set error %s", rs.Primary.Attributes["project_id"], job.ID, job.State)
		}
		return nil
	}
}

func testAccBackupRSOnlineArchiveConfigWithDailySchedule(orgID, projectName, clusterName string, startHour, deleteExpirationDays int) string {
	return fmt.Sprintf(`
	%[1]s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = mongodbatlas_cluster.online_archive_test.project_id
		cluster_name = mongodbatlas_cluster.online_archive_test.name
		coll_name = "listingsAndReviews"
		collection_type = "STANDARD"
		db_name = "sample_airbnb"
	
		criteria {
			type = "DATE"
			date_field = "last_review"
			date_format = "ISODATE"
			expire_after_days = 2
		}

		data_expiration_rule {
			expire_after_days = %[3]d
		}

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
	`, testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, clusterName), startHour, deleteExpirationDays)
}

func testAccBackupRSOnlineArchiveConfigWithoutSchedule(orgID, projectName, clusterName string) string {
	return fmt.Sprintf(`
	%s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = mongodbatlas_cluster.online_archive_test.project_id
		cluster_name = mongodbatlas_cluster.online_archive_test.name
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
	`, testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, clusterName))
}

func testAccBackupRSOnlineArchiveConfigWithProcessRegion(orgID, projectName, clusterName, cloudProvider, region string) string {
	return fmt.Sprintf(`
	%s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = mongodbatlas_cluster.online_archive_test.project_id
		cluster_name = mongodbatlas_cluster.online_archive_test.name
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
	`, testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, clusterName), cloudProvider, region)
}

func testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, clusterName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "cluster_project" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cluster" "online_archive_test" {
		project_id   = mongodbatlas_project.cluster_project.id
		name         = %[3]q
		disk_size_gb = 10

		cluster_type = "REPLICASET"
		replication_specs {
		  num_shards = 1
		  regions_config {
			 region_name     = "US_EAST_1"
			 electable_nodes = 3
			 priority        = 7
			 read_only_nodes = 0
		   }
		}

		cloud_backup                 = false
		auto_scaling_disk_gb_enabled = true

		// Provider Settings "block"
		provider_name               = "AWS"
		provider_instance_size_name = "M10"

		labels {
			key   = "ArchiveTest"
			value = "true"
		}
		labels {
			key   = "Owner"
			value = "acctest"
		}
	}

	
	`, orgID, projectName, clusterName)
}

func testAccBackupRSOnlineArchiveConfigWithWeeklySchedule(orgID, projectName, clusterName string, startHour int) string {
	return fmt.Sprintf(`
	%s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = mongodbatlas_cluster.online_archive_test.project_id
		cluster_name = mongodbatlas_cluster.online_archive_test.name
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
			start_hour = %d
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
	`, testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, clusterName), startHour)
}

func testAccBackupRSOnlineArchiveConfigWithMonthlySchedule(orgID, projectName, clusterName string, startHour int) string {
	return fmt.Sprintf(`
	%s
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = mongodbatlas_cluster.online_archive_test.project_id
		cluster_name = mongodbatlas_cluster.online_archive_test.name
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
			start_hour = %d
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
	`, testAccBackupRSOnlineArchiveConfigFirstStep(orgID, projectName, clusterName), startHour)
}
