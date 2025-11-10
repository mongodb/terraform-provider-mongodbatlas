package streamworkspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

// MoveState is used with moved block to upgrade from stream_instance to stream_workspace.
func (r *rs) MoveState(context.Context) []resource.StateMover {
	return []resource.StateMover{{StateMover: stateMover}}
}

func stateMover(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	if req.SourceTypeName != "mongodbatlas_stream_instance" || !strings.HasSuffix(req.SourceProviderAddress, "/mongodbatlas") {
		return
	}

	// Extract all fields from source state to preserve values during move.
	stateAttrs := map[string]tftypes.Type{
		"project_id":    tftypes.String,
		"instance_name": tftypes.String,
		"data_process_region": tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"cloud_provider": tftypes.String,
				"region":         tftypes.String,
			},
		},
		"stream_config": tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"tier": tftypes.String,
			},
		},
		"hostnames": tftypes.List{
			ElementType: tftypes.String,
		},
	}

	rawStateValue, err := req.SourceRawState.UnmarshalWithOpts(tftypes.Object{
		AttributeTypes: stateAttrs,
	}, tfprotov6.UnmarshalOpts{ValueFromJSONOpts: tftypes.ValueFromJSONOpts{IgnoreUndefinedAttributes: true}})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Unmarshal state", err.Error())
		return
	}

	var stateObj map[string]tftypes.Value
	if err := rawStateValue.As(&stateObj); err != nil {
		resp.Diagnostics.AddError("Unable to Parse state", err.Error())
		return
	}

	projectID := schemafunc.GetAttrFromStateObj[string](stateObj, "project_id")
	instanceName := schemafunc.GetAttrFromStateObj[string](stateObj, "instance_name")

	if !conversion.IsStringPresent(projectID) || !conversion.IsStringPresent(instanceName) {
		resp.Diagnostics.AddError("Unable to read project_id or instance_name from state",
			fmt.Sprintf("project_id: %s, instance_name: %s", conversion.SafeString(projectID), conversion.SafeString(instanceName)))
		return
	}

	// Create model with actual values from source state and map instance_name to workspace_name.
	model := &TFModel{
		ID:            types.StringNull(), // Will be computed during read
		ProjectID:     types.StringPointerValue(projectID),
		WorkspaceName: types.StringPointerValue(instanceName),
	}

	// Extract and preserve data_process_region if present.
	if dataProcessRegionVal, exists := stateObj["data_process_region"]; exists && !dataProcessRegionVal.IsNull() {
		var regionObj map[string]tftypes.Value
		if err := dataProcessRegionVal.As(&regionObj); err == nil {
			cloudProvider := schemafunc.GetAttrFromStateObj[string](regionObj, "cloud_provider")
			region := schemafunc.GetAttrFromStateObj[string](regionObj, "region")

			objValue, diags := types.ObjectValue(map[string]attr.Type{
				"cloud_provider": types.StringType,
				"region":         types.StringType,
			}, map[string]attr.Value{
				"cloud_provider": types.StringPointerValue(cloudProvider),
				"region":         types.StringPointerValue(region),
			})
			if !diags.HasError() {
				model.DataProcessRegion = objValue
			}
		}
	}
	if model.DataProcessRegion.IsNull() {
		model.DataProcessRegion = types.ObjectNull(map[string]attr.Type{
			"cloud_provider": types.StringType,
			"region":         types.StringType,
		})
	}

	// Extract and preserve stream_config if present.
	if streamConfigVal, exists := stateObj["stream_config"]; exists && !streamConfigVal.IsNull() {
		var configObj map[string]tftypes.Value
		if err := streamConfigVal.As(&configObj); err == nil {
			tier := schemafunc.GetAttrFromStateObj[string](configObj, "tier")

			objValue, diags := types.ObjectValue(map[string]attr.Type{
				"tier": types.StringType,
			}, map[string]attr.Value{
				"tier": types.StringPointerValue(tier),
			})
			if !diags.HasError() {
				model.StreamConfig = objValue
			}
		}
	}
	if model.StreamConfig.IsNull() {
		model.StreamConfig = types.ObjectNull(map[string]attr.Type{
			"tier": types.StringType,
		})
	}

	// Extract and preserve hostnames if present.
	if hostnamesVal, exists := stateObj["hostnames"]; exists && !hostnamesVal.IsNull() {
		var hostnamesList []tftypes.Value
		if err := hostnamesVal.As(&hostnamesList); err == nil {
			var hostnames []string
			for _, hostnameVal := range hostnamesList {
				var hostname string
				if err := hostnameVal.As(&hostname); err == nil {
					hostnames = append(hostnames, hostname)
				}
			}
			listValue, diags := types.ListValueFrom(ctx, types.StringType, hostnames)
			if !diags.HasError() {
				model.Hostnames = listValue
			}
		}
	}
	if model.Hostnames.IsNull() {
		model.Hostnames = types.ListNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.TargetState.Set(ctx, model)...)
}
