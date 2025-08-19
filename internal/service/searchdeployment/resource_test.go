package searchdeployment_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/searchdeployment"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceID   = "mongodbatlas_search_deployment.test"
	dataSourceID = "data.mongodbatlas_search_deployment.test"
)

func importStep(tfConfig string) resource.TestStep {
	return resource.TestStep{
		Config:            tfConfig,
		ResourceName:      resourceID,
		ImportStateIdFunc: importStateIDFunc(resourceID),
		ImportState:       true,
		ImportStateVerify: true,
	}
}
func TestAccSearchDeployment_basic(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 6)
		updateStep             = newSearchNodeTestStep(resourceID, projectID, clusterName, "S30_HIGHCPU_NVME", 4)
		updateStepNoWait       = configBasic(projectID, clusterName, "S30_HIGHCPU_NVME", 4, true)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			newSearchNodeTestStep(resourceID, projectID, clusterName, "S20_HIGHCPU_NVME", 3),
			// Do a no-wait update then expect next step to wait for the update to complete
			// We cannot check the state_name as the response of a PATCH can be IDLE
			{
				Config: updateStepNoWait,
				Check:  updateStep.Check,
			},
			// Changes: skip_wait_on_update true -> null
			updateStep,
			importStep(updateStep.Config),
		},
	})
}

const deleteTimeout = 30 * time.Minute

func TestAccSearchDeployment_timeoutTest(t *testing.T) {
	var (
		timeoutsStrShort = `
			timeouts = {
				create = "90s"
			}
			delete_on_create_timeout = true
		`
		timeoutsStrLong        = strings.ReplaceAll(timeoutsStrShort, "90s", "6000s")
		timeoutsStrLongFalse   = strings.ReplaceAll(timeoutsStrLong, "true", "false")
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 6)
		configWithTimeout      = func(timeoutsStr string) string {
			normalConfig := configBasic(projectID, clusterName, "S20_HIGHCPU_NVME", 3, false)
			return acc.ConfigAddResourceStr(t, normalConfig, resourceID, timeoutsStr)
		}
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      configWithTimeout(timeoutsStrShort),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
			{
				PreConfig: func() {
					timeoutConfig := searchdeployment.RetryTimeConfig(deleteTimeout, 30*time.Second)
					err := searchdeployment.WaitSearchNodeDelete(t.Context(), projectID, clusterName, acc.ConnV2().AtlasSearchApi, timeoutConfig)
					require.NoError(t, err)
				},
				Config: configWithTimeout(timeoutsStrLong),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "timeouts.create", "6000s"),
					resource.TestCheckResourceAttr(resourceID, "delete_on_create_timeout", "true"),
				),
			},
			{
				Config: configWithTimeout(timeoutsStrLongFalse),
				Check:  resource.TestCheckResourceAttr(resourceID, "delete_on_create_timeout", "false"),
			},
			{
				Config: configWithTimeout(""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceID, "delete_on_create_timeout"),
					resource.TestCheckNoResourceAttr(resourceID, "timeouts.create"),
				),
			},
			importStep(configWithTimeout("")),
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

func newSearchNodeTestStep(resourceName, projectID, clusterName, instanceSize string, searchNodeCount int) resource.TestStep {
	resourceChecks := searchNodeChecks(resourceName, clusterName, instanceSize, searchNodeCount)
	dataSourceChecks := searchNodeChecks(dataSourceID, clusterName, instanceSize, searchNodeCount)
	return resource.TestStep{
		Config: configBasic(projectID, clusterName, instanceSize, searchNodeCount, false),
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
		// checking if encryption_at_rest_provider is set is not possible because it takes 10-15 minutes after the cluster is created to apply the encryption(and for the API to return the value)
	}
}

func configBasic(projectID, clusterName, instanceSize string, searchNodeCount int, skipWaitOnUpdate bool) string {
	clusterConfig := acc.ConfigBasicDedicated(projectID, clusterName, "")
	var skipWaitOnUpdateStr string
	if skipWaitOnUpdate {
		skipWaitOnUpdateStr = fmt.Sprintf("skip_wait_on_update = %t", skipWaitOnUpdate)
	}
	return fmt.Sprintf(`
		%[1]s

		resource "mongodbatlas_search_deployment" "test" {
			project_id = %[2]q
			cluster_name = mongodbatlas_advanced_cluster.test.name # ensure dependency on cluster
			specs = [
				{
					instance_size = %[3]q
					node_count = %[4]d
				}
			]
			%[5]s
		}

		data "mongodbatlas_search_deployment" "test" {
			project_id = mongodbatlas_search_deployment.test.project_id
			cluster_name = mongodbatlas_search_deployment.test.cluster_name
		}
	`, clusterConfig, projectID, instanceSize, searchNodeCount, skipWaitOnUpdateStr)
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

		deploymentResp, _, err := acc.ConnV2().AtlasSearchApi.GetAtlasSearchDeployment(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"]).Execute()
		if err != nil || searchdeployment.IsNotFoundDeploymentResponse(deploymentResp) {
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
