package advancedcluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
			// "retain_backups_enabled": schema.BoolAttribute{
			// 	Optional:    true,
			// 	Description: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster",
			// },
			"cluster_type": schema.StringAttribute{
				Required: true,
			},
			"connection_strings": advancedClusterRSConnectionStringSchemaComputed(), //checked
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
				Computed: true, // exists in previous schema
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
			// "tags": schema.SetNestedAttribute{
			// 	Optional: true,
			// 	NestedObject: schema.NestedAttributeObject{
			// 		Attributes: map[string]schema.Attribute{
			// 			"key": schema.StringAttribute{
			// 				Required: true,
			// 			},
			// 			"value": schema.StringAttribute{
			// 				Required: true,
			// 			},
			// 		},
			// 	},
			// },
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
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
		// 	},
		// },
	}

	// if s.Blocks == nil {
	// 	s.Blocks = make(map[string]schema.Block)
	// }
	// s.Blocks["timeouts"] = timeouts.Block(context.Background(), timeouts.Opts{
	// 	Create: true,
	// 	Update: true,
	// 	Delete: true,
	// })
	return s
}

func advancedClusterRSReplicationSpecsSchemaV0() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		Optional: true,
		Computed: true,
		// Required: true,
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
					// Required: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"backing_provider_name": schema.StringAttribute{
								Optional: true,
								// Computed: true,
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
							// "analytics_auto_scaling": advancedClusterRSRegionConfigAutoScalingSpecsSchema(),
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
			// PlanModifiers: []planmodifier.Object{
			// 	objectplanmodifier.UseStateForUnknown(),
			// },
		},
	}
}

func advancedClusterRSRegionConfigSpecsSchemaV0() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		// Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"disk_iops": schema.Int64Attribute{
					Optional: true,
					Computed: true,
				},
				"ebs_volume_type": schema.StringAttribute{
					Optional: true,
					// Computed: true,
				},
				"instance_size": schema.StringAttribute{
					Required: true,
				},
				"node_count": schema.Int64Attribute{
					Optional: true,
					// Computed: true,
				},
			},
		},
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
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
	// arr := clusterV0.BiConnector
	// biConnectorConfig := mongodbatlas.BiConnector{
	// 	Enabled:        arr[0].Enabled.ValueBoolPointer(),
	// 	ReadPreference: arr[0].ReadPreference.ValueString(),
	// }

	// clusterV1 := tfAdvancedClusterRSModel{
	// 	ID: clusterV0.ID,
	// }
	// clusterV1.BiConnectorConfig, diags = types.ListValueFrom(ctx, TfBiConnectorConfigType, newTfBiConnectorConfigModel(&biConnectorConfig))
	// if diags.HasError() {
	// 	resp.Diagnostics.Append(diags...)
	// 	return
	// }

	diags = resp.State.Set(ctx, clusterV1)
	resp.Diagnostics.Append(diags...)
}

type tfAdvancedClusterRSModelV0 struct {
	ProjectID                    types.String   `tfsdk:"project_id"`
	ClusterID                    types.String   `tfsdk:"cluster_id"`
	BackupEnabled                types.Bool     `tfsdk:"backup_enabled"`
	ClusterType                  types.String   `tfsdk:"cluster_type"`
	CreateDate                   types.String   `tfsdk:"create_date"`
	DiskSizeGb                   types.Float64  `tfsdk:"disk_size_gb"`
	EncryptionAtRestProvider     types.String   `tfsdk:"encryption_at_rest_provider"`
	MongoDBMajorVersion          types.String   `tfsdk:"mongo_db_major_version"`
	MongoDBVersion               types.String   `tfsdk:"mongo_db_version"`
	Name                         types.String   `tfsdk:"name"`
	PitEnabled                   types.Bool     `tfsdk:"pit_enabled"`
	Paused                       types.Bool     `tfsdk:"paused"`
	RootCertType                 types.String   `tfsdk:"root_cert_type"`
	StateName                    types.String   `tfsdk:"state_name"`
	VersionReleaseSystem         types.String   `tfsdk:"version_release_system"`
	TerminationProtectionEnabled types.Bool     `tfsdk:"termination_protection_enabled"`
	ID                           types.String   `tfsdk:"id"`
	Timeouts                     timeouts.Value `tfsdk:"timeouts"`

	// ConnectionStrings            types.List     `tfsdk:"connection_strings"`
	// Labels                       types.Set      `tfsdk:"labels"`
	// AdvancedConfiguration        types.List     `tfsdk:"advanced_configuration"`
	// BiConnector                  types.List     `tfsdk:"bi_connector"`
	// ReplicationSpecs             types.Set      `tfsdk:"replication_specs"`
	ConnectionStrings     []tfConnectionStringModel      `tfsdk:"connection_strings"`
	Labels                []TfLabelModel                 `tfsdk:"labels"`
	AdvancedConfiguration []TfAdvancedConfigurationModel `tfsdk:"advanced_configuration"`
	BiConnector           []TfBiConnectorConfigModel     `tfsdk:"bi_connector"`
	ReplicationSpecs      []tfReplicationSpecRSModelV0   `tfsdk:"replication_specs"`

	// RetainBackupsEnabled         types.Bool     `tfsdk:"retain_backups_enabled"`

}

type tfReplicationSpecRSModelV0 struct {
	RegionsConfigs []tfRegionsConfigModelV0 `tfsdk:"region_configs"`
	ContainerID    types.Map                `tfsdk:"container_id"`
	ID             types.String             `tfsdk:"id"`
	ZoneName       types.String             `tfsdk:"zone_name"`
	NumShards      types.Int64              `tfsdk:"num_shards"`
}

