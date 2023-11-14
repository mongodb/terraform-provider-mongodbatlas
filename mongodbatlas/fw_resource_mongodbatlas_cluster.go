package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
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
	defaultTimeout           = (3 * time.Hour)
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_scaling_disk_gb_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				Description: "Clusters running MongoDB FCV 4.2 or later and any new Atlas clusters of any type do not support this parameter",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"retain_backups_enabled": schema.BoolAttribute{
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mongo_db_major_version_formatted": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"num_shards": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				// PlanModifiers: []planmodifier.Int64{
				// 	int64planmodifier.UseStateForUnknown(),
				// },
			},
			"cloud_backup": schema.BoolAttribute{
				Optional: true,
				Default:  booldefault.StaticBool(false),
				Computed: true,
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
					planmodifiers.UseNullForUnknownString(),
				},
			},
			"provider_disk_iops": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					planmodifiers.UseNullForUnknownInt64(),
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
					planmodifiers.UseNullForUnknownBool(),
				},
			},
			"provider_encrypt_ebs_volume_flag": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					planmodifiers.UseNullForUnknownBool(),
				},
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
					// int64planmodifier.UseStateForUnknown(),
					planmodifiers.UseNullForUnknownInt64(),
				},
			},
			"mongo_db_version": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mongo_uri": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mongo_uri_updated": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mongo_uri_with_options": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"paused": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"srv_address": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
				PlanModifiers: []planmodifier.String{
					// stringplanmodifier.UseStateForUnknown(),
					planmodifiers.UseNullForUnknownString(),
				},
			},
			"version_release_system": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("LTS", "CONTINUOUS"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_strings":     clusterRSConnectionStringSchemaAttribute(),
			"snapshot_backup_policy": clusterRSSnapshotBackupPolicySchemaAttribute(),
			// computed-only attributes:
			// "advanced_configuration_output": clusterRSAdvancedConfigurationSchemaAttribute(),
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
						},
						"value": schema.StringAttribute{
							Optional: true,
						},
					},
				},
				DeprecationMessage: fmt.Sprintf(DeprecationByDateWithReplacement, "September 2024", "tags"),
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
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
	}
}

