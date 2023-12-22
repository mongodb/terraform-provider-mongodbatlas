package advancedcluster

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"
	"golang.org/x/exp/slices"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mwielbut/pointy"
	"github.com/spf13/cast"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/planmodifiers"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/utility"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorInvalidCreateValues = "Invalid values. Unable to CREATE advanced_cluster"
	defaultTimeout           = (3 * time.Hour)
	defaultInt               = 0
	defaultString            = ""
	defaultZoneName          = "ZoneName managed by Terraform"
)

var _ resource.ResourceWithConfigure = &advancedClusterRS{}
var _ resource.ResourceWithImportState = &advancedClusterRS{}
var _ resource.ResourceWithUpgradeState = &advancedClusterRS{}

type advancedClusterRS struct {
	config.RSCommon
}

// UpgradeState implements resource.ResourceWithUpgradeState.
func (*advancedClusterRS) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	schemaV0 := TPFResourceV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &schemaV0,
			StateUpgrader: upgradeAdvClusterResourceStateV0toV1,
		},
	}
}

// TODO rename to Resource() after deleting old resource
func TPFResource() resource.Resource {
	return &advancedClusterRS{
		RSCommon: config.RSCommon{
			ResourceName: AdvancedClusterResourceName,
		},
	}
}

func (r *advancedClusterRS) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"backup_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				// PlanModifiers: []planmodifier.Bool{
				// 	boolplanmodifier.UseStateForUnknown(),
				// },
			},
			"retain_backups_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster",
			},
			"cluster_type": schema.StringAttribute{
				Required: true,
			},
			"connection_strings": advClusterRSConnectionStringSchemaComputed(),
			"create_date": schema.StringAttribute{
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"disk_size_gb": schema.Float64Attribute{
				Optional: true,
				Computed: true,
				// PlanModifiers: []planmodifier.Float64{
				// 	float64planmodifier.UseStateForUnknown(),
				// },
			},
			"encryption_at_rest_provider": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			// https://developer.hashicorp.com/terraform/plugin/framework/migrating/resources/crud#planned-value-does-not-match-config-value
			"mongo_db_major_version": schema.StringAttribute{
				CustomType: customtypes.DBVersionStringType{},
				Optional:   true,
				Computed:   true,
				PlanModifiers: []planmodifier.String{
					planmodifiers.DBVersion(),
				},
			},
			"mongo_db_version": schema.StringAttribute{
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
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
				// PlanModifiers: []planmodifier.Bool{
				// 	boolplanmodifier.UseStateForUnknown(),
				// },
			},
			"pit_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				// PlanModifiers: []planmodifier.Bool{
				// 	boolplanmodifier.UseStateForUnknown(),
				// },
			},
			"root_cert_type": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"state_name": schema.StringAttribute{
				Computed: true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"termination_protection_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				// PlanModifiers: []planmodifier.Bool{
				// 	boolplanmodifier.UseStateForUnknown(),
				// },
			},
			"version_release_system": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("LTS", "CONTINUOUS"),
				},
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
			"accept_data_risks_and_force_replica_set_reconfig": schema.StringAttribute{
				Optional:    true,
				Description: "Submit this field alongside your topology reconfiguration to request a new regional outage resistant topology",
			},
			"advanced_configuration": advClusterRSAdvancedConfigurationSchema(),
			"bi_connector_config":    advClusterRSBiConnectorConfigSchema(),
			"replication_specs":      advClusterRSReplicationSpecsSchema(),
			"labels": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
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
			"tags": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
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
		// Blocks: map[string]schema.Block{
		// 	"labels": schema.SetNestedBlock{
		// 		NestedObject: schema.NestedBlockObject{
		// 			Attributes: map[string]schema.Attribute{
		// 				"key": schema.StringAttribute{
		// 					Optional: true,
		// 				},
		// 				"value": schema.StringAttribute{
		// 					Optional: true,
		// 				},
		// 			},
		// 		},
		// 		DeprecationMessage: fmt.Sprintf(constant.DeprecationParamByDateWithReplacement, "September 2024", "tags"),
		// 	},
		// 	"tags": schema.SetNestedBlock{
		// 		NestedObject: schema.NestedBlockObject{
		// 			Attributes: map[string]schema.Attribute{
		// 				"key": schema.StringAttribute{
		// 					Required: true,
		// 				},
		// 				"value": schema.StringAttribute{
		// 					Required: true,
		// 				},
		// 			},
		// 		},
		// 	},
		// },
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

func advClusterRSConnectionStringSchemaComputed() schema.ListNestedAttribute {
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
		// remove
		// PlanModifiers: []planmodifier.List{
		// 	listplanmodifier.UseStateForUnknown(),
		// },
	}
}

func advClusterRSBiConnectorConfigSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Bool{
					// 	boolplanmodifier.UseStateForUnknown(),
					// },
				},
				"read_preference": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
			},
			// PlanModifiers: []planmodifier.Object{
			// 	objectplanmodifier.UseStateForUnknown(),
			// },
		},
		// Default: listdefault.StaticValue(defaultBiConnectorConfig(ctx)),
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		// PlanModifiers: []planmodifier.List{
		// 	listplanmodifier.UseStateForUnknown(),
		// },
	}
}

func advClusterRSAdvancedConfigurationSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"default_read_concern": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
				"default_write_concern": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
				"fail_index_key_too_long": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Bool{
					// 	boolplanmodifier.UseStateForUnknown(),
					// },
				},
				"javascript_enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Bool{
					// 	boolplanmodifier.UseStateForUnknown(),
					// },
				},
				"minimum_enabled_tls_protocol": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
				"no_table_scan": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Bool{
					// 	boolplanmodifier.UseStateForUnknown(),
					// },
				},
				"oplog_min_retention_hours": schema.Int64Attribute{
					Optional: true,
					// PlanModifiers: []planmodifier.Int64{
					// 	int64planmodifier.UseStateForUnknown(),
					// },
				},
				"oplog_size_mb": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Int64{
					// 	int64planmodifier.UseStateForUnknown(),
					// },
				},
				"sample_refresh_interval_bi_connector": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Int64{
					// 	int64planmodifier.UseStateForUnknown()
					// },
				},
				"sample_size_bi_connector": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Int64{
					// 	int64planmodifier.UseStateForUnknown(),
					// },
				},
				"transaction_lifetime_limit_seconds": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Int64{
					// 	int64planmodifier.UseStateForUnknown(),
					// },
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		// remove
		// PlanModifiers: []planmodifier.List{
		// 	listplanmodifier.UseStateForUnknown(),
		// },
	}
}

func advClusterRSReplicationSpecsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"container_id": schema.MapAttribute{
					ElementType: types.StringType,
					Computed:    true,
					// PlanModifiers: []planmodifier.Map{
					// 	mapplanmodifier.UseStateForUnknown(),
					// },
				},
				"id": schema.StringAttribute{
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
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
					// Default:  stringdefault.StaticString("ZoneName managed by Terraform"),
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
				"region_configs": schema.ListNestedAttribute{
					Optional: true,
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"backing_provider_name": schema.StringAttribute{
								Optional: true,
								Computed: true,
								// PlanModifiers: []planmodifier.String{
								// 	stringplanmodifier.UseStateForUnknown(),
								// },
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
							"analytics_auto_scaling": advClusterRSRegionConfigAutoScalingSpecsBlock(),
							"auto_scaling":           advClusterRSRegionConfigAutoScalingSpecsBlock(),
							"analytics_specs":        advClusterRSRegionConfigSpecsBlock(),
							"electable_specs":        advClusterRSRegionConfigSpecsBlock(),
							"read_only_specs":        advClusterRSRegionConfigSpecsBlock(),
						},
					},
					Validators: []validator.List{
						listvalidator.IsRequired(),
					},
					// PlanModifiers: []planmodifier.List{
					// 	listplanmodifier.UseStateForUnknown(),
					// },
				},
			},
			// PlanModifiers: []planmodifier.Object{
			// 	objectplanmodifier.UseStateForUnknown(),
			// },
		},
		Validators: []validator.List{
			listvalidator.IsRequired(),
		},
		// PlanModifiers: []planmodifier.List{
		// 	listplanmodifier.UseStateForUnknown(),
		// },
	}
}

func advClusterRSRegionConfigSpecsBlock() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"disk_iops": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Int64{
					// 	int64planmodifier.UseStateForUnknown(),
					// },
				},
				"ebs_volume_type": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
				"instance_size": schema.StringAttribute{
					Required: true,
				},
				"node_count": schema.Int64Attribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Int64{
					// 	int64planmodifier.UseStateForUnknown(),
					// },
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		// PlanModifiers: []planmodifier.List{
		// 	listplanmodifier.UseStateForUnknown(),
		// },
	}
}

func advClusterRSRegionConfigAutoScalingSpecsBlock() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		// PlanModifiers: []planmodifier.List{
		// 	listplanmodifier.UseStateForUnknown(),
		// },
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"compute_enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Bool{
					// 	boolplanmodifier.UseStateForUnknown(),
					// },
				},
				"compute_max_instance_size": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
				"compute_min_instance_size": schema.StringAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.String{
					// 	stringplanmodifier.UseStateForUnknown(),
					// },
				},
				"compute_scale_down_enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Bool{
					// 	boolplanmodifier.UseStateForUnknown(),
					// },
				},
				"disk_gb_enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
					// PlanModifiers: []planmodifier.Bool{
					// 	boolplanmodifier.UseStateForUnknown(),
					// },
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (r *advancedClusterRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	conn := r.Client.Atlas
	var plan, tfConfig tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ------validations
	if plan.AcceptDataRisksAndForceReplicaSetReconfig.ValueString() != "" {
		resp.Diagnostics.AddError(errorInvalidCreateValues, "accept_data_risks_and_force_replica_set_reconfig can not be set in creation, only in update")
		return
	}
	// We need to validate the oplog_size_mb attr of the advanced configuration option to show the error
	// before that the cluster is created
	var advConfig *matlas.ProcessArgs
	if v := plan.AdvancedConfiguration; !v.IsNull() {
		advConfig = newAdvancedConfiguration(ctx, v)
		if advConfig != nil && advConfig.OplogSizeMB != nil && *advConfig.OplogSizeMB <= 0 {
			resp.Diagnostics.AddError(errorInvalidCreateValues, "`advanced_configuration.oplog_size_mb` cannot be <= 0")
			return
		}
	}
	if v := plan.Labels; !v.IsNull() && ContainsLabelOrKey(newLabels(ctx, v), DefaultLabel) {
		resp.Diagnostics.AddError(errorInvalidCreateValues, "you should not set `Infrastructure Tool` label, it is used for internal purposes")
		return
	}
	// ------validations end

	projectID := plan.ProjectID.ValueString()

	request := &matlas.AdvancedCluster{
		Name:             plan.Name.ValueString(),
		ClusterType:      plan.ClusterType.ValueString(),
		ReplicationSpecs: newReplicationSpecs(ctx, plan.ReplicationSpecs),
	}

	if v := plan.BackupEnabled; !v.IsUnknown() {
		request.BackupEnabled = v.ValueBoolPointer()
	}

	if v := plan.BiConnectorConfig; !v.IsUnknown() {
		request.BiConnector = newBiConnectorConfig(ctx, plan.BiConnectorConfig)
	}

	if v := plan.DiskSizeGb; !v.IsUnknown() {
		request.DiskSizeGB = v.ValueFloat64Pointer()
	}

	if v := plan.EncryptionAtRestProvider; !v.IsUnknown() {
		request.EncryptionAtRestProvider = v.ValueString()
	}

	request.Labels = append(newLabels(ctx, plan.Labels), DefaultLabel)

	request.Tags = newTags(ctx, plan.Tags)

	if v := plan.MongoDBMajorVersion; !v.IsUnknown() {
		request.MongoDBMajorVersion = utility.FormatMongoDBMajorVersion(v.ValueString())
	}

	if v := plan.PitEnabled; !v.IsUnknown() {
		request.PitEnabled = v.ValueBoolPointer()
	}
	if v := plan.RootCertType; !v.IsUnknown() {
		request.RootCertType = v.ValueString()
	}
	if v := plan.TerminationProtectionEnabled; !v.IsUnknown() {
		request.TerminationProtectionEnabled = v.ValueBoolPointer()
	}
	if v := plan.VersionReleaseSystem; !v.IsUnknown() {
		request.VersionReleaseSystem = v.ValueString()
	}

	// TODO undo
	cluster, _, err := conn.AdvancedClusters.Create(ctx, projectID, request)
	// cluster, _, err := conn.AdvancedClusters.Get(ctx, projectID, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", fmt.Sprintf(errorClusterAdvancedCreate, err))
		return
	}

	timeout, diags := plan.Timeouts.Create(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING", "REPEATING", "PENDING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, cluster.Name, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 1 * time.Minute,
		Delay:      3 * time.Minute,
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to CREATE cluster. Error during create in Atlas", fmt.Sprintf(errorClusterAdvancedCreate, err))
		return
	}

	/*
		So far, the cluster has created correctly, so we need to set up
		the advanced configuration option to attach it
	*/
	if advConfig != nil {
		_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, cluster.Name, advConfig)
		if err != nil {
			resp.Diagnostics.AddError("Error during cluster CREATE", fmt.Sprintf(errorAdvancedClusterAdvancedConfUpdate, cluster.Name, err))
		}
	}

	// To pause a cluster
	if v := plan.Paused.ValueBool(); v {
		request = &matlas.AdvancedCluster{
			Paused: pointy.Bool(v),
		}

		_, _, err = updateAdvancedCluster(ctx, conn, request, projectID, cluster.Name, timeout)
		if err != nil {
			resp.Diagnostics.AddError("Error during cluster CREATE. An error occurred attempting to pause cluster in Atlas", fmt.Sprintf(errorClusterAdvancedCreate, err))
			return
		}
	}

	cluster, response, err := conn.AdvancedClusters.Get(ctx, projectID, cluster.Name)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error during cluster READ from Atlas", fmt.Sprintf(errorClusterAdvancedRead, cluster.Name, err.Error()))
		return
	}

	// during READ, mongodb_major_version should match what is in the config
	newClusterModel, diags := newTfAdvClusterRSModel(ctx, conn, cluster, &plan, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterModel)...)
}

