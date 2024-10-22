package testname

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"string_attr": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "string description",
			},
			"bool_attr": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "bool description",
			},
			"int_attr": schema.Int64Attribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "int description",
			},
			"float_attr": schema.Float64Attribute{
				Optional:            true,
				MarkdownDescription: "float description",
			},
			"number_attr": schema.NumberAttribute{
				Optional:            true,
				MarkdownDescription: "number description",
			},
			"simple_list_attr": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "simple arr description",
				ElementType:         types.StringType,
			},
			"simple_set_attr": schema.SetAttribute{
				Optional:            true,
				MarkdownDescription: "simple set description",
				ElementType:         types.Float64Type,
			},
			"simple_map_attr": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "simple map description",
				ElementType:         types.BoolType,
			},
		},
	}
}

type TFModel struct {
	StringAttr     types.String  `tfsdk:"string_attr"`
	BoolAttr       types.Bool    `tfsdk:"bool_attr"`
	IntAttr        types.Int64   `tfsdk:"int_attr"`
	FloatAttr      types.Float64 `tfsdk:"float_attr"`
	NumberAttr     types.Number  `tfsdk:"number_attr"`
	SimpleListAttr types.List    `tfsdk:"simple_list_attr"`
	SimpleSetAttr  types.Set     `tfsdk:"simple_set_attr"`
	SimpleMapAttr  types.Map     `tfsdk:"simple_map_attr"`
}
