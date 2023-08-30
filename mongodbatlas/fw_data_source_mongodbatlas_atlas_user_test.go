package mongodbatlas

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigDSAtlasUser_ByUserID(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_atlas_user.test"
		userID         = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testAccPreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUserByUserID(userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "username"),
					resource.TestCheckResourceAttr(dataSourceName, "user_id", userID),
					resource.TestCheckResourceAttrSet(dataSourceName, "email_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "first_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "country"),
				),
			},
		},
	})
}

func TestAccConfigDSAtlasUser_ByUsername(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_atlas_user.test"
		username       = os.Getenv("MONGODB_ATLAS_USERNAME_CLOUD_DEV")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t); testAccPreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUserByUsername(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "username", username),
					resource.TestCheckResourceAttrSet(dataSourceName, "user_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "email_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "first_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "last_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "country"),
				),
			},
		},
	})
}

func TestAccConfigDSAtlasUser_InvalidAttrCombination(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckBasic(t) },
		ProtoV6ProviderFactories: testAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDSMongoDBAtlasUserInvalidAttr(),
				ExpectError: regexp.MustCompile(`Attribute "username" cannot be specified when "user_id" is specified`),
			},
		},
	})
}

func testAccDSMongoDBAtlasUserByUsername(username string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_atlas_user" "test" {
			username = %[1]q
		}
	`, username)
}

func testAccDSMongoDBAtlasUserByUserID(userID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_atlas_user" "test" {
			user_id = %[1]q
		}
	`, userID)
}

func testAccDSMongoDBAtlasUserInvalidAttr() string {
	return `
		data "mongodbatlas_atlas_user" "test" {
			user_id = "64b6b9e71f89ae5ca0f0c866"
			username = "some@gmail.com"
		}
	`
}
