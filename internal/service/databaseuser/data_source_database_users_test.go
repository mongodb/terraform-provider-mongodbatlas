package databaseuser_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSDatabaseUsers_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
		roleName  = "atlasAdmin"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configDSPlural(projectID, username, roleName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.username"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.roles.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.labels.#"),
					resource.TestCheckNoResourceAttr(dataSourcePluralName, "results.0.description"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.1.username"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.1.roles.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.1.labels.#"),
					resource.TestCheckNoResourceAttr(dataSourcePluralName, "results.1.description"),
				),
			},
		},
	})
}

func configDSPlural(projectID, username, roleName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "db_user" {
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
		}

		resource "mongodbatlas_database_user" "db_user_1" {
			project_id         = %[1]q
			username           = "%[2]s-1"
			password           = "test-acc-password-1"
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

		data "mongodbatlas_database_users" "test" {
			project_id = %[1]q
			depends_on = [mongodbatlas_database_user.db_user, mongodbatlas_database_user.db_user_1]
		}
	`, projectID, username, roleName)
}
