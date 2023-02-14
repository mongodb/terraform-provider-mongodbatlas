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

func TestAccProjectRSAccesslistAPIKey_SettingIPAddress(t *testing.T) {
	resourceName := "mongodbatlas_access_list_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	description := fmt.Sprintf("test-acc-access_list-api_key-%s", acctest.RandString(5))
	updatedIPAddress := fmt.Sprintf("179.154.228.%d", acctest.RandIntRange(0, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAccessListAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAccessListAPIKeyConfigSettingIPAddress(orgID, description, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
				),
			},
			{
				Config: testAccMongoDBAtlasAccessListAPIKeyConfigSettingIPAddress(orgID, description, updatedIPAddress),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "ip_address", updatedIPAddress),
				),
			},
		},
	})
}

func TestAccProjectRSAccessListAPIKey_SettingCIDRBlock(t *testing.T) {
	resourceName := "mongodbatlas_access_list_api_key.test"
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	cidrBlock := fmt.Sprintf("179.154.226.%d/32", acctest.RandIntRange(0, 255))
	description := fmt.Sprintf("test-acc-access_list-api_key-%s", acctest.RandString(5))
	updatedCIDRBlock := fmt.Sprintf("179.154.228.%d/32", acctest.RandIntRange(0, 255))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAccessListAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAccessListAPIKeyConfigSettingCIDRBlock(orgID, description, cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", cidrBlock),
				),
			},
			{
				Config: testAccMongoDBAtlasAccessListAPIKeyConfigSettingCIDRBlock(orgID, description, updatedCIDRBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAccessListAPIKeyExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "org_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cidr_block"),

					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cidr_block", updatedCIDRBlock),
				),
			},
		},
	})
}

func TestAccProjectRSAccessListAPIKey_importBasic(t *testing.T) {
	orgID := os.Getenv("MONGODB_ATLAS_ORG_ID")
	ipAddress := fmt.Sprintf("179.154.226.%d", acctest.RandIntRange(0, 255))
	resourceName := "mongodbatlas_access_list_api_key.test"
	description := fmt.Sprintf("test-acc-access_list-api_key-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAccessListAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAccessListAPIKeyConfigSettingIPAddress(orgID, description, ipAddress),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasAccessListAPIKeyImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasAccessListAPIKeyExists(resourceName string) resource.TestCheckFunc {
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

func testAccCheckMongoDBAtlasAccessListAPIKeyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_access_list_api_key" {
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

func testAccCheckMongoDBAtlasAccessListAPIKeyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["org_id"], ids["api_key_id"], ids["entry"]), nil
	}
}

func testAccMongoDBAtlasAccessListAPIKeyConfigSettingIPAddress(orgID, description, ipAddress string) string {
	return fmt.Sprintf(`

	   resource "mongodbatlas_api_key" "test" {
		  org_id = %[1]q
		  description = %[2]q
		  role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	    }

		resource "mongodbatlas_access_list_api_key" "test" {
			org_id = %[1]q
			ip_address = %[3]q
			api_key_id = mongodbatlas_api_key.test.api_key_id
		}
	`, orgID, description, ipAddress)
}
func testAccMongoDBAtlasAccessListAPIKeyConfigSettingCIDRBlock(orgID, description, cidrBlock string) string {
	return fmt.Sprintf(`

	resource "mongodbatlas_api_key" "test" {
		org_id = %[1]q
		description = %[2]q
		role_names  = ["ORG_MEMBER","ORG_BILLING_ADMIN"]
	  }

		resource "mongodbatlas_access_list_api_key" "test" {
		  org_id = %[1]q
		  api_key_id = mongodbatlas_api_key.test.api_key_id
		  cidr_block = %[3]q
		}
	`, orgID, description, cidrBlock)
}
