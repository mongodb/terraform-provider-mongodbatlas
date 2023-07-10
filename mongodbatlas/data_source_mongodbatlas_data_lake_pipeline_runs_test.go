package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBackupDSDataLakePipelineRuns_basic(t *testing.T) {
	testCheckDataLakePipelineRuns(t)
	var (
		dataSourceName = "data.mongodbatlas_data_lake_pipeline_runs.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		pipelineName   = os.Getenv("MONGODB_ATLAS_DATA_LAKE_PIPELINE_NAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeDataSourcePipelineRunsConfig(projectID, pipelineName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "pipeline_name", pipelineName),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataLakeDataSourcePipelineRunsConfig(projectID, pipelineName string) string {
	return fmt.Sprintf(`

data "mongodbatlas_data_lake_pipeline_runs" "test" {
  project_id           = %[1]q
  pipeline_name        = %[2]q
}
	`, projectID, pipelineName)
}
