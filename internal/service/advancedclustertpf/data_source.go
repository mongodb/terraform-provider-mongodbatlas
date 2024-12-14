package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

var _ datasource.DataSource = &ds{}
var _ datasource.DataSourceWithConfigure = &ds{}

func DataSource() datasource.DataSource {
	return &ds{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type ds struct {
	config.DSCommon
}

func (d *ds) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "name"},
		OverridenFields: map[string]schema.Attribute{
			"use_replication_spec_per_shard": schema.BoolAttribute{ // TODO: added as in current resource
				Optional:            true,
				MarkdownDescription: "use_replication_spec_per_shard", // TODO: add documentation
			},
		},
	})
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModelDS
	diags := &resp.Diagnostics
	diags.Append(req.Config.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	model := d.readCluster(ctx, diags, &state)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *ds) readCluster(ctx context.Context, diags *diag.Diagnostics, modelDS *TFModelDS) *TFModelDS {
	clusterName := modelDS.Name.ValueString()
	projectID := modelDS.ProjectID.ValueString()
	useReplicationSpecPerShard := modelDS.UseReplicationSpecPerShard.ValueBool()
	api := d.Client.AtlasV2.ClustersApi
	clusterResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) {
			return nil
		}
		diags.AddError("errorRead", fmt.Sprintf(errorRead, clusterName, err.Error()))
		return nil
	}
	modelIn := &TFModel{
		ProjectID: modelDS.ProjectID,
		Name:      modelDS.Name,
	}
	// TODO: pass !UseReplicationSpecPerShard to overrideUsingLegacySchema
	modelOut, extraInfo := getBasicClusterModel(ctx, diags, d.Client, clusterResp, modelIn)
	if diags.HasError() {
		return nil
	}
	if extraInfo.AsymmetricShardUnsupportedError && !useReplicationSpecPerShard {
		diags.AddError("errorRead", "Please add `use_replication_spec_per_shard = true` to your data source configuration to enable asymmetric shard support. Refer to documentation for more details.")
		return nil
	}
	updateModelAdvancedConfig(ctx, diags, d.Client, modelOut, nil, nil)
	if diags.HasError() {
		return nil
	}
	modelOutDS, err := conversion.CopyModel[TFModelDS](modelOut)
	if err != nil {
		diags.AddError(errorRead, fmt.Sprintf("error setting model: %s", err.Error()))
		return nil
	}
	modelOutDS.UseReplicationSpecPerShard = modelDS.UseReplicationSpecPerShard // attrs not in resource model
	return modelOutDS
}

// TODO: difference with TFModel: misses timeouts, adds use_replication_spec_per_shard.
type TFModelDS struct {
	DiskSizeGB                                types.Float64 `tfsdk:"disk_size_gb"`
	Labels                                    types.Set     `tfsdk:"labels"`
	ReplicationSpecs                          types.List    `tfsdk:"replication_specs"`
	Tags                                      types.Set     `tfsdk:"tags"`
	ReplicaSetScalingStrategy                 types.String  `tfsdk:"replica_set_scaling_strategy"`
	Name                                      types.String  `tfsdk:"name"`
	AdvancedConfiguration                     types.Object  `tfsdk:"advanced_configuration"`
	BiConnectorConfig                         types.Object  `tfsdk:"bi_connector_config"`
	RootCertType                              types.String  `tfsdk:"root_cert_type"`
	ClusterType                               types.String  `tfsdk:"cluster_type"`
	MongoDBMajorVersion                       types.String  `tfsdk:"mongo_db_major_version"`
	ConfigServerType                          types.String  `tfsdk:"config_server_type"`
	VersionReleaseSystem                      types.String  `tfsdk:"version_release_system"`
	ConnectionStrings                         types.Object  `tfsdk:"connection_strings"`
	StateName                                 types.String  `tfsdk:"state_name"`
	MongoDBVersion                            types.String  `tfsdk:"mongo_db_version"`
	CreateDate                                types.String  `tfsdk:"create_date"`
	AcceptDataRisksAndForceReplicaSetReconfig types.String  `tfsdk:"accept_data_risks_and_force_replica_set_reconfig"`
	EncryptionAtRestProvider                  types.String  `tfsdk:"encryption_at_rest_provider"`
	ProjectID                                 types.String  `tfsdk:"project_id"`
	ClusterID                                 types.String  `tfsdk:"cluster_id"`
	ConfigServerManagementMode                types.String  `tfsdk:"config_server_management_mode"`
	PinnedFCV                                 types.Object  `tfsdk:"pinned_fcv"`
	UseReplicationSpecPerShard                types.Bool    `tfsdk:"use_replication_spec_per_shard"`
	RedactClientLogData                       types.Bool    `tfsdk:"redact_client_log_data"`
	GlobalClusterSelfManagedSharding          types.Bool    `tfsdk:"global_cluster_self_managed_sharding"`
	BackupEnabled                             types.Bool    `tfsdk:"backup_enabled"`
	RetainBackupsEnabled                      types.Bool    `tfsdk:"retain_backups_enabled"`
	Paused                                    types.Bool    `tfsdk:"paused"`
	TerminationProtectionEnabled              types.Bool    `tfsdk:"termination_protection_enabled"`
	PitEnabled                                types.Bool    `tfsdk:"pit_enabled"`
}