func newTfAdvClusterRSModel(ctx context.Context, conn *matlas.Client, cluster *matlas.AdvancedCluster, state *tfAdvancedClusterRSModel, isImport bool) (*tfAdvancedClusterRSModel, diag.Diagnostics) {
	var d, diags diag.Diagnostics
	projectID := cluster.GroupID
	name := cluster.Name

	clusterModel := tfAdvancedClusterRSModel{
		ClusterID:                    types.StringValue(cluster.ID),
		BackupEnabled:                types.BoolPointerValue(cluster.BackupEnabled),
		ClusterType:                  types.StringValue(cluster.ClusterType),
		CreateDate:                   types.StringValue(cluster.CreateDate),
		DiskSizeGb:                   types.Float64PointerValue(cluster.DiskSizeGB),
		EncryptionAtRestProvider:     types.StringValue(cluster.EncryptionAtRestProvider),
		MongoDBMajorVersion:          customtypes.DBVersionStringValue{StringValue: types.StringValue(cluster.MongoDBMajorVersion)},
		MongoDBVersion:               types.StringValue(cluster.MongoDBVersion),
		Name:                         types.StringValue(name),
		Paused:                       types.BoolPointerValue(cluster.Paused),
		PitEnabled:                   types.BoolPointerValue(cluster.PitEnabled),
		RootCertType:                 types.StringValue(cluster.RootCertType),
		StateName:                    types.StringValue(cluster.StateName),
		TerminationProtectionEnabled: types.BoolPointerValue(cluster.TerminationProtectionEnabled),
		VersionReleaseSystem:         types.StringValue(cluster.VersionReleaseSystem),
		AcceptDataRisksAndForceReplicaSetReconfig: conversion.StringNullIfEmpty(cluster.AcceptDataRisksAndForceReplicaSetReconfig),
		ProjectID:            types.StringValue(projectID),
		RetainBackupsEnabled: state.RetainBackupsEnabled,
	}

	clusterModel.ID = types.StringValue(conversion.EncodeStateID(map[string]string{
		"cluster_id":   cluster.ID,
		"project_id":   projectID,
		"cluster_name": name,
	}))

	clusterModel.BiConnectorConfig, d = types.ListValueFrom(ctx, TfBiConnectorConfigType, NewTfBiConnectorConfigModel(cluster.BiConnector))
	diags.Append(d...)

	clusterModel.ConnectionStrings, d = types.ListValueFrom(ctx, tfConnectionStringType, newTfConnectionStringsModel(ctx, cluster.ConnectionStrings))
	diags.Append(d...)

	clusterModel.Labels, d = types.SetValueFrom(ctx, TfLabelType, RemoveDefaultLabel(NewTfLabelsModel(cluster.Labels)))
	if len(clusterModel.Labels.Elements()) == 0 {
		// clusterModel.Labels, d = types.SetValueFrom(ctx, TfLabelType, []TfLabelModel{})
		clusterModel.Labels = types.SetNull(TfLabelType)
	}
	diags.Append(d...)

	clusterModel.Tags, d = types.SetValueFrom(ctx, TfTagType, NewTfTagsModel(&cluster.Tags))
	if len(clusterModel.Tags.Elements()) == 0 {
		// clusterModel.Tags, d = types.SetValueFrom(ctx, TfTagType, []TfTagModel{})
		clusterModel.Tags = types.SetNull(TfTagType)
	}
	diags.Append(d...)

	replicationSpecs, diags := newTfReplicationSpecsRS(ctx, conn, cluster.ReplicationSpecs, state.ReplicationSpecs, projectID)
	diags.Append(d...)

	if diags.HasError() {
		return nil, diags
	}
	clusterModel.ReplicationSpecs, diags = types.ListValueFrom(ctx, tfReplicationSpecRSType, replicationSpecs)
	diags.Append(d...)

	advancedConfiguration, err := NewTfAdvancedConfigurationModelDSFromAtlas(ctx, conn, projectID, name)
	if err != nil {
		diags.AddError("An error occurred when getting advanced_configuration from Atlas", err.Error())
		return nil, diags
	}
	clusterModel.AdvancedConfiguration, diags = types.ListValueFrom(ctx, tfAdvancedConfigurationType, advancedConfiguration)
	if diags.HasError() {
		return nil, diags
	}

	clusterModel.Timeouts = state.Timeouts

	return &clusterModel, diags
}

func (r *advancedClusterRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Client.Atlas

	var isImport bool
	var state tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use the ID only with the IMPORT operation
	if state.ID.ValueString() != "" && (state.ClusterID.ValueString() == "") {
		isImport = true
	}

	ids := conversion.DecodeStateID(state.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	cluster, response, err := conn.AdvancedClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error during cluster READ from Atlas", fmt.Sprintf(errorClusterAdvancedRead, clusterName, err.Error()))
		return
	}

	log.Printf("[DEBUG] GET ClusterAdvanced %+v", cluster)

	newClusterModel, diags := newTfAdvClusterRSModel(ctx, conn, cluster, &state, isImport)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if !isImport {
		newClusterModel.MongoDBMajorVersion = state.MongoDBMajorVersion
	}

	// save updated data into terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterModel)...)
}

