package streamworkspace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamDSStreamWorkspace_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_stream_workspace.test"
		projectID      = acc.ProjectIDExecution(t)
		workspaceName  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: streamWorkspaceDataSourceConfig(projectID, workspaceName, region, cloudProvider),
				Check: resource.ComposeAggregateTestCheckFunc(
					streamWorkspaceAttributeChecks(dataSourceName, workspaceName, region, cloudProvider),
					resource.TestCheckResourceAttr(dataSourceName, "stream_config.tier", "SP30"),
				),
			},
		},
	})
}

func streamWorkspaceDataSourceConfig(projectID, workspaceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_workspace" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
			workspace_name = mongodbatlas_stream_workspace.test.workspace_name
		}
	`, acc.StreamInstanceConfig(projectID, workspaceName, region, cloudProvider))
}
