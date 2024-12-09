package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113002/admin"
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
	model := d.readCluster(ctx, &state, &resp.State, diags, true)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *ds) readCluster(ctx context.Context, model *TFModelDS, state *tfsdk.State, diags *diag.Diagnostics, allowNotFound bool) *TFModelDS {
	clusterName := model.Name.ValueString()
	projectID := model.ProjectID.ValueString()
	api := d.Client.AtlasV2.ClustersApi
	readResp, _, err := api.GetCluster(ctx, projectID, clusterName).Execute()
	if err != nil {
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) && allowNotFound {
			state.RemoveResource(ctx)
			return nil
		}
		diags.AddError("errorRead", fmt.Sprintf(errorRead, clusterName, err.Error()))
		return nil
	}
	return d.convertClusterAddAdvConfig(ctx, nil, nil, readResp, model, nil, diags)
}

func (d *ds) convertClusterAddAdvConfig(ctx context.Context, legacyAdvConfig *admin20240530.ClusterDescriptionProcessArgs, advConfig *admin.ClusterDescriptionProcessArgs20240805, cluster *admin.ClusterDescription20240805, modelIn *TFModelDS, oldAdvConfig *types.Object, diags *diag.Diagnostics) *TFModelDS {
	apiInfo := resolveAPIInfoDS(ctx, modelIn, diags, cluster, d.Client)
	if diags.HasError() {
		return nil
	}
	modelOut := NewTFModelDS(ctx, cluster, diags, *apiInfo)
	if diags.HasError() {
		return nil
	}
	modelOut.UseReplicationSpecPerShard = modelIn.UseReplicationSpecPerShard // input param

	if oldAdvConfig != nil {
		modelOut.AdvancedConfiguration = *oldAdvConfig
	} else {
		legacyAdvConfig, advConfig = readUnsetAdvancedConfigurationDS(ctx, d.Client, modelOut, legacyAdvConfig, advConfig, diags)
		AddAdvancedConfigDS(ctx, modelOut, advConfig, legacyAdvConfig, diags)
		if diags.HasError() {
			return nil
		}
	}
	overrideKnowTPFIssueFieldsDS(modelIn, modelOut)
	return modelOut
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
	UseReplicationSpecPerShard                types.Bool    `tfsdk:"use_replication_spec_per_shard"`
	RedactClientLogData                       types.Bool    `tfsdk:"redact_client_log_data"`
	GlobalClusterSelfManagedSharding          types.Bool    `tfsdk:"global_cluster_self_managed_sharding"`
	BackupEnabled                             types.Bool    `tfsdk:"backup_enabled"`
	RetainBackupsEnabled                      types.Bool    `tfsdk:"retain_backups_enabled"`
	Paused                                    types.Bool    `tfsdk:"paused"`
	TerminationProtectionEnabled              types.Bool    `tfsdk:"termination_protection_enabled"`
	PitEnabled                                types.Bool    `tfsdk:"pit_enabled"`
}
