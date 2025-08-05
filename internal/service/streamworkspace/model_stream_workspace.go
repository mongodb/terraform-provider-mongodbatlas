package streamworkspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func NewStreamWorkspaceCreateReq(ctx context.Context, plan *TFStreamWorkspaceModel) (*admin.StreamsTenant, diag.Diagnostics) {
	dataProcessRegion := &TFWorkspaceProcessRegionSpecModel{}
	if diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	streamTenant := &admin.StreamsTenant{
		GroupId: plan.ProjectID.ValueStringPointer(),
		Name:    plan.WorkspaceName.ValueStringPointer(),
		DataProcessRegion: &admin.StreamsDataProcessRegion{
			CloudProvider: dataProcessRegion.CloudProvider.ValueString(),
			Region:        dataProcessRegion.Region.ValueString(),
		},
	}
	if !plan.StreamConfig.IsNull() && !plan.StreamConfig.IsUnknown() {
		streamConfig := new(TFWorkspaceStreamConfigSpecModel)
		if diags := plan.StreamConfig.As(ctx, streamConfig, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		streamTenant.StreamConfig = &admin.StreamConfig{
			Tier: streamConfig.Tier.ValueStringPointer(),
		}
	}
	return streamTenant, nil
}

func NewStreamWorkspaceUpdateReq(ctx context.Context, plan *TFStreamWorkspaceModel) (*admin.StreamsDataProcessRegion, diag.Diagnostics) {
	dataProcessRegion := &TFWorkspaceProcessRegionSpecModel{}
	if diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsDataProcessRegion{
		CloudProvider: dataProcessRegion.CloudProvider.ValueString(),
		Region:        dataProcessRegion.Region.ValueString(),
	}, nil
}

func NewTFStreamWorkspace(ctx context.Context, apiResp *admin.StreamsTenant) (*TFStreamWorkspaceModel, diag.Diagnostics) {
	hostnames, diags := types.ListValueFrom(ctx, types.StringType, apiResp.Hostnames)

	var dataProcessRegion = types.ObjectNull(ProcessRegionObjectType.AttrTypes)
	if apiResp.DataProcessRegion != nil {
		returnedProcessRegion, diagsProcessRegion := types.ObjectValueFrom(ctx, ProcessRegionObjectType.AttrTypes, TFWorkspaceProcessRegionSpecModel{
			CloudProvider: types.StringValue(apiResp.DataProcessRegion.CloudProvider),
			Region:        types.StringValue(apiResp.DataProcessRegion.Region),
		})
		dataProcessRegion = returnedProcessRegion
		diags.Append(diagsProcessRegion...)
	}
	var streamConfig = types.ObjectNull(StreamConfigObjectType.AttrTypes)
	apiStreamConfig := apiResp.StreamConfig
	if apiStreamConfig != nil && apiStreamConfig.Tier != nil {
		returnedStreamConfig, diagsStreamConfig := types.ObjectValueFrom(ctx, StreamConfigObjectType.AttrTypes, TFWorkspaceStreamConfigSpecModel{
			Tier: types.StringPointerValue(apiStreamConfig.Tier),
		})
		streamConfig = returnedStreamConfig
		diags.Append(diagsStreamConfig...)
	}
	if diags.HasError() {
		return nil, diags
	}

	return &TFStreamWorkspaceModel{
		ID:                types.StringPointerValue(apiResp.Id),
		WorkspaceName:     types.StringPointerValue(apiResp.Name),
		ProjectID:         types.StringPointerValue(apiResp.GroupId),
		DataProcessRegion: dataProcessRegion,
		StreamConfig:      streamConfig,
		Hostnames:         hostnames,
	}, nil
}

func NewTFStreamWorkspaces(ctx context.Context, streamWorkspacesConfig *TFStreamWorkspacesModel, paginatedResult *admin.PaginatedApiStreamsTenant) (*TFStreamWorkspacesModel, diag.Diagnostics) {
	input := paginatedResult.GetResults()
	results := make([]TFStreamWorkspaceModel, len(input))
	for i := range input {
		workspace, diags := NewTFStreamWorkspace(ctx, &input[i])
		if diags.HasError() {
			return nil, diags
		}
		results[i] = *workspace
	}
	return &TFStreamWorkspacesModel{
		ID:           types.StringValue(id.UniqueId()),
		ProjectID:    streamWorkspacesConfig.ProjectID,
		PageNum:      streamWorkspacesConfig.PageNum,
		ItemsPerPage: streamWorkspacesConfig.ItemsPerPage,
		TotalCount:   types.Int64PointerValue(conversion.IntPtrToInt64Ptr(paginatedResult.TotalCount)),
		Results:      results,
	}, nil
}
