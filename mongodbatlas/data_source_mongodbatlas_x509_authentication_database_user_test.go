package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGenericAdvDSX509AuthDBUser_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_x509_authentication_database_user.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	username := acctest.RandomWithPrefix("test-acc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckBasic(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoX509AuthDBUserDataSourceConfig(orgID, projectName, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
				),
			},
		},
	})
}

func TestAccGenericAdvDSX509AuthDBUser_WithCustomerX509(t *testing.T) {
	resourceName := "data.mongodbatlas_x509_authentication_database_user.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	cas := os.Getenv("CA_CERT")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoX509AuthDBUserDataSourceConfigWithCustomerX509(orgID, projectName, cas),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_x509_cas"),
				),
			},
		},
	})
}

func testAccMongoX509AuthDBUserDataSourceConfig(orgID, projectName, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = mongodbatlas_project.test.id
			username   = "%s"
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = "${mongodbatlas_x509_authentication_database_user.test.project_id}"
			username   = "${mongodbatlas_x509_authentication_database_user.test.username}"
		}
	`, orgID, projectName, username)
}

func testAccMongoX509AuthDBUserDataSourceConfigWithCustomerX509(orgID, projecName, cas string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id        = mongodbatlas_project.test.id
			customer_x509_cas = <<EOT
								%s
								EOT
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = "${mongodbatlas_x509_authentication_database_user.test.project_id}"
		}
	`, orgID, projecName, cas)
}
