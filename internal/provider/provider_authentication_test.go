package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccSTSAssumeRole_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_project.test"
		projectID      = acc.ProjectIDGlobal(t)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckSTSAssumeRole(t); acc.PreCheckRegularCredsAreEmpty(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", projectID),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "teams.#"),
				),
			},
		},
	})
}

func configBasic(projectID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_project" "test" {
			project_id = %q
		}
`, projectID)
}
