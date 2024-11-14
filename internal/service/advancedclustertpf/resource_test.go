package advancedclustertpf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_advanced_cluster.test"
)

func ChangeMockData(data *advancedclustertpf.MockData) resource.TestCheckFunc {
	changer := func(*terraform.State) error {
		data.NextResponse(true)
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

func CheckUpdatePayload(t *testing.T, requestName string) resource.TestCheckFunc {
	t.Helper()
	return func(state *terraform.State) error {
		g := goldie.New(t, goldie.WithNameSuffix(".json"))
		lastPayload, err := advancedclustertpf.ReadLastUpdatePayload()
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
		mockData    = &advancedclustertpf.MockData{
			ClusterResponse: "replicaset",
		}
	)
	err := advancedclustertpf.SetMockData(mockData)
	require.NoError(t, err)
	resource.Test(t, resource.TestCase{ // Sequential as it is using global variables
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "CREATING"),
					CheckRequestPayload(t, "replicaset_create"),
					ChangeMockData(mockData), // For the next test step
				),
			},
			{
				Config: configBasic(projectID, clusterName, "accept_data_risks_and_force_replica_set_reconfig = \"2006-01-02T15:04:05Z\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "accept_data_risks_and_force_replica_set_reconfig", "2006-01-02T15:04:05Z"),
					CheckUpdatePayload(t, "replicaset_update1"),
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

func TestAccAdvancedCluster_configSharded(t *testing.T) {
	var (
		projectID   = "111111111111111111111111"
		clusterName = "sharded-multi-replication"
		mockData    = &advancedclustertpf.MockData{
			ClusterResponse: "sharded",
		}
	)
	err := advancedclustertpf.SetMockData(mockData)
	require.NoError(t, err)
	resource.Test(t, resource.TestCase{ // Sequential as it is using global variables
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configSharded(projectID, clusterName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "CREATING"),
					CheckRequestPayload(t, "sharded_create"),
					ChangeMockData(mockData), // For the next test step
				),
			},
			{
				Config: configSharded(projectID, clusterName, true),
				Check: resource.ComposeTestCheckFunc(
					CheckUpdatePayload(t, "sharded_update1"),
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
func configSharded(projectID, clusterName string, withUpdate bool) string {
	var autoScaling, analyticsSpecs string
	if withUpdate {
		autoScaling = `
			auto_scaling = {
				disk_gb_enabled = true
			}`
		analyticsSpecs = `
			analytics_specs = {
				instance_size   = "M30"
				node_count      = 1
				ebs_volume_type = "PROVISIONED"
				disk_iops       = 2000
			}`
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "SHARDED"

			replication_specs = [
				{ # shard 1
				region_configs = [{
					electable_specs = {
						instance_size   = "M30"
						disk_iops       = 2000
						node_count      = 3
						ebs_volume_type = "PROVISIONED"
					}
					%[3]s
					%[4]s
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
					}]
					},
					{ # shard 2
				region_configs = [{
					electable_specs = {
						instance_size   = "M30"
						ebs_volume_type = "PROVISIONED"
						disk_iops       = 1000
						node_count      = 3
					}
					%[5]s
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}]
			}]
			}
		

	`, projectID, clusterName, autoScaling, analyticsSpecs, strings.ReplaceAll(analyticsSpecs, "2000", "1000"))
}
