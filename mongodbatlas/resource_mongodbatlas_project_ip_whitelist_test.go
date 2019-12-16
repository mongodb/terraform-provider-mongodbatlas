package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasProjectIPWhitelist_SettingIPAddress(t *testing.T) {

	resourceName := "mongodbatlas_project_ip_whitelist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipAddress (%s)", ipAddress)

	updatedIPAddress := fmt.Sprintf("179.154.228.%d", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for ipAddress (%s)", updatedIPAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, ipAddress, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, updatedIPAddress, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", updatedIPAddress),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})

}
func TestAccResourceMongoDBAtlasProjectIPWhitelist_SettingCIDRBlock(t *testing.T) {

	resourceName := "mongodbatlas_project_ip_whitelist.test"
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for cidrBlock (%s)", cidrBlock)

	updatedCIDRBlock := fmt.Sprintf("179.154.228.%d/32", acctest.RandIntRange(0, 255))
	updatedComment := fmt.Sprintf("TestAcc for cidrBlock (%s)", updatedCIDRBlock)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingCIDRBlock(projectID, cidrBlock, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingCIDRBlock(projectID, updatedCIDRBlock, updatedComment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),
					resource.TestCheckResourceAttrSet(resourceName, "comment"),

					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", updatedCIDRBlock),
					resource.TestCheckResourceAttr(resourceName, "comment", updatedComment),
				),
			},
		},
	})

}

func TestAccResourceMongoDBAtlasProjectIPWhitelist_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	comment := fmt.Sprintf("TestAcc for ipaddres (%s)", ipAddress)
	resourceName := "mongodbatlas_project_ip_whitelist.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasProjectIPWhitelistDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, ipAddress, comment),
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

func testAccCheckMongoDBAtlasProjectIPWhitelistExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"])
		if err != nil {
			return fmt.Errorf("project ip whitelist entry (%s) does not exist", ids["entry"])
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

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.ProjectIPWhitelist.Get(context.Background(), ids["project_id"], ids["entry"])
		if err == nil {
			return fmt.Errorf("project ip whitelist entry (%s) still exists", ids["entry"])
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
		return rs.Primary.Attributes["project_id"], nil
	}
}

func testAccMongoDBAtlasProjectIPWhitelistConfigSettingIPAddress(projectID, ipAddress, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id = "%s"
			ip_address = "%s"
			comment    = "%s"
		}
	`, projectID, ipAddress, comment)
}

func testAccMongoDBAtlasProjectIPWhitelistConfigSettingCIDRBlock(projectID, cidrBlock, comment string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project_ip_whitelist" "test" {
			project_id = "%s"
			cidr_block = "%s"
			comment    = "%s"
		}
	`, projectID, cidrBlock, comment)
}
