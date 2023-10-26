package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/planmodifiers"
)

const (
	clusterResourceName      = "cluster"
	errorClusterCreate       = "error creating MongoDB Cluster: %s"
	errorClusterRead         = "error reading MongoDB Cluster (%s): %s"
	errorClusterDelete       = "error deleting MongoDB Cluster (%s): %s"
	errorClusterUpdate       = "error updating MongoDB Cluster (%s): %s"
	errorClusterSetting      = "error setting `%s` for MongoDB Cluster (%s): %s"
	errorAdvancedConfUpdate  = "error updating Advanced Configuration Option form MongoDB Cluster (%s): %s"
	errorAdvancedConfRead    = "error reading Advanced Configuration Option form MongoDB Cluster (%s): %s"
	errorInvalidCreateValues = "Invalid values. Unable to CREATE cluster"
	// errorSnapshotBackupPolicyRead    = "error getting a Cloud Provider Snapshot Backup Policy for the cluster(%s): %s"
	// errorSnapshotBackupPolicySetting = "error setting `%s` for Cloud Provider Snapshot Backup Policy(%s): %s"
	defaultTimeout = (3 * time.Hour)
)

var _ resource.ResourceWithConfigure = &ClusterRS{}
var _ resource.ResourceWithImportState = &ClusterRS{}
var _ resource.ResourceWithModifyPlan = &ClusterRS{}

var defaultLabel = matlas.Label{Key: "Infrastructure Tool", Value: "MongoDB Atlas Terraform Provider"}

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
// TODO test labels
// TODO bug - tenantDisksize
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
			"connection_strings": clusterRSConnectionStringSchemaAttribute(),
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
			"snapshot_backup_policy": clusterRSSnapshotBackupPolicySchemaAttribute(),
		},
		Blocks: map[string]schema.Block{
			"advanced_configuration": clusterRSAdvancedConfigurationSchemaBlock(),
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
			"replication_specs": clusterRSReplicationSpecsSchemaBlock(),
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

func clusterRSReplicationSpecsSchemaBlock() schema.SetNestedBlock {
	return schema.SetNestedBlock{
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
	}
}

func clusterRSAdvancedConfigurationSchemaBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
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
	}
}

func clusterRSSnapshotBackupPolicySchemaAttribute() schema.ListNestedAttribute {
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
	}
}

func clusterRSConnectionStringSchemaAttribute() schema.ListNestedAttribute {
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
	}
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

