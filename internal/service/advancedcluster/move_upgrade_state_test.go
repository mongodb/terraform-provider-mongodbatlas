package advancedcluster_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccAdvancedCluster_moveBasic(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 3)
	)
	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirst(projectID, clusterName, 1),
			},
			{
				Config: configMoveSecond(projectID, clusterName, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccAdvancedCluster_moveMultisharding(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 9)
	)
	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirst(projectID, clusterName, 3),
			},
			{
				Config: configMoveSecond(projectID, clusterName, 3),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccAdvancedCluster_moveInvalid(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 0)
	)
	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMoveFirstInvalid(projectID, clusterName),
			},
			{
				Config:      configMoveSecondInvalid(projectID, clusterName),
				ExpectError: regexp.MustCompile("Unable to Move Resource State"),
			},
			{
				Config: configMoveFirstInvalid(projectID, clusterName),
			},
		},
	})
}

func configMoveFirst(projectID, clusterName string, numShards int) string {
	clusterTypeStr := "REPLICASET"
	if numShards > 1 {
		clusterTypeStr = "GEOSHARDED"
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "old" {
			project_id = %[1]q
			name = %[2]q
			disk_size_gb = 10
			cluster_type                = %[3]q
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			replication_specs {
				num_shards = %[4]d
				regions_config {
					region_name     = "US_EAST_1"
					electable_nodes = 3
					priority        = 7
				}
			}
		}
	`, projectID, clusterName, clusterTypeStr, numShards)
}

func configMoveBasic(projectID, clusterName string, numShards int) string {
	clusterTypeStr := "REPLICASET"
	if numShards > 1 {
		clusterTypeStr = "GEOSHARDED"
	}
	var replicationSpecsStr []string
	for range numShards {
		replicationSpecsStr = append(replicationSpecsStr, `
			{
				region_configs = [{
					priority        = 7
					provider_name = "AWS"
					region_name     = "US_EAST_1"
					auto_scaling = {
						compute_scale_down_enabled = false # necessary to have similar SDKv2 request
						compute_enabled = false # necessary to have similar SDKv2 request
						disk_gb_enabled = true
					}
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
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
			cluster_type = %[3]q
			replication_specs = [%[4]s]
		}
	`, projectID, clusterName, clusterTypeStr, strings.Join(replicationSpecsStr, ","))
}

func configMoveSecond(projectID, clusterName string, numShards int) string {
	return `
		moved {
			from = mongodbatlas_cluster.old
			to   = mongodbatlas_advanced_cluster.test
		}
	` + configMoveBasic(projectID, clusterName, numShards)
}

func configMoveFirstInvalid(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "old" {
			project_id         = %[1]q
			username           = %[2]q # use cluster name as username
			password           = "test-acc-password"
			auth_database_name = "admin"
			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}
	`, projectID, clusterName)
}

func configMoveSecondInvalid(projectID, clusterName string) string {
	return `
		moved {
			from = mongodbatlas_database_user.old
			to   = mongodbatlas_advanced_cluster.test
		}
	` + configMoveBasic(projectID, clusterName, 1)
}
