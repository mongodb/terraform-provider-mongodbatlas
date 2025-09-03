package customdbrole_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

const (
	resourceName         = "mongodbatlas_custom_db_role.test"
	dataSourceName       = "data.mongodbatlas_custom_db_role.test"
	dataSourcePluralName = "data.mongodbatlas_custom_db_roles.test"
	dataSourceSingular   = `
		data "mongodbatlas_custom_db_role" "test" {
			project_id = mongodbatlas_custom_db_role.test.project_id
			role_name  = mongodbatlas_custom_db_role.test.role_name
		}`
	dataSourcePlural = `
		data "mongodbatlas_custom_db_roles" "test" {
			project_id = mongodbatlas_custom_db_role.test.project_id

			depends_on = [mongodbatlas_custom_db_role.test]
		}
	`
)

func TestAccCustomDBRoles_Basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID       = acc.ProjectIDExecution(t)
		roleName        = acc.RandomName()
		databaseName1   = acc.RandomClusterName()
		databaseName2   = acc.RandomClusterName()
		pluralMapChecks = func(actionName, db string) map[string]knownvalue.Check {
			return map[string]knownvalue.Check{
				"actions.0.action":                    knownvalue.StringExact(actionName),
				"actions.0.resources.0.database_name": knownvalue.StringExact(db),
			}
		}
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:            configBasic(projectID, roleName, "INSERT", databaseName1) + dataSourceSingular + dataSourcePlural,
				Check:             checkAttrs(projectID, roleName, "INSERT", databaseName1),
				ConfigStateChecks: []statecheck.StateCheck{acc.PluralResultCheck(dataSourcePluralName, "role_name", knownvalue.StringExact(roleName), pluralMapChecks("INSERT", databaseName1))},
			},
			{
				Config:            configBasic(projectID, roleName, "UPDATE", databaseName2) + dataSourceSingular + dataSourcePlural,
				Check:             checkAttrs(projectID, roleName, "UPDATE", databaseName2),
				ConfigStateChecks: []statecheck.StateCheck{acc.PluralResultCheck(dataSourcePluralName, "role_name", knownvalue.StringExact(roleName), pluralMapChecks("UPDATE", databaseName2))},
			},
			{
				Config:            configBasic(projectID, roleName, "FIND", "") + dataSourceSingular + dataSourcePlural,
				Check:             checkAttrs(projectID, roleName, "FIND", ""),
				ConfigStateChecks: []statecheck.StateCheck{acc.PluralResultCheck(dataSourcePluralName, "role_name", knownvalue.StringExact(roleName), pluralMapChecks("FIND", ""))},
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCustomDBRolesImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"actions.0.resources.0.cluster"},
			},
		},
	}
}

func checkAttrs(projectID, roleName, action, databaseName string) resource.TestCheckFunc {
	return acc.CheckRSAndDS(
		resourceName,
		conversion.Pointer(dataSourceName),
		nil,
		nil,
		map[string]string{
			"project_id":                          projectID,
			"role_name":                           roleName,
			"actions.#":                           "1",
			"actions.0.action":                    action,
			"actions.0.resources.#":               "1",
			"actions.0.resources.0.database_name": databaseName,
		},
		checkExists(resourceName),
	)
}

