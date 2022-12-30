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

func TestAccProjectRSProjectIPAccesslistAPIKey_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_accesslist_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))

	updatedIPAddress := fmt.Sprintf("179.154.228.%d", acctest.RandIntRange(0, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListAPIKeyConfigSettingIPAddress(orgID, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPAccessListAPIKeyConfigSettingIPAddress(orgID, updatedIPAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", updatedIPAddress),
				),
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessListAPIKey_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_accesslist_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))

	updatedCIDRBlock := fmt.Sprintf("179.154.228.%d/32", acctest.RandIntRange(0, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListAPIKeyConfigSettingCIDRBlock(orgID, cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectIPAccessListAPIKeyConfigSettingCIDRBlock(orgID, updatedCIDRBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", updatedCIDRBlock),
				),
			},
		},
	})
}

func TestAccProjectRSProjectIPAccessListAPIKey_importBasic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	resourceName := "mongodbatlas_accesslist_api_key.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasProjectIPAccessListAPIKeyConfigSettingIPAddress(orgID, ipAddress),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyExists(resourceName string) resource.TestCheckFunc {
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

		_, _, err := conn.AccessListAPIKeys.Get(context.Background(), ids["org_id"], ids["api_key_id"], ids["entry"])
		if err != nil {
			return fmt.Errorf("access list API Key (%s) does not exist", ids["api_key_id"])
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_accesslist_api_key" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.AccessListAPIKeys.Get(context.Background(), ids["project_id"], ids["api_key_id"], ids["entry"])
		if err == nil {
			return fmt.Errorf("access list API Key (%s) still exists", ids["api_key_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasProjectIPAccessListAPIKeyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["org_id"], ids["api_key_id"], ids["entry"]), nil
	}
}

func testAccMongoDBAtlasProjectIPAccessListAPIKeyConfigSettingIPAddress(orgID, ipAddress string) string {
	return fmt.Sprintf(`

	   resource "mongodbatlas_api_key" "test" {
		  org_id = "%s"
		  description = "IPAccessList test key"
		  role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	    }

		resource "mongodbatlas_accesslist_api_key" "test" {
			org_id = "%s"
			ip_address = "%s"
			api_key_id = mongodbatlas_api_key.test.api_key_id
		}
	`, orgID, orgID, ipAddress)
}
func testAccMongoDBAtlasProjectIPAccessListAPIKeyConfigSettingCIDRBlock(orgID, cidrBlock string) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_api_key" "test" {
		org_id = "%s"
		description = "IPAccessList test key"
		role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	  }

		resource "mongodbatlas_accesslist_api_key" "test" {
			org_id = "%s"
			api_key_id = mongodbatlas_api_key.test.api_key_id
			cidr_block = "%s"
		}
	`, orgID, orgID, cidrBlock)
}
