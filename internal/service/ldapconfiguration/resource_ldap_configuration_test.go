package ldapconfiguration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
)

var (
	resourceName = "mongodbatlas_ldap_configuration.test"
)

func TestAccLDAPConfiguration_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		hostname    = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username    = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password    = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port        = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		authEnabled = true
		projectName = acc.RandomProjectName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyLDAPConfiguration,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, hostname, username, password, authEnabled, cast.ToInt(port)),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "bind_username"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
		},
	})
}

func TestAccLDAPConfiguration_withVerify_CACertificateComplete(t *testing.T) {
	var (
		resourceVerifyName = "mongodbatlas_ldap_verify.test"
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName()
		clusterName        = acc.RandomClusterName()
		hostname           = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username           = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password           = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port               = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		caCertificate      = os.Getenv("MONGODB_ATLAS_LDAP_CA_CERTIFICATE")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAPCert(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyLDAPConfiguration,
		Steps: []resource.TestStep{
			{
				Config: configWithVerify(projectName, orgID, clusterName, hostname, username, password, caCertificate, cast.ToInt(port), true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "bind_username"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttr(resourceVerifyName, "status", "SUCCESS"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.0.validation_type", "SERVER_SPECIFIED"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.0.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.1.validation_type", "CONNECT"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.1.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.2.validation_type", "AUTHENTICATE"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.2.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.3.validation_type", "AUTHORIZATION_ENABLED"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.3.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.4.validation_type", "PARSE_AUTHZ_QUERY_TEMPLATE"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.4.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.5.validation_type", "QUERY_SERVER"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.5.status", "OK"),
				),
			},
		},
	})
}

func TestAccLDAPConfiguration_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_ldap_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		hostname     = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username     = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password     = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port         = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		authEnabled  = true
		projectName  = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyLDAPConfiguration,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, hostname, username, password, authEnabled, cast.ToInt(port)),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "bind_username"),
					resource.TestCheckResourceAttrSet(resourceName, "authentication_enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasLDAPConfigurationImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bind_password"},
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		_, _, err := acc.Conn().LDAPConfigurations.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ldapConfiguration (%s) does not exist", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasLDAPConfigurationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func configBasic(projectName, orgID, hostname, username, password string, authEnabled bool, port int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}

		resource "mongodbatlas_ldap_configuration" "test" {
			project_id                  =  mongodbatlas_project.test.id
			authentication_enabled      =  %[6]t
			hostname					= "%[3]s"
			port                     	=  %[7]d
			bind_username               = "%[4]s"
			bind_password               = "%[5]s"
		}
		
		data "mongodbatlas_ldap_configuration" "test" {
			project_id = mongodbatlas_ldap_configuration.test.id
		}
	`, projectName, orgID, hostname, username, password, authEnabled, port)
}

func configWithVerify(projectName, orgID, clusterName, hostname, username, password, caCertificate string, port int, authEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = "%[3]s"
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			cloud_backup                = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_ldap_verify" "test" {
			project_id                  = mongodbatlas_project.test.id
			hostname = "%[4]s"
			port                     = %[7]d
			bind_username                     = "%[5]s"
			bind_password                     = "%[6]s"
			ca_certificate = <<-EOF
%[9]s
			EOF
			authz_query_template = "{USER}?memberOf?base"
			depends_on = [mongodbatlas_cluster.test]
		}

		resource "mongodbatlas_ldap_configuration" "test" {
			project_id                  = mongodbatlas_project.test.id
			authentication_enabled                = %[8]t
			authorization_enabled                = false
			hostname = "%[4]s"
			port                     = %[7]d
			bind_username                     = "%[5]s"
			bind_password                     = "%[6]s"
			ca_certificate = <<-EOF
%[9]s
			EOF
			authz_query_template = "{USER}?memberOf?base"
			user_to_dn_mapping{
				match = "(.+)"
				ldap_query = "DC=example,DC=com??sub?(userPrincipalName={0})"
			}
			depends_on = [mongodbatlas_ldap_verify.test]
		}`, projectName, orgID, clusterName, hostname, username, password, port, authEnabled, caCertificate)
}