func TestAccConfigRSCustomDBRoles_WithInheritedRoles(t *testing.T) {
	var (
		testRoleResourceName         = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceNameOne = "mongodbatlas_custom_db_role.inherited_role_one"
		InheritedRoleResourceNameTwo = "mongodbatlas_custom_db_role.inherited_role_two"
		orgID                        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName                  = acc.RandomProjectName()
	)

	inheritRole := []admin.UserCustomDBRole{
		{
			RoleName: acc.RandomName(),
			Actions: &[]admin.DatabasePrivilegeAction{{
				Action: "INSERT",
				Resources: &[]admin.DatabasePermittedNamespaceResource{{
					Db: acc.RandomClusterName(),
				}},
			}},
		},
		{
			RoleName: acc.RandomName(),
			Actions: &[]admin.DatabasePrivilegeAction{{
				Action: "SERVER_STATUS",
				Resources: &[]admin.DatabasePermittedNamespaceResource{{
					Cluster: true,
				}},
			}},
		},
	}

	testRole := &admin.UserCustomDBRole{
		RoleName: acc.RandomName(),
		Actions: &[]admin.DatabasePrivilegeAction{{
			Action: "UPDATE",
			Resources: &[]admin.DatabasePermittedNamespaceResource{{
				Db: acc.RandomClusterName(),
			}},
		}},
	}

	inheritRoleUpdated := []admin.UserCustomDBRole{
		{
			RoleName: inheritRole[0].RoleName,
			Actions: &[]admin.DatabasePrivilegeAction{{
				Action: "FIND",
				Resources: &[]admin.DatabasePermittedNamespaceResource{{
					Db: acc.RandomClusterName(),
				}},
			}},
		},
		{
			RoleName: inheritRole[1].RoleName,
			Actions: &[]admin.DatabasePrivilegeAction{{
				Action: "CONN_POOL_STATS",
				Resources: &[]admin.DatabasePermittedNamespaceResource{{
					Cluster: true,
				}},
			}},
		},
	}

	testRoleUpdated := &admin.UserCustomDBRole{
		RoleName: testRole.RoleName,
		Actions: &[]admin.DatabasePrivilegeAction{{
			Action: "REMOVE",
			Resources: &[]admin.DatabasePermittedNamespaceResource{{
				Db: acc.RandomClusterName(),
			}},
		}},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithInheritedRoles(orgID, projectName, inheritRole, testRole),
				Check: resource.ComposeAggregateTestCheckFunc(

					// For Inherited Roles
					// inherited Role [0]
					checkExists(InheritedRoleResourceNameOne),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "role_name", inheritRole[0].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRole[0].GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRole[0].GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRole[0].GetActions()[0].GetResources()))),

					// inherited Role [1]
					checkExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRole[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRole[1].GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRole[1].GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRole[1].GetActions()[0].GetResources()))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.GetActions()))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.GetActions()[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.GetActions()[0].GetResources()))),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "2"),
				),
			},
			{
				Config: configWithInheritedRoles(orgID, projectName, inheritRoleUpdated, testRoleUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(

					// For Inherited Role
					// inherited Role [0]
					checkExists(InheritedRoleResourceNameOne),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "role_name", inheritRoleUpdated[0].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRoleUpdated[0].GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRoleUpdated[0].GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[0].GetActions()[0].GetResources()))),

					// inherited Role [1]
					checkExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRoleUpdated[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRoleUpdated[1].GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRoleUpdated[1].GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[1].GetActions()[0].GetResources()))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRoleUpdated.GetActions()))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRoleUpdated.GetActions()[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRoleUpdated.GetActions()[0].GetResources()))),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "2"),
				),
			},
		},
	})
}

