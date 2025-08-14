package encryptionatrestprivateendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
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
			"delete_on_create_timeout": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					customplanmodifier.CreateOnlyBoolPlanModifier(),
				},
				MarkdownDescription: "Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.",
			},
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
	DeleteOnCreateTimeout         types.Bool     `tfsdk:"delete_on_create_timeout"`
}
