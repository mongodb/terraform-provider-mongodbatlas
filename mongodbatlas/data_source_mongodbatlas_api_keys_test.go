package mongodbatlas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc/todoacc"
)

func TestAccConfigDSAPIKeys_basic(t *testing.T) {
	resourceName := "mongodbatlas_api_key.test"
	dataSourceName := "data.mongodbatlas_api_keys.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	description := fmt.Sprintf("test-acc-api_key-%s", acctest.RandString(5))
	roleName := "ORG_MEMBER"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		CheckDestroy:             testAccCheckMongoDBAtlasAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAPIKeysConfig(orgID, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					// Test for Resource
					testAccCheckMongoDBAtlasAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", description),

					// Test for Data source
					resource.TestCheckResourceAttrSet(dataSourceName, "org_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasAPIKeysConfig(orgID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_key" "test" {
		  org_id = "%s"
		  description  = "%s"
		  role_names  = ["%s"]
		}

		data "mongodbatlas_api_keys" "test" {
		  org_id = "${mongodbatlas_api_key.test.org_id}"
		}
	`, orgID, description, roleNames)
}
