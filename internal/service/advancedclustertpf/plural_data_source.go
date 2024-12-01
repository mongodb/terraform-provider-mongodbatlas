package advancedclustertpf

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
	requiredFields := []string{"project_id"}
	overridenRootFields := map[string]schema.Attribute{
		"use_replication_spec_per_shard": schema.BoolAttribute{ // TODO: added as in current resource
			Optional:            true,
			MarkdownDescription: "use_replication_spec_per_shard", // TODO: add documentation
		},
		"include_deleted_with_retained_backups": schema.BoolAttribute{ // TODO: not in current resource, decide if keep
			Optional:            true,
			MarkdownDescription: "Flag that indicates whether to return Clusters with retain backups.",
		},
	}
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil, overridenRootFields)
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}

type AdvancedClustersModel struct {
	ProjectID                         types.String `tfsdk:"project_id"`
	UseReplicationSpecPerShard        types.Bool   `tfsdk:"use_replication_spec_per_shard"`        // TODO: added as in current resource
	IncludeDeletedWithRetainedBackups types.Bool   `tfsdk:"include_deleted_with_retained_backups"` // TODO: not in current resource, decide if keep
}