func (r *ClusterRS) Create(ctx context.Context, req resource.CreateRequest, response *resource.CreateResponse) {
	conn := r.client.Atlas

	var plan tfClusterRSModel
	var autoScaling *matlas.AutoScaling

	response.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	projectID := plan.ProjectID.ValueString()
	providerName := plan.ProviderName.ValueString()

	computeEnabled := plan.AutoScalingComputeEnabled.ValueBool()
	scaleDownEnabled := plan.AutoScalingComputeScaleDownEnabled.ValueBool()
	diskGBEnabled := plan.AutoScalingDiskGBEnabled.ValueBoolPointer()

	validateClusterConfig(&plan, response)
	if response.Diagnostics.HasError() {
		return
	}

	if providerName != "TENANT" {
		autoScaling = &matlas.AutoScaling{
			DiskGBEnabled: diskGBEnabled,
			Compute: &matlas.Compute{
				Enabled:          &computeEnabled,
				ScaleDownEnabled: &scaleDownEnabled,
			},
		}
	}

	providerSettings, err := newAtlasProviderSetting(&plan)
	if err != nil {
		response.Diagnostics.AddError("Unable to CREATE cluster. Error in translating provider settings", fmt.Sprintf(errorClusterCreate, err))
		return
	}

	replicationSpecs, err := newAtlasReplicationSpecs(&plan)
	if err != nil {
		response.Diagnostics.AddError("Unable to CREATE cluster. Error in translating replication specs", fmt.Sprintf(errorClusterCreate, err))
		return
	}

	clusterRequest := &matlas.Cluster{
		Name:                     plan.Name.ValueString(),
		EncryptionAtRestProvider: plan.EncryptionAtRestProvider.ValueString(),
		ClusterType:              plan.ClusterID.ValueString(),
		BackupEnabled:            plan.BackupEnabled.ValueBoolPointer(),
		PitEnabled:               plan.PitEnabled.ValueBoolPointer(),
		AutoScaling:              autoScaling,
		ProviderSettings:         providerSettings,
		ReplicationSpecs:         replicationSpecs,
	}

	if cloudBackup := plan.CloudBackup; !cloudBackup.IsNull() {
		clusterRequest.ProviderBackupEnabled = cloudBackup.ValueBoolPointer()
	}

	if biConnector := plan.BiConnectorConfig; len(biConnector) > 0 {
		biConnector, err := newAtlasBiConnectorConfig(&plan)
		if err != nil {
			response.Diagnostics.AddError("Unable to CREATE cluster. Error in translating bi_connector_config", fmt.Sprintf(errorClusterCreate, err))
			return
		}
		clusterRequest.BiConnector = biConnector
	}

	labels := newAtlasLabels(plan.Labels)
	if containsLabelOrKey(labels, defaultLabel) {
		response.Diagnostics.AddError("Unable to CREATE cluster. Incorrect labels", "you should not set `Infrastructure Tool` label, it is used for internal purposes")
		return
	}
	labels = append(labels, defaultLabel)
	clusterRequest.Labels = labels

	if tags := plan.Tags; len(tags) > 0 {
		tagsSlice := newAtlasTags(tags)
		clusterRequest.Tags = &tagsSlice
	}

	if v := plan.DiskSizeGb; !v.IsNull() {
		clusterRequest.DiskSizeGB = v.ValueFloat64Pointer()
	}

	tenantDisksize := pointy.Float64(0)
	if cast.ToFloat64(tenantDisksize) != 0 {
		clusterRequest.DiskSizeGB = tenantDisksize
	}
	if v := plan.MongoDBMajorVersion; !v.IsNull() {
		clusterRequest.MongoDBMajorVersion = formatMongoDBMajorVersion(v.ValueString())
	}
	if v := plan.ReplicationFactor; !v.IsNull() {
		clusterRequest.ReplicationFactor = v.ValueInt64Pointer()
	}
	if v := plan.NumShards; !v.IsNull() {
		clusterRequest.NumShards = v.ValueInt64Pointer()
	}
	if v := plan.TerminationProtectionEnabled; !v.IsNull() {
		clusterRequest.TerminationProtectionEnabled = v.ValueBoolPointer()
	}

	if v := plan.VersionReleaseSystem; !v.IsNull() {
		clusterRequest.VersionReleaseSystem = v.ValueString()
	}

	cluster, _, err := conn.Clusters.Create(ctx, projectID, clusterRequest)
	if err != nil {
		response.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", "you should not set `Infrastructure Tool` label, it is used for internal purposes")
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, defaultTimeout)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterRefreshFunc(ctx, plan.Name.ValueString(), projectID, conn),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		response.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", "you should not set `Infrastructure Tool` label, it is used for internal purposes")
		return
	}

	/*
		So far, the cluster has created correctly, so we need to set up
		the advanced configuration option to attach it
	*/
	// ac, ok := d.GetOk("advanced_configuration")
	if ac := plan.AdvancedConfiguration; len(ac) > 0 {
		advancedConfReq := newAtlasProcessArgs(&plan.AdvancedConfiguration[0])

		_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, cluster.Name, advancedConfReq)
		if err != nil {
			response.Diagnostics.AddError("Unable to CREATE cluster. Error when updating advanced_configuration in Atlas", fmt.Sprintf(errorClusterCreate, err))
			return
		}
	}

	// To pause a cluster
	if v := plan.Paused.ValueBool(); v {
		clusterRequest = &matlas.Cluster{
			Paused: pointy.Bool(v),
		}

		_, _, err = updateCluster(ctx, conn, clusterRequest, projectID, cluster.Name, timeout)
		if err != nil {
			response.Diagnostics.AddError("Unable to CREATE cluster. Error when attempting to pause cluster in Atlas", fmt.Sprintf(errorClusterCreate, err))
			return
		}
	}

	plan.ID = types.StringValue(encodeStateID(map[string]string{
		"cluster_id":    cluster.ID,
		"project_id":    projectID,
		"cluster_name":  cluster.Name,
		"provider_name": providerName,
	}))

	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *ClusterRS) Read(ctx context.Context, req resource.ReadRequest, response *resource.ReadResponse) {
	conn := r.client.Atlas

	var isImport bool
	var clusterState tfClusterRSModel
	response.Diagnostics.Append(req.State.Get(ctx, &clusterState)...)
	if response.Diagnostics.HasError() {
		return
	}

	// Use the ID only with the IMPORT operation
	if clusterState.ID.ValueString() != "" && (clusterState.ClusterID.ValueString() == "") {
		isImport = true
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

	newClusterState, err := newTFClusterModel(ctx, conn, isImport, cluster, &clusterState)
	if err != nil {
		response.Diagnostics.AddError("error during cluster READ when translating to model", fmt.Sprintf(errorClusterRead, clusterName, err.Error()))
		return
	}

	// save updated data into Terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &newClusterState)...)
}

func (r *ClusterRS) Update(ctx context.Context, req resource.UpdateRequest, response *resource.UpdateResponse) {
	// conn := r.client.Atlas
	var state, plan tfClusterRSModel

	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	// timeout, diags := plan.Timeouts.Update(ctx, defaultTimeout)
	// response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	// save updated data into terraform state
	// response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *ClusterRS) Delete(ctx context.Context, req resource.DeleteRequest, response *resource.DeleteResponse) {
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
func (r *ClusterRS) ImportState(ctx context.Context, req resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, response)
}

func newTFClusterModel(ctx context.Context, conn *matlas.Client, isImport bool, apiResp *matlas.Cluster, currState *tfClusterRSModel) (*tfClusterRSModel, error) {
	var err error
	projectID := apiResp.GroupID
	clusterName := apiResp.Name

	clusterModel := tfClusterRSModel{
		// ID:                                 currState.ID,
		// ClusterID:                          currState.ClusterID,
		ProjectID:                          types.StringValue(projectID),
		Name:                               types.StringValue(clusterName),
		ProviderName:                       types.StringValue(apiResp.ProviderSettings.ProviderName),
		AutoScalingComputeEnabled:          types.BoolPointerValue(apiResp.AutoScaling.Compute.Enabled),
		AutoScalingComputeScaleDownEnabled: types.BoolPointerValue(apiResp.AutoScaling.Compute.ScaleDownEnabled),
		ProviderAutoScalingComputeMinInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MinInstanceSize),
		ProviderAutoScalingComputeMaxInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MaxInstanceSize),
		BackupEnabled:                types.BoolPointerValue(apiResp.BackupEnabled),
		CloudBackup:                  types.BoolPointerValue(apiResp.ProviderBackupEnabled), //
		ClusterType:                  types.StringValue(apiResp.ClusterType),
		DiskSizeGb:                   types.Float64PointerValue(apiResp.DiskSizeGB),
		EncryptionAtRestProvider:     types.StringValue(apiResp.EncryptionAtRestProvider),
		MongoDBMajorVersion:          types.StringValue(apiResp.MongoDBMajorVersion), // TODO version formatting
		MongoDBVersion:               types.StringValue(apiResp.MongoDBVersion),
		MongoURI:                     types.StringValue(apiResp.MongoURI),
		MongoURIUpdated:              types.StringValue(apiResp.MongoURIUpdated),
		MongoURIWithOptions:          types.StringValue(apiResp.MongoURIWithOptions),
		PitEnabled:                   types.BoolPointerValue(apiResp.PitEnabled),
		Paused:                       types.BoolPointerValue(apiResp.Paused),
		SrvAddress:                   types.StringValue(apiResp.SrvAddress),
		StateName:                    types.StringValue(apiResp.StateName),
		TerminationProtectionEnabled: types.BoolPointerValue(apiResp.TerminationProtectionEnabled),
		ReplicationFactor:            types.Int64PointerValue(apiResp.ReplicationFactor),
		ConnectionStrings:            newTFConnectionStringsModel(apiResp.ConnectionStrings),
		BiConnectorConfig:            newTFBiConnectorConfigModel(apiResp.BiConnector),
		ReplicationSpecs:             newTFReplicationSpecsModel(apiResp.ReplicationSpecs),
		Labels:                       removeDefaultLabel(newTFLabelsModel(apiResp.Labels)),
		Tags:                         newTFTagsModel(apiResp.Tags),
		VersionReleaseSystem:         types.StringValue(apiResp.VersionReleaseSystem),
	}

	if isImport {
		clusterModel.ClusterID = types.StringValue(apiResp.ID)
		// clusterModel.ProjectID = types.StringValue(apiResp.GroupID)
		// clusterModel.Name = types.StringValue(apiResp.Name)
		//  clusterModel.ProviderName = types.StringValue(apiResp.ProviderSettings.ProviderName)
		clusterModel.CloudBackup = types.BoolPointerValue(apiResp.ProviderBackupEnabled)
	} else {
		clusterModel.ID = currState.ID
		// clusterModel.ClusterID = currState.ClusterID
		// clusterModel.ProjectID = currState.ProjectID

		if !currState.CloudBackup.IsNull() {
			clusterModel.CloudBackup = types.BoolPointerValue(apiResp.ProviderBackupEnabled)
		}
	}

	// Avoid Global Cluster issues. (NumShards is not present in Global Clusters)
	if numShards := apiResp.NumShards; numShards != nil {
		clusterModel.NumShards = types.Int64PointerValue(numShards)
	}

	if apiResp.ProviderSettings != nil {
		setTFProviderSettings(&clusterModel, apiResp.ProviderSettings)
	}

	if v := apiResp.ProviderSettings.ProviderName; v != "TENANT" {
		containers, _, err := conn.Containers.List(ctx, projectID,
			&matlas.ContainersListOptions{ProviderName: v})
		if err != nil {
			return nil, fmt.Errorf(errorClusterRead, clusterName, err)
		}

		clusterModel.ContainerID = types.StringValue(getContainerID(containers, apiResp))
		clusterModel.AutoScalingDiskGBEnabled = types.BoolPointerValue(apiResp.AutoScaling.DiskGBEnabled)
	}

	clusterModel.AdvancedConfiguration, err = newTFAdvancedConfigurationModelFromAtlas(ctx, conn, currState.ProjectID.ValueString(), apiResp.Name)
	if err != nil {
		return nil, err
	}

	clusterModel.SnapshotBackupPolicy, err = newTFSnapshotBackupPolicyModel(ctx, currState, conn, projectID, clusterName)
	if err != nil {
		return nil, err
	}

	return &clusterModel, nil
}

