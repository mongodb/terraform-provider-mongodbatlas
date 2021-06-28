package mongodbatlas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceMongoDBAtlasSearchIndexes_basic(t *testing.T) {
	var (
		clusterName    = acctest.RandomWithPrefix("test-acc-global")
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		collectionName = "collection_test"
		databaseName   = "database_test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexesDSConfig(projectID, clusterName, databaseName, collectionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test.0", "name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test.0", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_index.test.0", "cluster_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasSearchIndexesDSConfig(projectID string, clusterName string, databaseName string, collectionName string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_search_indexes" "test" {
			cluster_name           = mongodbatlas_search_index.test.cluster_name
			project_id         = mongodbatlas_search_index.test.project_id
			database_name   = "%s"
			collection_name = "%s"
			page_num = 1
			items_per_page = 100
			
		}
	`, testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName), databaseName, collectionName)
}
