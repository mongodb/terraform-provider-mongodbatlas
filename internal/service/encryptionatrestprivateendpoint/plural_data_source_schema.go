package encryptionatrestprivateendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func PluralDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Label that identifies the cloud provider for the Encryption At Rest private endpoint.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
			"results": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: DSAttributes(false),
				},
				Computed:            true,
				MarkdownDescription: "List of returned documents that MongoDB Cloud provides when completing this request.",
			},
		},
	}
}

type TFEncryptionAtRestPrivateEndpointsDSModel struct {
	CloudProvider types.String                `tfsdk:"cloud_provider"`
	ProjectID     types.String                `tfsdk:"project_id"`
	Results       []TFEarPrivateEndpointModel `tfsdk:"results"`
}
