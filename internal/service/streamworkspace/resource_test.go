package streamworkspace_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamWorkspaceRS_basic(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance, // Reuse the same destroy check
		Steps: []resource.TestStep{
			{
				Config: streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider),
				Check: resource.ComposeAggregateTestCheckFunc(
					streamsWorkspaceAttributeChecks(resourceName, workspaceName, region, cloudProvider),
					resource.TestCheckResourceAttr(resourceName, "stream_config.max_tier_size", "SP30"),
					resource.TestCheckResourceAttr(resourceName, "stream_config.tier", "SP10"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamsWorkspaceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStreamWorkspaceRS_withStreamConfig(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance, // Reuse the same destroy check
		Steps: []resource.TestStep{
			{
				Config: streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider),
				Check: resource.ComposeAggregateTestCheckFunc(
					streamsWorkspaceAttributeChecks(resourceName, workspaceName, region, cloudProvider),
					resource.TestCheckResourceAttr(resourceName, "stream_config.max_tier_size", "SP30"),
					resource.TestCheckResourceAttr(resourceName, "stream_config.tier", "SP10"),
				),
			},
		},
	})
}

func streamsWorkspaceAttributeChecks(resourceName, workspaceName, region, cloudProvider string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		checkStreamsWorkspaceExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
		resource.TestCheckResourceAttr(resourceName, "data_process_region.region", region),
		resource.TestCheckResourceAttr(resourceName, "data_process_region.cloud_provider", cloudProvider),
		resource.TestCheckResourceAttrSet(resourceName, "hostnames.#"),
	)
}

func checkStreamsWorkspaceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		projectID := rs.Primary.Attributes["project_id"]
		workspaceName := rs.Primary.Attributes["workspace_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamWorkspace(context.Background(), projectID, workspaceName).Execute()
		if err != nil {
			return fmt.Errorf("stream workspace (%s:%s) does not exist: %s", projectID, workspaceName, err)
		}
		return nil
	}
}

func checkStreamsWorkspaceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["workspace_name"]), nil
	}
}

func streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_workspace" "test" {
			project_id = %[1]q
			workspace_name = %[2]q
			data_process_region = {
				region = %[3]q
				cloud_provider = %[4]q
			}
			stream_config = {
				max_tier_size = "SP30"
				tier = "SP10"
			}
		}
	`, projectID, workspaceName, region, cloudProvider)
}
