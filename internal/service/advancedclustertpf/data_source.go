package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
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
	useReplicationSpecPerShard := modelDS.UseReplicationSpecPerShard.ValueBool()
	api := d.Client.AtlasV2.ClustersApi
	clusterResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) {
			return nil
		}
		diags.AddError("errorRead", fmt.Sprintf(errorRead, clusterName, err.Error()))
		return nil
	}
	modelIn := &TFModel{
		ProjectID: modelDS.ProjectID,
		Name:      modelDS.Name,
	}
	modelOut, extraInfo := getBasicClusterModel(ctx, diags, d.Client, clusterResp, modelIn, !useReplicationSpecPerShard)
	if diags.HasError() {
		return nil
	}
	if extraInfo.AsymmetricShardUnsupported && !useReplicationSpecPerShard {
		diags.AddError("errorRead", "Please add `use_replication_spec_per_shard = true` to your data source configuration to enable asymmetric shard support. Refer to documentation for more details.")
		return nil
	}
	updateModelAdvancedConfig(ctx, diags, d.Client, modelOut, nil, nil)
	if diags.HasError() {
		return nil
	}
	modelOutDS, err := conversion.CopyModel[TFModelDS](modelOut)
	if err != nil {
		diags.AddError(errorRead, fmt.Sprintf("error setting model: %s", err.Error()))
		return nil
	}
	modelOutDS.UseReplicationSpecPerShard = modelDS.UseReplicationSpecPerShard // attrs not in resource model
	return modelOutDS
}
