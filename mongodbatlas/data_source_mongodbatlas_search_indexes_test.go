package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasSearchIndexes_basic(t *testing.T) {
	var (
		clusterName    = acctest.RandomWithPrefix("test-acc-global")
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		collectionName = "collection_test"
		databaseName   = "database_test"
		datasourceName = "data.mongodbatlas_search_indexes.data_index"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexesDSConfig(projectID, clusterName, databaseName, collectionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "database"),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.index_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.analyzer"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasSearchIndexesDSConfig(projectID, clusterName, databaseName, collectionName string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_search_indexes" "data_index" {
			cluster_name           = mongodbatlas_search_index.test.cluster_name
			project_id         = mongodbatlas_search_index.test.project_id
			database   = "%s"
			collection_name = "%s"
			page_num = 1
			items_per_page = 100
			
		}
	`, testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName), databaseName, collectionName)
}
