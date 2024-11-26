package advancedclustertpf_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_advanced_cluster.test"
)

func TestAccClusterAdvancedCluster_basicTenant(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)
	testCase := basicTenantTestCase(t, projectID, clusterName, clusterNameUpdated)
	resource.ParallelTest(t, *testCase)
}

func basicTenantTestCase(t *testing.T, projectID, clusterName, clusterNameUpdated string) *resource.TestCase {
	t.Helper()
	return &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		Steps: []resource.TestStep{
			{
				Config: configTenant(projectID, clusterName),
				Check:  checkTenant(projectID, clusterName),
			},
			{
				Config: configTenant(projectID, clusterNameUpdated),
				Check:  checkTenant(projectID, clusterNameUpdated),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	}
}

func configTenant(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M5"
					}
					provider_name         = "TENANT"
					backing_provider_name = "AWS"
					region_name           = "US_EAST_1"
					priority              = 7
				}]
			}]
		}
	`, projectID, name)
}

func checkTenant(projectID, name string) resource.TestCheckFunc {
	attrsSet := []string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"}
	attrsMap := map[string]string{
		"project_id":                           projectID,
		"name":                                 name,
		"termination_protection_enabled":       "false",
		"global_cluster_self_managed_sharding": "false",
		"labels.#":                             "0",
	}
	checks := acc.AddAttrSetChecks(resourceName, nil, attrsSet...)
	checks = acc.AddAttrChecks(resourceName, checks, attrsMap)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}
