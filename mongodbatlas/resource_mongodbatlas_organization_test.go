package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSOrganization_Basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = fmt.Sprintf("test-acc-organization-%s", acctest.RandString(5))
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasOrganizationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
		},
	})
}

func TestAccConfigRSOrganization_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_organization.test"
		orgOwnerID   = os.Getenv("MONGODB_ATLAS_ORG_OWNER_ID")
		name         = fmt.Sprintf("test-acc-import-organization-%s", acctest.RandString(5))
		description  = "test Key for Acceptance tests"
		roleName     = "ORG_OWNER"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasOrganizationImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccCheckMongoDBAtlasOrganizationExists(resourceName string) resource.TestCheckFunc {
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

		/*config := Config{
			PublicKey:  rs.Primary.Attributes["public_key"],
			PrivateKey: rs.Primary.Attributes["private_key"],
		}

		clients, _ := config.NewClient(context.TODO())
		conn := clients.(*MongoDBClient).Atlas*/

		/*_, _, err := conn.Organizations.Get(context.Background(), ids["org_id"])
		if err != nil {
			return fmt.Errorf("Organization (%s) does not exist", ids["org_id"])
		}*/

		organizationOptions := &matlas.OrganizationsListOptions{}
		orgs, _, err := conn.Organizations.List(context.Background(), organizationOptions)
		if err == nil {
			for _, val := range orgs.Results {
				if val.ID == ids["org_id"] {
					return fmt.Errorf("Organization (%s) still exists", ids["role_name"])
				}
			}
			return nil
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasOrganizationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_organization" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		// config := Config{
		// 	PublicKey:  rs.Primary.Attributes["public_key"],
		// 	PrivateKey: rs.Primary.Attributes["private_key"],
		// }

		// clients, _ := config.NewClient(context.TODO())
		// conn := clients.(*MongoDBClient).Atlas

		organizationOptions := &matlas.OrganizationsListOptions{}
		orgs, _, err := conn.Organizations.List(context.Background(), organizationOptions)
		if err == nil {
			for _, val := range orgs.Results {
				if val.ID == ids["org_id"] {
					return fmt.Errorf("Organization (%s) still exists", ids["role_name"])
				}
			}
			return nil
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasOrganizationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.Attributes["org_id"], nil
	}
}

func testAccMongoDBAtlasOrganizationConfigBasic(orgOwnerID, name, description, roleNames string) string {
	return fmt.Sprintf(`
	  resource "mongodbatlas_organization" "test" {
		org_owner_id = "%s"
		name = "%s"
		description = "%s"
		role_names = ["%s"]
	  }
	`, orgOwnerID, name, description, roleNames)
}
