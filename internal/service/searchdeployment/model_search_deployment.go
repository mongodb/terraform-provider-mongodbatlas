package searchdeployment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
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

func NewTFSearchDeployment(ctx context.Context, clusterName string, deployResp *admin.ApiSearchDeploymentResponse, timeout *timeouts.Value, forceSpecLength int) (*TFSearchDeploymentRSModel, diag.Diagnostics) {
	result := TFSearchDeploymentRSModel{
		ID:          types.StringPointerValue(deployResp.Id),
		ClusterName: types.StringValue(clusterName),
		ProjectID:   types.StringPointerValue(deployResp.GroupId),
		StateName:   types.StringPointerValue(deployResp.StateName),
	}

	if timeout != nil {
		result.Timeouts = *timeout
	}

	allSpecs := deployResp.GetSpecs()
	fullLen := len(allSpecs)
	if forceSpecLength > 0 {
		allSpecs = allSpecs[:forceSpecLength]
	}
	specsList, diagnostics := types.ListValueFrom(ctx, SpecObjectType, newTFSpecsModel(allSpecs))
	if diagnostics.HasError() {
		return nil, diagnostics
	}
	if forceSpecLength > 0 && fullLen != forceSpecLength {
		diagnostics.Append(diag.NewWarningDiagnostic("spec length missmatch", fmt.Sprintf("your configuration has %d specs, but the actual deployment has %d specs, please update your configuration", forceSpecLength, fullLen)))
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
