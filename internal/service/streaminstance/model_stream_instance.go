package streaminstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func NewStreamInstanceCreateReq(ctx context.Context, plan *TFStreamInstanceRSModel) (*admin.StreamsTenant, diag.Diagnostics) {
	dataProcessRegion := &TFInstanceProcessRegionSpecModel{}
	diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsTenant{
		GroupId: plan.ProjectID.ValueStringPointer(),
		Name:    plan.InstanceName.ValueStringPointer(),
		DataProcessRegion: &admin.StreamsDataProcessRegion{
			CloudProvider: dataProcessRegion.CloudProvider.ValueString(),
			Region:        dataProcessRegion.Region.ValueString(),
		},
	}, nil
}

func NewStreamInstanceUpdateReq(ctx context.Context, plan *TFStreamInstanceRSModel) (*admin.StreamsDataProcessRegion, diag.Diagnostics) {
	dataProcessRegion := &TFInstanceProcessRegionSpecModel{}
	diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsDataProcessRegion{
		CloudProvider: dataProcessRegion.CloudProvider.ValueString(),
		Region:        dataProcessRegion.Region.ValueString(),
	}, nil
}

func NewTFStreamInstance(ctx context.Context, apiResp *admin.StreamsTenant) (*TFStreamInstanceRSModel, diag.Diagnostics) {
	hostnames, diags := types.ListValueFrom(ctx, types.StringType, apiResp.Hostnames)
	// TODO check dataRegionIsDefined
	dataProcessRegion, diagsProcessRegion := types.ObjectValueFrom(ctx, ProcessRegionObjectType.AttrTypes, TFInstanceProcessRegionSpecModel{
		CloudProvider: types.StringValue(apiResp.DataProcessRegion.CloudProvider),
		Region:        types.StringValue(apiResp.DataProcessRegion.Region),
	})
	diags.Append(diagsProcessRegion...)
	if diags.HasError() {
		return nil, diags
	}

	return &TFStreamInstanceRSModel{
		ID:                types.StringPointerValue(apiResp.Id),
		InstanceName:      types.StringPointerValue(apiResp.Name),
		ProjectID:         types.StringPointerValue(apiResp.GroupId),
		DataProcessRegion: dataProcessRegion,
		Hostnames:         hostnames,
	}, nil
}
