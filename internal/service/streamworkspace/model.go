package streamworkspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

// newStreamWorkspaceCreateReq creates an API request for creating a stream workspace.
func newStreamWorkspaceCreateReq(ctx context.Context, plan *TFModel) (*admin.StreamsTenant, diag.Diagnostics) {
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
		streamConfig := new(TFWorkspaceStreamConfigModel)
		if diags := plan.StreamConfig.As(ctx, streamConfig, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		var maxTierSize *string
		if !streamConfig.MaxTierSize.IsNull() && !streamConfig.MaxTierSize.IsUnknown() {
			maxTierSize = streamConfig.MaxTierSize.ValueStringPointer()
		}

		var tier *string
		if !streamConfig.Tier.IsNull() && !streamConfig.Tier.IsUnknown() {
			tier = streamConfig.Tier.ValueStringPointer()
		}
		streamTenant.StreamConfig = &admin.StreamConfig{
			MaxTierSize: maxTierSize,
			Tier:        tier,
		}
	}
	return streamTenant, nil
}

// newStreamWorkspaceUpdateReq creates an API request for updating a stream workspace.
func newStreamWorkspaceUpdateReq(ctx context.Context, plan *TFModel) (*admin.StreamsDataProcessRegion, diag.Diagnostics) {
	dataProcessRegion := &TFWorkspaceProcessRegionSpecModel{}
	if diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, diags
	}
	return &admin.StreamsDataProcessRegion{
		CloudProvider: dataProcessRegion.CloudProvider.ValueString(),
		Region:        dataProcessRegion.Region.ValueString(),
	}, nil
}

// FromInstanceModel populates this workspace model from a TFStreamInstanceModel and maps instance_name to workspace_name.
// This eliminates the need for conversion functions by directly updating fields.
func (m *TFModel) FromInstanceModel(instanceModel *streaminstance.TFStreamInstanceModel) {
	m.ID = instanceModel.ID
	m.WorkspaceName = instanceModel.InstanceName
	m.ProjectID = instanceModel.ProjectID
	m.DataProcessRegion = instanceModel.DataProcessRegion
	if instanceModel.StreamConfig.IsNull() {
		m.StreamConfig = types.ObjectNull(map[string]attr.Type{
			"max_tier_size": types.StringType,
			"tier":          types.StringType,
		})
	} else {
		instanceStreamConfigAttrs := instanceModel.StreamConfig.Attributes()
		tierValue := instanceStreamConfigAttrs["tier"]

		workspaceStreamConfig, _ := types.ObjectValue(
			map[string]attr.Type{
				"max_tier_size": types.StringType,
				"tier":          types.StringType,
			},
			map[string]attr.Value{
				"max_tier_size": types.StringNull(),
				"tier":          tierValue,
			},
		)
		m.StreamConfig = workspaceStreamConfig
	}
	m.Hostnames = instanceModel.Hostnames
}