func newTFSnapshotBackupPolicyModel(ctx context.Context, currState *tfClusterRSModel, conn *matlas.Client, projectID, clusterName string) ([]tfSnapshotBackupPolicyModel, error) {
	backupPolicy, res, err := conn.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)
	if err != nil {
		if res.StatusCode == http.StatusNotFound ||
			strings.Contains(err.Error(), "BACKUP_CONFIG_NOT_FOUND") ||
			strings.Contains(err.Error(), "Not Found") ||
			strings.Contains(err.Error(), "404") {
			return []tfSnapshotBackupPolicyModel{}, nil
		}

		return []tfSnapshotBackupPolicyModel{}, fmt.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
	}

	return []tfSnapshotBackupPolicyModel{
		{
			ClusterID:             conversion.StringNullIfEmpty(backupPolicy.ClusterID),
			ClusterName:           conversion.StringNullIfEmpty(backupPolicy.ClusterName),
			NextSnapshot:          conversion.StringNullIfEmpty(backupPolicy.NextSnapshot),
			ReferenceHourOfDay:    types.Int64PointerValue(backupPolicy.ReferenceHourOfDay),
			ReferenceMinuteOfHour: types.Int64PointerValue(backupPolicy.ReferenceMinuteOfHour),
			RestoreWindowDays:     types.Int64PointerValue(backupPolicy.RestoreWindowDays),
			UpdateSnapshots:       types.BoolPointerValue(backupPolicy.UpdateSnapshots),
			Policies:              newTFSnapshotPolicyModel(backupPolicy.Policies),
		},
	}, nil
}

func newTFSnapshotPolicyModel(policies []matlas.Policy) []tfSnapshotPolicyModel {
	res := make([]tfSnapshotPolicyModel, len(policies))

	for i, pe := range policies {
		res[i] = tfSnapshotPolicyModel{
			ID:         conversion.StringNullIfEmpty(pe.ID),
			PolicyItem: newTFSnapshotPolicyItemModel(pe.PolicyItems),
		}
	}
	return res
}

func newTFSnapshotPolicyItemModel(policyItems []matlas.PolicyItem) []tfSnapshotPolicyItemModel {
	res := make([]tfSnapshotPolicyItemModel, len(policyItems))

	for i, pe := range policyItems {
		res[i] = tfSnapshotPolicyItemModel{
			ID:                conversion.StringNullIfEmpty(pe.ID),
			FrequencyInterval: types.Int64Value(cast.ToInt64(pe.FrequencyInterval)),
			FrequencyType:     conversion.StringNullIfEmpty(pe.FrequencyType),
			RetentionUnit:     conversion.StringNullIfEmpty(pe.RetentionUnit),
			RetentionValue:    types.Int64Value(cast.ToInt64(pe.RetentionValue)),
		}
	}
	return res
}

