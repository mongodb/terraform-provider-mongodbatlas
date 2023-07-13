package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProjectDataSource_FrameworkMigration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckBasic(t) },
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.10.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: `data "mongodbatlas_project" "test" {
					name   = "framework-datasource-project"
					
				  }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mongodbatlas_project.test", "org_id", "63bec56c014da65b8f73c05e"),
				),
			},
			{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Config: `data "mongodbatlas_project" "test" {
					name   = "framework-datasource-project"
				  }`,
				PlanOnly: true,
			},
		},
	})
}
