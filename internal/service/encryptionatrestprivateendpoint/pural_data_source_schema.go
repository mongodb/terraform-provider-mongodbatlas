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
				Description:         "Human-readable label that identifies the cloud provider for the private endpoints to return.",
				MarkdownDescription: "Human-readable label that identifies the cloud provider for the private endpoints to return.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"results": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: DSAttributes(true),
				},
				Computed:            true,
				Description:         "List of returned documents that MongoDB Cloud providers when completing this request.",
				MarkdownDescription: "List of returned documents that MongoDB Cloud providers when completing this request.",
			},
		},
	}
}

type TFEncryptionAtRestPrivateEndpointsDSModel struct {
	CloudProvider types.String                `tfsdk:"cloud_provider"`
	ProjectID     types.String                `tfsdk:"project_id"`
	Results       []TFEarPrivateEndpointModel `tfsdk:"results"`
}
