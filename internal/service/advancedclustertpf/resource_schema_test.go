package advancedclustertpf_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAdvancedCluster_PlanModifierErrors(t *testing.T) {
	var (
		projectID   = "111111111111111111111111"
		clusterName = "test"
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configBasic(projectID, clusterName, "accept_data_risks_and_force_replica_set_reconfig = \"2006-01-02T15:04:05Z\""),
				ExpectError: regexp.MustCompile("Update only attribute set on create: accept_data_risks_and_force_replica_set_reconfig"),
			},
		},
	})
}

func configBasic(projectID, clusterName, extra string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			timeouts = {
				create = "20s"
			}
			project_id = %[1]q
			name = %[2]q
			cluster_type = "REPLICASET"
			replication_specs = [{
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
			}]
			%[3]s
		}
	`, projectID, clusterName, extra)
}
