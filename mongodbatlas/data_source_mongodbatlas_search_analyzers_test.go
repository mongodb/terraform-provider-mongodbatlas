package mongodbatlas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceMongoDBAtlasSearchAnalyzer_basic(t *testing.T) {
	var (
		clusterName = acctest.RandomWithPrefix("test-acc-global")
		projectID   = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasSearchAnalyzerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasSearchAnalyzerDSConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_search_analyzers.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_search_analyzers.test", "cluster_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasSearchAnalyzerDSConfig(projectID, clusterName string) string {
	return fmt.Sprintf(`
		%s

		data "mongodbatlas_search_analyzers" "test" {
			cluster_name           = mongodbatlas_search_analyzer.test.cluster_name
			project_id         = mongodbatlas_search_analyzer.test.project_id
		}
	`, testAccMongoDBAtlasSearchAnalyzerConfig(projectID, clusterName))
}
