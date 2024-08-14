package streamprocessor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var dataSourceName = "data.mongodbatlas_streamprocessor.test"
var pluralDataSourceName = "data.mongodbatlas_streamprocessors.test"

func TestAccStreamProcessorDS_readManual(t *testing.T) {
	acc.SkipTestForCI(t) // only for manual testing so far, will be moved to resource acceptance test
	var (
		projectID     = acc.ProjectIDExecution(t)
		instanceName  = os.Getenv("MONGODB_ATLAS_STREAM_INSTANCE_NAME")
		processorName = os.Getenv("MONGODB_ATLAS_STREAM_PROCESSOR_NAME")
		checks        = acc.AddAttrChecks(dataSourceName, nil, map[string]string{
			"project_id":     projectID,
			"instance_name":  instanceName,
			"processor_name": processorName,
			"state":          "CREATED",
		})
		checksPlural = acc.AddAttrChecks(pluralDataSourceName, nil, map[string]string{
			"project_id":               projectID,
			"instance_name":            instanceName,
			"results.#":                "1",
			"results.0.processor_name": processorName,
			"results.0.state":          "CREATED",
			"results.0.instance_name":  instanceName,
		})
	)

	checks = acc.AddAttrSetChecks(dataSourceName, checks, "id", "pipeline", "stats")
	checksPlural = acc.AddAttrSetChecks(pluralDataSourceName, checksPlural, "results.0.id", "results.0.pipeline", "results.0.stats")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: streamProcessorConfigDS(projectID, instanceName, processorName),
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
			{
				Config: streamProcessorsConfigDS(projectID, instanceName),
				Check:  resource.ComposeAggregateTestCheckFunc(checksPlural...),
			},
		},
	},
	)
}

func streamProcessorsConfigDS(projectID, instanceName string) string {
	return fmt.Sprintf(`
	data "mongodbatlas_streamprocessors" "test" {
		project_id = %[1]q
		instance_name = %[2]q
	}`, projectID, instanceName)
}

func streamProcessorConfigDS(projectID, instanceName, processorName string) string {
	return fmt.Sprintf(`
	data "mongodbatlas_streamprocessor" "test" {
		project_id = %[1]q
		instance_name = %[2]q
		processor_name = %[3]q
	}`, projectID, instanceName, processorName)
}
