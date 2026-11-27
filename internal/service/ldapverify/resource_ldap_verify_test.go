package ldapverify_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/ldapverify"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	resourceName   = "mongodbatlas_ldap_verify.test"
	dataSourceName = "data.mongodbatlas_ldap_verify.test"
)

func TestAccLDAPVerify_basic(t *testing.T) {
	resource.Test(t, *basicTestCase(t)) // cannot run in parallel as creating multiple ldap_verify resources for the same project at the same time leads to 500 errors.
}

func TestAccLDAPVerify_withConfiguration_CACertificate(t *testing.T) {
	var (
		hostname      = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username      = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password      = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port          = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		caCertificate = os.Getenv("MONGODB_ATLAS_LDAP_CA_CERTIFICATE")
		projectID, _  = acc.ClusterNameExecution(t, false)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAPCert(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configWithConfiguration(projectID, hostname, username, password, caCertificate, cast.ToInt(port), true),
				Check: resource.ComposeAggregateTestCheckFunc(
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

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		hostname     = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username     = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password     = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port         = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
		projectID, _ = acc.ClusterNameExecution(tb, false) // cluster needed to avoid error: No clusters available to perform LDAP connectivity verification
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckLDAP(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, hostname, username, password, cast.ToInt(port)),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "hostname", hostname),
					resource.TestCheckResourceAttr(resourceName, "bind_username", username),
					resource.TestCheckResourceAttr(resourceName, "port", port),

					checkExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "request_id"),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
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
		_, _, err := acc.ConnV2().LDAPConfigurationAPI.GetUserSecurityVerify(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["request_id"]).Execute()
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

func configBasic(projectID, hostname, username, password string, port int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_ldap_verify" "test" {
			project_id               =  %[1]q
			hostname 				 = %[2]q
			port                     =  %[5]d
			bind_username            = %[3]q
			bind_password            = %[4]q
		}
	
		data "mongodbatlas_ldap_verify" "test" {
			project_id = mongodbatlas_ldap_verify.test.project_id
			request_id = mongodbatlas_ldap_verify.test.request_id
		}		
	`, projectID, hostname, username, password, port)
}

func configWithConfiguration(projectID, hostname, username, password, caCertificate string, port int, authEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_ldap_verify" "test" {
			project_id                  = %[1]q
			hostname = %[2]q
			port                     = %[5]d
			bind_username                     = %[3]q
			bind_password                     = %[4]q
			ca_certificate = <<-EOF
%[7]s
			EOF
		}

		resource "mongodbatlas_ldap_configuration" "test" {
			project_id                  = %[1]q
			authentication_enabled                = %[6]t
			authorization_enabled                = false
			hostname = %[2]q
			port                     = %[5]d
			bind_username                     = %[3]q
			bind_password                     = %[4]q
			ca_certificate = <<-EOF
%[7]s
			EOF
			depends_on = [mongodbatlas_ldap_verify.test]
		}
	`, projectID, hostname, username, password, port, authEnabled, caCertificate)
}

func TestSplitImportID(t *testing.T) {
	const errorString = "import format error: to import a LDAP Verify use the format {project_id}/{request_id} or the legacy format {project_id}-{request_id}"
	tests := map[string]struct {
		importID    string
		projectID   string
		requestID   string
		errorString string
		legacy      bool
	}{
		"slash format": {
			importID:  "664619d870c247237f4b86a6/5cfe95b7e3f0e4b0f4b3d123",
			projectID: "664619d870c247237f4b86a6",
			requestID: "5cfe95b7e3f0e4b0f4b3d123",
		},
		"slash format with more parts": {
			importID:    "part1/part2/part3",
			errorString: errorString,
		},
		"legacy format": {
			importID:  "664619d870c247237f4b86a6-5cfe95b7e3f0e4b0f4b3d123",
			projectID: "664619d870c247237f4b86a6",
			requestID: "5cfe95b7e3f0e4b0f4b3d123",
			legacy:    true,
		},
		"legacy format with hyphens in request ID": {
			importID:  "664619d870c247237f4b86a6-5cfe95b7-3f0e-4b0f-4b3d-123",
			projectID: "664619d870c247237f4b86a6",
			requestID: "5cfe95b7-3f0e-4b0f-4b3d-123",
			legacy:    true,
		},
		"missing request ID": {
			importID:    "664619d870c247237f4b86a6",
			errorString: errorString,
		},
		"empty input": {
			importID:    "",
			errorString: errorString,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			projectID, requestID, legacy, err := ldapverify.SplitImportID(tc.importID)
			if tc.errorString != "" {
				require.EqualError(t, err, tc.errorString)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.projectID, projectID)
			assert.Equal(t, tc.requestID, requestID)
			assert.Equal(t, tc.legacy, legacy)
		})
	}
}
