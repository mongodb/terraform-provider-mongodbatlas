package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/planmodifiers"
)

const (
	clusterResourceName     = "cluster"
	errorClusterCreate      = "error creating MongoDB Cluster: %s"
	errorClusterRead        = "error reading MongoDB Cluster (%s): %s"
	errorClusterDelete      = "error deleting MongoDB Cluster (%s): %s"
	errorClusterUpdate      = "error updating MongoDB Cluster (%s): %s"
	errorClusterSetting     = "error setting `%s` for MongoDB Cluster (%s): %s"
	errorAdvancedConfUpdate = "error updating Advanced Configuration Option form MongoDB Cluster (%s): %s"
	errorAdvancedConfRead   = "error reading Advanced Configuration Option form MongoDB Cluster (%s): %s"
)

var _ resource.ResourceWithConfigure = &ClusterRS{}
var _ resource.ResourceWithImportState = &ClusterRS{}
var _ resource.ResourceWithModifyPlan = &ClusterRS{}

type ClusterRS struct {
	RSCommon
}

func NewClusterRS() resource.Resource {
	return &ClusterRS{
		RSCommon: RSCommon{
			resourceName: clusterResourceName,
		},
	}
}

// TODO stateUpgrader
// TODO StateFunc: formatMongoDBMajorVersion mongo_db_major_version
// https://discuss.hashicorp.com/t/is-it-possible-to-have-statefunc-like-behavior-with-the-plugin-framework/58377/2
// TODO timeouts
// TODO provider name change from TENANT -
func (r *ClusterRS) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Computed: true,
			},
			"auto_scaling_disk_gb_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"auto_scaling_compute_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_scaling_compute_scale_down_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Clusters running MongoDB FCV 4.2 or later and any new Atlas clusters of any type do not support this parameter",
			},
			"retain_backups_enabled": schema.BoolAttribute{
				// TODO make sure this is handled for null values
				Optional:    true,
				Description: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster",
			},
			"cluster_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disk_size_gb": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"encryption_at_rest_provider": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// https://discuss.hashicorp.com/t/is-it-possible-to-have-statefunc-like-behavior-with-the-plugin-framework/58377/2
			"mongo_db_major_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// TODO StateFunc: formatMongoDBMajorVersion,

			},
			"mongo_db_major_version_formatted": schema.StringAttribute{
				Computed: true,
				// TODO StateFunc: formatMongoDBMajorVersion,
			},
			"num_shards": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(1),
			},
			"cloud_backup": schema.BoolAttribute{
				Optional: true,
				Default:  booldefault.StaticBool(false),
				Validators: []validator.Bool{
					boolvalidator.ConflictsWith(path.MatchRoot("backup_enabled")),
				},
			},
			"provider_instance_size_name": schema.StringAttribute{
				Required: true,
			},
			"provider_name": schema.StringAttribute{
				Required: true,
			},
			"pit_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"backing_provider_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provider_disk_iops": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"provider_disk_type_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provider_encrypt_ebs_volume": schema.BoolAttribute{
				Optional:           true,
				Computed:           true,
				DeprecationMessage: "All EBS volumes are encrypted by default, the option to disable encryption has been removed",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"provider_encrypt_ebs_volume_flag": schema.BoolAttribute{
				Computed: true,
			},

			"provider_region_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provider_volume_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// https://github.com/mongodb/terraform-provider-mongodbatlas/pull/515
			// TODO ensure test cases cover this
			"provider_auto_scaling_compute_max_instance_size": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.ClusterAutoScalingMaxInstanceModifier(),
				},
			},
			"provider_auto_scaling_compute_min_instance_size": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.ClusterAutoScalingMinInstanceModifier(),
				},
			},
			"replication_factor": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
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
				Optional: true,
				Default:  booldefault.StaticBool(false),
			},
			"srv_address": schema.StringAttribute{
				Computed: true,
			},
			"state_name": schema.StringAttribute{
				Computed: true,
			},
			"connection_strings": schema.ListNestedAttribute{
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
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"connection_string": schema.StringAttribute{
										Computed: true,
									},
									"endpoints": schema.ListNestedAttribute{
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
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"container_id": schema.StringAttribute{
				Computed: true,
			},
			"version_release_system": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("LTS"),
				Validators: []validator.String{
					stringvalidator.OneOf("LTS", "CONTINUOUS"),
				},
			},
			"snapshot_backup_policy": schema.ListNestedAttribute{
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
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"policy_item": schema.ListNestedAttribute{
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
			},
		},
		Blocks: map[string]schema.Block{
			"advanced_configuration": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"default_read_concern": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"default_write_concern": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"fail_index_key_too_long": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"javascript_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"minimum_enabled_tls_protocol": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"no_table_scan": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"oplog_min_retention_hours": schema.Int64Attribute{
							Optional: true,
							// TODO check this during testing
							// Computed: true,
						},
						"oplog_size_mb": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"sample_refresh_interval_bi_connector": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"sample_size_bi_connector": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"transaction_lifetime_limit_seconds": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"bi_connector_config": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"read_preference": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
			},
			"labels": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"value": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
				DeprecationMessage: "this parameter is deprecated and will be removed by September 2024, please transition to tags",
			},
			"replication_specs": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"num_shards": schema.Int64Attribute{
							Required: true,
						},
						"zone_name": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("ZoneName managed by Terraform"),
						},
					},
					Blocks: map[string]schema.Block{
						"regions_config": schema.SetNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"analytics_nodes": schema.Int64Attribute{
										Optional: true,
										Computed: true,
										Default:  int64default.StaticInt64(0),
									},
									"electable_nodes": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"priority": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"read_only_nodes": schema.Int64Attribute{
										Optional: true,
										Computed: true,
										Default:  int64default.StaticInt64(0),
									},
									"region_name": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"tags": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
		Version: 1,
	}

	if s.Blocks == nil {
		s.Blocks = make(map[string]schema.Block)
	}
	s.Blocks["timeouts"] = timeouts.Block(ctx, timeouts.Opts{
		Create: true,
		Update: true,
		Delete: true,
	})

	response.Schema = s
}

