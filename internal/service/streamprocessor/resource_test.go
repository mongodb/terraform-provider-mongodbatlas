package streamprocessor_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName = "mongodbatlas_stream_processor.processor"
)

func TestAccStreamProcessorRS_basic(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		processorName = "new-processor"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config: streamProcessorConfigWithSampleConnection(projectID, processorName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "processor_name"),
					resource.TestCheckResourceAttr(resourceName, "processor_name", processorName),
					resource.TestCheckResourceAttr(resourceName, "state", "CREATED"),
				),
			},
			{
				Config: streamProcessorConfigWithSampleConnection(projectID, processorName, "STARTED"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "processor_name"),
					resource.TestCheckResourceAttr(resourceName, "processor_name", processorName),
					resource.TestCheckResourceAttr(resourceName, "state", "STARTED"),
				),
			},
			{
				Config:            streamProcessorConfigWithSampleConnection(projectID, processorName, ""),
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		}})
}

func TestAccStreamProcessorRS_createWithAutoStart(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		processorName = "new-processor"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config: streamProcessorConfigWithSampleConnection(projectID, processorName, "STARTED"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "processor_name"),
					resource.TestCheckResourceAttr(resourceName, "processor_name", processorName),
					resource.TestCheckResourceAttr(resourceName, "state", "STARTED"),
				),
			},
			{
				Config: streamProcessorConfigWithSampleConnection(projectID, processorName, "STOPPED"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "processor_name"),
					resource.TestCheckResourceAttr(resourceName, "processor_name", processorName),
					resource.TestCheckResourceAttr(resourceName, "state", "STOPPED"),
				),
			},
		}})
}

func TestAccStreamProcessorRS_failWithInvalidStateOnCreation(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		processorName = "new-processor"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:      streamProcessorConfigWithSampleConnection(projectID, processorName, "STOPPED"),
				ExpectError: regexp.MustCompile("When creating a stream processor, the only valid states are CREATED and STARTED"),
			},
		}})
}

func streamProcessorConfigWithSampleConnection(projectID, processorName, state string) string {
	// Add mongodbatlas_stream_connection once sample stream connection is not created by default
	// resource "mongodbatlas_stream_connection" "sample" {
	// 	project_id      = %[1]q
	// 	instance_name   = mongodbatlas_stream_instance.instance.instance_name
	// 	connection_name = "sample_stream_solar_2"
	// 	type            = "Sample"
	// }
	stateConfig := ""
	if state != "" {
		stateConfig = fmt.Sprintf(`state = %[1]q`, state)
	}
	return fmt.Sprintf(`
	resource "mongodbatlas_stream_instance" "instance" {
		project_id    = %[1]q
		instance_name = "test-instance"
		data_process_region = {
			region         = "VIRGINIA_USA"
			cloud_provider = "AWS"
		}
	}

	resource "mongodbatlas_stream_processor" "processor" {
		project_id     = %[1]q
		instance_name  = mongodbatlas_stream_instance.instance.instance_name
		processor_name = %[2]q
		pipeline       = "[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"
		%[3]s
	}
	`, projectID, processorName, stateConfig)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		instanceName := rs.Primary.Attributes["instance_name"]
		processorName := rs.Primary.Attributes["processor_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamProcessor(context.Background(), projectID, instanceName, processorName).Execute()

		if err != nil {
			return fmt.Errorf("Stream processor (%s) does not exist", processorName)
		}

		return nil
	}
}

func checkDestroyStreamProcessor(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_stream_processor" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().StreamsApi.GetStreamProcessor(context.Background(), ids["project_id"], ids["instance_name"], ids["processor_name"]).Execute()
		if err == nil {
			return fmt.Errorf("Stream processor (%s) still exists", ids["processor_name"])
		}
	}

	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["instance_name"], rs.Primary.Attributes["project_id"], rs.Primary.Attributes["processor_name"]), nil
	}
}
