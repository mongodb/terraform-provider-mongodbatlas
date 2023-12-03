package advancedcluster

import (
	"context"
	"fmt"
	"net/http"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	AdvancedClusterResourceName = "advanced_cluster"
)

var _ datasource.DataSource = &advancedClusterDS{}
var _ datasource.DataSourceWithConfigure = &advancedClusterDS{}

type advancedClusterDS struct {
	config.DSCommon
}

func DataSource() datasource.DataSource {
	return &advancedClusterDS{
		DSCommon: config.DSCommon{
			DataSourceName: AdvancedClusterResourceName,
		},
	}
}

func (d *advancedClusterDS) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"advanced_configuration": ClusterDSAdvancedConfigurationSchemaAttr(),
			"backup_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"bi_connector_config": ClusterDSBiConnectorConfigSchemaAttr(),
			"cluster_type": schema.StringAttribute{
				Computed: true,
			},
			"connection_strings": advClusterDSConnectionStringSchemaAttr(),
			"create_date": schema.StringAttribute{
				Computed: true,
			},
			"disk_size_gb": schema.Float64Attribute{
				Computed: true,
			},
			"encryption_at_rest_provider": schema.StringAttribute{
				Computed: true,
			},
			"labels": ClusterDSLabelsSchemaAttr(),
			"tags":   ClusterDSTagsSchemaAttr(),
			"mongo_db_major_version": schema.StringAttribute{
				Computed: true,
			},
			"mongo_db_version": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"paused": schema.BoolAttribute{
				Computed: true,
			},
			"pit_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"replication_specs": advClusterDSReplicationSpecsSchemaAttr(),
			"root_cert_type": schema.StringAttribute{
				Computed: true,
			},
			"state_name": schema.StringAttribute{
				Computed: true,
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"version_release_system": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *advancedClusterDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Client.Atlas
	var clusterConfig tfAdvancedClusterDSModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &clusterConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := clusterConfig.ProjectID.ValueString()
	clusterName := clusterConfig.Name.ValueString()

	cluster, response, err := conn.AdvancedClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("cluster not found in Atlas", fmt.Sprintf(errorClusterAdvancedRead, clusterName, err))
			return
		}
		resp.Diagnostics.AddError("error in getting cluster details from Atlas", fmt.Sprintf(errorClusterAdvancedRead, clusterName, err))
		return
	}

	newClusterState, err := newTfAdvClusterDSModel(ctx, conn, cluster)
	if err != nil {
		resp.Diagnostics.AddError("error in getting cluster details from Atlas", fmt.Sprintf(errorClusterAdvancedRead, clusterName, err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTfAdvClusterDSModel(ctx context.Context, conn *matlas.Client, apiResp *matlas.AdvancedCluster) (*tfAdvancedClusterDSModel, error) {
	var err error
	projectID := apiResp.GroupID

	clusterModel := tfAdvancedClusterDSModel{
		ID:                           types.StringValue(apiResp.ID),
		BackupEnabled:                types.BoolPointerValue(apiResp.BackupEnabled),
		BiConnectorConfig:            NewTFBiConnectorConfigModel(apiResp.BiConnector),
		ClusterType:                  types.StringValue(apiResp.ClusterType),
		ConnectionStrings:            newTFConnectionStringsModel(ctx, apiResp.ConnectionStrings),
		CreateDate:                   types.StringValue(apiResp.CreateDate),
		DiskSizeGb:                   types.Float64PointerValue(apiResp.DiskSizeGB),
		EncryptionAtRestProvider:     types.StringValue(apiResp.EncryptionAtRestProvider),
		Labels:                       RemoveDefaultLabel(NewTFLabelsModel(apiResp.Labels)),
		Tags:                         NewTFTagsModel(&apiResp.Tags),
		MongoDBMajorVersion:          types.StringValue(apiResp.MongoDBMajorVersion),
		MongoDBVersion:               types.StringValue(apiResp.MongoDBVersion),
		Name:                         types.StringValue(apiResp.Name),
		Paused:                       types.BoolPointerValue(apiResp.Paused),
		PitEnabled:                   types.BoolPointerValue(apiResp.PitEnabled),
		ReplicationSpecs:             newTFReplicationSpecsDSModel(apiResp.ReplicationSpecs, projectID),
		RootCertType:                 types.StringValue(apiResp.RootCertType),
		StateName:                    types.StringValue(apiResp.StateName),
		TerminationProtectionEnabled: types.BoolPointerValue(apiResp.TerminationProtectionEnabled),
		VersionReleaseSystem:         types.StringValue(apiResp.VersionReleaseSystem),
	}

	clusterModel.AdvancedConfiguration, err = NewTFAdvancedConfigurationModelDSFromAtlas(ctx, conn, projectID, apiResp.Name)
	if err != nil {
		return nil, err
	}

	return &clusterModel, nil
}

func newTFReplicationSpecsDSModel(replicationSpecs []*matlas.AdvancedReplicationSpec, projectID string) []*tfReplicationSpecModel {
	res := make([]*tfReplicationSpecModel, len(replicationSpecs))

	for i, rSpec := range replicationSpecs {
		tfRepSpec := &tfReplicationSpecModel{
			ID:        types.StringValue(rSpec.ID),
			NumShards: types.Int64Value(cast.ToInt64(rSpec.NumShards)),
			ZoneName:  types.StringValue(rSpec.ZoneName),
		}
		// regionConfigs, containerIDs := tfRegionConfigs(rSpec.RegionConfigs, projectID)
		res[i] = tfRepSpec
	}
	return res
}

// func tfRegionConfigs(advancedRegionConfig []*matlas.AdvancedRegionConfig, projectID string) {

// }

func advClusterDSConnectionStringSchemaAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"standard": schema.StringAttribute{
					Computed: true,
				},
				"standard_srv": schema.StringAttribute{
					Computed: true,
				},
				"private": schema.StringAttribute{
					Computed: true,
				},
				"private_srv": schema.StringAttribute{
					Computed: true,
				},
				"private_endpoint": schema.ListNestedAttribute{
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"connection_string": schema.StringAttribute{
								Computed: true,
							},
							"endpoints": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"endpoint_id": schema.StringAttribute{
											Computed: true,
										},
										"provider_name": schema.StringAttribute{
											Computed: true,
										},
										"region": schema.StringAttribute{
											Computed: true,
										},
									},
								},
							},
							"srv_connection_string": schema.StringAttribute{
								Computed: true,
							},
							"srv_shard_optimized_connection_string": schema.StringAttribute{
								Computed: true,
							},
							"type": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func advClusterDSReplicationSpecsSchemaAttr() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},
				"num_shards": schema.Int64Attribute{
					Computed: true,
				},
				"container_id": schema.MapAttribute{
					ElementType: types.StringType,
					Computed:    true,
				},
				"zone_name": schema.StringAttribute{
					Computed: true,
				},
				"region_configs": schema.SetNestedAttribute{
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"analytics_specs":        advancedClusterRegionConfigSpecsAttr(),
							"auto_scaling":           advancedClusterRegionConfigAutoScalingSpecsAttr(),
							"analytics_auto_scaling": advancedClusterRegionConfigAutoScalingSpecsAttr(),
							"backing_provider_name": schema.StringAttribute{
								Computed: true,
							},
							"electable_specs": advancedClusterRegionConfigSpecsAttr(),
							"priority": schema.Int64Attribute{
								Computed: true,
							},
							"provider_name": schema.StringAttribute{
								Computed: true,
							},
							"read_only_specs": advancedClusterRegionConfigSpecsAttr(),
							"region_name": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func advancedClusterRegionConfigAutoScalingSpecsAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"disk_gb_enabled": schema.BoolAttribute{
					Computed: true,
				},
				"compute_enabled": schema.BoolAttribute{
					Computed: true,
				},
				"compute_scale_down_enabled": schema.BoolAttribute{
					Computed: true,
				},
				"compute_min_instance_size": schema.StringAttribute{
					Computed: true,
				},
				"compute_max_instance_size": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func advancedClusterRegionConfigSpecsAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"disk_iops": schema.Int64Attribute{
					Optional: true,
					Computed: true,
				},
				"ebs_volume_type": schema.StringAttribute{
					Optional: true,
				},
				"instance_size": schema.BoolAttribute{
					Required: true,
				},
				"node_count": schema.BoolAttribute{
					Optional: true,
				},
			},
		},
	}
}

