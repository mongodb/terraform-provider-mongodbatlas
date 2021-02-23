package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasLDAPConfiguration_basic(t *testing.T) {
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
				Config: testAccMongoDBAtlasDataSourceLDAPConfigurationConfig(projectName, orgID, hostname, username, password, authEnabled, cast.ToInt(port)),
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

func testAccMongoDBAtlasDataSourceLDAPConfigurationConfig(projectName, orgID, hostname, username, password string, authEnabled bool, port int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = "%[1]s"
			org_id = "%[2]s"
		}

		resource "mongodbatlas_ldap_configuration" "test" {
			project_id                  =  mongodbatlas_project.test.id
			authentication_enabled      =  %[6]t
			hostname 					= "%[3]s"
			port                        =  %[7]d
			bind_username               = "%[4]s"
			bind_password               = "%[5]s"
		}
		
		data "mongodbatlas_ldap_configuration" "test" {
			project_id = mongodbatlas_ldap_configuration.test.id
		}
`, projectName, orgID, hostname, username, password, authEnabled, port)
}
