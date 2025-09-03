package streamprocessor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
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

func NewStreamProcessorUpdateReq(ctx context.Context, plan *TFStreamProcessorRSModel) (*admin.ModifyStreamProcessorApiParams, diag.Diagnostics) {
	pipeline, diags := convertPipelineToSdk(plan.Pipeline.ValueString())
	if diags != nil {
		return nil, diags
	}

	streamProcessorAPIParams := &admin.ModifyStreamProcessorApiParams{
		GroupId:       plan.ProjectID.ValueString(),
		TenantName:    plan.InstanceName.ValueString(),
		ProcessorName: plan.ProcessorName.ValueString(),
		StreamsModifyStreamProcessor: &admin.StreamsModifyStreamProcessor{
			Name:     plan.ProcessorName.ValueStringPointer(),
			Pipeline: &pipeline,
		},
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
		streamProcessorAPIParams.StreamsModifyStreamProcessor.Options = &admin.StreamsModifyStreamProcessorOptions{
			Dlq: &admin.StreamsDLQ{
				Coll:           dlqModel.Coll.ValueStringPointer(),
				ConnectionName: dlqModel.ConnectionName.ValueStringPointer(),
				Db:             dlqModel.DB.ValueStringPointer(),
			},
		}
	}

	return streamProcessorAPIParams, nil
}

func NewStreamProcessorWithStats(ctx context.Context, projectID, instanceName string, apiResp *admin.StreamsProcessorWithStats) (*TFStreamProcessorRSModel, diag.Diagnostics) {
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
	optionsTF, diags := ConvertOptionsToTF(ctx, apiResp.Options)
	if diags.HasError() {
		return nil, diags
	}
	tfModel := &TFStreamProcessorRSModel{
		InstanceName:  types.StringPointerValue(&instanceName),
		Options:       *optionsTF,
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
	optionsTF, diags := ConvertOptionsToTF(ctx, apiResp.Options)
	if diags.HasError() {
		return nil, diags
	}
	tfModel := &TFStreamProcessorDSModel{
		ID:            types.StringPointerValue(&apiResp.Id),
		InstanceName:  types.StringPointerValue(&instanceName),
		Options:       *optionsTF,
		Pipeline:      types.StringValue(pipelineTF.ValueString()),
		ProcessorName: types.StringPointerValue(&apiResp.Name),
		ProjectID:     types.StringPointerValue(&projectID),
		State:         types.StringPointerValue(&apiResp.State),
		Stats:         statsTF,
	}
	return tfModel, nil
}

func ConvertOptionsToTF(ctx context.Context, options *admin.StreamsOptions) (*types.Object, diag.Diagnostics) {
	if options == nil || !options.HasDlq() {
		optionsTF := types.ObjectNull(OptionsObjectType.AttributeTypes())
		return &optionsTF, nil
	}
	dlqTF, diags := convertDlqToTF(ctx, options.Dlq)
	if diags.HasError() {
		return nil, diags
	}
	optionsTF := &TFOptionsModel{
		Dlq: *dlqTF,
	}
	optionsObject, diags := types.ObjectValueFrom(ctx, OptionsObjectType.AttributeTypes(), optionsTF)
	if diags.HasError() {
		return nil, diags
	}
	return &optionsObject, nil
}

func convertDlqToTF(ctx context.Context, dlq *admin.StreamsDLQ) (*types.Object, diag.Diagnostics) {
	if dlq == nil {
		dlqTF := types.ObjectNull(DlqObjectType.AttributeTypes())
		return &dlqTF, nil
	}
	dlqModel := TFDlqModel{
		Coll:           types.StringPointerValue(dlq.Coll),
		ConnectionName: types.StringPointerValue(dlq.ConnectionName),
		DB:             types.StringPointerValue(dlq.Db),
	}
	dlqObject, diags := types.ObjectValueFrom(ctx, DlqObjectType.AttributeTypes(), dlqModel)
	if diags.HasError() {
		return nil, diags
	}
	return &dlqObject, nil
}
func convertPipelineToTF(pipeline []any) (jsontypes.Normalized, diag.Diagnostics) {
	pipelineJSON, err := json.Marshal(pipeline)
	if err != nil {
		return jsontypes.NewNormalizedValue(""), diag.Diagnostics{diag.NewErrorDiagnostic("failed to marshal pipeline", err.Error())}
	}
	return jsontypes.NewNormalizedValue(string(pipelineJSON)), nil
}

func convertStatsToTF(stats any) (types.String, diag.Diagnostics) {
	if stats == nil {
		return types.StringNull(), nil
	}
	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return types.StringValue(""), diag.Diagnostics{diag.NewErrorDiagnostic("failed to marshal stats", err.Error())}
	}
	return types.StringValue(string(statsJSON)), nil
}

func NewTFStreamProcessors(ctx context.Context,
	streamProcessorsConfig *TFStreamProcessorsDSModel,
	sdkResults []admin.StreamsProcessorWithStats) (*TFStreamProcessorsDSModel, diag.Diagnostics) {
	results := make([]TFStreamProcessorDSModel, len(sdkResults))
	projectID := streamProcessorsConfig.ProjectID.ValueString()
	instanceName := streamProcessorsConfig.InstanceName.ValueString()
	for i := range sdkResults {
		processorModel, diags := NewTFStreamprocessorDSModel(ctx, projectID, instanceName, &sdkResults[i])
		if diags.HasError() {
			return nil, diags
		}
		results[i] = *processorModel
	}
	return &TFStreamProcessorsDSModel{
		ProjectID:    streamProcessorsConfig.ProjectID,
		InstanceName: streamProcessorsConfig.InstanceName,
		Results:      results,
	}, nil
}

func convertPipelineToSdk(pipeline string) ([]any, diag.Diagnostics) {
	var pipelineSliceOfMaps []any
	err := json.Unmarshal([]byte(pipeline), &pipelineSliceOfMaps)
	if err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("failed to unmarshal pipeline", err.Error())}
	}
	return pipelineSliceOfMaps, nil
}
