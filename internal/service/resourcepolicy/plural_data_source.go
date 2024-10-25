package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &resourcePolicysDS{}
var _ datasource.DataSourceWithConfigure = &resourcePolicysDS{}

const (
	dataSourcePluralName     = "resource_policies"
	fullDataSourcePluralName = "mongodbatlas_" + dataSourcePluralName
	errorReadDSP             = "error reading plural data source " + fullDataSourcePluralName
)

func PluralDataSource() datasource.DataSource {
	return &resourcePolicysDS{
		DSCommon: config.DSCommon{
			DataSourceName: dataSourcePluralName,
		},
	}
}

type resourcePolicysDS struct {
	config.DSCommon
}

func (d *resourcePolicysDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourcePluralSchema(ctx)
}

func (d *resourcePolicysDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg TFModelDSP
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	orgID := cfg.OrgID.ValueString()
	apiResp, _, err := connV2.ResourcePoliciesApi.GetAtlasResourcePolicies(ctx, orgID).Execute()

	if err != nil {
		resp.Diagnostics.AddError(errorReadDSP, err.Error())
		return
	}

	newResourcePolicysModel, diags := NewTFModelDSP(ctx, orgID, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicysModel)...)
}
