package organization_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSOrganization_basic(t *testing.T) {
	var (
		orgID = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigWithDS(orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.mongodbatlas_organization.test", "name"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_organization.test", "id"),
				),
			},
		},
	})
}
func testAccMongoDBAtlasOrganizationConfigWithDS(orgID string) string {
	config := fmt.Sprintf(`
		
		data "mongodbatlas_organization" "test" {
			org_id = %[1]q
		}
	`, orgID)
	return config
}
