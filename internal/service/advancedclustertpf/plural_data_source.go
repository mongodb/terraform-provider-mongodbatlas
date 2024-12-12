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
)

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

func PluralDataSource() datasource.DataSource {
	return &pluralDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", resourceName),
		},
	}
}

type pluralDS struct {
	config.DSCommon
}

func (d *pluralDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields: []string{"project_id"},
		OverridenRootFields: map[string]schema.Attribute{
			"use_replication_spec_per_shard": schema.BoolAttribute{ // TODO: added as in current resource
				Optional:            true,
				MarkdownDescription: "use_replication_spec_per_shard", // TODO: add documentation
			},
			"include_deleted_with_retained_backups": schema.BoolAttribute{ // TODO: not in current resource, decide if keep
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether to return Clusters with retain backups.",
			},
		},
		OverridenFields: map[string]schema.Attribute{
			"use_replication_spec_per_shard": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "use_replication_spec_per_shard", // TODO: add documentation
			},
		},
	})
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModelPluralDS
	diags := &resp.Diagnostics
	diags.Append(req.Config.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	model := readClustersDS(ctx, diags, d.Client, &state, &resp.State)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func readClustersDS(ctx context.Context, diags *diag.Diagnostics, client *config.MongoDBClient, pluralModel *TFModelPluralDS, stateObj *tfsdk.State) *TFModelPluralDS {
	projectID := pluralModel.ProjectID.ValueString()
	api := client.AtlasV2.ClustersApi
	list, _, err := api.ListClusters(ctx, projectID).Execute()
	if err != nil {
		diags.AddError("errorList", fmt.Sprintf(errorList, projectID, err.Error()))
		return nil
	}
	outs := &TFModelPluralDS{
		ProjectID:                         pluralModel.ProjectID,
		UseReplicationSpecPerShard:        pluralModel.UseReplicationSpecPerShard,
		IncludeDeletedWithRetainedBackups: pluralModel.IncludeDeletedWithRetainedBackups,
	}
	// TODO: use result cluster info to avoid extra API calls
	for i := range list.GetResults() {
		stateDS := &TFModelDS{
			ProjectID:                  pluralModel.ProjectID,
			Name:                       types.StringPointerValue(list.GetResults()[i].Name),
			UseReplicationSpecPerShard: pluralModel.UseReplicationSpecPerShard,
		}
		if modelDS := readClusterDS(ctx, diags, client, stateDS, stateObj); modelDS != nil {
			outs.Results = append(outs.Results, modelDS)
		}
	}
	return outs
}

type TFModelPluralDS struct {
	ProjectID                         types.String `tfsdk:"project_id"`
	Results                           []*TFModelDS `tfsdk:"results"`
	UseReplicationSpecPerShard        types.Bool   `tfsdk:"use_replication_spec_per_shard"`        // TODO: added as in current resource
	IncludeDeletedWithRetainedBackups types.Bool   `tfsdk:"include_deleted_with_retained_backups"` // TODO: not in current resource, decide if keep
}