func (r *advancedClusterRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	conn := r.Client.Atlas
	var state, plan, tfconfig tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &tfconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids := conversion.DecodeStateID(state.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	timeout, _ := plan.Timeouts.Update(ctx, defaultTimeout)

	if upgradeRequest := TPFgetUpgradeRequest(ctx, &state, &plan); upgradeRequest != nil {
		_, _, err := UpgradeCluster(ctx, conn, upgradeRequest, projectID, clusterName, timeout)

		if err != nil {
			resp.Diagnostics.AddError("Unable to UPDATE cluster. An error occurred while upgrading cluster.", err.Error())
			return
		}
	} else {
		resp.Diagnostics.Append(updateCluster(ctx, conn, &state, &plan, timeout)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// READ
	cluster, response, err := conn.AdvancedClusters.Get(ctx, projectID, clusterName)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error during cluster READ from Atlas", fmt.Sprintf(errorClusterAdvancedRead, clusterName, err.Error()))
		return
	}

	log.Printf("[DEBUG] GET ClusterAdvanced %+v", cluster)
	newClusterModel, diags := newTfAdvClusterRSModel(ctx, conn, cluster, &plan, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterModel)...)
}

func TPFgetUpgradeRequest(ctx context.Context, state, plan *tfAdvancedClusterRSModel) *matlas.Cluster {
	if reflect.DeepEqual(plan.ReplicationSpecs, state.ReplicationSpecs) {
		return nil
	}

	currentSpecs := newReplicationSpecs(ctx, state.ReplicationSpecs)
	updatedSpecs := newReplicationSpecs(ctx, plan.ReplicationSpecs)

	if len(currentSpecs) != 1 || len(updatedSpecs) != 1 || len(currentSpecs[0].RegionConfigs) != 1 || len(updatedSpecs[0].RegionConfigs) != 1 {
		return nil
	}

	currentRegion := currentSpecs[0].RegionConfigs[0]
	updatedRegion := updatedSpecs[0].RegionConfigs[0]
	currentSize := currentRegion.ElectableSpecs.InstanceSize

	if currentRegion.ElectableSpecs.InstanceSize == updatedRegion.ElectableSpecs.InstanceSize || !IsSharedTier(currentSize) {
		return nil
	}

	return &matlas.Cluster{
		ProviderSettings: &matlas.ProviderSettings{
			ProviderName:     updatedRegion.ProviderName,
			InstanceSizeName: updatedRegion.ElectableSpecs.InstanceSize,
			RegionName:       updatedRegion.RegionName,
		},
	}
}

func updateCluster(ctx context.Context, conn *matlas.Client, state, plan *tfAdvancedClusterRSModel, timeout time.Duration) diag.Diagnostics {
	var diags diag.Diagnostics

	ids := conversion.DecodeStateID(state.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	cluster := new(matlas.AdvancedCluster)
	clusterChangeDetect := new(matlas.AdvancedCluster)

	if !plan.BackupEnabled.Equal(state.BackupEnabled) {
		cluster.BackupEnabled = plan.BackupEnabled.ValueBoolPointer()
	}

	if !reflect.DeepEqual(plan.BiConnectorConfig, state.BiConnectorConfig) {
		cluster.BiConnector = newBiConnectorConfig(ctx, plan.BiConnectorConfig)
	}

	if !plan.ClusterType.Equal(state.ClusterType) {
		cluster.ClusterType = plan.ClusterType.ValueString()
	}
	if !plan.BackupEnabled.Equal(state.BackupEnabled) {
		cluster.BackupEnabled = plan.BackupEnabled.ValueBoolPointer()
	}
	if !plan.DiskSizeGb.Equal(state.DiskSizeGb) {
		cluster.DiskSizeGB = plan.DiskSizeGb.ValueFloat64Pointer()
	}
	if !plan.EncryptionAtRestProvider.Equal(state.EncryptionAtRestProvider) {
		cluster.EncryptionAtRestProvider = plan.EncryptionAtRestProvider.ValueString()
	}

	if !reflect.DeepEqual(plan.Labels, state.Labels) {
		if ContainsLabelOrKey(newLabels(ctx, plan.Labels), defaultLabel) {
			diags.AddError("Unable to UPDATE cluster. An error occurred when updating labels.", "you should not set `Infrastructure Tool` label, it is used for internal purposes")
			return diags
		}
		cluster.Labels = newLabels(ctx, plan.Labels)
	}

	if !reflect.DeepEqual(plan.Tags, state.Tags) {
		cluster.Tags = newTags(ctx, plan.Tags)
	}

	if !plan.MongoDBMajorVersion.Equal(state.MongoDBMajorVersion) {
		cluster.MongoDBMajorVersion = utility.FormatMongoDBMajorVersion(plan.MongoDBMajorVersion.ValueString())
	}
	if !plan.PitEnabled.Equal(state.PitEnabled) {
		cluster.PitEnabled = plan.PitEnabled.ValueBoolPointer()
	}

	var tfRepSpecsPlan, tfRepSpecsState []tfReplicationSpecRSModel

	if !reflect.DeepEqual(plan.ReplicationSpecs, state.ReplicationSpecs) {
		// TODO remove:
		plan.ReplicationSpecs.ElementsAs(ctx, &tfRepSpecsPlan, true)
		state.ReplicationSpecs.ElementsAs(ctx, &tfRepSpecsState, true)

		if !reflect.DeepEqual(tfRepSpecsPlan, tfRepSpecsState) {
			if !reflect.DeepEqual(tfRepSpecsPlan[0].RegionsConfigs, tfRepSpecsState[0].RegionsConfigs) {
				var tfRegionConfigsPlan, tfRegionConfigsState []tfRegionsConfigModel
				tfRepSpecsPlan[0].RegionsConfigs.ElementsAs(ctx, &tfRegionConfigsPlan, true)
				tfRepSpecsState[0].RegionsConfigs.ElementsAs(ctx, &tfRegionConfigsState, true)
			}
		}
		cluster.ReplicationSpecs = newReplicationSpecs(ctx, plan.ReplicationSpecs)
	}

	if !plan.RootCertType.Equal(state.RootCertType) {
		cluster.RootCertType = plan.RootCertType.ValueString()
	}
	if !plan.TerminationProtectionEnabled.Equal(state.TerminationProtectionEnabled) {
		cluster.TerminationProtectionEnabled = plan.TerminationProtectionEnabled.ValueBoolPointer()
	}
	if !plan.AcceptDataRisksAndForceReplicaSetReconfig.Equal(state.AcceptDataRisksAndForceReplicaSetReconfig) {
		cluster.AcceptDataRisksAndForceReplicaSetReconfig = plan.AcceptDataRisksAndForceReplicaSetReconfig.ValueString()
	}
	if !plan.Paused.Equal(state.Paused) {
		cluster.Paused = plan.Paused.ValueBoolPointer()
	}

	if !reflect.DeepEqual(plan.AdvancedConfiguration, state.AdvancedConfiguration) {
		ac := plan.AdvancedConfiguration
		if len(ac.Elements()) > 0 {
			advancedConfReq := newAdvancedConfiguration(ctx, ac)
			if !reflect.DeepEqual(advancedConfReq, matlas.ProcessArgs{}) {
				_, _, err := conn.Clusters.UpdateProcessArgs(ctx, projectID, clusterName, advancedConfReq)
				if err != nil {
					diags.AddError("Unable to UPDATE cluster. An error occurred when updating advanced_configuration.", err.Error())
					return diags
				}
			}
		}
	}

	// Has changes
	if !reflect.DeepEqual(cluster, clusterChangeDetect) {
		err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, resp, err := updateAdvancedCluster(ctx, conn, cluster, projectID, clusterName, timeout)
			if err != nil {
				if resp == nil || resp.StatusCode == 400 {
					return retry.NonRetryableError(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
				}
				return retry.RetryableError(fmt.Errorf(errorClusterAdvancedUpdate, clusterName, err))
			}
			return nil
		})
		if err != nil {
			diags.AddError("Unable to UPDATE cluster. An error occurred when updating cluster in Atlas.", err.Error())
			return diags
		}
	}

	if plan.Paused.ValueBool() {
		clusterRequest := &matlas.AdvancedCluster{
			Paused: pointy.Bool(true),
		}

		_, _, err := updateAdvancedCluster(ctx, conn, clusterRequest, projectID, clusterName, timeout)
		if err != nil {
			diags.AddError("Unable to UPDATE cluster. An error occurred when attempting to pause cluster in Atlas.", err.Error())
			return diags
		}
	}

	return diags
}

func (r *advancedClusterRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	conn := r.Client.Atlas
	var state tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ids := conversion.DecodeStateID(state.ID.ValueString())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	var options *matlas.DeleteAdvanceClusterOptions
	if v := state.RetainBackupsEnabled; !v.IsNull() {
		options = &matlas.DeleteAdvanceClusterOptions{
			RetainBackups: v.ValueBoolPointer(),
		}
	}

	_, err := conn.AdvancedClusters.Delete(ctx, projectID, clusterName, options)
	if err != nil {
		resp.Diagnostics.AddError("Unable to DELETE cluster. An error occurred when deleting cluster in Atlas", fmt.Sprintf(errorClusterAdvancedDelete, clusterName, err))
		return
	}

	timeout, diags := state.Timeouts.Delete(ctx, defaultTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	log.Println("[INFO] Waiting for MongoDB ClusterAdvanced to be destroyed")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"IDLE", "CREATING", "UPDATING", "REPAIRING", "DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, clusterName, projectID, conn),
		Timeout:    timeout,
		MinTimeout: 30 * time.Second,
		Delay:      1 * time.Minute, // Wait 30 secs before starting
	}

	// Wait, catching any errors
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to DELETE cluster. An error occurred when deleting cluster in Atlas", fmt.Sprintf(errorClusterAdvancedDelete, clusterName, err))
		return
	}
}

