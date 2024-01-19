package searchdeployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &searchDeploymentDS{}
var _ datasource.DataSourceWithConfigure = &searchDeploymentDS{}

func DataSource() datasource.DataSource {
	return &searchDeploymentDS{
		DSCommon: config.DSCommon{
			DataSourceName: searchDeploymentName,
		},
	}
}

type searchDeploymentDS struct {
	config.DSCommon
}

func (d *searchDeploymentDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema(ctx)
}

func (d *searchDeploymentDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var searchDeploymentConfig TFSearchDeploymentDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &searchDeploymentConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := searchDeploymentConfig.ProjectID.ValueString()
	clusterName := searchDeploymentConfig.ClusterName.ValueString()
	deploymentResp, _, err := connV2.AtlasSearchApi.GetAtlasSearchDeployment(ctx, projectID, clusterName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error getting search node information", err.Error())
		return
	}

	newSearchDeploymentModel, diagnostics := NewTFSearchDeployment(ctx, clusterName, deploymentResp, nil)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	dsModel := convertToDSModel(newSearchDeploymentModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, dsModel)...)
}

func convertToDSModel(inputModel *TFSearchDeploymentRSModel) TFSearchDeploymentDSModel {
	return TFSearchDeploymentDSModel{
		ID:          inputModel.ID,
		ClusterName: inputModel.ClusterName,
		ProjectID:   inputModel.ProjectID,
		Specs:       inputModel.Specs,
		StateName:   inputModel.StateName,
	}
}
