package advancedclustertpf

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312001/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
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
	resp.Schema = pluralDataSourceSchema(ctx)
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModelPluralDS
	diags := &resp.Diagnostics
	diags.Append(req.Config.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	model, diags := d.readClusters(ctx, diags, &state)
	resp.Diagnostics = *diags
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *pluralDS) readClusters(ctx context.Context, diags *diag.Diagnostics, pluralModel *TFModelPluralDS) (*TFModelPluralDS, *diag.Diagnostics) {
	projectID := pluralModel.ProjectID.ValueString()
	useReplicationSpecPerShard := pluralModel.UseReplicationSpecPerShard.ValueBool()
	api := d.Client.AtlasV2.ClustersApi
	params := admin.ListClustersApiParams{
		GroupId: projectID,
	}
	list, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ClusterDescription20240805], *http.Response, error) {
		request := api.ListClustersWithParams(ctx, &params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		diags.AddError(errorList, fmt.Sprintf(errorListDetail, projectID, err.Error()))
		return nil, diags
	}
	outs := &TFModelPluralDS{
		ProjectID:                  pluralModel.ProjectID,
		UseReplicationSpecPerShard: pluralModel.UseReplicationSpecPerShard,
	}
	for i := range list {
		clusterResp := &list[i]
		modelOut, extraInfo := getBasicClusterModel(ctx, diags, d.Client, clusterResp, useReplicationSpecPerShard)
		if diags.HasError() {
			if DiagsHasOnlyClusterNotFoundErrors(diags) {
				diags = ResetClusterNotFoundErrors(diags)
				continue
			}
			return nil, diags
		}
		if extraInfo.UseOldShardingConfigFailed {
			continue
		}
		updateModelAdvancedConfig(ctx, diags, d.Client, modelOut, nil, nil)
		if diags.HasError() {
			if DiagsHasOnlyClusterNotFoundErrors(diags) {
				diags = ResetClusterNotFoundErrors(diags)
				continue
			}
			return nil, diags
		}
		modelOutDS := conversion.CopyModel[TFModelDS](modelOut)
		modelOutDS.UseReplicationSpecPerShard = pluralModel.UseReplicationSpecPerShard // attrs not in resource model
		outs.Results = append(outs.Results, modelOutDS)
	}
	flexModels := d.getFlexClustersModels(ctx, diags, projectID)
	if diags.HasError() {
		return nil, diags
	}
	outs.Results = append(outs.Results, flexModels...)
	return outs, diags
}
func DiagsHasOnlyClusterNotFoundErrors(diags *diag.Diagnostics) bool {
	for _, d := range *diags {
		if d.Severity() == diag.SeverityError && !strings.Contains(d.Detail(), "CLUSTER_NOT_FOUND") {
			return false
		}
	}
	return true
}

func ResetClusterNotFoundErrors(diags *diag.Diagnostics) *diag.Diagnostics {
	newDiags := &diag.Diagnostics{}
	for _, d := range *diags {
		if d.Severity() == diag.SeverityError && strings.Contains(d.Detail(), "CLUSTER_NOT_FOUND") {
			continue
		}
		newDiags.Append(d)
	}
	return newDiags
}

func (d *pluralDS) getFlexClustersModels(ctx context.Context, diags *diag.Diagnostics, projectID string) []*TFModelDS {
	var results []*TFModelDS

	listFlexClusters, err := flexcluster.ListFlexClusters(ctx, projectID, d.Client.AtlasV2.FlexClustersApi)
	if err != nil {
		diags.AddError(errorList, fmt.Sprintf(errorListDetail, projectID, err.Error()))
		return nil
	}

	for i := range *listFlexClusters {
		flexClusterResp := (*listFlexClusters)[i]
		modelOut := NewTFModelFlex(ctx, diags, &flexClusterResp, nil)
		if diags.HasError() {
			if DiagsHasOnlyClusterNotFoundErrors(diags) {
				diags = ResetClusterNotFoundErrors(diags)
				continue
			}
			return nil
		}
		modelOutDS := conversion.CopyModel[TFModelDS](modelOut)
		results = append(results, modelOutDS)
	}
	return results
}
