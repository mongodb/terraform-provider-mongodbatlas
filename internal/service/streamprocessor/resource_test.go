package streamprocessor_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

type connectionConfig struct {
	connectionType       string
	clusterName          string
	pipelineStepIsSource bool
	useAsDLQ             bool
	extraWhitespace      bool
	invalidJSON          bool
}

var (
	resourceName         = "mongodbatlas_stream_processor.processor"
	dataSourceName       = "data.mongodbatlas_stream_processor.test"
	pluralDataSourceName = "data.mongodbatlas_stream_processors.test"
	connTypeSample       = "Sample"
	connTypeCluster      = "Cluster"
	connTypeKafka        = "Kafka"
	connTypeTestLog      = "TestLog"
	sampleSrcConfig      = connectionConfig{connectionType: connTypeSample, pipelineStepIsSource: true}
	testLogDestConfig    = connectionConfig{connectionType: connTypeTestLog, pipelineStepIsSource: false}
)

func TestAccStreamProcessor_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		randomSuffix             = acctest.RandString(5)
		processorName            = "new-processor" + randomSuffix
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:            config(t, projectID, workspaceName, processorName, "", randomSuffix, sampleSrcConfig, testLogDestConfig, "", nil),
				Check:             composeStreamProcessorChecks(projectID, workspaceName, processorName, streamprocessor.CreatedState, false, false),
				ConfigStateChecks: pluralConfigStateChecks(processorName, streamprocessor.CreatedState, workspaceName, false, false),
			},
			{
				Config:            config(t, projectID, workspaceName, processorName, streamprocessor.StartedState, randomSuffix, sampleSrcConfig, testLogDestConfig, "", nil),
				Check:             composeStreamProcessorChecks(projectID, workspaceName, processorName, streamprocessor.StartedState, true, false),
				ConfigStateChecks: pluralConfigStateChecks(processorName, streamprocessor.StartedState, workspaceName, true, false),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stats"},
			},
		}}
}

// basicTestCaseMigration is the same as basicTestCase but uses instance_name instead of workspace_name to test compatibility of the deprecated instance_name
func basicTestCaseMigration(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID, instanceName = acc.ProjectIDExecutionWithStreamInstance(t)
		randomSuffix            = acctest.RandString(5)
		processorName           = "new-processor" + randomSuffix
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:            configMigration(t, projectID, instanceName, processorName, "", randomSuffix, sampleSrcConfig, testLogDestConfig, "", nil),
				Check:             composeStreamProcessorChecksMigration(projectID, instanceName, processorName, streamprocessor.CreatedState, false, false),
				ConfigStateChecks: pluralConfigStateChecksMigration(processorName, streamprocessor.CreatedState, instanceName, false, false),
			},
			{
				Config:            configMigration(t, projectID, instanceName, processorName, streamprocessor.StartedState, randomSuffix, sampleSrcConfig, testLogDestConfig, "", nil),
				Check:             composeStreamProcessorChecksMigration(projectID, instanceName, processorName, streamprocessor.StartedState, true, false),
				ConfigStateChecks: pluralConfigStateChecksMigration(processorName, streamprocessor.StartedState, instanceName, true, false),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"stats"},
			},
		}}
}

func TestAccStreamProcessor_JSONWhiteSpaceFormat(t *testing.T) {
	var (
		projectID, workspaceName   = acc.ProjectIDExecutionWithStreamInstance(t)
		randomSuffix               = acctest.RandString(5)
		processorName              = "new-processor-json-unchanged"
		sampleSrcConfigExtraSpaces = connectionConfig{connectionType: connTypeSample, pipelineStepIsSource: true, extraWhitespace: true}
	)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:            config(t, projectID, workspaceName, processorName, streamprocessor.CreatedState, randomSuffix, sampleSrcConfigExtraSpaces, testLogDestConfig, "", nil),
				Check:             composeStreamProcessorChecks(projectID, workspaceName, processorName, streamprocessor.CreatedState, false, false),
				ConfigStateChecks: pluralConfigStateChecks(processorName, streamprocessor.CreatedState, workspaceName, false, false),
			},
		}})
}

