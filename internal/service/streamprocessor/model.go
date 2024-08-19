package streamprocessor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func NewStreamProcessorReq(ctx context.Context, plan *TFStreamProcessorRSModel) (*admin.StreamsProcessor, diag.Diagnostics) {
	pipeline, diags := convertPipelineToSdk(plan.Pipeline.ValueString())
	if diags != nil {
		return nil, diags
	}
	streamProcessor := &admin.StreamsProcessor{
		Name:     plan.ProcessorName.ValueStringPointer(),
		Pipeline: &pipeline,
	}

	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		optionsModel := &TFOptionsModel{}
		if diags := plan.Options.As(ctx, optionsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		dlqModel := &TFDlqModel{}
		if diags := optionsModel.Dlq.As(ctx, dlqModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamProcessor.Options = &admin.StreamsOptions{
			Dlq: &admin.StreamsDLQ{
				Coll:           dlqModel.Coll.ValueStringPointer(),
				ConnectionName: dlqModel.ConnectionName.ValueStringPointer(),
				Db:             dlqModel.DB.ValueStringPointer(),
			},
		}
	}

	return streamProcessor, nil
}

func NewStreamProcessorWithStats(ctx context.Context, projectID, instanceName string, apiResp *admin.StreamsProcessorWithStats, stateOptions types.Object) (*TFStreamProcessorRSModel, diag.Diagnostics) {
	if apiResp == nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("streamProcessor API response is nil", "")}
	}
	if stateOptions.IsUnknown() {
		stateOptions = types.ObjectNull(OptionsObjectType.AttrTypes)
	}
	pipelineTF, diags := convertPipelineToTF(apiResp.GetPipeline())
	if diags.HasError() {
		return nil, diags
	}
	statsTF, diags := convertStatsToTF(apiResp.GetStats())
	if diags.HasError() {
		return nil, diags
	}
	tfModel := &TFStreamProcessorRSModel{
		InstanceName:  types.StringPointerValue(&instanceName),
		Options:       stateOptions,
		Pipeline:      pipelineTF,
		ProcessorID:   types.StringPointerValue(&apiResp.Id),
		ProcessorName: types.StringPointerValue(&apiResp.Name),
		ProjectID:     types.StringPointerValue(&projectID),
		State:         types.StringPointerValue(&apiResp.State),
		Stats:         statsTF,
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

func convertPipelineToSdk(pipeline string) ([]any, diag.Diagnostics) {
	var pipelineSliceOfMaps []any
	err := json.Unmarshal([]byte(pipeline), &pipelineSliceOfMaps)
	if err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("failed to unmarshal pipeline", err.Error())}
	}
	return pipelineSliceOfMaps, nil
}
