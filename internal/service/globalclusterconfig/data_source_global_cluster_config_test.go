package globalclusterconfig_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccClusterRSGlobalClusterDS_basic(t *testing.T) {
	acc.SkipTestForCI(t) // needs to be fixed: 404 (request "GROUP_NOT_FOUND") No group with ID
	var (
		dataSourceName = "data.mongodbatlas_global_cluster_config.config"
		name           = fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configGlobalCluster(orgID, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_name"),
				),
			},
		},
	})
}

func configGlobalCluster(orgID, projectName, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "project" {
		org_id = %[1]q
		name   = %[2]q
	}

	resource "mongodbatlas_cluster" "test" {
		project_id              = mongodbatlas_project.project.id
		name                    = %[3]q
		disk_size_gb            = 80
		cloud_backup            = false
		cluster_type            = "GEOSHARDED"

		// Provider Settings "block"
		provider_name               = "AWS"
		provider_instance_size_name = "M30"

		replication_specs {
			zone_name  = "Zone 1"
			num_shards = 2
			regions_config {
			region_name     = "US_EAST_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
			}
		}

		replication_specs {
			zone_name  = "Zone 2"
			num_shards = 2
			regions_config {
			region_name     = "US_EAST_2"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
			}
		}
	}

	resource "mongodbatlas_global_cluster_config" "config" {
		project_id = mongodbatlas_cluster.test.project_id
		cluster_name = mongodbatlas_cluster.test.name

		managed_namespaces {
			db 				 = "mydata"
			collection 		 = "publishers"
			custom_shard_key = "city"
		}

		custom_zone_mappings {
			location ="CA"
			zone =  "Zone 1"
		}
	}

	data "mongodbatlas_global_cluster_config" "config" {
		project_id = mongodbatlas_global_cluster_config.config.project_id
		cluster_name = mongodbatlas_global_cluster_config.config.cluster_name
	}
	`, orgID, projectName, name)
}
