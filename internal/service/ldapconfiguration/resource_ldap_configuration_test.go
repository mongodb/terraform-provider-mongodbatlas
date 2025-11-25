package ldapconfiguration_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
)

const (
	resourceName   = "mongodbatlas_ldap_configuration.test"
	dataSourceName = "data.mongodbatlas_ldap_configuration.test"
)

func TestAccLDAPConfiguration_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t))
}

func TestAccLDAPConfiguration_withVerify_CACertificateComplete(t *testing.T) {
	var (
		resourceVerifyName = "mongodbatlas_ldap_verify.test"
		hostname           = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username           = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password           = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port               = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		caCertificate      = os.Getenv("MONGODB_ATLAS_LDAP_CA_CERTIFICATE")
		clusterInfo        = acc.GetClusterInfo(t, &acc.ClusterRequest{
			CloudBackup: true,
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{Region: "US_EAST_2"},
			},
		})
		projectID           = clusterInfo.ProjectID
		clusterTerraformStr = clusterInfo.TerraformStr
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAPCert(t); acc.SerialSleep(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithVerify(clusterTerraformStr, clusterInfo.ResourceName, projectID, hostname, username, password, caCertificate, cast.ToInt(port), true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceName, "bind_username", username),
					resource.TestCheckResourceAttr(resourceName, "authentication_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "port", port),
					resource.TestCheckResourceAttr(resourceVerifyName, "status", "SUCCESS"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.#", "5"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.0.validation_type", "CONNECT"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.0.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.1.validation_type", "AUTHENTICATE"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.1.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.2.validation_type", "AUTHORIZATION_ENABLED"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.2.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.3.validation_type", "PARSE_AUTHZ_QUERY_TEMPLATE"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.3.status", "OK"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.4.validation_type", "QUERY_SERVER"),
					resource.TestCheckResourceAttr(resourceVerifyName, "validations.4.status", "OK"),
				),
			},
		},
	})
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		hostname    = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username    = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password    = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port        = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		authEnabled = true
		projectID   = acc.ProjectIDExecution(tb)
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAP(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, hostname, username, password, authEnabled, cast.ToInt(port)),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceName, "bind_username", username),
					resource.TestCheckResourceAttr(resourceName, "authentication_enabled", strconv.FormatBool(authEnabled)),
					resource.TestCheckResourceAttr(resourceName, "port", port),

					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(dataSourceName, "bind_username", username),
					resource.TestCheckResourceAttr(dataSourceName, "authentication_enabled", strconv.FormatBool(authEnabled)),
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
	}
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
		_, _, err := acc.ConnV2().LDAPConfigurationApi.GetUserSecurity(context.Background(), rs.Primary.ID).Execute()
		if err != nil {
			return fmt.Errorf("ldapConfiguration (%s) does not exist", rs.Primary.ID)
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

		return rs.Primary.ID, nil
	}
}

func configBasic(projectID, hostname, username, password string, authEnabled bool, port int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_ldap_configuration" "test" {
			project_id                  =  %[1]q
			hostname								= %[2]q
			bind_username           = %[3]q
			bind_password           = %[4]q
			authentication_enabled  =  %[5]t
			port                   	=  %[6]d
		}
		
		data "mongodbatlas_ldap_configuration" "test" {
			project_id = mongodbatlas_ldap_configuration.test.id
		}
	`, projectID, hostname, username, password, authEnabled, port)
}

func configWithVerify(clusterTerraformStr, clusterResourceName, projectID, hostname, username, password, caCertificate string, port int, authEnabled bool) string {
	return fmt.Sprintf(`
%[8]s

		resource "mongodbatlas_ldap_verify" "test" {
			project_id    = %[1]q
			hostname      = %[2]q
			bind_username = %[3]q
			bind_password = %[4]q
			port          = %[5]d
			ca_certificate = <<-EOF
%[7]s
			EOF
			authz_query_template = "{USER}?memberOf?base"
			depends_on = [%[9]s]
		}

		resource "mongodbatlas_ldap_configuration" "test" {
			project_id             = %[1]q
			authorization_enabled  = false
			hostname               = %[2]q
			bind_username          = %[3]q
			bind_password          = %[4]q
			port                   = %[5]d
			authentication_enabled = %[6]t
			ca_certificate = <<-EOF
%[7]s
			EOF
			authz_query_template = "{USER}?memberOf?base"
			user_to_dn_mapping{
				match = "(.+)"
				ldap_query = "DC=example,DC=com??sub?(userPrincipalName={0})"
			}
			depends_on = [mongodbatlas_ldap_verify.test]
		}`, projectID, hostname, username, password, port, authEnabled, caCertificate, clusterTerraformStr, clusterResourceName)
}
