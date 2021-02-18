package mongodbatlas

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasLDAPConfiguration_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		ldapConfiguration matlas.LDAPConfiguration
		resourceName      = "mongodbatlas_ldap_configuration.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName       = acctest.RandomWithPrefix("test-acc")
		hostname          = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username          = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password          = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		authEnabled       = true
		port              = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkLDAP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasLDAPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasLDAPConfigurationConfig(projectName, orgID, hostname, username, password, authEnabled, cast.ToInt(port)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasLDAPConfigurationExists(resourceName, &ldapConfiguration),

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

func TestAccResourceMongoDBAtlasLDAPConfigurationWithVerify_CACertificateComplete(t *testing.T) {
	SkipTestExtCred(t)
	var (
		ldapConfiguration     matlas.LDAPConfiguration
		resourceName          = "mongodbatlas_ldap_configuration.test"
		resourceVerifyName    = "mongodbatlas_ldap_verify.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName           = acctest.RandomWithPrefix("test-acc")
		clusterName           = acctest.RandomWithPrefix("test-acc")
		hostname              = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username              = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password              = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port                  = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		caCertificateFilePath = os.Getenv("MONGODB_ATLAS_LDAP_CA_FILEPATH")
		caCertificate         string
	)
	fileContents, err := ioutil.ReadFile(caCertificateFilePath)
	if err != nil {
		fmt.Println(err)
	}
	caCertificate = string(fileContents)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkLDAP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasLDAPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasLDAPConfigurationWithVerifyConfig(projectName, orgID, clusterName, hostname, username, password, caCertificate, cast.ToInt(port), true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasLDAPConfigurationExists(resourceName, &ldapConfiguration),

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

func TestAccResourceMongoDBAtlasLDAPConfiguration_importBasic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		ldapConf     = matlas.LDAPConfiguration{}
		resourceName = "mongodbatlas_ldap_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		hostname     = "3.138.245.82"
		username     = "cn=admin,dc=space,dc=intern"
		password     = "neuewelt32"
		authEnabled  = true
		port         = 7001
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkLDAP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasLDAPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasLDAPConfigurationConfig(projectName, orgID, hostname, username, password, authEnabled, port),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasLDAPConfigurationExists(resourceName, &ldapConf),

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
				ImportStateVerifyIgnore: []string{"project_id", "bind_password"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasLDAPConfigurationExists(resourceName string, ldapConf *matlas.LDAPConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ldapConfRes, _, err := conn.LDAPConfigurations.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("ldapConfiguration (%s) does not exist", rs.Primary.ID)
		}

		ldapConf = ldapConfRes

		return nil
	}
}

func testAccCheckMongoDBAtlasLDAPConfigurationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_ldap_configuration" {
			continue
		}

		_, _, err := conn.LDAPConfigurations.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ldapConfiguration (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
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

func testAccMongoDBAtlasLDAPConfigurationConfig(projectName, orgID, hostname, username, password string, authEnabled bool, port int) string {
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
		}`, projectName, orgID, hostname, username, password, authEnabled, port)
}

func testAccMongoDBAtlasLDAPConfigurationWithVerifyConfig(projectName, orgID, clusterName, hostname, username, password, caCertificate string, port int, authEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = "%[3]s"
			disk_size_gb = 5

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true //enable cloud provider snapshots
			provider_disk_iops          = 100
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
