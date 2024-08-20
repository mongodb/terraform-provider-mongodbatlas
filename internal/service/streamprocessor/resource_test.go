package streamprocessor_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName = "mongodbatlas_stream_processor.processor"
)

func TestAccStreamProcessorRS_basic(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		processorName = "new-processor"
		instanceName  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config: streamProcessorConfigWithSampleConnection(projectID, instanceName, processorName, ""),
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
				Config: streamProcessorConfigWithSampleConnection(projectID, instanceName, processorName, "STARTED"),
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
				Config:                  streamProcessorConfigWithSampleConnection(projectID, instanceName, processorName, ""),
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stats"},
			},
		}})
}

func TestAccStreamProcessorRS_createWithAutoStart(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		processorName = "new-processor"
		instanceName  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config: streamProcessorConfigWithSampleConnection(projectID, instanceName, processorName, "STARTED"),
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
				Config: streamProcessorConfigWithSampleConnection(projectID, instanceName, processorName, "STOPPED"),
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

func TestAccStreamProcessorRS_clusterType(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		processorName          = "new-processor"
		instanceName           = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config: streamProcessorConfigWithClusterConnection(projectID, clusterName, instanceName, processorName, "STARTED"),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "processor_name"),
					resource.TestCheckResourceAttr(resourceName, "processor_name", processorName),
					resource.TestCheckResourceAttr(resourceName, "state", "STARTED"),
					resource.TestCheckResourceAttrSet(resourceName, "stats"),
				),
			},
		}})
}

func TestAccStreamProcessorRS_failWithInvalidStateOnCreation(t *testing.T) {
	var (
		projectID     = acc.ProjectIDExecution(t)
		processorName = "new-processor"
		instanceName  = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:      streamProcessorConfigWithSampleConnection(projectID, instanceName, processorName, "STOPPED"),
				ExpectError: regexp.MustCompile("When creating a stream processor, the only valid states are CREATED and STARTED"),
			},
		}})
}

func streamProcessorConfigWithSampleConnection(projectID, instanceName, processorName, state string) string {
	stateConfig := ""
	if state != "" {
		stateConfig = fmt.Sprintf(`state = %[1]q`, state)
	}
	return fmt.Sprintf(`
	resource "mongodbatlas_stream_instance" "instance" {
		project_id    = %[1]q
		instance_name = %[2]q
		data_process_region = {
			region         = "VIRGINIA_USA"
			cloud_provider = "AWS"
		}
	}

	resource "mongodbatlas_stream_connection" "sample" {
		project_id      = %[1]q
		instance_name   = mongodbatlas_stream_instance.instance.instance_name
		connection_name = "sample_stream_solar"
		type            = "Sample"
		depends_on = [mongodbatlas_stream_instance.instance] 
	}

	resource "mongodbatlas_stream_processor" "processor" {
		project_id     = %[1]q
		instance_name  = mongodbatlas_stream_instance.instance.instance_name
		processor_name = %[3]q
		pipeline       = "[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"
		%[4]s
		depends_on = [mongodbatlas_stream_connection.sample] 

	}
	`, projectID, instanceName, processorName, stateConfig)
}

func streamProcessorConfigWithClusterConnection(projectID, clusterName, instanceName, processorName, state string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_stream_instance" "instance" {
		project_id    = %[1]q
		instance_name = %[2]q
		data_process_region = {
			region         = "VIRGINIA_USA"
			cloud_provider = "AWS"
		}
	}

	resource "mongodbatlas_stream_connection" "cluster" {
		project_id      = %[1]q
		instance_name   = mongodbatlas_stream_instance.instance.instance_name
		connection_name = "ClusterConnection"
		type            = "Cluster"
		cluster_name    = %[3]q
		db_role_to_execute = {
			role = "atlasAdmin"
			type = "BUILT_IN"
		}
		depends_on = [mongodbatlas_stream_instance.instance] 
	}

	resource "mongodbatlas_stream_processor" "processor" {
		project_id     = %[1]q
		instance_name  = mongodbatlas_stream_instance.instance.instance_name
		processor_name = %[4]q
		pipeline       = "[{\"$source\":{\"connectionName\":\"ClusterConnection\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"
		state 		   = %[5]q
		depends_on = [mongodbatlas_stream_connection.cluster] 
	}
	`, projectID, instanceName, clusterName, processorName, state)
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
		projectID := rs.Primary.Attributes["project_id"]
		instanceName := rs.Primary.Attributes["instance_name"]
		processorName := rs.Primary.Attributes["processor_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamProcessor(context.Background(), projectID, instanceName, processorName).Execute()
		if err == nil {
			return fmt.Errorf("Stream processor (%s) still exists", processorName)
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