func TestAccConfigRSCustomDBRoles_MultipleCustomRoles(t *testing.T) {
	var (
		testRoleResourceName      = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceName = "mongodbatlas_custom_db_role.inherited_role"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acc.RandomProjectName()
	)

	inheritRole := &admin.UserCustomDBRole{
		RoleName: acc.RandomName(),
		Actions: &[]admin.DatabasePrivilegeAction{
			{
				Action: "REMOVE",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
			{
				Action: "FIND",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
		},
	}

	testRole := &admin.UserCustomDBRole{
		RoleName: acc.RandomName(),
		Actions: &[]admin.DatabasePrivilegeAction{
			{
				Action: "UPDATE",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
		},
		InheritedRoles: &[]admin.DatabaseInheritedRole{
			{
				Role: inheritRole.RoleName,
				Db:   "admin",
			},
		},
	}

	inheritRoleUpdated := &admin.UserCustomDBRole{
		RoleName: inheritRole.RoleName,
		Actions: &[]admin.DatabasePrivilegeAction{
			{
				Action: "UPDATE",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
			{
				Action: "FIND",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
		},
	}

	testRoleUpdated := &admin.UserCustomDBRole{
		RoleName: testRole.RoleName,
		Actions: &[]admin.DatabasePrivilegeAction{
			{
				Action: "REMOVE",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
		},
		InheritedRoles: &[]admin.DatabaseInheritedRole{
			{
				Role: inheritRole.RoleName,
				Db:   "admin",
			},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithMultiple(orgID, projectName, inheritRole, testRole),
				Check: resource.ComposeAggregateTestCheckFunc(

					// For Inherited Role
					checkExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRole.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.GetActions()[0].GetResources()))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.GetActions()))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.GetActions()[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.GetActions()[0].GetResources()))),
				),
			},
			{
				Config: configWithMultiple(orgID, projectName, inheritRoleUpdated, testRoleUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(

					// For Inherited Role
					checkExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.GetActions()[0].GetResources()))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRoleUpdated.GetActions()))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRoleUpdated.GetActions()[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRoleUpdated.GetActions()[0].GetResources()))),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSCustomDBRoles_UpdatedInheritRoles(t *testing.T) {
	var (
		testRoleResourceName      = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceName = "mongodbatlas_custom_db_role.inherited_role"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acc.RandomProjectName()
	)

	inheritRole := &admin.UserCustomDBRole{
		RoleName: acc.RandomName(),
		Actions: &[]admin.DatabasePrivilegeAction{
			{
				Action: "REMOVE",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
			{
				Action: "FIND",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
		},
	}

	inheritRoleUpdated := &admin.UserCustomDBRole{
		RoleName: inheritRole.RoleName,
		Actions: &[]admin.DatabasePrivilegeAction{
			{
				Action: "UPDATE",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
			{
				Action: "FIND",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: &[]admin.DatabasePermittedNamespaceResource{
					{
						Db: acc.RandomClusterName(),
					},
					{
						Db: acc.RandomClusterName(),
					},
				},
			},
		},
	}

	testRole := &admin.UserCustomDBRole{
		RoleName: acc.RandomName(),
		InheritedRoles: &[]admin.DatabaseInheritedRole{
			{
				Role: inheritRole.RoleName,
				Db:   "admin",
			},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithMultiple(orgID, projectName, inheritRole, testRole),
				Check: resource.ComposeAggregateTestCheckFunc(

					// For Inherited Role
					checkExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRole.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.GetActions()[0].GetResources()))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
			{
				Config: configWithMultiple(orgID, projectName, inheritRoleUpdated, testRole),
				Check: resource.ComposeAggregateTestCheckFunc(

					// For Inherited Role
					checkExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.GetActions()))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.GetActions()[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.GetActions()[0].GetResources()))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().CustomDatabaseRolesApi.GetCustomDatabaseRole(context.Background(), ids["project_id"], ids["role_name"]).Execute()
		if err != nil {
			return fmt.Errorf("custom DB Role (%s) does not exist", ids["role_name"])
		}
		return nil
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_custom_db_role" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().CustomDatabaseRolesApi.GetCustomDatabaseRole(context.Background(), ids["project_id"], ids["role_name"]).Execute()
		if err == nil {
			return fmt.Errorf("custom DB Role (%s) still exists", ids["role_name"])
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasCustomDBRolesImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["role_name"]), nil
	}
}

func configBasic(projectID, roleName, action, databaseName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_custom_db_role" "test" {
			project_id = %[1]q
			role_name  = %[2]q

			actions {
				action = %[3]q
				resources {
					collection_name = ""
					database_name   = %[4]q
				}
			}
		}
	`, projectID, roleName, action, databaseName)
}

func configWithInheritedRoles(orgID, projectName string, inheritedRole []admin.UserCustomDBRole, testRole *admin.UserCustomDBRole) string {
	return fmt.Sprintf(`

		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_custom_db_role" "inherited_role_one" {
		 	project_id = mongodbatlas_project.test.id
		 	role_name  = %[3]q

			actions {
				action = %[4]q
				resources {
					collection_name = ""
					database_name   = %[5]q
				}
			}
		}

		resource "mongodbatlas_custom_db_role" "inherited_role_two" {
			project_id = mongodbatlas_custom_db_role.inherited_role_one.project_id
		 	role_name  = %[6]q

			actions {
				action = %[7]q
				resources {
					cluster = %[8]t
				}
			}
		}

		resource "mongodbatlas_custom_db_role" "test_role" {
			project_id = mongodbatlas_custom_db_role.inherited_role_one.project_id
			role_name  = %[9]q

			actions {
				action = %[10]q
				resources {
					collection_name = ""
					database_name   = %[11]q
				}
			}

			inherited_roles {
				role_name     = mongodbatlas_custom_db_role.inherited_role_one.role_name
				database_name = "admin"
			}

			inherited_roles {
				role_name     = mongodbatlas_custom_db_role.inherited_role_two.role_name
				database_name = "admin"
			}
		}
	`, orgID, projectName,
		inheritedRole[0].RoleName, inheritedRole[0].GetActions()[0].Action, inheritedRole[0].GetActions()[0].GetResources()[0].Db,
		inheritedRole[1].RoleName, inheritedRole[1].GetActions()[0].Action, inheritedRole[1].GetActions()[0].GetResources()[0].Cluster,
		testRole.RoleName, testRole.GetActions()[0].Action, testRole.GetActions()[0].GetResources()[0].Db,
	)
}

func configWithMultiple(orgID, projectName string, inheritedRole, testRole *admin.UserCustomDBRole) string {
	getCustomRoleFields := func(customRole *admin.UserCustomDBRole) map[string]string {
		var (
			actions        string
			inheritedRoles string
		)

		for _, a := range customRole.GetActions() {
			var resources string

			// get the resources
			for _, r := range a.GetResources() {
				resources += fmt.Sprintf(`
					resources {
						collection_name = ""
						database_name   = "%s"
					}
			`, r.Db)
			}

			// get the actions and set the resources
			actions += fmt.Sprintf(`
				actions {
					action = "%s"
					%s
				}
			`, a.Action, resources)
		}

		for _, in := range customRole.GetInheritedRoles() {
			inheritedRoles += fmt.Sprintf(`
				inherited_roles {
					role_name     = "%s"
					database_name = "%s"
				}
			`, in.Role, in.Db)
		}

		return map[string]string{
			"actions":         actions,
			"inherited_roles": inheritedRoles,
		}
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_custom_db_role" "inherited_role" {
		 	project_id = mongodbatlas_project.test.id
		 	role_name  = %[3]q

			 %[4]s
		}

		resource "mongodbatlas_custom_db_role" "test_role" {
			project_id = mongodbatlas_custom_db_role.inherited_role.project_id
			role_name  = %[5]q

			%[6]s

			%[7]s
		}
	`, orgID, projectName,
		inheritedRole.RoleName, getCustomRoleFields(inheritedRole)["actions"],
		testRole.RoleName, getCustomRoleFields(testRole)["actions"], getCustomRoleFields(testRole)["inherited_roles"],
	)
}