func clusterRSAdvancedConfigurationSchemaAttribute() schema.ListNestedAttribute {
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
				"oplog_min_retention_hours": schema.Int64Attribute{
					Computed: true,
				},
				"oplog_size_mb": schema.Int64Attribute{
					Computed: true,
				},
				"sample_refresh_interval_bi_connector": schema.Int64Attribute{
					Computed: true,
				},
				"sample_size_bi_connector": schema.Int64Attribute{
					Computed: true,
				},
				"transaction_lifetime_limit_seconds": schema.Int64Attribute{
					Computed: true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
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
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
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
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
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
		// TODO test for inconsistent state error
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

	validateClusterConfig(ctx, &plan, response)
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

	replicationSpecs, err := updateAtlasReplicationSpecs(&plan)
	if err != nil {
		response.Diagnostics.AddError("Unable to CREATE cluster. Error in translating replication specs", fmt.Sprintf(errorClusterCreate, err))
		return
	}

	clusterRequest := &matlas.Cluster{
		Name:                     plan.Name.ValueString(),
		EncryptionAtRestProvider: plan.EncryptionAtRestProvider.ValueString(),
		ClusterType:              plan.ClusterType.ValueString(),
		BackupEnabled:            plan.BackupEnabled.ValueBoolPointer(),
		PitEnabled:               plan.PitEnabled.ValueBoolPointer(),
		AutoScaling:              autoScaling,
		ProviderSettings:         providerSettings,
		ReplicationSpecs:         replicationSpecs,
	}

	if cloudBackup := plan.CloudBackup; !cloudBackup.IsUnknown() {
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

	if v := plan.DiskSizeGb; !v.IsUnknown() {
		clusterRequest.DiskSizeGB = v.ValueFloat64Pointer()
	}

	tenantDisksize := pointy.Float64(0)
	if cast.ToFloat64(tenantDisksize) != 0 {
		clusterRequest.DiskSizeGB = tenantDisksize
	}
	if v := plan.MongoDBMajorVersion; !v.IsUnknown() {
		clusterRequest.MongoDBMajorVersion = formatMongoDBMajorVersion(v.ValueString())
	}
	if v := plan.ReplicationFactor; !v.IsUnknown() {
		clusterRequest.ReplicationFactor = v.ValueInt64Pointer()
	}
	if v := plan.NumShards; !v.IsUnknown() {
		clusterRequest.NumShards = v.ValueInt64Pointer()
	}
	if v := plan.TerminationProtectionEnabled; !v.IsUnknown() {
		clusterRequest.TerminationProtectionEnabled = v.ValueBoolPointer()
	}

	if v := plan.VersionReleaseSystem; !v.IsUnknown() {
		clusterRequest.VersionReleaseSystem = v.ValueString()
	}

	cluster, _, err := conn.Clusters.Create(ctx, projectID, clusterRequest)

	if err != nil {
		response.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", fmt.Sprintf(errorClusterCreate, err))
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
		response.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", fmt.Sprintf(errorClusterCreate, err))
		return
	}

	var acmodel []tfAdvancedConfigurationModel
	plan.AdvancedConfiguration.ElementsAs(ctx, &acmodel, true)
	if len(acmodel) > 0 {
		advancedConfReq := newAtlasProcessArgs(&acmodel[0])

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

	// get latest state from Atlas
	cluster, resp, err := conn.Clusters.Get(ctx, projectID, cluster.Name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("error in getting cluster during CREATE from Atlas", fmt.Sprintf(errorClusterCreate, err.Error()))
		return
	}

	newClusterState, err := newTFClusterModel(ctx, conn, false, cluster, &plan)
	if err != nil {
		response.Diagnostics.AddError("error in getting cluster during CREATE when translating to model", fmt.Sprintf(errorClusterCreate, err.Error()))
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, &newClusterState)...)
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
	conn := r.client.Atlas
	var state, plan tfClusterRSModel

	response.Diagnostics.Append(req.State.Get(ctx, &state)...)
	response.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(state.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	clusterReq := new(matlas.Cluster)
	clusterChangeDetect := new(matlas.Cluster)
	clusterChangeDetect.AutoScaling = &matlas.AutoScaling{Compute: &matlas.Compute{}}

	if !plan.Name.Equal(state.Name) {
		clusterReq.Name = plan.Name.ValueString()
	}

	err := updateBiConnectorConfig(clusterReq, &plan, &state)
	if err != nil {
		response.Diagnostics.AddError("Unable to UPDATE cluster. Error in updating bi_connector_config", fmt.Sprintf(errorClusterUpdate, clusterReq.Name, err))
	}

	err = updateProviderSettings(clusterReq, &plan, &state)
	if err != nil {
		response.Diagnostics.AddError("Unable to UPDATE cluster. Error in updating provider settings for cluster", fmt.Sprintf(errorClusterUpdate, clusterReq.Name, err))
	}

	err = updateReplicationSpecs(clusterReq, &plan, &state, plan.Name.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Unable to UPDATE cluster. Error in updating replication_specs for cluster", fmt.Sprintf(errorClusterUpdate, clusterReq.Name, err))
	}

	clusterReq.AutoScaling = &matlas.AutoScaling{Compute: &matlas.Compute{}}
	updateAutoScaling(clusterReq, &plan, &state)
	if err != nil {
		response.Diagnostics.AddError("Unable to UPDATE cluster. Error in updating auto_scaling_* properties for cluster", fmt.Sprintf(errorClusterUpdate, clusterReq.Name, err))
	}

	err = updateOtherClusterProps(clusterReq, &plan, &state)
	if err != nil {
		response.Diagnostics.AddError("Unable to UPDATE cluster. Error in updating properties for cluster", fmt.Sprintf(errorClusterUpdate, clusterReq.Name, err))
	}

	if response.Diagnostics.HasError() {
		return
	}

	timeout, diags := plan.Timeouts.Update(ctx, defaultTimeout)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	err = updateAdvancedConfiguration(ctx, conn, clusterReq, &plan, &state)
	if err != nil {
		response.Diagnostics.AddError("Unable to UPDATE cluster. Error in updating advanced_configuration for cluster", fmt.Sprintf(errorClusterUpdate, clusterReq.Name, err))
		return
	}

	if isUpgradeRequired2(plan.ProviderInstanceSizeName.ValueString(), state.ProviderInstanceSizeName.ValueString()) {
		updatedCluster, _, err := upgradeCluster(ctx, conn, clusterReq, projectID, clusterName, timeout)

		if err != nil {
			response.Diagnostics.AddError("Unable to UPDATE cluster. Error in upgrading the cluster", fmt.Sprintf(errorClusterUpdate, clusterReq.Name, err))
			return
		}

		plan.ID = types.StringValue(encodeStateID(map[string]string{
			"cluster_id":    updatedCluster.ID,
			"project_id":    projectID,
			"cluster_name":  updatedCluster.Name,
			"provider_name": updatedCluster.ProviderSettings.ProviderName,
		}))
	} else if !reflect.DeepEqual(clusterReq, clusterChangeDetect) {
		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, _, err := updateCluster(ctx, conn, clusterReq, projectID, clusterName, timeout)

			if didErrOnPausedCluster(err) {
				clusterRequest := &matlas.Cluster{
					Paused: pointy.Bool(false),
				}

				_, _, err = updateCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
			}

			if err != nil {
				return retry.NonRetryableError(fmt.Errorf(errorClusterUpdate, clusterName, err))
			}

			return nil
		})

		if err != nil {
			response.Diagnostics.AddError("Unable to UPDATE cluster. Error in updating cluster", fmt.Sprintf(errorClusterUpdate, clusterName, err))
			return
		}
	}

	if plan.Paused.ValueBool() && !isSharedTier(plan.ProviderInstanceSizeName.ValueString()) {
		clusterRequest := &matlas.Cluster{
			Paused: pointy.Bool(true),
		}

		_, _, err := updateCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
		if err != nil {
			response.Diagnostics.AddError("Unable to UPDATE(PAUSE) cluster. Error in PAUSING cluster", fmt.Sprintf(errorClusterUpdate, clusterName, err))
		}
	}

	// get latest state from Atlas
	cluster, resp, err := conn.Clusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			response.State.RemoveResource(ctx)
			return
		}
		response.Diagnostics.AddError("error in getting cluster during UPDATE from Atlas", fmt.Sprintf(errorClusterUpdate, clusterName, err.Error()))
		return
	}

	newClusterState, err := newTFClusterModel(ctx, conn, false, cluster, &plan)
	if err != nil {
		response.Diagnostics.AddError("error in getting cluster during UPDATE when translating to model", fmt.Sprintf(errorClusterUpdate, clusterName, err.Error()))
		return
	}

	// save updated data into terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &newClusterState)...)
}

func didErrOnPausedCluster(err error) bool {
	if err == nil {
		return false
	}

	var target *matlas.ErrorResponse

	return errors.As(err, &target) && target.ErrorCode == "CANNOT_UPDATE_PAUSED_CLUSTER"
}

