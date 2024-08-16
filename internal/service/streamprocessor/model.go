package streamprocessor

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

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
	totalCount := len(sdkResults)
	return &TFStreamProcessorsDSModel{
		ProjectID:    streamProcessorsConfig.ProjectID,
		InstanceName: streamProcessorsConfig.InstanceName,
		Results:      results,
		TotalCount:   types.Int64PointerValue(conversion.IntPtrToInt64Ptr(&totalCount)),
	}, nil
}
