package cluster

import (
	"context"
	"fmt"
	"log"
	"net/http"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

const (
	clustersDataSourceName = "clusters"
)

var _ datasource.DataSource = &clustersDS{}
var _ datasource.DataSourceWithConfigure = &clustersDS{}

func PluralDataSource() datasource.DataSource {
	return &clustersDS{
		DSCommon: config.DSCommon{
			DataSourceName: clustersDataSourceName,
		},
	}
}

type clustersDS struct {
	config.DSCommon
}

func (d *clustersDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
					Attributes: clusterDSAttributes(),
				},
			},
		},
	}
}

func (d *clustersDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Client.Atlas
	var clustersConfig tfClustersDSModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &clustersConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := clustersConfig.ProjectID.ValueString()

	clusters, response, err := conn.Clusters.List(ctx, projectID, nil)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("clusters not found in Atlas", fmt.Sprintf("error reading cluster list for project(%s): %s", projectID, err))
			return
		}
		resp.Diagnostics.AddError("error in getting clusters from Atlas", fmt.Sprintf("error reading cluster list for project(%s): %s", projectID, err))
		return
	}

	newClustersState, err := newTFClustersDSModel(ctx, conn, clusters, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error while getting clusters results from Atlas", fmt.Sprintf("error reading cluster list for project(%s): %s", projectID, err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClustersState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newClustersState.ID = conversion.StringNullIfEmpty(id.UniqueId())
}

func newTFClustersDSModel(ctx context.Context, conn *matlas.Client, clusters []matlas.Cluster, projectID string) (tfClustersDSModel, error) {
	tfClustersModel := tfClustersDSModel{
		ID:        conversion.StringNullIfEmpty(id.UniqueId()),
		ProjectID: conversion.StringNullIfEmpty(projectID),
	}

	res, err := newTFClustersDSModelResults(ctx, conn, clusters)
	if err != nil {
		return tfClustersModel, fmt.Errorf("error while getting clusters results from Atlas")
	}

	tfClustersModel.Results = res
	return tfClustersModel, nil
}

