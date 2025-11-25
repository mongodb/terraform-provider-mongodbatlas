package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &resourcePolicyDS{}
var _ datasource.DataSourceWithConfigure = &resourcePolicyDS{}

const (
	errorReadDS = "error reading data source " + fullResourceName
)

func DataSource() datasource.DataSource {
	return &resourcePolicyDS{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type resourcePolicyDS struct {
	config.DSCommon
}

func (d *resourcePolicyDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"org_id", "id"},
	})
}

func (d *resourcePolicyDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var cfg TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	apiResp, _, err := connV2.ResourcePoliciesApi.GetOrgResourcePolicy(ctx, cfg.OrgID.ValueString(), cfg.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorReadDS, err.Error())
		return
	}

	out, diags := NewTFModel(ctx, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, out)...)
}
