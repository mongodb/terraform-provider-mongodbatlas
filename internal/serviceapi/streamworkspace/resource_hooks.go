package streamworkspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

// Manual-parity custom code for mongodbatlas_stream_workspace (spike pilot, lift-and-shift).
//
// The manual resource (TPF) has no equivalent spec surface for: the `id` attribute (the API
// response carries `_id`, which cannot be a Terraform attribute name), requires-replace on
// `workspace_name`/`project_id`, and the write-once semantics of `failover_regions`. All are
// handled here (expanded_model: true provides the model slot for `id`).
//
// MoveState (from the 1.x stream_instance rename) is also hand-written below:
// autogen.HandleMove only copies id attributes by name, and this move renames
// instance_name -> workspace_name, so the generated mover cannot express it.

type TFExpandedModel struct {
	// id comes from the API `_id` (UUID) via apiname, like the manual resource's read.
	ID types.String `tfsdk:"id" apiname:"_id" autogen:"omitjson"`
}

var _ autogen.ResourceSchemaHook = (*rs)(nil)

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	// The generator marks path-param/omit-update attributes as CreateOnly (hard error on
	// change) because the PATCH request body lacks them. The manual resource replaces
	// instead, so rebuild both with RequiresReplace only.
	for _, name := range []string{"workspace_name", "project_id"} {
		if stringAttr, ok := baseSchema.Attributes[name].(schema.StringAttribute); ok {
			stringAttr.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
			baseSchema.Attributes[name] = stringAttr
		}
	}
	if listAttr, ok := baseSchema.Attributes["failover_regions"].(schema.ListNestedAttribute); ok {
		listAttr.PlanModifiers = append(listAttr.PlanModifiers, FailoverRegionsWriteOnce{})
		baseSchema.Attributes["failover_regions"] = listAttr
	}
	// id: the API _id (UUID), populated on read/create from the response body.
	baseSchema.Attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	return baseSchema
}

// FailoverRegionsWriteOnce requires replacement when failover_regions is already configured
// and being changed. The backend treats failover_regions as write-once: it can be set via PATCH
// exactly once, but cannot be modified or removed after that. Ported from the manual resource.
type FailoverRegionsWriteOnce struct{}

func (m FailoverRegionsWriteOnce) Description(_ context.Context) string {
	return "Requires replacement when failover_regions is already configured and being changed."
}

func (m FailoverRegionsWriteOnce) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m FailoverRegionsWriteOnce) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if len(req.StateValue.Elements()) == 0 {
		return
	}
	if req.PlanValue.IsUnknown() {
		return
	}
	if req.PlanValue.Equal(req.StateValue) {
		return
	}
	resp.RequiresReplace = true
}

// MoveState upgrades from the 1.x mongodbatlas_stream_instance resource (the manual
// implementation's `moved` support). Hand-written because autogen.HandleMove only copies
// id attributes by name, while this move renames instance_name -> workspace_name and
// reshapes the rest of the state. Ported from the manual resource's move_state.go.
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
			fmt.Sprintf("project_id: %s, instance_name: %s", conversion.SafeValue(projectID), conversion.SafeValue(instanceName)))
		return
	}

	// Create model with actual values from source state and map instance_name to workspace_name.
	model := &TFModel{
		Project_id:     types.StringPointerValue(projectID),
		Workspace_name: types.StringPointerValue(instanceName),
	}

	// Extract and preserve data_process_region if present.
	if dataProcessRegionVal, exists := stateObj["data_process_region"]; exists && !dataProcessRegionVal.IsNull() {
		var regionObj map[string]tftypes.Value
		if err := dataProcessRegionVal.As(&regionObj); err == nil {
			cloudProvider := schemafunc.GetAttrFromStateObj[string](regionObj, "cloud_provider")
			region := schemafunc.GetAttrFromStateObj[string](regionObj, "region")
			model.DataProcessRegion = customtypes.NewObjectValue[TFDataProcessRegionModel](ctx, TFDataProcessRegionModel{
				CloudProvider: types.StringPointerValue(cloudProvider),
				Region:        types.StringPointerValue(region),
			})
		}
	}
	if model.DataProcessRegion.IsNull() {
		model.DataProcessRegion = customtypes.NewObjectValueNull[TFDataProcessRegionModel](ctx)
	}

	// Extract and preserve stream_config if present.
	if streamConfigVal, exists := stateObj["stream_config"]; exists && !streamConfigVal.IsNull() {
		var configObj map[string]tftypes.Value
		if err := streamConfigVal.As(&configObj); err == nil {
			tier := schemafunc.GetAttrFromStateObj[string](configObj, "tier")
			model.StreamConfig = customtypes.NewObjectValue[TFStreamConfigModel](ctx, TFStreamConfigModel{
				MaxTierSize: types.StringNull(),
				Tier:        types.StringPointerValue(tier),
			})
		}
	}
	if model.StreamConfig.IsNull() {
		model.StreamConfig = customtypes.NewObjectValueNull[TFStreamConfigModel](ctx)
	}

	model.FailoverRegions = customtypes.NewNestedListValueNull[TFFailoverRegionsModel](ctx)

	// Extract and preserve hostnames if present.
	if hostnamesVal, exists := stateObj["hostnames"]; exists && !hostnamesVal.IsNull() {
		var hostnamesList []tftypes.Value
		if err := hostnamesVal.As(&hostnamesList); err == nil {
			var hostnames []attr.Value
			for _, hostnameVal := range hostnamesList {
				var hostname string
				if err := hostnameVal.As(&hostname); err == nil {
					hostnames = append(hostnames, types.StringValue(hostname))
				}
			}
			model.Hostnames = customtypes.NewListValue[types.String](ctx, hostnames)
		}
	}
	if model.Hostnames.IsNull() {
		model.Hostnames = customtypes.NewListValueNull[types.String](ctx)
	}

	// id is Computed and populated on the read that follows the move.
	model.ID = types.StringNull()

	resp.Diagnostics.Append(resp.TargetState.Set(ctx, model)...)
}
