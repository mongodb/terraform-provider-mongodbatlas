package streamprocessor

import (
	"context"

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
