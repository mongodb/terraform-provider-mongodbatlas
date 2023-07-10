package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccConfigDSProjectAPIKey_basic(t *testing.T) {
	resourceName := "mongodbatlas_project_api_key.test"
	dataSourceName := "data.mongodbatlas_project_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	projectName := acctest.RandomWithPrefix("test-acc")
	description := fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
	roleName := "GROUP_OWNER"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasProjectAPIKeyConfig(orgID, projectName, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasProjectAPIKeyConfig(orgID, projectName, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_project_api_key" "test" {
		  project_id = mongodbatlas_project.test.id
		  description  = %[3]q
		  role_names  = [%[4]q]	
		}

		data "mongodbatlas_project_api_key" "test" {
		  project_id      = mongodbatlas_project.test.id
		  api_key_id  = "${mongodbatlas_project_api_key.test.api_key_id}"
		}
	`, orgID, projectName, description, roleNames)
}
