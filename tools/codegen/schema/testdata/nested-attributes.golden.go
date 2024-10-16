package test_name

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
			"set_nested_attribute": schema.SetNestedAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "set nested attribute",
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
	}
}

type TFModel struct {
	NestedSingleAttr   types.Object `tfsdk:"nested_single_attr"`
	NestedListAttr     types.List   `tfsdk:"nested_list_attr"`
	SetNestedAttribute types.Set    `tfsdk:"set_nested_attribute"`
}
type TFNestedListAttrModel struct {
	StringAttr types.String `tfsdk:"string_attr"`
	IntAttr    types.Int64  `tfsdk:"int_attr"`
}

var TFNestedListAttrModelObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"string_attr": types.StringType,
	"int_attr":    types.StringType,
}}

type TFSetNestedAttributeModel struct {
	StringAttr types.String `tfsdk:"string_attr"`
	IntAttr    types.Int64  `tfsdk:"int_attr"`
}

var TFSetNestedAttributeModelObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"string_attr": types.StringType,
	"int_attr":    types.StringType,
}}
