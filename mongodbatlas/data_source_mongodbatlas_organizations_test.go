package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSOrganizations_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationsConfigWithDS(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_organizations.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_organizations.test", "id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccConfigDSOrganizations_withPagination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationsConfigWithPagination(2, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_organizations.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_organizations.test", "id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_organizations.test", "results.#"),
				),
				ExpectNonEmptyPlan: true,
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
