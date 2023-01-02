package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSProjectAPIKeys_basic(t *testing.T) {
	resourceName := "mongodbatlas_project_api_key.test"
	dataSourceName := "data.mongodbatlas_project_api_keys.test"
	orgID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	description := fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
	roleName := "GROUP_OWNER"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasProjectAPIKeysConfig(orgID, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					//testAccCheckMongoDBAtlasProjectAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),

					resource.TestCheckResourceAttr(resourceName, "project_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", description),

					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasProjectAPIKeysConfig(projectID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
		  project_id = %[1]q
		  description  = %[2]q
		  role_names  = [%[3]q]
		}

		data "mongodbatlas_project_api_keys" "test" {
		  project_id = %[1]q
		}
	`, projectID, description, roleNames)
}
