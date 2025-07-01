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
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

var (
	projectID                 = "661fe3ad234b02027dabcabc"
	instanceName              = "test-instance-name"
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
	if err != nil {
		t.Fatal(err)
	}
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
	if diags.HasError() {
		t.Fatal(diags)
	}
	assert.NotNil(t, result)
	return *result
}

func TestDSSDKToTFModel(t *testing.T) {
	testCases := []struct {
		sdkModel        *admin.StreamsProcessorWithStats
		expectedTFModel *streamprocessor.TFStreamProcessorDSModel
		name            string
	}{
		{
			name: "afterCreate",
			sdkModel: admin.NewStreamsProcessorWithStats(
				processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, stateCreated,
			),
			expectedTFModel: streamProcessorDSTFModel(t, stateCreated, "", optionsToTFModel(t, nil)),
		},
		{
			name:            "afterStarted",
			sdkModel:        streamProcessorWithStats(t, nil),
			expectedTFModel: streamProcessorDSTFModel(t, stateStarted, statsExample, optionsToTFModel(t, nil)),
		},
		{
			name:            "withOptions",
			sdkModel:        streamProcessorWithStats(t, &streamOptionsExample),
			expectedTFModel: streamProcessorDSTFModel(t, stateStarted, statsExample, optionsToTFModel(t, &streamOptionsExample)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			resultModel, diags := streamprocessor.NewTFStreamprocessorDSModel(t.Context(), projectID, instanceName, sdkModel)
			if diags.HasError() {
				t.Fatalf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel.Options, resultModel.Options)
			if sdkModel.Stats != nil {
				assert.True(t, schemafunc.EqualJSON(resultModel.Pipeline.String(), tc.expectedTFModel.Pipeline.String(), "test stream processor schema"))
				var statsResult any
				err := json.Unmarshal([]byte(resultModel.Stats.ValueString()), &statsResult)
				if err != nil {
					t.Fatal(err)
				}
				assert.Len(t, sdkModel.Stats, 15)
				assert.Len(t, statsResult, 15)
			} else {
				assert.Equal(t, tc.expectedTFModel, resultModel)
			}
		})
	}
}

func TestSDKToTFModel(t *testing.T) {
	testCases := []struct {
		sdkModel        *admin.StreamsProcessorWithStats
		expectedTFModel *streamprocessor.TFStreamProcessorRSModel
		name            string
	}{
		{
			name: "afterCreate",
			sdkModel: admin.NewStreamsProcessorWithStats(
				processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, "CREATED",
			),
			expectedTFModel: &streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringValue(instanceName),
				Options:       types.ObjectNull(streamprocessor.OptionsObjectType.AttrTypes),
				ProcessorID:   types.StringValue(processorID),
				Pipeline:      jsontypes.NewNormalizedValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
				State:         types.StringValue("CREATED"),
				Stats:         types.StringNull(),
			},
		},
		{
			name:     "afterStarted",
			sdkModel: streamProcessorWithStats(t, nil),
			expectedTFModel: &streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringValue(instanceName),
				Options:       types.ObjectNull(streamprocessor.OptionsObjectType.AttrTypes),
				ProcessorID:   types.StringValue(processorID),
				Pipeline:      jsontypes.NewNormalizedValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
				State:         types.StringValue("STARTED"),
				Stats:         types.StringValue(statsExample),
			},
		},
		{
			name:     "withOptions",
			sdkModel: streamProcessorWithStats(t, &streamOptionsExample),
			expectedTFModel: &streamprocessor.TFStreamProcessorRSModel{
				InstanceName:  types.StringValue(instanceName),
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			resultModel, diags := streamprocessor.NewStreamProcessorWithStats(t.Context(), projectID, instanceName, sdkModel)
			if diags.HasError() {
				t.Fatalf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel.Options, resultModel.Options)
			if sdkModel.Stats != nil {
				assert.True(t, schemafunc.EqualJSON(resultModel.Pipeline.String(), tc.expectedTFModel.Pipeline.String(), "test stream processor schema"))
				var statsResult any
				err := json.Unmarshal([]byte(resultModel.Stats.ValueString()), &statsResult)
				if err != nil {
					t.Fatal(err)
				}
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
				*streamProcessorDSTFModel(t, stateCreated, "", optionsToTFModel(t, nil)),
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
			if diags.HasError() {
				t.Fatalf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel)
		})
	}
}
