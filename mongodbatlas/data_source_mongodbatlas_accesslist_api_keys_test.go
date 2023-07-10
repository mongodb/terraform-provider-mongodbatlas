package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigDSAccesslistAPIKeys_basic(t *testing.T) {
	resourceName := "mongodbatlas_access_list_api_key.test"
	dataSourceName := "data.mongodbatlas_access_list_api_keys.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	description := fmt.Sprintf("test-acc-accesslist-api_keys-%s", acctest.RandString(5))
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAccessListAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAccesslistAPIKeysConfig(orgID, description, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					testAccCheckMongoDBAtlasAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasAccesslistAPIKeysConfig(orgID, description, ipAddress string) string {
	return fmt.Sprintf(`
	data "mongodbatlas_access_list_api_keys" "test" {
		org_id     = %[1]q
		api_key_id = mongodbatlas_access_list_api_key.test.api_key_id
	  }
	  
	  resource "mongodbatlas_api_key" "test" {
		org_id = %[1]q
		description = %[2]q
		role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	  }
	  
	  resource "mongodbatlas_access_list_api_key" "test" {
		org_id     = %[1]q
		ip_address = %[3]q
	    api_key_id = mongodbatlas_api_key.test.api_key_id
	  }
	`, orgID, description, ipAddress)
}
