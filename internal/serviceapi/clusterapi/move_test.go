package clusterapi_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccClusterAPI_moveBasic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyClusters,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirst(projectID, clusterName, 2),
			},
			{
				Config: configMoveSecond(projectID, clusterName, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccClusterAPI_moveFromUnsupportedSource(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyDatabaseUser,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirstUnsupported(projectID, clusterName),
			},
			{
				Config:      configMoveSecondUnsupported(projectID, clusterName),
				ExpectError: regexp.MustCompile("Unable to Move Resource State"),
			},
			{
				Config: configMoveFirstUnsupported(projectID, clusterName),
			},
		},
	})
}

func configMoveFirst(projectID, clusterName string, numShards int) string {
	clusterTypeStr := "REPLICASET"
	if numShards > 1 {
		clusterTypeStr = "SHARDED"
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster_old_api" "old" {
			group_id   = %[1]q
			name         = %[2]q
			cluster_type = %[3]q
			disk_size_gb = 10

			replication_specs = [{
				num_shards = %[4]d

				region_configs = [{
					provider_name = "AWS"
					region_name   = "US_EAST_1"
					priority = 7

					electable_specs = {
						node_count = 3
						instance_size = "M10" 
					}
				}]
			}]
		}
	`, projectID, clusterName, clusterTypeStr, numShards)
}

func configMoveSecond(projectID, clusterName string, numShards int) string {
	clusterTypeStr := "REPLICASET"
	if numShards > 1 {
		clusterTypeStr = "SHARDED"
	}
	var replicationSpecsStr []string
	for range numShards {
		replicationSpecsStr = append(replicationSpecsStr, `
			{
				region_configs = [{
					provider_name = "AWS"
					region_name   = "US_EAST_1"
					priority      = 7

					electable_specs = {
						node_count = 3
						instance_size = "M10"
						disk_size_gb = 10
					}
				}]
			}
		`)
	}

	return fmt.Sprintf(`
		moved {
			from = mongodbatlas_cluster_old_api.old
			to   = mongodbatlas_cluster_api.test
		}

		resource "mongodbatlas_cluster_api" "test" {
			group_id          = %[1]q
			name		      = %[2]q
			cluster_type 	  = %[3]q
			replication_specs = [%[4]s]
		}
	`, projectID, clusterName, clusterTypeStr, strings.Join(replicationSpecsStr, ","))
}

func configMoveFirstUnsupported(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user_api" "old" {
			group_id      = %[1]q
			username      = %[2]q
			password      = "test-acc-password"
			database_name = "admin"

			roles = [{
				role_name     = "atlasAdmin"
				database_name = "admin"
			}]
		}
	`, projectID, clusterName)
}

func configMoveSecondUnsupported(projectID, clusterName string) string {
	return fmt.Sprintf(`
		moved {
			from = mongodbatlas_database_user_api.old
			to   = mongodbatlas_cluster_api.test
		}

		resource "mongodbatlas_cluster_api" "test" {
			group_id          = %[1]q
			name		      = %[2]q
			cluster_type 	  = "REPLICASET"
			replication_specs = [{
				region_configs = [{
					provider_name = "AWS"
					region_name   = "US_EAST_1"
					priority      = 7

					electable_specs = {
						node_count = 3
						instance_size = "M10"
						disk_size_gb = 10
					}
				}]
			}]
		}
	`, projectID, clusterName)
}

func checkDestroyClusters(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster_old_api" && rs.Type != "mongodbatlas_cluster_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		clusterName := rs.Primary.Attributes["name"]
		if groupID == "" || clusterName == "" {
			return fmt.Errorf("groupID or clusterName is empty: %s, %s", groupID, clusterName)
		}
		resp, _, _ := acc.ConnV2().ClustersApi.GetCluster(context.Background(), groupID, clusterName).Execute()
		if resp.GetId() != "" {
			return fmt.Errorf("cluster (%s:%s) still exists", clusterName, rs.Primary.ID)
		}
	}
	return nil
}

func checkDestroyDatabaseUser(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_database_user_api" {
			continue
		}
		groupID := rs.Primary.Attributes["group_id"]
		databaseName := rs.Primary.Attributes["database_name"]
		username := rs.Primary.Attributes["username"]
		if groupID == "" || databaseName == "" || username == "" {
			return fmt.Errorf("groupID, databaseName or username is empty: %s, %s, %s", groupID, databaseName, username)
		}
		_, _, err := acc.ConnV2().DatabaseUsersApi.GetDatabaseUser(context.Background(), groupID, databaseName, username).Execute()
		if err == nil {
			return fmt.Errorf("database user (%s) still exists", username)
		}
	}
	return nil
}