func setTFProviderSettings(clusterModel *tfClusterRSModel, settings *matlas.ProviderSettings) {
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

func newTFAdvancedConfigurationModelFromAtlas(ctx context.Context, conn *matlas.Client, projectID, clusterName string) ([]tfAdvancedConfigurationModel, error) {
	processArgs, _, err := conn.Clusters.GetProcessArgs(ctx, projectID, clusterName)

	return newTfAdvancedConfigurationModel(processArgs), err
}

func newTfAdvancedConfigurationModel(p *matlas.ProcessArgs) []tfAdvancedConfigurationModel {
	return []tfAdvancedConfigurationModel{
		{
			DefaultReadConcern:  conversion.StringNullIfEmpty(p.DefaultReadConcern),
			DefaultWriteConcern: conversion.StringNullIfEmpty(p.DefaultWriteConcern),

			FailIndexKeyTooLong:              types.BoolPointerValue(p.FailIndexKeyTooLong),
			JavascriptEnabled:                types.BoolPointerValue(p.JavascriptEnabled),
			MinimumEnabledTLSProtocol:        conversion.StringNullIfEmpty(p.MinimumEnabledTLSProtocol),
			NoTableScan:                      types.BoolPointerValue(p.NoTableScan),
			OplogSizeMB:                      types.Int64PointerValue(p.OplogSizeMB),
			OplogMinRetentionHours:           types.Int64Value(cast.ToInt64(p.OplogMinRetentionHours)),
			SampleSizeBiConnector:            types.Int64PointerValue(p.SampleSizeBIConnector),
			SampleRefreshIntervalBiConnector: types.Int64PointerValue(p.SampleRefreshIntervalBIConnector),
			TransactionLifetimeLimitSeconds:  types.Int64PointerValue(p.TransactionLifetimeLimitSeconds),
		},
	}
}

func removeDefaultLabel(labels []tfLabelModel) []tfLabelModel {
	var result []tfLabelModel

	for _, item := range labels {
		if item.Key.ValueString() == defaultLabel.Key && item.Value.ValueString() == defaultLabel.Value {
			continue
		}
		result = append(result, item)
	}

	return result
}

func newTFTagsModel(tags *[]*matlas.Tag) []tfTagModel {
	res := make([]tfTagModel, len(*tags))

	for i, v := range *tags {
		res[i] = tfTagModel{
			Key:   types.StringValue(v.Key),
			Value: types.StringValue(v.Value),
		}
	}

	return res
}

func newTFReplicationSpecsModel(replicationSpecs []matlas.ReplicationSpec) []tfReplicationSpecModel {
	res := make([]tfReplicationSpecModel, len(replicationSpecs))

	for i, rSpec := range replicationSpecs {
		res[i] = tfReplicationSpecModel{
			ID:            conversion.StringNullIfEmpty(rSpec.ID),
			NumShards:     types.Int64PointerValue(rSpec.NumShards),
			ZoneName:      conversion.StringNullIfEmpty(rSpec.ZoneName),
			RegionsConfig: newTFRegionsConfigModel(rSpec.RegionsConfig),
		}
	}
	return res
}

func newTFRegionsConfigModel(regionsConfig map[string]matlas.RegionsConfig) []tfRegionConfigModel {
	res := make([]tfRegionConfigModel, len(regionsConfig))

	for regionName, regionConfig := range regionsConfig {
		region := tfRegionConfigModel{
			RegionName:     conversion.StringNullIfEmpty(regionName),
			Priority:       types.Int64PointerValue(regionConfig.Priority),
			AnalyticsNodes: types.Int64PointerValue(regionConfig.AnalyticsNodes),
			ElectableNodes: types.Int64PointerValue(regionConfig.ElectableNodes),
			ReadOnlyNodes:  types.Int64PointerValue(regionConfig.ReadOnlyNodes),
		}
		res = append(res, region)
	}
	return res
}

func newTFBiConnectorConfigModel(biConnector *matlas.BiConnector) []tfBiConnectorConfigModel {
	if biConnector == nil {
		return []tfBiConnectorConfigModel{}
	}

	return []tfBiConnectorConfigModel{
		{
			Enabled:        types.BoolPointerValue(biConnector.Enabled),
			ReadPreference: conversion.StringNullIfEmpty(biConnector.ReadPreference),
		},
	}
}

func newTFConnectionStringsModel(connString *matlas.ConnectionStrings) []tfConnectionStringModel {
	if connString == nil {
		return []tfConnectionStringModel{}
	}

	return []tfConnectionStringModel{
		{
			Standard:        conversion.StringNullIfEmpty(connString.Standard),
			StandardSrv:     conversion.StringNullIfEmpty(connString.StandardSrv),
			Private:         conversion.StringNullIfEmpty(connString.Private),
			PrivateSrv:      conversion.StringNullIfEmpty(connString.PrivateSrv),
			PrivateEndpoint: newTFPrivateEndpointModel(connString.PrivateEndpoint),
		},
	}
}

func newTFPrivateEndpointModel(privateEndpoints []matlas.PrivateEndpoint) []tfPrivateEndpointModel {
	// if len(privateEndpoints) == 0 {
	// 	return []tfPrivateEndpointModel{}
	// }

	res := make([]tfPrivateEndpointModel, len(privateEndpoints))

	for i, pe := range privateEndpoints {
		res[i] = tfPrivateEndpointModel{
			ConnectionString:                  conversion.StringNullIfEmpty(pe.ConnectionString),
			SrvConnectionString:               conversion.StringNullIfEmpty(pe.SRVConnectionString),
			SrvShardOptimizedConnectionString: conversion.StringNullIfEmpty(pe.SRVShardOptimizedConnectionString),
			EndpointType:                      conversion.StringNullIfEmpty(pe.Type),
			Endpoints:                         newTFEndpointModel(pe.Endpoints),
		}
	}
	return res
}

func newTFEndpointModel(endpoints []matlas.Endpoint) []tfEndpointModel {
	res := make([]tfEndpointModel, len(endpoints))

	for i, e := range endpoints {
		res[i] = tfEndpointModel{
			Region:       conversion.StringNullIfEmpty(e.Region),
			ProviderName: conversion.StringNullIfEmpty(e.ProviderName),
			EndpointID:   conversion.StringNullIfEmpty(e.EndpointID),
		}
	}
	return res
}

func newAtlasProcessArgs(tfModel *tfAdvancedConfigurationModel) *matlas.ProcessArgs {
	res := &matlas.ProcessArgs{}

	if v := tfModel.DefaultReadConcern; !v.IsNull() {
		res.DefaultReadConcern = v.ValueString()
	}
	if v := tfModel.DefaultWriteConcern; !v.IsNull() {
		res.DefaultWriteConcern = v.ValueString()
	}

	if v := tfModel.FailIndexKeyTooLong; !v.IsNull() {
		res.FailIndexKeyTooLong = v.ValueBoolPointer()
	}

	if v := tfModel.JavascriptEnabled; !v.IsNull() {
		res.JavascriptEnabled = v.ValueBoolPointer()
	}

	if v := tfModel.MinimumEnabledTLSProtocol; !v.IsNull() {
		res.MinimumEnabledTLSProtocol = v.ValueString()
	}

	if v := tfModel.NoTableScan; !v.IsNull() {
		res.NoTableScan = v.ValueBoolPointer()
	}

	if v := tfModel.SampleSizeBiConnector; !v.IsNull() {
		res.SampleSizeBIConnector = v.ValueInt64Pointer()
	}

	if v := tfModel.SampleRefreshIntervalBiConnector; !v.IsNull() {
		res.SampleRefreshIntervalBIConnector = v.ValueInt64Pointer()
	}

	if v := tfModel.OplogSizeMB; !v.IsNull() {
		if sizeMB := v.ValueInt64(); sizeMB != 0 {
			res.OplogSizeMB = v.ValueInt64Pointer()
		} else {
			log.Printf(errorClusterSetting, `oplog_size_mb`, "", cast.ToString(sizeMB))
		}
	}

	if v := tfModel.OplogMinRetentionHours; !v.IsNull() {
		if minRetentionHours := v.ValueInt64(); minRetentionHours >= 0 {
			res.OplogMinRetentionHours = pointy.Float64(cast.ToFloat64(v.ValueInt64()))
		} else {
			log.Printf(errorClusterSetting, `oplog_min_retention_hours`, "", cast.ToString(minRetentionHours))
		}
	}

	if v := tfModel.TransactionLifetimeLimitSeconds; !v.IsNull() {
		if transactionLimitSeconds := v.ValueInt64(); transactionLimitSeconds > 0 {
			res.TransactionLifetimeLimitSeconds = v.ValueInt64Pointer()
		} else {
			log.Printf(errorClusterSetting, `transaction_lifetime_limit_seconds`, "", cast.ToString(transactionLimitSeconds))
		}
	}

	return res
}

func newAtlasTags(list []tfTagModel) []*matlas.Tag {
	res := make([]*matlas.Tag, len(list))
	for i, v := range list {
		res[i] = &matlas.Tag{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}
	return res
}

func newAtlasLabels(list []tfLabelModel) []matlas.Label {
	res := make([]matlas.Label, len(list))

	for i, v := range list {
		res[i] = matlas.Label{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}

	return res
}

func newAtlasBiConnectorConfig(plan *tfClusterRSModel) (*matlas.BiConnector, error) {
	var biConnector matlas.BiConnector

	if v := plan.BiConnectorConfig; v != nil {
		if len(v) > 0 {
			biConnMap := v[0]

			biConnector = matlas.BiConnector{
				Enabled:        biConnMap.Enabled.ValueBoolPointer(),
				ReadPreference: biConnMap.ReadPreference.ValueString(),
			}
		}
	}

	return &biConnector, nil
}

func newAtlasProviderSetting(tfClusterModel *tfClusterRSModel) (*matlas.ProviderSettings, error) {
	var (
		region, _          = valRegion(tfClusterModel.ProviderRegionName.ValueString())
		minInstanceSize    = getInstanceSizeToInt(tfClusterModel.ProviderAutoScalingComputeMinInstanceSize.ValueString())
		maxInstanceSize    = getInstanceSizeToInt(tfClusterModel.ProviderAutoScalingComputeMaxInstanceSize.ValueString())
		instanceSize       = getInstanceSizeToInt(tfClusterModel.ProviderInstanceSizeName.ValueString())
		compute            *matlas.Compute
		autoScalingEnabled = tfClusterModel.AutoScalingComputeEnabled.ValueBool()
		providerName       = tfClusterModel.ProviderName.ValueString()
	)

	if minInstanceSize != 0 && autoScalingEnabled {
		if instanceSize < minInstanceSize {
			return nil, fmt.Errorf("`provider_auto_scaling_compute_min_instance_size` must be lower than `provider_instance_size_name`")
		}

		compute = &matlas.Compute{
			MinInstanceSize: tfClusterModel.ProviderAutoScalingComputeMinInstanceSize.ValueString(),
		}
	}

	if maxInstanceSize != 0 && autoScalingEnabled {
		if instanceSize > maxInstanceSize {
			return nil, fmt.Errorf("`provider_auto_scaling_compute_max_instance_size` must be higher than `provider_instance_size_name`")
		}

		if compute == nil {
			compute = &matlas.Compute{}
		}
		compute.MaxInstanceSize = tfClusterModel.ProviderAutoScalingComputeMaxInstanceSize.ValueString()
	}

	providerSettings := &matlas.ProviderSettings{
		InstanceSizeName: tfClusterModel.ProviderInstanceSizeName.ValueString(),
		ProviderName:     providerName,
		RegionName:       region,
		VolumeType:       tfClusterModel.ProviderVolumeType.ValueString(),
	}

	// TODO include in update()
	// if d.HasChange("provider_disk_type_name") {
	// 	_, newdiskTypeName := d.GetChange("provider_disk_type_name")
	// 	diskTypeName := cast.ToString(newdiskTypeName)
	// 	if diskTypeName != "" { // ensure disk type is not included in request if attribute is removed, prevents errors in NVME intances
	// 		providerSettings.DiskTypeName = diskTypeName
	// 	}
	// }

	if providerName == "TENANT" {
		providerSettings.BackingProviderName = tfClusterModel.BackingProviderName.ValueString()
	}

	if autoScalingEnabled {
		providerSettings.AutoScaling = &matlas.AutoScaling{Compute: compute}
	}

	if tfClusterModel.ProviderName.ValueString() == "AWS" {
		// Check if the Provider Disk IOS sets in the Terraform configuration and if the instance size name is not NVME.
		// If it didn't, the MongoDB Atlas server would set it to the default for the amount of storage.
		if providerDiskIops := tfClusterModel.ProviderDiskIops; !providerDiskIops.IsNull() && !strings.Contains(providerSettings.InstanceSizeName, "NVME") {
			providerSettings.DiskIOPS = providerDiskIops.ValueInt64Pointer()
		}

		providerSettings.EncryptEBSVolume = pointy.Bool(true)
	}

	return providerSettings, nil
}

func newAtlasReplicationSpecs(tfClusterModel *tfClusterRSModel) ([]matlas.ReplicationSpec, error) {
	rSpecs := make([]matlas.ReplicationSpec, 0)

	vRSpecs := tfClusterModel.ReplicationSpecs
	// vPRName := tfClusterModel.ProviderRegionName

	if len(vRSpecs) == 0 {
		return rSpecs, nil
	}

	for _, spec := range vRSpecs {
		// spec := s.(map[string]interface{})

		replaceRegion := ""
		originalRegion := ""
		id := ""

		// TODO update() logic
		// if okPRName && d.Get("provider_name").(string) == "GCP" && cast.ToString(d.Get("cluster_type")) == "REPLICASET" {
		// 	if d.HasChange("provider_region_name") {
		// 		replaceRegion = vPRName.(string)
		// 		original, _ := d.GetChange("provider_region_name")
		// 		originalRegion = original.(string)
		// 	}
		// }

		// TODO update() logic
		// if d.HasChange("replication_specs") {
		// 	// Get original and new object
		// 	var oldSpecs map[string]interface{}
		// 	original, _ := d.GetChange("replication_specs")
		// 	for _, s := range original.(*schema.Set).List() {
		// 		oldSpecs = s.(map[string]interface{})
		// 		if spec["zone_name"].(string) == cast.ToString(oldSpecs["zone_name"]) {
		// 			id = oldSpecs["id"].(string)
		// 			break
		// 		}
		// 	}
		// 	if id == "" && oldSpecs != nil {
		// 		id = oldSpecs["id"].(string)
		// 	}
		// }

		// regionsConfig, err := expandRegionsConfig(spec["regions_config"].(*schema.Set).List(), originalRegion, replaceRegion)
		regionsConfig, err := newAtlasRegionsConfig(spec.RegionsConfig, originalRegion, replaceRegion)
		if err != nil {
			return rSpecs, err
		}

		rSpec := matlas.ReplicationSpec{
			ID:            id,
			NumShards:     spec.NumShards.ValueInt64Pointer(),
			ZoneName:      spec.ZoneName.ValueString(),
			RegionsConfig: regionsConfig,
		}
		rSpecs = append(rSpecs, rSpec)
	}
	return nil, nil // TODO complete
}

func newAtlasRegionsConfig(regions []tfRegionConfigModel, originalRegion, replaceRegion string) (map[string]matlas.RegionsConfig, error) {
	regionsConfig := make(map[string]matlas.RegionsConfig)

	for _, region := range regions {
		// region := r.(map[string]interface{})

		r, err := valRegion(region.RegionName.ValueString())
		if err != nil {
			return regionsConfig, err
		}

		if replaceRegion != "" && r == originalRegion {
			r, err = valRegion(replaceRegion)
		}
		if err != nil {
			return regionsConfig, err
		}

		regionsConfig[r] = matlas.RegionsConfig{
			AnalyticsNodes: region.AnalyticsNodes.ValueInt64Pointer(),
			ElectableNodes: region.ElectableNodes.ValueInt64Pointer(),
			Priority:       region.Priority.ValueInt64Pointer(),
			ReadOnlyNodes:  region.ReadOnlyNodes.ValueInt64Pointer(),
		}
	}

	return regionsConfig, nil
}

func validateClusterConfig(plan *tfClusterRSModel, response *resource.CreateResponse) {
	providerName := plan.ProviderName.ValueString()
	computeEnabled := plan.AutoScalingComputeEnabled.ValueBool()
	scaleDownEnabled := plan.AutoScalingComputeScaleDownEnabled.ValueBool()
	minInstanceSize := plan.ProviderAutoScalingComputeMinInstanceSize.ValueString()
	maxInstanceSize := plan.ProviderAutoScalingComputeMaxInstanceSize.ValueString()

	if scaleDownEnabled && !computeEnabled {
		response.Diagnostics.AddError(errorInvalidCreateValues, "`auto_scaling_compute_scale_down_enabled` must be set when `auto_scaling_compute_enabled` is set")
	}

	if computeEnabled && maxInstanceSize == "" {
		response.Diagnostics.AddError(errorInvalidCreateValues, "`provider_auto_scaling_compute_max_instance_size` must be set when `auto_scaling_compute_enabled` is set")
	}

	if scaleDownEnabled && minInstanceSize == "" {
		response.Diagnostics.AddError(errorInvalidCreateValues, "`provider_auto_scaling_compute_min_instance_size` must be set when `auto_scaling_compute_scale_down_enabled` is set")
	}

	if plan.ReplicationSpecs != nil && len(plan.ReplicationSpecs) > 0 {
		if plan.ClusterType.IsNull() {
			response.Diagnostics.AddError(errorInvalidCreateValues, "`cluster_type` should be set when `replication_specs` is set")
		}
		if plan.NumShards.IsNull() {
			response.Diagnostics.AddError(errorInvalidCreateValues, "`num_shards` should be set when `replication_specs` is set")
		}
	}

	if providerName != "AWS" {
		if plan.ProviderDiskIops.IsNull() {
			response.Diagnostics.AddError(errorInvalidCreateValues, "`provider_disk_iops` shouldn't be set when provider name is `GCP` or `AZURE`")
		}
		if plan.ProviderVolumeType.IsNull() {
			response.Diagnostics.AddError(errorInvalidCreateValues, "`provider_volume_type` shouldn't be set when provider name is `GCP` or `AZURE`")
		}
	}

	if providerName != "AZURE" {
		if plan.ProviderDiskTypeName.IsNull() {
			response.Diagnostics.AddError(errorInvalidCreateValues, "`provider_volume_type` shouldn't be set when provider name is `GCP` or `AZURE`")
		}
	}

	if providerName == "AZURE" {
		if plan.DiskSizeGb.IsNull() {
			response.Diagnostics.AddError(errorInvalidCreateValues, "`provider_disk_type_name` shouldn't be set when provider name is `GCP` or `AWS`")
		}
	}

	if providerName == "TENANT" {
		if instanceSizeName := plan.ProviderInstanceSizeName; !instanceSizeName.IsNull() {
			if instanceSizeName.ValueString() == "M2" {
				if diskSizeGB := plan.DiskSizeGb; !diskSizeGB.IsNull() {
					if cast.ToFloat64(diskSizeGB.ValueFloat64()) != 2 {
						response.Diagnostics.AddError(errorInvalidCreateValues, "`disk_size_gb` must be 2 for M2 shared tier")
					}
				}
			}
			if instanceSizeName.ValueString() == "M5" {
				if diskSizeGB := plan.DiskSizeGb; !diskSizeGB.IsNull() {
					if cast.ToFloat64(diskSizeGB.ValueFloat64()) != 5 {
						response.Diagnostics.AddError(errorInvalidCreateValues, "`disk_size_gb` must be 5 for M5 shared tier")
					}
				}
			}
		}
	}

	// We need to validate the oplog_size_mb attr of the advanced configuration option to show the error
	// before that the cluster is created
	if plan.AdvancedConfiguration != nil && len(plan.AdvancedConfiguration) > 0 {
		if oplogSizeMB := plan.AdvancedConfiguration[0].OplogSizeMB; !oplogSizeMB.IsNull() {
			if cast.ToInt64(oplogSizeMB.ValueInt64()) <= 0 {
				response.Diagnostics.AddError(errorInvalidCreateValues, "`advanced_configuration.oplog_size_mb` cannot be <= 0")
			}
		}
	}
}

type tfClusterRSModel struct {
	DiskSizeGb                                types.Float64                  `tfsdk:"disk_size_gb"`
	ProjectID                                 types.String                   `tfsdk:"project_id"`
	SrvAddress                                types.String                   `tfsdk:"srv_address"`
	BackingProviderName                       types.String                   `tfsdk:"backing_provider_name"`
	Timeouts                                  timeouts.Value                 `tfsdk:"timeouts"`
	ProviderAutoScalingComputeMinInstanceSize types.String                   `tfsdk:"provider_auto_scaling_compute_min_instance_size"`
	ClusterID                                 types.String                   `tfsdk:"cluster_id"`
	ClusterType                               types.String                   `tfsdk:"cluster_type"`
	ContainerID                               types.String                   `tfsdk:"container_id"`
	VersionReleaseSystem                      types.String                   `tfsdk:"version_release_system"`
	EncryptionAtRestProvider                  types.String                   `tfsdk:"encryption_at_rest_provider"`
	ID                                        types.String                   `tfsdk:"id"`
	MongoDBMajorVersion                       types.String                   `tfsdk:"mongo_db_major_version"`
	MongoDBMajorVersionFormatted              types.String                   `tfsdk:"mongo_db_major_version_formatted"`
	MongoDBVersion                            types.String                   `tfsdk:"mongo_db_version"`
	MongoURI                                  types.String                   `tfsdk:"mongo_uri"`
	MongoURIUpdated                           types.String                   `tfsdk:"mongo_uri_updated"`
	ProviderAutoScalingComputeMaxInstanceSize types.String                   `tfsdk:"provider_auto_scaling_compute_max_instance_size"`
	Name                                      types.String                   `tfsdk:"name"`
	StateName                                 types.String                   `tfsdk:"state_name"`
	ProviderVolumeType                        types.String                   `tfsdk:"provider_volume_type"`
	ProviderRegionName                        types.String                   `tfsdk:"provider_region_name"`
	ProviderName                              types.String                   `tfsdk:"provider_name"`
	ProviderInstanceSizeName                  types.String                   `tfsdk:"provider_instance_size_name"`
	ProviderDiskTypeName                      types.String                   `tfsdk:"provider_disk_type_name"`
	MongoURIWithOptions                       types.String                   `tfsdk:"mongo_uri_with_options"`
	Labels                                    []tfLabelModel                 `tfsdk:"labels"`
	AdvancedConfiguration                     []tfAdvancedConfigurationModel `tfsdk:"advanced_configuration"`
	Tags                                      []tfTagModel                   `tfsdk:"tags"`
	ReplicationSpecs                          []tfReplicationSpecModel       `tfsdk:"replication_specs"`
	BiConnectorConfig                         []tfBiConnectorConfigModel     `tfsdk:"bi_connector_config"`
	SnapshotBackupPolicy                      []tfSnapshotBackupPolicyModel  `tfsdk:"snapshot_backup_policy"`
	ConnectionStrings                         []tfConnectionStringModel      `tfsdk:"links"`
	ReplicationFactor                         types.Int64                    `tfsdk:"replication_factor"`
	NumShards                                 types.Int64                    `tfsdk:"num_shards"`
	ProviderDiskIops                          types.Int64                    `tfsdk:"provider_disk_iops"`
	PitEnabled                                types.Bool                     `tfsdk:"pit_enabled"`
	ProviderEncryptEbsVolume                  types.Bool                     `tfsdk:"provider_encrypt_ebs_volume"`
	AutoScalingComputeScaleDownEnabled        types.Bool                     `tfsdk:"auto_scaling_compute_scale_down_enabled"`
	BackupEnabled                             types.Bool                     `tfsdk:"backup_enabled"`
	RetainBackupsEnabled                      types.Bool                     `tfsdk:"retain_backups_enabled"`
	TerminationProtectionEnabled              types.Bool                     `tfsdk:"termination_protection_enabled"`
	Paused                                    types.Bool                     `tfsdk:"paused"`
	AutoScalingComputeEnabled                 types.Bool                     `tfsdk:"auto_scaling_compute_enabled"`
	ProviderEncryptEbsVolumeFlag              types.Bool                     `tfsdk:"provider_encrypt_ebs_volume_flag"`
	CloudBackup                               types.Bool                     `tfsdk:"cloud_backup"`
	AutoScalingDiskGBEnabled                  types.Bool                     `tfsdk:"auto_scaling_disk_gb_enabled"`
}

type tfSnapshotBackupPolicyModel struct {
	ClusterID             types.String            `tfsdk:"cluster_id"`
	ClusterName           types.String            `tfsdk:"cluster_name"`
	NextSnapshot          types.String            `tfsdk:"next_snapshot"`
	Policies              []tfSnapshotPolicyModel `tfsdk:"policies"`
	ReferenceHourOfDay    types.Int64             `tfsdk:"reference_hour_of_day"`
	ReferenceMinuteOfHour types.Int64             `tfsdk:"reference_minute_of_hour"`
	RestoreWindowDays     types.Int64             `tfsdk:"restore_window_days"`
	UpdateSnapshots       types.Bool              `tfsdk:"update_snapshots"`
}

type tfSnapshotPolicyModel struct {
	ID         types.String                `tfsdk:"id"`
	PolicyItem []tfSnapshotPolicyItemModel `tfsdk:"policy_item"`
}

type tfSnapshotPolicyItemModel struct {
	ID                types.String `tfsdk:"id"`
	FrequencyType     types.String `tfsdk:"frequency_type"`
	RetentionUnit     types.String `tfsdk:"retention_unit"`
	FrequencyInterval types.Int64  `tfsdk:"frequency_interval"`
	RetentionValue    types.Int64  `tfsdk:"retention_value"`
}

type tfTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type tfAdvancedConfigurationModel struct {
	DefaultReadConcern               types.String `tfsdk:"default_read_concern"`
	DefaultWriteConcern              types.String `tfsdk:"default_write_concern"`
	MinimumEnabledTLSProtocol        types.String `tfsdk:"minimum_enabled_tls_protocol"`
	OplogSizeMB                      types.Int64  `tfsdk:"oplog_size_mb"`
	OplogMinRetentionHours           types.Int64  `tfsdk:"oplog_min_retention_hours"`
	SampleSizeBiConnector            types.Int64  `tfsdk:"sample_size_bi_connector"`
	SampleRefreshIntervalBiConnector types.Int64  `tfsdk:"sample_refresh_interval_bi_connector"`
	TransactionLifetimeLimitSeconds  types.Int64  `tfsdk:"transaction_lifetime_limit_seconds"`
	FailIndexKeyTooLong              types.Bool   `tfsdk:"fail_index_key_too_long"`
	JavascriptEnabled                types.Bool   `tfsdk:"javascript_enabled"`
	NoTableScan                      types.Bool   `tfsdk:"no_table_scan"`
}

type tfReplicationSpecModel struct {
	ID            types.String          `tfsdk:"id"`
	ZoneName      types.String          `tfsdk:"zone_name"`
	RegionsConfig []tfRegionConfigModel `tfsdk:"regions_config"`
	NumShards     types.Int64           `tfsdk:"num_shards"`
}

type tfRegionConfigModel struct {
	RegionName     types.String `tfsdk:"region_name"`
	ElectableNodes types.Int64  `tfsdk:"electable_nodes"`
	Priority       types.Int64  `tfsdk:"priority"`
	ReadOnlyNodes  types.Int64  `tfsdk:"read_only_nodes"`
	AnalyticsNodes types.Int64  `tfsdk:"analytics_nodes"`
}

type tfBiConnectorConfigModel struct {
	ReadPreference types.String `tfsdk:"read_preference"`
	Enabled        types.Bool   `tfsdk:"enabled"`
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
	EndpointID   types.String `tfsdk:"endpoint_id"`
	ProviderName types.String `tfsdk:"provider_name"`
	Region       types.String `tfsdk:"region"`
}
