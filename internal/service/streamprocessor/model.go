package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func NewTFStreamProcessor(ctx context.Context, apiResp *admin.StreamsProcessor) (*TFStreamProcessorRSModel, diag.Diagnostics) {
	return &TFStreamProcessorRSModel{}, nil
}

func NewStreamProcessorReq(ctx context.Context, plan *TFStreamProcessorRSModel) (*admin.StreamsProcessor, diag.Diagnostics) {
	return &admin.StreamsProcessor{}, nil
}

func NewStreamProcessorWithStats(ctx context.Context, apiResp *admin.StreamsProcessorWithStats) (*TFStreamProcessorRSModel, diag.Diagnostics) {
	return &TFStreamProcessorRSModel{}, nil
}
