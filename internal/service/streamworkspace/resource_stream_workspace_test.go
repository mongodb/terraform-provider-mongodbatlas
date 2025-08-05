package streamworkspace_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamRSStreamWorkspace_basic(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: acc.StreamInstanceConfig(projectID, workspaceName, region, cloudProvider), // as of now there are no values that can be updated because only one region is supported
				Check: resource.ComposeAggregateTestCheckFunc(
					streamWorkspaceAttributeChecks(resourceName, workspaceName, region, cloudProvider),
					resource.TestCheckResourceAttr(resourceName, "stream_config.tier", "SP30"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamWorkspaceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccStreamRSStreamWorkspace_withStreamConfig(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: acc.StreamInstanceWithStreamConfigConfig(projectID, workspaceName, region, cloudProvider, "SP10"), // as of now there are no values that can be updated because only one region is supported
				Check: resource.ComposeAggregateTestCheckFunc(
					streamWorkspaceAttributeChecks(resourceName, workspaceName, region, cloudProvider),
					resource.TestCheckResourceAttr(resourceName, "stream_config.tier", "SP10"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: checkStreamWorkspaceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func streamWorkspaceAttributeChecks(resourceName, workspaceName, region, cloudProvider string) resource.TestCheckFunc {
	resourceChecks := []resource.TestCheckFunc{
		checkSearchWorkspaceExists(),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "project_id"),
		resource.TestCheckResourceAttr(resourceName, "workspace_name", workspaceName),
		resource.TestCheckResourceAttr(resourceName, "data_process_region.region", region),
		resource.TestCheckResourceAttr(resourceName, "data_process_region.cloud_provider", cloudProvider),
		resource.TestCheckResourceAttr(resourceName, "hostnames.#", "1"),
	}
	return resource.ComposeAggregateTestCheckFunc(resourceChecks...)
}

func checkStreamWorkspaceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["workspace_name"]), nil
	}
}

func checkSearchWorkspaceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "mongodbatlas_stream_workspace" {
				_, _, err := acc.ConnV2().StreamsApi.GetStreamInstance(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["workspace_name"]).Execute()
				if err != nil {
					return fmt.Errorf("stream workspace (%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["workspace_name"])
				}
			}
		}
		return nil
	}
}
