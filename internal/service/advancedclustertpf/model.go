package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

func NewTFModel(ctx context.Context, apiResp *admin.ClusterDescription20240805) (*TFModel, diag.Diagnostics) {
	return &TFModel{}, nil
}

func NewAtlasReq(ctx context.Context, plan *TFModel) (*admin.ClusterDescription20240805, diag.Diagnostics) {
	return &admin.ClusterDescription20240805{}, nil
}