func TestAccStreamProcessor_withOptions(t *testing.T) {
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		_, clusterName           = acc.ClusterNameExecution(t, false)
		src                      = connectionConfig{connectionType: connTypeCluster, clusterName: clusterName, pipelineStepIsSource: true, useAsDLQ: true}
		dest                     = connectionConfig{connectionType: connTypeKafka, pipelineStepIsSource: false}
		randomSuffix             = acctest.RandString(5)
		processorName            = "new-processor" + randomSuffix
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:            config(t, projectID, workspaceName, processorName, streamprocessor.CreatedState, randomSuffix, src, dest, "", nil),
				Check:             composeStreamProcessorChecks(projectID, workspaceName, processorName, streamprocessor.CreatedState, false, true),
				ConfigStateChecks: pluralConfigStateChecks(processorName, streamprocessor.CreatedState, workspaceName, false, true),
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

func TestAccStreamProcessor_StateTransitionsUpdates(t *testing.T) {
	transitions := map[string]struct {
		setupState   string // Optional: Initial setup state (needed for STOPPED tests since we can't create in a STOPPED state)
		initialState string // State to create or transition to after setup
		targetState  string // Final state to transition to
		description  string // Description of what the test validates
	}{
		"CreatedToCreated": {
			initialState: CreatedState,
			targetState:  CreatedState,
			description:  "Verifies a processor in CREATED state can be updated while remaining in CREATED state",
		},
		"CreatedToStarted": {
			initialState: CreatedState,
			targetState:  StartedState,
			description:  "Verifies a processor can transition from CREATED to STARTED state",
		},
		"StartedToStopped": {
			initialState: StartedState,
			targetState:  StoppedState,
			description:  "Verifies a processor can transition from STARTED to STOPPED state",
		},
		"StartedToStarted": {
			initialState: StartedState,
			targetState:  StartedState,
			description:  "Verifies a processor in STARTED state can be updated while remaining in STARTED state",
		},
		"StoppedToStarted": {
			setupState:   StartedState, // Must first get to STARTED before we can test STOPPED→STARTED
			initialState: StoppedState,
			targetState:  StartedState,
			description:  "Verifies a processor can transition from STOPPED to STARTED state",
		},
		"StoppedToStopped": {
			setupState:   StartedState, // Must first get to STARTED before we can test STOPPED→STOPPED
			initialState: StoppedState,
			targetState:  StoppedState,
			description:  "Verifies a processor in STOPPED state can be updated while remaining in STOPPED state",
		},
	}

	for name, tc := range transitions {
		t.Run(name, func(t *testing.T) {
			t.Logf("Testing: %s", tc.description)
			testAccStreamProcessorStateTransitionForUpdates(t, tc.setupState, tc.initialState, tc.targetState, "")
		})
	}
}

// when an empty state is provided in a stream processor update, it should use the existing current state
func TestAccStreamProcessor_EmptyStateUpdates(t *testing.T) {
	transitions := map[string]struct {
		setupState   string // Optional: Initial setup state (needed for STOPPED tests since we can't create in a STOPPED state)
		initialState string // State to create or transition to after setup
		targetState  string // Final state to transition to
		description  string // Description of what the test validates
	}{
		"CreatedToEmptyState": {
			initialState: CreatedState,
			targetState:  "",
			description:  "Verifies that a processor in CREATED state can be updated while remaining in a derived CREATED state from empty state",
		},
		"StartedToEmptyState": {
			initialState: StartedState,
			targetState:  "",
			description:  "Verifies that a processor in STARTED state can be updated while remaining in a derived STARTED state from empty state",
		},
		"StoppedToEmptyState": {
			setupState:   StartedState, // Must first get to STARTED before we can test STOPPED→EMPTY
			initialState: StoppedState,
			targetState:  "",
			description:  "Verifies that a processor in STOPPED state can be updated while remaining in a derived STOPPED state from empty state",
		},
	}

	for name, tc := range transitions {
		t.Run(name, func(t *testing.T) {
			t.Logf("Testing: %s", tc.description)
			testAccStreamProcessorStateTransitionForUpdates(t, tc.setupState, tc.initialState, tc.targetState, "")
		})
	}
}

func TestAccStreamProcessor_InvalidStateTransitionUpdates(t *testing.T) {
	transitions := map[string]struct {
		setupState    string // Optional: Initial setup state (needed for STOPPED tests)
		initialState  string // State to create or transition to after setup
		targetState   string // Final state to transition to
		expectedError string
		description   string // Description of what the test validates
	}{
		"CreatedToStopped": {
			initialState:  CreatedState,
			targetState:   StoppedState,
			expectedError: fmt.Sprintf(streamprocessor.ErrorUpdateStateTransition, StartedState, StoppedState),
			description:   "Verifies a processor cannot transition from CREATED to STOPPED state",
		},
		"StoppedToCreated": {
			setupState:    StartedState, // Must first get to STARTED before we can test STOPPED→CREATED
			initialState:  StoppedState,
			targetState:   CreatedState,
			expectedError: fmt.Sprintf(streamprocessor.ErrorUpdateToCreatedState, StoppedState),
			description:   "Verifies a processor cannot transition from STOPPED to CREATED state",
		},
		"StartedToCreated": {
			initialState:  StartedState,
			targetState:   CreatedState,
			expectedError: fmt.Sprintf(streamprocessor.ErrorUpdateToCreatedState, StartedState),
			description:   "Verifies a processor cannot transition from STARTED to CREATED state",
		},
	}

	for name, tc := range transitions {
		t.Run(name, func(t *testing.T) {
			t.Logf("Testing: %s", tc.description)
			testAccStreamProcessorStateTransitionForUpdates(t, tc.setupState, tc.initialState, tc.targetState, tc.expectedError)
		})
	}
}

func TestAccStreamProcessor_clusterType(t *testing.T) {
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		_, clusterName           = acc.ClusterNameExecution(t, false)
		randomSuffix             = acctest.RandString(5)
		processorName            = "new-processor" + randomSuffix
		srcConfig                = connectionConfig{connectionType: connTypeCluster, clusterName: clusterName, pipelineStepIsSource: true}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:            config(t, projectID, workspaceName, processorName, streamprocessor.StartedState, randomSuffix, srcConfig, testLogDestConfig, "", nil),
				Check:             composeStreamProcessorChecks(projectID, workspaceName, processorName, streamprocessor.StartedState, true, false),
				ConfigStateChecks: pluralConfigStateChecks(processorName, streamprocessor.StartedState, workspaceName, true, false),
			},
		}})
}

