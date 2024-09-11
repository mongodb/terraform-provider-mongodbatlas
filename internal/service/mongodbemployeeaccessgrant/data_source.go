package mongodbemployeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
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
	resp.Schema = DataSourceSchema(ctx)
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	connV2 := d.Client.AtlasV2
	projectID := tfModel.ProjectID.ValueString()
	clusterName := tfModel.ClusterName.ValueString()
	cluster, _, err := connV2.ClustersApi.GetCluster(ctx, projectID, clusterName).Execute()
	atlasResp, _ := cluster.GetMongoDBEmployeeAccessGrantOk()
	if err != nil || atlasResp == nil {
		msg := "employee access grant not defined for that cluster"
		if err != nil {
			msg = err.Error()
		}
		resp.Diagnostics.AddError(errorDataSource, msg)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModel(projectID, clusterName, atlasResp))...)
}
