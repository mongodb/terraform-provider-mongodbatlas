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
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		username    = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckBasic(t)
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, username),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "username", username),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
				),
			},
		},
	})
}

func TestAccGenericX509AuthDBUser_withCustomerX509(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_x509_authentication_database_user.test"
		cas            = os.Getenv("CA_CERT")
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCert(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithCustomerX509(projectName, orgID, cas),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "customer_x509_cas"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "customer_x509_cas"),
				),
			},
		},
	})
}

func TestAccGenericX509AuthDBUser_importBasic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		username    = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheckBasic(t)
		},
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, username),
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
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username    = acc.RandomName()
		months      = acctest.RandIntRange(1, 24)
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithDatabaseUser(projectName, orgID, username, months),
				Check: resource.ComposeTestCheckFunc(
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

func TestAccGenericX509AuthDBUser_importWithCustomerX509(t *testing.T) {
	var (
		cas         = os.Getenv("CA_CERT")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckCert(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithCustomerX509(projectName, orgID, cas),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFuncBasic(resourceName),
				ImportState:       true,
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
			if _, _, err := acc.ConnV2().X509AuthenticationApi.ListDatabaseUserCertificates(context.Background(), ids["project_id"], ids["username"]).Execute(); err == nil {
				return nil
			}
			return fmt.Errorf("the X509 Authentication Database User(%s) does not exist in the project(%s)", ids["username"], ids["project_id"])
		}
		if _, _, err := acc.ConnV2().LDAPConfigurationApi.GetLDAPConfiguration(context.Background(), ids["project_id"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("the Customer X509 Authentication does not exist in the project(%s)", ids["project_id"])
	}
}

func configBasic(projectName, orgID, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_database_user" "basic_ds" {
			username           = "%s"
			project_id         = mongodbatlas_project.test.id
			auth_database_name = "$external"
			x509_type          = "MANAGED"

			roles {
				role_name     = "atlasAdmin"
				database_name = "admin"
			}
		}

		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id              = mongodbatlas_project.test.id
			username                = mongodbatlas_database_user.basic_ds.username
			months_until_expiration = 5
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = mongodbatlas_x509_authentication_database_user.test.project_id
			username   = mongodbatlas_x509_authentication_database_user.test.username
		}
	`, projectName, orgID, username)
}

func configWithCustomerX509(projectName, orgID, cas string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%s"
			org_id = "%s"
		}

		resource "mongodbatlas_x509_authentication_database_user" "test" {
			project_id        = mongodbatlas_project.test.id
			customer_x509_cas = <<-EOT
			%s
			EOT
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = mongodbatlas_x509_authentication_database_user.test.project_id
			username   = mongodbatlas_x509_authentication_database_user.test.username
		}
	`, projectName, orgID, cas)
}

func configWithDatabaseUser(projectName, orgID, username string, months int) string {
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
			project_id              = mongodbatlas_database_user.user.project_id
			username                = mongodbatlas_database_user.user.username
			months_until_expiration = %d
		}

		data "mongodbatlas_x509_authentication_database_user" "test" {
			project_id = mongodbatlas_x509_authentication_database_user.test.project_id
			username   = mongodbatlas_x509_authentication_database_user.test.username
		}
	`, projectName, orgID, username, months)
}