func (r *advancedClusterRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	conn := r.Client.Atlas

	projectID, name, err := splitSClusterAdvancedImportID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to IMPORT cluster. An error occurred when attempting to read resource ID", err.Error())
		return
	}

	u, _, err := conn.AdvancedClusters.Get(ctx, *projectID, *name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to IMPORT cluster. An error occurred when getting cluster details from Atlas.",
			fmt.Sprintf("couldn't import cluster %s in project %s, error: %s", *name, *projectID, err))
		return
	}
	id := conversion.EncodeStateID(map[string]string{
		"cluster_id":   u.ID,
		"project_id":   u.GroupID,
		"cluster_name": u.Name,
	})

	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(id))
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTfReplicationSpecsRS(ctx context.Context, conn *matlas.Client,
	rawAPIObjects []*matlas.AdvancedReplicationSpec,
	configSpecsList types.List,
	projectID string) ([]tfReplicationSpecRSModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var configSpecs []tfReplicationSpecRSModel

	if !configSpecsList.IsNull() { // create return to state - filter by config, read/tf plan - filter by config, update - filter by config, import - return everything from API
		configSpecsList.ElementsAs(ctx, &configSpecs, true)
	}

	var apiObjects []*matlas.AdvancedReplicationSpec

	for _, advancedReplicationSpec := range rawAPIObjects {
		if advancedReplicationSpec != nil {
			apiObjects = append(apiObjects, advancedReplicationSpec)
		}
	}

	if len(apiObjects) == 0 {
		return nil, diags
	}

	tfList := make([]tfReplicationSpecRSModel, len(apiObjects))
	wasAPIObjectUsed := make([]bool, len(apiObjects))

	for i := 0; i < len(tfList); i++ {
		var tfMapObject tfReplicationSpecRSModel

		if len(configSpecs) > i {
			tfMapObject = configSpecs[i]
		}

		for j := 0; j < len(apiObjects); j++ {
			if wasAPIObjectUsed[j] {
				continue
			}

			if !TPFdoesAdvancedReplicationSpecMatchAPI(&tfMapObject, apiObjects[j]) {
				continue
			}

			advancedReplicationSpec, diags := newTfReplicationSpecRS(ctx, apiObjects[j], &tfMapObject, conn, projectID)
			if diags.HasError() {
				return nil, diags
			}

			tfList[i] = *advancedReplicationSpec
			wasAPIObjectUsed[j] = true
			break
		}
	}

	for i := range tfList {
		tfo := tfList[i]
		var tfMapObject *tfReplicationSpecRSModel
		if !reflect.DeepEqual(tfo, (tfReplicationSpecRSModel{})) {
			continue
		}

		if len(configSpecs) > i {
			tfMapObject = &configSpecs[i]
		}

		j := slices.IndexFunc(wasAPIObjectUsed, func(isUsed bool) bool { return !isUsed })
		advancedReplicationSpec, diags := newTfReplicationSpecRS(ctx, apiObjects[j], tfMapObject, conn, projectID)

		if diags.HasError() {
			return nil, diags
		}

		tfList[i] = *advancedReplicationSpec
		wasAPIObjectUsed[j] = true
	}

	return tfList, nil
}

func newTfReplicationSpecRS(ctx context.Context, apiObject *matlas.AdvancedReplicationSpec, configSpec *tfReplicationSpecRSModel,
	conn *matlas.Client, projectID string) (*tfReplicationSpecRSModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return nil, diags
	}

	tfMap := tfReplicationSpecRSModel{}
	tfMap.NumShards = types.Int64Value(cast.ToInt64(apiObject.NumShards))
	tfMap.ID = types.StringValue(apiObject.ID)
	if configSpec != nil {
		object, containerIds, diags := newTfRegionsConfigs(ctx, apiObject.RegionConfigs, configSpec.RegionsConfigs, conn, projectID)
		if diags.HasError() {
			return nil, diags
		}
		l, diags := types.ListValueFrom(ctx, tfRegionsConfigType, object)
		if diags.HasError() {
			return nil, diags
		}
		tfMap.RegionsConfigs = l
		tfMap.ContainerID = containerIds
	} else {
		object, containerIds, diags := newTfRegionsConfigs(ctx, apiObject.RegionConfigs, types.ListNull(tfRegionsConfigType), conn, projectID)
		if diags.HasError() {
			return nil, diags
		}
		l, diags := types.ListValueFrom(ctx, tfRegionsConfigType, object)
		if diags.HasError() {
			return nil, diags
		}
		tfMap.RegionsConfigs = l
		tfMap.ContainerID = containerIds
	}
	tfMap.ZoneName = types.StringValue(apiObject.ZoneName)

	return &tfMap, diags
}

func newTfRegionsConfigs(ctx context.Context, apiObjects []*matlas.AdvancedRegionConfig, configRegionConfigsList types.List,
	conn *matlas.Client, projectID string) (tfResult []tfRegionsConfigModel, containersIDs types.Map, diags1 diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(apiObjects) == 0 {
		return nil, types.MapNull(types.StringType), diags
	}

	var configRegionConfigs []*tfRegionsConfigModel
	containerIDsMap := map[string]attr.Value{}

	if !configRegionConfigsList.IsNull() { // create return to state - filter by config, read/tf plan - filter by config, update - filter by config, import - return everything from API
		configRegionConfigsList.ElementsAs(ctx, &configRegionConfigs, true)
	}

	var tfList []tfRegionsConfigModel

	for i, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		if len(configRegionConfigs) > i {
			tfMapObject := configRegionConfigs[i]
			rc, diags := newTfRegionsConfig(ctx, apiObject, tfMapObject)
			if diags.HasError() {
				break
			}

			tfList = append(tfList, *rc)
		} else {
			rc, diags := newTfRegionsConfig(ctx, apiObject, nil)
			if diags.HasError() {
				break
			}

			tfList = append(tfList, *rc)
		}

		if apiObject.ProviderName != "TENANT" {
			containers, _, err := conn.Containers.List(ctx, projectID,
				&matlas.ContainersListOptions{ProviderName: apiObject.ProviderName})
			if err != nil {
				diags.AddError("error when getting containers list from Atlas", err.Error())
				return nil, types.MapNull(types.StringType), diags
			}
			if result := getAdvancedClusterContainerID(containers, apiObject); result != "" {
				// Will print as "providerName:regionName" = "containerId" in terraform show
				key := fmt.Sprintf("%s:%s", apiObject.ProviderName, apiObject.RegionName)
				containerIDsMap[key] = types.StringValue(result)
			}
		}
	}
	tfContainersIDsMap, _ := types.MapValue(types.StringType, containerIDsMap)

	return tfList, tfContainersIDsMap, diags
}

