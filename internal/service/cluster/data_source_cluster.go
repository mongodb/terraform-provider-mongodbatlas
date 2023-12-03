package cluster

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	clusterResourceName = "cluster"
)

var _ datasource.DataSource = &clusterDS{}
var _ datasource.DataSourceWithConfigure = &clusterDS{}

func DataSource() datasource.DataSource {
	return &clusterDS{
		DSCommon: config.DSCommon{
			DataSourceName: clusterResourceName,
		},
	}
}

type clusterDS struct {
	config.DSCommon
}

func (d *clusterDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: clusterDSAttributes(),
	}
}

func clusterDSAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed: true,
		},
		"project_id": schema.StringAttribute{
			Required: true,
		},
		"name": schema.StringAttribute{
			Required: true,
		},
		"advanced_configuration": clusterDSAdvancedConfigurationSchemaAttr(),
		"auto_scaling_disk_gb_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"auto_scaling_compute_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"auto_scaling_compute_scale_down_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"backup_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"bi_connector_config": clusterDSBiConnectorConfigSchemaAttr(),
		"cluster_type": schema.StringAttribute{
			Computed: true,
		},
		"connection_strings": clusterDSConnectionStringSchemaAttr(),
		"disk_size_gb": schema.Float64Attribute{
			Computed: true,
		},
		"encryption_at_rest_provider": schema.StringAttribute{
			Computed: true,
		},
		"mongo_db_major_version": schema.StringAttribute{
			Computed: true,
		},
		"num_shards": schema.Int64Attribute{
			Computed: true,
		},
		"pit_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"provider_backup_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"provider_instance_size_name": schema.StringAttribute{
			Computed: true,
		},
		"provider_name": schema.StringAttribute{
			Computed: true,
		},
		"backing_provider_name": schema.StringAttribute{
			Computed: true,
		},
		"provider_disk_iops": schema.Int64Attribute{
			Computed: true,
		},
		"provider_disk_type_name": schema.StringAttribute{
			Computed: true,
		},
		"provider_encrypt_ebs_volume": schema.BoolAttribute{
			Computed: true,
		},
		"provider_encrypt_ebs_volume_flag": schema.BoolAttribute{
			Computed: true,
		},
		"provider_region_name": schema.StringAttribute{
			Computed: true,
		},
		"provider_volume_type": schema.StringAttribute{
			Computed: true,
		},
		"provider_auto_scaling_compute_max_instance_size": schema.StringAttribute{
			Computed: true,
		},
		"provider_auto_scaling_compute_min_instance_size": schema.StringAttribute{
			Computed: true,
		},
		"replication_factor": schema.Int64Attribute{
			Computed: true,
		},
		"replication_specs": clusterDSReplicationSpecsSchemaAttr(),
		"mongo_db_version": schema.StringAttribute{
			Computed: true,
		},
		"mongo_uri": schema.StringAttribute{
			Computed: true,
		},
		"mongo_uri_updated": schema.StringAttribute{
			Computed: true,
		},
		"mongo_uri_with_options": schema.StringAttribute{
			Computed: true,
		},
		"paused": schema.BoolAttribute{
			Computed: true,
		},
		"srv_address": schema.StringAttribute{
			Computed: true,
		},
		"state_name": schema.StringAttribute{
			Computed: true,
		},
		"labels":                 clusterDSLabelsSchemaAttr(),
		"tags":                   clusterDSTagsSchemaAttr(),
		"snapshot_backup_policy": clusterDSSnapshotBackupPolicySchemaAttr(),
		"termination_protection_enabled": schema.BoolAttribute{
			Computed: true,
		},
		"container_id": schema.StringAttribute{
			Computed: true,
		},
		"version_release_system": schema.StringAttribute{
			Computed: true,
		},
	}
}

