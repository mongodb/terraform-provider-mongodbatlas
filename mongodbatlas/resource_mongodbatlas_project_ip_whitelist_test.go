package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasProjectIPWhitelist_basic(t *testing.T) {
	var projectIPEntry matlas.ProjectIPWhitelist

	resourceName := "mongodbatlas_project_ip_whitelist.test"
	projectID := "5cf5a45a9ccf6400e60981b6" //Modify until project data source is created.
	cidrBlock := "179.154.224.129/32"
	ipAddress := "179.154.224.129"

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
				Config: testAccMongoDBAtlasProjectIPWhitelistConfig(projectID, "179.154.224.128/32"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName, &projectIPEntry),
					testAccCheckMongoDBAtlasProjectIPWhitelistAttributes(&projectIPEntry, "179.154.224.128"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", "179.154.224.128/32"),
					resource.TestCheckResourceAttr(resourceName, "comment", "for tf acc testing"),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasProjectIPWhitelist_basicInvalid(t *testing.T) {
	projectID := "5cf5a45a9ccf6400e60981b6" //Modify until project data source is created.
	cidrBlock := "179.154.224.129/32"

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
	projectID := "5cf5a45a9ccf6400e60981b6"
	ipEntry := "179.154.224.130/32"
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

		log.Printf("[DEBUG] projectID: %s", rs.Primary.Attributes["project_id"])

		if resp, _, err := conn.ProjectIPWhitelist.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.ID); err == nil {
			*ipEntry = *resp
			return nil
		}

		return fmt.Errorf("project ip whitelist entry (%s) does not exist", rs.Primary.ID)
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
		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("project ip whitelist entry (%s) still exists", rs.Primary.ID)
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
