package advancedclustertpf

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TODO: see if we can leverage on resource or singular data source schema, e.g. have a func to add computed
func PluralDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"use_replication_spec_per_shard": schema.BoolAttribute{ // TODO: added as in current resource
				Optional:            true,
				MarkdownDescription: "use_replication_spec_per_shard", // TODO: add documentation
			},
			"include_deleted_with_retained_backups": schema.BoolAttribute{ // TODO: not in current resource, decide if keep
				Optional:            true,
				MarkdownDescription: "Flag that indicates whether to return Clusters with retain backups.",
			},
			"results": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of returned documents that MongoDB Cloud provides when completing this request.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"project_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Human-readable label that identifies this cluster.",
						},
					},
				},
			},
			"total_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Total number of documents available. MongoDB Cloud omits this value if `includeCount` is set to `false`.",
			},
		},
	}
}

type AdvancedClustersModel struct {
	ProjectID                         types.String `tfsdk:"project_id"`
	ItemsPerPage                      types.Int64  `tfsdk:"items_per_page"`
	PageNum                           types.Int64  `tfsdk:"page_num"`
	TotalCount                        types.Int64  `tfsdk:"total_count"`
	IncludeCount                      types.Bool   `tfsdk:"include_count"`
	IncludeDeletedWithRetainedBackups types.Bool   `tfsdk:"include_deleted_with_retained_backups"`
}