func newTfRegionsConfig(ctx context.Context, apiObject *matlas.AdvancedRegionConfig, configRegionConfig *tfRegionsConfigModel) (*tfRegionsConfigModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	if apiObject == nil {
		return nil, diags
	}

	tfMap := tfRegionsConfigModel{}
	if configRegionConfig != nil {
		if v := configRegionConfig.AnalyticsSpecs; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.AnalyticsSpecs, d = newTfRegionsConfigSpec(ctx, apiObject.AnalyticsSpecs, apiObject.ProviderName, configRegionConfig.AnalyticsSpecs)
		} else {
			tfMap.AnalyticsSpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, []tfRegionsConfigSpecsModel{})
		}
		diags.Append(d...)
		if v := configRegionConfig.ElectableSpecs; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.ElectableSpecs, d = newTfRegionsConfigSpec(ctx, apiObject.ElectableSpecs, apiObject.ProviderName, configRegionConfig.ElectableSpecs)
		} else {
			tfMap.ElectableSpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, []tfRegionsConfigSpecsModel{})
		}
		diags.Append(d...)
		if v := configRegionConfig.ReadOnlySpecs; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.ReadOnlySpecs, d = newTfRegionsConfigSpec(ctx, apiObject.ReadOnlySpecs, apiObject.ProviderName, configRegionConfig.ReadOnlySpecs)
		} else {
			tfMap.ReadOnlySpecs, d = types.ListValueFrom(ctx, tfRegionsConfigSpecType, []tfRegionsConfigSpecsModel{})
		}
		diags.Append(d...)
		if v := configRegionConfig.AutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.AutoScaling, d = newTfRegionsConfigAutoScalingSpecs(ctx, apiObject.AutoScaling)
		} else {
			tfMap.AutoScaling, d = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, []tfRegionsConfigAutoScalingSpecsModel{})
		}
		diags.Append(d...)
		if v := configRegionConfig.AnalyticsAutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
			tfMap.AnalyticsAutoScaling, d = newTfRegionsConfigAutoScalingSpecs(ctx, apiObject.AnalyticsAutoScaling)
		} else {
			tfMap.AnalyticsAutoScaling, d = types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, []tfRegionsConfigAutoScalingSpecsModel{})
		}
		diags.Append(d...)
	} else {
		nilSpecList := types.ListNull(tfRegionsConfigSpecType)
		tfMap.AnalyticsSpecs, d = newTfRegionsConfigSpec(ctx, apiObject.AnalyticsSpecs, apiObject.ProviderName, nilSpecList)
		diags.Append(d...)
		tfMap.ElectableSpecs, d = newTfRegionsConfigSpec(ctx, apiObject.ElectableSpecs, apiObject.ProviderName, nilSpecList)
		diags.Append(d...)
		tfMap.ReadOnlySpecs, d = newTfRegionsConfigSpec(ctx, apiObject.ReadOnlySpecs, apiObject.ProviderName, nilSpecList)
		diags.Append(d...)
		tfMap.AutoScaling, d = newTfRegionsConfigAutoScalingSpecs(ctx, apiObject.AutoScaling)
		diags.Append(d...)
		tfMap.AnalyticsAutoScaling, d = newTfRegionsConfigAutoScalingSpecs(ctx, apiObject.AnalyticsAutoScaling)
		diags.Append(d...)
	}

	tfMap.RegionName = types.StringValue(apiObject.RegionName)
	tfMap.ProviderName = types.StringValue(apiObject.ProviderName)
	tfMap.BackingProviderName = types.StringValue(apiObject.BackingProviderName)
	tfMap.Priority = types.Int64Value(cast.ToInt64(apiObject.Priority))

	return &tfMap, diags
}

func newTfRegionsConfigSpec(ctx context.Context, apiObject *matlas.Specs, providerName string, tfMapObjects types.List) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ListNull(tfRegionsConfigSpecType), diags
	}

	var configRegionConfigSpecs []*tfRegionsConfigSpecsModel

	if !tfMapObjects.IsNull() { // create return to state - filter by config, read/tf plan - filter by config, update - filter by config, import - return everything from API
		tfMapObjects.ElementsAs(ctx, &configRegionConfigSpecs, true)
	}

	var tfList []tfRegionsConfigSpecsModel

	tfMap := tfRegionsConfigSpecsModel{}

	if len(configRegionConfigSpecs) > 0 {
		tfMapObject := configRegionConfigSpecs[0]

		if providerName == "AWS" {
			if cast.ToInt64(apiObject.DiskIOPS) > 0 {
				tfMap.DiskIOPS = types.Int64PointerValue(apiObject.DiskIOPS)
			} else {
				tfMap.DiskIOPS = types.Int64Null()
			}
			// if v := tfMapObject.EBSVolumeType; !v.IsNull() && v.ValueString() != "" {
			// 	tfMap.EBSVolumeType = types.StringValue(apiObject.EbsVolumeType)
			// }
			if v := tfMapObject.EBSVolumeType; !v.IsNull() {
				tfMap.EBSVolumeType = types.StringValue(apiObject.EbsVolumeType)
			}

		}
		if v := tfMapObject.NodeCount; !v.IsNull() {
			tfMap.NodeCount = types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiObject.NodeCount))
		}
		if v := tfMapObject.InstanceSize; !v.IsNull() && v.ValueString() != "" {
			tfMap.InstanceSize = types.StringValue(apiObject.InstanceSize)
		}

		// if tfMap.DiskIOPS.IsNull() {
		// 	tfMap.DiskIOPS = types.Int64Value(defaultInt)
		// }
		// if tfMap.NodeCount.IsNull() {
		// 	tfMap.NodeCount = types.Int64Value(defaultInt)
		// }
		// if tfMap.EBSVolumeType.IsNull() {
		// 	tfMap.EBSVolumeType = types.StringValue(defaultString)
		// }
		tfList = append(tfList, tfMap)
	} else {
		tfMap.DiskIOPS = types.Int64PointerValue(apiObject.DiskIOPS)
		tfMap.EBSVolumeType = types.StringValue(apiObject.EbsVolumeType)
		tfMap.NodeCount = types.Int64PointerValue(conversion.IntPtrToInt64Ptr(apiObject.NodeCount))
		tfMap.InstanceSize = types.StringValue(apiObject.InstanceSize)
		// if tfMap.DiskIOPS.IsNull() {
		// 	tfMap.DiskIOPS = types.Int64Value(defaultInt)
		// }
		// if tfMap.NodeCount.IsNull() {
		// 	tfMap.NodeCount = types.Int64Value(defaultInt)
		// }
		// if tfMap.EBSVolumeType.IsNull() {
		// 	tfMap.EBSVolumeType = types.StringValue(defaultString)
		// }
		tfList = append(tfList, tfMap)
	}

	return types.ListValueFrom(ctx, tfRegionsConfigSpecType, tfList)
}

