package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccConfigRSProjectAPIKey_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		description  = fmt.Sprintf("test-acc-project-api_key-%s", acctest.RandString(5))
		roleName     = "GROUP_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(projectID, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func TestAccConfigRSProjectAPIKey_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_project_api_key.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		description  = fmt.Sprintf("test-acc-import-project-api_key-%s", acctest.RandString(5))
		roleName     = "GROUP_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectAPIKeyConfigBasic(projectID, description, roleName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectAPIKeyImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectAPIKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_api_key" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		projectAPIKeys, _, err := conn.ProjectAPIKeys.List(context.Background(), ids["project_id"], nil)
		if err != nil {
			return nil
		}

		for _, val := range projectAPIKeys {
			if val.ID == ids["api_key_id"] {
				return fmt.Errorf("Project API Key (%s) still exists", ids["role_name"])
			}
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasProjectAPIKeyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["api_key_id"]), nil
	}
}

func testAccMongoDBAtlasProjectAPIKeyConfigBasic(projectID, description, roleNames string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_api_key" "test" {
			project_id     = %[1]q
			description  = %[2]q
			role_names  = [%[3]q]
		}
	`, projectID, description, roleNames)
}
