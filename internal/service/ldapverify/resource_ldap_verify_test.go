package ldapverify_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	resourceName   = "mongodbatlas_ldap_verify.test"
	dataSourceName = "data.mongodbatlas_ldap_verify.test"
)

func TestAccLDAPVerify_basic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		hostname    = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username    = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password    = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port        = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID, clusterName, hostname, username, password, cast.ToInt(port)),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceName, "bind_username", username),
					resource.TestCheckResourceAttr(resourceName, "port", port),

					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "request_id"),
					resource.TestCheckResourceAttr(dataSourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(dataSourceName, "bind_username", username),
					resource.TestCheckResourceAttr(dataSourceName, "port", port),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "bind_password"},
			},
		},
	})
}

func TestAccLDAPVerify_withConfiguration_CACertificate(t *testing.T) {
	var (
		orgID         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		hostname      = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username      = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password      = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port          = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		caCertificate = os.Getenv("MONGODB_ATLAS_LDAP_CA_CERTIFICATE")
		projectName   = acc.RandomProjectName()
		clusterName   = acc.RandomClusterName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAPCert(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithConfiguration(projectName, orgID, clusterName, hostname, username, password, caCertificate, cast.ToInt(port), true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceName, "bind_username", username),
					resource.TestCheckResourceAttr(resourceName, "port", port),
					resource.TestCheckResourceAttr(resourceName, "status", "SUCCESS"),
					resource.TestCheckResourceAttr(resourceName, "validations.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "validations.0.validation_type", "CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "validations.0.status", "OK"),
					resource.TestCheckResourceAttr(resourceName, "validations.1.validation_type", "AUTHENTICATE"),
					resource.TestCheckResourceAttr(resourceName, "validations.1.status", "OK"),
				),
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
		_, _, err := acc.ConnV2().LDAPConfigurationApi.GetLDAPConfigurationStatus(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["request_id"]).Execute()
		if err != nil {
			return fmt.Errorf("ldapVerify (%s) does not exist", rs.Primary.ID)
		}
		return nil
	}
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["request_id"]), nil
	}
}

func configBasic(projectName, orgID, clusterName, hostname, username, password string, port int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[1]q
			org_id = %[2]q
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = %[3]q
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			cloud_backup                = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_ldap_verify" "test" {
			project_id               =  mongodbatlas_project.test.id
			hostname 				 = %[4]q
			port                     =  %[7]d
			bind_username            = %[5]q
			bind_password            = %[6]q
			depends_on = ["mongodbatlas_cluster.test"]
		}
	
		data "mongodbatlas_ldap_verify" "test" {
			project_id = mongodbatlas_ldap_verify.test.project_id
			request_id = mongodbatlas_ldap_verify.test.request_id
		}		
	`, projectName, orgID, clusterName, hostname, username, password, port)
}

func configWithConfiguration(projectName, orgID, clusterName, hostname, username, password, caCertificate string, port int, authEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[1]q
			org_id = %[2]q
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = %[3]q
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			
		}

		resource "mongodbatlas_ldap_verify" "test" {
			project_id                  = mongodbatlas_project.test.id
			hostname = %[4]q
			port                     = %[7]d
			bind_username                     = %[5]q
			bind_password                     = %[6]q
			ca_certificate = <<-EOF
%[9]s
			EOF
			depends_on = [mongodbatlas_cluster.test]
		}

		resource "mongodbatlas_ldap_configuration" "test" {
			project_id                  = mongodbatlas_project.test.id
			authentication_enabled                = %[8]t
			authorization_enabled                = false
			hostname = %[4]q
			port                     = %[7]d
			bind_username                     = %[5]q
			bind_password                     = %[6]q
			ca_certificate = <<-EOF
%[9]s
			EOF
			depends_on = [mongodbatlas_ldap_verify.test]
		}
	`, projectName, orgID, clusterName, hostname, username, password, port, authEnabled, caCertificate)
}
