package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// func TestProjectDataSource_FrameworkMigration(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck: func() { testAccPreCheckBasic(t) },
// 		Steps: []resource.TestStep{
// 			{
// 				ExternalProviders: map[string]resource.ExternalProvider{
// 					"mongodbatlas": {
// 						VersionConstraint: "1.10.0",
// 						Source:            "mongodb/mongodbatlas",
// 					},
// 				},
// 				Config: `data "mongodbatlas_project" "test" {
// 					name   = "framework-datasource-project"

// 				  }`,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("data.mongodbatlas_project.test", "org_id", "63bec56c014da65b8f73c05e"),
// 				),
// 			},
// 			{
// 				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
// 				Config: `data "mongodbatlas_project" "test" {
// 					name   = "framework-datasource-project"
// 				  }`,
// 				PlanOnly: true,
// 			},
// 		},
// 	})
// }

func TestProjectDataSource_Basic_Migration(t *testing.T) {
	var (
		resourceName = "data.mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("test-acc-migration")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.10.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: fmt.Sprintf(`data "mongodbatlas_project" "test" {
					name   = "%s"
				  }`, projectName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Config: fmt.Sprintf(`data "mongodbatlas_project" "test" {
					name   = "%s"
					org_id = "%s"
				  }`, projectName, orgID),
				PlanOnly: true,
			},
		},
	})
}

func TestProjectDatasource_Migration(t *testing.T) {
	var (
		project      matlas.Project
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("test-acc-migration")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount = "0"
		teamsIds     = strings.Split(os.Getenv("MONGODB_ATLAS_TEAMS_IDS"), ",")
		apiKeysIds   = strings.Split(os.Getenv("MONGODB_ATLAS_API_KEYS_IDS"), ",")
	)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckBasic(t); testCheckTeamsIds(t) },
		CheckDestroy: testAccCheckMongoDBAtlasProjectDestroy,

		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: "1.10.0",
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
					[]*testAPIKey{
						{
							id:    apiKeysIds[0],
							roles: []string{"GROUP_READ_ONLY"},
						},
						{
							id:    apiKeysIds[1],
							roles: []string{"GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				ProtoV6ProviderFactories: testProtoV6ProviderFactories,
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,
					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
					[]*testAPIKey{
						{
							id:    apiKeysIds[0],
							roles: []string{"GROUP_READ_ONLY"},
						},
						{
							id:    apiKeysIds[1],
							roles: []string{"GROUP_OWNER"},
						},
					},
				),
				PlanOnly: true,
			},
		},
	})
}
