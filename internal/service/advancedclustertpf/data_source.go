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

const (
	errorReadDatasource                      = "Error reading  advanced cluster datasource"
	errorReadDatasourceForceAsymmetric       = "Error reading advanced cluster datasource, was expecting symmetric shards but found asymmetric shards"
	errorReadDatasourceForceAsymmetricDetail = "Cluster name %s. Please add `use_replication_spec_per_shard = true` to your data source configuration to enable asymmetric shard support. Refer to documentation for more details."
)

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
		diags.AddError(errorReadDatasource, defaultAPIErrorDetails(clusterName, err))
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
	if extraInfo.ForceLegacySchemaFailed {
		diags.AddError(errorReadDatasourceForceAsymmetric, fmt.Sprintf(errorReadDatasourceForceAsymmetricDetail, clusterName))
		return nil
	}
	updateModelAdvancedConfig(ctx, diags, d.Client, modelOut, nil, nil)
	if diags.HasError() {
		return nil
	}
	modelOutDS := conversion.CopyModel[TFModelDS](modelOut)
	modelOutDS.UseReplicationSpecPerShard = modelDS.UseReplicationSpecPerShard // attrs not in resource model
	return modelOutDS
}
