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

const (
	clusterConfig = `
	resource "mongodbatlas_cluster" "online_archive_test" {
		project_id   = "%s"
		name         = "%s"
		disk_size_gb = 10
		num_shards   = 1

		replication_factor           = 3
		provider_backup_enabled      = %s
		auto_scaling_disk_gb_enabled = true

		// Provider Settings "block"
		provider_name               = "AWS"
		provider_encrypt_ebs_volume = false
		provider_instance_size_name = "M10"
		provider_region_name        = "US_EAST_2"

		labels {
			key   = "ArchiveTest"
			value = "true"
		}
		labels {
			key   = "Owner"
			value = "acctest"
		}
	}

	data "mongodbatlas_clusters" "online_archive_test" {
		project_id = mongodbatlas_cluster.online_archive_test.project_id
	}
`

	onlineArchiveConfig = `
	resource "mongodbatlas_online_archive" "users_archive" {
		project_id = mongodbatlas_cluster.online_archive_test.project_id
		cluster_name = mongodbatlas_cluster.online_archive_test.name
		coll_name = "listingsAndReviews"
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
`
)

func TestAccResourceMongoDBAtlasOnlineArchive(t *testing.T) {
	var (
		cluster                   matlas.Cluster
		resourceName              = "mongodbatlas_cluster.online_archive_test"
		onlineArchiveResourceName = "mongodbatlas_online_archive.users_archive"
		projectID                 = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name                      = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	initialConfig := fmt.Sprintf(clusterConfig, projectID, name, "false")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					populateWithSampleData(resourceName, &cluster),
				),
			},
			{
				Config: initialConfig + onlineArchiveConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "state"),
					resource.TestCheckResourceAttrSet(onlineArchiveResourceName, "archive_id"),
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
