package streamprocessor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

// GetWorkspaceOrInstanceName returns the workspace name from workspace_name or instance_name field. Assumes exactly one of the two is set.
func GetWorkspaceOrInstanceName(workspaceName, instanceName types.String) string {
	if !workspaceName.IsNull() && !workspaceName.IsUnknown() {
		return workspaceName.ValueString()
	}
	return instanceName.ValueString()
}

func NewStreamProcessorReq(ctx context.Context, plan *TFStreamProcessorRSModel) (*admin.StreamsProcessor, diag.Diagnostics) {
	pipeline, diags := convertPipelineToSdk(plan.Pipeline.ValueString())
	if diags != nil {
		return nil, diags
	}
	streamProcessor := &admin.StreamsProcessor{
		Name:     plan.ProcessorName.ValueStringPointer(),
		Pipeline: &pipeline,
	}

	if !plan.FailoverEnabled.IsNull() && !plan.FailoverEnabled.IsUnknown() {
		streamProcessor.FailoverEnabled = plan.FailoverEnabled.ValueBoolPointer()
	}

	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		optionsModel := &TFOptionsModel{}
		if diags := plan.Options.As(ctx, optionsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		dlq, diags := newDlqReq(ctx, optionsModel.Dlq)
		if diags.HasError() {
			return nil, diags
		}
		autoscaling, diags := newAutoscalingReq(ctx, optionsModel.Autoscaling)
		if diags.HasError() {
			return nil, diags
		}
		streamProcessor.Options = &admin.StreamsOptions{
			Dlq:         dlq,
			Autoscaling: autoscaling,
		}
	}

	if !plan.Tier.IsNull() && !plan.Tier.IsUnknown() {
		streamProcessor.Tier = plan.Tier.ValueStringPointer()
	}

	return streamProcessor, nil
}

func NewStreamProcessorUpdateReq(ctx context.Context, plan, state *TFStreamProcessorRSModel) (*admin.UpdateStreamProcessorApiParams, diag.Diagnostics) {
	pipeline, diags := convertPipelineToSdk(plan.Pipeline.ValueString())
	if diags != nil {
		return nil, diags
	}

	workspaceOrInstanceName := GetWorkspaceOrInstanceName(plan.WorkspaceName, plan.InstanceName)

	streamProcessorAPIParams := &admin.UpdateStreamProcessorApiParams{
		GroupId:       plan.ProjectID.ValueString(),
		TenantName:    workspaceOrInstanceName,
		ProcessorName: plan.ProcessorName.ValueString(),
		StreamsModifyStreamProcessor: &admin.StreamsModifyStreamProcessor{
			Name:     plan.ProcessorName.ValueStringPointer(),
			Pipeline: &pipeline,
		},
	}

	if !plan.FailoverEnabled.IsNull() && !plan.FailoverEnabled.IsUnknown() {
		streamProcessorAPIParams.StreamsModifyStreamProcessor.FailoverEnabled = plan.FailoverEnabled.ValueBoolPointer()
	}

	// Resolve the autoscaling operation with PATCH tri-state semantics.
	autoscaling, diags := resolveAutoscalingForUpdate(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	if !plan.Options.IsNull() && !plan.Options.IsUnknown() {
		optionsModel := &TFOptionsModel{}
		if diags := plan.Options.As(ctx, optionsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		dlq, diags := newDlqReq(ctx, optionsModel.Dlq)
		if diags.HasError() {
			return nil, diags
		}
		streamProcessorAPIParams.StreamsModifyStreamProcessor.Options = &admin.StreamsModifyStreamProcessorOptions{
			Dlq:         dlq,
			Autoscaling: autoscaling,
		}
	} else if autoscaling != nil {
		// The whole options block was removed but autoscaling still needs an explicit
		// disable. Dlq is left nil (omitted => preserved by the API).
		streamProcessorAPIParams.StreamsModifyStreamProcessor.Options = &admin.StreamsModifyStreamProcessorOptions{
			Autoscaling: autoscaling,
		}
	}

	// Baseline tier is settable on the PATCH body; when autoscaling is enabled the API treats
	// it as the initial tier only and reports the running tier via effectiveTier.
	if !plan.Tier.IsNull() && !plan.Tier.IsUnknown() {
		streamProcessorAPIParams.StreamsModifyStreamProcessor.Tier = plan.Tier.ValueStringPointer()
	}

	return streamProcessorAPIParams, nil
}

func NewStreamProcessorWithStats(ctx context.Context, projectID, instanceName, workspaceName string, apiResp *admin.StreamsProcessorWithStats, timeout *timeouts.Value, deleteOnCreateTimeout, failoverEnabled *types.Bool) (*TFStreamProcessorRSModel, diag.Diagnostics) {
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
		Options:       *optionsTF,
		Pipeline:      pipelineTF,
		ProcessorID:   types.StringPointerValue(&apiResp.Id),
		ProcessorName: types.StringPointerValue(&apiResp.Name),
		ProjectID:     types.StringPointerValue(&projectID),
		State:         types.StringPointerValue(&apiResp.State),
		Stats:         statsTF,
		Tier:          types.StringPointerValue(apiResp.Tier),
		EffectiveTier: types.StringValue(apiResp.EffectiveTier),
	}

	if workspaceName != "" {
		tfModel.WorkspaceName = types.StringValue(workspaceName)
		tfModel.InstanceName = types.StringNull()
	} else {
		// Default to instance_name for backward compatibility
		tfModel.InstanceName = types.StringValue(instanceName)
		tfModel.WorkspaceName = types.StringNull()
	}

	if timeout != nil {
		tfModel.Timeouts = *timeout
	}
	if deleteOnCreateTimeout != nil {
		tfModel.DeleteOnCreateTimeout = *deleteOnCreateTimeout
	}
	if failoverEnabled != nil {
		tfModel.FailoverEnabled = *failoverEnabled
	}
	return tfModel, nil
}

func NewTFStreamprocessorDSModel(ctx context.Context, projectID, instanceName, workspaceName string, apiResp *admin.StreamsProcessorWithStats) (*TFStreamProcessorDSModel, diag.Diagnostics) {
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
		ID:              types.StringPointerValue(&apiResp.Id),
		Options:         *optionsTF,
		Pipeline:        types.StringValue(pipelineTF.ValueString()),
		ProcessorName:   types.StringPointerValue(&apiResp.Name),
		ProjectID:       types.StringPointerValue(&projectID),
		State:           types.StringPointerValue(&apiResp.State),
		Stats:           statsTF,
		Tier:            types.StringPointerValue(apiResp.Tier),
		EffectiveTier:   types.StringValue(apiResp.EffectiveTier),
		FailoverEnabled: types.BoolValue(apiResp.GetFailoverEnabled()),
	}

	if workspaceName != "" {
		tfModel.WorkspaceName = types.StringValue(workspaceName)
		tfModel.InstanceName = types.StringNull()
	} else {
		// Default to instance_name for backward compatibility
		tfModel.InstanceName = types.StringValue(instanceName)
		tfModel.WorkspaceName = types.StringNull()
	}
	return tfModel, nil
}

