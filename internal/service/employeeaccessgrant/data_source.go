package employeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &employeeAccessGrantDS{}
var _ datasource.DataSourceWithConfigure = &employeeAccessGrantDS{}

func DataSource() datasource.DataSource {
	return &employeeAccessGrantDS{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type employeeAccessGrantDS struct {
	config.DSCommon
}

func (d *employeeAccessGrantDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *employeeAccessGrantDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
		msg := "info not found"
		if err != nil {
			msg = err.Error()
		}
		resp.Diagnostics.AddError(errorDataSource, msg)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModel(projectID, clusterName, atlasResp))...)
}