func TestAccStreamProcessor_createErrors(t *testing.T) {
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		processorName            = "new-processor"
		invalidJSONConfig        = connectionConfig{connectionType: connTypeSample, pipelineStepIsSource: true, invalidJSON: true}
		randomSuffix             = acctest.RandString(5)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:      config(t, projectID, workspaceName, processorName, streamprocessor.StoppedState, randomSuffix, invalidJSONConfig, testLogDestConfig, "", nil),
				ExpectError: regexp.MustCompile("Invalid JSON String Value"),
			},
			{
				Config:      config(t, projectID, workspaceName, processorName, streamprocessor.StoppedState, randomSuffix, sampleSrcConfig, testLogDestConfig, "", nil),
				ExpectError: regexp.MustCompile("When creating a stream processor, the only valid states are CREATED and STARTED"),
			},
		}})
}

func TestAccStreamProcessor_createTimeoutWithDeleteOnCreate(t *testing.T) {
	acc.SkipTestForCI(t) // Creation of stream processor for testing is too fast to force the creation timeout
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		processorName            = "new-processor"
		randomSuffix             = acctest.RandString(5)
		createTimeout            = "1s"
		deleteOnCreateTimeout    = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps: []resource.TestStep{
			{
				Config:      config(t, projectID, workspaceName, processorName, streamprocessor.StartedState, randomSuffix, sampleSrcConfig, testLogDestConfig, acc.TimeoutConfig(&createTimeout, nil, nil), &deleteOnCreateTimeout),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
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
		workspaceName := rs.Primary.Attributes["workspace_name"]
		processorName := rs.Primary.Attributes["processor_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamProcessor(context.Background(), projectID, workspaceName, processorName).Execute()

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
		workspaceName := rs.Primary.Attributes["workspace_name"]
		processorName := rs.Primary.Attributes["processor_name"]
		_, _, err := acc.ConnV2().StreamsApi.GetStreamProcessor(context.Background(), projectID, workspaceName, processorName).Execute()
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

		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["workspace_name"], rs.Primary.Attributes["project_id"], rs.Primary.Attributes["processor_name"]), nil
	}
}

// testAccStreamProcessorStateTransition tests a state transition with optional setup state
// If errorPattern is not empty, the final transition is expected to fail with an error matching that pattern
func testAccStreamProcessorStateTransitionForUpdates(t *testing.T, setupState, initialState, targetState, errorPattern string) {
	t.Helper()
	var (
		projectID, workspaceName = acc.ProjectIDExecutionWithStreamInstance(t)
		processorName            = fmt.Sprintf("processor-%s-to-%s", strings.ToLower(initialState), strings.ToLower(targetState))
	)

	initialPipeline := `[
		{
			"$source": {
				"connectionName":"sample_stream_solar"
			}
		},
		{
			"$emit": {
				"connectionName":"__testLog"
			}
		}
	]`

	updatedPipelineWithTumblingWindow := `[
		{
			"$source": {
				"connectionName": "sample_stream_solar"
			}
		},
		{
			"$tumblingWindow": {
				"interval": { 
					"size": 10, 
					"unit": "second" 
				},
				"pipeline": [
					{
						"$group": {
							"_id": "$group_id",
							"max_temp": { "$avg": "$obs.temp" },
							"avg_watts": { "$min": "$obs.watts" }
						}
					}
				]
			}
		},
		{
			"$emit": {
				"connectionName": "__testLog"
			}
		}
	]`

	steps := []resource.TestStep{}
	// Add setup step if needed (required for testing from STOPPED state)
	if setupState != "" {
		setupStateConfig := fmt.Sprintf(`state = %q`, setupState)
		steps = append(steps, resource.TestStep{
			Config: configToUpdateStreamProcessor(projectID, workspaceName, processorName, setupStateConfig, initialPipeline),
			Check:  checkAttributesFromBasicUpdateFlow(projectID, workspaceName, processorName, setupState, initialPipeline),
		})
	}

	var initialStateConfig string
	if initialState != "" {
		initialStateConfig = fmt.Sprintf(`state = %q`, initialState)
	}
	steps = append(steps, resource.TestStep{
		Config: configToUpdateStreamProcessor(projectID, workspaceName, processorName, initialStateConfig, initialPipeline),
		Check:  checkAttributesFromBasicUpdateFlow(projectID, workspaceName, processorName, initialState, initialPipeline),
	})

	var targetStateConfig string
	if targetState != "" {
		targetStateConfig = fmt.Sprintf(`state = %q`, targetState)
	}
	// Add target state step, with error checking if applicable
	finalStep := resource.TestStep{
		Config: configToUpdateStreamProcessor(projectID, workspaceName, processorName, targetStateConfig, updatedPipelineWithTumblingWindow),
	}

	// Configure the step based on whether we expect success or failure
	if errorPattern != "" {
		// For invalid transitions, expect an error
		finalStep.ExpectError = regexp.MustCompile(errorPattern)
	} else {
		// if the empty state is passed in the target config, this should be the same initial state of the config
		if targetState == "" {
			targetState = initialState
		}
		finalStep.Check = checkAttributesFromBasicUpdateFlow(projectID, workspaceName, processorName, targetState, updatedPipelineWithTumblingWindow)
	}
	steps = append(steps, finalStep)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroyStreamProcessor,
		Steps:                    steps,
	})
}

