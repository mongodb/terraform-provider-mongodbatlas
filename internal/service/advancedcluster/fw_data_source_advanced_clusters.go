package advancedcluster

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"
	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	advancedClustersDatasourceName = "advanced_clusters"
)

var _ datasource.DataSource = &advancedClustersDS{}
var _ datasource.DataSourceWithConfigure = &advancedClustersDS{}

type advancedClustersDS struct {
	config.DSCommon
}

func PluralDataSource() datasource.DataSource {
	return &advancedClustersDS{
		DSCommon: config.DSCommon{
			DataSourceName: advancedClustersDatasourceName,
		},
	}
}

func (d *advancedClustersDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework: https://github.com/hashicorp/terraform-plugin-testing/issues/84#issuecomment-1480006432
				DeprecationMessage: "Please use each cluster's id attribute instead",
				Computed:           true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: advancedClusterDSAttributes(),
				},
			},
		},
	}
}

func (d *advancedClustersDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Client.Atlas
	var clustersConfig tfAdvancedClustersDSModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &clustersConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := clustersConfig.ProjectID.ValueString()

	clusters, response, err := conn.AdvancedClusters.List(ctx, projectID, nil)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("advanced_clusters not found in Atlas", fmt.Sprintf("error reading advanced_clusters list for project(%s): %s", projectID, err))
			return
		}
		resp.Diagnostics.AddError("error in getting advanced_clusters from Atlas", fmt.Sprintf("error reading advanced_clusters list for project(%s): %s", projectID, err))
		return
	}

	newClustersState, diags := newTfAdvancedClustersDSModel(ctx, conn, clusters, projectID)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClustersState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	newClustersState.ID = types.StringValue(id.UniqueId())
}

func newTfAdvancedClustersDSModel(ctx context.Context, conn *mongodbatlas.Client, clusters *mongodbatlas.AdvancedClustersResponse, projectID string) (tfAdvancedClustersDSModel, diag.Diagnostics) {
	tfAdvClustersModel := tfAdvancedClustersDSModel{
		ID:        types.StringValue(id.UniqueId()),
		ProjectID: conversion.StringNullIfEmpty(projectID),
	}

	res, diags := newTfAdvancedClustersDSModelResults(ctx, conn, clusters, projectID)
	if diags.HasError() {
		return tfAdvClustersModel, diags
	}

	tfAdvClustersModel.Results = res
	return tfAdvClustersModel, nil
}

func newTfAdvancedClustersDSModelResults(ctx context.Context, conn *mongodbatlas.Client,
	apiResp *mongodbatlas.AdvancedClustersResponse,
	projectID string) ([]*tfAdvancedClusterDSModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	results := make([]*tfAdvancedClusterDSModel, 0)

	for i := range apiResp.Results {
		tfAdvCluster, diags := newTfAdvClustersDSModel(ctx, conn, apiResp.Results[i])
		if diags.HasError() {
			return nil, diags
		}

		results = append(results, tfAdvCluster)
	}
	return results, diags
}

func newTfAdvClustersDSModel(ctx context.Context, conn *matlas.Client, apiResp *matlas.AdvancedCluster) (*tfAdvancedClusterDSModel, diag.Diagnostics) {
	var err error
	projectID := apiResp.GroupID
	var diags diag.Diagnostics
	var d diag.Diagnostics

	clusterModel := tfAdvancedClusterDSModel{
		ID:                           types.StringValue(apiResp.ID),
		BackupEnabled:                types.BoolPointerValue(apiResp.BackupEnabled),
		ClusterType:                  types.StringValue(apiResp.ClusterType),
		CreateDate:                   types.StringValue(apiResp.CreateDate),
		DiskSizeGb:                   types.Float64PointerValue(apiResp.DiskSizeGB),
		EncryptionAtRestProvider:     types.StringValue(apiResp.EncryptionAtRestProvider),
		MongoDBMajorVersion:          types.StringValue(apiResp.MongoDBMajorVersion),
		MongoDBVersion:               types.StringValue(apiResp.MongoDBVersion),
		Name:                         types.StringValue(apiResp.Name),
		Paused:                       types.BoolPointerValue(apiResp.Paused),
		PitEnabled:                   types.BoolPointerValue(apiResp.PitEnabled),
		RootCertType:                 types.StringValue(apiResp.RootCertType),
		StateName:                    types.StringValue(apiResp.StateName),
		TerminationProtectionEnabled: types.BoolPointerValue(apiResp.TerminationProtectionEnabled),
		VersionReleaseSystem:         types.StringValue(apiResp.VersionReleaseSystem),
		ProjectID:                    types.StringValue(projectID),
	}
	clusterModel.BiConnectorConfig, d = types.ListValueFrom(ctx, TfBiConnectorConfigType, NewTfBiConnectorConfigModel(apiResp.BiConnector))
	diags.Append(d...)

	clusterModel.ConnectionStrings, d = types.ListValueFrom(ctx, tfConnectionStringType, newTfConnectionStringsModel(ctx, apiResp.ConnectionStrings))
	diags.Append(d...)

	clusterModel.Labels, d = types.SetValueFrom(ctx, TfLabelType, NewTfLabelsModel(apiResp.Labels))
	diags.Append(d...)

	clusterModel.Tags, d = types.SetValueFrom(ctx, TfTagType, NewTfTagsModel(&apiResp.Tags))
	diags.Append(d...)

	replicationSpecs, d := newTfReplicationSpecsDSModel(ctx, conn, apiResp.ReplicationSpecs, projectID)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}
	clusterModel.ReplicationSpecs, diags = types.SetValueFrom(ctx, tfReplicationSpecType, replicationSpecs)

	advancedConfiguration, err := NewTfAdvancedConfigurationModelDSFromAtlas(ctx, conn, projectID, apiResp.Name)
	if err != nil {
		diags.AddError("error when getting advanced_configuration from Atlas", err.Error())
		return nil, diags
	}
	clusterModel.AdvancedConfiguration, diags = types.ListValueFrom(ctx, tfAdvancedConfigurationType, advancedConfiguration)

	return &clusterModel, nil
}

type tfAdvancedClustersDSModel struct {
	ID        types.String                `tfsdk:"id"`
	ProjectID types.String                `tfsdk:"project_id"`
	Results   []*tfAdvancedClusterDSModel `tfsdk:"results"`
}
