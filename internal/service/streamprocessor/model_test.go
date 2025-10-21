package streamprocessor_test

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

var (
	projectID                 = "661fe3ad234b02027dabcabc"
	instanceName              = "test-instance-name"
	workspaceName             = "test-workspace-name"
	pipelineStageSourceSample = map[string]any{
		"$source": map[string]any{
			"connectionName": "sample_stream_solar",
		},
	}
	pipelineStageEmitLog = map[string]any{
		"$emit": map[string]any{
			"connectionName": "__testLog",
		},
	}
	processorName        = "processor1"
	processorID          = "66b39806187592e8d721215d"
	stateCreated         = streamprocessor.CreatedState
	stateStarted         = streamprocessor.StartedState
	streamOptionsExample = admin.StreamsOptions{
		Dlq: &admin.StreamsDLQ{
			Coll:           conversion.StringPtr("testColl"),
			ConnectionName: conversion.StringPtr("testConnection"),
			Db:             conversion.StringPtr("testDB"),
		},
	}
)

var statsExample = `
{
	"dlqMessageCount": 0,
	"dlqMessageSize": 0.0,
	"inputMessageCount": 12,
	"inputMessageSize": 4681.0,
	"memoryTrackerBytes": 0.0,
	"name": "processor1",
	"ok": 1.0,
	"changeStreamState": { "_data": "8266C37388000000012B0429296E1404" },
	"operatorStats": [
		{
			"dlqMessageCount": 0,
			"dlqMessageSize": 0.0,
			"executionTimeSecs": 0,
			"inputMessageCount": 12,
			"inputMessageSize": 4681.0,
			"maxMemoryUsage": 0.0,
			"name": "SampleDataSourceOperator",
			"outputMessageCount": 12,
			"outputMessageSize": 0.0,
			"stateSize": 0.0,
			"timeSpentMillis": 0
		},
		{
			"dlqMessageCount": 0,
			"dlqMessageSize": 0.0,
			"executionTimeSecs": 0,
			"inputMessageCount": 12,
			"inputMessageSize": 4681.0,
			"maxMemoryUsage": 0.0,
			"name": "LogSinkOperator",
			"outputMessageCount": 12,
			"outputMessageSize": 4681.0,
			"stateSize": 0.0,
			"timeSpentMillis": 0
		}
	],
	"outputMessageCount": 12,
	"outputMessageSize": 4681.0,
	"processorId": "66b3941109bbccf048ccff06",
	"scaleFactor": 1,
	"stateSize": 0.0,
	"status": "running"
}`

func streamProcessorWithStats(t *testing.T, options *admin.StreamsOptions) *admin.StreamsProcessorWithStats {
	t.Helper()
	processor := admin.NewStreamsProcessorWithStats(
		processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, stateStarted,
	)
	var stats any
	err := json.Unmarshal([]byte(statsExample), &stats)
	require.NoError(t, err)
	processor.SetStats(stats)
	if options != nil {
		processor.SetOptions(*options)
	}
	return processor
}

func streamProcessorDSTFModel(t *testing.T, state, stats string, options types.Object) *streamprocessor.TFStreamProcessorDSModel {
	t.Helper()
	return &streamprocessor.TFStreamProcessorDSModel{
		ID:            types.StringValue(processorID),
		WorkspaceName: types.StringValue(workspaceName),
		Options:       options,
		Pipeline:      types.StringValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
		ProcessorName: types.StringValue(processorName),
		ProjectID:     types.StringValue(projectID),
		State:         conversion.StringNullIfEmpty(state),
		Stats:         conversion.StringNullIfEmpty(stats),
	}
}

func streamProcessorDSTFModelWithInstanceName(t *testing.T, state, stats string, options types.Object) *streamprocessor.TFStreamProcessorDSModel {
	t.Helper()
	return &streamprocessor.TFStreamProcessorDSModel{
		ID:            types.StringValue(processorID),
		InstanceName:  types.StringValue(instanceName),
		Options:       options,
		Pipeline:      types.StringValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
		ProcessorName: types.StringValue(processorName),
		ProjectID:     types.StringValue(projectID),
		State:         conversion.StringNullIfEmpty(state),
		Stats:         conversion.StringNullIfEmpty(stats),
	}
}

