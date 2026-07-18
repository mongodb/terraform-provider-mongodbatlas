package encryptionatrestprivateendpoint

import "github.com/hashicorp/terraform-plugin-framework/types"

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
