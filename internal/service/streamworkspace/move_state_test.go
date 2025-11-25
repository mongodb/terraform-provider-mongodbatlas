package streamworkspace_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccStreamWorkspace_moveInstance(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0), // moved blocks require Terraform 1.8+
		},
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				Config: configMoveToStreamWorkspace(projectID, workspaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("mongodbatlas_stream_workspace.test", "project_id", projectID),
					resource.TestCheckResourceAttr("mongodbatlas_stream_workspace.test", "workspace_name", workspaceName),
					resource.TestCheckResourceAttrSet("mongodbatlas_stream_workspace.test", "data_process_region.cloud_provider"),
					resource.TestCheckResourceAttrSet("mongodbatlas_stream_workspace.test", "data_process_region.region"),
					resource.TestCheckResourceAttrSet("mongodbatlas_stream_workspace.test", "stream_config.tier"),
				),
			},
		},
	})
}

func configMoveToStreamWorkspace(projectID, workspaceName string) string {
	return fmt.Sprintf(`
moved {
  from = mongodbatlas_stream_instance.test
  to   = mongodbatlas_stream_workspace.test
}

resource "mongodbatlas_stream_workspace" "test" {
  project_id = "%s"
  workspace_name = "%s"
  data_process_region = {
    region = "VIRGINIA_USA"
    cloud_provider = "AWS"
  }
  stream_config = {
    tier = "SP30"
  }
}
`, projectID, workspaceName)
}