func (r ClusterRS) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var currentProvider types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("provider_name"), &currentProvider)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updatedProvider types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("provider_name"), &updatedProvider)...)
	if resp.Diagnostics.HasError() {
		return
	}

	willProviderChange := currentProvider != updatedProvider
	if !willProviderChange {
		return // do nothing
	}

	willLeaveTenant := willProviderChange && currentProvider.ValueString() == "TENANT"

	if willLeaveTenant {
		// this might throw inconsistent state error
		resp.Plan.SetAttribute(ctx, path.Root("backing_provider_name"), types.StringNull())

	} else if willProviderChange {
		resp.RequiresReplace = path.Paths{path.Root("provider_name")}
	}
}

func (r *ClusterRS) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	// conn := r.client.Atlas
	// var plan tfClusterRSModel

	// response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	// if response.Diagnostics.HasError() {
	// 	return
	// }
	// createTimeout := r.CreateTimeout(ctx, plan.Timeouts)

	// plan.ID = types.StringValue("TODO")

	// // TODO: initialize and set newState

	// // set state to fully populated data
	// response.Diagnostics.Append(response.State.Set(ctx, &newState)...)
}

func (r *ClusterRS) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	conn := r.client.Atlas

	var clusterState tfClusterRSModel
	response.Diagnostics.Append(request.State.Get(ctx, &clusterState)...)
	if response.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(clusterState.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	// providerName := ids["provider_name"]

	cluster, resp, err := conn.Clusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("error during cluster READ from Atlas", fmt.Sprintf(errorClusterRead, clusterName, err.Error()))
		return
	}

	newClusterState := newTFClusterModel(cluster, &clusterState)

	// save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &newClusterState)...)
}

