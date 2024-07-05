package customdbrole_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const resourceName = "mongodbatlas_custom_db_role.test"

func TestAccConfigRSCustomDBRoles_Basic(t *testing.T) {
	var (
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName   = acc.RandomProjectName()
		roleName      = acc.RandomName()
		databaseName1 = acc.RandomClusterName()
		databaseName2 = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, roleName, "INSERT", databaseName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "INSERT"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),
				),
			},
			{
				Config: configBasic(orgID, projectName, roleName, "UPDATE", databaseName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "UPDATE"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCustomDBRolesImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"actions.0.resources.0.cluster"},
			},
		},
	})
}

func TestAccConfigRSCustomDBRoles_WithInheritedRoles(t *testing.T) {
	var (
		testRoleResourceName         = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceNameOne = "mongodbatlas_custom_db_role.inherited_role_one"
		InheritedRoleResourceNameTwo = "mongodbatlas_custom_db_role.inherited_role_two"
		orgID                        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName                  = acc.RandomProjectName()
	)

	inheritRole := []matlas.CustomDBRole{
		{
			RoleName: acc.RandomName(),
			Actions: []matlas.Action{{
				Action: "INSERT",
				Resources: []matlas.Resource{{
					DB: conversion.Pointer(acc.RandomClusterName()),
				}},
			}},
		},
		{
			RoleName: acc.RandomName(),
			Actions: []matlas.Action{{
				Action: "SERVER_STATUS",
				Resources: []matlas.Resource{{
					Cluster: conversion.Pointer(true),
				}},
			}},
		},
	}

	testRole := &matlas.CustomDBRole{
		RoleName: acc.RandomName(),
		Actions: []matlas.Action{{
			Action: "UPDATE",
			Resources: []matlas.Resource{{
				DB: conversion.Pointer(acc.RandomClusterName()),
			}},
		}},
	}

	inheritRoleUpdated := []matlas.CustomDBRole{
		{
			RoleName: inheritRole[0].RoleName,
			Actions: []matlas.Action{{
				Action: "FIND",
				Resources: []matlas.Resource{{
					DB: conversion.Pointer(acc.RandomClusterName()),
				}},
			}},
		},
		{
			RoleName: inheritRole[1].RoleName,
			Actions: []matlas.Action{{
				Action: "CONN_POOL_STATS",
				Resources: []matlas.Resource{{
					Cluster: conversion.Pointer(true),
				}},
			}},
		},
	}

	testRoleUpdated := &matlas.CustomDBRole{
		RoleName: testRole.RoleName,
		Actions: []matlas.Action{{
			Action: "REMOVE",
			Resources: []matlas.Resource{{
				DB: conversion.Pointer(acc.RandomClusterName()),
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
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRole[0].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRole[0].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRole[0].Actions[0].Resources))),

					// inherited Role [1]
					checkExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRole[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRole[1].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRole[1].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRole[1].Actions[0].Resources))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.Actions[0].Resources))),
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
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRoleUpdated[0].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRoleUpdated[0].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[0].Actions[0].Resources))),

					// inherited Role [1]
					checkExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRoleUpdated[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRoleUpdated[1].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRoleUpdated[1].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[1].Actions[0].Resources))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRoleUpdated.Actions[0].Resources))),
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

	inheritRole := &matlas.CustomDBRole{
		RoleName: acc.RandomName(),
		Actions: []matlas.Action{
			{
				Action: "REMOVE",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
		},
	}

	testRole := &matlas.CustomDBRole{
		RoleName: acc.RandomName(),
		Actions: []matlas.Action{
			{
				Action: "UPDATE",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
		},
		InheritedRoles: []matlas.InheritedRole{
			{
				Role: inheritRole.RoleName,
				Db:   "admin",
			},
		},
	}

	inheritRoleUpdated := &matlas.CustomDBRole{
		RoleName: inheritRole.RoleName,
		Actions: []matlas.Action{
			{
				Action: "UPDATE",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
		},
	}

	testRoleUpdated := &matlas.CustomDBRole{
		RoleName: testRole.RoleName,
		Actions: []matlas.Action{
			{
				Action: "REMOVE",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
		},
		InheritedRoles: []matlas.InheritedRole{
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
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.Actions[0].Resources))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.Actions[0].Resources))),
				),
			},
			{
				Config: configWithMultiple(orgID, projectName, inheritRoleUpdated, testRoleUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(

					// For Inherited Role
					checkExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.Actions[0].Resources))),

					// For Test Role
					checkExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRoleUpdated.Actions[0].Resources))),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSCustomDBRoles_MultipleResources(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)

	for i := 0; i < 5; i++ {
		roleName := fmt.Sprintf("test-acc-custom_role-%d", i)
		projectName := acc.RandomProjectName()
		t.Run(roleName, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { acc.PreCheckBasic(t) },
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				CheckDestroy:             checkDestroy,
				Steps: []resource.TestStep{
					{
						Config: configBasic(orgID, projectName, roleName, "INSERT", acc.RandomClusterName()),
						Check: resource.ComposeAggregateTestCheckFunc(
							checkExists(resourceName),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
							resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
							resource.TestCheckResourceAttr(resourceName, "actions.0.action", "INSERT"),
							resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),
						),
					},
				},
			})
		})
	}
}

func TestAccConfigRSCustomDBRoles_UpdatedInheritRoles(t *testing.T) {
	var (
		testRoleResourceName      = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceName = "mongodbatlas_custom_db_role.inherited_role"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acc.RandomProjectName()
	)

	inheritRole := &matlas.CustomDBRole{
		RoleName: acc.RandomName(),
		Actions: []matlas.Action{
			{
				Action: "REMOVE",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
		},
	}

	inheritRoleUpdated := &matlas.CustomDBRole{
		RoleName: inheritRole.RoleName,
		Actions: []matlas.Action{
			{
				Action: "UPDATE",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
					{
						DB: conversion.Pointer(acc.RandomClusterName()),
					},
				},
			},
		},
	}

	testRole := &matlas.CustomDBRole{
		RoleName: acc.RandomName(),
		InheritedRoles: []matlas.InheritedRole{
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
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.Actions[0].Resources))),

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
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.Actions[0].Resources))),

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

func configBasic(orgID, projectName, roleName, action, databaseName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_custom_db_role" "test" {
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
	`, orgID, projectName, roleName, action, databaseName)
}

func configWithInheritedRoles(orgID, projectName string, inheritedRole []matlas.CustomDBRole, testRole *matlas.CustomDBRole) string {
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
		inheritedRole[0].RoleName, inheritedRole[0].Actions[0].Action, *inheritedRole[0].Actions[0].Resources[0].DB,
		inheritedRole[1].RoleName, inheritedRole[1].Actions[0].Action, *inheritedRole[1].Actions[0].Resources[0].Cluster,
		testRole.RoleName, testRole.Actions[0].Action, *testRole.Actions[0].Resources[0].DB,
	)
}

func configWithMultiple(orgID, projectName string, inheritedRole, testRole *matlas.CustomDBRole) string {
	getCustomRoleFields := func(customRole *matlas.CustomDBRole) map[string]string {
		var (
			actions        string
			inheritedRoles string
		)

		for _, a := range customRole.Actions {
			var resources string

			// get the resources
			for _, r := range a.Resources {
				resources += fmt.Sprintf(`
					resources {
						collection_name = ""
						database_name   = "%s"
					}
			`, *r.DB)
			}

			// get the actions and set the resources
			actions += fmt.Sprintf(`
				actions {
					action = "%s"
					%s
				}
			`, a.Action, resources)
		}

		for _, in := range customRole.InheritedRoles {
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
