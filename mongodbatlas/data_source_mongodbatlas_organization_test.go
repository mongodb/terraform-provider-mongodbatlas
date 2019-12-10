package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasOrganization_basic(t *testing.T) {
	var organization matlas.Organization
	dataSourceName := "data.mongodbatlas_organization.test"

	name := fmt.Sprintf("testacc-organization-%s", acctest.RandString(5))
	nameUpdated := fmt.Sprintf("testacc-organization-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceOrganizationConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists("mongodbatlas_organization.test", &organization),
					testAccCheckMongoDBAtlasOrganizationAttributes(&organization, name),
					resource.TestCheckResourceAttrSet("mongodbatlas_organization.test", "name"),

					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasDataSourceOrganizationConfig(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists("mongodbatlas_organization.test", &organization),
					testAccCheckMongoDBAtlasOrganizationAttributes(&organization, nameUpdated),
					resource.TestCheckResourceAttrSet("mongodbatlas_organization.test", "name"),

					resource.TestCheckResourceAttr(dataSourceName, "name", nameUpdated),
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceOrganizationConfig(name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_organization" "test" {
			name = "%s"
		}

		data "mongodbatlas_organization" "test" {
			org_id = "${mongodbatlas_organization.test.id}"
		}
	`, name)
}