func clusterDSAdvancedConfigurationSchemaAttr() schema.ListNestedAttribute {
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

func clusterDSSnapshotBackupPolicySchemaAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"cluster_id": schema.StringAttribute{
					Computed: true,
				},
				"cluster_name": schema.StringAttribute{
					Computed: true,
				},
				"next_snapshot": schema.StringAttribute{
					Computed: true,
				},
				"reference_hour_of_day": schema.Int64Attribute{
					Computed: true,
				},
				"reference_minute_of_hour": schema.Int64Attribute{
					Computed: true,
				},
				"restore_window_days": schema.Int64Attribute{
					Computed: true,
				},
				"update_snapshots": schema.BoolAttribute{
					Computed: true,
				},
				"policies": schema.ListNestedAttribute{
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								Computed: true,
							},
							"policy_item": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											Computed: true,
										},
										"frequency_interval": schema.Int64Attribute{
											Computed: true,
										},
										"frequency_type": schema.StringAttribute{
											Computed: true,
										},
										"retention_unit": schema.StringAttribute{
											Computed: true,
										},
										"retention_value": schema.Int64Attribute{
											Computed: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func clusterDSTagsSchemaAttr() schema.SetNestedAttribute {
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

func clusterDSLabelsSchemaAttr() schema.SetNestedAttribute {
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

func clusterDSReplicationSpecsSchemaAttr() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},
				"num_shards": schema.Int64Attribute{
					Computed: true,
				},
				"regions_config": schema.SetNestedAttribute{
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"region_name": schema.StringAttribute{
								Computed: true,
							},
							"electable_nodes": schema.Int64Attribute{
								Computed: true,
							},
							"priority": schema.Int64Attribute{
								Computed: true,
							},
							"read_only_nodes": schema.Int64Attribute{
								Computed: true,
							},
							"analytics_nodes": schema.Int64Attribute{
								Computed: true,
							},
						},
					},
				},
				"zone_name": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func clusterDSBiConnectorConfigSchemaAttr() schema.ListNestedAttribute {
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

func clusterDSConnectionStringSchemaAttr() schema.ListNestedAttribute {
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
				"aws_private_link": schema.MapAttribute{
					ElementType: types.StringType,
					Computed:    true,
				},
				"aws_private_link_srv": schema.MapAttribute{
					ElementType: types.StringType,
					Computed:    true,
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

func (d *clusterDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Client.Atlas
	var clusterConfig tfClusterDSModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &clusterConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := clusterConfig.ProjectID.ValueString()
	clusterName := clusterConfig.Name.ValueString()

	cluster, response, err := conn.Clusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError("cluster not found in Atlas", fmt.Sprintf(errorClusterRead, clusterName, err.Error()))
			return
		}
		resp.Diagnostics.AddError("error in getting cluster details from Atlas", fmt.Sprintf(errorClusterRead, clusterName, err.Error()))
		return
	}

	newClusterState, err := newTFClusterDSModel(ctx, conn, cluster)
	if err != nil {
		resp.Diagnostics.AddError("error in getting cluster details from Atlas", fmt.Sprintf(errorClusterRead, clusterName, err.Error()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTFClusterDSModel(ctx context.Context, conn *matlas.Client, apiResp *matlas.Cluster) (*tfClusterDSModel, error) {
	var err error
	projectID := apiResp.GroupID
	clusterName := apiResp.Name

	clusterModel := tfClusterDSModel{
		AutoScalingComputeEnabled:                 types.BoolPointerValue(apiResp.AutoScaling.Compute.Enabled),
		AutoScalingComputeScaleDownEnabled:        types.BoolPointerValue(apiResp.AutoScaling.Compute.ScaleDownEnabled),
		ProviderAutoScalingComputeMinInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MinInstanceSize),
		ProviderAutoScalingComputeMaxInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MaxInstanceSize),
		AutoScalingDiskGbEnabled:                  types.BoolPointerValue(apiResp.AutoScaling.DiskGBEnabled),
		BackupEnabled:                             types.BoolPointerValue(apiResp.BackupEnabled),
		PitEnabled:                                types.BoolPointerValue(apiResp.PitEnabled),
		ProviderBackupEnabled:                     types.BoolPointerValue(apiResp.ProviderBackupEnabled),
		ClusterType:                               types.StringValue(apiResp.ClusterType),
		ConnectionStrings:                         newTFConnectionStringsModelDS(ctx, apiResp.ConnectionStrings),
		DiskSizeGb:                                types.Float64PointerValue(apiResp.DiskSizeGB),
		EncryptionAtRestProvider:                  types.StringValue(apiResp.EncryptionAtRestProvider),
		MongoDBMajorVersion:                       types.StringValue(apiResp.MongoDBMajorVersion),
		MongoDBVersion:                            types.StringValue(apiResp.MongoDBVersion),
		MongoURI:                                  types.StringValue(apiResp.MongoURI),
		MongoURIUpdated:                           types.StringValue(apiResp.MongoURIUpdated),
		MongoURIWithOptions:                       types.StringValue(apiResp.MongoURIWithOptions),
		Paused:                                    types.BoolPointerValue(apiResp.Paused),
		SrvAddress:                                types.StringValue(apiResp.SrvAddress),
		StateName:                                 types.StringValue(apiResp.StateName),
		BiConnectorConfig:                         newTFBiConnectorConfigModel(apiResp.BiConnector),
		ReplicationFactor:                         types.Int64PointerValue(apiResp.ReplicationFactor),
		ReplicationSpecs:                          newTFReplicationSpecsModel(apiResp.ReplicationSpecs),
		Labels:                                    newTFLabelsModel(apiResp.Labels),
		Tags:                                      newTFTagsModel(apiResp.Tags),
		TerminationProtectionEnabled:              types.BoolPointerValue(apiResp.TerminationProtectionEnabled),
		VersionReleaseSystem:                      types.StringValue(apiResp.VersionReleaseSystem),
		ProjectID:                                 types.StringValue(projectID),
		Name:                                      types.StringValue(clusterName),
		ID:                                        types.StringValue(apiResp.ID),
	}

	// Avoid Global Cluster issues. (NumShards is not present in Global Clusters)
	if numShards := apiResp.NumShards; numShards != nil {
		clusterModel.NumShards = types.Int64PointerValue(numShards)
	}

	if apiResp.ProviderSettings != nil {
		setTFProviderSettingsDS(&clusterModel, apiResp.ProviderSettings)

		if pName := apiResp.ProviderSettings.ProviderName; pName != "TENANT" {
			containers, _, err := conn.Containers.List(ctx, projectID,
				&matlas.ContainersListOptions{ProviderName: pName})
			if err != nil {
				return nil, fmt.Errorf(errorClusterRead, clusterName, err)
			}

			clusterModel.ContainerID = types.StringValue(getContainerID(containers, apiResp))
		}
	}

	clusterModel.AdvancedConfiguration, err = newTFAdvancedConfigurationModelDSFromAtlas(ctx, conn, projectID, apiResp.Name)
	if err != nil {
		return nil, err
	}

	clusterModel.SnapshotBackupPolicy, err = newTFSnapshotBackupPolicyModelFromAtlas(ctx, conn, projectID, clusterName)
	if err != nil {
		return nil, err
	}

	return &clusterModel, nil
}

func newTFTagsModel(tags *[]*matlas.Tag) []*tfTagModel {
	res := make([]*tfTagModel, len(*tags))

	for i, v := range *tags {
		res[i] = &tfTagModel{
			Key:   types.StringValue(v.Key),
			Value: types.StringValue(v.Value),
		}
	}

	return res
}

func newTFReplicationSpecsModel(replicationSpecs []matlas.ReplicationSpec) []*tfReplicationSpecModel {
	res := make([]*tfReplicationSpecModel, len(replicationSpecs))

	for i, rSpec := range replicationSpecs {
		res[i] = &tfReplicationSpecModel{
			ID:            types.StringValue(rSpec.ID),
			NumShards:     types.Int64PointerValue(rSpec.NumShards),
			ZoneName:      types.StringValue(rSpec.ZoneName),
			RegionsConfig: newTFRegionsConfigModel(rSpec.RegionsConfig),
		}
	}
	return res
}

func newTFRegionsConfigModel(regionsConfig map[string]matlas.RegionsConfig) []tfRegionConfigModel {
	res := []tfRegionConfigModel{}

	for regionName, regionConfig := range regionsConfig {
		region := tfRegionConfigModel{
			RegionName:     types.StringValue(regionName),
			Priority:       types.Int64PointerValue(regionConfig.Priority),
			AnalyticsNodes: types.Int64PointerValue(regionConfig.AnalyticsNodes),
			ElectableNodes: types.Int64PointerValue(regionConfig.ElectableNodes),
			ReadOnlyNodes:  types.Int64PointerValue(regionConfig.ReadOnlyNodes),
		}
		res = append(res, region)
	}
	return res
}

func newTFSnapshotBackupPolicyModelFromAtlas(ctx context.Context, conn *matlas.Client, projectID, clusterName string) ([]*tfSnapshotBackupPolicyModel, error) {
	res := []*tfSnapshotBackupPolicyModel{}

	backupPolicy, response, err := conn.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)

	if err != nil {
		if response.StatusCode == http.StatusNotFound ||
			strings.Contains(err.Error(), "BACKUP_CONFIG_NOT_FOUND") ||
			strings.Contains(err.Error(), "Not Found") ||
			strings.Contains(err.Error(), "404") {
			return res, nil
		}

		return nil, fmt.Errorf(ErrorSnapshotBackupPolicyRead, clusterName, err)
	}

	res = append(res, &tfSnapshotBackupPolicyModel{
		ClusterID:             types.StringValue(backupPolicy.ClusterID),
		ClusterName:           types.StringValue(backupPolicy.ClusterName),
		NextSnapshot:          types.StringValue(backupPolicy.NextSnapshot),
		ReferenceHourOfDay:    types.Int64PointerValue(backupPolicy.ReferenceHourOfDay),
		ReferenceMinuteOfHour: types.Int64PointerValue(backupPolicy.ReferenceMinuteOfHour),
		RestoreWindowDays:     types.Int64PointerValue(backupPolicy.RestoreWindowDays),
		UpdateSnapshots:       types.BoolPointerValue(backupPolicy.UpdateSnapshots),
		Policies:              newTFSnapshotPolicyModel(ctx, backupPolicy.Policies),
	})
	return res, nil
}

func newTFSnapshotPolicyModel(ctx context.Context, policies []matlas.Policy) types.List {
	res := make([]tfSnapshotPolicyModel, len(policies))

	for i, pe := range policies {
		res[i] = tfSnapshotPolicyModel{
			ID:         types.StringValue(pe.ID),
			PolicyItem: newTFSnapshotPolicyItemModel(ctx, pe.PolicyItems),
		}
	}
	s, _ := types.ListValueFrom(ctx, tfSnapshotPolicyType, res)
	return s
}

func newTFSnapshotPolicyItemModel(ctx context.Context, policyItems []matlas.PolicyItem) types.List {
	res := make([]tfSnapshotPolicyItemModel, len(policyItems))

	for i, pe := range policyItems {
		res[i] = tfSnapshotPolicyItemModel{
			ID:                types.StringValue(pe.ID),
			FrequencyInterval: types.Int64Value(cast.ToInt64(pe.FrequencyInterval)),
			FrequencyType:     types.StringValue(pe.FrequencyType),
			RetentionUnit:     types.StringValue(pe.RetentionUnit),
			RetentionValue:    types.Int64Value(cast.ToInt64(pe.RetentionValue)),
		}
	}
	s, _ := types.ListValueFrom(ctx, tfSnapshotPolicyItemType, res)
	return s
}

func newTFAdvancedConfigurationModelDSFromAtlas(ctx context.Context, conn *matlas.Client, projectID, clusterName string) ([]*tfAdvancedConfigurationModel, error) {
	processArgs, _, err := conn.Clusters.GetProcessArgs(ctx, projectID, clusterName)
	if err != nil {
		return nil, err
	}

	advConfigModel := newTfAdvancedConfigurationModel(ctx, processArgs)
	return advConfigModel, err
}

func newTfAdvancedConfigurationModel(ctx context.Context, p *matlas.ProcessArgs) []*tfAdvancedConfigurationModel {
	res := []*tfAdvancedConfigurationModel{
		{
			DefaultReadConcern:               types.StringValue(p.DefaultReadConcern),
			DefaultWriteConcern:              types.StringValue(p.DefaultWriteConcern),
			FailIndexKeyTooLong:              types.BoolPointerValue(p.FailIndexKeyTooLong),
			JavascriptEnabled:                types.BoolPointerValue(p.JavascriptEnabled),
			MinimumEnabledTLSProtocol:        types.StringValue(p.MinimumEnabledTLSProtocol),
			NoTableScan:                      types.BoolPointerValue(p.NoTableScan),
			OplogSizeMB:                      types.Int64PointerValue(p.OplogSizeMB),
			OplogMinRetentionHours:           types.Int64Value(cast.ToInt64(p.OplogMinRetentionHours)),
			SampleSizeBiConnector:            types.Int64PointerValue(p.SampleSizeBIConnector),
			SampleRefreshIntervalBiConnector: types.Int64PointerValue(p.SampleRefreshIntervalBIConnector),
			TransactionLifetimeLimitSeconds:  types.Int64PointerValue(p.TransactionLifetimeLimitSeconds),
		},
	}
	return res
}

func newTFConnectionStringsModelDS(ctx context.Context, connString *matlas.ConnectionStrings) []*tfConnectionStringDSModel {
	res := []*tfConnectionStringDSModel{}

	if connString != nil {
		res = append(res, &tfConnectionStringDSModel{
			Standard:          types.StringValue(connString.Standard),
			StandardSrv:       types.StringValue(connString.StandardSrv),
			Private:           types.StringValue(connString.Private),
			PrivateSrv:        types.StringValue(connString.PrivateSrv),
			PrivateEndpoint:   newTFPrivateEndpointModel(ctx, connString.PrivateEndpoint),
			AwsPrivateLink:    newTFAwsPrivateLinkMap(connString.AwsPrivateLink),
			AwsPrivateLinkSrv: newTFAwsPrivateLinkMap(connString.AwsPrivateLinkSrv),
		})
	}
	return res
}

func newTFAwsPrivateLinkMap(mp map[string]string) basetypes.MapValue {
	mapValue, _ := types.MapValue(types.StringType, map[string]attr.Value{})
	return mapValue
}

func newTFPrivateEndpointModel(ctx context.Context, privateEndpoints []matlas.PrivateEndpoint) types.List {
	res := make([]tfPrivateEndpointModel, len(privateEndpoints))

	for i, pe := range privateEndpoints {
		res[i] = tfPrivateEndpointModel{
			ConnectionString:                  types.StringValue(pe.ConnectionString),
			SrvConnectionString:               types.StringValue(pe.SRVConnectionString),
			SrvShardOptimizedConnectionString: types.StringValue(pe.SRVShardOptimizedConnectionString),
			EndpointType:                      types.StringValue(pe.Type),
			Endpoints:                         newTFEndpointModel(ctx, pe.Endpoints),
		}
	}
	s, _ := types.ListValueFrom(ctx, tfPrivateEndpointType, res)
	return s
}

func newTFEndpointModel(ctx context.Context, endpoints []matlas.Endpoint) types.List {
	res := make([]tfEndpointModel, len(endpoints))

	for i, e := range endpoints {
		res[i] = tfEndpointModel{
			Region:       types.StringValue(e.Region),
			ProviderName: types.StringValue(e.ProviderName),
			EndpointID:   types.StringValue(e.EndpointID),
		}
	}
	s, _ := types.ListValueFrom(ctx, tfEndpointType, res)
	return s
}

func newTFLabelsModel(labels []matlas.Label) []tfLabelModel {
	out := make([]tfLabelModel, len(labels))
	for i, v := range labels {
		out[i] = tfLabelModel{
			Key:   types.StringValue(v.Key),
			Value: types.StringValue(v.Value),
		}
	}

	return out
}

func setTFProviderSettingsDS(clusterModel *tfClusterDSModel, settings *matlas.ProviderSettings) {
	if settings.ProviderName == "TENANT" {
		clusterModel.BackingProviderName = types.StringValue(settings.BackingProviderName)
	}

	if settings.DiskIOPS != nil && *settings.DiskIOPS != 0 {
		clusterModel.ProviderDiskIops = types.Int64PointerValue(settings.DiskIOPS)
	}
	if settings.EncryptEBSVolume != nil {
		clusterModel.ProviderEncryptEbsVolumeFlag = types.BoolPointerValue(settings.EncryptEBSVolume)
	}
	clusterModel.ProviderDiskTypeName = types.StringValue(settings.DiskTypeName)
	clusterModel.ProviderInstanceSizeName = types.StringValue(settings.InstanceSizeName)
	clusterModel.ProviderName = types.StringValue(settings.ProviderName)
	clusterModel.ProviderRegionName = types.StringValue(settings.RegionName)
	clusterModel.ProviderVolumeType = types.StringValue(settings.VolumeType)
}

func newTFBiConnectorConfigModel(biConnector *matlas.BiConnector) []*tfBiConnectorConfigModel {
	if biConnector == nil {
		return []*tfBiConnectorConfigModel{}
	}

	return []*tfBiConnectorConfigModel{
		{
			Enabled:        types.BoolPointerValue(biConnector.Enabled),
			ReadPreference: types.StringValue(biConnector.ReadPreference),
		},
	}
}

type tfClusterDSModel struct {
	DiskSizeGb                                types.Float64                   `tfsdk:"disk_size_gb"`
	ProviderAutoScalingComputeMaxInstanceSize types.String                    `tfsdk:"provider_auto_scaling_compute_max_instance_size"`
	EncryptionAtRestProvider                  types.String                    `tfsdk:"encryption_at_rest_provider"`
	VersionReleaseSystem                      types.String                    `tfsdk:"version_release_system"`
	StateName                                 types.String                    `tfsdk:"state_name"`
	ClusterType                               types.String                    `tfsdk:"cluster_type"`
	ContainerID                               types.String                    `tfsdk:"container_id"`
	SrvAddress                                types.String                    `tfsdk:"srv_address"`
	ProviderVolumeType                        types.String                    `tfsdk:"provider_volume_type"`
	ID                                        types.String                    `tfsdk:"id"`
	MongoDBMajorVersion                       types.String                    `tfsdk:"mongo_db_major_version"`
	MongoDBVersion                            types.String                    `tfsdk:"mongo_db_version"`
	MongoURI                                  types.String                    `tfsdk:"mongo_uri"`
	ProviderAutoScalingComputeMinInstanceSize types.String                    `tfsdk:"provider_auto_scaling_compute_min_instance_size"`
	MongoURIWithOptions                       types.String                    `tfsdk:"mongo_uri_with_options"`
	Name                                      types.String                    `tfsdk:"name"`
	ProviderRegionName                        types.String                    `tfsdk:"provider_region_name"`
	ProviderName                              types.String                    `tfsdk:"provider_name"`
	ProviderInstanceSizeName                  types.String                    `tfsdk:"provider_instance_size_name"`
	ProjectID                                 types.String                    `tfsdk:"project_id"`
	ProviderDiskTypeName                      types.String                    `tfsdk:"provider_disk_type_name"`
	MongoURIUpdated                           types.String                    `tfsdk:"mongo_uri_updated"`
	BackingProviderName                       types.String                    `tfsdk:"backing_provider_name"`
	ConnectionStrings                         []*tfConnectionStringDSModel    `tfsdk:"connection_strings"`
	SnapshotBackupPolicy                      []*tfSnapshotBackupPolicyModel  `tfsdk:"snapshot_backup_policy"`
	AdvancedConfiguration                     []*tfAdvancedConfigurationModel `tfsdk:"advanced_configuration"`
	ReplicationSpecs                          []*tfReplicationSpecModel       `tfsdk:"replication_specs"`
	Tags                                      []*tfTagModel                   `tfsdk:"tags"`
	Labels                                    []tfLabelModel                  `tfsdk:"labels"`
	BiConnectorConfig                         []*tfBiConnectorConfigModel     `tfsdk:"bi_connector_config"`
	ProviderDiskIops                          types.Int64                     `tfsdk:"provider_disk_iops"`
	NumShards                                 types.Int64                     `tfsdk:"num_shards"`
	ReplicationFactor                         types.Int64                     `tfsdk:"replication_factor"`
	Paused                                    types.Bool                      `tfsdk:"paused"`
	ProviderEncryptEbsVolume                  types.Bool                      `tfsdk:"provider_encrypt_ebs_volume"`
	ProviderEncryptEbsVolumeFlag              types.Bool                      `tfsdk:"provider_encrypt_ebs_volume_flag"`
	AutoScalingComputeEnabled                 types.Bool                      `tfsdk:"auto_scaling_compute_enabled"`
	ProviderBackupEnabled                     types.Bool                      `tfsdk:"provider_backup_enabled"`
	AutoScalingDiskGbEnabled                  types.Bool                      `tfsdk:"auto_scaling_disk_gb_enabled"`
	PitEnabled                                types.Bool                      `tfsdk:"pit_enabled"`
	BackupEnabled                             types.Bool                      `tfsdk:"backup_enabled"`
	TerminationProtectionEnabled              types.Bool                      `tfsdk:"termination_protection_enabled"`
	AutoScalingComputeScaleDownEnabled        types.Bool                      `tfsdk:"auto_scaling_compute_scale_down_enabled"`
}
