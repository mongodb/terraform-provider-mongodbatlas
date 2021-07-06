package mongodbatlas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccResourceMongoDBAtlasSearchAnalyzer_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_search_analyzer.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-global")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchAnalyzerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchAnalyzerConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasSearchAnalyzerExists(resourceName),
				),
			},
		},
	})
}
func TestAccResourceMongoDBAtlasSearchAnalyzer_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_search_analyzer.test"
		clusterName  = acctest.RandomWithPrefix("test-acc-global")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchAnalyzerDestroy,
		Steps: []resource.TestStep{
			{
				Config:            testAccMongoDBAtlasSearchAnalyzerConfig(projectID, clusterName),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasSearchAnalyzerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:            testAccMongoDBAtlasSearchAnalyzerConfigAdvanced(projectID, clusterName),
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasSearchAnalyzerImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func testAccCheckMongoDBAtlasSearchAnalyzerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.Search.ListAnalyzers(context.Background(), ids["project_id"], ids["cluster_name"], nil)
		if err != nil {
			return fmt.Errorf("index (%s) does not exist", ids["index_id"])
		}

		return nil
	}
}

func testAccMongoDBAtlasSearchAnalyzerConfig(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_analyzer" "test_analyzer" {
			project_id   = "%[1]s"
			cluster_name = "%[2]s"

			search_analyzers = {
				name = "test_analyzer_1"
				base_analyzer = "lucene.standard"
			},
		}
		data "mongodbatlas_search_analyzer" "test_analyzer" {
			cluster_name = mongodbatlas_search_index.test.cluster_name
			project_id   = mongodbatlas_search_index.test.project_id
		}
	`, projectID, clusterName)
}

func testAccMongoDBAtlasSearchAnalyzerConfigAdvanced(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_search_analyzer" "test_analyzer" {
			project_id   = "%[1]s"
			cluster_name = "%[2]s"

			search_analyzers = {
				name = "test_analyzer_1"
				base_analyzer = "lucene.standard"
				ignore_case = "false"
				stem_exclusion_set = ["foo", "bar", "baz"]
				stopwords = ["foo", "bar", "baz"]
			}

			search_analyzers = {
				name = "test_analyzer_2"
				base_analyzer = "lucenene.standard"
				ignore_case = true
			}
		}
		data "mongodbatlas_search_analyzer" "test_analyzer" {
			cluster_name = mongodbatlas_search_index.test.cluster_name
			project_id   = mongodbatlas_search_index.test.project_id
		}
	`, projectID, clusterName)
}

func testAccCheckMongoDBAtlasSearchAnalyzerDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_search_analyzer" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.Search.ListAnalyzers(context.Background(), ids["project_id"], ids["cluster_name"], nil)
		if err == nil {
			return fmt.Errorf("index id (%s) still exists", ids["index_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasSearchAnalyzerImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.ID, nil
	}
}
