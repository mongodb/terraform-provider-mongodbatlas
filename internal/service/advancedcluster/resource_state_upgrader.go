package advancedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tfAdvancedClusterRSModelV0 struct {
	DiskSizeGb                   types.Float64                  `tfsdk:"disk_size_gb"`
	MongoDBVersion               types.String                   `tfsdk:"mongo_db_version"`
	ClusterType                  types.String                   `tfsdk:"cluster_type"`
	ProjectID                    types.String                   `tfsdk:"project_id"`
	Name                         types.String                   `tfsdk:"name"`
	ClusterID                    types.String                   `tfsdk:"cluster_id"`
	EncryptionAtRestProvider     types.String                   `tfsdk:"encryption_at_rest_provider"`
	MongoDBMajorVersion          types.String                   `tfsdk:"mongo_db_major_version"`
	Timeouts                     timeouts.Value                 `tfsdk:"timeouts"`
	ID                           types.String                   `tfsdk:"id"`
	VersionReleaseSystem         types.String                   `tfsdk:"version_release_system"`
	CreateDate                   types.String                   `tfsdk:"create_date"`
	RootCertType                 types.String                   `tfsdk:"root_cert_type"`
	StateName                    types.String                   `tfsdk:"state_name"`
	ReplicationSpecs             []tfReplicationSpecRSModelV0   `tfsdk:"replication_specs"`
	AdvancedConfiguration        []TfAdvancedConfigurationModel `tfsdk:"advanced_configuration"`
	ConnectionStrings            []tfConnectionStringModel      `tfsdk:"connection_strings"`
	Labels                       []TfLabelModel                 `tfsdk:"labels"`
	BiConnector                  []TfBiConnectorConfigModel     `tfsdk:"bi_connector"`
	PitEnabled                   types.Bool                     `tfsdk:"pit_enabled"`
	TerminationProtectionEnabled types.Bool                     `tfsdk:"termination_protection_enabled"`
	Paused                       types.Bool                     `tfsdk:"paused"`
	BackupEnabled                types.Bool                     `tfsdk:"backup_enabled"`
}

type tfReplicationSpecRSModelV0 struct {
	ContainerID    types.Map                `tfsdk:"container_id"`
	ID             types.String             `tfsdk:"id"`
	ZoneName       types.String             `tfsdk:"zone_name"`
	RegionsConfigs []tfRegionsConfigModelV0 `tfsdk:"region_configs"`
	NumShards      types.Int64              `tfsdk:"num_shards"`
}

type tfRegionsConfigModelV0 struct {
	AnalyticsSpecs      types.List   `tfsdk:"analytics_specs"`
	AutoScaling         types.List   `tfsdk:"auto_scaling"`
	ReadOnlySpecs       types.List   `tfsdk:"read_only_specs"`
	ElectableSpecs      types.List   `tfsdk:"electable_specs"`
	BackingProviderName types.String `tfsdk:"backing_provider_name"`
	ProviderName        types.String `tfsdk:"provider_name"`
	RegionName          types.String `tfsdk:"region_name"`
	Priority            types.Int64  `tfsdk:"priority"`
}

func TPFResourceV0(ctx context.Context) schema.Schema {
	s := schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
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
			"cluster_type": schema.StringAttribute{
				Required: true,
			},
			"connection_strings": advancedClusterRSConnectionStringSchemaComputed(), // checked
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
			"labels": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
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
			},
			"advanced_configuration": advancedClusterRSAdvancedConfigurationSchema(), // checked
			"bi_connector":           advancedClusterRSBiConnectorConfigSchema(),
			"replication_specs":      advancedClusterRSReplicationSpecsSchemaV0(),
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}

	return s
}

func advancedClusterRSReplicationSpecsSchemaV0() schema.SetNestedAttribute {
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
							"auto_scaling":    advancedClusterRSRegionConfigAutoScalingSpecsSchema(), // checked
							"analytics_specs": advancedClusterRSRegionConfigSpecsSchemaV0(),          // checked
							"electable_specs": advancedClusterRSRegionConfigSpecsSchemaV0(),          // checked
							"read_only_specs": advancedClusterRSRegionConfigSpecsSchemaV0(),          // checked
						},
					},
					Validators: []validator.Set{
						setvalidator.IsRequired(),
					},
				},
			},
		},
	}
}

func advancedClusterRSRegionConfigSpecsSchemaV0() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
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
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
	}
}

func (*advancedClusterRS) UpgradeState(ctx context.Context) map[int64]resource.StateUpgrader {
	schemaV0 := TPFResourceV0(ctx)

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &schemaV0,
			StateUpgrader: upgradeAdvancedClusterResourceStateV0toV1,
		},
	}
}

