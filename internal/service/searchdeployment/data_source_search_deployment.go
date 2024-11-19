package searchdeployment

import (
	"context"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
	// TODO: THIS WILL BE REMOVED BEFORE MERGING, check old data source schema and new auto-generated schema are the same
	ds1 := DataSourceSchemaDelete(ctx)
	conversion.UpdateSchemaDescription(&ds1)
	ds2 := conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), "project_id", "cluster_name")
	if diff := cmp.Diff(ds1, ds2); diff != "" {
		log.Fatal(diff)
	}
	resp.Schema = ds2
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

	newSearchDeploymentModel, diagnostics := NewTFSearchDeployment(ctx, clusterName, deploymentResp, nil, true)
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
