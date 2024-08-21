package streamprocessor_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

type connectionConfig struct {
	connectionType       string
	clusterName          string
	pipelineStepIsSource bool
}

var (
	resourceName      = "mongodbatlas_stream_processor.processor"
	connTypeSample    = "Sample"
	connTypeCluster   = "Cluster"
	connTypeKafka     = "Kafka"
	connTypeTestLog   = "TestLog"
	sampleSrcConfig   = connectionConfig{connectionType: connTypeSample, pipelineStepIsSource: true}
	testLogDestConfig = connectionConfig{connectionType: connTypeTestLog, pipelineStepIsSource: false}
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
				Config: config(t, projectID, instanceName, processorName, "", sampleSrcConfig, testLogDestConfig),
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
				Config: config(t, projectID, instanceName, processorName, streamprocessor.StartedState, sampleSrcConfig, testLogDestConfig),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "processor_name"),
					resource.TestCheckResourceAttr(resourceName, "processor_name", processorName),
					resource.TestCheckResourceAttr(resourceName, "state", streamprocessor.StartedState),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stats"},
			},
		}})
}

func TestAccStreamProcessorRS_withOptions(t *testing.T) {
	var (
		projectID, clusterName = acc.ClusterNameExecution(t)
		processorName          = "new-processor"
		instanceName           = acc.RandomName()
		src                    = connectionConfig{connectionType: connTypeCluster, clusterName: clusterName, pipelineStepIsSource: true}
		dest                   = connectionConfig{connectionType: connTypeKafka, pipelineStepIsSource: false}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config: config(t, projectID, instanceName, processorName, streamprocessor.CreatedState, src, dest),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttr(resourceName, "processor_name", processorName),
					resource.TestCheckResourceAttr(resourceName, "state", "CREATED"),
				),
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
				Config: config(t, projectID, instanceName, processorName, streamprocessor.StartedState, sampleSrcConfig, testLogDestConfig),
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
				Config: config(t, projectID, instanceName, processorName, streamprocessor.StoppedState, sampleSrcConfig, testLogDestConfig),
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
		srcConfig              = connectionConfig{connectionType: connTypeCluster, clusterName: clusterName, pipelineStepIsSource: true}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config: config(t, projectID, instanceName, processorName, streamprocessor.StartedState, srcConfig, testLogDestConfig),
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
				Config:      config(t, projectID, instanceName, processorName, streamprocessor.StoppedState, sampleSrcConfig, testLogDestConfig),
				ExpectError: regexp.MustCompile("When creating a stream processor, the only valid states are CREATED and STARTED"),
			},
		}})
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

func config(t *testing.T, projectID, instanceName, processorName, state string, src, dest connectionConfig) string {
	t.Helper()
	stateConfig := ""
	if state != "" {
		stateConfig = fmt.Sprintf(`state = %[1]q`, state)
	}

	connectionConfigSrc, connectionIDSrc, pipelineStepSrc := configConnection(t, projectID, src)
	connectionConfigDest, connectionIDDest, pipelineStepDest := configConnection(t, projectID, dest)
	dependsOn := []string{}
	if connectionIDSrc != "" {
		dependsOn = append(dependsOn, connectionIDSrc)
	}
	if connectionIDDest != "" {
		dependsOn = append(dependsOn, connectionIDDest)
	}
	dependsOnStr := strings.Join(dependsOn, ", ")
	pipeline := fmt.Sprintf("[{\"$source\":%1s},{\"$emit\":%2s}]", pipelineStepSrc, pipelineStepDest)
	fmt.Println("\nPIPELINE:")
	fmt.Println(pipeline)

	configStr := fmt.Sprintf(`
		resource "mongodbatlas_stream_instance" "instance" {
			project_id    = %[1]q
			instance_name = %[2]q
			data_process_region = {
				region         = "VIRGINIA_USA"
				cloud_provider = "AWS"
			}
		}

		%[3]s
		%[4]s

		resource "mongodbatlas_stream_processor" "processor" {
			project_id     = %[1]q
			instance_name  = mongodbatlas_stream_instance.instance.instance_name
			processor_name = %[5]q
			pipeline       = %[6]q
			%[7]s
			depends_on = [%[8]s]
		}
	`, projectID, instanceName, connectionConfigSrc, connectionConfigDest, processorName, pipeline, stateConfig, dependsOnStr)
	fmt.Println("\nCONFIG:")
	fmt.Println(configStr)
	return configStr
}

