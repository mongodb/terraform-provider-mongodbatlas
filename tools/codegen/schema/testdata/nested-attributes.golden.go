package test_name

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
				},
			},
		},
	}
}