func newTfRegionsConfigAutoScalingSpecs(ctx context.Context, apiObject *matlas.AdvancedAutoScaling) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if apiObject == nil {
		return types.ListNull(tfRegionsConfigAutoScalingSpecType), diags
	}

	var tfList []tfRegionsConfigAutoScalingSpecsModel

	tfMap := tfRegionsConfigAutoScalingSpecsModel{}
	if apiObject.DiskGB != nil {
		tfMap.DiskGBEnabled = types.BoolPointerValue(apiObject.DiskGB.Enabled)
	}
	if apiObject.Compute != nil {
		tfMap.ComputeEnabled = types.BoolPointerValue(apiObject.Compute.Enabled)
		tfMap.ComputeScaleDownEnabled = types.BoolPointerValue(apiObject.Compute.ScaleDownEnabled)
		tfMap.ComputeMinInstanceSize = types.StringValue(apiObject.Compute.MinInstanceSize)
		tfMap.ComputeMaxInstanceSize = types.StringValue(apiObject.Compute.MaxInstanceSize)
	}

	tfList = append(tfList, tfMap)

	return types.ListValueFrom(ctx, tfRegionsConfigAutoScalingSpecType, tfList)
}

func TPFdoesAdvancedReplicationSpecMatchAPI(tfObject *tfReplicationSpecRSModel, apiObject *matlas.AdvancedReplicationSpec) bool {
	return tfObject.ID.ValueString() == apiObject.ID || (tfObject.ID.IsNull() && tfObject.ZoneName.ValueString() == apiObject.ZoneName)
}

func newAdvancedConfiguration(ctx context.Context, tfList basetypes.ListValue) *matlas.ProcessArgs {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var tfAdvancedConfigArr []TfAdvancedConfigurationModel
	tfList.ElementsAs(ctx, &tfAdvancedConfigArr, true)

	if len(tfAdvancedConfigArr) == 0 {
		return nil
	}
	tfModel := tfAdvancedConfigArr[0]

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
			log.Printf(ErrorClusterSetting, `oplog_size_mb`, "", cast.ToString(sizeMB))
		}
	}

	if v := tfModel.OplogMinRetentionHours; !v.IsNull() {
		if minRetentionHours := v.ValueInt64(); minRetentionHours >= 0 {
			res.OplogMinRetentionHours = pointy.Float64(cast.ToFloat64(v.ValueInt64()))
		} else {
			log.Printf(ErrorClusterSetting, `oplog_min_retention_hours`, "", cast.ToString(minRetentionHours))
		}
	}

	if v := tfModel.TransactionLifetimeLimitSeconds; !v.IsUnknown() {
		if transactionLimitSeconds := v.ValueInt64(); transactionLimitSeconds > 0 {
			res.TransactionLifetimeLimitSeconds = v.ValueInt64Pointer()
		} else {
			log.Printf(ErrorClusterSetting, `transaction_lifetime_limit_seconds`, "", cast.ToString(transactionLimitSeconds))
		}
	}

	return res
}

func newTags(ctx context.Context, tfSet basetypes.SetValue) []*matlas.Tag {
	if tfSet.IsNull() || len(tfSet.Elements()) == 0 {
		return nil
	}
	var tfArr []TfTagModel
	tfSet.ElementsAs(ctx, &tfArr, true)

	res := make([]*matlas.Tag, len(tfArr))
	for i, v := range tfArr {
		res[i] = &matlas.Tag{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}
	return res
}

func newLabels(ctx context.Context, tfSet basetypes.SetValue) []matlas.Label {
	if tfSet.IsNull() || len(tfSet.Elements()) == 0 {
		return nil
	}

	var tfArr []TfLabelModel
	tfSet.ElementsAs(ctx, &tfArr, true)

	res := make([]matlas.Label, len(tfArr))

	for i, v := range tfArr {
		res[i] = matlas.Label{
			Key:   v.Key.ValueString(),
			Value: v.Value.ValueString(),
		}
	}

	return res
}

func newBiConnectorConfig(ctx context.Context, tfList basetypes.ListValue) *matlas.BiConnector {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var tfArr []TfBiConnectorConfigModel
	tfList.ElementsAs(ctx, &tfArr, true)

	tfBiConnector := tfArr[0]

	biConnector := matlas.BiConnector{
		Enabled:        tfBiConnector.Enabled.ValueBoolPointer(),
		ReadPreference: tfBiConnector.ReadPreference.ValueString(),
	}

	return &biConnector
}

func newReplicationSpecs(ctx context.Context, tfList basetypes.ListValue) []*matlas.AdvancedReplicationSpec {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var tfRepSpecs []tfReplicationSpecRSModel
	tfList.ElementsAs(ctx, &tfRepSpecs, true)

	var repSpecs []*matlas.AdvancedReplicationSpec

	for i := range tfRepSpecs {
		rs := newReplicationSpec(ctx, &tfRepSpecs[i])
		repSpecs = append(repSpecs, rs)
	}
	return repSpecs
}

func newReplicationSpec(ctx context.Context, tfRepSpec *tfReplicationSpecRSModel) *matlas.AdvancedReplicationSpec {
	if tfRepSpec == nil {
		return nil
	}

	zoneName := tfRepSpec.ZoneName.ValueString()
	if !conversion.IsStringPresent(&zoneName) {
		zoneName = defaultZoneName
	}
	return &matlas.AdvancedReplicationSpec{
		NumShards:     int(tfRepSpec.NumShards.ValueInt64()),
		ZoneName:      tfRepSpec.ZoneName.ValueString(),
		RegionConfigs: newRegionConfigs(ctx, tfRepSpec.RegionsConfigs),
	}
}

func newRegionConfigs(ctx context.Context, tfList basetypes.ListValue) []*matlas.AdvancedRegionConfig {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var tfRegionConfigs []tfRegionsConfigModel
	tfList.ElementsAs(ctx, &tfRegionConfigs, true)

	var regionConfigs []*matlas.AdvancedRegionConfig

	for i := range tfRegionConfigs {
		rc := newRegionConfig(ctx, &tfRegionConfigs[i])

		regionConfigs = append(regionConfigs, rc)
	}

	return regionConfigs
}

func newRegionConfig(ctx context.Context, tfRegionConfig *tfRegionsConfigModel) *matlas.AdvancedRegionConfig {
	if tfRegionConfig == nil {
		return nil
	}

	providerName := tfRegionConfig.ProviderName.ValueString()
	apiObject := &matlas.AdvancedRegionConfig{
		Priority:     conversion.Int64PtrToIntPtr(tfRegionConfig.Priority.ValueInt64Pointer()),
		ProviderName: providerName,
		RegionName:   tfRegionConfig.RegionName.ValueString(),
	}

	if v := tfRegionConfig.AnalyticsSpecs; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.AnalyticsSpecs = newRegionConfigSpec(ctx, v, providerName)
	}
	if v := tfRegionConfig.ElectableSpecs; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.ElectableSpecs = newRegionConfigSpec(ctx, v, providerName)
	}
	if v := tfRegionConfig.ReadOnlySpecs; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.ReadOnlySpecs = newRegionConfigSpec(ctx, v, providerName)
	}
	if v := tfRegionConfig.AutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.AutoScaling = newRegionConfigAutoScalingSpec(ctx, v)
	}
	if v := tfRegionConfig.AnalyticsAutoScaling; !v.IsNull() && len(v.Elements()) > 0 {
		apiObject.AnalyticsAutoScaling = newRegionConfigAutoScalingSpec(ctx, v)
	}
	if v := tfRegionConfig.BackingProviderName; !v.IsNull() && v.ValueString() != defaultString {
		apiObject.BackingProviderName = v.ValueString()
	}

	return apiObject
}