func upgradeCluster(ctx context.Context, conn *matlas.Client, request *matlas.Cluster, projectID, name string, timeout time.Duration) (*matlas.Cluster, *matlas.Response, error) {
	request.Name = name

	cluster, resp, err := conn.Clusters.Upgrade(ctx, projectID, request)
	if err != nil {
		return nil, nil, err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterRefreshFunc(ctx, name, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	return cluster, resp, nil
}

func isUpgradeRequired2(plannedInstanceSizeName, stateInstanceSizeName string) bool {
	return stateInstanceSizeName != plannedInstanceSizeName && isSharedTier(stateInstanceSizeName)
}

func isSharedTier(instanceSize string) bool {
	return instanceSize == "M0" || instanceSize == "M2" || instanceSize == "M5"
}

func updateAdvancedConfiguration(ctx context.Context, conn *matlas.Client, clusterReq *matlas.Cluster, plan, state *tfClusterRSModel) error {
	if !reflect.DeepEqual(plan.AdvancedConfiguration, state.AdvancedConfiguration) {
		var acmodel []tfAdvancedConfigurationModel
		plan.AdvancedConfiguration.ElementsAs(ctx, &acmodel, true)

		if len(acmodel) > 0 {
			advancedConfReq := newAtlasProcessArgs(&acmodel[0])
			if !reflect.DeepEqual(advancedConfReq, matlas.ProcessArgs{}) {
				clusterName := plan.Name.ValueString()
				argResp, _, err := conn.Clusters.UpdateProcessArgs(ctx, plan.ProjectID.ValueString(), clusterName, advancedConfReq)
				if err != nil {
					return fmt.Errorf(errorAdvancedConfUpdate, clusterName+argResp.DefaultReadConcern, err)
				}
			}
		}
	}

	return nil
}

func updateOtherClusterProps(clusterReq *matlas.Cluster, plan, state *tfClusterRSModel) error {
	if !plan.EncryptionAtRestProvider.Equal(state.EncryptionAtRestProvider) {
		clusterReq.EncryptionAtRestProvider = plan.EncryptionAtRestProvider.ValueString()
	}

	if !plan.MongoDBMajorVersion.Equal(state.MongoDBMajorVersion) {
		clusterReq.MongoDBMajorVersion = formatMongoDBMajorVersion(plan.MongoDBMajorVersion.ValueString())
	}

	if !plan.ClusterType.Equal(state.ClusterType) {
		clusterReq.ClusterType = plan.ClusterType.ValueString()
	}

	if !plan.BackupEnabled.Equal(state.BackupEnabled) {
		clusterReq.BackupEnabled = plan.BackupEnabled.ValueBoolPointer()
	}

	if !plan.DiskSizeGb.Equal(state.DiskSizeGb) {
		clusterReq.DiskSizeGB = plan.DiskSizeGb.ValueFloat64Pointer()
	}

	if !plan.CloudBackup.Equal(state.CloudBackup) {
		clusterReq.ProviderBackupEnabled = plan.CloudBackup.ValueBoolPointer()
	}

	if !plan.PitEnabled.Equal(state.PitEnabled) {
		clusterReq.PitEnabled = plan.PitEnabled.ValueBoolPointer()
	}

	if !plan.ReplicationFactor.Equal(state.ReplicationFactor) {
		clusterReq.ReplicationFactor = plan.ReplicationFactor.ValueInt64Pointer()
	}

	if !plan.NumShards.Equal(state.NumShards) {
		clusterReq.NumShards = plan.NumShards.ValueInt64Pointer()
	}

	if !plan.VersionReleaseSystem.Equal(state.VersionReleaseSystem) {
		clusterReq.VersionReleaseSystem = plan.VersionReleaseSystem.ValueString()
	}

	if !plan.TerminationProtectionEnabled.Equal(state.TerminationProtectionEnabled) {
		clusterReq.TerminationProtectionEnabled = plan.TerminationProtectionEnabled.ValueBoolPointer()
	}

	if hasLabelsChanged(plan.Labels, state.Labels) {
		if containsLabelOrKey(newAtlasLabels(plan.Labels), defaultLabel) {
			return fmt.Errorf("you should not set `Infrastructure Tool` label, it is used for internal purposes")
		}
		clusterReq.Labels = append(newAtlasLabels(plan.Labels), defaultLabel)
	}

	if hasTagsChanged(plan.Tags, state.Tags) {
		tagsSlice := newAtlasTags(plan.Tags)
		clusterReq.Tags = &tagsSlice
	}

	// when Provider instance type changes this argument must be passed explicitly in patch request
	if !plan.ProviderInstanceSizeName.Equal(state.ProviderInstanceSizeName) {
		if !plan.CloudBackup.IsNull() {
			clusterReq.ProviderBackupEnabled = plan.CloudBackup.ValueBoolPointer()
		}
	}

	if !plan.Paused.Equal(state.Paused) && !plan.Paused.ValueBool() {
		clusterReq.Paused = plan.Paused.ValueBoolPointer()
	}
	return nil
}

func hasLabelsChanged(planLabels, stateLables []tfLabelModel) bool {
	sort.Slice(planLabels, func(i, j int) bool {
		return planLabels[i].Key.ValueString() < planLabels[j].Key.ValueString()
	})
	sort.Slice(stateLables, func(i, j int) bool {
		return stateLables[i].Key.ValueString() < stateLables[j].Key.ValueString()
	})
	return !reflect.DeepEqual(planLabels, stateLables)
}

func hasTagsChanged(planTags, stateTags []*tfTagModel) bool {
	sort.Slice(planTags, func(i, j int) bool {
		return planTags[i].Key.ValueString() < planTags[j].Key.ValueString()
	})
	sort.Slice(stateTags, func(i, j int) bool {
		return stateTags[i].Key.ValueString() < stateTags[j].Key.ValueString()
	})
	return !reflect.DeepEqual(planTags, stateTags)
}

// TODO implement replicationSpecs create and update
func updateReplicationSpecs(clusterReq *matlas.Cluster, plan, state *tfClusterRSModel, clusterName string) error {
	// if areTFReplicationSpecSlicesEqual(plan.ReplicationSpecs, state.ReplicationSpecs) {
	// 	return nil
	// }

	rSpecs := make([]matlas.ReplicationSpec, 0)

	vRSpecs := plan.ReplicationSpecs
	// vPRName := plan.ProviderRegionName

	if len(vRSpecs) > 0 {
		for _, newSpec := range vRSpecs { // for each plan.replication_specs
			replaceRegion := ""
			originalRegion := ""
			id := ""

			// if plan.provider_name is GCP and plan.cluster_type is REPLICASET
			// and plan.provider_region_name != state.provider_region_name (has changed) then
			// get newProviderRegion(plan) and oldProviderRegion(state)
			if plan.ProviderRegionName.ValueString() == "GCP" && plan.ClusterType.ValueString() == "REPLICASET" {
				if !state.ProviderRegionName.Equal(plan.ProviderRegionName) {
					replaceRegion = plan.ProviderRegionName.ValueString()
					originalRegion = state.ProviderRegionName.ValueString()
				}
			}

			// Get original and new object
			var oldSpecs *tfReplicationSpecModel
			original := state.ReplicationSpecs
			for _, oldSpecsPtr := range original { // iterate over state.replication_specs
				oldSpecs = oldSpecsPtr
				if newSpec.ZoneName.ValueString() == oldSpecs.ZoneName.ValueString() { // find plan.replication_specs with matching zone_name
					id = oldSpecs.ID.ValueString() // and get it's id
					break
				}
			}
			if id == "" && oldSpecs != nil { // if match not found with same zone name
				id = oldSpecs.ID.ValueString() // then id = plan.replication_specs[i]
			}

			regionsConfig, err := updateRegionConfigs(newSpec.RegionsConfig, originalRegion, replaceRegion)
			if err != nil {
				return err
			}

			rSpec := matlas.ReplicationSpec{
				ID:            id,
				NumShards:     newSpec.NumShards.ValueInt64Pointer(),
				ZoneName:      newSpec.ZoneName.ValueString(),
				RegionsConfig: regionsConfig,
			}
			rSpecs = append(rSpecs, rSpec)
		}
	}
	clusterReq.ReplicationSpecs = rSpecs

	return nil
}

func updateRegionConfigs(regions []tfRegionConfigModel, originalRegion, replaceRegion string) (map[string]matlas.RegionsConfig, error) {
	regionsConfig := make(map[string]matlas.RegionsConfig)

	for _, region := range regions { // for each plan.region_configs
		r, err := valRegion(region.RegionName.ValueString()) // r = plan.region_configs[i].region_name
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

func updateBiConnectorConfig(clusterReq *matlas.Cluster, plan, state *tfClusterRSModel) error {
	if !reflect.DeepEqual(plan.BiConnectorConfig, state.BiConnectorConfig) {
		biConnector, err := newAtlasBiConnectorConfig(plan)
		if err != nil {
			return fmt.Errorf(errorClusterCreate, err)
		}
		clusterReq.BiConnector = biConnector
	}
	return nil
}

func updateAutoScaling(clusterReq *matlas.Cluster, plan, state *tfClusterRSModel) {
	if !plan.AutoScalingDiskGBEnabled.Equal(state.AutoScalingDiskGBEnabled) {
		clusterReq.AutoScaling.DiskGBEnabled = plan.AutoScalingDiskGBEnabled.ValueBoolPointer()
	}

	if !plan.AutoScalingComputeEnabled.Equal(state.AutoScalingComputeEnabled) {
		clusterReq.AutoScaling.Compute.Enabled = plan.AutoScalingComputeEnabled.ValueBoolPointer()
	}

	if !plan.AutoScalingComputeScaleDownEnabled.Equal(state.AutoScalingComputeScaleDownEnabled) {
		clusterReq.AutoScaling.Compute.ScaleDownEnabled = plan.AutoScalingComputeScaleDownEnabled.ValueBoolPointer()
	}
}

func updateProviderSettings(cluster *matlas.Cluster, plan, state *tfClusterRSModel) error {
	properties := []string{"ProviderDiskIops", "BackingProviderName", "ProviderDiskTypeName",
		"ProviderInstanceSizeName", "ProviderName", "ProviderRegionName", "ProviderVolumeType",
		"ProviderAutoScalingComputeMaxInstanceSize", "ProviderAutoScalingComputeMinInstanceSize"}

	// If at least one of the provider settings argument has changed, update all provider settings
	vplan := reflect.ValueOf(plan).Elem()
	vstate := reflect.ValueOf(state).Elem()

	for _, prop := range properties {
		fieldPlan := vplan.FieldByName(prop).Interface()
		fieldState := vstate.FieldByName(prop).Interface()

		if !reflect.ValueOf(fieldPlan).MethodByName("Equal").IsValid() {
			return fmt.Errorf("the property %s in ProviderSettings does not have an Equal method", prop)
		}

		result := reflect.ValueOf(fieldPlan).MethodByName("Equal").Call([]reflect.Value{reflect.ValueOf(fieldState)})

		if len(result) == 0 || !result[0].Bool() { // not equal
			var err error
			cluster.ProviderSettings, err = newAtlasProviderSetting(plan)
			if err != nil {
				return fmt.Errorf(errorClusterUpdate, plan.Name.ValueString(), err)
			}
		}
	}

	return nil
}

func (r *ClusterRS) Delete(ctx context.Context, req resource.DeleteRequest, response *resource.DeleteResponse) {
	conn := r.client.Atlas

	var clusterState tfClusterRSModel
	response.Diagnostics.Append(req.State.Get(ctx, &clusterState)...)
	if response.Diagnostics.HasError() {
		return
	}

	ids := decodeStateID(clusterState.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	var options *matlas.DeleteAdvanceClusterOptions
	if v := clusterState.RetainBackupsEnabled; !v.IsNull() {
		options = &matlas.DeleteAdvanceClusterOptions{
			RetainBackups: v.ValueBoolPointer(),
		}
	}

	_, err := conn.Clusters.Delete(ctx, projectID, clusterName, options)
	if err != nil {
		response.Diagnostics.AddError("error during cluster DELETE in Atlas", fmt.Sprintf(errorClusterDelete, clusterName, err.Error()))
		return
	}

	log.Println("[INFO] Waiting for MongoDB Cluster to be destroyed")

	timeout, diags := clusterState.Timeouts.Delete(ctx, defaultTimeout)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceClusterRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		response.Diagnostics.AddError("error during cluster DELETE in Atlas", fmt.Sprintf(errorClusterDelete, clusterName, err.Error()))
		return
	}
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
		ClusterID:                          types.StringValue(apiResp.ID),
		ProjectID:                          types.StringValue(projectID),
		Name:                               types.StringValue(clusterName),
		ProviderName:                       types.StringValue(apiResp.ProviderSettings.ProviderName),
		AutoScalingComputeEnabled:          types.BoolPointerValue(apiResp.AutoScaling.Compute.Enabled),
		AutoScalingComputeScaleDownEnabled: types.BoolPointerValue(apiResp.AutoScaling.Compute.ScaleDownEnabled),
		ProviderAutoScalingComputeMinInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MinInstanceSize),
		ProviderAutoScalingComputeMaxInstanceSize: types.StringValue(apiResp.ProviderSettings.AutoScaling.Compute.MaxInstanceSize),
		BackupEnabled:                types.BoolPointerValue(apiResp.BackupEnabled),
		CloudBackup:                  types.BoolPointerValue(apiResp.ProviderBackupEnabled),
		ClusterType:                  types.StringValue(apiResp.ClusterType),
		DiskSizeGb:                   types.Float64PointerValue(apiResp.DiskSizeGB),
		EncryptionAtRestProvider:     types.StringValue(apiResp.EncryptionAtRestProvider),
		MongoDBMajorVersion:          types.StringValue(apiResp.MongoDBMajorVersion),
		MongoDBMajorVersionFormatted: types.StringValue(apiResp.MongoDBMajorVersion),
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
		ConnectionStrings:            newTFConnectionStringsModelList(ctx, apiResp.ConnectionStrings),
		BiConnectorConfig:            newTFBiConnectorConfigModel(apiResp.BiConnector),
		ReplicationSpecs:             newTFReplicationSpecsModel(apiResp.ReplicationSpecs),
		Labels:                       removeDefaultLabel(newTFLabelsModel(apiResp.Labels)),
		Tags:                         newTFTagsModel(apiResp.Tags),
		VersionReleaseSystem:         types.StringValue(apiResp.VersionReleaseSystem),
		Timeouts:                     currState.Timeouts,
	}

	if isImport {
		clusterModel.CloudBackup = types.BoolPointerValue(apiResp.ProviderBackupEnabled)
	} else {
		clusterModel.ID = currState.ID

		if currState.MongoDBMajorVersion.ValueString() == "" {
			clusterModel.MongoDBMajorVersion = types.StringValue("")
		}

		if !currState.CloudBackup.IsNull() {
			clusterModel.CloudBackup = types.BoolPointerValue(apiResp.ProviderBackupEnabled)
		}

		if !currState.BackingProviderName.IsNull() {
			clusterModel.BackingProviderName = currState.BackingProviderName
		}
		// clusterModel.ProviderEncryptEbsVolume = types.BoolUnknown()
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

	// clusterModel.AdvancedConfigurationOutput, err = newTFAdvancedConfigurationModelFromAtlas(ctx, conn, currState.ProjectID.ValueString(), apiResp.Name)
	// clusterModel.AdvancedConfiguration = currState.AdvancedConfiguration
	clusterModel.AdvancedConfiguration, err = newTFAdvancedConfigurationModelFromAtlas(ctx, conn, currState.ProjectID.ValueString(), apiResp.Name)
	if err != nil {
		return nil, err
	}

	clusterModel.SnapshotBackupPolicy, err = newTFSnapshotBackupPolicyRSModel(ctx, conn, projectID, clusterName)
	if err != nil {
		return nil, err
	}

	clusterModel.ID = types.StringValue(encodeStateID(map[string]string{
		"cluster_id":    currState.ClusterID.ValueString(),
		"project_id":    projectID,
		"cluster_name":  currState.Name.ValueString(),
		"provider_name": currState.ProviderName.ValueString(),
	}))

	return &clusterModel, nil
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
		clusterModel.ProviderEncryptEbsVolume = types.BoolPointerValue(settings.EncryptEBSVolume)
	}
	clusterModel.ProviderDiskTypeName = types.StringValue(settings.DiskTypeName)
	clusterModel.ProviderInstanceSizeName = types.StringValue(settings.InstanceSizeName)
	clusterModel.ProviderName = types.StringValue(settings.ProviderName)
	clusterModel.ProviderRegionName = types.StringValue(settings.RegionName)
	clusterModel.ProviderVolumeType = types.StringValue(settings.VolumeType)
}

func newTFAdvancedConfigurationModelFromAtlas(ctx context.Context, conn *matlas.Client, projectID, clusterName string) (types.List, error) {
	processArgs, _, err := conn.Clusters.GetProcessArgs(ctx, projectID, clusterName)
	if err != nil {
		return types.ListNull(tfAdvancedConfigurationType), err
	}

	advConfigModel := newTfAdvancedConfigurationModel(ctx, processArgs)
	l, _ := types.ListValueFrom(ctx, tfAdvancedConfigurationType, advConfigModel)

	return l, err
}

func newTfAdvancedConfigurationModel(ctx context.Context, p *matlas.ProcessArgs) []*tfAdvancedConfigurationModel {
	res := []*tfAdvancedConfigurationModel{
		{
			DefaultReadConcern:               conversion.StringNullIfEmpty(p.DefaultReadConcern),
			DefaultWriteConcern:              conversion.StringNullIfEmpty(p.DefaultWriteConcern),
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
	if p.OplogMinRetentionHours != nil {
		res[0].OplogMinRetentionHours = types.Int64PointerValue(p.OplogSizeMB)
	}
	return res
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
			ID:            conversion.StringNullIfEmpty(rSpec.ID),
			NumShards:     types.Int64PointerValue(rSpec.NumShards),
			ZoneName:      conversion.StringNullIfEmpty(rSpec.ZoneName),
			RegionsConfig: newTFRegionsConfigModel(rSpec.RegionsConfig),
		}
	}
	return res
}

func newTFRegionsConfigModel(regionsConfig map[string]matlas.RegionsConfig) []tfRegionConfigModel {
	res := []tfRegionConfigModel{}

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

func newTFBiConnectorConfigModel(biConnector *matlas.BiConnector) []*tfBiConnectorConfigModel {
	if biConnector == nil {
		return []*tfBiConnectorConfigModel{}
	}

	return []*tfBiConnectorConfigModel{
		{
			Enabled:        types.BoolPointerValue(biConnector.Enabled),
			ReadPreference: conversion.StringNullIfEmpty(biConnector.ReadPreference),
		},
	}
}

func newTFSnapshotBackupPolicyRSModel(ctx context.Context, conn *matlas.Client, projectID, clusterName string) (types.List, error) {
	res, err := newTFSnapshotBackupPolicyModel(ctx, conn, projectID, clusterName)
	if err != nil {
		return types.ListNull(tfSnapshotBackupPolicyType), fmt.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
	}

	s, _ := types.ListValueFrom(ctx, tfSnapshotBackupPolicyType, res)
	return s, nil
}

func newTFSnapshotBackupPolicyModel(ctx context.Context, conn *matlas.Client, projectID, clusterName string) ([]*tfSnapshotBackupPolicyModel, error) {
	res := []*tfSnapshotBackupPolicyModel{}

	backupPolicy, response, err := conn.CloudProviderSnapshotBackupPolicies.Get(ctx, projectID, clusterName)

	if err != nil {
		if response.StatusCode == http.StatusNotFound ||
			strings.Contains(err.Error(), "BACKUP_CONFIG_NOT_FOUND") ||
			strings.Contains(err.Error(), "Not Found") ||
			strings.Contains(err.Error(), "404") {
			return res, nil
		}

		return nil, fmt.Errorf(errorSnapshotBackupPolicyRead, clusterName, err)
	}

	res = append(res, &tfSnapshotBackupPolicyModel{
		ClusterID:             conversion.StringNullIfEmpty(backupPolicy.ClusterID),
		ClusterName:           conversion.StringNullIfEmpty(backupPolicy.ClusterName),
		NextSnapshot:          conversion.StringNullIfEmpty(backupPolicy.NextSnapshot),
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
			ID:         conversion.StringNullIfEmpty(pe.ID),
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
			ID:                conversion.StringNullIfEmpty(pe.ID),
			FrequencyInterval: types.Int64Value(cast.ToInt64(pe.FrequencyInterval)),
			FrequencyType:     conversion.StringNullIfEmpty(pe.FrequencyType),
			RetentionUnit:     conversion.StringNullIfEmpty(pe.RetentionUnit),
			RetentionValue:    types.Int64Value(cast.ToInt64(pe.RetentionValue)),
		}
	}
	s, _ := types.ListValueFrom(ctx, tfSnapshotPolicyItemType, res)
	return s
}

func newTFConnectionStringsModel(ctx context.Context, connString *matlas.ConnectionStrings) []tfConnectionStringModel {
	res := []tfConnectionStringModel{}

	if connString != nil {
		res = append(res, tfConnectionStringModel{
			Standard:        conversion.StringNullIfEmpty(connString.Standard),
			StandardSrv:     conversion.StringNullIfEmpty(connString.StandardSrv),
			Private:         conversion.StringNullIfEmpty(connString.Private),
			PrivateSrv:      conversion.StringNullIfEmpty(connString.PrivateSrv),
			PrivateEndpoint: newTFPrivateEndpointModel(ctx, connString.PrivateEndpoint),
		})
	}
	return res
}

func newTFConnectionStringsModelList(ctx context.Context, connString *matlas.ConnectionStrings) types.List {
	res := newTFConnectionStringsModel
	s, _ := types.ListValueFrom(ctx, tfConnectionStringType, res)
	return s
}

func newTFPrivateEndpointModel(ctx context.Context, privateEndpoints []matlas.PrivateEndpoint) types.List {
	res := make([]tfPrivateEndpointModel, len(privateEndpoints))

	for i, pe := range privateEndpoints {
		res[i] = tfPrivateEndpointModel{
			ConnectionString:                  conversion.StringNullIfEmpty(pe.ConnectionString),
			SrvConnectionString:               conversion.StringNullIfEmpty(pe.SRVConnectionString),
			SrvShardOptimizedConnectionString: conversion.StringNullIfEmpty(pe.SRVShardOptimizedConnectionString),
			EndpointType:                      conversion.StringNullIfEmpty(pe.Type),
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
			Region:       conversion.StringNullIfEmpty(e.Region),
			ProviderName: conversion.StringNullIfEmpty(e.ProviderName),
			EndpointID:   conversion.StringNullIfEmpty(e.EndpointID),
		}
	}
	s, _ := types.ListValueFrom(ctx, tfEndpointType, res)
	return s
}

func newAtlasProcessArgs(tfModel *tfAdvancedConfigurationModel) *matlas.ProcessArgs {
	res := &matlas.ProcessArgs{}

	if v := tfModel.DefaultReadConcern; !v.IsUnknown() {
		res.DefaultReadConcern = v.ValueString()
	}
	if v := tfModel.DefaultWriteConcern; !v.IsUnknown() {
		res.DefaultWriteConcern = v.ValueString()
	}

	if v := tfModel.FailIndexKeyTooLong; !v.IsUnknown() {
		res.FailIndexKeyTooLong = v.ValueBoolPointer()
	}

	if v := tfModel.JavascriptEnabled; !v.IsUnknown() {
		res.JavascriptEnabled = v.ValueBoolPointer()
	}

	if v := tfModel.MinimumEnabledTLSProtocol; !v.IsUnknown() {
		res.MinimumEnabledTLSProtocol = v.ValueString()
	}

	if v := tfModel.NoTableScan; !v.IsUnknown() {
		res.NoTableScan = v.ValueBoolPointer()
	}

	if v := tfModel.SampleSizeBiConnector; !v.IsUnknown() {
		res.SampleSizeBIConnector = v.ValueInt64Pointer()
	}

	if v := tfModel.SampleRefreshIntervalBiConnector; !v.IsUnknown() {
		res.SampleRefreshIntervalBIConnector = v.ValueInt64Pointer()
	}

	if v := tfModel.OplogSizeMB; !v.IsUnknown() {
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

	if v := tfModel.TransactionLifetimeLimitSeconds; !v.IsUnknown() {
		if transactionLimitSeconds := v.ValueInt64(); transactionLimitSeconds > 0 {
			res.TransactionLifetimeLimitSeconds = v.ValueInt64Pointer()
		} else {
			log.Printf(errorClusterSetting, `transaction_lifetime_limit_seconds`, "", cast.ToString(transactionLimitSeconds))
		}
	}

	return res
}

func newAtlasTags(list []*tfTagModel) []*matlas.Tag {
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
		if providerDiskIops := tfClusterModel.ProviderDiskIops; !providerDiskIops.IsUnknown() && !strings.Contains(providerSettings.InstanceSizeName, "NVME") {
			providerSettings.DiskIOPS = providerDiskIops.ValueInt64Pointer()
		}

		providerSettings.EncryptEBSVolume = pointy.Bool(true)
	}

	return providerSettings, nil
}

// https://github.com/mongodb/terraform-provider-mongodbatlas/pull/463
func updateAtlasReplicationSpecs(tfClusterModel *tfClusterRSModel) ([]matlas.ReplicationSpec, error) {
	rSpecs := make([]matlas.ReplicationSpec, 0)

	for _, repSpec := range tfClusterModel.ReplicationSpecs {
		regionsConfig, err := newAtlasRegionsConfig(repSpec.RegionsConfig)
		if err != nil {
			return nil, err
		}

		rSpec := matlas.ReplicationSpec{
			ID:            repSpec.ID.ValueString(),
			NumShards:     repSpec.NumShards.ValueInt64Pointer(),
			ZoneName:      repSpec.ZoneName.ValueString(),
			RegionsConfig: regionsConfig,
		}
		rSpecs = append(rSpecs, rSpec)
	}

	return rSpecs, nil
}

func updatedAtlasRegionsConfig(regions []tfRegionConfigModel, originalRegion, replaceRegion string) (map[string]matlas.RegionsConfig, error) {
	regionsConfig := make(map[string]matlas.RegionsConfig)

	for _, region := range regions {
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

func newAtlasRegionsConfig(regions []tfRegionConfigModel) (map[string]matlas.RegionsConfig, error) {
	regionsConfig := make(map[string]matlas.RegionsConfig)

	for _, region := range regions {
		r, err := valRegion(region.RegionName.ValueString())
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

func validateClusterConfig(ctx context.Context, plan *tfClusterRSModel, response *resource.CreateResponse) {
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

	var acmodel []tfAdvancedConfigurationModel
	plan.AdvancedConfiguration.ElementsAs(ctx, &acmodel, true)
	if len(acmodel) > 0 {
		if oplogSizeMB := acmodel[0].OplogSizeMB; !oplogSizeMB.IsUnknown() {
			if cast.ToInt64(oplogSizeMB.ValueInt64()) <= 0 {
				response.Diagnostics.AddError(errorInvalidCreateValues, "`advanced_configuration.oplog_size_mb` cannot be <= 0")
			}
		}
	}
}

type tfClusterRSModel struct {
	DiskSizeGb                                types.Float64               `tfsdk:"disk_size_gb"`
	AdvancedConfiguration                     types.List                  `tfsdk:"advanced_configuration"`
	ConnectionStrings                         types.List                  `tfsdk:"connection_strings"`
	SnapshotBackupPolicy                      types.List                  `tfsdk:"snapshot_backup_policy"`
	ProviderName                              types.String                `tfsdk:"provider_name"`
	ClusterType                               types.String                `tfsdk:"cluster_type"`
	ClusterID                                 types.String                `tfsdk:"cluster_id"`
	MongoURIWithOptions                       types.String                `tfsdk:"mongo_uri_with_options"`
	ContainerID                               types.String                `tfsdk:"container_id"`
	VersionReleaseSystem                      types.String                `tfsdk:"version_release_system"`
	EncryptionAtRestProvider                  types.String                `tfsdk:"encryption_at_rest_provider"`
	ID                                        types.String                `tfsdk:"id"`
	MongoDBMajorVersion                       types.String                `tfsdk:"mongo_db_major_version"`
	MongoDBMajorVersionFormatted              types.String                `tfsdk:"mongo_db_major_version_formatted"`
	MongoDBVersion                            types.String                `tfsdk:"mongo_db_version"`
	MongoURI                                  types.String                `tfsdk:"mongo_uri"`
	MongoURIUpdated                           types.String                `tfsdk:"mongo_uri_updated"`
	ProviderAutoScalingComputeMaxInstanceSize types.String                `tfsdk:"provider_auto_scaling_compute_max_instance_size"`
	Name                                      types.String                `tfsdk:"name"`
	ProjectID                                 types.String                `tfsdk:"project_id"`
	ProviderVolumeType                        types.String                `tfsdk:"provider_volume_type"`
	ProviderRegionName                        types.String                `tfsdk:"provider_region_name"`
	Timeouts                                  timeouts.Value              `tfsdk:"timeouts"`
	ProviderInstanceSizeName                  types.String                `tfsdk:"provider_instance_size_name"`
	SrvAddress                                types.String                `tfsdk:"srv_address"`
	ProviderAutoScalingComputeMinInstanceSize types.String                `tfsdk:"provider_auto_scaling_compute_min_instance_size"`
	StateName                                 types.String                `tfsdk:"state_name"`
	ProviderDiskTypeName                      types.String                `tfsdk:"provider_disk_type_name"`
	BackingProviderName                       types.String                `tfsdk:"backing_provider_name"`
	Labels                                    []tfLabelModel              `tfsdk:"labels"`
	BiConnectorConfig                         []*tfBiConnectorConfigModel `tfsdk:"bi_connector_config"`
	ReplicationSpecs                          []*tfReplicationSpecModel   `tfsdk:"replication_specs"`
	Tags                                      []*tfTagModel               `tfsdk:"tags"`
	ReplicationFactor                         types.Int64                 `tfsdk:"replication_factor"`
	ProviderDiskIops                          types.Int64                 `tfsdk:"provider_disk_iops"`
	NumShards                                 types.Int64                 `tfsdk:"num_shards"`
	TerminationProtectionEnabled              types.Bool                  `tfsdk:"termination_protection_enabled"`
	PitEnabled                                types.Bool                  `tfsdk:"pit_enabled"`
	AutoScalingDiskGBEnabled                  types.Bool                  `tfsdk:"auto_scaling_disk_gb_enabled"`
	CloudBackup                               types.Bool                  `tfsdk:"cloud_backup"`
	Paused                                    types.Bool                  `tfsdk:"paused"`
	RetainBackupsEnabled                      types.Bool                  `tfsdk:"retain_backups_enabled"`
	BackupEnabled                             types.Bool                  `tfsdk:"backup_enabled"`
	ProviderEncryptEbsVolume                  types.Bool                  `tfsdk:"provider_encrypt_ebs_volume"`
	AutoScalingComputeScaleDownEnabled        types.Bool                  `tfsdk:"auto_scaling_compute_scale_down_enabled"`
	AutoScalingComputeEnabled                 types.Bool                  `tfsdk:"auto_scaling_compute_enabled"`
	ProviderEncryptEbsVolumeFlag              types.Bool                  `tfsdk:"provider_encrypt_ebs_volume_flag"`
}

type tfSnapshotBackupPolicyModel struct {
	ClusterID             types.String `tfsdk:"cluster_id"`
	ClusterName           types.String `tfsdk:"cluster_name"`
	NextSnapshot          types.String `tfsdk:"next_snapshot"`
	Policies              types.List   `tfsdk:"policies"`
	ReferenceHourOfDay    types.Int64  `tfsdk:"reference_hour_of_day"`
	ReferenceMinuteOfHour types.Int64  `tfsdk:"reference_minute_of_hour"`
	RestoreWindowDays     types.Int64  `tfsdk:"restore_window_days"`
	UpdateSnapshots       types.Bool   `tfsdk:"update_snapshots"`
}

var tfSnapshotBackupPolicyType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"cluster_id":               types.StringType,
	"cluster_name":             types.StringType,
	"next_snapshot":            types.StringType,
	"policies":                 types.ListType{ElemType: tfSnapshotPolicyType},
	"reference_hour_of_day":    types.Int64Type,
	"reference_minute_of_hour": types.Int64Type,
	"restore_window_days":      types.Int64Type,
	"update_snapshots":         types.BoolType,
}}

type tfSnapshotPolicyModel struct {
	ID         types.String `tfsdk:"id"`
	PolicyItem types.List   `tfsdk:"policy_item"`
}

var tfSnapshotPolicyType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":          types.StringType,
	"policy_item": types.ListType{ElemType: tfSnapshotPolicyItemType},
}}

type tfSnapshotPolicyItemModel struct {
	ID                types.String `tfsdk:"id"`
	FrequencyType     types.String `tfsdk:"frequency_type"`
	RetentionUnit     types.String `tfsdk:"retention_unit"`
	FrequencyInterval types.Int64  `tfsdk:"frequency_interval"`
	RetentionValue    types.Int64  `tfsdk:"retention_value"`
}

var tfSnapshotPolicyItemType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":                 types.StringType,
	"frequency_type":     types.StringType,
	"retention_unit":     types.StringType,
	"frequency_interval": types.Int64Type,
	"retention_value":    types.Int64Type,
}}

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

var tfAdvancedConfigurationType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"default_read_concern":                 types.StringType,
	"default_write_concern":                types.StringType,
	"minimum_enabled_tls_protocol":         types.StringType,
	"oplog_size_mb":                        types.Int64Type,
	"oplog_min_retention_hours":            types.Int64Type,
	"sample_size_bi_connector":             types.Int64Type,
	"sample_refresh_interval_bi_connector": types.Int64Type,
	"transaction_lifetime_limit_seconds":   types.Int64Type,

	"fail_index_key_too_long": types.BoolType,
	"javascript_enabled":      types.BoolType,
	"no_table_scan":           types.BoolType,
}}

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
	Standard    types.String `tfsdk:"standard"`
	StandardSrv types.String `tfsdk:"standard_srv"`
	Private     types.String `tfsdk:"private"`
	PrivateSrv  types.String `tfsdk:"private_srv"`
	// PrivateEndpoint []tfPrivateEndpointModel `tfsdk:"private_endpoint"`
	PrivateEndpoint types.List `tfsdk:"private_endpoint"`
}

