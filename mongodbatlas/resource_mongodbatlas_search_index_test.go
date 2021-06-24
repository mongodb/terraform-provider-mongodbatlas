package mongodbatlas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
	"os"
	"strings"
	"testing"
)

func TestAccResourceMongoDBAtlasSearchIndex_basic(t *testing.T) {
	var (
		index        matlas.SearchIndex
		resourceName = "mongodbatlas_search_index.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-global")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	a:= matlas.CustomAnalyzer{}

	b:= matlas.AnalyzerCharFilter{Mappings: }

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,

					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_OWNER"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_DATA_ACCESS_READ_WRITE"},
						},
						{
							TeamID:    teamsIds[2],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
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
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID,

					[]*matlas.ProjectTeam{
						{
							TeamID:    teamsIds[0],
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_READ_ONLY"},
						},
						{
							TeamID:    teamsIds[1],
							RoleNames: []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"},
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
				Config: testAccMongoDBAtlasProjectConfig(projectName, orgID, []*matlas.ProjectTeam{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasProjectExists(resourceName, &project),
					testAccCheckMongoDBAtlasProjectAttributes(&project, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
		},
	})
}

func testAccMongoDBAtlasSearchIndexConfig(projectID string, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_index" "test" {
			project_id         = "%[1]s"
			cluster_name       = "%[2]s"

			analyzer = "lucene.standard"
			collectionName = "collection_test"
			database = "database_test"
			mappings{
				dynamic = true
			}
			name = "name_test"
			searchAnalyzer = "lucene.standard"

		data "mongodbatlas_search_index" "test" {
			cluster_name           = mongodbatlas_search_index.test.cluster_name
			project_id         = mongodbatlas_search_index.test.project_id
		}
	`, projectID, clusterName)
}

func testAccCheckMongoDBAtlasSearchIndexDestroy(state *terraform.State) error {

}
