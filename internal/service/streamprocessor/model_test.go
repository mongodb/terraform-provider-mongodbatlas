package streamprocessor_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamprocessor"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
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
	processorName = "processor1"
	processorID   = "66b39806187592e8d721215d"
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

func streamProcessorWithStats(t *testing.T) *admin.StreamsProcessorWithStats {
	t.Helper()
	processor := admin.NewStreamsProcessorWithStats(
		processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, "STARTED",
	)
	var stats any
	err := json.Unmarshal([]byte(statsExample), &stats)
	if err != nil {
		t.Fatal(err)
	}
	processor.SetStats(stats)
	return processor
}

func TestTeamsDataSourceSDKToTFModel(t *testing.T) {
	testCases := []struct {
		sdkModel        *admin.StreamsProcessorWithStats
		expectedTFModel *streamprocessor.TFStreamprocessorDSModel
		name            string
	}{
		{
			name: "afterCreate",
			sdkModel: admin.NewStreamsProcessorWithStats(
				processorID, processorName, []any{pipelineStageSourceSample, pipelineStageEmitLog}, "CREATED",
			),
			expectedTFModel: &streamprocessor.TFStreamprocessorDSModel{
				ID:            types.StringValue(processorID),
				InstanceName:  types.StringValue(instanceName),
				Pipeline:      types.StringValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
				State:         types.StringValue("CREATED"),
				Stats:         types.StringValue("{}"),
			},
		},
		{
			name:     "afterStarted",
			sdkModel: streamProcessorWithStats(t),
			expectedTFModel: &streamprocessor.TFStreamprocessorDSModel{
				ID:            types.StringValue(processorID),
				InstanceName:  types.StringValue(instanceName),
				Pipeline:      types.StringValue("[{\"$source\":{\"connectionName\":\"sample_stream_solar\"}},{\"$emit\":{\"connectionName\":\"__testLog\"}}]"),
				ProcessorName: types.StringValue(processorName),
				ProjectID:     types.StringValue(projectID),
				State:         types.StringValue("STARTED"),
				Stats:         types.StringValue(statsExample),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sdkModel := tc.sdkModel
			resultModel, diags := streamprocessor.NewTFStreamprocessorDSModel(context.Background(), projectID, instanceName, sdkModel)
			if diags.HasError() {
				t.Fatalf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			if sdkModel.Stats != nil {
				assert.True(t, schemafunc.EqualJSON(resultModel.Pipeline.String(), tc.expectedTFModel.Pipeline.String(), "test stream processor schema"))
				var statsResult any
				err := json.Unmarshal([]byte(resultModel.Stats.ValueString()), &statsResult)
				if err != nil {
					t.Fatal(err)
				}
				assert.Len(t, sdkModel.Stats, 14)
				assert.Len(t, statsResult, 14)
			} else {
				assert.Equal(t, tc.expectedTFModel, resultModel)
			}
		})
	}
}
