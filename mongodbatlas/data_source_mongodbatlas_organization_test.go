package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSOrganization_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); checkTeamsIds(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigWithDS(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_organization.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_organization.test", "id"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func testAccMongoDBAtlasOrganizationConfigWithDS(IncludeDeletedOrgs bool) string {
	config := fmt.Sprintf(`
		
		data "mongodbatlas_organization" "test" {
			include_deleted_orgs = %t
		}
	`, IncludeDeletedOrgs)
	return config
}
