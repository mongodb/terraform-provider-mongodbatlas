package streaminstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"
)

func NewStreamInstanceCreateReq(ctx context.Context, plan *TFStreamInstanceModel) (*admin.StreamsTenant, diag.Diagnostics) {
	dataProcessRegion := &TFInstanceProcessRegionSpecModel{}
	if diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{}); diags.HasError() {
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

func NewStreamInstanceUpdateReq(ctx context.Context, plan *TFStreamInstanceModel) (*admin.StreamsDataProcessRegion, diag.Diagnostics) {
	dataProcessRegion := &TFInstanceProcessRegionSpecModel{}
	if diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsDataProcessRegion{
		CloudProvider: dataProcessRegion.CloudProvider.ValueString(),
		Region:        dataProcessRegion.Region.ValueString(),
	}, nil
}

func NewTFStreamInstance(ctx context.Context, apiResp *admin.StreamsTenant) (*TFStreamInstanceModel, diag.Diagnostics) {
	hostnames, diags := types.ListValueFrom(ctx, types.StringType, apiResp.Hostnames)

	var dataProcessRegion = types.ObjectNull(ProcessRegionObjectType.AttrTypes)
	if apiResp.DataProcessRegion != nil {
		returnedProcessRegion, diagsProcessRegion := types.ObjectValueFrom(ctx, ProcessRegionObjectType.AttrTypes, TFInstanceProcessRegionSpecModel{
			CloudProvider: types.StringValue(apiResp.DataProcessRegion.CloudProvider),
			Region:        types.StringValue(apiResp.DataProcessRegion.Region),
		})
		dataProcessRegion = returnedProcessRegion
		diags.Append(diagsProcessRegion...)
	}
	if diags.HasError() {
		return nil, diags
	}

	return &TFStreamInstanceModel{
		ID:                types.StringPointerValue(apiResp.Id),
		InstanceName:      types.StringPointerValue(apiResp.Name),
		ProjectID:         types.StringPointerValue(apiResp.GroupId),
		DataProcessRegion: dataProcessRegion,
		Hostnames:         hostnames,
	}, nil
}

func NewTFStreamInstances(ctx context.Context, streamInstancesConfig *TFStreamInstancesModel, paginatedResult *admin.PaginatedApiStreamsTenant) (*TFStreamInstancesModel, diag.Diagnostics) {
	input := paginatedResult.GetResults()
	results := make([]TFStreamInstanceModel, len(input))
	for i := range input {
		instance, diags := NewTFStreamInstance(ctx, &input[i])
		if diags.HasError() {
			return nil, diags
		}
		results[i] = *instance
	}
	return &TFStreamInstancesModel{
		ID:           types.StringValue(id.UniqueId()),
		ProjectID:    streamInstancesConfig.ProjectID,
		PageNum:      streamInstancesConfig.PageNum,
		ItemsPerPage: streamInstancesConfig.ItemsPerPage,
		TotalCount:   types.Int64PointerValue(conversion.IntPtrToInt64Ptr(paginatedResult.TotalCount)),
		Results:      results,
	}, nil
}
