package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSProjectAPIKey_basic(t *testing.T) {
	resourceName := "mongodbatlas_project_api_key.test"
	dataSourceName := "data.mongodbatlas_project_api_key.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	description := fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
	roleName := "GROUP_OWNER"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasProjectAPIKeyConfig(projectID, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasProjectAPIKeyConfig(projectID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
		  project_id = %[1]q
		  description  = %[2]q
		  role_names  = [%[3]q]	
		}

		data "mongodbatlas_project_api_key" "test" {
		  project_id      = %[1]q
		  api_key_id  = "${mongodbatlas_project_api_key.test.api_key_id}"
		}
	`, projectID, description, roleNames)
}
