package privateendpointregionalmode_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccNetworkDSPrivateEndpointRegionalMode_basic(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_private_endpoint_regional_mode.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDS(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "enabled", "true"),
				),
			},
		},
	})
}

func configDS(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id = mongodbatlas_project.project.id
			enabled = true
		}

		data "mongodbatlas_private_endpoint_regional_mode" "test" {
			project_id = mongodbatlas_project.project.id
			depends_on = [ mongodbatlas_private_endpoint_regional_mode.test ]
		}
	`, orgID, projectName)
}
