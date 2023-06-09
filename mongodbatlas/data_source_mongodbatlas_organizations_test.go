package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSOrganizations_basic(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_organizations.test"
	)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationsConfigWithDS(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.id"),
				),
			},
		},
	})
}

func TestAccConfigDSOrganizations_withPagination(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_organizations.test"
	)
	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationsConfigWithPagination(2, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasOrganizationsConfigWithDS(includedeletedorgs bool) string {
	config := fmt.Sprintf(`
		
		data "mongodbatlas_organizations" "test" {
			include_deleted_orgs = %t
		}
	`, includedeletedorgs)
	return config
}

func testAccMongoDBAtlasOrganizationsConfigWithPagination(pageNum, itemPage int) string {
	return fmt.Sprintf(`
		data "mongodbatlas_organizations" "test" {
			page_num = %d
			items_per_page = %d
		}
	`, pageNum, itemPage)
}
