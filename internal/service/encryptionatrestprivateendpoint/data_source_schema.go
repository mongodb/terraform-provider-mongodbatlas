package encryptionatrestprivateendpoint

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DSAttributes(withArguments bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cloud_provider": schema.StringAttribute{
			Required:            withArguments,
			Computed:            !withArguments,
			MarkdownDescription: "Label that identifies the cloud provider for the Encryption At Rest private endpoint.",
		},
		"error_message": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Error message for failures associated with the Encryption At Rest private endpoint.",
		},
		"project_id": schema.StringAttribute{
			Required:            withArguments,
			Computed:            !withArguments,
			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
		},
		"id": schema.StringAttribute{
			Required:            withArguments,
			Computed:            !withArguments,
			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the Private Endpoint Service.",
		},
		"private_endpoint_connection_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Connection name of the Azure Private Endpoint.",
		},
		"region_name": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Cloud provider region in which the Encryption At Rest private endpoint is located.",
		},
		"status": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "State of the Encryption At Rest private endpoint.",
		},
	}
}

// TFEarPrivateEndpointModelDS represents the model for data sources (without timeout fields)
type TFEarPrivateEndpointModelDS struct {
	CloudProvider                 types.String `tfsdk:"cloud_provider"`
	ErrorMessage                  types.String `tfsdk:"error_message"`
	ProjectID                     types.String `tfsdk:"project_id"`
	ID                            types.String `tfsdk:"id"`
	PrivateEndpointConnectionName types.String `tfsdk:"private_endpoint_connection_name"`
	RegionName                    types.String `tfsdk:"region_name"`
	Status                        types.String `tfsdk:"status"`
}

type TFEncryptionAtRestPrivateEndpointsDSModel struct {
	CloudProvider types.String                  `tfsdk:"cloud_provider"`
	ProjectID     types.String                  `tfsdk:"project_id"`
	Results       []TFEarPrivateEndpointModelDS `tfsdk:"results"`
}
