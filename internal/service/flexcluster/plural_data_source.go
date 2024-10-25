package flexcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"
)

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

func PluralDataSource() datasource.DataSource {
	return &pluralDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", resourceName),
		},
	}
}

type pluralDS struct {
	config.DSCommon
}

func (d *pluralDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = PluralDataSourceSchema(ctx)
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfModel TFModelDSP
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	apiResp, err := getFlexClusterList(ctx, connV2, tfModel.ProjectId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("error reading plural data source", err.Error())
		return
	}

	newFlexClustersModel, diags := NewTFModelDSP(ctx, tfModel.ProjectId.ValueString(), apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newFlexClustersModel)...)
}

func getFlexClusterList(ctx context.Context, connV2 *admin.APIClient, projectId string) ([]admin.FlexClusterDescription20250101, error) {
	var list []admin.FlexClusterDescription20250101
	apiResp, _, err := connV2.FlexClustersApi.ListFlexClusters(ctx, projectId).Execute()
	if err != nil {
		return nil, fmt.Errorf("error reading plural data source: %s", err)
	}

	for _, result := range apiResp.GetResults() {
		if cluster, ok := result.(admin.FlexClusterDescription20250101); ok {
			list = append(list, cluster)
		} else {
			return nil, fmt.Errorf("error reading plural data source: %s", err)
		}
	}

	return list, nil
}