var tfReplicationSpecRSTypeV0 = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":             types.StringType,
	"zone_name":      types.StringType,
	"num_shards":     types.Int64Type,
	"container_id":   types.MapType{ElemType: types.StringType},
	"region_configs": types.SetType{ElemType: tfRegionsConfigTypeV0},
},
}

type tfRegionsConfigModelV0 struct {
	AnalyticsSpecs types.List `tfsdk:"analytics_specs"`
	AutoScaling    types.List `tfsdk:"auto_scaling"`
	// AnalyticsAutoScaling types.List   `tfsdk:"analytics_auto_scaling"`
	ReadOnlySpecs       types.List   `tfsdk:"read_only_specs"`
	ElectableSpecs      types.List   `tfsdk:"electable_specs"`
	BackingProviderName types.String `tfsdk:"backing_provider_name"`
	ProviderName        types.String `tfsdk:"provider_name"`
	RegionName          types.String `tfsdk:"region_name"`
	Priority            types.Int64  `tfsdk:"priority"`
}

var tfRegionsConfigTypeV0 = types.ObjectType{AttrTypes: map[string]attr.Type{
	"backing_provider_name": types.StringType,
	"priority":              types.Int64Type,
	"provider_name":         types.StringType,
	"region_name":           types.StringType,
	"analytics_specs":       types.ListType{ElemType: tfRegionsConfigSpecType},
	"electable_specs":       types.ListType{ElemType: tfRegionsConfigSpecType},
	"read_only_specs":       types.ListType{ElemType: tfRegionsConfigSpecType},
	"auto_scaling":          types.ListType{ElemType: tfRegionsConfigAutoScalingSpecType},
	// "analytics_auto_scaling": types.ListType{ElemType: tfRegionsConfigAutoScalingSpecType},
}}

func newTfAdvancedClusterRSModelV1(ctx context.Context, clusterV0 *tfAdvancedClusterRSModelV0) (*tfAdvancedClusterRSModel, diag.Diagnostics) {
	var d, diags diag.Diagnostics
	// projectID := cluster.GroupID
	// name := cluster.Name

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
		// AcceptDataRisksAndForceReplicaSetReconfig: conversion.StringNullIfEmpty(cluster.AcceptDataRisksAndForceReplicaSetReconfig),
		ProjectID: clusterV0.ProjectID,
		Timeouts:  clusterV0.Timeouts,
		// RetainBackupsEnabled: state.RetainBackupsEnabled,
	}

	// clusterModel.ID = types.StringValue(conversion.EncodeStateID(map[string]string{
	// 	"cluster_id":   cluster.ID,
	// 	"project_id":   projectID,
	// 	"cluster_name": name,
	// }))

	clusterV1.BiConnectorConfig, d = newTfBiConnectorConfigFromV0(ctx, clusterV0.BiConnector)
	diags.Append(d...)

	clusterV1.ConnectionStrings, d = newTfConnectionStringsFromV0(ctx, clusterV0.ConnectionStrings)
	diags.Append(d...)

	clusterV1.Labels, d = newTfLabelsFromV0(ctx, clusterV0.Labels)
	// if len(clusterV1.Labels.Elements()) == 0 {
	// 	clusterV1.Labels = types.SetNull(TfLabelType)
	// }
	diags.Append(d...)

	clusterV1.Tags = types.SetNull(TfTagType)

	// clusterV1.Tags, d = types.SetValueFrom(ctx, TfTagType, newTfTagsModel(&cluster.Tags))
	// if len(clusterV1.Tags.Elements()) == 0 {
	// 	clusterV1.Tags = types.SetNull(TfTagType)
	// }
	// diags.Append(d...)

	// repSpecs, d := newTfReplicationSpecsFromV0(ctx, clusterV0.ReplicationSpecs)
	// diags.Append(d...)
	// if diags.HasError() {
	// 	return nil, diags
	// }
	clusterV1.ReplicationSpecs, d = newTfReplicationSpecsFromV0(ctx, clusterV0.ReplicationSpecs)
	diags.Append(d...)

	// advancedConfiguration, err := newTfAdvancedConfigurationModelDSFromAtlas(ctx, conn, projectID, name)
	// if err != nil {
	// 	diags.AddError("An error occurred when getting advanced_configuration from Atlas", err.Error())
	// 	return nil, diags
	// }
	clusterV1.AdvancedConfiguration, d = newTfAdvancedConfigurationFromV0(ctx, clusterV0.AdvancedConfiguration)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	clusterV1.Timeouts = clusterV0.Timeouts

	return &clusterV1, diags
}

func newTfBiConnectorConfigFromV0(ctx context.Context, configV0 []TfBiConnectorConfigModel) (types.List, diag.Diagnostics) {
	// biConnectorConfig := mongodbatlas.BiConnector{
	// 	Enabled:        configV0[0].Enabled.ValueBoolPointer(),
	// 	ReadPreference: configV0[0].ReadPreference.ValueString(),
	// }

	// biConnectorConfigV1, diags := types.ListValueFrom(ctx, TfBiConnectorConfigType, newTfBiConnectorConfigModel(&biConnectorConfig))
	biConnectorConfigV1, diags := types.ListValueFrom(ctx, TfBiConnectorConfigType, configV0)

	return biConnectorConfigV1, diags
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
