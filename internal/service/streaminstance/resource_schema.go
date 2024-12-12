package streaminstance

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"instance_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_process_region": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"cloud_provider": schema.StringAttribute{
						Required: true,
					},
					"region": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"hostnames": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"stream_config": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"tier": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
}

type TFStreamInstanceModel struct {
	ID                types.String `tfsdk:"id"`
	InstanceName      types.String `tfsdk:"instance_name"`
	ProjectID         types.String `tfsdk:"project_id"`
	DataProcessRegion types.Object `tfsdk:"data_process_region"`
	StreamConfig      types.Object `tfsdk:"stream_config"`
	Hostnames         types.List   `tfsdk:"hostnames"`
}

type TFInstanceProcessRegionSpecModel struct {
	CloudProvider types.String `tfsdk:"cloud_provider"`
	Region        types.String `tfsdk:"region"`
}

type TFInstanceStreamConfigSpecModel struct {
	Tier types.String `tfsdk:"tier"`
}

var ProcessRegionObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"cloud_provider": types.StringType,
	"region":         types.StringType,
}}

var StreamConfigObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"tier": types.StringType,
}}
