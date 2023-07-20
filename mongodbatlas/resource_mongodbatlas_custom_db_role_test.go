package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSCustomDBRoles_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_custom_db_role.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		roleName     = fmt.Sprintf("test-acc-custom_role-%s", acctest.RandString(5))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(orgID, projectName, roleName, "INSERT", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_name"),
					resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "INSERT"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(orgID, projectName, roleName, "UPDATE", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_name"),
					resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "UPDATE"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSCustomDBRoles_WithInheritedRoles(t *testing.T) {
	testRoleResourceName := "mongodbatlas_custom_db_role.test_role"
	InheritedRoleResourceNameOne := "mongodbatlas_custom_db_role.inherited_role_one"
	InheritedRoleResourceNameTwo := "mongodbatlas_custom_db_role.inherited_role_two"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")

	inheritRole := []matlas.CustomDBRole{
		{
			RoleName: fmt.Sprintf("test-acc-INHERITED_ROLE-%s", acctest.RandString(5)),
			Actions: []matlas.Action{{
				Action: "INSERT",
				Resources: []matlas.Resource{{
					DB: pointy.String(fmt.Sprintf("b_test-acc-ddb_name-%s", acctest.RandString(5))),
				}},
			}},
		},
		{
			RoleName: fmt.Sprintf("test-acc-INHERITED_ROLE-%s", acctest.RandString(5)),
			Actions: []matlas.Action{{
				Action: "SERVER_STATUS",
				Resources: []matlas.Resource{{
					Cluster: pointy.Bool(true),
				}},
			}},
		},
	}

	testRole := &matlas.CustomDBRole{
		RoleName: fmt.Sprintf("test-acc-TEST_ROLE-%s", acctest.RandString(5)),
		Actions: []matlas.Action{{
			Action: "UPDATE",
			Resources: []matlas.Resource{{
				DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
			}},
		}},
	}

	inheritRoleUpdated := []matlas.CustomDBRole{
		{
			RoleName: inheritRole[0].RoleName,
			Actions: []matlas.Action{{
				Action: "FIND",
				Resources: []matlas.Resource{{
					DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
				}},
			}},
		},
		{
			RoleName: inheritRole[1].RoleName,
			Actions: []matlas.Action{{
				Action: "CONN_POOL_STATS",
				Resources: []matlas.Resource{{
					Cluster: pointy.Bool(true),
				}},
			}},
		},
	}

	testRoleUpdated := &matlas.CustomDBRole{
		RoleName: testRole.RoleName,
		Actions: []matlas.Action{{
			Action: "REMOVE",
			Resources: []matlas.Resource{{
				DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
			}},
		}},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigWithInheritedRoles(orgID, projectName, inheritRole, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Roles
					// inherited Role [0]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameOne),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "role_name", inheritRole[0].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRole[0].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRole[0].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRole[0].Actions[0].Resources))),

					// inherited Role [1]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRole[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRole[1].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRole[1].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRole[1].Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.Actions[0].Resources))),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigWithInheritedRoles(orgID, projectName, inheritRoleUpdated, testRoleUpdated),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					// inherited Role [0]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameOne),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "role_name", inheritRoleUpdated[0].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRoleUpdated[0].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRoleUpdated[0].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[0].Actions[0].Resources))),

					// inherited Role [1]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRoleUpdated[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRoleUpdated[1].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRoleUpdated[1].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[1].Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),
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
		projectName               = acctest.RandomWithPrefix("test-acc")
	)

	inheritRole := &matlas.CustomDBRole{
		RoleName: fmt.Sprintf("test-acc-INHERITED_ROLE-%s", acctest.RandString(5)),
		Actions: []matlas.Action{
			{
				Action: "REMOVE",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
		},
	}

	testRole := &matlas.CustomDBRole{
		RoleName: fmt.Sprintf("test-acc-TEST_ROLE-%s", acctest.RandString(5)),
		Actions: []matlas.Action{
			{
				Action: "UPDATE",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
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
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
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
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
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
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(orgID, projectName, inheritRole, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRole.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.Actions[0].Resources))),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(orgID, projectName, inheritRoleUpdated, testRoleUpdated),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),
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
	t.Skip()
	var (
		resourceName = "mongodbatlas_custom_db_role.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		roleName     string
	)

	for i := 0; i < 100; i++ {
		roleName = fmt.Sprintf("test-acc-custom_role-%d", i)

		t.Run(roleName, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheckBasic(t) },
				ProviderFactories: testAccProviderFactories,
				CheckDestroy:      testAccCheckMongoDBAtlasCustomDBRolesDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(orgID, projectName, roleName, "INSERT", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttrSet(resourceName, "role_name"),
							resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),
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

func TestAccConfigRSCustomDBRoles_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_custom_db_role.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		roleName     = fmt.Sprintf("test-acc-custom_role-%s", acctest.RandString(5))
		databaseName = fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(orgID, projectName, roleName, "INSERT", databaseName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasCustomDBRolesImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConfigRSCustomDBRoles_UpdatedInheritRoles(t *testing.T) {
	var (
		testRoleResourceName      = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceName = "mongodbatlas_custom_db_role.inherited_role"
		orgID                     = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName               = acctest.RandomWithPrefix("test-acc")
	)

	inheritRole := &matlas.CustomDBRole{
		RoleName: fmt.Sprintf("test-acc-INHERITED_ROLE-%s", acctest.RandString(5)),
		Actions: []matlas.Action{
			{
				Action: "REMOVE",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
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
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
					{
						DB: pointy.String(fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
					},
				},
			},
		},
	}

	testRole := &matlas.CustomDBRole{
		RoleName: fmt.Sprintf("test-acc-TEST_ROLE-%s", acctest.RandString(5)),
		InheritedRoles: []matlas.InheritedRole{
			{
				Role: inheritRole.RoleName,
				Db:   "admin",
			},
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(orgID, projectName, inheritRole, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRole.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(orgID, projectName, inheritRoleUpdated, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.CustomDBRoles.Get(context.Background(), ids["project_id"], ids["role_name"])
		if err != nil {
			return fmt.Errorf("custom DB Role (%s) does not exist", ids["role_name"])
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasCustomDBRolesDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_custom_db_role" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.CustomDBRoles.Get(context.Background(), ids["project_id"], ids["role_name"])
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

func testAccMongoDBAtlasCustomDBRolesConfigBasic(orgID, projectName, roleName, action, databaseName string) string {
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

func testAccMongoDBAtlasCustomDBRolesConfigWithInheritedRoles(orgID, projectName string, inheritedRole []matlas.CustomDBRole, testRole *matlas.CustomDBRole) string {
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
			project_id = "${mongodbatlas_custom_db_role.inherited_role_one.project_id}"
		 	role_name  = %[6]q

			actions {
				action = %[7]q
				resources {
					cluster = %[8]t
				}
			}
		}

		resource "mongodbatlas_custom_db_role" "test_role" {
			project_id = "${mongodbatlas_custom_db_role.inherited_role_one.project_id}"
			role_name  = %[9]q

			actions {
				action = %[10]q
				resources {
					collection_name = ""
					database_name   = %[11]q
				}
			}

			inherited_roles {
				role_name     = "${mongodbatlas_custom_db_role.inherited_role_one.role_name}"
				database_name = "admin"
			}

			inherited_roles {
				role_name     = "${mongodbatlas_custom_db_role.inherited_role_two.role_name}"
				database_name = "admin"
			}
		}
	`, orgID, projectName,
		inheritedRole[0].RoleName, inheritedRole[0].Actions[0].Action, *inheritedRole[0].Actions[0].Resources[0].DB,
		inheritedRole[1].RoleName, inheritedRole[1].Actions[0].Action, *inheritedRole[1].Actions[0].Resources[0].Cluster,
		testRole.RoleName, testRole.Actions[0].Action, *testRole.Actions[0].Resources[0].DB,
	)
}

func testAccMongoDBAtlasCustomDBRolesConfigMultiple(orgID, projectName string, inheritedRole, testRole *matlas.CustomDBRole) string {
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
			project_id = "${mongodbatlas_custom_db_role.inherited_role.project_id}"
			role_name  = %[5]q

			%[6]s

			%[7]s
		}
	`, orgID, projectName,
		inheritedRole.RoleName, getCustomRoleFields(inheritedRole)["actions"],
		testRole.RoleName, getCustomRoleFields(testRole)["actions"], getCustomRoleFields(testRole)["inherited_roles"],
	)
}
