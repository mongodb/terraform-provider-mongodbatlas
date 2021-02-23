package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
	cas := `
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
			customer_x509_cas = <<EOT
								%s
								EOT
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = "${mongodbatlas_x509_authentication_database_user.test.project_id}"
		}
	`, projectID, username)
}
