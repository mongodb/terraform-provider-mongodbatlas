package streamprocessor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"go.mongodb.org/atlas-sdk/v20250312023/admin"
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

	// options.resume_from_checkpoint is only honored by the modify endpoint, so it is ignored here.
	dlq, diags := newDlqReq(ctx, &plan.Options)
	if diags.HasError() {
		return nil, diags
	}
	if dlq != nil {
		streamProcessor.Options = &admin.StreamsOptions{Dlq: dlq}
	}

	return streamProcessor, nil
}

// parseOptions returns the options model, or nil when options is not set.
func parseOptions(ctx context.Context, options *types.Object) (*TFOptionsModel, diag.Diagnostics) {
	if options == nil || options.IsNull() || options.IsUnknown() {
		return nil, nil
	}
	optionsModel := &TFOptionsModel{}
	if diags := options.As(ctx, optionsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return optionsModel, nil
}

// newDlqReq converts options.dlq to the SDK type, returning nil when options or dlq is not set.
func newDlqReq(ctx context.Context, options *types.Object) (*admin.StreamsDLQ, diag.Diagnostics) {
	optionsModel, diags := parseOptions(ctx, options)
	if diags.HasError() || optionsModel == nil {
		return nil, diags
	}
	if optionsModel.Dlq.IsNull() || optionsModel.Dlq.IsUnknown() {
		return nil, nil
	}
	dlqModel := &TFDlqModel{}
	if diags := optionsModel.Dlq.As(ctx, dlqModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsDLQ{
		Coll:           dlqModel.Coll.ValueStringPointer(),
		ConnectionName: dlqModel.ConnectionName.ValueStringPointer(),
		Db:             dlqModel.DB.ValueStringPointer(),
	}, nil
}

// ResumeFromCheckpointFromOptions returns options.resume_from_checkpoint, or a null Bool when it is
// not set. The Atlas Admin API never returns this value, so it is carried over from configuration or
// prior state rather than read from the API response.
func ResumeFromCheckpointFromOptions(ctx context.Context, options *types.Object) (types.Bool, diag.Diagnostics) {
	optionsModel, diags := parseOptions(ctx, options)
	if diags.HasError() || optionsModel == nil {
		return types.BoolNull(), diags
	}
	return optionsModel.ResumeFromCheckpoint, diags
}

// pipelineChanged reports whether the plan changes the pipeline, comparing semantically so that
// formatting-only differences in the JSON string are not treated as a change. Without prior state
// there is nothing to compare, so it reports no change and the caller omits resume_from_checkpoint,
// the safer default for a flag that discards checkpoints.
func pipelineChanged(plan, state *TFStreamProcessorRSModel) bool {
	if state == nil {
		return false
	}
	return !schemafunc.EqualJSON(state.Pipeline.ValueString(), plan.Pipeline.ValueString(), "stream processor pipeline")
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

	apiOptions := &admin.StreamsModifyStreamProcessorOptions{}
	optionsSet := false

	dlq, diags := newDlqReq(ctx, &plan.Options)
	if diags.HasError() {
		return nil, diags
	}
	if dlq != nil {
		apiOptions.Dlq = dlq
		optionsSet = true
	}

	resumeFromCheckpoint, diags := ResumeFromCheckpointFromOptions(ctx, &plan.Options)
	if diags.HasError() {
		return nil, diags
	}
	// resume_from_checkpoint is only sent when the pipeline changes, the only updates the API rejects
	// for a checkpoint incompatible with the new pipeline. Omitting it for other updates, for example a
	// tier or dlq change, avoids discarding the checkpoint when the value is left in the
	// configuration. The field is never sent as true either, omission lets the API apply its default.
	if !resumeFromCheckpoint.IsNull() && !resumeFromCheckpoint.IsUnknown() && pipelineChanged(plan, state) {
		apiOptions.ResumeFromCheckpoint = resumeFromCheckpoint.ValueBoolPointer()
		optionsSet = true
	}

	if optionsSet {
		streamProcessorAPIParams.StreamsModifyStreamProcessor.Options = apiOptions
	}

	return streamProcessorAPIParams, nil
}

// NewStreamProcessorWithStats builds the resource model from the API response. configOptions is the
// options object from the plan or prior state, needed to preserve resume_from_checkpoint which the
// API does not return.
func NewStreamProcessorWithStats(ctx context.Context, projectID, instanceName, workspaceName string, apiResp *admin.StreamsProcessorWithStats, timeout *timeouts.Value, deleteOnCreateTimeout, failoverEnabled *types.Bool, configOptions *types.Object) (*TFStreamProcessorRSModel, diag.Diagnostics) {
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
	resumeFromCheckpoint, diags := ResumeFromCheckpointFromOptions(ctx, configOptions)
	if diags.HasError() {
		return nil, diags
	}
	optionsTF, diags := ConvertOptionsToTF(ctx, apiResp.Options, resumeFromCheckpoint)
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
	optionsTF, diags := ConvertOptionsToTFDS(ctx, apiResp.Options)
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

// ConvertOptionsToTFDS builds the data source options object, which only holds dlq.
func ConvertOptionsToTFDS(ctx context.Context, options *admin.StreamsOptions) (*types.Object, diag.Diagnostics) {
	if options == nil || !options.HasDlq() {
		return new(types.ObjectNull(DSOptionsObjectType.AttributeTypes())), nil
	}
	dlqTF, diags := convertDlqToTF(ctx, options.Dlq)
	if diags.HasError() {
		return nil, diags
	}
	optionsObject, diags := types.ObjectValueFrom(ctx, DSOptionsObjectType.AttributeTypes(), &TFDSOptionsModel{Dlq: *dlqTF})
	if diags.HasError() {
		return nil, diags
	}
	return &optionsObject, nil
}

// ConvertOptionsToTF builds the resource options object. resumeFromCheckpoint is not returned by the
// API, so callers pass the value from configuration or prior state to preserve it.
func ConvertOptionsToTF(ctx context.Context, options *admin.StreamsOptions, resumeFromCheckpoint types.Bool) (*types.Object, diag.Diagnostics) {
	hasDlq := options != nil && options.HasDlq()
	if !hasDlq && resumeFromCheckpoint.IsNull() {
		return new(types.ObjectNull(OptionsObjectType.AttributeTypes())), nil
	}
	dlqTF := new(types.ObjectNull(DlqObjectType.AttributeTypes()))
	if hasDlq {
		var diags diag.Diagnostics
		dlqTF, diags = convertDlqToTF(ctx, options.Dlq)
		if diags.HasError() {
			return nil, diags
		}
	}
	optionsTF := &TFOptionsModel{
		Dlq:                  *dlqTF,
		ResumeFromCheckpoint: resumeFromCheckpoint,
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
	var pipelineSliceOfMaps []any
	err := json.Unmarshal([]byte(pipeline), &pipelineSliceOfMaps)
	if err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("failed to unmarshal pipeline", err.Error())}
	}
	return pipelineSliceOfMaps, nil
}
