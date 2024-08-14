package streamprocessor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func NewStreamProcessorReq(ctx context.Context, plan *TFStreamProcessorRSModel) (*admin.StreamsProcessor, diag.Diagnostics) {
	return &admin.StreamsProcessor{}, nil
}

func NewStreamProcessorWithStats(ctx context.Context, projectID, instanceName string, apiResp *admin.StreamsProcessorWithStats, stateOptions types.Object) (*TFStreamProcessorRSModel, diag.Diagnostics) {
	if apiResp == nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("streamProcessor API response is nil", "")}
	}
	pipelineTF, diags := convertPipelineToTF(apiResp.GetPipeline())
	if diags.HasError() {
		return nil, diags
	}
	changeStreamToken, diags := extractChangeStreamTokenFromStats(apiResp.GetStats())
	if diags.HasError() {
		return nil, diags
	}
	tfModel := &TFStreamProcessorRSModel{
		InstanceName:      types.StringPointerValue(&instanceName),
		Options:           stateOptions,
		Pipeline:          pipelineTF,
		ProcessorID:       types.StringPointerValue(&apiResp.Id),
		ProcessorName:     types.StringPointerValue(&apiResp.Name),
		ProjectID:         types.StringPointerValue(&projectID),
		State:             types.StringPointerValue(&apiResp.State),
		ChangeStreamToken: changeStreamToken,
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

func extractChangeStreamTokenFromStats(stats any) (types.String, diag.Diagnostics) {
	if stats == nil {
		return types.StringValue("{}"), nil
	}
	var statsMap map[string]interface{}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return types.StringValue(""), diag.Diagnostics{diag.NewErrorDiagnostic("failed to marshal stats", err.Error())}
	}

	err = json.Unmarshal(statsJSON, &statsMap)
	if err != nil {
		return types.StringValue(""), diag.Diagnostics{diag.NewErrorDiagnostic("failed to unmarshal stats", err.Error())}
	}

	changeStreamToken := statsMap["data"].(map[string]interface{})["changeStreamToken"]

	return types.StringValue(changeStreamToken.(string)), nil
}
