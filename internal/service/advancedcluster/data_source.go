package advancedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

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
	resp.Schema = dataSourceSchema(ctx)
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModelDS
	diags := &resp.Diagnostics
	diags.Append(req.Config.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	model := d.readCluster(ctx, diags, &state)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *ds) readCluster(ctx context.Context, diags *diag.Diagnostics, modelDS *TFModelDS) *TFModelDS {
	clusterName := modelDS.Name.ValueString()
	projectID := modelDS.ProjectID.ValueString()
	clusterResp, flexClusterResp := GetClusterDetails(ctx, diags, projectID, clusterName, d.Client, false, modelDS.UseEffectiveFields.ValueBool())
	if diags.HasError() {
		return nil
	}
	if flexClusterResp == nil && clusterResp == nil {
		return nil
	}
	if flexClusterResp != nil {
		modelOutDS := newTFModelFlexDS(ctx, diags, flexClusterResp, nil)
		if diags.HasError() {
			return nil
		}
		modelOutDS.UseEffectiveFields = modelDS.UseEffectiveFields
		return modelOutDS
	}
	modelOutDS := getBasicClusterModelDS(ctx, diags, d.Client, clusterResp)
	if diags.HasError() {
		return nil
	}
	updateModelAdvancedConfigDS(ctx, diags, d.Client, modelOutDS, &ProcessArgs{
		ArgsDefault:           nil,
		ClusterAdvancedConfig: clusterResp.AdvancedConfiguration,
	})
	if diags.HasError() {
		return nil
	}
	modelOutDS.UseEffectiveFields = modelDS.UseEffectiveFields // Set Optional Terraform-only attribute.
	return modelOutDS
}
