package encryptionatrestprivateendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				Required:            true,
				Description:         "Label that identifies the cloud provider for the Encryption At Rest private endpoint.",
				MarkdownDescription: "Label that identifies the cloud provider for the Encryption At Rest private endpoint.",
			},
			"error_message": schema.StringAttribute{
				Computed:            true,
				Description:         "Error message for failures associated with the Encryption At Rest private endpoint.",
				MarkdownDescription: "Error message for failures associated with the Encryption At Rest private endpoint.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies the Private Endpoint Service.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the Private Endpoint Service.",
			},
			"private_endpoint_connection_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Connection name of the Azure Private Endpoint.",
				MarkdownDescription: "Connection name of the Azure Private Endpoint.",
			},
			"region_name": schema.StringAttribute{
				Required:            true,
				Description:         "Cloud provider region in which the Encryption At Rest private endpoint is located.",
				MarkdownDescription: "Cloud provider region in which the Encryption At Rest private endpoint is located.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				Description:         "State of the Encryption At Rest private endpoint.",
				MarkdownDescription: "State of the Encryption At Rest private endpoint.",
			},
		},
	}
}

type TFEarPrivateEndpointModel struct {
	CloudProvider                 types.String `tfsdk:"cloud_provider"`
	ErrorMessage                  types.String `tfsdk:"error_message"`
	ProjectID                     types.String `tfsdk:"project_id"`
	ID                            types.String `tfsdk:"id"`
	PrivateEndpointConnectionName types.String `tfsdk:"private_endpoint_connection_name"`
	RegionName                    types.String `tfsdk:"region_name"`
	Status                        types.String `tfsdk:"status"`
}
