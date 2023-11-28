package streaminstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func newStreamInstanceCreateReq(ctx context.Context, plan *tfStreamInstanceRSModel) (*admin.StreamsTenant, diag.Diagnostics) {
	dataProcessRegion := &tfInstanceProcessRegionSpecModel{}
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

func newStreamInstanceUpdateReq(ctx context.Context, plan *tfStreamInstanceRSModel) (*admin.StreamsDataProcessRegion, diag.Diagnostics) {
	dataProcessRegion := &tfInstanceProcessRegionSpecModel{}
	diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsDataProcessRegion{
		CloudProvider: dataProcessRegion.CloudProvider.ValueString(),
		Region:        dataProcessRegion.Region.ValueString(),
	}, nil
}

func newTFStreamInstance(ctx context.Context, apiResp *admin.StreamsTenant) (*tfStreamInstanceRSModel, diag.Diagnostics) {
	hostnames, diags := types.ListValueFrom(ctx, types.StringType, apiResp.Hostnames)
	// TODO check dataRegionIsDefined
	dataProcessRegion, diagsProcessRegion := types.ObjectValueFrom(ctx, ProcessRegionObjectType.AttrTypes, tfInstanceProcessRegionSpecModel{
		CloudProvider: types.StringValue(apiResp.DataProcessRegion.CloudProvider),
		Region:        types.StringValue(apiResp.DataProcessRegion.Region),
	})
	diags.Append(diagsProcessRegion...)
	if diags.HasError() {
		return nil, diags
	}

	return &tfStreamInstanceRSModel{
		ID:                types.StringPointerValue(apiResp.Id),
		InstanceName:      types.StringPointerValue(apiResp.Name),
		ProjectID:         types.StringPointerValue(apiResp.GroupId),
		DataProcessRegion: dataProcessRegion,
		Hostnames:         hostnames,
	}, nil
}
