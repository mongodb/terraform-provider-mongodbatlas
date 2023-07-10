package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccConfigRSAPIKey_Basic(t *testing.T) {
	var (
		resourceName      = "mongodbatlas_api_key.test"
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		description       = fmt.Sprintf("test-acc-api_key-%s", acctest.RandString(5))
		descriptionUpdate = fmt.Sprintf("test-acc-api_key-%s", acctest.RandString(5))
		roleName          = "ORG_MEMBER"
		roleNameUpdated   = "ORG_BILLING_ADMIN"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAPIKeyConfigBasic(orgID, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				Config: testAccMongoDBAtlasAPIKeyConfigBasic(orgID, descriptionUpdate, roleNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdate),
				),
			},
		},
	})
}

func TestAccConfigRSAPIKey_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_api_key.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		description  = fmt.Sprintf("test-acc-import-api_key-%s", acctest.RandString(5))
		roleName     = "ORG_MEMBER"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAPIKeyConfigBasic(orgID, description, roleName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasAPIKeyImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckMongoDBAtlasAPIKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.APIKeys.Get(context.Background(), ids["org_id"], ids["api_key_id"])
		if err != nil {
			return fmt.Errorf("API Key (%s) does not exist", ids["api_key_id"])
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasAPIKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_api_key" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.APIKeys.Get(context.Background(), ids["org_id"], ids["role_name"])
		if err == nil {
			return fmt.Errorf("API Key (%s) still exists", ids["role_name"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasAPIKeyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["org_id"], rs.Primary.Attributes["api_key_id"]), nil
	}
}

func testAccMongoDBAtlasAPIKeyConfigBasic(orgID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_api_key" "test" {
			org_id     = "%s"
			description  = "%s"

			role_names  = ["%s"]
		}
	`, orgID, description, roleNames)
}