func newTFClusterModel(apiResp *mongodbatlas.Cluster, currState *tfClusterRSModel) tfClusterRSModel {
	clusterModel := tfClusterRSModel{
		ID:                                 currState.ID,
		ProjectID:                          currState.ProjectID,
		ClusterID:                          currState.ClusterID,
		AutoScalingComputeEnabled:          types.BoolPointerValue(apiResp.AutoScaling.Compute.Enabled),
		AutoScalingComputeScaleDownEnabled: types.BoolPointerValue(apiResp.AutoScaling.Compute.ScaleDownEnabled),

		ProviderAutoScalingComputeMinInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MinInstanceSize),
		ProviderAutoScalingComputeMaxInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MaxInstanceSize),
		BackupEnabled: types.BoolPointerValue(apiResp.BackupEnabled),
		CloudBackup:   types.BoolPointerValue(apiResp.ProviderBackupEnabled), //
		ClusterType:   types.StringValue(apiResp.ClusterType),

		DiskSizeGb:               types.Float64PointerValue(apiResp.DiskSizeGB),
		EncryptionAtRestProvider: types.StringValue(apiResp.EncryptionAtRestProvider),
		MongoDbMajorVersion:      types.StringValue(apiResp.MongoDBMajorVersion), // version formatting
		MongoDbVersion:           types.StringValue(apiResp.MongoDBVersion),
		MongoUri:                 types.StringValue(apiResp.MongoURI),

		MongoUriUpdated:              types.StringValue(apiResp.MongoURIUpdated),
		MongoUriWithOptions:          types.StringValue(apiResp.MongoURIWithOptions),
		PitEnabled:                   types.BoolPointerValue(apiResp.PitEnabled),
		Paused:                       types.BoolPointerValue(apiResp.Paused),
		SrvAddress:                   types.StringValue(apiResp.SrvAddress),
		StateName:                    types.StringValue(apiResp.StateName),
		TerminationProtectionEnabled: types.BoolPointerValue(apiResp.TerminationProtectionEnabled),
		BiConnectorConfig:            nil,
		ConnectionStrings:            nil,
		ReplicationSpecs:             nil,
		Labels:                       nil,
		Tags:                         nil,
		AdvancedConfiguration:        nil,
		SnapshotBackupPolicy:         nil,
		VersionReleaseSystem:         types.StringValue(apiResp.VersionReleaseSystem),

		// MetricThresholdConfig: newTFMetricThresholdConfigModel(apiRespConfig.MetricThreshold, currState.MetricThresholdConfig),

	}
	// connection_strings
	// numshards
	// bi connector
	// flattenProviderSettings
	// if providerName != "TENANT" {
	// processArgs - adv config
	// snapshotBackupPolicy

	return clusterModel
}

func (r *ClusterRS) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// conn := r.client.Atlas
	// var state, plan tfClusterRSModel

	// response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	// if response.Diagnostics.HasError() {
	// 	return
	// }

	// response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	// if response.Diagnostics.HasError() {
	// 	return
	// }
	// updateTimeout := r.UpdateTimeout(ctx, plan.Timeouts)

	// // save updated data into terraform state
	// response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *ClusterRS) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// conn := r.client.Atlas
	// var state tfClusterRSModel

	// response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	// if response.Diagnostics.HasError() {
	// 	return
	// }
	// deleteTimeout := r.DeleteTimeout(ctx, state.Timeouts)

	// tflog.Debug(ctx, "deleting TODO", map[string]interface{}{
	// 	"id": state.ID.ValueString(),
	// })
}

// ImportState is called when the provider must import the state of a resource instance.
// This method must return enough state so the Read method can properly refresh the full resource.
//
// If setting an attribute with the import identifier, it is recommended to use the ImportStatePassthroughID() call in this method.
func (r *ClusterRS) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

