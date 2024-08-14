package dsschema

import (
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func PaginatedDSSchema(arguments, resultAttributes map[string]schema.Attribute, skipID bool) schema.Schema {
	result := schema.Schema{
		Attributes: map[string]schema.Attribute{
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
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: resultAttributes,
				},
			},
		},
	}
	if !skipID {
		result.Attributes["id"] = schema.StringAttribute{
			Computed: true,
		}
	}
	maps.Copy(result.Attributes, arguments)
	return result
}