func upgradeAdvancedClusterResourceStateV0toV1(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	var clusterV0 tfAdvancedClusterRSModelV0
	var diags diag.Diagnostics

	resp.Diagnostics.Append(req.State.Get(ctx, &clusterV0)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterV1, d := newTfAdvancedClusterRSModelV1(ctx, &clusterV0)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	diags = resp.State.Set(ctx, clusterV1)
	resp.Diagnostics.Append(diags...)
}

func newTfAdvancedClusterRSModelV1(ctx context.Context, clusterV0 *tfAdvancedClusterRSModelV0) (*tfAdvancedClusterRSModel, diag.Diagnostics) {
	var d, diags diag.Diagnostics

	clusterV1 := tfAdvancedClusterRSModel{
		ID:                           clusterV0.ID,
		ClusterID:                    clusterV0.ClusterID,
		BackupEnabled:                clusterV0.BackupEnabled,
		ClusterType:                  clusterV0.ClusterType,
		CreateDate:                   clusterV0.CreateDate,
		DiskSizeGb:                   clusterV0.DiskSizeGb,
		EncryptionAtRestProvider:     clusterV0.EncryptionAtRestProvider,
		MongoDBMajorVersion:          clusterV0.MongoDBMajorVersion,
		MongoDBVersion:               clusterV0.MongoDBVersion,
		Name:                         clusterV0.Name,
		Paused:                       clusterV0.Paused,
		PitEnabled:                   clusterV0.PitEnabled,
		RootCertType:                 clusterV0.RootCertType,
		StateName:                    clusterV0.StateName,
		TerminationProtectionEnabled: clusterV0.TerminationProtectionEnabled,
		VersionReleaseSystem:         clusterV0.VersionReleaseSystem,
		ProjectID:                    clusterV0.ProjectID,
		Timeouts:                     clusterV0.Timeouts,
	}

	clusterV1.BiConnectorConfig, d = newTfBiConnectorConfigFromV0(ctx, clusterV0.BiConnector)
	diags.Append(d...)

	clusterV1.ConnectionStrings, d = newTfConnectionStringsFromV0(ctx, clusterV0.ConnectionStrings)
	diags.Append(d...)

	clusterV1.Labels, d = newTfLabelsFromV0(ctx, clusterV0.Labels)
	diags.Append(d...)

	clusterV1.Tags = types.SetNull(TfTagType)

	clusterV1.ReplicationSpecs, d = newTfReplicationSpecsFromV0(ctx, clusterV0.ReplicationSpecs)
	diags.Append(d...)

	clusterV1.AdvancedConfiguration, d = newTfAdvancedConfigurationFromV0(ctx, clusterV0.AdvancedConfiguration)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	clusterV1.Timeouts = clusterV0.Timeouts

	return &clusterV1, diags
}

func newTfBiConnectorConfigFromV0(ctx context.Context, configV0 []TfBiConnectorConfigModel) (types.List, diag.Diagnostics) {
	return types.ListValueFrom(ctx, TfBiConnectorConfigType, configV0)
}

func newTfConnectionStringsFromV0(ctx context.Context, connStringsV0 []tfConnectionStringModel) (types.List, diag.Diagnostics) {
	return types.ListValueFrom(ctx, tfConnectionStringType, connStringsV0)
}

func newTfLabelsFromV0(ctx context.Context, labelsV0 []TfLabelModel) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, TfLabelType, removeDefaultLabel(labelsV0))
}

func newTfAdvancedConfigurationFromV0(ctx context.Context, advConfigV0 []TfAdvancedConfigurationModel) (types.List, diag.Diagnostics) {
	return types.ListValueFrom(ctx, tfAdvancedConfigurationType, advConfigV0)
}

func newTfReplicationSpecsFromV0(ctx context.Context, repSpecsV0 []tfReplicationSpecRSModelV0) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	repSpecsV1 := make([]tfReplicationSpecRSModel, len(repSpecsV0))

	for i := range repSpecsV0 {
		repSpecV0 := repSpecsV0[i]

		specV1 := tfReplicationSpecRSModel{
			ID:          repSpecV0.ID,
			ZoneName:    repSpecV0.ZoneName,
			NumShards:   repSpecV0.NumShards,
			ContainerID: repSpecV0.ContainerID,
		}
		specV1.RegionsConfigs, diags = newTfRegionConfigsFromV0(ctx, repSpecV0.RegionsConfigs)
		if diags.HasError() {
			return types.ListNull(tfRegionsConfigType), diags
		}
		repSpecsV1[i] = specV1
	}
	return types.ListValueFrom(ctx, tfReplicationSpecRSType, repSpecsV1)
}

func newTfRegionConfigsFromV0(ctx context.Context, regionConfigslV0 []tfRegionsConfigModelV0) (types.List, diag.Diagnostics) {
	regionConfigsV1 := make([]tfRegionsConfigModel, len(regionConfigslV0))

	for i := range regionConfigslV0 {
		configV0 := regionConfigslV0[i]

		configV1 := tfRegionsConfigModel{
			AnalyticsSpecs:       configV0.AnalyticsSpecs,
			AutoScaling:          configV0.AutoScaling,
			AnalyticsAutoScaling: types.ListNull(tfRegionsConfigAutoScalingSpecType),
			ReadOnlySpecs:        configV0.ReadOnlySpecs,
			ElectableSpecs:       configV0.ElectableSpecs,
			BackingProviderName:  configV0.BackingProviderName,
			ProviderName:         configV0.ProviderName,
			RegionName:           configV0.RegionName,
			Priority:             configV0.Priority,
		}

		regionConfigsV1[i] = configV1
	}
	return types.ListValueFrom(ctx, tfRegionsConfigType, regionConfigsV1)
}
