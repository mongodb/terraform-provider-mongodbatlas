package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccBackupRSOnlineArchive(t *testing.T) {
	var (
		cluster                   matlas.Cluster
		resourceName              = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acctest.RandomWithPrefix("test-acc")
		name                      = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasClusterDestroy,
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
				Config: testAccBackupRSOnlineArchiveConfig(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "collection_type"),
				),
			},
		},
	})
}

func populateWithSampleData(resourceName string, cluster *matlas.Cluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

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

func testAccBackupRSOnlineArchiveConfig(orgID, projectName, clusterName string) string {
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
			field_name = "maximum_nights"
			order = 0
		}
	
		partition_fields {
			field_name = "name"
			order = 1
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
