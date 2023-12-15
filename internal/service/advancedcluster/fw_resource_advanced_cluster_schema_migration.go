package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/planmodifiers"
)

func fw_ResourceV0() schema.Schema {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			// "id": schema.StringAttribute{
			// 	Computed: true,
			// 	PlanModifiers: []planmodifier.String{
			// 		stringplanmodifier.UseStateForUnknown(),
			// 	},
			// },
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
			"connection_strings": advClusterRSConnectionStringSchemaAttr(),
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
				Default:  booldefault.StaticBool(false),
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
				Default:  stringdefault.StaticString("LTS"),
				Validators: []validator.String{
					stringvalidator.OneOf("LTS", "CONTINUOUS"),
				},
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
		},
		Blocks: map[string]schema.Block{
			"advanced_configuration": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"default_read_concern": schema.StringAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.String{
							// 	planmodifiers.UseNullForUnknownString(),
							// 	// stringplanmodifier.UseStateForUnknown(),
							// },
						},
						"default_write_concern": schema.StringAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.String{
							// 	planmodifiers.UseNullForUnknownString(),
							// 	// stringplanmodifier.UseStateForUnknown(),
							// },
						},
						"fail_index_key_too_long": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Bool{
							// 	planmodifiers.UseNullForUnknownBool(),
							// 	// boolplanmodifier.UseStateForUnknown(),
							// },
						},
						"javascript_enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Bool{
							// 	planmodifiers.UseNullForUnknownBool(),
							// 	// boolplanmodifier.UseStateForUnknown(),
							// },
						},
						"minimum_enabled_tls_protocol": schema.StringAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.String{
							// 	planmodifiers.UseNullForUnknownString(),
							// 	// stringplanmodifier.UseStateForUnknown(),
							// },
						},
						"no_table_scan": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Bool{
							// 	planmodifiers.UseNullForUnknownBool(),
							// 	boolplanmodifier.UseStateForUnknown(),
							// },
						},
						"oplog_min_retention_hours": schema.Int64Attribute{
							Optional: true,
						},
						"oplog_size_mb": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Int64{
							// 	planmodifiers.UseNullForUnknownInt64(),
							// 	// int64planmodifier.UseStateForUnknown(),
							// },
						},
						"sample_refresh_interval_bi_connector": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Int64{
							// 	planmodifiers.UseNullForUnknownInt64(),
							// 	// int64planmodifier.UseStateForUnknown(),
							// },
						},
						"sample_size_bi_connector": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Int64{
							// 	planmodifiers.UseNullForUnknownInt64(),
							// 	// int64planmodifier.UseStateForUnknown(),
							// },
						},
						"transaction_lifetime_limit_seconds": schema.Int64Attribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Int64{
							// 	planmodifiers.UseNullForUnknownInt64(),
							// 	// int64planmodifier.UseStateForUnknown(),
							// },
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				// PlanModifiers: []planmodifier.List{
				// 	// planmodifiers.UseNullForUnknownInt64(),
				// 	listplanmodifier.UseStateForUnknown(),
				// },
			},
			"bi_connector": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.Bool{
							// 	// planmodifiers.UseNullForUnknownBool(),
							// 	boolplanmodifier.UseStateForUnknown(),
							// },
						},
						"read_preference": schema.StringAttribute{
							Optional: true,
							Computed: true,
							// PlanModifiers: []planmodifier.String{
							// 	// planmodifiers.UseNullForUnknownString(),
							// 	stringplanmodifier.UseStateForUnknown(),
							// },
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
				},
				// PlanModifiers: []planmodifier.List{
				// 	// planmodifiers.UseNullForUnknownInt64(),
				// 	listplanmodifier.UseStateForUnknown(),
				// },
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
			"replication_specs": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
							// PlanModifiers: []planmodifier.String{
							// 	// planmodifiers.UseNullForUnknownBool(),
							// 	stringplanmodifier.UseStateForUnknown(),
							// },
						},
						"container_id": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
							// PlanModifiers: []planmodifier.Map{
							// 	// planmodifiers.UseNullForUnknownBool(),
							// 	mapplanmodifier.UseStateForUnknown(),
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
						},
					},
				},
				Validators: []validator.Set{
					setvalidator.IsRequired(),
				},
			},
		},
		Version: 1,
	}

	if s.Blocks == nil {
		s.Blocks = make(map[string]schema.Block)
	}
	s.Blocks["timeouts"] = timeouts.Block(context.Background(), timeouts.Opts{
		Create: true,
		Update: true,
		Delete: true,
	})
	return s
}

func upgradeAdvClusterResourceStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var clusterV0 tfAdvancedClusterRSModelV0

	resp.Diagnostics.Append(req.State.Get(ctx, &clusterV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterV1 := tfAdvancedClusterRSModel{
		ID: clusterV0.ID,
	}

	clusterV1.BiConnectorConfig = clusterV0.BiConnector

	diags := resp.State.Set(ctx, clusterV1)
	resp.Diagnostics.Append(diags...)
}

type tfAdvancedClusterRSModelV0 struct {
	BackupEnabled            types.Bool    `tfsdk:"backup_enabled"`
	ClusterID                types.String  `tfsdk:"cluster_id"`
	ClusterType              types.String  `tfsdk:"cluster_type"`
	CreateDate               types.String  `tfsdk:"create_date"`
	DiskSizeGb               types.Float64 `tfsdk:"disk_size_gb"`
	EncryptionAtRestProvider types.String  `tfsdk:"encryption_at_rest_provider"`
	ID                       types.String  `tfsdk:"id"`
	// MongoDBMajorVersion                       types.String  `tfsdk:"mongo_db_major_version"`
	MongoDBMajorVersion          customtypes.DBVersionStringValue `tfsdk:"mongo_db_major_version"`
	MongoDBVersion               types.String                     `tfsdk:"mongo_db_version"`
	Name                         types.String                     `tfsdk:"name"`
	Paused                       types.Bool                       `tfsdk:"paused"`
	PitEnabled                   types.Bool                       `tfsdk:"pit_enabled"`
	ProjectID                    types.String                     `tfsdk:"project_id"`
	RetainBackupsEnabled         types.Bool                       `tfsdk:"retain_backups_enabled"`
	RootCertType                 types.String                     `tfsdk:"root_cert_type"`
	StateName                    types.String                     `tfsdk:"state_name"`
	TerminationProtectionEnabled types.Bool                       `tfsdk:"termination_protection_enabled"`
	VersionReleaseSystem         types.String                     `tfsdk:"version_release_system"`

	Labels                types.Set  `tfsdk:"labels"`
	ReplicationSpecs      types.Set  `tfsdk:"replication_specs"`
	BiConnector           types.List `tfsdk:"bi_connector"`
	ConnectionStrings     types.List `tfsdk:"connection_strings"`
	AdvancedConfiguration types.List `tfsdk:"advanced_configuration"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
