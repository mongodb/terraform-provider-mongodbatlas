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
		username     = os.Getenv("MONGODB_ATLAS_DB_USERNAME")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			func() {
				if os.Getenv("MONGODB_ATLAS_DB_USERNAME") == "" {
					t.Fatal("`MONGODB_ATLAS_DB_USERNAME` must be set for MongoDB Atlas X509 Authentication Database users  testing")
				}
			}()
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
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

func TestAccResourceMongoDBAtlasX509AuthDBUser_WithCustomerX509(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		cas          = `
-----BEGIN CERTIFICATE-----
MIIEBzCCAu+gAwIBAgIUbwfQS97LZAIc/tPPVRYEWKkDETkwDQYJKoZIhvcNAQEL
BQAwgZIxCzAJBgNVBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5l
eTEQMA4GA1UECgwHTW9uZ29EQjENMAsGA1UECwwEcm9vdDEbMBkGA1UEAwwScG9w
LW9zLmxvY2FsZG9tYWluMSYwJAYJKoZIhvcNAQkBFhdlZGdhci5sb3BlekBtb25n
b2RiLmNvbTAeFw0yMTAyMjIxODE4NTFaFw0yMTAzMjQxODE4NTFaMIGSMQswCQYD
VQQGEwJBVTEMMAoGA1UECAwDTlNXMQ8wDQYDVQQHDAZTeWRuZXkxEDAOBgNVBAoM
B01vbmdvREIxDTALBgNVBAsMBHJvb3QxGzAZBgNVBAMMEnBvcC1vcy5sb2NhbGRv
bWFpbjEmMCQGCSqGSIb3DQEJARYXZWRnYXIubG9wZXpAbW9uZ29kYi5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQClBpr0bTP5AElcNmSLC5ioMYJZ
3LyBCTtcz2YDvrFoaN4UKUvD5pXTkkhSdHRgIpVbvibHQl118haj+gcN9s1GP0lb
6Lz5XPOs6QdjO4fGz9M8WOPFWuAiXmDqqGhobVcSdEFtddtCPE0PAsKmzVBuDd/1
RYGskzLC94f0SL9YYmF6kqXKXTH+D7JHpUWqCms3RCKIc2AYdlU0LD1dqyjabbWN
I2PS4j6xQca9ZfpqlHvUxwAzLuaMAZYHDUQ++uVJi/iHY7Dd2/PA41sUT/ymwmJH
4Zc4Nd73WFtYUBQxHa3sNfhiNFZ4BW6LkBGcPV6+r5AZIe3ZiEP1MuKim91XAgMB
AAGjUzBRMB0GA1UdDgQWBBRTwbI8Tx+JDkNUn7k+JHZ/HHQAnTAfBgNVHSMEGDAW
gBRTwbI8Tx+JDkNUn7k+JHZ/HHQAnTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQAwObUHhCqt9UZAuqvke+HMU9FRiNehEKm+1JygZ2SQSPAnlR2S
+ttihCxjyU5wVgEs8lo/VoPxtc8TfA/rKYn3xhTdqSo9nSOZVS+F8OA2A5ClnTc8
U1l6t5KeQTKOsGFdyc5xzRM2P0nAY6WWB2PpFyPlwCJtPIq2l0my3W2za8m7DCLO
XXJ5LYUZ8kFJKFEnS7F++A7102+tOs/GSeXwg3u9aYhhjwgsHneWzW5YLOtDqIPg
ulnNdinFfGNo57BqKbRlwqdU0HIHLAZPXnftQbKMamxNw2IN169yPJPRNfdPvL/S
5niRVVHHXMoEYp8L9KZ3aJQODnxY05IbTEjP
-----END CERTIFICATE-----
`
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
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

func TestAccResourceMongoDBAtlasX509AuthDBUser_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		username     = os.Getenv("MONGODB_ATLAS_DB_USERNAME")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			func() {
				if os.Getenv("MONGODB_ATLAS_DB_USERNAME") == "" {
					t.Fatal("`MONGODB_ATLAS_DB_USERNAME` must be set for MongoDB Atlas X509 Authentication Database users  testing")
				}
			}()
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
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

func TestAccResourceMongoDBAtlasX509AuthDBUser_WithDatabaseUser(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		months       = acctest.RandIntRange(1, 24)
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
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

func TestAccResourceMongoDBAtlasX509AuthDBUser_importWithCustomerX509(t *testing.T) {
	var (
		resourceName = "mongodbatlas_x509_authentication_database_user.test"
		cas          = `
-----BEGIN CERTIFICATE-----
MIIEBzCCAu+gAwIBAgIUbwfQS97LZAIc/tPPVRYEWKkDETkwDQYJKoZIhvcNAQEL
BQAwgZIxCzAJBgNVBAYTAkFVMQwwCgYDVQQIDANOU1cxDzANBgNVBAcMBlN5ZG5l
eTEQMA4GA1UECgwHTW9uZ29EQjENMAsGA1UECwwEcm9vdDEbMBkGA1UEAwwScG9w
LW9zLmxvY2FsZG9tYWluMSYwJAYJKoZIhvcNAQkBFhdlZGdhci5sb3BlekBtb25n
b2RiLmNvbTAeFw0yMTAyMjIxODE4NTFaFw0yMTAzMjQxODE4NTFaMIGSMQswCQYD
VQQGEwJBVTEMMAoGA1UECAwDTlNXMQ8wDQYDVQQHDAZTeWRuZXkxEDAOBgNVBAoM
B01vbmdvREIxDTALBgNVBAsMBHJvb3QxGzAZBgNVBAMMEnBvcC1vcy5sb2NhbGRv
bWFpbjEmMCQGCSqGSIb3DQEJARYXZWRnYXIubG9wZXpAbW9uZ29kYi5jb20wggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQClBpr0bTP5AElcNmSLC5ioMYJZ
3LyBCTtcz2YDvrFoaN4UKUvD5pXTkkhSdHRgIpVbvibHQl118haj+gcN9s1GP0lb
6Lz5XPOs6QdjO4fGz9M8WOPFWuAiXmDqqGhobVcSdEFtddtCPE0PAsKmzVBuDd/1
RYGskzLC94f0SL9YYmF6kqXKXTH+D7JHpUWqCms3RCKIc2AYdlU0LD1dqyjabbWN
I2PS4j6xQca9ZfpqlHvUxwAzLuaMAZYHDUQ++uVJi/iHY7Dd2/PA41sUT/ymwmJH
4Zc4Nd73WFtYUBQxHa3sNfhiNFZ4BW6LkBGcPV6+r5AZIe3ZiEP1MuKim91XAgMB
AAGjUzBRMB0GA1UdDgQWBBRTwbI8Tx+JDkNUn7k+JHZ/HHQAnTAfBgNVHSMEGDAW
gBRTwbI8Tx+JDkNUn7k+JHZ/HHQAnTAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQAwObUHhCqt9UZAuqvke+HMU9FRiNehEKm+1JygZ2SQSPAnlR2S
+ttihCxjyU5wVgEs8lo/VoPxtc8TfA/rKYn3xhTdqSo9nSOZVS+F8OA2A5ClnTc8
U1l6t5KeQTKOsGFdyc5xzRM2P0nAY6WWB2PpFyPlwCJtPIq2l0my3W2za8m7DCLO
XXJ5LYUZ8kFJKFEnS7F++A7102+tOs/GSeXwg3u9aYhhjwgsHneWzW5YLOtDqIPg
ulnNdinFfGNo57BqKbRlwqdU0HIHLAZPXnftQbKMamxNw2IN169yPJPRNfdPvL/S
5niRVVHHXMoEYp8L9KZ3aJQODnxY05IbTEjP
-----END CERTIFICATE-----`
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasX509AuthDBUserDestroy,
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
			project_id         = "${mongodbatlas_project.test.id}"
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
