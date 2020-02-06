package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoX509AuthDBUser_basic(t *testing.T) {
	resourceName := "data.mongodbatlas_x509_authentication_database_user.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	username := os.Getenv("DB_USERNAME")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			func() {
				if os.Getenv("DB_USERNAME") == "" {
					t.Fatal("`DB_USERNAME` must be set for MongoDB Atlas X509 Authentication Database users  testing")
				}
			}()
		},
		Providers: testAccProviders,
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
	cas := `-----BEGIN CERTIFICATE-----\nMIICmTCCAgICCQDZnHzklxsT9TANBgkqhkiG9w0BAQsFADCBkDELMAkGA1UEBhMC\nVVMxDjAMBgNVBAgMBVRleGFzMQ8wDQYDVQQHDAZBdXN0aW4xETAPBgNVBAoMCHRl\nc3QuY29tMQ0wCwYDVQQLDARUZXN0MREwDwYDVQQDDAh0ZXN0LmNvbTErMCkGCSqG\nSIb3DQEJARYcbWVsaXNzYS5wbHVua2V0dEBtb25nb2RiLmNvbTAeFw0yMDAyMDQy\nMDQ2MDFaFw0yMTAyMDMyMDQ2MDFaMIGQMQswCQYDVQQGEwJVUzEOMAwGA1UECAwF\nVGV4YXMxDzANBgNVBAcMBkF1c3RpbjERMA8GA1UECgwIdGVzdC5jb20xDTALBgNV\nBAsMBFRlc3QxETAPBgNVBAMMCHRlc3QuY29tMSswKQYJKoZIhvcNAQkBFhxtZWxp\nc3NhLnBsdW5rZXR0QG1vbmdvZGIuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB\niQKBgQCf1LRqr1zftzdYx2Aj9G76tb0noMPtj6faGLlPji1+m6Rn7RWD9L0ntWAr\ncURxvypa9jZ9MXFzDtLevvd3tHEmfrUT3ukNDX6+Jtc4kWm+Dh2A70Pd+deKZ2/O\nFh8audEKAESGXnTbeJCeQa1XKlIkjqQHBNwES5h1b9vJtFoLJwIDAQABMA0GCSqG\nSIb3DQEBCwUAA4GBADMUncjEPV/MiZUcVNGmktP6BPmEqMXQWUDpdGW2+Tg2JtUA\n7MMILtepBkFzLO+GlpZxeAlXO0wxiNgEmCRONgh4+t2w3e7a8GFijYQ99FHrAC5A\niul59bdl18gVqXia1Yeq/iK7Ohfy/Jwd7Hsm530elwkM/ZEkYDjBlZSXYdyz\n-----END CERTIFICATE-----`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
			customer_x509_cas = "%s"
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = "${mongodbatlas_x509_authentication_database_user.test.project_id}"
		}
	`, projectID, username)
}
