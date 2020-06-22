package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasDatabaseUser_basic(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "mongodbatlas_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectID, "read", username, "Second Key", "Second value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasDatabaseUser_withX509Type(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "mongodbatlas_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithX509TypeConfig(projectID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", "MANAGED"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasDatabaseUser_WithLabels(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "mongodbatlas_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	username := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectID, "atlasAdmin", username, []matlas.Label{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectID, "atlasAdmin", username,
					[]matlas.Label{
						{
							Key:   "key 1",
							Value: "value 1",
						},
						{
							Key:   "key 2",
							Value: "value 2",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectID, "read", username,
					[]matlas.Label{
						{
							Key:   "key 4",
							Value: "value 4",
						},
						{
							Key:   "key 3",
							Value: "value 3",
						},
						{
							Key:   "key 2",
							Value: "value 2",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "3"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasDatabaseUser_withRoles(t *testing.T) {
	var dbUser matlas.DatabaseUser

	resourceName := "mongodbatlas_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	username := acctest.RandomWithPrefix("test-acc-user-")
	password := acctest.RandomWithPrefix("test-acc-pass-")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithRoles(username, password, projectID,
					[]*matlas.Role{
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: "stir",
						},
						{
							RoleName:       "read",
							DatabaseName:   "admin",
							CollectionName: "unpledged",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUserWithRoles(username, password, projectID,
					[]*matlas.Role{
						{
							RoleName:     "read",
							DatabaseName: "admin",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "roles.#", "1"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasDatabaseUser_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	username := fmt.Sprintf("test-username-%s", acctest.RandString(5))
	resourceName := "mongodbatlas_database_user.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectID, "read", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", "test-acc-password"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasDatabaseUserImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasDatabaseUserImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["username"], ids["auth_database_name"]), nil
	}
}

func testAccCheckMongoDBAtlasDatabaseUserExists(resourceName string, dbUser *matlas.DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		if dbUserResp, _, err := conn.DatabaseUsers.Get(context.Background(), ids["auth_database_name"], ids["project_id"], ids["username"]); err == nil {
			*dbUser = *dbUserResp
			return nil
		}
		return fmt.Errorf("database user(%s) does not exist", ids["project_id"])
	}
}

func testAccCheckMongoDBAtlasDatabaseUserAttributes(dbUser *matlas.DatabaseUser, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if dbUser.Username != username {
			return fmt.Errorf("bad username: %s", dbUser.Username)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasDatabaseUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_database_user" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		// Try to find the database user
		_, _, err := conn.DatabaseUsers.Get(context.Background(), ids["auth_database_name"], ids["project_id"], ids["username"])
		if err == nil {
			return fmt.Errorf("database user (%s) still exists", ids["project_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasDatabaseUserConfig(projectID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			username           = "%[3]s"
			password           = "test-acc-password"
			project_id         = "%[1]s"
			auth_database_name = "admin"

			roles {
				role_name     = "%[2]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectID, roleName, username, keyLabel, valueLabel)
}

func testAccMongoDBAtlasDatabaseUserWithX509TypeConfig(projectID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			username           = "%[3]s"
			x509_type          = "MANAGED"
			project_id         = "%[1]s"
			auth_database_name = "$external"

			roles {
				role_name     = "%[2]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectID, roleName, username, keyLabel, valueLabel)
}

func testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectID, roleName, username string, labels []matlas.Label) string {

	var labelsConf string
	for _, label := range labels {
		labelsConf += fmt.Sprintf(`
			labels {
				key   = "%s"
				value = "%s"
			}
		`, label.Key, label.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			username           = "%[3]s"
			password           = "test-acc-password"
			project_id         = "%[1]s"
			auth_database_name = "admin"

			roles {
				role_name     = "%[2]s"
				database_name = "admin"
			}

			%[4]s

		}
	`, projectID, roleName, username, labelsConf)
}

func testAccMongoDBAtlasDatabaseUserWithRoles(username, password, projectID string, rolesArr []*matlas.Role) string {
	var roles string
	for _, role := range rolesArr {
		var roleName, databaseName, collection string

		if role.RoleName != "" {
			roleName = fmt.Sprintf(`role_name = "%s"`, role.RoleName)
		}
		if role.DatabaseName != "" {
			databaseName = fmt.Sprintf(`database_name = "%s"`, role.DatabaseName)
		}
		if role.CollectionName != "" {
			collection = fmt.Sprintf(`collection_name = "%s"`, role.CollectionName)
		}

		roles += fmt.Sprintf(`
			roles {
				%s
				%s
				%s
			}
		`, roleName, databaseName, collection)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "test" {
			username           = "%s"
			password           = "%s"
			project_id         = "%s"
			auth_database_name = "admin"

			%s

		}
	`, username, password, projectID, roles)
}
