package atlasuser_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/atlasuser"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSAtlasUsers_ByOrgID(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_atlas_users.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	// We can only ensure count > 0 because this test relies on all users in the organization to be stable during its test run.
	// Our test suite is running multiple jobs in parallel, so that guarantee is not fulfilled.
	// We should isolate these runs by using separate organizations.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUsersByOrgID(orgID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("org_id"),
						knownvalue.StringExact(orgID),
					),
					statecheck.ExpectKnownValue(
						dataSourceName,
						tfjsonpath.New("total_count"),
						knownvalue.Int64Func(func(v int64) error {
							if v > 0 {
								return nil
							}
							return fmt.Errorf("expected total_count to be > 0, got %d", v)
						}),
					),
				},
			},
		},
	})
}

func TestAccConfigDSAtlasUsers_ByProjectID(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_atlas_users.test"
		projectName    = acc.RandomProjectName()
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUsersByProjectID(projectName, orgID, projectOwnerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "total_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"), // we know project will only have the project owner
					resource.TestCheckResourceAttr(dataSourceName, "results.0.user_id", projectOwnerID),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.email_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.first_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.last_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.created_at"),
				),
			},
		},
	})
}

func TestAccConfigDSAtlasUsers_ByTeamID(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_atlas_users.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username       = os.Getenv("MONGODB_ATLAS_USERNAME")
		teamName       = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUsersByTeamID(orgID, teamName, username),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "team_id"),
					resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(dataSourceName, "total_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "1"), // we know created team has only 1 user
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.user_id"),
					resource.TestCheckResourceAttr(dataSourceName, "results.0.username", username),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.email_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.first_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.last_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.created_at"),
				),
			},
		},
	})
}

func TestAccConfigDSAtlasUsers_UsingPagination(t *testing.T) {
	var (
		dataSourceName = "data.mongodbatlas_atlas_users.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		username       = os.Getenv("MONGODB_ATLAS_USERNAME")
		teamName       = acc.RandomName()
		pageNum        = 2
		itemsPerPage   = 1
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckAtlasUsername(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyTeam,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasUsersByTeamWithPagination(orgID, teamName, username, itemsPerPage, pageNum),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "team_id"),
					resource.TestCheckResourceAttr(dataSourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(dataSourceName, "total_count", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "items_per_page", fmt.Sprintf("%d", itemsPerPage)),
					resource.TestCheckResourceAttr(dataSourceName, "page_num", fmt.Sprintf("%d", pageNum)),
					resource.TestCheckResourceAttr(dataSourceName, "results.#", "0"),
				),
			},
		},
	})
}

func TestAccConfigDSAtlasUsers_InvalidAttrCombinations(t *testing.T) {
	tests := []struct {
		name          string
		config        string
		expectedError string
	}{
		{
			name: "invalid all three attributes defined",
			config: `
				data "mongodbatlas_atlas_users" "test" {
					org_id = "64c0f3f5ce752426ab9f506b"
					project_id = "64c0f3f5ce752426ab9f506b"
					team_id = "64c0f3f5ce752426ab9f506b"
				}
			`,
			expectedError: "Invalid Attribute Combination",
		},
		{
			name: "invalid org and project attributes defined",
			config: `
				data "mongodbatlas_atlas_users" "test" {
					org_id = "64c0f3f5ce752426ab9f506b"
					project_id = "64c0f3f5ce752426ab9f506b"
				}
			`,
			expectedError: "Invalid Attribute Combination",
		},
		{
			name: "invalid project and team attributes defined",
			config: `
				data "mongodbatlas_atlas_users" "test" {
					team_id = "64c0f3f5ce752426ab9f506b"
					project_id = "64c0f3f5ce752426ab9f506b"
				}
			`,
			expectedError: "Invalid Attribute Combination",
		},
		{
			name: "invalid team attribute defined",
			config: `
				data "mongodbatlas_atlas_users" "test" {
					team_id = "64c0f3f5ce752426ab9f506b"
				}
			`,
			expectedError: atlasuser.ErrorMissingAttributesDetail,
		},
		{
			name: "invalid empty attributes defined",
			config: `
				data "mongodbatlas_atlas_users" "test" {
				}
			`,
			expectedError: atlasuser.ErrorMissingAttributesDetail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.ParallelTest(t, resource.TestCase{
				PreCheck:                 func() { acc.PreCheckBasic(t) },
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Steps: []resource.TestStep{
					{
						Config:      tt.config,
						ExpectError: regexp.MustCompile(tt.expectedError),
					},
				},
			})
		})
	}
}

func testAccDSMongoDBAtlasUsersByOrgID(orgID string) string {
	return fmt.Sprintf(`
		data "mongodbatlas_atlas_users" "test" {
			org_id = %[1]q
		}
	`, orgID)
}

func testAccDSMongoDBAtlasUsersByProjectID(projectName, orgID, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   			 = %[1]q
			org_id 			 = %[2]q
			project_owner_id = %[3]q
		}

		data "mongodbatlas_atlas_users" "test" {
			project_id = mongodbatlas_project.test.id
		}
	`, projectName, orgID, projectOwnerID)
}

func testAccDSMongoDBAtlasUsersByTeamID(orgID, teamName, username string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "test" {
			org_id     = %[1]q
			name       = %[2]q
			usernames  = [%[3]q]
		}
		
		data "mongodbatlas_atlas_users" "test" {
			org_id = %[1]q
			team_id = mongodbatlas_team.test.team_id
		}
	`, orgID, teamName, username)
}

func testAccDSMongoDBAtlasUsersByTeamWithPagination(orgID, teamName, username string, itemsPerPage, pageNum int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_team" "test" {
			org_id     = %[1]q
			name       = %[2]q
			usernames  = [%[3]q]
		}
		
		data "mongodbatlas_atlas_users" "test" {
			org_id = %[1]q
			team_id = mongodbatlas_team.test.team_id
			items_per_page = %[4]d
			page_num = %[5]d
		}
	`, orgID, teamName, username, itemsPerPage, pageNum)
}
