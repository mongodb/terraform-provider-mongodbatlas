package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConfigRSDatabaseUser_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: NewUnitTestProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.basic_ds", "project_id"),
					resource.TestCheckResourceAttr("mongodbatlas_database_user.basic_ds", "username", "test"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDatabaseUserConfig() string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "basic_ds" {
			username           = "test"
			project_id         = "test"
			auth_database_name = "test"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}
	`)
}
