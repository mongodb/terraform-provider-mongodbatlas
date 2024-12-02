package dsschema

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// TODO: DECIDE IF DELETE THIS AND TEST FILE
func PaginatedDSSchema(arguments, resultAttributes map[string]schema.Attribute) schema.Schema {
	result := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"page_num": schema.Int64Attribute{
				Optional: true,
			},
			"items_per_page": schema.Int64Attribute{
				Optional: true,
			},
			"total_count": schema.Int64Attribute{
				Computed: true,
			},
			"results": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of returned documents that MongoDB Cloud provides when completing this request.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: resultAttributes,
				},
			},
		},
	}
	maps.Copy(result.Attributes, arguments)
	return result
}
