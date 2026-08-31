package streamprocessor

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"

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

// IsAliasOnlyTransition reports whether a plan preserves the effective workspace
// and changes no processor configuration other than its alias representation.
func IsAliasOnlyTransition(plan, state *TFStreamProcessorRSModel) bool {
	if plan.WorkspaceName.IsUnknown() || plan.InstanceName.IsUnknown() || state.WorkspaceName.IsUnknown() || state.InstanceName.IsUnknown() {
		return false
	}
	if plan.WorkspaceName.Equal(state.WorkspaceName) && plan.InstanceName.Equal(state.InstanceName) {
		return false
	}
	if GetWorkspaceOrInstanceName(plan.WorkspaceName, plan.InstanceName) != GetWorkspaceOrInstanceName(state.WorkspaceName, state.InstanceName) {
		return false
	}

	return plan.Options.Equal(state.Options) &&
		plan.Pipeline.Equal(state.Pipeline) &&
		plan.ProcessorName.Equal(state.ProcessorName) &&
		plan.ProjectID.Equal(state.ProjectID) &&
		plan.FailoverEnabled.Equal(state.FailoverEnabled) &&
		equalOptionalComputedString(plan.State, state.State) &&
		equalOptionalComputedString(plan.Tier, state.Tier)
}

// equalOptionalComputedString treats an unknown plan value as unchanged because
// Terraform may defer Optional+Computed attributes during an update plan.
func equalOptionalComputedString(plan, state types.String) bool {
	return plan.IsUnknown() || plan.Equal(state)
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

// ResumeFromCheckpointFromOptions returns options.resume_from_checkpoint, or a null Bool when it is
// not set. The Atlas Admin API never returns this value, so it is carried over from configuration or
// prior state rather than read from the API response.
func ResumeFromCheckpointFromOptions(ctx context.Context, options *types.Object) (types.Bool, diag.Diagnostics) {
	if options == nil || options.IsNull() || options.IsUnknown() {
		return types.BoolNull(), nil
	}
	optionsModel := &TFOptionsModel{}
	if diags := options.As(ctx, optionsModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return types.BoolNull(), diags
	}
	return optionsModel.ResumeFromCheckpoint, nil
}

// checkpointIncompatibleChange reports whether the plan changes something the API refuses to modify
// while resuming from a checkpoint: the $source stage, or a window's type, interval, hopSize or
// allowedLateness, including adding or removing a window. Any other edit, for example to a $match or
// $emit stage, can resume from the existing checkpoint, so resume_from_checkpoint is not sent for it
// and a value left in the configuration cannot discard the checkpoint again.
//
// When the change cannot be determined it reports true, so that an explicitly configured
// resume_from_checkpoint is still sent rather than silently dropped. This is defensive: the plan
// pipeline is already valid JSON by the time it gets here, both from the jsontypes.Normalized
// attribute type and from convertPipelineToSdk above.
func checkpointIncompatibleChange(plan, state *TFStreamProcessorRSModel) bool {
	if state == nil {
		return true
	}
	planStages, ok := unmarshalPipelineStages(plan.Pipeline.ValueString())
	if !ok {
		return true
	}
	stateStages, ok := unmarshalPipelineStages(state.Pipeline.ValueString())
	if !ok {
		return true
	}
	if !reflect.DeepEqual(sourceStage(planStages), sourceStage(stateStages)) {
		return true
	}
	return !reflect.DeepEqual(windowStages(planStages), windowStages(stateStages))
}

func unmarshalPipelineStages(pipeline string) (stages []any, ok bool) {
	if pipeline == "" {
		return nil, false
	}
	if err := json.Unmarshal([]byte(pipeline), &stages); err != nil {
		return nil, false
	}
	return stages, true
}

// sourceStage returns the $source stage, always the first stage of a stream processor pipeline.
func sourceStage(stages []any) any {
	if len(stages) == 0 {
		return nil
	}
	return stages[0]
}

// windowStages projects the window stages down to the fields the API rejects a checkpoint resume
// for, so that changes to a window's inner pipeline are not mistaken for incompatible ones. Window
// operators are matched by name suffix to cover $tumblingWindow, $hoppingWindow and any later
// addition.
func windowStages(stages []any) []any {
	var windows []any
	for _, stage := range stages {
		obj, ok := stage.(map[string]any)
		if !ok {
			continue
		}
		for name, value := range obj {
			if !strings.HasSuffix(name, "Window") {
				continue
			}
			config, _ := value.(map[string]any)
			windows = append(windows, map[string]any{
				"name":            name,
				"interval":        config["interval"],
				"hopSize":         config["hopSize"],
				"allowedLateness": config["allowedLateness"],
			})
		}
	}
	return windows
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

	// Resolve the options PATCH operations with tri-state semantics.
	autoscaling, clearAutoscaling, diags := resolveAutoscalingForUpdate(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}
	dlq, clearDLQ, diags := resolveDlqForUpdate(ctx, plan, state)
	if diags.HasError() {
		return nil, diags
	}

	resumeFromCheckpoint, diags := ResumeFromCheckpointFromOptions(ctx, &plan.Options)
	if diags.HasError() {
		return nil, diags
	}
	// resume_from_checkpoint is only sent for the changes the API rejects while resuming from a
	// checkpoint, see checkpointIncompatibleChange. Omitting it otherwise means a value left in the
	// configuration cannot discard the checkpoint on an unrelated update. A configured value is sent
	// as-is, including true; when the attribute is not set nothing is sent, so the API applies its own
	// default rather than the provider assuming one.
	sendResumeFromCheckpoint := !resumeFromCheckpoint.IsNull() && !resumeFromCheckpoint.IsUnknown() && checkpointIncompatibleChange(plan, state)

	if autoscaling != nil || clearAutoscaling || dlq != nil || clearDLQ || sendResumeFromCheckpoint {
		options := &admin.StreamsModifyStreamProcessorOptions{
			Dlq:         dlq,
			Autoscaling: autoscaling,
		}
		if clearAutoscaling {
			options.SetAutoscalingNil()
		}
		if clearDLQ {
			// SPM uses an empty object (rather than null) as the DLQ clear signal.
			options.Dlq = &admin.StreamsDLQ{}
		}
		if sendResumeFromCheckpoint {
			options.ResumeFromCheckpoint = resumeFromCheckpoint.ValueBoolPointer()
		}
		streamProcessorAPIParams.StreamsModifyStreamProcessor.Options = options
	}

	// Baseline tier is settable on the PATCH body; when autoscaling is enabled the API treats
	// it as the initial tier only and reports the running tier via effectiveTier.
	if !plan.Tier.IsNull() && !plan.Tier.IsUnknown() {
		streamProcessorAPIParams.StreamsModifyStreamProcessor.Tier = plan.Tier.ValueStringPointer()
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

// ConvertOptionsToTFDS builds the data source options object, which has no resume_from_checkpoint.
func ConvertOptionsToTFDS(ctx context.Context, options *admin.StreamsOptions) (*types.Object, diag.Diagnostics) {
	if options == nil || (!options.HasDlq() && options.Autoscaling == nil) {
		return new(types.ObjectNull(DSOptionsObjectType.AttributeTypes())), nil
	}
	dlqTF, diags := convertDlqToTF(ctx, options.Dlq)
	if diags.HasError() {
		return nil, diags
	}
	autoscalingTF, diags := convertAutoscalingToTF(ctx, options.Autoscaling)
	if diags.HasError() {
		return nil, diags
	}
	optionsObject, diags := types.ObjectValueFrom(ctx, DSOptionsObjectType.AttributeTypes(), &TFDSOptionsModel{
		Dlq:         *dlqTF,
		Autoscaling: autoscalingTF,
	})
	if diags.HasError() {
		return nil, diags
	}
	return &optionsObject, nil
}

// ConvertOptionsToTF builds the resource options object. resumeFromCheckpoint is not returned by the
// API, so callers pass the value from configuration or prior state to preserve it.
func ConvertOptionsToTF(ctx context.Context, options *admin.StreamsOptions, resumeFromCheckpoint types.Bool) (*types.Object, diag.Diagnostics) {
	hasOptions := options != nil && (options.HasDlq() || options.Autoscaling != nil)
	if !hasOptions && resumeFromCheckpoint.IsNull() {
		return new(types.ObjectNull(OptionsObjectType.AttributeTypes())), nil
	}
	dlqTF := new(types.ObjectNull(DlqObjectType.AttributeTypes()))
	autoscalingTF := types.ObjectNull(AutoscalingObjectType.AttributeTypes())
	if hasOptions {
		var diags diag.Diagnostics
		dlqTF, diags = convertDlqToTF(ctx, options.Dlq)
		if diags.HasError() {
			return nil, diags
		}
		autoscalingTF, diags = convertAutoscalingToTF(ctx, options.Autoscaling)
		if diags.HasError() {
			return nil, diags
		}
	}
	optionsTF := &TFOptionsModel{
		Dlq:                  *dlqTF,
		Autoscaling:          autoscalingTF,
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
