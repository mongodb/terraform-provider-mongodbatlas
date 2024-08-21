package datalakepipeline_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccDataLakeRunDSPlural_basic(t *testing.T) {
	acc.SkipTestForCI(t) // needs a data lake pipeline, can be joined to resource test

	var (
		dataSourceName = "data.mongodbatlas_data_lake_pipeline_runs.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		pipelineName   = os.Getenv("MONGODB_ATLAS_DATA_LAKE_PIPELINE_NAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckDataLakePipelineRuns(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configRunDSPlural(projectID, pipelineName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "pipeline_name", pipelineName),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
				),
			},
		},
	})
}

func configRunDSPlural(projectID, pipelineName string) string {
	return fmt.Sprintf(`

data "mongodbatlas_data_lake_pipeline_runs" "test" {
  project_id           = %[1]q
  pipeline_name        = %[2]q
}
	`, projectID, pipelineName)
}
