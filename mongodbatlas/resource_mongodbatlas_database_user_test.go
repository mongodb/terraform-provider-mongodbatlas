package mongodbatlas

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDatabaseUserUnitTest_basic(t *testing.T) {
	os.Setenv("MONGODB_ATLAS_BASE_URL", "http://localhost:8080/")
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBAtlasDatabaseUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.basic_ds", "project_id"),
					resource.TestCheckResourceAttr("mongodbatlas_database_user.basic_ds", "username", "test"),
				),
			},
		},
	})
}

func testMongoDBAtlasDatabaseUserConfig() string {
	return `
		resource "mongodbatlas_database_user" "basic_ds" {
			username           = "test"
			project_id         = "64f59c4f8db4346b7322b6df"
			auth_database_name = "admin"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}
	`
}
