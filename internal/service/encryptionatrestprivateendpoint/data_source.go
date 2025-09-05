package encryptionatrestprivateendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &encryptionAtRestPrivateEndpointDS{}
var _ datasource.DataSourceWithConfigure = &encryptionAtRestPrivateEndpointDS{}

func DataSource() datasource.DataSource {
	return &encryptionAtRestPrivateEndpointDS{
		DSCommon: config.DSCommon{
			DataSourceName: encryptionAtRestPrivateEndpointName,
		},
	}
}

type encryptionAtRestPrivateEndpointDS struct {
	config.DSCommon
}

func (d *encryptionAtRestPrivateEndpointDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"id", "project_id", "cloud_provider"},
	})
}

func (d *encryptionAtRestPrivateEndpointDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var earPrivateEndpointConfig TFEarPrivateEndpointModelDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &earPrivateEndpointConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := earPrivateEndpointConfig.ProjectID.ValueString()
	cloudProvider := earPrivateEndpointConfig.CloudProvider.ValueString()
	endpointID := earPrivateEndpointConfig.ID.ValueString()

	endpointModel, _, err := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetRestPrivateEndpoint(ctx, projectID, cloudProvider, endpointID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFEarPrivateEndpointDS(*endpointModel, projectID))...)
}