func newTFClustersDSModelResults(ctx context.Context, conn *matlas.Client, clusters []matlas.Cluster) ([]*tfClusterDSModel, error) {
	results := make([]*tfClusterDSModel, len(clusters))

	for i := range clusters {
		cluster := clusters[i]

		snapshotBackupPolicy, err := newTFSnapshotBackupPolicyModelFromAtlas(ctx, conn, cluster.GroupID, cluster.Name)
		if err != nil {
			return nil, err
		}

		advancedConfiguration, err := advancedcluster.NewTFAdvancedConfigurationModelDSFromAtlas(ctx, conn, cluster.GroupID, cluster.Name)
		if err != nil {
			return nil, err
		}

		var containerID string
		if cluster.ProviderSettings != nil && cluster.ProviderSettings.ProviderName != "TENANT" {
			containers, _, err := conn.Containers.List(ctx, cluster.GroupID,
				&matlas.ContainersListOptions{ProviderName: cluster.ProviderSettings.ProviderName})
			if err != nil {
				log.Printf(errorClusterRead, cluster.Name, err)
			}

			containerID = getContainerID(containers, &cluster)
		}
		result := &tfClusterDSModel{
			AdvancedConfiguration:              advancedConfiguration,
			AutoScalingComputeEnabled:          types.BoolPointerValue(cluster.AutoScaling.Compute.Enabled),
			AutoScalingComputeScaleDownEnabled: types.BoolPointerValue(cluster.AutoScaling.Compute.ScaleDownEnabled),
			AutoScalingDiskGbEnabled:           types.BoolPointerValue(cluster.AutoScaling.DiskGBEnabled),
			BackupEnabled:                      types.BoolPointerValue(cluster.BackupEnabled),
			ProviderBackupEnabled:              types.BoolPointerValue(cluster.ProviderBackupEnabled),
			ClusterType:                        conversion.StringNullIfEmpty(cluster.ClusterType),
			ConnectionStrings:                  newTFConnectionStringsModelDS(ctx, cluster.ConnectionStrings),
			DiskSizeGb:                         types.Float64PointerValue(cluster.DiskSizeGB),
			EncryptionAtRestProvider:           conversion.StringNullIfEmpty(cluster.EncryptionAtRestProvider),
			MongoDBMajorVersion:                conversion.StringNullIfEmpty(cluster.MongoDBMajorVersion),
			Name:                               conversion.StringNullIfEmpty(cluster.Name),
			NumShards:                          types.Int64PointerValue(cluster.NumShards),
			MongoDBVersion:                     conversion.StringNullIfEmpty(cluster.MongoDBVersion),
			MongoURI:                           conversion.StringNullIfEmpty(cluster.MongoURI),
			MongoURIUpdated:                    conversion.StringNullIfEmpty(cluster.MongoURIUpdated),
			MongoURIWithOptions:                conversion.StringNullIfEmpty(cluster.MongoURIWithOptions),
			PitEnabled:                         types.BoolPointerValue(cluster.PitEnabled),
			Paused:                             types.BoolPointerValue(cluster.Paused),
			SrvAddress:                         conversion.StringNullIfEmpty(cluster.SrvAddress),
			StateName:                          conversion.StringNullIfEmpty(cluster.StateName),
			ReplicationFactor:                  types.Int64PointerValue(cluster.ReplicationFactor),

			ProviderAutoScalingComputeMinInstanceSize: conversion.StringNullIfEmpty(cluster.ProviderSettings.AutoScaling.Compute.MinInstanceSize),
			ProviderAutoScalingComputeMaxInstanceSize: conversion.StringNullIfEmpty(cluster.ProviderSettings.AutoScaling.Compute.MaxInstanceSize),
			BackingProviderName:                       conversion.StringNullIfEmpty(cluster.ProviderSettings.BackingProviderName),
			ProviderDiskIops:                          types.Int64PointerValue(cluster.ProviderSettings.DiskIOPS),
			ProviderDiskTypeName:                      conversion.StringNullIfEmpty(cluster.ProviderSettings.DiskTypeName),
			ProviderEncryptEbsVolume:                  types.BoolPointerValue(cluster.ProviderSettings.EncryptEBSVolume),
			ProviderInstanceSizeName:                  conversion.StringNullIfEmpty(cluster.ProviderSettings.InstanceSizeName),
			ProviderName:                              conversion.StringNullIfEmpty(cluster.ProviderSettings.ProviderName),
			ProviderRegionName:                        conversion.StringNullIfEmpty(cluster.ProviderSettings.RegionName),

			BiConnectorConfig:            advancedcluster.NewTFBiConnectorConfigModel(cluster.BiConnector),
			ReplicationSpecs:             newTFReplicationSpecsModel(cluster.ReplicationSpecs),
			Labels:                       advancedcluster.NewTFLabelsModel(cluster.Labels),
			Tags:                         advancedcluster.NewTFTagsModel(cluster.Tags),
			SnapshotBackupPolicy:         snapshotBackupPolicy,
			TerminationProtectionEnabled: types.BoolPointerValue(cluster.TerminationProtectionEnabled),
			VersionReleaseSystem:         conversion.StringNullIfEmpty(cluster.VersionReleaseSystem),
			ContainerID:                  conversion.StringNullIfEmpty(containerID),
			ProjectID:                    conversion.StringNullIfEmpty(cluster.GroupID),
			ID:                           conversion.StringNullIfEmpty(cluster.ID),
		}
		results[i] = result
	}
	return results, nil
}

type tfClustersDSModel struct {
	ID        types.String        `tfsdk:"id"`
	ProjectID types.String        `tfsdk:"project_id"`
	Results   []*tfClusterDSModel `tfsdk:"results"`
}
