package mongodbatlas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
	"os"
	"testing"
)

func TestAccResourceMongoDBAtlasSearchIndex_basic(t *testing.T) {
	var (
		index        matlas.SearchIndex
		resourceName = "mongodbatlas_search_index.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-global")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name         = "name_test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasSearchIndexExists(resourceName, &index),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "project_id", name),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", name),
				),
			},
		},
	})
}
func TestAccResourceMongoDBAtlasSearchIndex_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_search_index.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-global")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchIndexDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchIndexConfig(projectID, clusterName),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasSearchIndexImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckMongoDBAtlasSearchIndexExists(resourceName string, index *matlas.SearchIndex) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		indexResponse, _, err := conn.Search.GetIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"])
		if err != nil {
			return fmt.Errorf("index (%s) does not exist", ids["index_id"])
		}

		*index = *indexResponse

		return nil
	}
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
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_search_index" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.Search.GetIndex(context.Background(), ids["project_id"], ids["cluster_name"], ids["index_id"])
		if err == nil {
			return fmt.Errorf("index id (%s) still exists", ids["index_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasSearchIndexImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