type tfAdvancedClusterDSModel struct {
	DiskSizeGb                   types.Float64                   `tfsdk:"disk_size_gb"`
	VersionReleaseSystem         types.String                    `tfsdk:"version_release_system"`
	ProjectID                    types.String                    `tfsdk:"project_id"`
	ID                           types.String                    `tfsdk:"id"`
	StateName                    types.String                    `tfsdk:"state_name"`
	RootCertType                 types.String                    `tfsdk:"root_cert_type"`
	EncryptionAtRestProvider     types.String                    `tfsdk:"encryption_at_rest_provider"`
	MongoDBMajorVersion          types.String                    `tfsdk:"mongo_db_major_version"`
	MongoDBVersion               types.String                    `tfsdk:"mongo_db_version"`
	Name                         types.String                    `tfsdk:"name"`
	CreateDate                   types.String                    `tfsdk:"create_date"`
	ClusterType                  types.String                    `tfsdk:"cluster_type"`
	Tags                         []*TfTagModel                   `tfsdk:"tags"`
	Labels                       []TfLabelModel                  `tfsdk:"labels"`
	BiConnectorConfig            []*TfBiConnectorConfigModel     `tfsdk:"bi_connector_config"`
	AdvancedConfiguration        []*TfAdvancedConfigurationModel `tfsdk:"advanced_configuration"`
	ConnectionStrings            []*tfConnectionStringModel      `tfsdk:"connection_strings"`
	ReplicationSpecs             []*tfReplicationSpecModel       `tfsdk:"replication_specs"`
	BackupEnabled                types.Bool                      `tfsdk:"backup_enabled"`
	Paused                       types.Bool                      `tfsdk:"paused"`
	TerminationProtectionEnabled types.Bool                      `tfsdk:"termination_protection_enabled"`
	PitEnabled                   types.Bool                      `tfsdk:"pit_enabled"`
}
