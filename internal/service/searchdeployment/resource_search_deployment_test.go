package searchdeployment_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceID   = "mongodbatlas_search_deployment.test"
	dataSourceID = "data.mongodbatlas_search_deployment.test"
)

func TestAccSearchDeployment_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			newSearchNodeTestStep(resourceID, orgID, projectName, clusterName, "S20_HIGHCPU_NVME", 3),
			newSearchNodeTestStep(resourceID, orgID, projectName, clusterName, "S30_HIGHCPU_NVME", 4),
			{
				Config:            configBasic(orgID, projectName, clusterName, "S30_HIGHCPU_NVME", 4),
				ResourceName:      resourceID,
				ImportStateIdFunc: importStateIDFunc(resourceID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSearchDeployment_multiRegion(t *testing.T) {
	var (
		clusterInfo = acc.GetClusterInfo(t, &acc.ClusterRequest{
			ClusterName: "multi-region-cluster",
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{
					Region: "US_EAST_1",
					ExtraRegionConfigs: []acc.ReplicationSpecRequest{
						{Region: "US_WEST_2", Priority: 6, InstanceSize: "M10", NodeCount: 2},
					},
				},
			},
		})
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: clusterInfo.TerraformStr + configSearchDeployment(clusterInfo.ProjectID, clusterInfo.TerraformNameRef, "S20_HIGHCPU_NVME", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "specs.#", "1"),
				),
			},
		},
	})
}

func newSearchNodeTestStep(resourceName, orgID, projectName, clusterName, instanceSize string, searchNodeCount int) resource.TestStep {
	resourceChecks := searchNodeChecks(resourceName, clusterName, instanceSize, searchNodeCount)
	dataSourceChecks := searchNodeChecks(dataSourceID, clusterName, instanceSize, searchNodeCount)
	return resource.TestStep{
		Config: configBasic(orgID, projectName, clusterName, instanceSize, searchNodeCount),
		Check:  resource.ComposeAggregateTestCheckFunc(append(resourceChecks, dataSourceChecks...)...),
	}
}

func searchNodeChecks(targetName, clusterName, instanceSize string, searchNodeCount int) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		checkExists(targetName),
		resource.TestCheckResourceAttrSet(targetName, "id"),
		resource.TestCheckResourceAttrSet(targetName, "project_id"),
		resource.TestCheckResourceAttr(targetName, "cluster_name", clusterName),
		resource.TestCheckResourceAttr(targetName, "specs.0.instance_size", instanceSize),
		resource.TestCheckResourceAttr(targetName, "specs.0.node_count", fmt.Sprintf("%d", searchNodeCount)),
		resource.TestCheckResourceAttrSet(targetName, "state_name"),
		resource.TestCheckResourceAttrSet(targetName, "encryption_at_rest_provider"),
	}
}

func configBasic(orgID, projectName, clusterName, instanceSize string, searchNodeCount int) string {
	clusterConfig := advancedClusterConfig(orgID, projectName, clusterName)
	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_search_deployment" "test" {
			project_id = mongodbatlas_project.test.id
			cluster_name = mongodbatlas_advanced_cluster.test.name
			specs = [
				{
					instance_size = %[2]q
					node_count = %[3]d
				}
			]
		}

		data "mongodbatlas_search_deployment" "test" {
			project_id = mongodbatlas_search_deployment.test.project_id
			cluster_name = mongodbatlas_search_deployment.test.cluster_name
		}
	`, clusterConfig, instanceSize, searchNodeCount)
}

func configSearchDeployment(projectID, clusterNameRef, instanceSize string, searchNodeCount int) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_search_deployment" "test" {
		project_id = %[1]q
		cluster_name = %[2]s
		specs = [
			{
				instance_size = %[3]q
				node_count = %[4]d
			}
		]
	}
	`, projectID, clusterNameRef, instanceSize, searchNodeCount)
}

func advancedClusterConfig(orgID, projectName, clusterName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "test" {
		org_id = %[1]q
		name   = %[2]q
	}
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = mongodbatlas_project.test.id
		name         = %[3]q
		cluster_type = "REPLICASET"
		retain_backups_enabled = "true"

		replication_specs {
			region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
			}
		}
	}
	`, orgID, projectName, clusterName)
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"]), nil
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		_, _, err := acc.ConnV2().AtlasSearchApi.GetAtlasSearchDeployment(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"]).Execute()
		if err != nil {
			return fmt.Errorf("search deployment (%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"])
		}
		return nil
	}
}

func checkDestroy(state *terraform.State) error {
	if projectDestroyedErr := acc.CheckDestroyProject(state); projectDestroyedErr != nil {
		return projectDestroyedErr
	}
	if clusterDestroyedErr := acc.CheckDestroyCluster(state); clusterDestroyedErr != nil {
		return clusterDestroyedErr
	}
	for _, rs := range state.RootModule().Resources {
		if rs.Type == "mongodbatlas_search_deployment" {
			_, _, err := acc.ConnV2().AtlasSearchApi.GetAtlasSearchDeployment(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"]).Execute()
			if err == nil {
				return fmt.Errorf("search deployment (%s:%s) still exists", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"])
			}
		}
	}
	return nil
}
