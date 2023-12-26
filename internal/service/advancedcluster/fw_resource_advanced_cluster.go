package advancedcluster

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customtypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/planmodifiers"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	errorClusterAdvancedCreate             = "error creating MongoDB ClusterAdvanced: %s"
	errorClusterAdvancedRead               = "error reading MongoDB ClusterAdvanced (%s): %s"
	errorClusterAdvancedDelete             = "error deleting MongoDB ClusterAdvanced (%s): %s"
	errorClusterAdvancedUpdate             = "error updating MongoDB ClusterAdvanced (%s): %s"
	errorAdvancedClusterAdvancedConfUpdate = "error updating Advanced Configuration Option for MongoDB Cluster (%s): %s"
	errorAdvancedClusterAdvancedConfRead   = "error reading Advanced Configuration Option for MongoDB Cluster (%s): %s"
	errorInvalidCreateValues               = "Invalid values. Unable to CREATE advanced_cluster"
	defaultTimeout                         = (3 * time.Hour)
	defaultString                          = ""
	DefaultZoneName                        = "ZoneName managed by Terraform"
)

var _ resource.ResourceWithConfigure = &advancedClusterRS{}
var _ resource.ResourceWithImportState = &advancedClusterRS{}
var _ resource.ResourceWithUpgradeState = &advancedClusterRS{}

type advancedClusterRS struct {
	config.RSCommon
}

func Resource() resource.Resource {
	return &advancedClusterRS{
		RSCommon: config.RSCommon{
			ResourceName: AdvancedClusterResourceName,
		},
	}
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

func (r *advancedClusterRS) Schema(ctx context.Context, request resource.SchemaRequest, response *resource.SchemaResponse) {
	s := schema.Schema{
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
			"retain_backups_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Flag that indicates whether to retain backup snapshots for the deleted dedicated cluster",
			},
			"cluster_type": schema.StringAttribute{
				Required: true,
			},
			"connection_strings": advancedClusterRSConnectionStringSchemaComputed(),
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
			"advanced_configuration": advancedClusterRSAdvancedConfigurationSchema(),
			"bi_connector_config":    advancedClusterRSBiConnectorConfigSchema(),
			"replication_specs":      advancedClusterRSReplicationSpecsSchema(),
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

func advancedClusterRSConnectionStringSchemaComputed() schema.ListNestedAttribute {
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
	}
}

func advancedClusterRSBiConnectorConfigSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
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
	}
}

func advancedClusterRSAdvancedConfigurationSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
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
	}
}

func advancedClusterRSReplicationSpecsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
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
				},
				"region_configs": schema.ListNestedAttribute{
					Optional: true,
					Computed: true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"backing_provider_name": schema.StringAttribute{
								Optional: true,
								Computed: true,
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
							"analytics_auto_scaling": advancedClusterRSRegionConfigAutoScalingSpecsSchema(),
							"auto_scaling":           advancedClusterRSRegionConfigAutoScalingSpecsSchema(),
							"analytics_specs":        advancedClusterRSRegionConfigSpecsSchema(),
							"electable_specs":        advancedClusterRSRegionConfigSpecsSchema(),
							"read_only_specs":        advancedClusterRSRegionConfigSpecsSchema(),
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
	}
}

func advancedClusterRSRegionConfigSpecsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"disk_iops": schema.Int64Attribute{
					Optional: true,
					Computed: true,
				},
				"ebs_volume_type": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},
				"instance_size": schema.StringAttribute{
					Required: true,
				},
				"node_count": schema.Int64Attribute{
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

func advancedClusterRSRegionConfigAutoScalingSpecsSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
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

func (r *advancedClusterRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	conn := r.Client.Atlas

	var state tfAdvancedClusterRSModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
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
		resp.Diagnostics.AddError("Unable to READ cluster. An error occurred when getting cluster details from Atlas", err.Error())
		return
	}

	newClusterModel, diags := newTfAdvancedClusterRSModel(ctx, conn, cluster, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newClusterModel)...)
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

// TODO UpgradeState implements resource.ResourceWithUpgradeState.
func (*advancedClusterRS) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	schemaV0 := TPFResourceV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &schemaV0,
			StateUpgrader: upgradeAdvancedClusterResourceStateV0toV1,
		},
	}
}

func resourceClusterAdvancedRefreshFunc(ctx context.Context, name, projectID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		c, resp, err := client.AdvancedClusters.Get(ctx, projectID, name)

		if err != nil && strings.Contains(err.Error(), "reset by peer") {
			return nil, "REPEATING", nil
		}

		if err != nil && c == nil && resp == nil {
			return nil, "", err
		}

		if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		if c.StateName != "" {
			log.Printf("[DEBUG] status for MongoDB cluster: %s: %s", name, c.StateName)
		}

		return c, c.StateName, nil
	}
}

