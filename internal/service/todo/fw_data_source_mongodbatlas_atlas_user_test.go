package todo_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc/todoacc"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/util"
	"go.mongodb.org/atlas-sdk/v20231115001/admin"
)

func TestAccConfigDSAtlasUser_ByUserID(t *testing.T) {
	acc.SkipIfTFAccNotDefined(t)
	var (
		dataSourceName = "data.mongodbatlas_atlas_user.test"
		userID         = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		user           = fetchUser(userID, t)
	)
	resource.Test(t, resource.TestCase{ // does not run in parallel to avoid changes in fetched user during execution
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUserByUserID(userID),
				Check: resource.ComposeTestCheckFunc(
					dataSourceChecksForUser(dataSourceName, "", user)...,
				),
			},
		},
	})
}

func TestAccConfigDSAtlasUser_ByUsername(t *testing.T) {
	acc.SkipIfTFAccNotDefined(t)
	var (
		dataSourceName = "data.mongodbatlas_atlas_user.test"
		username       = os.Getenv("MONGODB_ATLAS_USERNAME_CLOUD_DEV")
		user           = fetchUserByUsername(username, t)
	)
	resource.Test(t, resource.TestCase{ // does not run in parallel to avoid changes in fetched user during execution
		PreCheck:                 func() { acc.PreCheckBasic(t); testAccPreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUserByUsername(username),
				Check: resource.ComposeTestCheckFunc(
					dataSourceChecksForUser(dataSourceName, "", user)...,
				),
			},
		},
	})
}

func dataSourceChecksForUser(dataSourceName, attrPrefix string, user *admin.CloudAppUser) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%susername", attrPrefix), user.Username),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%suser_id", attrPrefix), *user.Id),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%semail_address", attrPrefix), user.EmailAddress),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%sfirst_name", attrPrefix), user.FirstName),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%slast_name", attrPrefix), user.LastName),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%smobile_number", attrPrefix), user.MobileNumber),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%scountry", attrPrefix), user.Country),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%screated_at", attrPrefix), *util.TimePtrToStringPtr(user.CreatedAt)),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%sroles.#", attrPrefix), fmt.Sprintf("%d", len(user.Roles))),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%steam_ids.#", attrPrefix), fmt.Sprintf("%d", len(user.TeamIds))),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%slinks.#", attrPrefix), fmt.Sprintf("%d", len(user.Links))),
	}
}

func TestAccConfigDSAtlasUser_InvalidAttrCombination(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: todoacc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDSMongoDBAtlasUserInvalidAttr(),
				ExpectError: regexp.MustCompile(`Attribute "username" cannot be specified when "user_id" is specified`),
			},
		},
	})
}

func fetchUser(userID string, t *testing.T) *admin.CloudAppUser {
	connV2 := todoacc.TestMongoDBClient.(*config.MongoDBClient).AtlasV2
	userResp, _, err := connV2.MongoDBCloudUsersApi.GetUser(context.Background(), userID).Execute()
	if err != nil {
		t.Fatalf("the Atlas User (%s) could not be fetched: %v", userID, err)
	}
	return userResp
}

func fetchUserByUsername(username string, t *testing.T) *admin.CloudAppUser {
	connV2 := testMongoDBClient.(*MongoDBClient).AtlasV2

	userResp, _, err := connV2.MongoDBCloudUsersApi.GetUserByUsername(context.Background(), username).Execute()
	if err != nil {
		t.Fatalf("the Atlas User (%s) could not be fetched: %v", username, err)
	}
	return userResp
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