// configToUpdateStreamProcessor generates Terraform configuration for Stream Processor state transition tests.
// It creates a minimal test environment with a stream instance, sample source connection and pipelines that can be updated
func configToUpdateStreamProcessor(projectID, workspaceName, processorName, state, pipeline string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_stream_processor" "processor" {
			project_id     = %[1]q
			workspace_name  = %[2]q
			processor_name = %[3]q
			pipeline       = %[4]q
			%[5]s
		}
		`, projectID, workspaceName, processorName, pipeline, state)
}

func checkAttributesFromBasicUpdateFlow(projectID, workspaceName, processorName, state, expectedPipelineStr string) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	attributes := map[string]string{
		"project_id":     projectID,
		"workspace_name": workspaceName,
		"processor_name": processorName,
		"state":          state,
		"pipeline":       expectedPipelineStr,
	}

	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

// pluralConfigStateChecks allows checking one of the results returned by the plural data source
func pluralConfigStateChecks(processorName, state, workspaceName string, includeStats, includeOptions bool) []statecheck.StateCheck {
	return []statecheck.StateCheck{
		acc.PluralResultCheck(pluralDataSourceName, "processor_name", knownvalue.StringExact(processorName), pluralValueChecks(processorName, state, workspaceName, includeStats, includeOptions)),
	}
}

func pluralValueChecks(processorName, state, workspaceName string, includeStats, includeOptions bool) map[string]knownvalue.Check {
	checks := map[string]knownvalue.Check{
		"processor_name": knownvalue.StringExact(processorName),
		"state":          knownvalue.StringExact(state),
		"workspace_name": knownvalue.StringExact(workspaceName),
	}
	if includeStats {
		checks["stats"] = knownvalue.NotNull()
	}
	if includeOptions {
		checks["options"] = knownvalue.NotNull()
	}
	return checks
}

func pluralConfigStateChecksMigration(processorName, state, instanceName string, includeStats, includeOptions bool) []statecheck.StateCheck {
	checks := map[string]knownvalue.Check{
		"processor_name": knownvalue.StringExact(processorName),
		"state":          knownvalue.StringExact(state),
		"instance_name":  knownvalue.StringExact(instanceName),
	}
	if includeStats {
		checks["stats"] = knownvalue.NotNull()
	}
	if includeOptions {
		checks["options"] = knownvalue.NotNull()
	}

	return []statecheck.StateCheck{
		acc.PluralResultCheck(pluralDataSourceName, "processor_name", knownvalue.StringExact(processorName), checks),
	}
}

func composeStreamProcessorChecks(projectID, workspaceName, processorName, state string, includeStats, includeOptions bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	attributes := map[string]string{
		"project_id":     projectID,
		"workspace_name": workspaceName,
		"processor_name": processorName,
		"state":          state,
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrChecks(dataSourceName, checks, attributes)
	if includeStats {
		checks = acc.AddAttrSetChecks(resourceName, checks, "stats", "pipeline")
		checks = acc.AddAttrSetChecks(dataSourceName, checks, "stats", "pipeline")
	}
	if includeOptions {
		checks = acc.AddAttrSetChecks(resourceName, checks, "options.dlq.db", "options.dlq.coll", "options.dlq.connection_name")
		checks = acc.AddAttrSetChecks(dataSourceName, checks, "options.dlq.db", "options.dlq.coll", "options.dlq.connection_name")
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func composeStreamProcessorChecksMigration(projectID, instanceName, processorName, state string, includeStats, includeOptions bool) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	attributes := map[string]string{
		"project_id":     projectID,
		"instance_name":  instanceName,
		"processor_name": processorName,
		"state":          state,
	}
	checks = acc.AddAttrChecks(resourceName, checks, attributes)
	checks = acc.AddAttrChecks(dataSourceName, checks, attributes)
	if includeStats {
		checks = acc.AddAttrSetChecks(resourceName, checks, "stats", "pipeline")
		checks = acc.AddAttrSetChecks(dataSourceName, checks, "stats", "pipeline")
	}
	if includeOptions {
		checks = acc.AddAttrSetChecks(resourceName, checks, "options.dlq.db", "options.dlq.coll", "options.dlq.connection_name")
		checks = acc.AddAttrSetChecks(dataSourceName, checks, "options.dlq.db", "options.dlq.coll", "options.dlq.connection_name")
	}
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func config(t *testing.T, projectID, workspaceName, processorName, state, nameSuffix string, src, dest connectionConfig, timeoutConfig string, deleteOnCreateTimeout *bool) string {
	t.Helper()
	stateConfig := ""
	if state != "" {
		stateConfig = fmt.Sprintf(`state = %[1]q`, state)
	}
	deleteOnCreateTimeoutConfig := ""
	if deleteOnCreateTimeout != nil {
		deleteOnCreateTimeoutConfig = fmt.Sprintf(`delete_on_create_timeout = %[1]t`, *deleteOnCreateTimeout)
	}

	connectionConfigSrc, connectionIDSrc, pipelineStepSrc := configConnection(t, projectID, workspaceName, src, nameSuffix)
	connectionConfigDest, connectionIDDest, pipelineStepDest := configConnection(t, projectID, workspaceName, dest, nameSuffix)
	dependsOn := []string{}
	if connectionIDSrc != "" && !strings.HasPrefix(connectionIDSrc, "data.") {
		dependsOn = append(dependsOn, connectionIDSrc)
	}
	if connectionIDDest != "" && !strings.HasPrefix(connectionIDDest, "data.") {
		dependsOn = append(dependsOn, connectionIDDest)
	}
	dependsOnStr := strings.Join(dependsOn, ", ")
	pipeline := fmt.Sprintf("[{\"$source\":%1s},{\"$emit\":%2s}]", pipelineStepSrc, pipelineStepDest)
	optionsStr := ""
	if src.useAsDLQ {
		assert.Equal(t, connTypeCluster, src.connectionType)
		optionsStr = fmt.Sprintf(`
			options = {
				dlq = {
					coll = "dlq_coll"
					connection_name = %[1]s.connection_name
					db = "dlq_db"
				}
			}`, connectionIDSrc)
	}

	dataSource := fmt.Sprintf(`
	data "mongodbatlas_stream_processor" "test" {
		project_id = %[1]q
		workspace_name = %[2]q
		processor_name = %[3]q
		depends_on = [%4s]
	}`, projectID, workspaceName, processorName, resourceName)
	dataSourcePlural := fmt.Sprintf(`
	data "mongodbatlas_stream_processors" "test" {
		project_id = %[1]q
		workspace_name = %[2]q
		depends_on = [%3s]
	}`, projectID, workspaceName, resourceName)
	otherConfig := connectionConfigSrc + connectionConfigDest + dataSource + dataSourcePlural

	return fmt.Sprintf(`
	resource "mongodbatlas_stream_processor" "processor" {
		project_id     = %[1]q
		workspace_name  = %[2]q
		processor_name = %[3]q
		pipeline       = %[4]q
		%[5]s
		%[6]s
		depends_on = [%[7]s]
		%[8]s
		%[9]s
		}
		
	`, projectID, workspaceName, processorName, pipeline, stateConfig, optionsStr, dependsOnStr, timeoutConfig, deleteOnCreateTimeoutConfig) + otherConfig
}

func configMigration(t *testing.T, projectID, instanceName, processorName, state, nameSuffix string, src, dest connectionConfig, timeoutConfig string, deleteOnCreateTimeout *bool) string {
	t.Helper()
	stateConfig := ""
	if state != "" {
		stateConfig = fmt.Sprintf(`state = %[1]q`, state)
	}
	deleteOnCreateTimeoutConfig := ""
	if deleteOnCreateTimeout != nil {
		deleteOnCreateTimeoutConfig = fmt.Sprintf(`delete_on_create_timeout = %[1]t`, *deleteOnCreateTimeout)
	}

	connectionConfigSrc, connectionIDSrc, pipelineStepSrc := configConnectionMigration(t, projectID, instanceName, src, nameSuffix)
	connectionConfigDest, connectionIDDest, pipelineStepDest := configConnectionMigration(t, projectID, instanceName, dest, nameSuffix)
	dependsOn := []string{}
	if connectionIDSrc != "" && !strings.HasPrefix(connectionIDSrc, "data.") {
		dependsOn = append(dependsOn, connectionIDSrc)
	}
	if connectionIDDest != "" && !strings.HasPrefix(connectionIDDest, "data.") {
		dependsOn = append(dependsOn, connectionIDDest)
	}
	dependsOnStr := strings.Join(dependsOn, ", ")
	pipeline := fmt.Sprintf("[{\"$source\":%1s},{\"$emit\":%2s}]", pipelineStepSrc, pipelineStepDest)
	optionsStr := ""
	if src.useAsDLQ {
		assert.Equal(t, connTypeCluster, src.connectionType)
		optionsStr = fmt.Sprintf(`
			options = {
				dlq = {
					coll = "dlq_coll"
					connection_name = %[1]s.connection_name
					db = "dlq_db"
				}
			}`, connectionIDSrc)
	}

	dataSource := fmt.Sprintf(`
	data "mongodbatlas_stream_processor" "test" {
		project_id = %[1]q
		instance_name = %[2]q
		processor_name = %[3]q
		depends_on = [%4s]
	}`, projectID, instanceName, processorName, resourceName)
	dataSourcePlural := fmt.Sprintf(`
	data "mongodbatlas_stream_processors" "test" {
		project_id = %[1]q
		instance_name = %[2]q
		depends_on = [%3s]
	}`, projectID, instanceName, resourceName)
	otherConfig := connectionConfigSrc + connectionConfigDest + dataSource + dataSourcePlural

	return fmt.Sprintf(`
	resource "mongodbatlas_stream_processor" "processor" {
		project_id     = %[1]q
		instance_name  = %[2]q
		processor_name = %[3]q
		pipeline       = %[4]q
		%[5]s
		%[6]s
		depends_on = [%[7]s]
		%[8]s
		%[9]s
		}
		
	`, projectID, instanceName, processorName, pipeline, stateConfig, optionsStr, dependsOnStr, timeoutConfig, deleteOnCreateTimeoutConfig) + otherConfig
}

func configConnection(t *testing.T, projectID, workspaceName string, config connectionConfig, nameSuffix string) (connectionConfig, resourceID, pipelineStep string) {
	t.Helper()
	assert.False(t, config.extraWhitespace && config.connectionType != connTypeSample, "extraWhitespace is only supported for Sample connection")
	assert.False(t, config.invalidJSON && config.connectionType != connTypeSample, "invalidJson is only supported for Sample connection")
	connectionType := config.connectionType
	pipelineStepIsSource := config.pipelineStepIsSource
	switch connectionType {
	case "Cluster":
		var connectionName, resourceName string
		clusterName := config.clusterName
		assert.NotEmpty(t, clusterName)
		if pipelineStepIsSource {
			connectionName = "ClusterConnectionSrc" + nameSuffix
			resourceName = "cluster_src"
		} else {
			connectionName = "ClusterConnectionDest" + nameSuffix
			resourceName = "cluster_dest"
		}
		connectionConfig = fmt.Sprintf(`
            resource "mongodbatlas_stream_connection" %[4]q {
                project_id      = %[1]q
                cluster_name    = %[2]q
                workspace_name   = %[5]q
                connection_name = %[3]q
                type            = "Cluster"
                db_role_to_execute = {
                    role = "atlasAdmin"
                    type = "BUILT_IN"
                }
            }
        `, projectID, clusterName, connectionName, resourceName, workspaceName)
		resourceID = fmt.Sprintf("mongodbatlas_stream_connection.%s", resourceName)
		pipelineStep = fmt.Sprintf("{\"connectionName\":%q}", connectionName)
		return connectionConfig, resourceID, pipelineStep
	case "Kafka":
		var connectionName, resourceName, pipelineStep string
		if pipelineStepIsSource {
			connectionName = "KafkaConnectionSrc" + nameSuffix
			resourceName = "kafka_src"
			pipelineStep = fmt.Sprintf("{\"connectionName\":%q}", connectionName)
		} else {
			connectionName = "KafkaConnectionDest" + nameSuffix
			resourceName = "kafka_dest"
			pipelineStep = fmt.Sprintf("{\"connectionName\":%q,\"topic\":\"random_topic\"}", connectionName)
		}
		connectionConfig = fmt.Sprintf(`
            resource "mongodbatlas_stream_connection" %[3]q{
                project_id      = %[1]q
                workspace_name   = %[4]q
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
                    protocol = "SASL_PLAINTEXT"
                }
            }
        `, projectID, connectionName, resourceName, workspaceName)
		resourceID = fmt.Sprintf("mongodbatlas_stream_connection.%s", resourceName)
		return connectionConfig, resourceID, pipelineStep
	case "Sample":
		if !pipelineStepIsSource {
			t.Fatal("Sample connection must be used as a source")
		}
		connectionConfig = fmt.Sprintf(`
            data "mongodbatlas_stream_connection" "sample" {
                project_id      = %[1]q
                workspace_name   = %[2]q
                connection_name = "sample_stream_solar"
            }
        `, projectID, workspaceName)
		resourceID = "data.mongodbatlas_stream_connection.sample"
		if config.extraWhitespace {
			pipelineStep = "{\"connectionName\": \"sample_stream_solar\"}"
		} else {
			pipelineStep = "{\"connectionName\":\"sample_stream_solar\"}"
		}
		if config.invalidJSON {
			pipelineStep = "{\"connectionName\": \"sample_stream_solar\"" // missing closing bracket
		}
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

func configConnectionMigration(t *testing.T, projectID, instanceName string, config connectionConfig, nameSuffix string) (connectionConfig, resourceID, pipelineStep string) {
	t.Helper()
	assert.False(t, config.extraWhitespace && config.connectionType != connTypeSample, "extraWhitespace is only supported for Sample connection")
	assert.False(t, config.invalidJSON && config.connectionType != connTypeSample, "invalidJson is only supported for Sample connection")
	connectionType := config.connectionType
	pipelineStepIsSource := config.pipelineStepIsSource
	switch connectionType {
	case "Cluster":
		var connectionName, resourceName string
		clusterName := config.clusterName
		assert.NotEmpty(t, clusterName)
		if pipelineStepIsSource {
			connectionName = "ClusterConnectionSrc" + nameSuffix
			resourceName = "cluster_src"
		} else {
			connectionName = "ClusterConnectionDest" + nameSuffix
			resourceName = "cluster_dest"
		}
		connectionConfig = fmt.Sprintf(`
            resource "mongodbatlas_stream_connection" %[4]q {
                project_id      = %[1]q
                cluster_name    = %[2]q
                workspace_name   = %[5]q
                connection_name = %[3]q
                type            = "Cluster"
                db_role_to_execute = {
                    role = "atlasAdmin"
                    type = "BUILT_IN"
                }
            }
        `, projectID, clusterName, connectionName, resourceName, workspaceName)
		resourceID = fmt.Sprintf("mongodbatlas_stream_connection.%s", resourceName)
		pipelineStep = fmt.Sprintf("{\"connectionName\":%q}", connectionName)
		return connectionConfig, resourceID, pipelineStep
	case "Kafka":
		var connectionName, resourceName, pipelineStep string
		if pipelineStepIsSource {
			connectionName = "KafkaConnectionSrc" + nameSuffix
			resourceName = "kafka_src"
			pipelineStep = fmt.Sprintf("{\"connectionName\":%q}", connectionName)
		} else {
			connectionName = "KafkaConnectionDest" + nameSuffix
			resourceName = "kafka_dest"
			pipelineStep = fmt.Sprintf("{\"connectionName\":%q,\"topic\":\"random_topic\"}", connectionName)
		}
		connectionConfig = fmt.Sprintf(`
            resource "mongodbatlas_stream_connection" %[3]q{
                project_id      = %[1]q
                instanceName   = %[4]q
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
                    protocol = "SASL_PLAINTEXT"
                }
            }
        `, projectID, connectionName, resourceName, workspaceName)
		resourceID = fmt.Sprintf("mongodbatlas_stream_connection.%s", resourceName)
		return connectionConfig, resourceID, pipelineStep
	case "Sample":
		if !pipelineStepIsSource {
			t.Fatal("Sample connection must be used as a source")
		}
		connectionConfig = fmt.Sprintf(`
            data "mongodbatlas_stream_connection" "sample" {
                project_id      = %[1]q
                instance_name   = %[2]q
                connection_name = "sample_stream_solar"
            }
        `, projectID, workspaceName)
		resourceID = "data.mongodbatlas_stream_connection.sample"
		if config.extraWhitespace {
			pipelineStep = "{\"connectionName\": \"sample_stream_solar\"}"
		} else {
			pipelineStep = "{\"connectionName\":\"sample_stream_solar\"}"
		}
		if config.invalidJSON {
			pipelineStep = "{\"connectionName\": \"sample_stream_solar\"" // missing closing bracket
		}
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
