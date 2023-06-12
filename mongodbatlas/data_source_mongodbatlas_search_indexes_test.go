package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConfigDSSearchIndexes_basic(t *testing.T) {
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

func TestAccConfigDSSearchIndexes_WithSynonyms(t *testing.T) {
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
				Config: testAccMongoDBAtlasSearchIndexesDSConfigSynonyms(projectID, clusterName, databaseName, collectionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "database"),
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.index_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "results.0.analyzer"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.0.analyzer", "lucene.simple"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.0.name", "synonym_test"),
					resource.TestCheckResourceAttr(datasourceName, "results.0.synonyms.0.source_collection", "collection_test"),
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
			analyzer = "lucene.simple"
			
		}
	`, testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName), databaseName, collectionName)
}

func testAccMongoDBAtlasSearchIndexesDSConfigSynonyms(projectID, clusterName, databaseName, collectionName string) string {
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
	`, testAccMongoDBAtlasSearchIndexConfigSynonyms(projectID, clusterName), databaseName, collectionName)
}
