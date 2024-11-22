package mongodbemployeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	requiredFields := []string{"project_id", "cluster_name"}
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil)
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
	apiResp, _ := cluster.GetMongoDBEmployeeAccessGrantOk()
	if err != nil || apiResp == nil {
		msg := "employee access grant not defined for that cluster"
		if err != nil {
			msg = err.Error()
		}
		resp.Diagnostics.AddError(errorDataSource, msg)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModel(projectID, clusterName, apiResp))...)
}