type tfClusterRSModel struct {
	AutoScalingComputeEnabled                 types.Bool    `tfsdk:"auto_scaling_compute_enabled"`
	AutoScalingComputeScaleDownEnabled        types.Bool    `tfsdk:"auto_scaling_compute_scale_down_enabled"`
	AutoScalingDiskGbEnabled                  types.Bool    `tfsdk:"auto_scaling_disk_gb_enabled"`
	BackingProviderName                       types.String  `tfsdk:"backing_provider_name"`
	BackupEnabled                             types.Bool    `tfsdk:"backup_enabled"`
	CloudBackup                               types.Bool    `tfsdk:"cloud_backup"`
	ClusterID                                 types.String  `tfsdk:"cluster_id"`
	ClusterType                               types.String  `tfsdk:"cluster_type"`
	ContainerID                               types.String  `tfsdk:"container_id"`
	DiskSizeGb                                types.Float64 `tfsdk:"disk_size_gb"`
	EncryptionAtRestProvider                  types.String  `tfsdk:"encryption_at_rest_provider"`
	ID                                        types.String  `tfsdk:"id"`
	MongoDbMajorVersion                       types.String  `tfsdk:"mongo_db_major_version"`
	MongoDbMajorVersionFormatted              types.String  `tfsdk:"mongo_db_major_version_formatted"`
	MongoDbVersion                            types.String  `tfsdk:"mongo_db_version"`
	MongoUri                                  types.String  `tfsdk:"mongo_uri"`
	MongoUriUpdated                           types.String  `tfsdk:"mongo_uri_updated"`
	MongoUriWithOptions                       types.String  `tfsdk:"mongo_uri_with_options"`
	Name                                      types.String  `tfsdk:"name"`
	NumShards                                 types.Int64   `tfsdk:"num_shards"`
	Paused                                    types.Bool    `tfsdk:"paused"`
	PitEnabled                                types.Bool    `tfsdk:"pit_enabled"`
	ProjectID                                 types.String  `tfsdk:"project_id"`
	ProviderAutoScalingComputeMaxInstanceSize types.String  `tfsdk:"provider_auto_scaling_compute_max_instance_size"`
	ProviderAutoScalingComputeMinInstanceSize types.String  `tfsdk:"provider_auto_scaling_compute_min_instance_size"`
	ProviderDiskIops                          types.Int64   `tfsdk:"provider_disk_iops"`
	ProviderDiskTypeName                      types.String  `tfsdk:"provider_disk_type_name"`
	ProviderEncryptEbsVolume                  types.Bool    `tfsdk:"provider_encrypt_ebs_volume"`
	ProviderEncryptEbsVolumeFlag              types.Bool    `tfsdk:"provider_encrypt_ebs_volume_flag"`
	ProviderInstanceSizeName                  types.String  `tfsdk:"provider_instance_size_name"`
	ProviderName                              types.String  `tfsdk:"provider_name"`
	ProviderRegionName                        types.String  `tfsdk:"provider_region_name"`
	ProviderVolumeType                        types.String  `tfsdk:"provider_volume_type"`
	ReplicationFactor                         types.Int64   `tfsdk:"replication_factor"`
	RetainBackupsEnabled                      types.Bool    `tfsdk:"retain_backups_enabled"`
	// SnapshotBackupPolicy                      types.List    `tfsdk:"snapshot_backup_policy"`
	SrvAddress                   types.String `tfsdk:"srv_address"`
	StateName                    types.String `tfsdk:"state_name"`
	TerminationProtectionEnabled types.Bool   `tfsdk:"termination_protection_enabled"`
	VersionReleaseSystem         types.String `tfsdk:"version_release_system"`

	// Computed list attributes
	ConnectionStrings    []tfConnectionStringModel     `tfsdk:"links"`
	SnapshotBackupPolicy []tfSnapshotBackupPolicyModel `tfsdk:"snapshot_backup_policy"`

	// Optional Computed
	BiConnectorConfig []tfBiConnectorConfigModel `tfsdk:"bi_connector_config"`
	// TODO remove this comment: https://discuss.hashicorp.com/t/set-list-attribute-migration-from-sdkv2-to-framework/56472/2
	ReplicationSpecs      []tfReplicationSpecModel       `tfsdk:"replication_specs"`      // SetNestedAttribute
	AdvancedConfiguration []tfAdvancedConfigurationModel `tfsdk:"advanced_configuration"` // SetNestedAttribute
	Labels                []tfLabelModel                 `tfsdk:"labels"`                 // SetNestedAttribute
	Tags                  []tfTagModel                   `tfsdk:"tags"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type tfSnapshotBackupPolicyModel struct {
	ClusterID             types.String            `tfsdk:"cluster_id"`
	ClusterName           types.String            `tfsdk:"cluster_name"`
	NextSnapshot          types.String            `tfsdk:"next_snapshot"`
	ReferenceHourOfDay    types.Int64             `tfsdk:"reference_hour_of_day"`
	ReferenceMinuteOfHour types.Int64             `tfsdk:"reference_minute_of_hour"`
	RestoreWindowDays     types.Int64             `tfsdk:"restore_window_days"`
	UpdateSnapshots       types.Bool              `tfsdk:"update_snapshots"`
	Policies              []tfSnapshotPolicyModel `tfsdk:"policies"`
}

type tfSnapshotPolicyModel struct {
	ID         types.String                `tfsdk:"id"`
	PolicyItem []tfSnapshotPolicyItemModel `tfsdk:"policy_item"`
}

type tfSnapshotPolicyItemModel struct {
	ID                types.String `tfsdk:"id"`
	FrequencyInterval types.Int64  `tfsdk:"frequency_interval"`
	FrequencyType     types.String `tfsdk:"frequency_type"`
	RetentionUnit     types.String `tfsdk:"retention_unit"`
	RetentionValue    types.Int64  `tfsdk:"retention_value"`
}

type tfTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type tfAdvancedConfigurationModel struct {
	DefaultReadConcern               types.String `tfsdk:"default_read_concern"`
	DefaultWriteConcern              types.String `tfsdk:"default_write_concern"`
	FailIndexKeyTooLong              types.String `tfsdk:"fail_index_key_too_long"`
	JavascriptEnabled                types.Bool   `tfsdk:"javascript_enabled"`
	MinimumEnabledTlsProtocol        types.String `tfsdk:"minimum_enabled_tls_protocol"`
	NoTableScan                      types.Bool   `tfsdk:"no_table_scan"`
	OlogSizeMB                       types.Int64  `tfsdk:"oplog_size_mb"`
	OplogMinRetentionHours           types.Int64  `tfsdk:"oplog_min_retention_hours"`
	SampleSizeBiConnector            types.Int64  `tfsdk:"sample_size_bi_connector"`
	SampleRefreshIntervalBiConnector types.Int64  `tfsdk:"sample_refresh_interval_bi_connector"`
	TransactionLifetimeLimitSeconds  types.Int64  `tfsdk:"transaction_lifetime_limit_seconds"`
}

type tfReplicationSpecModel struct {
	ID            types.String          `tfsdk:"id"`
	num_shards    types.Int64           `tfsdk:"num_shards"`
	RegionsConfig []tfRegionConfigModel `tfsdk:"regions_config"`
	zone_name     types.String          `tfsdk:"zone_name"`
}

type tfRegionConfigModel struct {
	RegionName     types.String `tfsdk:"region_name"`
	ElectableNodes types.Int64  `tfsdk:"electable_nodes"`
	Priority       types.Int64  `tfsdk:"priority"`
	ReadOnlyNodes  types.Int64  `tfsdk:"read_only_nodes"`
	AnalyticsNodes types.Int64  `tfsdk:"analytics_nodes"`
}

type tfBiConnectorConfigModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	ReadPreference types.String `tfsdk:"read_preference"`
}

type tfConnectionStringModel struct {
	Standard        types.String             `tfsdk:"standard"`
	StandardSrv     types.String             `tfsdk:"standard_srv"`
	Private         types.String             `tfsdk:"private"`
	PrivateSrv      types.String             `tfsdk:"private_srv"`
	PrivateEndpoint []tfPrivateEndpointModel `tfsdk:"private_endpoint"`
}

type tfPrivateEndpointModel struct {
	ConnectionString                  types.String      `tfsdk:"connection_string"`
	SrvConnectionString               types.String      `tfsdk:"srv_connection_string"`
	SrvShardOptimizedConnectionString types.String      `tfsdk:"srv_shard_optimized_connection_string"`
	EndpointType                      types.String      `tfsdk:"type"`
	Endpoints                         []tfEndpointModel `tfsdk:"endpoints"`
}

type tfEndpointModel struct {
	endpoint_id   types.String `tfsdk:"endpoint_id"`
	provider_name types.String `tfsdk:"provider_name"`
	region        types.String `tfsdk:"region"`
}
