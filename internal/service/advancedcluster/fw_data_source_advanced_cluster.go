package advancedcluster

import (
	"context"
	"fmt"
	"net/http"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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
		Attributes: advancedClusterDSAttributes(),
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
		resp.Diagnostics.AddError("An error occurred while getting cluster details from Atlas", fmt.Sprintf(errorClusterAdvancedRead, clusterName, err))
		return
	}

	newClusterState, diags := newTfAdvancedClusterDSModel(ctx, conn, cluster)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTfAdvancedClusterDSModel(ctx context.Context, conn *matlas.Client, apiResp *matlas.AdvancedCluster) (*tfAdvancedClusterDSModel, diag.Diagnostics) {
	var err error
	projectID := apiResp.GroupID
	var diags diag.Diagnostics
	var d diag.Diagnostics

	clusterModel := tfAdvancedClusterDSModel{
		ID:                           conversion.StringNullIfEmpty(apiResp.ID),
		BackupEnabled:                types.BoolPointerValue(apiResp.BackupEnabled),
		ClusterType:                  conversion.StringNullIfEmpty(apiResp.ClusterType),
		CreateDate:                   conversion.StringNullIfEmpty(apiResp.CreateDate),
		DiskSizeGb:                   types.Float64PointerValue(apiResp.DiskSizeGB),
		EncryptionAtRestProvider:     conversion.StringNullIfEmpty(apiResp.EncryptionAtRestProvider),
		MongoDBMajorVersion:          conversion.StringNullIfEmpty(apiResp.MongoDBMajorVersion),
		MongoDBVersion:               conversion.StringNullIfEmpty(apiResp.MongoDBVersion),
		Name:                         conversion.StringNullIfEmpty(apiResp.Name),
		Paused:                       types.BoolPointerValue(apiResp.Paused),
		PitEnabled:                   types.BoolPointerValue(apiResp.PitEnabled),
		RootCertType:                 conversion.StringNullIfEmpty(apiResp.RootCertType),
		StateName:                    conversion.StringNullIfEmpty(apiResp.StateName),
		TerminationProtectionEnabled: types.BoolPointerValue(apiResp.TerminationProtectionEnabled),
		VersionReleaseSystem:         conversion.StringNullIfEmpty(apiResp.VersionReleaseSystem),
		ProjectID:                    conversion.StringNullIfEmpty(projectID),
	}
	clusterModel.BiConnectorConfig, d = types.ListValueFrom(ctx, TfBiConnectorConfigType, NewTfBiConnectorConfigModel(apiResp.BiConnector))
	diags.Append(d...)

	clusterModel.ConnectionStrings, d = types.ListValueFrom(ctx, tfConnectionStringType, newTfConnectionStringsModel(ctx, apiResp.ConnectionStrings))
	diags.Append(d...)

	clusterModel.Labels, d = types.SetValueFrom(ctx, TfLabelType, newTfLabelsModel(apiResp.Labels))
	diags.Append(d...)

	clusterModel.Tags, d = types.SetValueFrom(ctx, TfTagType, newTfTagsModel(&apiResp.Tags))
	diags.Append(d...)

	replicationSpecs, d := newTfReplicationSpecsDSModel(ctx, conn, apiResp.ReplicationSpecs, projectID)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}
	clusterModel.ReplicationSpecs, diags = types.SetValueFrom(ctx, tfReplicationSpecType, replicationSpecs)

	advancedConfiguration, err := newTfAdvancedConfigurationModelDSFromAtlas(ctx, conn, projectID, apiResp.Name)
	if err != nil {
		diags.AddError("An error occurred while getting advanced_configuration from Atlas", err.Error())
		return nil, diags
	}
	clusterModel.AdvancedConfiguration, diags = types.ListValueFrom(ctx, tfAdvancedConfigurationType, advancedConfiguration)
	if diags.HasError() {
		return nil, diags
	}

	return &clusterModel, nil
}

