package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ resource.ResourceWithConfigure = &advancedClusterRS{}
var _ resource.ResourceWithImportState = &advancedClusterRS{}

type advancedClusterRS struct {
	config.RSCommon
}

// TODO rename to Resource() after deleting old resource
func Fw_Resource() resource.Resource {
	return &advancedClusterRS{
		RSCommon: config.RSCommon{
			ResourceName: AdvancedClusterResourceName,
		},
	}
}

func (r *advancedClusterRS) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": // TODO framework.IDAttribute()
			schema.StringAttribute{
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
			"backup_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"retain_backups_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster",
			},
			"cluster_type": schema.StringAttribute{
				Required: true,
			},
			"connection_strings": advClusterDSConnectionStringSchemaAttr(),
			"create_date": schema.StringAttribute{
				Computed: true,
			},
			"disk_size_gb": schema.Float64Attribute{
				Optional: true,
				Computed: true,
			},
			"encryption_at_rest_provider": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			// https://discuss.hashicorp.com/t/is-it-possible-to-have-statefunc-like-behavior-with-the-plugin-framework/58377/2
			"mongo_db_major_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// StateFunc: FormatMongoDBMajorVersion,
			},
			"mongo_db_version": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"paused": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"pit_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},

			"root_cert_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"state_name": schema.StringAttribute{
				Computed: true,
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"version_release_system": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("LTS", "CONTINUOUS"),
				},
			},
			"accept_data_risks_and_force_replica_set_reconfig": schema.StringAttribute{
				Optional:    true,
				Description: "Submit this field alongside your topology reconfiguration to request a new regional outage resistant topology",
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
						},
						"value": schema.StringAttribute{
							Optional: true,
						},
					},
				},
				DeprecationMessage: fmt.Sprintf(constant.DeprecationParamByDateWithReplacement, "September 2024", "tags"),
			},
			"replication_specs": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"container_id": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"num_shards": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							Default:  int64default.StaticInt64(1),
							Validators: []validator.Int64{
								int64validator.Between(1, 50),
							},
						},
						"zone_name": schema.StringAttribute{
							Optional: true,
							Computed: true,
							Default:  stringdefault.StaticString("ZoneName managed by Terraform"),
						},
					},
					Blocks: map[string]schema.Block{
						"region_configs": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"backing_provider_name": schema.StringAttribute{
										Optional: true,
									},
									"priority": schema.Int64Attribute{
										Required: true,
									},
									"provider_name": schema.StringAttribute{
										Required: true,
									},
									"region_name": schema.StringAttribute{
										Required: true,
									},
								},
								Blocks: map[string]schema.Block{
									"analytics_auto_scaling": advancedClusterRSRegionConfigAutoScalingSpecsBlock(),
									"auto_scaling":           advancedClusterRSRegionConfigAutoScalingSpecsBlock(),
									"analytics_specs":        advancedClusterRSRegionConfigSpecsBlock(),
									"electable_specs":        advancedClusterRSRegionConfigSpecsBlock(),
									"read_only_specs":        advancedClusterRSRegionConfigSpecsBlock(),
								},
							},
							Validators: []validator.List{
								listvalidator.IsRequired(),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.IsRequired(),
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

func advancedClusterRSRegionConfigSpecsBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
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
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func advancedClusterRSRegionConfigAutoScalingSpecsBlock() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"compute_enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				"compute_max_instance_size": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"compute_min_instance_size": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"compute_scale_down_enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},
				"disk_gb_enabled": schema.BoolAttribute{
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

func (r *advancedClusterRS) Create(ctx context.Context, request resource.CreateRequest, response *resource.CreateResponse) {
	conn := r.client.Atlas
	var plan tfAdvancedClusterRSModel

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	createTimeout := r.CreateTimeout(ctx, plan.Timeouts)

	plan.ID = types.StringValue("TODO")

	// TODO: initialize and set newState

	// set state to fully populated data
	response.Diagnostics.Append(response.State.Set(ctx, &newState)...)
}

func (r *advancedClusterRS) Read(ctx context.Context, request resource.ReadRequest, response *resource.ReadResponse) {
	conn := r.client.Atlas
	var state tfAdvancedClusterRSModel

	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	// TODO: initialize and set newState

	// save updated data into terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &newState)...)
}

func (r *advancedClusterRS) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	conn := r.client.Atlas
	var state, plan tfAdvancedClusterRSModel

	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}
	updateTimeout := r.UpdateTimeout(ctx, plan.Timeouts)

	// save updated data into terraform state
	response.Diagnostics.Append(response.State.Set(ctx, &plan)...)
}

func (r *advancedClusterRS) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	conn := r.client.Atlas
	var state tfAdvancedClusterRSModel

	response.Diagnostics.Append(request.State.Get(ctx, &state)...)

	if response.Diagnostics.HasError() {
		return
	}
	deleteTimeout := r.DeleteTimeout(ctx, state.Timeouts)

	tflog.Debug(ctx, "deleting TODO", map[string]interface{}{
		"id": state.ID.ValueString(),
	})
}

// ImportState is called when the provider must import the state of a resource instance.
// This method must return enough state so the Read method can properly refresh the full resource.
//
// If setting an attribute with the import identifier, it is recommended to use the ImportStatePassthroughID() call in this method.
func (r *advancedClusterRS) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), request, response)
}

type tfAdvancedClusterRSModel struct {
	BackupEnabled                types.Bool    `tfsdk:"backup_enabled"`
	ClusterID                    types.String  `tfsdk:"cluster_id"`
	ClusterType                  types.String  `tfsdk:"cluster_type"`
	CreateDate                   types.String  `tfsdk:"create_date"`
	DiskSizeGb                   types.Float64 `tfsdk:"disk_size_gb"`
	EncryptionAtRestProvider     types.String  `tfsdk:"encryption_at_rest_provider"`
	ID                           types.String  `tfsdk:"id"`
	MongoDbMajorVersion          types.String  `tfsdk:"mongo_db_major_version"`
	MongoDbVersion               types.String  `tfsdk:"mongo_db_version"`
	Name                         types.String  `tfsdk:"name"`
	Paused                       types.Bool    `tfsdk:"paused"`
	PitEnabled                   types.Bool    `tfsdk:"pit_enabled"`
	ProjectID                    types.String  `tfsdk:"project_id"`
	RetainBackupsEnabled         types.Bool    `tfsdk:"retain_backups_enabled"`
	RootCertType                 types.String  `tfsdk:"root_cert_type"`
	StateName                    types.String  `tfsdk:"state_name"`
	TerminationProtectionEnabled types.Bool    `tfsdk:"termination_protection_enabled"`
	VersionReleaseSystem         types.String  `tfsdk:"version_release_system"`

	Labels                types.Set  `tfsdk:"labels"`
	Tags                  types.Set  `tfsdk:"tags"`
	ReplicationSpecs      types.List `tfsdk:"replication_specs"`
	BiConnectorConfig     types.List `tfsdk:"bi_connector_config"`
	ConnectionStrings     types.List `tfsdk:"connection_strings"`
	AdvancedConfiguration types.List `tfsdk:"advanced_configuration"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

type tfReplicationSpecRSModel struct {
	RegionsConfigs types.List   `tfsdk:"region_configs"`
	ContainerID    types.Map    `tfsdk:"container_id"`
	ID             types.String `tfsdk:"id"`
	ZoneName       types.String `tfsdk:"zone_name"`
	NumShards      types.Int64  `tfsdk:"num_shards"`
}
