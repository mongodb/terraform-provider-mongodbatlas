package advancedcluster

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customtypes"
)

func TPFResourceV0() schema.Schema {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
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
			"connection_strings": advClusterRSConnectionStringSchemaComputed(),
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
			"mongo_db_major_version": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"mongo_db_version": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"paused": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
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
				Default:  stringdefault.StaticString("LTS"),
				Validators: []validator.String{
					stringvalidator.OneOf("LTS", "CONTINUOUS"),
				},
			},
			"advanced_configuration": advClusterRSAdvancedConfigurationSchema(),
			"bi_connector_config":    advClusterRSBiConnectorConfigSchema(),
			"replication_specs":      advClusterRSReplicationSpecsSchemaV0(),
		},
		Blocks: map[string]schema.Block{
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
		},
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

func advClusterRSReplicationSpecsSchemaV0() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"container_id": schema.MapAttribute{
					ElementType: types.StringType,
					Computed:    true,
					PlanModifiers: []planmodifier.Map{
						mapplanmodifier.UseStateForUnknown(),
					},
				},
				"id": schema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
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
				"region_configs": schema.SetNestedAttribute{
					Optional: true,
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"backing_provider_name": schema.StringAttribute{
								Optional: true,
								Computed: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
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
					Validators: []validator.Set{
						setvalidator.IsRequired(),
					},
				},
			},
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
	}
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
	DiskSizeGb                   types.Float64                    `tfsdk:"disk_size_gb"`
	Labels                       types.Set                        `tfsdk:"labels"`
	AdvancedConfiguration        types.List                       `tfsdk:"advanced_configuration"`
	ConnectionStrings            types.List                       `tfsdk:"connection_strings"`
	BiConnector                  types.List                       `tfsdk:"bi_connector"`
	ReplicationSpecs             types.Set                        `tfsdk:"replication_specs"`
	ID                           types.String                     `tfsdk:"id"`
	EncryptionAtRestProvider     types.String                     `tfsdk:"encryption_at_rest_provider"`
	MongoDBVersion               types.String                     `tfsdk:"mongo_db_version"`
	Name                         types.String                     `tfsdk:"name"`
	Timeouts                     timeouts.Value                   `tfsdk:"timeouts"`
	ClusterID                    types.String                     `tfsdk:"cluster_id"`
	ProjectID                    types.String                     `tfsdk:"project_id"`
	ClusterType                  types.String                     `tfsdk:"cluster_type"`
	RootCertType                 types.String                     `tfsdk:"root_cert_type"`
	StateName                    types.String                     `tfsdk:"state_name"`
	CreateDate                   types.String                     `tfsdk:"create_date"`
	VersionReleaseSystem         types.String                     `tfsdk:"version_release_system"`
	MongoDBMajorVersion          customtypes.DBVersionStringValue `tfsdk:"mongo_db_major_version"`
	BackupEnabled                types.Bool                       `tfsdk:"backup_enabled"`
	TerminationProtectionEnabled types.Bool                       `tfsdk:"termination_protection_enabled"`
	RetainBackupsEnabled         types.Bool                       `tfsdk:"retain_backups_enabled"`
	PitEnabled                   types.Bool                       `tfsdk:"pit_enabled"`
	Paused                       types.Bool                       `tfsdk:"paused"`
}
