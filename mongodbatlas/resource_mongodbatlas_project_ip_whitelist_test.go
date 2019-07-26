package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasProjectIPWhitelist_basic(t *testing.T) {
	var projectIPEntry matlas.ProjectIPWhitelist

	randInt := acctest.RandIntRange(0, 255)

	resourceName := "mongodbatlas_project_ip_whitelist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("179.154.224.%d/32", randInt)
	ipAddress := fmt.Sprintf("179.154.224.%d", randInt)
	cidrBlockUpdated := fmt.Sprintf("179.154.224.%d/32", randInt)
	ipAddressUpdated := fmt.Sprintf("179.154.224.%d", randInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName, &projectIPEntry),
					testAccCheckMongoDBAtlasProjectIPWhitelistAttributes(&projectIPEntry, ipAddress),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", "for tf acc testing"),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, cidrBlockUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName, &projectIPEntry),
					testAccCheckMongoDBAtlasProjectIPWhitelistAttributes(&projectIPEntry, ipAddressUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlockUpdated),
					resource.TestCheckResourceAttr(resourceName, "comment", "for tf acc testing"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasProjectIPWhitelist_basicBadCIDR(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := "179.154.224.256/32"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, cidrBlock),
				ExpectError: regexp.MustCompile("expected cidr_block to contain a valid CIDR"),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasProjectIPWhitelist_basicInvalid(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := "179.154.224.256/32"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,

		Steps: []resource.TestStep{
			{
				Config:      testAccMongoDBAtlasProjectIPWhitelistConfigInvalid(projectID, cidrBlock),
				ExpectError: regexp.MustCompile(`"cidr_block": conflicts with ip_address`),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasProjectIPWhitelist_importBasic(t *testing.T) {
	randInt := acctest.RandIntRange(0, 255)
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipEntry := fmt.Sprintf("179.154.224.%d/32", randInt)
	importStateID := fmt.Sprintf("%s-%s", projectID, ipEntry)

	resourceName := "mongodbatlas_project_ip_whitelist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, ipEntry),
			},
			{
				ResourceName:      resourceName,
				ImportStateId:     importStateID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName string, ipEntry *matlas.ProjectIPWhitelist) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if resp, _, err := conn.ProjectIPWhitelist.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cidr_block"]); err == nil {
			*ipEntry = *resp
			return nil
		}
		return fmt.Errorf("project ip whitelist entry (%s) does not exist", rs.Primary.Attributes["cidr_block"])
	}
}

func testAccCheckMongoDBAtlasProjectIPWhitelistAttributes(ipEntry *matlas.ProjectIPWhitelist, ipAddress string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if ipEntry.IPAddress != ipAddress {
			return fmt.Errorf("bad ipAddress: %s", ipEntry.IPAddress)
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
		// Try to find the project ip whitelist entry
		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cidr_block"])

		if err == nil {
			return fmt.Errorf("project ip whitelist entry (%s) still exists", rs.Primary.Attributes["cidr_block"])
		}
	}
	return nil
}

func testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, cidrBlock string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id    = "%s"
			cidr_block    = "%s"
			comment = "for tf acc testing"
		}
`, projectID, cidrBlock)
}

func testAccMongoDBAtlasProjectIPWhitelistConfigInvalid(projectID, cidrBlock string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id    = "%s"
			cidr_block    = "%s"
			ip_address    = "0.0.0.0"
			comment = "for tf acc testing"
		}
`, projectID, cidrBlock)
}
