package advancedclustertpf_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/sebdah/goldie/v2"

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

func CheckRequestPayload(t *testing.T, requestName string) resource.TestCheckFunc {
	t.Helper()
	return func(state *terraform.State) error {
		g := goldie.New(t, goldie.WithNameSuffix(".json"))
		lastPayload, err := advancedclustertpf.ReadLastCreatePayload()
		if err != nil {
			return err
		}
		g.Assert(t, requestName, []byte(lastPayload))
		return nil
	}
}

func TestAccAdvancedCluster_basic(t *testing.T) {
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
			{
				Config: configBasic(projectID, clusterName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "CREATING"),
					ChangeReponseNumber(2),
					CheckRequestPayload(t, "create_payload_check1"),
				),
			},
			{
				Config: configBasic(projectID, clusterName, "accept_data_risks_and_force_replica_set_reconfig = \"2006-01-02T15:04:05Z\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "accept_data_risks_and_force_replica_set_reconfig", "2006-01-02T15:04:05Z"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateIdFunc:                    acc.ImportStateIDFuncProjectIDClusterName(resourceName, "project_id", "name"),
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "name",
			},
		},
	})
}

func configBasic(projectID, clusterName, extra string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = %[1]q
			name = %[2]q
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
			%[3]s
		}
	`, projectID, clusterName, extra)
}
