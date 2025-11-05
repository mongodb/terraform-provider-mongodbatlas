package streamaccountdetails

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

const resourceName = "stream_account_details"

var _ datasource.DataSource = &ds{}
var _ datasource.DataSourceWithConfigure = &ds{}

func DataSource() datasource.DataSource {
	return &ds{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type ds struct {
	config.DSCommon
}

func (d *ds) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	connV2 := d.Client.AtlasV2
	var streamAccountDetailsModel *TFStreamAccountDetailsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamAccountDetailsModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	accountDetailsRequestParams := connV2.StreamsApi.GetAccountDetailsWithParams(ctx, &admin.GetAccountDetailsApiParams{
		GroupId:       streamAccountDetailsModel.ProjectId.ValueString(),
		CloudProvider: streamAccountDetailsModel.CloudProvider.ValueStringPointer(),
		RegionName:    streamAccountDetailsModel.RegionName.ValueStringPointer(),
	})

	accountDetails, _, err := accountDetailsRequestParams.Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting account details for project", err.Error())
		return
	}

	newStreamAccountDetailsModel, diags := NewTFStreamAccountDetails(
		ctx,
		streamAccountDetailsModel.ProjectId.ValueString(),
		streamAccountDetailsModel.CloudProvider.ValueString(),
		streamAccountDetailsModel.RegionName.ValueString(),
		accountDetails,
	)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamAccountDetailsModel)...)
}
