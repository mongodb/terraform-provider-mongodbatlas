package flexcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func PluralDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the project.",
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: dataSourceSchema(true),
				},
				MarkdownDescription: "List of returned documents that MongoDB Cloud provides when completing this request.",
			},
		},
	}
}

type TFModelDSP struct {
	ProjectId types.String `tfsdk:"project_id"`
	Results   []TFModel    `tfsdk:"results"`
}
