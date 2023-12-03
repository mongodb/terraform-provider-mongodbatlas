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

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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

	newClustersState.ID = types.StringValue(id.UniqueId())
}

func newTFClustersDSModel(ctx context.Context, conn *matlas.Client, clusters []matlas.Cluster, projectID string) (tfClustersDSModel, error) {
	tfClustersModel := tfClustersDSModel{
		ID:        types.StringValue(id.UniqueId()),
		ProjectID: types.StringValue(projectID),
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

		advancedConfiguration, err := newTFAdvancedConfigurationModelDSFromAtlas(ctx, conn, cluster.GroupID, cluster.Name)
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
			ClusterType:                        types.StringValue(cluster.ClusterType),
			ConnectionStrings:                  newTFConnectionStringsModelDS(ctx, cluster.ConnectionStrings),
			DiskSizeGb:                         types.Float64PointerValue(cluster.DiskSizeGB),
			EncryptionAtRestProvider:           types.StringValue(cluster.EncryptionAtRestProvider),
			MongoDBMajorVersion:                types.StringValue(cluster.MongoDBMajorVersion),
			Name:                               types.StringValue(cluster.Name),
			NumShards:                          types.Int64PointerValue(cluster.NumShards),
			MongoDBVersion:                     types.StringValue(cluster.MongoDBVersion),
			MongoURI:                           types.StringValue(cluster.MongoURI),
			MongoURIUpdated:                    types.StringValue(cluster.MongoURIUpdated),
			MongoURIWithOptions:                types.StringValue(cluster.MongoURIWithOptions),
			PitEnabled:                         types.BoolPointerValue(cluster.PitEnabled),
			Paused:                             types.BoolPointerValue(cluster.Paused),
			SrvAddress:                         types.StringValue(cluster.SrvAddress),
			StateName:                          types.StringValue(cluster.StateName),
			ReplicationFactor:                  types.Int64PointerValue(cluster.ReplicationFactor),

			ProviderAutoScalingComputeMinInstanceSize: types.StringValue(cluster.ProviderSettings.AutoScaling.Compute.MinInstanceSize),
			ProviderAutoScalingComputeMaxInstanceSize: types.StringValue(cluster.ProviderSettings.AutoScaling.Compute.MaxInstanceSize),
			BackingProviderName:                       types.StringValue(cluster.ProviderSettings.BackingProviderName),
			ProviderDiskIops:                          types.Int64PointerValue(cluster.ProviderSettings.DiskIOPS),
			ProviderDiskTypeName:                      types.StringValue(cluster.ProviderSettings.DiskTypeName),
			ProviderEncryptEbsVolume:                  types.BoolPointerValue(cluster.ProviderSettings.EncryptEBSVolume),
			ProviderInstanceSizeName:                  types.StringValue(cluster.ProviderSettings.InstanceSizeName),
			ProviderName:                              types.StringValue(cluster.ProviderSettings.ProviderName),
			ProviderRegionName:                        types.StringValue(cluster.ProviderSettings.RegionName),

			BiConnectorConfig:            newTFBiConnectorConfigModel(cluster.BiConnector),
			ReplicationSpecs:             newTFReplicationSpecsModel(cluster.ReplicationSpecs),
			Labels:                       newTFLabelsModel(cluster.Labels),
			Tags:                         newTFTagsModel(cluster.Tags),
			SnapshotBackupPolicy:         snapshotBackupPolicy,
			TerminationProtectionEnabled: types.BoolPointerValue(cluster.TerminationProtectionEnabled),
			VersionReleaseSystem:         types.StringValue(cluster.VersionReleaseSystem),
			ContainerID:                  types.StringValue(containerID),
			ProjectID:                    types.StringValue(cluster.GroupID),
			ID:                           types.StringValue(cluster.ID),
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
