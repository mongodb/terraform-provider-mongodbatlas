//nolint:gocritic
package encryptionatrestprivateendpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

var _ datasource.DataSource = &encryptionAtRestPrivateEndpointsDS{}
var _ datasource.DataSourceWithConfigure = &encryptionAtRestPrivateEndpointsDS{}

func PluralDataSource() datasource.DataSource {
	return &encryptionAtRestPrivateEndpointsDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", encryptionAtRestPrivateEndpointName),
		},
	}
}

type encryptionAtRestPrivateEndpointsDS struct {
	config.DSCommon
}

func (d *encryptionAtRestPrivateEndpointsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = PluralDataSourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *encryptionAtRestPrivateEndpointsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var earPrivateEndpointConfig TFEncryptionAtRestPrivateEndpointsDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &earPrivateEndpointConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := earPrivateEndpointConfig.ProjectID.ValueString()
	cloudProvider := earPrivateEndpointConfig.CloudProvider.ValueString()

	connV2 := d.Client.AtlasV2

	params := admin.GetEncryptionAtRestPrivateEndpointsForCloudProviderApiParams{
		GroupId:       projectID,
		CloudProvider: cloudProvider,
	}

	privateEndpoints, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.EARPrivateEndpoint], *http.Response, error) {
		request := connV2.EncryptionAtRestUsingCustomerKeyManagementApi.GetEncryptionAtRestPrivateEndpointsForCloudProviderWithParams(ctx, &params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		resp.Diagnostics.AddError("error fetching results", err.Error())
		return
	}
	newEarPrivateEndpointsModel := NewTFEarPrivateEndpoints(projectID, cloudProvider, privateEndpoints)
	resp.Diagnostics.Append(resp.State.Set(ctx, newEarPrivateEndpointsModel)...)
}
