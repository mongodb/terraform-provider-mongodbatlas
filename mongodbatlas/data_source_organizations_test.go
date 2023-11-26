package mongodbatlas_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSOrganizations_basic(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_organizations.test"
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
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
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
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
