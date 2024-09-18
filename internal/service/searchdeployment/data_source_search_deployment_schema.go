package searchdeployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the search deployment.",
			},
			"cluster_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Label that identifies the cluster to return the search nodes for.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
			"specs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_size": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Hardware specification for the search node instance sizes. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Atlas-Search/operation/createAtlasSearchDeployment) describes the valid values. More details can also be found in the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-tier).",
						},
						"node_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "Number of search nodes in the cluster.",
						},
					},
				},
				MarkdownDescription: "List of settings that configure the search nodes for your cluster. This list is currently limited to defining a single element.",
			},
			"state_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable label that indicates the current operating condition of this search deployment.",
			},
		},
	}
}

type TFSearchDeploymentDSModel struct {
	ID          types.String `tfsdk:"id"`
	ClusterName types.String `tfsdk:"cluster_name"`
	ProjectID   types.String `tfsdk:"project_id"`
	Specs       types.List   `tfsdk:"specs"`
	StateName   types.String `tfsdk:"state_name"`
}
