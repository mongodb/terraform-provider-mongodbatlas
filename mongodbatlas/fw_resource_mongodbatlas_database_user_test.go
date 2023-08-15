package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSDatabaseUser_basic(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.basic_ds"
		username     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
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
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectName, orgID, "read", username, "Second Key", "Second value"),
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

func TestAccConfigRSDatabaseUser_withX509TypeCustomer(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=ellen@example.com,OU=users,DC=example,DC=com"
		x509Type     = "CUSTOMER"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithX509TypeConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withX509TypeManaged(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc")
		x509Type     = "MANAGED"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithX509TypeConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withAWSIAMType(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithAWSIAMTypeConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "aws_iam_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withAWSIAMType_import(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = os.Getenv("TEST_DB_USER_IAM_ARN")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	if username == "" {
		username = "arn:aws:iam::358363220050:user/mongodb-aws-iam-auth-test-user"
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithAWSIAMTypeConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "aws_iam_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
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

func TestAccConfigRSDatabaseUser_WithLabels(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectName, orgID, "atlasAdmin", username, []matlas.Label{}),
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
				Config: testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectName, orgID, "atlasAdmin", username,
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
				Config: testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectName, orgID, "read", username,
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

func TestAccConfigRSDatabaseUser_withRoles(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithRoles(username, password, projectName, orgID,
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
				Config: testAccMongoDBAtlasDatabaseUserWithRoles(username, password, projectName, orgID,
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

func TestAccConfigRSDatabaseUser_withScopes(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
						},
						{
							Name: "test-acc-nurk4llu2z",
							Type: "DATA_LAKE",
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
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
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
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withScopesAndEmpty(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc-user-")
		password     = acctest.RandomWithPrefix("test-acc-pass-")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{
						{
							Name: "test-acc-nurk4llu2z",
							Type: "CLUSTER",
						},
						{
							Name: "test-acc-nurk4llu2z",
							Type: "DATA_LAKE",
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
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasDatabaseUserWithScopes(username, password, projectName, orgID, "atlasAdmin", clusterName,
					[]*matlas.Scope{},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "password", password),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "admin"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "0"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_withLDAPAuthType(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=david@example.com,OU=users,DC=example,DC=com"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithLDAPAuthTypeConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "ldap_auth_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigRSDatabaseUser_importBasic(t *testing.T) {
	var (
		username     = fmt.Sprintf("test-username-%s", acctest.RandString(5))
		resourceName = "mongodbatlas_database_user.basic_ds"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserConfig(projectName, orgID, "read", username, "First Key", "First value"),
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

func TestAccConfigRSDatabaseUser_importX509TypeCustomer(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=ellen@example.com,OU=users,DC=example,DC=com"
		x509Type     = "CUSTOMER"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithX509TypeConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value", x509Type),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "x509_type", x509Type),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
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

func TestAccConfigRSDatabaseUser_importLDAPAuthType(t *testing.T) {
	var (
		dbUser       matlas.DatabaseUser
		resourceName = "mongodbatlas_database_user.test"
		username     = "CN=david@example.com,OU=users,DC=example,DC=com"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasDatabaseUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDatabaseUserWithLDAPAuthTypeConfig(projectName, orgID, "atlasAdmin", username, "First Key", "First value"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasDatabaseUserExists(resourceName, &dbUser),
					testAccCheckMongoDBAtlasDatabaseUserAttributes(&dbUser, username),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "ldap_auth_type", "USER"),
					resource.TestCheckResourceAttr(resourceName, "auth_database_name", "$external"),
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
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["username"], ids["auth_database_name"]), nil
	}
}

func testAccCheckMongoDBAtlasDatabaseUserExists(resourceName string, dbUser *matlas.DatabaseUser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProviderSdkV2.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		username := ids["username"]

		if dbUserResp, _, err := conn.DatabaseUsers.Get(context.Background(), ids["auth_database_name"], ids["project_id"], username); err == nil {
			*dbUser = *dbUserResp
			return nil
		}

		return fmt.Errorf("database user(%s) does not exist", ids["project_id"])
	}
}

func testAccCheckMongoDBAtlasDatabaseUserAttributes(dbUser *matlas.DatabaseUser, username string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[DEBUG] difference dbUser.Username: %s , username : %s", dbUser.Username, username)
		if dbUser.Username != username {
			return fmt.Errorf("bad username: %s", dbUser.Username)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasDatabaseUserDestroy(s *terraform.State) error {
	conn := testAccProviderSdkV2.Meta().(*MongoDBClient).Atlas

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

func testAccMongoDBAtlasDatabaseUserConfig(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "basic_ds" {
			username           = "%[4]s"
			password           = "test-acc-password"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel)
}

func testAccMongoDBAtlasDatabaseUserWithX509TypeConfig(projectName, orgID, roleName, username, keyLabel, valueLabel, x509Type string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			x509_type          = "%[7]s"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "$external"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel, x509Type)
}

func testAccMongoDBAtlasDatabaseUserWithLabelsConfig(projectName, orgID, roleName, username string, labels []matlas.Label) string {
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
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			password           = "test-acc-password"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			%[5]s

		}
	`, projectName, orgID, roleName, username, labelsConf)
}

func testAccMongoDBAtlasDatabaseUserWithRoles(username, password, projectName, orgID string, rolesArr []*matlas.Role) string {
	var roles string

	for _, role := range rolesArr {
		var roleName, databaseName, collection string

		if role.RoleName != "" {
			roleName = fmt.Sprintf(`role_name = %q`, role.RoleName)
		}

		if role.DatabaseName != "" {
			databaseName = fmt.Sprintf(`database_name = %q`, role.DatabaseName)
		}

		if role.CollectionName != "" {
			collection = fmt.Sprintf(`collection_name = %q`, role.CollectionName)
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
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%s"
			password           = "%s"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			%s

		}
	`, projectName, orgID, username, password, roles)
}

func testAccMongoDBAtlasDatabaseUserWithAWSIAMTypeConfig(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			aws_iam_type       = "USER"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "$external"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel)
}

func testAccMongoDBAtlasDatabaseUserWithScopes(username, password, projectName, orgID, roleName, clusterName string, scopesArr []*matlas.Scope) string {
	var scopes string

	for _, scope := range scopesArr {
		var scopeType string

		if scope.Type != "" {
			scopeType = fmt.Sprintf(`type = %q`, scope.Type)
		}

		scopes += fmt.Sprintf(`
			scopes {
				name = "${mongodbatlas_cluster.my_cluster.name}"
				%s
			}
		`, scopeType)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "${mongodbatlas_project.test.id}"
			name         = "%s"
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%s"
			password           = "%s"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "admin"

			roles {
				role_name     = "%s"
				database_name = "admin"
			}

			%s

		}
	`, projectName, orgID, clusterName, username, password, roleName, scopes)
}

func testAccMongoDBAtlasDatabaseUserWithLDAPAuthTypeConfig(projectName, orgID, roleName, username, keyLabel, valueLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "test" {
			username           = "%[4]s"
			ldap_auth_type     = "USER"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "$external"

			roles {
				role_name     = "%[3]s"
				database_name = "admin"
			}

			labels {
				key   = "%s"
				value = "%s"
			}
		}
	`, projectName, orgID, roleName, username, keyLabel, valueLabel)
}