func optionsToTFModel(t *testing.T, options *admin.StreamsOptions) types.Object {
	t.Helper()
	result, diags := streamprocessor.ConvertOptionsToTF(t.Context(), options)
	assert.False(t, diags.HasError())
	assert.NotNil(t, result)
	return *result
}

func TestDSSDKToTFModel(t *testing.T) {
	testCases := map[string]struct {
		sdkModel        *admin.StreamsProcessorWithStats
		expectedTFModel *streamprocessor.TFStreamProcessorDSModel
	}{
		"afterCreate": {
			sdkModel: admin.NewStreamsProcessorWithStats(
				processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, stateCreated,
			),
			expectedTFModel: streamProcessorDSTFModel(t, stateCreated, "", optionsToTFModel(t, nil)),
		},
		"afterStarted": {
			sdkModel:        streamProcessorWithStats(t, nil),
			expectedTFModel: streamProcessorDSTFModel(t, stateStarted, statsExample, optionsToTFModel(t, nil)),
		},
		"withOptions": {
			sdkModel:        streamProcessorWithStats(t, &streamOptionsExample),
			expectedTFModel: streamProcessorDSTFModel(t, stateStarted, statsExample, optionsToTFModel(t, &streamOptionsExample)),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			resultModel, diags := streamprocessor.NewTFStreamprocessorDSModel(t.Context(), projectID, "", workspaceName, sdkModel)
			assert.False(t, diags.HasError())
			assert.Equal(t, tc.expectedTFModel.Options, resultModel.Options)
			if sdkModel.Stats != nil {
				assert.True(t, schemafunc.EqualJSON(resultModel.Pipeline.String(), tc.expectedTFModel.Pipeline.String(), "test stream processor schema"))
				var statsResult any
				err := json.Unmarshal([]byte(resultModel.Stats.ValueString()), &statsResult)
				require.NoError(t, err)
				assert.Len(t, sdkModel.Stats, 15)
				assert.Len(t, statsResult, 15)
			} else {
				assert.Equal(t, tc.expectedTFModel, resultModel)
			}
		})
	}
}

