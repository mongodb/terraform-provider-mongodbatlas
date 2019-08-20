package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasProjectIPWhitelist_basic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var whitelist = []*matlas.ProjectIPWhitelist{
		{
			IPAddress: fmt.Sprintf("179.154.224.%d", acctest.RandIntRange(0, 255)),
		},
		{
			CIDRBlock: fmt.Sprintf("179.154.225.%d/32", acctest.RandIntRange(0, 255)),
		},
		{
			IPAddress: fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255)),
		},
		{
			CIDRBlock: fmt.Sprintf("179.154.227.%d/32", acctest.RandIntRange(0, 255)),
		},
	}

	resourceName := "mongodbatlas_project_ip_whitelist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, whitelist),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName, whitelist),
					testAccCheckMongoDBAtlasProjectIPWhitelistAttributes(whitelist[0].IPAddress, whitelist[0].IPAddress),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasProjectIPWhitelist_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var whitelist = []*matlas.ProjectIPWhitelist{
		{
			IPAddress: fmt.Sprintf("179.154.224.%d", acctest.RandIntRange(0, 255)),
		},
		{
			CIDRBlock: fmt.Sprintf("179.154.225.%d/32", acctest.RandIntRange(0, 255)),
		},
		{
			IPAddress: fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255)),
		},
		{
			CIDRBlock: fmt.Sprintf("179.154.227.%d/32", acctest.RandIntRange(0, 255)),
		},
	}

	resourceName := "mongodbatlas_project_ip_whitelist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, whitelist),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasProjectIPWhitelistImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName string, whitelist []*matlas.ProjectIPWhitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		entries, _ := getProjectIPWhitelist(decodeStateID(rs.Primary.ID), conn)
		if len(entries) > 0 {
			whitelist = entries
			return nil
		}
		return fmt.Errorf("project ip whitelist entry (%s) does not exist", rs.Primary.Attributes["project_id"])
	}
}

func testAccCheckMongoDBAtlasProjectIPWhitelistAttributes(cidrBlockEntry string, cidrBlock string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cidrBlockEntry != cidrBlock {
			return fmt.Errorf("bad cidrBlock: %s", cidrBlockEntry)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasProjectIPWhitelistDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_project_ip_whitelist" {
			continue
		}

		entry, err := getProjectIPWhitelist(decodeStateID(rs.Primary.ID), conn)
		if len(entry) > 0 {
			return fmt.Errorf("project ip whitelist entry (%s) still exists %s", rs.Primary.Attributes["project_id"], err)
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasProjectIPWhitelistImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.ID, nil
	}
}

func testAccMongoDBAtlasProjectIPWhitelistConfig(projectID string, entry []*matlas.ProjectIPWhitelist) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id    = "%s"
			whitelist {
				ip_address  = "%s"
				comment = "ip_address for tf acc testing"
			}
			whitelist {
				cidr_block  = "%s"
				comment = "cidr_block for tf acc testing"
			}
			whitelist {
				ip_address  = "%s"
				comment = "ip_address for tf acc testing"
			}
			whitelist {
				cidr_block  = "%s"
				comment = "ip_address for tf acc testing"
			}
		}
	`, projectID, entry[0].IPAddress, entry[1].CIDRBlock, entry[2].IPAddress, entry[3].CIDRBlock)
}