var tfConnectionStringType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"standard":         types.StringType,
	"standard_srv":     types.StringType,
	"private":          types.StringType,
	"private_srv":      types.StringType,
	"private_endpoint": types.ListType{ElemType: tfPrivateEndpointType},
}}

type tfPrivateEndpointModel struct {
	ConnectionString                  types.String `tfsdk:"connection_string"`
	SrvConnectionString               types.String `tfsdk:"srv_connection_string"`
	SrvShardOptimizedConnectionString types.String `tfsdk:"srv_shard_optimized_connection_string"`
	EndpointType                      types.String `tfsdk:"type"`
	// Endpoints                         []tfEndpointModel `tfsdk:"endpoints"`
	Endpoints types.List `tfsdk:"endpoints"`
}

var tfPrivateEndpointType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"connection_string":                     types.StringType,
	"endpoints":                             types.ListType{ElemType: tfEndpointType},
	"srv_connection_string":                 types.StringType,
	"srv_shard_optimized_connection_string": types.StringType,
	"type":                                  types.StringType,
}}

type tfEndpointModel struct {
	EndpointID   types.String `tfsdk:"endpoint_id"`
	ProviderName types.String `tfsdk:"provider_name"`
	Region       types.String `tfsdk:"region"`
}

var tfEndpointType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"endpoint_id":   types.StringType,
	"provider_name": types.StringType,
	"region":        types.StringType,
},
}

