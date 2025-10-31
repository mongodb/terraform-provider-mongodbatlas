package streamworkspace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	region        = "VIRGINIA_USA"
	cloudProvider = "AWS"
)

func TestAccStreamsWorkspaceDS_basic(t *testing.T) {
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
				Config: streamsWorkspaceDataSourceConfig(projectID, workspaceName, region, cloudProvider),
				Check: resource.ComposeAggregateTestCheckFunc(
					streamsWorkspaceAttributeChecks(dataSourceName, workspaceName, region, cloudProvider),
					resource.TestCheckResourceAttr(dataSourceName, "stream_config.tier", "SP30"),
				),
			},
		},
	})
}

func streamsWorkspaceDataSourceConfig(projectID, workspaceName, region, cloudProvider string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_stream_workspace" "test" {
			project_id = mongodbatlas_stream_workspace.test.project_id
			workspace_name = mongodbatlas_stream_workspace.test.workspace_name
		}
	`, streamsWorkspaceConfig(projectID, workspaceName, region, cloudProvider))
}
