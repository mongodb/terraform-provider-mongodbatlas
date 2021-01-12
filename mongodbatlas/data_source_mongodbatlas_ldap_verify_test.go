package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasLDAPVerify_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		ldapVerify   matlas.LDAPConfiguration
		resourceName = "mongodbatlas_ldap_verify.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = acctest.RandomWithPrefix("test-acc")
		hostname     = os.Getenv("MONGODB_ATLAS_LDAP_HOSTNAME")
		username     = os.Getenv("MONGODB_ATLAS_LDAP_USERNAME")
		password     = os.Getenv("MONGODB_ATLAS_LDAP_PASSWORD")
		port         = os.Getenv("MONGODB_ATLAS_LDAP_PORT")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); checkLDAP(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasLDAPConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceLDAPVerifyConfig(projectName, orgID, clusterName, hostname, username, password, cast.ToInt(port)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasLDAPVerifyExists(resourceName, &ldapVerify),

					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "hostname"),
					resource.TestCheckResourceAttrSet(resourceName, "bind_username"),
					resource.TestCheckResourceAttrSet(resourceName, "request_id"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceLDAPVerifyConfig(projectName, orgID, clusterName, hostname, username, password string, port int) string {
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
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_ldap_verify" "test" {
			project_id               =  mongodbatlas_project.test.id
			hostname 				 = "%[4]s"
			port                     =  %[7]d
			bind_username            = "%[5]s"
			bind_password            = "%[6]s"
			depends_on = ["mongodbatlas_cluster.test"]
		}
		
		data "mongodbatlas_ldap_verify" "test" {
			project_id = mongodbatlas_ldap_verify.test.project_id
			request_id = mongodbatlas_ldap_verify.test.request_id
		}
`, projectName, orgID, clusterName, hostname, username, password, port)
}