func newDlqReq(ctx context.Context, dlq types.Object) (*admin.StreamsDLQ, diag.Diagnostics) {
	if dlq.IsNull() || dlq.IsUnknown() {
		return nil, nil
	}
	dlqModel := &TFDlqModel{}
	if diags := dlq.As(ctx, dlqModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsDLQ{
		Coll:           dlqModel.Coll.ValueStringPointer(),
		ConnectionName: dlqModel.ConnectionName.ValueStringPointer(),
		Db:             dlqModel.DB.ValueStringPointer(),
	}, nil
}

func ConvertOptionsToTF(ctx context.Context, options *admin.StreamsOptions) (*types.Object, diag.Diagnostics) {
	if options == nil || (!options.HasDlq() && options.Autoscaling == nil) {
		return new(types.ObjectNull(OptionsObjectType.AttributeTypes())), nil
	}
	dlqTF, diags := convertDlqToTF(ctx, options.Dlq)
	if diags.HasError() {
		return nil, diags
	}
	autoscalingTF, diags := convertAutoscalingToTF(ctx, options.Autoscaling)
	if diags.HasError() {
		return nil, diags
	}
	optionsTF := &TFOptionsModel{
		Dlq:         *dlqTF,
		Autoscaling: autoscalingTF,
	}
	optionsObject, diags := types.ObjectValueFrom(ctx, OptionsObjectType.AttributeTypes(), optionsTF)
	if diags.HasError() {
		return nil, diags
	}
	return &optionsObject, nil
}

func convertDlqToTF(ctx context.Context, dlq *admin.StreamsDLQ) (*types.Object, diag.Diagnostics) {
	if dlq == nil {
		return new(types.ObjectNull(DlqObjectType.AttributeTypes())), nil
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
	workspaceName := streamProcessorsConfig.WorkspaceName.ValueString()
	for i := range sdkResults {
		processorModel, diags := NewTFStreamprocessorDSModel(ctx, projectID, instanceName, workspaceName, &sdkResults[i])
		if diags.HasError() {
			return nil, diags
		}
		results[i] = *processorModel
	}
	return &TFStreamProcessorsDSModel{
		ProjectID:     streamProcessorsConfig.ProjectID,
		InstanceName:  streamProcessorsConfig.InstanceName,
		WorkspaceName: streamProcessorsConfig.WorkspaceName,
		Results:       results,
	}, nil
}

func convertPipelineToSdk(pipeline string) ([]any, diag.Diagnostics) {
	// Each stage is kept as raw JSON instead of being decoded into map[string]any: Go maps are
	// unordered and encoding/json sorts their keys on marshal, which would send Atlas a pipeline
	// whose subdocument fields are alphabetized rather than in the order the practitioner wrote.
	// MongoDB documents are ordered, so that silently changes the meaning of sort specifications
	// and of equality comparisons against document literals. json.Marshal emits json.RawMessage
	// verbatim, including when boxed in []any, so the authored order is preserved end to end.
	var stages []json.RawMessage
	if err := json.Unmarshal([]byte(pipeline), &stages); err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("failed to unmarshal pipeline", err.Error())}
	}
	rawStages := make([]any, len(stages))
	for i := range stages {
		rawStages[i] = stages[i]
	}
	return rawStages, nil
}