func splitSClusterAdvancedImportID(id string) (projectID, clusterName *string, err error) {
	var re = regexp.MustCompile(`(?s)^([0-9a-fA-F]{24})-(.*)$`)
	parts := re.FindStringSubmatch(id)

	if len(parts) != 3 {
		err = errors.New("import format error: to import a advanced cluster, use the format {project_id}-{name}")
		return
	}

	projectID = &parts[1]
	clusterName = &parts[2]

	return
}

func getAdvancedClusterContainerID(containers []matlas.Container, cluster *matlas.AdvancedRegionConfig) string {
	if len(containers) != 0 {
		for i := range containers {
			if cluster.ProviderName == "GCP" {
				return containers[i].ID
			}

			if containers[i].ProviderName == cluster.ProviderName &&
				containers[i].Region == cluster.RegionName || // For Azure
				containers[i].RegionName == cluster.RegionName { // For AWS
				return containers[i].ID
			}
		}
	}

	return ""
}

func doesAdvancedReplicationSpecMatchAPI(tfObject *tfReplicationSpecRSModel, apiObject *matlas.AdvancedReplicationSpec) bool {
	return tfObject.ID.ValueString() == apiObject.ID || (tfObject.ID.IsNull() && tfObject.ZoneName.ValueString() == apiObject.ZoneName)
}

func removeDefaultLabel(labels []TfLabelModel) []TfLabelModel {
	result := make([]TfLabelModel, 0)

	for _, item := range labels {
		if item.Key.ValueString() == DefaultLabel.Key && item.Value.ValueString() == DefaultLabel.Value {
			continue
		}
		result = append(result, item)
	}

	return result
}

func newTfReplicationSpecsRSModel(ctx context.Context, conn *matlas.Client,
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

			if !doesAdvancedReplicationSpecMatchAPI(&tfMapObject, apiObjects[j]) {
				continue
			}

			advancedReplicationSpec, diags := newTfReplicationSpecRSModel(ctx, apiObjects[j], &tfMapObject, conn, projectID)
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
		advancedReplicationSpec, diags := newTfReplicationSpecRSModel(ctx, apiObjects[j], tfMapObject, conn, projectID)

		if diags.HasError() {
			return nil, diags
		}

		tfList[i] = *advancedReplicationSpec
		wasAPIObjectUsed[j] = true
	}

	return tfList, nil
}

func newTfAdvancedClusterRSModel(ctx context.Context, conn *matlas.Client, cluster *matlas.AdvancedCluster, state *tfAdvancedClusterRSModel) (*tfAdvancedClusterRSModel, diag.Diagnostics) {
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

	clusterModel.BiConnectorConfig, d = types.ListValueFrom(ctx, TfBiConnectorConfigType, newTfBiConnectorConfigModel(cluster.BiConnector))
	diags.Append(d...)

	clusterModel.ConnectionStrings, d = types.ListValueFrom(ctx, tfConnectionStringType, newTfConnectionStringsModel(ctx, cluster.ConnectionStrings))
	diags.Append(d...)

	clusterModel.Labels, d = types.SetValueFrom(ctx, TfLabelType, removeDefaultLabel(newTfLabelsModel(cluster.Labels)))
	if len(clusterModel.Labels.Elements()) == 0 {
		// clusterModel.Labels, d = types.SetValueFrom(ctx, TfLabelType, []TfLabelModel{})
		clusterModel.Labels = types.SetNull(TfLabelType)
	}
	diags.Append(d...)

	clusterModel.Tags, d = types.SetValueFrom(ctx, TfTagType, newTfTagsModel(&cluster.Tags))
	if len(clusterModel.Tags.Elements()) == 0 {
		// clusterModel.Tags, d = types.SetValueFrom(ctx, TfTagType, []TfTagModel{})
		clusterModel.Tags = types.SetNull(TfTagType)
	}
	diags.Append(d...)

	// var repSpecs []tfReplicationSpecRSModel
	// if isCreate{
	// 	repSpecs, d = newTfReplicationSpecsRS(ctx, conn, cluster.ReplicationSpecs, types.ListNull(tfReplicationSpecRSType), projectID)

	// }else{
	// 	repSpecs, d = newTfReplicationSpecsRS(ctx, conn, cluster.ReplicationSpecs, state.ReplicationSpecs, projectID)

	// }
	repSpecs, d := newTfReplicationSpecsRSModel(ctx, conn, cluster.ReplicationSpecs, state.ReplicationSpecs, projectID)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	clusterModel.ReplicationSpecs, diags = types.ListValueFrom(ctx, tfReplicationSpecRSType, repSpecs)
	diags.Append(d...)

	advancedConfiguration, err := newTfAdvancedConfigurationModelDSFromAtlas(ctx, conn, projectID, name)
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
