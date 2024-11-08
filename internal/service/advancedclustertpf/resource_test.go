package advancedclustertpf_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_advanced_cluster.test"
)

func ChangeReponseNumber(responseNumber int) resource.TestCheckFunc {
	changer := func(*terraform.State) error {
		advancedclustertpf.SetCurrentClusterResponse(responseNumber)
		return nil
	}
	return changer
}

func TestAccAdvancedCluster_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(""),
			},
			{
				Config: configBasic("accept_data_risks_and_force_replica_set_reconfig = \"2006-01-02T15:04:05Z\""),
				Check:  resource.ComposeTestCheckFunc(ChangeReponseNumber(2)),
			},
			{
				Config:                               configBasic(""),
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func configBasic(extra string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = "111111111111111111111111"
			%[1]s
			name = "test"
			cluster_type = "REPLICASET"
			replication_specs = [{
				region_configs = [{
					priority        = 7
					provider_name = "AWS"
					region_name     = "US_EAST_1"
					electable_specs = {
						node_count = 3
						instance_size = "M10"
						disk_size_gb = 10
					}
				}]
			}]
		}
	`, extra)
}
