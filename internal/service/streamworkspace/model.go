package streamworkspace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
	"go.mongodb.org/atlas-sdk/v20250312021/admin"
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
	if !plan.FailoverRegions.IsNull() && !plan.FailoverRegions.IsUnknown() {
		failoverDataRegions, diags := failoverRegionsToSDK(ctx, plan.FailoverRegions)
		if diags.HasError() {
			return nil, diags
		}
		streamTenant.FailoverRegions = &failoverDataRegions
	}
	if !plan.StreamConfig.IsNull() && !plan.StreamConfig.IsUnknown() {
		streamConfig := new(TFWorkspaceStreamConfigModel)
		if diags := plan.StreamConfig.As(ctx, streamConfig, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		var maxTierSize *string
		if streamConfig.MaxTierSize.ValueString() != "" {
			maxTierSize = streamConfig.MaxTierSize.ValueStringPointer()
		}
		var tier *string
		if streamConfig.Tier.ValueString() != "" {
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
// dataProcessRegion and failoverRegions are mutually exclusive in the PATCH body.
// failoverRegions is sent only when it is being set for the first time (state was null/empty).
func newStreamWorkspaceUpdateReq(ctx context.Context, plan, state *TFModel) (*admin.StreamsTenantUpdateRequest, diag.Diagnostics) {
	updateReq := &admin.StreamsTenantUpdateRequest{}
	if failoverRegionsChanging(plan, state) {
		failoverDataRegions, diags := failoverRegionsToSDK(ctx, plan.FailoverRegions)
		if diags.HasError() {
			return nil, diags
		}
		updateReq.FailoverRegions = &failoverDataRegions
	} else {
		dataProcessRegion := &TFWorkspaceProcessRegionSpecModel{}
		if diags := plan.DataProcessRegion.As(ctx, dataProcessRegion, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, diags
		}
		updateReq.CloudProvider = dataProcessRegion.CloudProvider.ValueStringPointer()
		updateReq.Region = dataProcessRegion.Region.ValueStringPointer()
		if !plan.StreamConfig.IsNull() && !plan.StreamConfig.IsUnknown() {
			streamConfig := new(TFWorkspaceStreamConfigModel)
			if diags := plan.StreamConfig.As(ctx, streamConfig, basetypes.ObjectAsOptions{}); diags.HasError() {
				return nil, diags
			}
			var maxTierSize *string
			if streamConfig.MaxTierSize.ValueString() != "" {
				maxTierSize = streamConfig.MaxTierSize.ValueStringPointer()
			}
			var tier *string
			if streamConfig.Tier.ValueString() != "" {
				tier = streamConfig.Tier.ValueStringPointer()
			}
			updateReq.StreamConfig = &admin.StreamConfig{
				MaxTierSize: maxTierSize,
				Tier:        tier,
			}
		}
	}
	return updateReq, nil
}

// failoverRegionsChanging reports whether failover_regions is being newly configured in this update
// (state was null/empty and plan now has values). This is distinct from unchanged inherited state.
func failoverRegionsChanging(plan, state *TFModel) bool {
	stateHasNoRegions := state.FailoverRegions.IsNull() || len(state.FailoverRegions.Elements()) == 0
	planHasRegions := !plan.FailoverRegions.IsNull() && len(plan.FailoverRegions.Elements()) > 0
	return stateHasNoRegions && planHasRegions
}

func failoverRegionsToSDK(ctx context.Context, list types.List) ([]admin.StreamsDataProcessRegion, diag.Diagnostics) {
	var regions []TFWorkspaceProcessRegionSpecModel
	if diags := list.ElementsAs(ctx, &regions, false); diags.HasError() {
		return nil, diags
	}
	result := make([]admin.StreamsDataProcessRegion, 0, len(regions))
	for _, r := range regions {
		result = append(result, admin.StreamsDataProcessRegion{
			CloudProvider: r.CloudProvider.ValueString(),
			Region:        r.Region.ValueString(),
		})
	}
	return result, nil
}

// newTFWorkspaceModel builds a workspace TFModel from an API response, including failover_regions.
func newTFWorkspaceModel(ctx context.Context, apiResp *admin.StreamsTenant) (*TFModel, diag.Diagnostics) {
	instanceModel, diags := streaminstance.NewTFStreamInstance(ctx, apiResp)
	if diags.HasError() {
		return nil, diags
	}
	var model TFModel
	model.FromInstanceModel(instanceModel)
	if apiResp.FailoverRegions != nil && len(*apiResp.FailoverRegions) > 0 {
		failoverRegions, diags := types.ListValueFrom(ctx, failoverRegionObjectType, toTFFailoverRegions(*apiResp.FailoverRegions))
		if diags.HasError() {
			return nil, diags
		}
		model.FailoverRegions = failoverRegions
	}
	return &model, nil
}

func toTFFailoverRegions(regions []admin.StreamsDataProcessRegion) []TFWorkspaceProcessRegionSpecModel {
	result := make([]TFWorkspaceProcessRegionSpecModel, len(regions))
	for i, r := range regions {
		result[i] = TFWorkspaceProcessRegionSpecModel{
			CloudProvider: types.StringValue(r.CloudProvider),
			Region:        types.StringValue(r.Region),
		}
	}
	return result
}

// FromInstanceModel populates this workspace model from a TFStreamInstanceModel and maps instance_name to workspace_name.
// This eliminates the need for conversion functions by directly updating fields.
func (m *TFModel) FromInstanceModel(instanceModel *streaminstance.TFStreamInstanceModel) {
	m.ID = instanceModel.ID
	m.WorkspaceName = instanceModel.InstanceName
	m.ProjectID = instanceModel.ProjectID
	m.DataProcessRegion = instanceModel.DataProcessRegion
	m.FailoverRegions = types.ListNull(failoverRegionObjectType)
	if instanceModel.StreamConfig.IsNull() {
		m.StreamConfig = types.ObjectNull(map[string]attr.Type{
			"max_tier_size": types.StringType,
			"tier":          types.StringType,
		})
	} else {
		instanceStreamConfigAttrs := instanceModel.StreamConfig.Attributes()
		tierValue := instanceStreamConfigAttrs["tier"]
		maxTierSizeValue := instanceStreamConfigAttrs["max_tier_size"]

		workspaceStreamConfig, _ := types.ObjectValue(
			map[string]attr.Type{
				"max_tier_size": types.StringType,
				"tier":          types.StringType,
			},
			map[string]attr.Value{
				"max_tier_size": maxTierSizeValue,
				"tier":          tierValue,
			},
		)
		m.StreamConfig = workspaceStreamConfig
	}
	m.Hostnames = instanceModel.Hostnames
}
