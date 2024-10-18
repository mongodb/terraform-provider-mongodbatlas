package testname

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nested_single_attr": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "nested single attribute",
				Attributes: map[string]schema.Attribute{
					"string_attr": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "string attribute",
					},
					"int_attr": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "int attribute",
					},
				},
			},
			"nested_list_attr": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "nested list attribute",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attr": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "string attribute",
						},
						"int_attr": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "int attribute",
						},
					},
				},
			},
			"set_nested_attribute": schema.SetNestedAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "set nested attribute",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attr": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "string attribute",
						},
						"int_attr": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "int attribute",
						},
					},
				},
			},
			"map_nested_attribute": schema.MapNestedAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "map nested attribute",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"string_attr": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "string attribute",
						},
						"int_attr": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "int attribute",
						},
					},
				},
			},
		},
	}
}

type TFModel struct {
	NestedSingleAttr   types.Object `tfsdk:"nested_single_attr"`
	NestedListAttr     types.List   `tfsdk:"nested_list_attr"`
	SetNestedAttribute types.Set    `tfsdk:"set_nested_attribute"`
	MapNestedAttribute types.Map    `tfsdk:"map_nested_attribute"`
}
type TFNestedSingleAttrModel struct {
	StringAttr types.String `tfsdk:"string_attr"`
	IntAttr    types.Int64  `tfsdk:"int_attr"`
}

var NestedSingleAttrObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"string_attr": types.StringType,
	"int_attr":    types.Int64Type,
}}

type TFNestedListAttrModel struct {
	StringAttr types.String `tfsdk:"string_attr"`
	IntAttr    types.Int64  `tfsdk:"int_attr"`
}

var NestedListAttrObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"string_attr": types.StringType,
	"int_attr":    types.Int64Type,
}}

type TFSetNestedAttributeModel struct {
	StringAttr types.String `tfsdk:"string_attr"`
	IntAttr    types.Int64  `tfsdk:"int_attr"`
}

var SetNestedAttributeObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"string_attr": types.StringType,
	"int_attr":    types.Int64Type,
}}

type TFMapNestedAttributeModel struct {
	StringAttr types.String `tfsdk:"string_attr"`
	IntAttr    types.Int64  `tfsdk:"int_attr"`
}

var MapNestedAttributeObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"string_attr": types.StringType,
	"int_attr":    types.Int64Type,
}}
