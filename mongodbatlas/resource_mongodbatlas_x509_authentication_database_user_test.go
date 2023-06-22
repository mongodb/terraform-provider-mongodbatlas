package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/spf13/cast"
)

func TestAccGenericAdvRSX509AuthDBUser_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckBasic(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfig(projectName, orgID, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
			},
		},
	})
}

func TestAccGenericAdvRSX509AuthDBUser_WithCustomerX509(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		cas          = os.Getenv("CA_CERT")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfigWithCustomerX509(projectName, orgID, cas),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_x509_cas"),
				),
			},
		},
	})
}

func TestAccGenericAdvRSX509AuthDBUser_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		username     = acctest.RandomWithPrefix("test-acc")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckBasic(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfig(projectName, orgID, username),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasX509AuthDBUserImportStateIDFuncBasic(resourceName),
				ImportState:       true,
			},
		},
	})
}

func TestAccGenericAdvRSX509AuthDBUser_WithDatabaseUser(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		months       = acctest.RandIntRange(1, 24)
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfigWithDatabaseUser(projectName, orgID, username, months),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "months_until_expiration"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "months_until_expiration", cast.ToString(months)),
				),
			},
		},
	})
}

func TestAccGenericAdvRSX509AuthDBUser_importWithCustomerX509(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		cas          = os.Getenv("CA_CERT")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfigWithCustomerX509(projectName, orgID, cas),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasX509AuthDBUserImportStateIDFuncBasic(resourceName),
				ImportState:       true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasX509AuthDBUserImportStateIDFuncBasic(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["username"]), nil
	}
}

func testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		if ids["current_certificate"] != "" {
			if _, _, err := conn.X509AuthDBUsers.GetUserCertificates(context.Background(), ids["project_id"], ids["username"], nil); err == nil {
				return nil
			}

			return fmt.Errorf("the X509 Authentication Database User(%s) does not exist in the project(%s)", ids["username"], ids["project_id"])
		}

		if _, _, err := conn.X509AuthDBUsers.GetCurrentX509Conf(context.Background(), ids["project_id"]); err == nil {
			return nil
		}

		return fmt.Errorf("the Customer X509 Authentication does not exist in the project(%s)", ids["project_id"])
	}
}

func testAccCheckMongoDBAtlasX509AuthDBUserDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_x509_authentication_database_user" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		if ids["current_certificate"] != "" {
			_, _, err := conn.X509AuthDBUsers.GetUserCertificates(context.Background(), ids["project_id"], ids["username"], nil)
			if err == nil {
				/*
					There is no way to remove one user certificate so until this comes it will keep in this way
				*/
				return nil
			}
		}

		if _, _, err := conn.X509AuthDBUsers.GetCurrentX509Conf(context.Background(), ids["project_id"]); err == nil {
			return nil
		}
	}

	return nil
}

func testAccMongoDBAtlasX509AuthDBUserConfig(projectName, orgID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "basic_ds" {
			username           = "%s"
			project_id         = "${mongodbatlas_project.test.id}"
			auth_database_name = "$external"
			x509_type          = "MANAGED"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}

		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id              = "${mongodbatlas_project.test.id}"
			username                = "${mongodbatlas_database_user.basic_ds.username}"
			months_until_expiration = 5
		}
	`, projectName, orgID, username)
}

func testAccMongoDBAtlasX509AuthDBUserConfigWithCustomerX509(projectName, orgID, cas string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id        = "${mongodbatlas_project.test.id}"
			customer_x509_cas = <<-EOT
			%s
			EOT
		}
	`, projectName, orgID, cas)
}

func testAccMongoDBAtlasX509AuthDBUserConfigWithDatabaseUser(projectName, orgID, username string, months int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "user" {
			project_id         = mongodbatlas_project.test.id
			username           = "%s"
			x509_type          = "MANAGED"
			auth_database_name = "$external"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}

			labels {
				key   = "My Key"
				value = "My Value"
			}
		}

		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id              = "${mongodbatlas_database_user.user.project_id}"
			username                = "${mongodbatlas_database_user.user.username}"
			months_until_expiration = %d
		}
	`, projectName, orgID, username, months)
}
