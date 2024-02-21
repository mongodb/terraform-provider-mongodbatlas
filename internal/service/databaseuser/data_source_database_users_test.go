package databaseuser_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSDatabaseUsers_basic(t *testing.T) {
	var (
		resourceName = "data.mongodbatlas_database_users.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		roleName     = "atlasAdmin"
		projectName  = acc.RandomProjectName()
		username     = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfig(orgID, projectName, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user", "id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_database_user.db_user_1", "id"),
					resource.TestCheckResourceAttr("mongodbatlas_database_user.db_user_1", "labels.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(orgID, projectName, roleName, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "results.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.username"),
					resource.TestCheckResourceAttrSet(resourceName, "results.0.roles.#"),
					resource.TestCheckResourceAttrSet(resourceName, "results.1.labels.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfig(orgID, projectName, roleName, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}
		resource "mongodbatlas_database_user" "db_user" {
			username           = %[4]q
			password           = "test-acc-password"
			project_id         = mongodbatlas_project.test.id
			auth_database_name = "admin"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}

			labels {
				key   = "key 1"
				value = "value 1"
			}
		}

		resource "mongodbatlas_database_user" "db_user_1" {
			username           = "%[4]s-1"
			password           = "test-acc-password-1"
			project_id         = mongodbatlas_project.test.id
			auth_database_name = "admin"

			roles {
				role_name     = %[3]q
				database_name = "admin"
			}

			labels {
				key   = "key 1"
				value = "value 1"
			}

			labels {
				key   = "key 2"
				value = "value 2"
			}
		}
	`, orgID, projectName, roleName, username)
}

func testAccMongoDBAtlasDatabaseUsersDataSourceConfigWithDS(orgID, projectName, roleName, username string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_database_users" "test" {
			project_id = mongodbatlas_project.test.id
		}
	`, testAccMongoDBAtlasDatabaseUsersDataSourceConfig(orgID, projectName, roleName, username))
}
