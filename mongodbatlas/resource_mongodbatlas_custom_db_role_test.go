package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCustomDBRoles_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_custom_db_role.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		roleName     = fmt.Sprintf("test-acc-custom_role-%s", acctest.RandString(5))
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(projectID, roleName, "INSERT", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_name"),
					resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "INSERT"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(projectID, roleName, "UPDATE", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "role_name"),
					resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "role_name", roleName),
					resource.TestCheckResourceAttr(resourceName, "actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.action", "UPDATE"),
					resource.TestCheckResourceAttr(resourceName, "actions.0.resources.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCustomDBRoles_WithInheritedRoles(t *testing.T) {
	testRoleResourceName := "mongodbatlas_custom_db_role.test_role"
	InheritedRoleResourceNameOne := "mongodbatlas_custom_db_role.inherited_role_one"
	InheritedRoleResourceNameTwo := "mongodbatlas_custom_db_role.inherited_role_two"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	inheritRole := []matlas.CustomDBRole{
		{
			RoleName: fmt.Sprintf("test-acc-INHERITED_ROLE-%s", acctest.RandString(5)),
			Actions: []matlas.Action{{
				Action: "INSERT",
				Resources: []matlas.Resource{{
					Db: fmt.Sprintf("b_test-acc-ddb_name-%s", acctest.RandString(5)),
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
				Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
			}},
		}},
	}

	inheritRoleUpdated := []matlas.CustomDBRole{
		{
			RoleName: inheritRole[0].RoleName,
			Actions: []matlas.Action{{
				Action: "FIND",
				Resources: []matlas.Resource{{
					Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
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
				Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
			}},
		}},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigWithInheritedRoles(projectID, inheritRole, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Roles
					// inherited Role [0]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameOne),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "role_name", inheritRole[0].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRole[0].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRole[0].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRole[0].Actions[0].Resources))),

					// inherited Role [1]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRole[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRole[1].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRole[1].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRole[1].Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(testRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.Actions[0].Resources))),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigWithInheritedRoles(projectID, inheritRoleUpdated, testRoleUpdated),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					// inherited Role [0]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameOne),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameOne, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "role_name", inheritRoleUpdated[0].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.#", cast.ToString(len(inheritRoleUpdated[0].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.action", inheritRoleUpdated[0].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameOne, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[0].Actions[0].Resources))),

					// inherited Role [1]
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceNameTwo),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceNameTwo, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "role_name", inheritRoleUpdated[1].RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.#", cast.ToString(len(inheritRoleUpdated[1].Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.action", inheritRoleUpdated[1].Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceNameTwo, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated[1].Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(testRoleResourceName, "project_id", projectID),
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

func TestAccResourceMongoDBAtlasCustomDBRoles_MultipleCustomRoles(t *testing.T) {
	var (
		testRoleResourceName      = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceName = "mongodbatlas_custom_db_role.inherited_role"
		projectID                 = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	inheritRole := &matlas.CustomDBRole{
		RoleName: fmt.Sprintf("test-acc-INHERITED_ROLE-%s", acctest.RandString(5)),
		Actions: []matlas.Action{
			{
				Action: "REMOVE",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
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
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
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
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
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
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(projectID, inheritRole, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRole.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(testRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.#", cast.ToString(len(testRole.Actions))),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.action", testRole.Actions[0].Action),
					resource.TestCheckResourceAttr(testRoleResourceName, "actions.0.resources.#", cast.ToString(len(testRole.Actions[0].Resources))),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(projectID, inheritRoleUpdated, testRoleUpdated),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(testRoleResourceName, "project_id", projectID),
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

func TestAccResourceMongoDBAtlasCustomDBRoles_MultipleResources(t *testing.T) {
	t.Skip()// The error seems appear to be similar to whitelist behavior, skip it temporally
	var (
		resourceName = "mongodbatlas_custom_db_role.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		roleName     string
	)

	for i := 0; i < 100; i++ {
		roleName = fmt.Sprintf("test-acc-custom_role-%d", i)

		t.Run(roleName, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckMongoDBAtlasCustomDBRolesDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(projectID, roleName, "INSERT", fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName),
							resource.TestCheckResourceAttrSet(resourceName, "project_id"),
							resource.TestCheckResourceAttrSet(resourceName, "role_name"),
							resource.TestCheckResourceAttrSet(resourceName, "actions.0.action"),

							resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

func TestAccResourceMongoDBAtlasCustomDBRoles_importBasic(t *testing.T) {
	SkipTestImport(t)
	var (
		resourceName = "mongodbatlas_custom_db_role.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		roleName     = fmt.Sprintf("test-acc-custom_role-%s", acctest.RandString(5))
		databaseName = fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5))
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigBasic(projectID, roleName, "INSERT", databaseName),
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

func TestAccResourceMongoDBAtlasCustomDBRoles_UpdatedInheritRoles(t *testing.T) {
	var (
		testRoleResourceName      = "mongodbatlas_custom_db_role.test_role"
		InheritedRoleResourceName = "mongodbatlas_custom_db_role.inherited_role"
		projectID                 = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	inheritRole := &matlas.CustomDBRole{
		RoleName: fmt.Sprintf("test-acc-INHERITED_ROLE-%s", acctest.RandString(5)),
		Actions: []matlas.Action{
			{
				Action: "REMOVE",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
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
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
				},
			},
			{
				Action: "FIND",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
				},
			},
			{
				Action: "INSERT",
				Resources: []matlas.Resource{
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
					},
					{
						Db: fmt.Sprintf("test-acc-db_name-%s", acctest.RandString(5)),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCustomDBRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(projectID, inheritRole, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRole.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRole.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRole.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRole.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),

					resource.TestCheckResourceAttr(testRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasCustomDBRolesConfigMultiple(projectID, inheritRoleUpdated, testRole),
				Check: resource.ComposeTestCheckFunc(

					// For Inherited Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(InheritedRoleResourceName),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "role_name"),
					resource.TestCheckResourceAttrSet(InheritedRoleResourceName, "actions.0.action"),

					resource.TestCheckResourceAttr(InheritedRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "role_name", inheritRoleUpdated.RoleName),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.#", cast.ToString(len(inheritRoleUpdated.Actions))),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.action", inheritRoleUpdated.Actions[0].Action),
					resource.TestCheckResourceAttr(InheritedRoleResourceName, "actions.0.resources.#", cast.ToString(len(inheritRoleUpdated.Actions[0].Resources))),

					// For Test Role
					testAccCheckMongoDBAtlasCustomDBRolesExists(testRoleResourceName),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "project_id"),
					resource.TestCheckResourceAttrSet(testRoleResourceName, "role_name"),

					resource.TestCheckResourceAttr(testRoleResourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(testRoleResourceName, "role_name", testRole.RoleName),
					resource.TestCheckResourceAttr(testRoleResourceName, "inherited_roles.#", "1"),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasCustomDBRolesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

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
	conn := testAccProvider.Meta().(*matlas.Client)

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

func testAccMongoDBAtlasCustomDBRolesConfigBasic(projectID, roleName, action, databaseName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_custom_db_role" "test" {
			project_id = "%s"
			role_name  = "%s"

			actions {
				action = "%s"
				resources {
					collection_name = ""
					database_name   = "%s"
				}
			}
		}
	`, projectID, roleName, action, databaseName)
}

func testAccMongoDBAtlasCustomDBRolesConfigWithInheritedRoles(projectID string, inheritedRole []matlas.CustomDBRole, testRole *matlas.CustomDBRole) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_custom_db_role" "inherited_role_one" {
		 	project_id = "%s"
		 	role_name  = "%s"

			actions {
				action = "%s"
				resources {
					collection_name = ""
					database_name   = "%s"
				}
			}
		}

		resource "mongodbatlas_custom_db_role" "inherited_role_two" {
			project_id = "${mongodbatlas_custom_db_role.inherited_role_one.project_id}"
		 	role_name  = "%s"

			actions {
				action = "%s"
				resources {
					cluster = %t
				}
			}
		}

		resource "mongodbatlas_custom_db_role" "test_role" {
			project_id = "${mongodbatlas_custom_db_role.inherited_role_one.project_id}"
			role_name  = "%s"

			actions {
				action = "%s"
				resources {
					collection_name = ""
					database_name   = "%s"
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
	`, projectID,
		inheritedRole[0].RoleName, inheritedRole[0].Actions[0].Action, inheritedRole[0].Actions[0].Resources[0].Db,
		inheritedRole[1].RoleName, inheritedRole[1].Actions[0].Action, *inheritedRole[1].Actions[0].Resources[0].Cluster,
		testRole.RoleName, testRole.Actions[0].Action, testRole.Actions[0].Resources[0].Db,
	)
}

func testAccMongoDBAtlasCustomDBRolesConfigMultiple(projectID string, inheritedRole, testRole *matlas.CustomDBRole) string {
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
		resource "mongodbatlas_custom_db_role" "inherited_role" {
		 	project_id = "%s"
		 	role_name  = "%s"

			%s
		}

		resource "mongodbatlas_custom_db_role" "test_role" {
			project_id = "${mongodbatlas_custom_db_role.inherited_role.project_id}"
			role_name  = "%s"

			%s

			%s
		}
	`, projectID,
		inheritedRole.RoleName, getCustomRoleFields(inheritedRole)["actions"],
		testRole.RoleName, getCustomRoleFields(testRole)["actions"], getCustomRoleFields(testRole)["inherited_roles"],
	)
}
