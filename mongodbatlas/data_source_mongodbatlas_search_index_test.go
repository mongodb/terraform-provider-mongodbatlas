package mongodbatlas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceMongoDBAtlasSearchIndex_byID(t *testing.T) {
	var (
		clusterName = acctest.RandomWithPrefix("test-acc-global")
		projectID   = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexDSConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test", "collection_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test", "database_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test", "search_analyzer"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasSearchIndexDSConfig(projectID string, clusterName string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_search_index" "test_two" {
			cluster_name        = mongodbatlas_search_index.test.cluster_name
			project_id          = mongodbatlas_search_index.test.project_id
			index_id 			= mongodbatlas_search_index.test.index_id
		}
	`, testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName))
}
