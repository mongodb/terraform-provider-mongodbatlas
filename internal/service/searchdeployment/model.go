package searchdeployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func NewSearchDeploymentReq(ctx context.Context, searchDeploymentPlan *TFSearchDeploymentRSModel) admin.ApiSearchDeploymentRequest {
	var specs []TFSearchNodeSpecModel
	searchDeploymentPlan.Specs.ElementsAs(ctx, &specs, true)

	resultSpecs := make([]admin.ApiSearchDeploymentSpec, len(specs))
	for i, spec := range specs {
		resultSpecs[i] = admin.ApiSearchDeploymentSpec{
			InstanceSize: spec.InstanceSize.ValueString(),
			NodeCount:    int(spec.NodeCount.ValueInt64()),
		}
	}
	return admin.ApiSearchDeploymentRequest{
		Specs: resultSpecs,
	}
}

func NewTFSearchDeployment(ctx context.Context, clusterName string, deployResp *admin.ApiSearchDeploymentResponse, timeout *timeouts.Value, allowMultipleSpecs bool) (*TFSearchDeploymentRSModel, diag.Diagnostics) {
	result := TFSearchDeploymentRSModel{
		ID:                       types.StringPointerValue(deployResp.Id),
		ClusterName:              types.StringValue(clusterName),
		ProjectID:                types.StringPointerValue(deployResp.GroupId),
		StateName:                types.StringPointerValue(deployResp.StateName),
		EncryptionAtRestProvider: types.StringPointerValue(deployResp.EncryptionAtRestProvider),
	}

	if timeout != nil {
		result.Timeouts = *timeout
	}

	specs := deployResp.GetSpecs()
	if !allowMultipleSpecs && len(specs) > 1 {
		specs = specs[:1]
	}
	specsList, diagnostics := types.ListValueFrom(ctx, SpecObjectType, newTFSpecsModel(specs))
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	result.Specs = specsList
	return &result, diagnostics
}

func newTFSpecsModel(specs []admin.ApiSearchDeploymentSpec) []TFSearchNodeSpecModel {
	result := make([]TFSearchNodeSpecModel, len(specs))
	for i, v := range specs {
		result[i] = TFSearchNodeSpecModel{
			InstanceSize: types.StringValue(v.InstanceSize),
			NodeCount:    types.Int64Value(int64(v.NodeCount)),
		}
	}

	return result
}