// TestDSSDKToTFModelInstanceName ensures that deprecated instance_name functionality is still supported
func TestDSSDKToTFModelInstanceName(t *testing.T) {
	testCases := map[string]struct {
		sdkModel        *admin.StreamsProcessorWithStats
		expectedTFModel *streamprocessor.TFStreamProcessorDSModel
	}{
		"afterCreate": {
			sdkModel: admin.NewStreamsProcessorWithStats(
				processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, stateCreated,
			),
			expectedTFModel: streamProcessorDSTFModelWithInstanceName(t, stateCreated, "", optionsToTFModel(t, nil)),
		},
		"afterStarted": {
			sdkModel:        streamProcessorWithStats(t, nil),
			expectedTFModel: streamProcessorDSTFModelWithInstanceName(t, stateStarted, statsExample, optionsToTFModel(t, nil)),
		},
		"withOptions": {
			sdkModel:        streamProcessorWithStats(t, &streamOptionsExample),
			expectedTFModel: streamProcessorDSTFModelWithInstanceName(t, stateStarted, statsExample, optionsToTFModel(t, &streamOptionsExample)),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			resultModel, diags := streamprocessor.NewTFStreamprocessorDSModel(t.Context(), projectID, instanceName, "", sdkModel)
			assert.False(t, diags.HasError())
			assert.Equal(t, tc.expectedTFModel.Options, resultModel.Options)
			if sdkModel.Stats != nil {
				assert.True(t, schemafunc.EqualJSON(resultModel.Pipeline.String(), tc.expectedTFModel.Pipeline.String(), "test stream processor schema"))
				var statsResult any
				err := json.Unmarshal([]byte(resultModel.Stats.ValueString()), &statsResult)
				require.NoError(t, err)
				assert.Len(t, sdkModel.Stats, 15)
				assert.Len(t, statsResult, 15)
			} else {
				assert.Equal(t, tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestSDKToTFModel(t *testing.T) {
	testCases := map[string]struct {
		sdkModel        *admin.StreamsProcessorWithStats
		expectedTFModel *streamprocessor.TFStreamProcessorRSModel
	}{
		"afterCreate": {
			sdkModel: admin.NewStreamsProcessorWithStats(
				processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, "CREATED",
			),
			expectedTFModel: &streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringValue(workspaceName),
				Options:       types.ObjectNull(streamprocessor.OptionsObjectType.AttrTypes),
				ProcessorID:   types.StringValue(processorID),
				Pipeline:      jsontypes.NewNormalizedValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
				State:         types.StringValue("CREATED"),
				Stats:         types.StringNull(),
			},
		},
		"afterStarted": {
			sdkModel: streamProcessorWithStats(t, nil),
			expectedTFModel: &streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringValue(workspaceName),
				Options:       types.ObjectNull(streamprocessor.OptionsObjectType.AttrTypes),
				ProcessorID:   types.StringValue(processorID),
				Pipeline:      jsontypes.NewNormalizedValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
				State:         types.StringValue("STARTED"),
				Stats:         types.StringValue(statsExample),
			},
		},
		"withOptions": {
			sdkModel: streamProcessorWithStats(t, &streamOptionsExample),
			expectedTFModel: &streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringValue(workspaceName),
				Options:       optionsToTFModel(t, &streamOptionsExample),
				ProcessorID:   types.StringValue(processorID),
				Pipeline:      jsontypes.NewNormalizedValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
				State:         types.StringValue("STARTED"),
				Stats:         types.StringNull(),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			resultModel, diags := streamprocessor.NewStreamProcessorWithStats(t.Context(), projectID, workspaceName, "", sdkModel, nil, nil)
			assert.False(t, diags.HasError())
			assert.Equal(t, tc.expectedTFModel.Options, resultModel.Options)
			if sdkModel.Stats != nil {
				assert.True(t, schemafunc.EqualJSON(resultModel.Pipeline.String(), tc.expectedTFModel.Pipeline.String(), "test stream processor schema"))
				var statsResult any
				err := json.Unmarshal([]byte(resultModel.Stats.ValueString()), &statsResult)
				require.NoError(t, err)
				assert.Len(t, sdkModel.Stats, 15)
				assert.Len(t, statsResult, 15)
			} else {
				assert.Equal(t, tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestPluralDSSDKToTFModel(t *testing.T) {
	testCases := map[string]struct {
		sdkModel         *admin.PaginatedApiStreamsStreamProcessorWithStats
		expectedTFModel  *streamprocessor.TFStreamProcessorsDSModel
		useWorkspaceName bool
	}{
		"noResults_with_workspace_name": {
			sdkModel: &admin.PaginatedApiStreamsStreamProcessorWithStats{
				Results:    &[]admin.StreamsProcessorWithStats{},
				TotalCount: admin.PtrInt(0),
			},
			expectedTFModel: &streamprocessor.TFStreamProcessorsDSModel{
				ProjectID:     types.StringValue(projectID),
				WorkspaceName: types.StringValue(workspaceName),
				Results:       []streamprocessor.TFStreamProcessorDSModel{},
			},
		},
		"oneResult_with_workspace_name": {
			sdkModel: &admin.PaginatedApiStreamsStreamProcessorWithStats{
				Results: &[]admin.StreamsProcessorWithStats{*admin.NewStreamsProcessorWithStats(
					processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, stateCreated,
				)},
				TotalCount: admin.PtrInt(1),
			},
			expectedTFModel: &streamprocessor.TFStreamProcessorsDSModel{
				ProjectID:     types.StringValue(projectID),
				WorkspaceName: types.StringValue(workspaceName),
				Results: []streamprocessor.TFStreamProcessorDSModel{
					*streamProcessorDSTFModel(t, stateCreated, "", optionsToTFModel(t, nil)),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			existingConfig := &streamprocessor.TFStreamProcessorsDSModel{
				ProjectID:     types.StringValue(projectID),
				WorkspaceName: types.StringValue(workspaceName),
			}
			resultModel, diags := streamprocessor.NewTFStreamProcessors(t.Context(), existingConfig, sdkModel.GetResults())
			assert.False(t, diags.HasError())
			assert.Equal(t, tc.expectedTFModel, resultModel)
		})
	}
}

// TestPluralDSSDKToTFModelWithInstanceName ensures that deprecated instance_name functionality is still supported
func TestPluralDSSDKToTFModelWithInstanceName(t *testing.T) {
	testCases := map[string]struct {
		sdkModel        *admin.PaginatedApiStreamsStreamProcessorWithStats
		expectedTFModel *streamprocessor.TFStreamProcessorsDSModel
	}{
		"noResults": {sdkModel: &admin.PaginatedApiStreamsStreamProcessorWithStats{
			Results:    &[]admin.StreamsProcessorWithStats{},
			TotalCount: admin.PtrInt(0),
		}, expectedTFModel: &streamprocessor.TFStreamProcessorsDSModel{
			ProjectID:    types.StringValue(projectID),
			InstanceName: types.StringValue(instanceName),
			Results:      []streamprocessor.TFStreamProcessorDSModel{},
		}},
		"oneResult": {sdkModel: &admin.PaginatedApiStreamsStreamProcessorWithStats{
			Results: &[]admin.StreamsProcessorWithStats{*admin.NewStreamsProcessorWithStats(
				processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, stateCreated,
			)},
			TotalCount: admin.PtrInt(1),
		}, expectedTFModel: &streamprocessor.TFStreamProcessorsDSModel{
			ProjectID:    types.StringValue(projectID),
			InstanceName: types.StringValue(instanceName),
			Results: []streamprocessor.TFStreamProcessorDSModel{
				*streamProcessorDSTFModelWithInstanceName(t, stateCreated, "", optionsToTFModel(t, nil)),
			},
		}},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			existingConfig := &streamprocessor.TFStreamProcessorsDSModel{
				ProjectID:    types.StringValue(projectID),
				InstanceName: types.StringValue(instanceName),
			}
			resultModel, diags := streamprocessor.NewTFStreamProcessors(t.Context(), existingConfig, sdkModel.GetResults())
			assert.False(t, diags.HasError())
			assert.Equal(t, tc.expectedTFModel, resultModel)
		})
	}
}

func TestNewStreamProcessorUpdateReq(t *testing.T) {
	validPipeline := jsontypes.NewNormalizedValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]")

	testCases := map[string]struct {
		model          *streamprocessor.TFStreamProcessorRSModel
		expectedResult string
	}{
		"workspace_name provided": {
			model: &streamprocessor.TFStreamProcessorRSModel{
				WorkspaceName: types.StringValue(workspaceName),
				InstanceName:  types.StringNull(),
				Pipeline:      validPipeline,
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
			},
			expectedResult: workspaceName,
		},
		"instance_name provided": {
			model: &streamprocessor.TFStreamProcessorRSModel{
				WorkspaceName: types.StringNull(),
				InstanceName:  types.StringValue(instanceName),
				Pipeline:      validPipeline,
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
			},
			expectedResult: instanceName,
		},
		"workspace_name and instance_name provided": {
			model: &streamprocessor.TFStreamProcessorRSModel{
				WorkspaceName: types.StringValue(workspaceName),
				InstanceName:  types.StringValue(instanceName),
				Pipeline:      validPipeline,
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
			},
			expectedResult: workspaceName,
		},
		"neither provided": {
			model: &streamprocessor.TFStreamProcessorRSModel{
				WorkspaceName: types.StringNull(),
				InstanceName:  types.StringNull(),
				Pipeline:      validPipeline,
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
			},
			expectedResult: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			updateReq, diags := streamprocessor.NewStreamProcessorUpdateReq(t.Context(), tc.model)
			if tc.expectedResult == "" {
				assert.False(t, diags.HasError())
				assert.Empty(t, updateReq.TenantName)
			} else {
				assert.False(t, diags.HasError())
				assert.Equal(t, tc.expectedResult, updateReq.TenantName)
			}
		})
	}
}

func TestGetWorkspaceOrInstanceName(t *testing.T) {
	testCases := map[string]struct {
		workspaceName types.String
		instanceName  types.String
		expected      string
	}{
		"workspace_name provided": {
			workspaceName: types.StringValue(workspaceName),
			instanceName:  types.StringNull(),
			expected:      workspaceName,
		},
		"instance_name provided": {
			workspaceName: types.StringNull(),
			instanceName:  types.StringValue(instanceName),
			expected:      instanceName,
		},
		"both provided": {
			workspaceName: types.StringValue(workspaceName),
			instanceName:  types.StringValue(instanceName),
			expected:      workspaceName,
		},
		"workspace_name empty_string": {
			workspaceName: types.StringValue(""),
			instanceName:  types.StringNull(),
			expected:      "",
		},
		"instance_namee empty_string": {
			workspaceName: types.StringNull(),
			instanceName:  types.StringValue(""),
			expected:      "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := streamprocessor.GetWorkspaceOrInstanceName(tc.workspaceName, tc.instanceName)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}
