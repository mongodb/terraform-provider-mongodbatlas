package databaseuser_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSDatabaseUser_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		roleName  = "atlasAdmin"
		username  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDS(projectID, username, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
					resource.TestCheckResourceAttr(dataSourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(dataSourceName, "x509_type", "NONE"),
					resource.TestCheckResourceAttr(dataSourceName, "roles.0.role_name", roleName),
					resource.TestCheckResourceAttr(dataSourceName, "roles.0.database_name", "admin"),
					resource.TestCheckResourceAttr(dataSourceName, "labels.#", "2"),
					resource.TestCheckNoResourceAttr(dataSourceName, "description"),
				),
			},
		},
	})
}

func configDS(projectID, username, roleName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			project_id         = %[1]q
			username           = %[2]q
			password           = "test-acc-password"
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

		data "mongodbatlas_database_user" "test" {
			username           = mongodbatlas_database_user.test.username
			project_id         = mongodbatlas_database_user.test.project_id
			auth_database_name = "admin"
		}
	`, projectID, username, roleName)
}
