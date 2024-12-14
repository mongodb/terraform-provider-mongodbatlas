package advancedclustertpf

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
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
	model := d.readClusters(ctx, diags, &state)
	if model != nil {
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *pluralDS) readClusters(ctx context.Context, diags *diag.Diagnostics, pluralModel *TFModelPluralDS) *TFModelPluralDS {
	projectID := pluralModel.ProjectID.ValueString()
	useReplicationSpecPerShard := pluralModel.UseReplicationSpecPerShard.ValueBool()
	api := d.Client.AtlasV2.ClustersApi
	params := admin.ListClustersApiParams{
		GroupId: projectID,
	}
	list, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ClusterDescription20240805], *http.Response, error) {
		request := api.ListClustersWithParams(ctx, &params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		diags.AddError("errorList", fmt.Sprintf(errorList, projectID, err.Error()))
		return nil
	}
	outs := &TFModelPluralDS{
		ProjectID:                         pluralModel.ProjectID,
		UseReplicationSpecPerShard:        pluralModel.UseReplicationSpecPerShard,
		IncludeDeletedWithRetainedBackups: pluralModel.IncludeDeletedWithRetainedBackups,
	}
	for i := range list {
		clusterResp := &list[i]
		modelIn := &TFModel{
			ProjectID: pluralModel.ProjectID,
			Name:      types.StringValue(clusterResp.GetName()),
		}
		// TODO: pass !UseReplicationSpecPerShard to overrideUsingLegacySchema
		modelOut, extraInfo := getBasicClusterModel(ctx, diags, d.Client, clusterResp, modelIn)
		if diags.HasError() {
			return nil
		}
		if extraInfo.AsymmetricShardUnsupported && !useReplicationSpecPerShard {
			continue
		}
		updateModelAdvancedConfig(ctx, diags, d.Client, modelOut, nil, nil)
		if diags.HasError() {
			return nil
		}
		modelOutDS, err := conversion.CopyModel[TFModelDS](modelOut)
		if err != nil {
			diags.AddError(errorList, fmt.Sprintf("error setting model: %s", err.Error()))
			return nil
		}
		modelOutDS.UseReplicationSpecPerShard = pluralModel.UseReplicationSpecPerShard // attrs not in resource model
		outs.Results = append(outs.Results, modelOutDS)
	}
	return outs
}

type TFModelPluralDS struct {
	ProjectID                         types.String `tfsdk:"project_id"`
	Results                           []*TFModelDS `tfsdk:"results"`
	UseReplicationSpecPerShard        types.Bool   `tfsdk:"use_replication_spec_per_shard"`        // TODO: added as in current resource
	IncludeDeletedWithRetainedBackups types.Bool   `tfsdk:"include_deleted_with_retained_backups"` // TODO: not in current resource, decide if keep
}
