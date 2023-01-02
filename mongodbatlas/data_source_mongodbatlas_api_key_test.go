package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSAPIKey_basic(t *testing.T) {
	resourceName := "mongodbatlas_api_key.test"
	dataSourceName := "data.mongodbatlas_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	description := fmt.Sprintf("test-acc-api_key-%s", acctest.RandString(5))
	roleName := "ORG_MEMBER"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		//CheckDestroy:      testAccCheckMongoDBAtlasNetworkPeeringDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAPIKeyConfig(orgID, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					testAccCheckMongoDBAtlasAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasAPIKeyConfig(orgID, apiKeyID, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_key" "test" {
		  org_id = "%s"
		  description  = "%s"
		  role_names  = ["%s"]	
		}

		data "mongodbatlas_api_key" "test" {
		  org_id      = "${mongodbatlas_api_key.test.org_id}"
		  api_key_id  = "${mongodbatlas_api_key.test.api_key_id}"
		}
	`, orgID, apiKeyID, roleNames)
}
