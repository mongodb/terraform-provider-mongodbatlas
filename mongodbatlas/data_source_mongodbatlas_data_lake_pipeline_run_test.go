package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBackupDSDataLakePipelineRun_basic(t *testing.T) {
	testCheckDataLakePipelineRun(t)
	var (
		dataSourceName = "data.mongodbatlas_data_lake_pipeline_run.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		pipelineName   = os.Getenv("MONGODB_ATLAS_DATA_LAKE_PIPELINE_NAME")
		runID          = os.Getenv("MONGODB_ATLAS_DATA_LAKE_PIPELINE_RUN_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataLakeDataSourcePipelineRunConfig(projectID, pipelineName, runID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "pipeline_name", pipelineName),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "phase"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pipeline_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "dataset_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataLakeDataSourcePipelineRunConfig(projectID, pipelineName, runID string) string {
	return fmt.Sprintf(`

data "mongodbatlas_data_lake_pipeline_run" "test" {
  project_id           = %[1]q
  pipeline_name        = %[2]q
  pipeline_run_id      = %[3]q
}
	`, projectID, pipelineName, runID)
}
