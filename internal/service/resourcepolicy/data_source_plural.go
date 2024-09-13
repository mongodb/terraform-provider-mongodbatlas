package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &resourcePolicysDS{}
var _ datasource.DataSourceWithConfigure = &resourcePolicysDS{}

const (
	dataSourceNamePlural = "resource_policies"
	errorReadDSP         = "error reading plural data source " + dataSourceNamePlural
)

func PluralDataSource() datasource.DataSource {
	return &resourcePolicysDS{
		DSCommon: config.DSCommon{
			DataSourceName: dataSourceNamePlural,
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
	var cfg TFResourcePoliciesDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	orgID := cfg.OrgID.ValueString()
	apiResp, _, err := connV2.AtlasResourcePoliciesApi.GetAtlasResourcePolicies(ctx, orgID).Execute()

	if err != nil {
		resp.Diagnostics.AddError(errorReadDSP, err.Error())
		return
	}

	newResourcePolicysModel, diags := NewTFResourcePoliciesModel(ctx, orgID, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newResourcePolicysModel)...)
}
