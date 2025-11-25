package x509authenticationdatabaseuser_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
)

const (
	resourceName   = "mongodbatlas_x509_authentication_database_user.test"
	dataSourceName = "data.mongodbatlas_x509_authentication_database_user.test"
)

func TestAccGenericX509AuthDBUser_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		username  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckBasic(t)
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, username),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFuncBasic(resourceName),
				ImportState:       true,
			},
		},
	})
}

func TestAccGenericX509AuthDBUser_withCustomerX509(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_x509_authentication_database_user.test"
		cas            = os.Getenv("CA_CERT")
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid CANNOT_GENERATE_CERT_IF_ADVANCED_X509
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCert(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithCustomerX509(orgID, projectName, cas),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_x509_cas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "customer_x509_cas"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFuncBasic(resourceName),
				ImportState:       true,
			},
		},
	})
}

func TestAccGenericX509AuthDBUser_withDatabaseUser(t *testing.T) {
	var (
		username  = acc.RandomName()
		months    = acctest.RandIntRange(1, 24)
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithDatabaseUser(projectID, username, months),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
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

func importStateIDFuncBasic(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["username"]), nil
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		if ids["current_certificate"] != "" {
			if _, _, err := acc.ConnV2().X509AuthenticationApi.ListDatabaseUserCerts(context.Background(), ids["project_id"], ids["username"]).Execute(); err == nil {
				return nil
			}
			return fmt.Errorf("the X509 Authentication Database User(%s) does not exist in the project(%s)", ids["username"], ids["project_id"])
		}
		if _, _, err := acc.ConnV2().LDAPConfigurationApi.GetUserSecurity(context.Background(), ids["project_id"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("the Customer X509 Authentication does not exist in the project(%s)", ids["project_id"])
	}
}

func configBasic(projectID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "basic_ds" {
			project_id         = %[1]q
			username           = %[2]q
			auth_database_name = "$external"
			x509_type          = "MANAGED"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}

		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id         = %[1]q
			username                = mongodbatlas_database_user.basic_ds.username
			months_until_expiration = 5
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id         = %[1]q
			username   = mongodbatlas_x509_authentication_database_user.test.username
		}
	`, projectID, username)
}

func configWithDatabaseUser(projectID, username string, months int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_database_user" "user" {
			project_id         = %[1]q
			username           = %[2]q
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
			project_id              = mongodbatlas_database_user.user.project_id
			username                = mongodbatlas_database_user.user.username
			months_until_expiration = %[3]d
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = mongodbatlas_x509_authentication_database_user.test.project_id
			username   = mongodbatlas_x509_authentication_database_user.test.username
		}
	`, projectID, username, months)
}

func configWithCustomerX509(orgID, projectName, cas string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id        = mongodbatlas_project.test.id
			customer_x509_cas = <<-EOT
			%[3]s
			EOT
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = mongodbatlas_x509_authentication_database_user.test.project_id
			username   = mongodbatlas_x509_authentication_database_user.test.username
		}
	`, orgID, projectName, cas)
}
