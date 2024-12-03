package advancedclustertpf_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_move_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirst(projectID, clusterName),
			},
			{
				Config: configMoveSecond(projectID, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
				),
			},
		},
	})
}

func TestAccAdvancedCluster_move_invalid(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirst(projectID, clusterName),
			},
			{
				Config:      configMoveSecond(projectID, clusterName),
				ExpectError: regexp.MustCompile("Unable to Move Resource State"),
			},
			{
				Config: configMoveFirst(projectID, clusterName),
			},
		},
	})
}

func configMoveFirst(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "old" {
			project_id = %[1]q
			name = %[2]q
			disk_size_gb = 10
			cluster_type                = "REPLICASET"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_EAST_1"
					electable_nodes = 3
					priority        = 7
				}
			}
		}
	`, projectID, clusterName)
}

func configMoveSecond(projectID, clusterName string) string {
	return `
		moved {
			from = mongodbatlas_cluster.old
			to   = mongodbatlas_advanced_cluster.test
		}
	` + configBasic(projectID, clusterName, "")
}
