package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/spf13/cast"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasX509AuthDBUser_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		username     = os.Getenv("DB_USERNAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			func() {
				if os.Getenv("DB_USERNAME") == "" {
					t.Fatal("`DB_USERNAME` must be set for MongoDB Atlas X509 Authentication Database users  testing")
				}
			}()
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfig(projectID, username),
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

func TestAccResourceMongoDBAtlasX509AuthDBUser_WithCustomerX509(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		cas          = `
		-----BEGIN CERTIFICATE-----
		MIICmTCCAgICCQDZnHzklxsT9TANBgkqhkiG9w0BAQsFADCBkDELMAkGA1UEBhMC
		VVMxDjAMBgNVBAgMBVRleGFzMQ8wDQYDVQQHDAZBdXN0aW4xETAPBgNVBAoMCHRl
		c3QuY29tMQ0wCwYDVQQLDARUZXN0MREwDwYDVQQDDAh0ZXN0LmNvbTErMCkGCSqG
		SIb3DQEJARYcbWVsaXNzYS5wbHVua2V0dEBtb25nb2RiLmNvbTAeFw0yMDAyMDQy
		MDQ2MDFaFw0yMTAyMDMyMDQ2MDFaMIGQMQswCQYDVQQGEwJVUzEOMAwGA1UECAwF
		VGV4YXMxDzANBgNVBAcMBkF1c3RpbjERMA8GA1UECgwIdGVzdC5jb20xDTALBgNV
		BAsMBFRlc3QxETAPBgNVBAMMCHRlc3QuY29tMSswKQYJKoZIhvcNAQkBFhxtZWxp
		c3NhLnBsdW5rZXR0QG1vbmdvZGIuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
		iQKBgQCf1LRqr1zftzdYx2Aj9G76tb0noMPtj6faGLlPji1+m6Rn7RWD9L0ntWAr
		cURxvypa9jZ9MXFzDtLevvd3tHEmfrUT3ukNDX6+Jtc4kWm+Dh2A70Pd+deKZ2/O
		Fh8audEKAESGXnTbeJCeQa1XKlIkjqQHBNwES5h1b9vJtFoLJwIDAQABMA0GCSqG
		SIb3DQEBCwUAA4GBADMUncjEPV/MiZUcVNGmktP6BPmEqMXQWUDpdGW2+Tg2JtUA
		7MMILtepBkFzLO+GlpZxeAlXO0wxiNgEmCRONgh4+t2w3e7a8GFijYQ99FHrAC5A
		iul59bdl18gVqXia1Yeq/iK7Ohfy/Jwd7Hsm530elwkM/ZEkYDjBlZSXYdyz
		-----END CERTIFICATE-----`
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfigWithCustomerX509(projectID, cas),
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

func TestAccResourceMongoDBAtlasX509AuthDBUser_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		username     = os.Getenv("DB_USERNAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			func() {
				if os.Getenv("DB_USERNAME") == "" {
					t.Fatal("`DB_USERNAME` must be set for MongoDB Atlas X509 Authentication Database users  testing")
				}
			}()
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfig(projectID, username),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasX509AuthDBUserImportStateIDFuncBasic(resourceName),
				ImportState:       true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasX509AuthDBUser_WithDatabaseUser(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		username     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		months       = acctest.RandIntRange(1, 24)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfigWithDatabaseUser(projectID, username, months),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasX509AuthDBUserExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "username"),
					resource.TestCheckResourceAttrSet(resourceName, "months_until_expiration"),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttr(resourceName, "months_until_expiration", cast.ToString(months)),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasX509AuthDBUser_importWithCustomerX509(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		cas          = `
		-----BEGIN CERTIFICATE-----
		MIICmTCCAgICCQDZnHzklxsT9TANBgkqhkiG9w0BAQsFADCBkDELMAkGA1UEBhMC
		VVMxDjAMBgNVBAgMBVRleGFzMQ8wDQYDVQQHDAZBdXN0aW4xETAPBgNVBAoMCHRl
		c3QuY29tMQ0wCwYDVQQLDARUZXN0MREwDwYDVQQDDAh0ZXN0LmNvbTErMCkGCSqG
		SIb3DQEJARYcbWVsaXNzYS5wbHVua2V0dEBtb25nb2RiLmNvbTAeFw0yMDAyMDQy
		MDQ2MDFaFw0yMTAyMDMyMDQ2MDFaMIGQMQswCQYDVQQGEwJVUzEOMAwGA1UECAwF
		VGV4YXMxDzANBgNVBAcMBkF1c3RpbjERMA8GA1UECgwIdGVzdC5jb20xDTALBgNV
		BAsMBFRlc3QxETAPBgNVBAMMCHRlc3QuY29tMSswKQYJKoZIhvcNAQkBFhxtZWxp
		c3NhLnBsdW5rZXR0QG1vbmdvZGIuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
		iQKBgQCf1LRqr1zftzdYx2Aj9G76tb0noMPtj6faGLlPji1+m6Rn7RWD9L0ntWAr
		cURxvypa9jZ9MXFzDtLevvd3tHEmfrUT3ukNDX6+Jtc4kWm+Dh2A70Pd+deKZ2/O
		Fh8audEKAESGXnTbeJCeQa1XKlIkjqQHBNwES5h1b9vJtFoLJwIDAQABMA0GCSqG
		SIb3DQEBCwUAA4GBADMUncjEPV/MiZUcVNGmktP6BPmEqMXQWUDpdGW2+Tg2JtUA
		7MMILtepBkFzLO+GlpZxeAlXO0wxiNgEmCRONgh4+t2w3e7a8GFijYQ99FHrAC5A
		iul59bdl18gVqXia1Yeq/iK7Ohfy/Jwd7Hsm530elwkM/ZEkYDjBlZSXYdyz
		-----END CERTIFICATE-----`
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasX509AuthDBUserConfigWithCustomerX509(projectID, cas),
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
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		if ids["current_certificate"] != "" {
			if _, _, err := conn.X509AuthDBUsers.GetUserCertificates(context.Background(), ids["project_id"], ids["username"]); err == nil {
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
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_x509_authentication_database_user" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		if ids["current_certificate"] != "" {
			_, _, err := conn.X509AuthDBUsers.GetUserCertificates(context.Background(), ids["project_id"], ids["username"])
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

func testAccMongoDBAtlasX509AuthDBUserConfig(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id              = "%s"
			username                = "%s"
			months_until_expiration = 5
		}
	`, projectID, username)
}

func testAccMongoDBAtlasX509AuthDBUserConfigWithCustomerX509(projectID, cas string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id        = "%s"
			customer_x509_cas = <<-EOT
			%s
			EOT
		}
	`, projectID, cas)
}

func testAccMongoDBAtlasX509AuthDBUserConfigWithDatabaseUser(projectID, username string, months int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "user" {
			project_id         = "%s"
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
	`, projectID, username, months)
}
