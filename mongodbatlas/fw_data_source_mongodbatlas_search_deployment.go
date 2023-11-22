package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &SearchDeploymentDS{}
var _ datasource.DataSourceWithConfigure = &SearchDeploymentDS{}

func NewSearchDeploymentDS() datasource.DataSource {
	return &SearchDeploymentDS{
		DSCommon: config.DSCommon{
			DataSourceName: searchDeploymentName,
		},
	}
}

type tfSearchDeploymentDSModel struct {
	ID          types.String `tfsdk:"id"`
	ClusterName types.String `tfsdk:"cluster_name"`
	ProjectID   types.String `tfsdk:"project_id"`
	Specs       types.List   `tfsdk:"specs"`
	StateName   types.String `tfsdk:"state_name"`
}

type SearchDeploymentDS struct {
	config.DSCommon
}

func (d *SearchDeploymentDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"cluster_name": schema.StringAttribute{
				Required: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"specs": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_size": schema.StringAttribute{
							Computed: true,
						},
						"node_count": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"state_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *SearchDeploymentDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var searchDeploymentConfig tfSearchDeploymentDSModel
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

	newSearchDeploymentModel, diagnostics := newTFSearchDeployment(ctx, clusterName, deploymentResp, nil)
	resp.Diagnostics.Append(diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	dsModel := convertToDSModel(newSearchDeploymentModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, dsModel)...)
}

func convertToDSModel(inputModel *tfSearchDeploymentRSModel) tfSearchDeploymentDSModel {
	return tfSearchDeploymentDSModel{
		ID:          inputModel.ID,
		ClusterName: inputModel.ClusterName,
		ProjectID:   inputModel.ProjectID,
		Specs:       inputModel.Specs,
		StateName:   inputModel.StateName,
	}
}
