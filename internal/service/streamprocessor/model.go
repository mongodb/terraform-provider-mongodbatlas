package streamprocessor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func NewTFStreamProcessor(ctx context.Context, apiResp *admin.StreamsProcessor) (*TFStreamProcessorRSModel, diag.Diagnostics) {
	return &TFStreamProcessorRSModel{}, nil
}

func NewStreamProcessorReq(ctx context.Context, plan *TFStreamProcessorRSModel) (*admin.StreamsProcessor, diag.Diagnostics) {
	return &admin.StreamsProcessor{}, nil
}

func NewStreamProcessorWithStats(ctx context.Context, projectID, instanceName string, apiResp *admin.StreamsProcessorWithStats) (*TFStreamProcessorRSModel, diag.Diagnostics) {
	tfModel := &TFStreamProcessorRSModel{
		InstanceName:      types.StringValue(instanceName),
		ProcessorName:     types.StringPointerValue(&apiResp.Name),
		ProjectID:         types.StringValue(projectID),
		State:             types.StringPointerValue(&apiResp.State),
		ProcessorID:       types.StringNull(),
		ChangeStreamToken: types.StringNull(),
	}
	return tfModel, nil
}

func NewTFStreamprocessorDSModel(ctx context.Context, projectID, instanceName string, apiResp *admin.StreamsProcessorWithStats) (*TFStreamProcessorDSModel, diag.Diagnostics) {
	if apiResp == nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("streamProcessor API response is nil", "")}
	}
	pipelineTF, diags := convertPipelineToTF(apiResp.GetPipeline())
	if diags.HasError() {
		return nil, diags
	}
	statsTF, diags := convertStatsToTF(apiResp.GetStats())
	if diags.HasError() {
		return nil, diags
	}
	tfModel := &TFStreamProcessorDSModel{
		ID:            types.StringPointerValue(&apiResp.Id),
		InstanceName:  types.StringPointerValue(&instanceName),
		Pipeline:      pipelineTF,
		ProcessorName: types.StringPointerValue(&apiResp.Name),
		ProjectID:     types.StringPointerValue(&projectID),
		State:         types.StringPointerValue(&apiResp.State),
		Stats:         statsTF,
	}
	return tfModel, nil
}

func convertPipelineToTF(pipeline []any) (types.String, diag.Diagnostics) {
	pipelineJSON, err := json.Marshal(pipeline)
	if err != nil {
		return types.StringValue(""), diag.Diagnostics{diag.NewErrorDiagnostic("failed to marshal pipeline", err.Error())}
	}
	return types.StringValue(string(pipelineJSON)), nil
}

func convertStatsToTF(stats any) (types.String, diag.Diagnostics) {
	if stats == nil {
		return types.StringValue("{}"), nil
	}
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return types.StringValue(""), diag.Diagnostics{diag.NewErrorDiagnostic("failed to marshal stats", err.Error())}
	}
	return types.StringValue(string(statsJSON)), nil
}
