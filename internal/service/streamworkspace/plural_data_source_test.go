package streamworkspace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

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
