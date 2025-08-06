package encryptionatrestprivateendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Label that identifies the cloud provider for the Encryption At Rest private endpoint.",
			},
			"error_message": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Error message for failures associated with the Encryption At Rest private endpoint.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the Private Endpoint Service.",
			},
			"private_endpoint_connection_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Connection name of the Azure Private Endpoint.",
			},
			"region_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Cloud provider region in which the Encryption At Rest private endpoint is located.",
			},
			"status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "State of the Encryption At Rest private endpoint.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Delete: true,
			}),
		},
	}
}

type TFEarPrivateEndpointModel struct {
	CloudProvider                 types.String   `tfsdk:"cloud_provider"`
	ErrorMessage                  types.String   `tfsdk:"error_message"`
	ProjectID                     types.String   `tfsdk:"project_id"`
	ID                            types.String   `tfsdk:"id"`
	PrivateEndpointConnectionName types.String   `tfsdk:"private_endpoint_connection_name"`
	RegionName                    types.String   `tfsdk:"region_name"`
	Status                        types.String   `tfsdk:"status"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
}
