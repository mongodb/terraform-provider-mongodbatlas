package atlasuser_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	admin20241113 "go.mongodb.org/atlas-sdk/v20241113005/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSAtlasUser_ByUserID(t *testing.T) {
	acc.SkipInUnitTest(t) // needed while fetchUser is called from the test
	var (
		dataSourceName = "data.mongodbatlas_atlas_user.test"
		userID         = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		user           = fetchUser(t, userID)
	)
	resource.Test(t, resource.TestCase{ // does not run in parallel to avoid changes in fetched user during execution
		PreCheck:                 func() { acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUserByUserID(userID),
				Check: resource.ComposeAggregateTestCheckFunc(
					dataSourceChecksForUser(dataSourceName, "", user)...,
				),
			},
		},
	})
}

func TestAccConfigDSAtlasUser_ByUsername(t *testing.T) {
	acc.SkipInUnitTest(t) // needed while fetchUserByUsername is called from the test
	var (
		dataSourceName = "data.mongodbatlas_atlas_user.test"
		username       = os.Getenv("MONGODB_ATLAS_USERNAME")
		user           = fetchUserByUsername(t, username)
	)
	resource.Test(t, resource.TestCase{ // does not run in parallel to avoid changes in fetched user during execution
		PreCheck:                 func() { acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUserByUsername(username),
				Check: resource.ComposeAggregateTestCheckFunc(
					dataSourceChecksForUser(dataSourceName, "", user)...,
				),
			},
		},
	})
}

func dataSourceChecksForUser(dataSourceName, attrPrefix string, user *admin20241113.CloudAppUser) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%susername", attrPrefix), user.Username),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%suser_id", attrPrefix), *user.Id),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%semail_address", attrPrefix), user.EmailAddress),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%sfirst_name", attrPrefix), user.FirstName),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%slast_name", attrPrefix), user.LastName),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%smobile_number", attrPrefix), user.MobileNumber),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%scountry", attrPrefix), user.Country),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%screated_at", attrPrefix), *conversion.TimePtrToStringPtr(user.CreatedAt)),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%steam_ids.#", attrPrefix), fmt.Sprintf("%d", len(*user.TeamIds))),
		resource.TestCheckResourceAttr(dataSourceName, fmt.Sprintf("%slinks.#", attrPrefix), fmt.Sprintf("%d", len(*user.Links))),
		// for assertion of roles the values of `user.Roles` must not be used as it has the risk of flaky executions. CLOUDP-220377
		resource.TestCheckResourceAttrWith(dataSourceName, fmt.Sprintf("%sroles.#", attrPrefix), acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrSet(dataSourceName, fmt.Sprintf("%sroles.0.role_name", attrPrefix)),
	}
}

func TestAccConfigDSAtlasUser_InvalidAttrCombination(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDSMongoDBAtlasUserInvalidAttr(),
				ExpectError: regexp.MustCompile(`Attribute "username" cannot be specified when "user_id" is specified`),
			},
		},
	})
}

func fetchUser(t *testing.T, userID string) *admin20241113.CloudAppUser {
	t.Helper()
	userResp, _, err := acc.ConnV220241113().MongoDBCloudUsersApi.GetUser(context.Background(), userID).Execute()
	if err != nil {
		t.Fatalf("the Atlas User (%s) could not be fetched: %v", userID, err)
	}
	return userResp
}

func fetchUserByUsername(t *testing.T, username string) *admin20241113.CloudAppUser {
	t.Helper()
	userResp, _, err := acc.ConnV220241113().MongoDBCloudUsersApi.GetUserByUsername(context.Background(), username).Execute()
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
