package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoX509AuthDBUser_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_x509_authentication_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	username := os.Getenv("MONGODB_ATLAS_DB_USERNAME")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			func() {
				if os.Getenv("MONGODB_ATLAS_DB_USERNAME") == "" {
					t.Fatal("`MONGODB_ATLAS_DB_USERNAME` must be set for MongoDB Atlas X509 Authentication Database users  testing")
				}
			}()
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoX509AuthDBUserDataSourceConfig(projectID, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
				),
			},
		},
	})
}

func TestAccDataSourceMongoX509AuthDBUser_WithCustomerX509(t *testing.T) {
	resourceName := "data.mongodbatlas_x509_authentication_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cas := os.Getenv("CA_CERT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoX509AuthDBUserDataSourceConfigWithCustomerX509(projectID, cas),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_x509_cas"),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func testAccMongoX509AuthDBUserDataSourceConfig(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = "%s"
			username   = "%s"
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = "${mongodbatlas_x509_authentication_database_user.test.project_id}"
			username   = "${mongodbatlas_x509_authentication_database_user.test.username}"
		}
	`, projectID, username)
}

func testAccMongoX509AuthDBUserDataSourceConfigWithCustomerX509(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id        = "%s"
			customer_x509_cas = <<EOT
								%s
								EOT
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = "${mongodbatlas_x509_authentication_database_user.test.project_id}"
		}
	`, projectID, username)
}