// handling replication_specs
func areTFReplicationSpecSlicesEqual(a, b []tfReplicationSpecModel) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Sort(ByReplicationSpecModel(a))
	sort.Sort(ByReplicationSpecModel(b))

	for i := range a {
		// Also sort the underlying RegionsConfig
		sort.Sort(ByRegionConfigModel(a[i].RegionsConfig))
		sort.Sort(ByRegionConfigModel(b[i].RegionsConfig))

		if !a[i].Equal(b[i]) {
			return false
		}
	}

	return true
}

func (r tfReplicationSpecModel) Equal(other tfReplicationSpecModel) bool {
	if !r.ID.Equal(other.ID) || !r.ZoneName.Equal(other.ZoneName) || !r.NumShards.Equal(other.NumShards) {
		return false
	}

	if len(r.RegionsConfig) != len(other.RegionsConfig) {
		return false
	}

	for i, region := range r.RegionsConfig {
		if !region.Equal(other.RegionsConfig[i]) {
			return false
		}
	}

	return true
}

func (r tfRegionConfigModel) Equal(other tfRegionConfigModel) bool {
	return r.RegionName.Equal(other.RegionName) &&
		r.ElectableNodes.Equal(other.ElectableNodes) &&
		r.Priority.Equal(other.Priority) &&
		r.ReadOnlyNodes.Equal(other.ReadOnlyNodes) &&
		r.AnalyticsNodes.Equal(other.AnalyticsNodes)
}

type ByReplicationSpecModel []tfReplicationSpecModel

func (a ByReplicationSpecModel) Len() int      { return len(a) }
func (a ByReplicationSpecModel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByReplicationSpecModel) Less(i, j int) bool {
	return a[i].ID.ValueString() < a[j].ID.ValueString()
}

type ByRegionConfigModel []tfRegionConfigModel

func (a ByRegionConfigModel) Len() int      { return len(a) }
func (a ByRegionConfigModel) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByRegionConfigModel) Less(i, j int) bool {
	return a[i].RegionName.ValueString() < a[j].RegionName.ValueString()
}