func newTfReplicationSpecsDSModel(ctx context.Context, conn *matlas.Client, replicationSpecs []*matlas.AdvancedReplicationSpec, projectID string) ([]*tfReplicationSpecModel, diag.Diagnostics) {
	res := make([]*tfReplicationSpecModel, len(replicationSpecs))
	var diags diag.Diagnostics

	for i, rSpec := range replicationSpecs {
		tfRepSpec := &tfReplicationSpecModel{
			ID:        conversion.StringNullIfEmpty(rSpec.ID),
			NumShards: types.Int64Value(cast.ToInt64(rSpec.NumShards)),
			ZoneName:  conversion.StringNullIfEmpty(rSpec.ZoneName),
		}
		regionConfigs, containerIDs, diags := getTfRegionConfigsAndContainerIDs(ctx, conn, rSpec.RegionConfigs, projectID)
		if diags.HasError() {
			return nil, diags
		}

		regionConfigsSet, diags := types.SetValueFrom(ctx, tfRegionsConfigType, regionConfigs)
		if diags.HasError() {
			return nil, diags
		}

		tfRepSpec.ContainerID = containerIDs
		tfRepSpec.RegionsConfigs = regionConfigsSet

		res[i] = tfRepSpec
	}
	return res, diags
}

func getTfRegionConfigsAndContainerIDs(ctx context.Context, conn *matlas.Client, apiObjects []*matlas.AdvancedRegionConfig, projectID string) ([]tfRegionsConfigModel, types.Map, diag.Diagnostics) {
	var tfContainersIDsMap basetypes.MapValue
	var diags diag.Diagnostics
	containerIDsMap := map[string]attr.Value{}

	tfRegionConfigs := make([]tfRegionsConfigModel, len(apiObjects))

	for i, apiObject := range apiObjects {
		tfRegionConfig, diags := newTfRegionConfig(ctx, conn, apiObject, projectID)

		tfRegionConfigs[i] = tfRegionConfig

		if apiObject.ProviderName != "TENANT" {
			containers, _, err := conn.Containers.List(ctx, projectID,
				&matlas.ContainersListOptions{ProviderName: apiObject.ProviderName})
			if err != nil {
				diags.AddError("An error occurred while getting Containers list from Atlas", err.Error())
				return nil, types.MapNull(types.StringType), diags
			}
			if result := getAdvancedClusterContainerID(containers, apiObject); result != "" {
				// Will print as "providerName:regionName" = "containerId" in terraform show
				key := fmt.Sprintf("%s:%s", apiObject.ProviderName, apiObject.RegionName)
				containerIDsMap[key] = types.StringValue(result)
			}
		}
	}

	tfContainersIDsMap, diags = types.MapValue(types.StringType, containerIDsMap)

	return tfRegionConfigs, tfContainersIDsMap, diags
}

func advancedClusterDSAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"project_id": schema.StringAttribute{
			Required: true,
		},
		"advanced_configuration": advancedConfigDSSchema(),
		"backup_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"bi_connector_config": biConnectorConfigDSSchema(),
		"cluster_type": schema.StringAttribute{
			Computed: true,
		},
		"connection_strings": connectionStringDSSchema(),
		"create_date": schema.StringAttribute{
			Computed: true,
		},
		"disk_size_gb": schema.Float64Attribute{
			Computed: true,
		},
		"encryption_at_rest_provider": schema.StringAttribute{
			Computed: true,
		},
		"labels": labelsDSSchema(),
		"tags":   tagsDSSchema(),
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
		"replication_specs": replicationSpecsDSSchema(),
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
	}
}

func connectionStringDSSchema() schema.ListNestedAttribute {
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

func replicationSpecsDSSchema() schema.SetNestedAttribute {
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
							"analytics_specs":        regionConfigSpecsDSSchema(),
							"auto_scaling":           regionConfigAutoScalingSpecsDSSchema(),
							"analytics_auto_scaling": regionConfigAutoScalingSpecsDSSchema(),
							"backing_provider_name": schema.StringAttribute{
								Computed: true,
							},
							"electable_specs": regionConfigSpecsDSSchema(),
							"priority": schema.Int64Attribute{
								Computed: true,
							},
							"provider_name": schema.StringAttribute{
								Computed: true,
							},
							"read_only_specs": regionConfigSpecsDSSchema(),
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

func regionConfigAutoScalingSpecsDSSchema() schema.ListNestedAttribute {
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

func regionConfigSpecsDSSchema() schema.ListNestedAttribute {
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
				"instance_size": schema.StringAttribute{
					Required: true,
				},
				"node_count": schema.Int64Attribute{
					Optional: true,
				},
			},
		},
	}
}

func advancedConfigDSSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"default_read_concern": schema.StringAttribute{
					Computed: true,
				},
				"default_write_concern": schema.StringAttribute{
					Computed: true,
				},
				"fail_index_key_too_long": schema.BoolAttribute{
					Computed: true,
				},
				"javascript_enabled": schema.BoolAttribute{
					Computed: true,
				},
				"minimum_enabled_tls_protocol": schema.StringAttribute{
					Computed: true,
				},
				"no_table_scan": schema.BoolAttribute{
					Computed: true,
				},
				"oplog_size_mb": schema.Int64Attribute{
					Computed: true,
				},
				"sample_size_bi_connector": schema.Int64Attribute{
					Computed: true,
				},
				"sample_refresh_interval_bi_connector": schema.Int64Attribute{
					Computed: true,
				},
				"oplog_min_retention_hours": schema.Int64Attribute{
					Computed: true,
				},
				"transaction_lifetime_limit_seconds": schema.Int64Attribute{
					Computed: true,
				},
			},
		},
	}
}

func biConnectorConfigDSSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Computed: true,
				},
				"read_preference": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func labelsDSSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Computed:           true,
		DeprecationMessage: fmt.Sprintf(constant.DeprecationParamByDateWithReplacement, "September 2024", "tags"),
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Computed: true,
				},
				"value": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func tagsDSSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"key": schema.StringAttribute{
					Computed: true,
				},
				"value": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func newTfAdvancedConfigurationModelDSFromAtlas(ctx context.Context, conn *matlas.Client, projectID, clusterName string) ([]*TfAdvancedConfigurationModel, error) {
	processArgs, _, err := conn.Clusters.GetProcessArgs(ctx, projectID, clusterName)
	if err != nil {
		return nil, err
	}

	advConfigModel := newTfAdvancedConfigurationModel(processArgs)
	return advConfigModel, err
}

type tfAdvancedClusterDSModel struct {
	DiskSizeGb                   types.Float64 `tfsdk:"disk_size_gb"`
	VersionReleaseSystem         types.String  `tfsdk:"version_release_system"`
	ProjectID                    types.String  `tfsdk:"project_id"`
	ID                           types.String  `tfsdk:"id"`
	StateName                    types.String  `tfsdk:"state_name"`
	RootCertType                 types.String  `tfsdk:"root_cert_type"`
	EncryptionAtRestProvider     types.String  `tfsdk:"encryption_at_rest_provider"`
	MongoDBMajorVersion          types.String  `tfsdk:"mongo_db_major_version"`
	MongoDBVersion               types.String  `tfsdk:"mongo_db_version"`
	Name                         types.String  `tfsdk:"name"`
	CreateDate                   types.String  `tfsdk:"create_date"`
	ClusterType                  types.String  `tfsdk:"cluster_type"`
	Tags                         types.Set     `tfsdk:"tags"`
	Labels                       types.Set     `tfsdk:"labels"`
	BiConnectorConfig            types.List    `tfsdk:"bi_connector_config"`
	AdvancedConfiguration        types.List    `tfsdk:"advanced_configuration"`
	ConnectionStrings            types.List    `tfsdk:"connection_strings"`
	ReplicationSpecs             types.Set     `tfsdk:"replication_specs"`
	BackupEnabled                types.Bool    `tfsdk:"backup_enabled"`
	Paused                       types.Bool    `tfsdk:"paused"`
	TerminationProtectionEnabled types.Bool    `tfsdk:"termination_protection_enabled"`
	PitEnabled                   types.Bool    `tfsdk:"pit_enabled"`
}