func newRegionConfigAutoScalingSpec(ctx context.Context, tfList basetypes.ListValue) *matlas.AdvancedAutoScaling {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var specs []tfRegionsConfigAutoScalingSpecsModel
	tfList.ElementsAs(ctx, &specs, true)

	spec := specs[0]
	advancedAutoScaling := &matlas.AdvancedAutoScaling{}
	diskGB := &matlas.DiskGB{}
	compute := &matlas.Compute{}

	if v := spec.DiskGBEnabled; !v.IsUnknown() {
		diskGB.Enabled = v.ValueBoolPointer()
	}
	if v := spec.ComputeEnabled; !v.IsUnknown() {
		compute.Enabled = v.ValueBoolPointer()
	}
	if v := spec.ComputeScaleDownEnabled; !v.IsUnknown() {
		compute.ScaleDownEnabled = v.ValueBoolPointer()
	}
	if v := spec.ComputeMinInstanceSize; !v.IsUnknown() {
		value := compute.ScaleDownEnabled
		if *value {
			compute.MinInstanceSize = v.ValueString()
		}
	}
	if v := spec.ComputeMaxInstanceSize; !v.IsUnknown() {
		value := compute.Enabled
		if *value {
			compute.MaxInstanceSize = v.ValueString()
		}
	}

	advancedAutoScaling.DiskGB = diskGB
	advancedAutoScaling.Compute = compute

	return advancedAutoScaling
}

func newRegionConfigSpec(ctx context.Context, tfList basetypes.ListValue, providerName string) *matlas.Specs {
	if tfList.IsNull() || len(tfList.Elements()) == 0 {
		return nil
	}

	var specs []tfRegionsConfigSpecsModel
	tfList.ElementsAs(ctx, &specs, true)

	spec := specs[0]
	apiObject := &matlas.Specs{}

	if providerName == "AWS" {
		if v := spec.DiskIOPS; !v.IsNull() && v.ValueInt64() > 0 {
			apiObject.DiskIOPS = v.ValueInt64Pointer()
		}
		if v := spec.EBSVolumeType; !v.IsNull() && v.ValueString() != defaultString {
			apiObject.EbsVolumeType = v.ValueString()
		}
	}

	if v := spec.InstanceSize; !v.IsNull() {
		apiObject.InstanceSize = v.ValueString()
	}
	if v := spec.NodeCount; !v.IsNull() && v.ValueInt64() > 0 {
		apiObject.NodeCount = conversion.Int64PtrToIntPtr(v.ValueInt64Pointer())
	}
	return apiObject
}

func updateAdvancedCluster(ctx context.Context, conn *matlas.Client, request *matlas.AdvancedCluster, projectID, name string, timeout time.Duration,
) (*matlas.AdvancedCluster, *matlas.Response, error) {
	cluster, resp, err := conn.AdvancedClusters.Update(ctx, projectID, name, request)
	if err != nil {
		return nil, nil, err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"CREATING", "UPDATING", "REPAIRING"},
		Target:     []string{"IDLE"},
		Refresh:    resourceClusterAdvancedRefreshFunc(ctx, name, projectID, conn),
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

type tfAdvancedClusterRSModel struct {
	DiskSizeGb                                types.Float64                    `tfsdk:"disk_size_gb"`
	Labels                                    types.Set                        `tfsdk:"labels"`
	AdvancedConfiguration                     types.List                       `tfsdk:"advanced_configuration"`
	ConnectionStrings                         types.List                       `tfsdk:"connection_strings"`
	BiConnectorConfig                         types.List                       `tfsdk:"bi_connector_config"`
	ReplicationSpecs                          types.List                       `tfsdk:"replication_specs"`
	Tags                                      types.Set                        `tfsdk:"tags"`
	ProjectID                                 types.String                     `tfsdk:"project_id"`
	RootCertType                              types.String                     `tfsdk:"root_cert_type"`
	Name                                      types.String                     `tfsdk:"name"`
	Timeouts                                  timeouts.Value                   `tfsdk:"timeouts"`
	ClusterID                                 types.String                     `tfsdk:"cluster_id"`
	MongoDBVersion                            types.String                     `tfsdk:"mongo_db_version"`
	ClusterType                               types.String                     `tfsdk:"cluster_type"`
	EncryptionAtRestProvider                  types.String                     `tfsdk:"encryption_at_rest_provider"`
	StateName                                 types.String                     `tfsdk:"state_name"`
	CreateDate                                types.String                     `tfsdk:"create_date"`
	VersionReleaseSystem                      types.String                     `tfsdk:"version_release_system"`
	AcceptDataRisksAndForceReplicaSetReconfig types.String                     `tfsdk:"accept_data_risks_and_force_replica_set_reconfig"`
	MongoDBMajorVersion                       customtypes.DBVersionStringValue `tfsdk:"mongo_db_major_version"`
	ID                                        types.String                     `tfsdk:"id"`
	BackupEnabled                             types.Bool                       `tfsdk:"backup_enabled"`
	TerminationProtectionEnabled              types.Bool                       `tfsdk:"termination_protection_enabled"`
	RetainBackupsEnabled                      types.Bool                       `tfsdk:"retain_backups_enabled"`
	PitEnabled                                types.Bool                       `tfsdk:"pit_enabled"`
	Paused                                    types.Bool                       `tfsdk:"paused"`
}

type tfReplicationSpecRSModel struct {
	RegionsConfigs types.List   `tfsdk:"region_configs"`
	ContainerID    types.Map    `tfsdk:"container_id"`
	ID             types.String `tfsdk:"id"`
	ZoneName       types.String `tfsdk:"zone_name"`
	NumShards      types.Int64  `tfsdk:"num_shards"`
}

var tfReplicationSpecRSType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":             types.StringType,
	"zone_name":      types.StringType,
	"num_shards":     types.Int64Type,
	"container_id":   types.MapType{ElemType: types.StringType},
	"region_configs": types.ListType{ElemType: tfRegionsConfigType},
},
}
