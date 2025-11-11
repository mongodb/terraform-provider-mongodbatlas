package streamworkspace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamsWorkspacesDS_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_workspaces.test"
		projectID      = acc.ProjectIDExecution(t)
		workspaceName  = acc.RandomName()
	)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
		resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
		resource.TestCheckResourceAttrSet(dataSourceName, "results.0.workspace_name"),
		resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_process_region.region"),
		resource.TestCheckResourceAttrSet(dataSourceName, "results.0.data_process_region.cloud_provider"),
		resource.TestCheckResourceAttrSet(dataSourceName, "results.0.hostnames.#"),
		resource.TestCheckResourceAttr(dataSourceName, "results.0.stream_config.max_tier_size", "SP30"),
		resource.TestCheckResourceAttr(dataSourceName, "results.0.stream_config.tier", "SP10"),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamsWorkspacesDataSourceConfig(projectID, workspaceName, region, cloudProvider),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

func TestAccStreamsWorkspacesDS_withPageConfig(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_workspaces.test"
		projectID      = acc.ProjectIDExecution(t)
		workspaceName  = acc.RandomName()
		pageNumber     = 1000 // high page number so no results are returned
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamsWorkspacesWithPageAttrDataSourceConfig(projectID, workspaceName, region, cloudProvider, pageNumber),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "0"),
				),
			},
		},
	})
}

func streamsWorkspacesDataSourceConfig(projectID, workspaceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_workspaces" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
		}
	`, streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider))
}

func streamsWorkspacesWithPageAttrDataSourceConfig(projectID, workspaceName, region, cloudProvider string, pageNum int) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_workspaces" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
			page_num = %d
			items_per_page = 1
		}
	`, streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider), pageNum)
}