func configConnection(t *testing.T, projectID string, config connectionConfig) (connectionConfig, resourceID, pipelineStep string) {
	t.Helper()
	connectionType := config.connectionType
	pipelineStepIsSource := config.pipelineStepIsSource
	switch connectionType {
	case "Cluster":
		var connectionName, resourceName string
		clusterName := config.clusterName
		assert.NotEqual(t, "", clusterName)
		if pipelineStepIsSource {
			connectionName = "ClusterConnectionSrc"
			resourceName = "cluster_src"
		} else {
			connectionName = "ClusterConnectionDest"
			resourceName = "cluster_dest"
		}
		connectionConfig = fmt.Sprintf(`
            resource "mongodbatlas_stream_connection" %[4]q {
                project_id      = %[1]q
                cluster_name    = %[2]q
                instance_name   = mongodbatlas_stream_instance.instance.instance_name
                connection_name = %[3]q
                type            = "Cluster"
                db_role_to_execute = {
                    role = "atlasAdmin"
                    type = "BUILT_IN"
                }
                depends_on = [mongodbatlas_stream_instance.instance] 
            }
        `, projectID, clusterName, connectionName, resourceName)
		resourceID = fmt.Sprintf("mongodbatlas_stream_connection.%s", resourceName)
		pipelineStep = fmt.Sprintf("{\"connectionName\":%q}", connectionName)
		return connectionConfig, resourceID, pipelineStep
	case "Kafka":
		var connectionName, resourceName, pipelineStep string
		if pipelineStepIsSource {
			connectionName = "KafkaConnectionSrc"
			resourceName = "kafka_src"
			pipelineStep = fmt.Sprintf("{\"connectionName\":%q}", connectionName)
		} else {
			connectionName = "KafkaConnectionDest"
			resourceName = "kafka_dest"
			pipelineStep = fmt.Sprintf("{\"connectionName\":%q,\"topic\":\"random_topic\"}", connectionName)
		}
		connectionConfig = fmt.Sprintf(`
            resource "mongodbatlas_stream_connection" %[3]q{
                project_id      = %[1]q
                instance_name   = mongodbatlas_stream_instance.instance.instance_name
                connection_name = %[2]q
                type            = "Kafka"
                authentication = {
                    mechanism = "PLAIN"
                    username  = "user"
                    password  = "rawpassword"
                }
                bootstrap_servers = "localhost:9092,localhost:9092"
                config = {
                    "auto.offset.reset" : "earliest"
                }
                security = {
                    protocol = "PLAINTEXT"
                }
                depends_on = [mongodbatlas_stream_instance.instance] 
            }
        `, projectID, connectionName, resourceName)
		resourceID = fmt.Sprintf("mongodbatlas_stream_connection.%s", resourceName)
		return connectionConfig, resourceID, pipelineStep
	case "Sample":
		if !pipelineStepIsSource {
			t.Fatal("Sample connection must be used as a source")
		}
		connectionConfig = fmt.Sprintf(`
            resource "mongodbatlas_stream_connection" "sample" {
                project_id      = %[1]q
                instance_name   = mongodbatlas_stream_instance.instance.instance_name
                connection_name = "sample_stream_solar"
                type            = "Sample"
                depends_on = [mongodbatlas_stream_instance.instance] 
            }
        `, projectID)
		resourceID = "mongodbatlas_stream_connection.sample"
		pipelineStep = "{\"connectionName\":\"sample_stream_solar\"}"
		return connectionConfig, resourceID, pipelineStep

	case "TestLog":
		if pipelineStepIsSource {
			t.Fatal("TestLog connection must be used as a destination")
		}
		connectionConfig = ""
		resourceID = ""
		pipelineStep = "{\"connectionName\":\"__testLog\"}"
		return connectionConfig, resourceID, pipelineStep
	}
	t.Fatalf("Unknown connection type: %s", connectionType)
	return connectionConfig, resourceID, pipelineStep
}
