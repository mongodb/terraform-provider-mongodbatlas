package searchdeployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides a Search Deployment data source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies the search deployment.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the search deployment.",
			},
			"cluster_name": schema.StringAttribute{
				Required:            true,
				Description:         "Label that identifies the cluster to return the search nodes for.",
				MarkdownDescription: "Label that identifies the cluster to return the search nodes for.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
			"specs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_size": schema.StringAttribute{
							Computed:            true,
							Description:         "Hardware specification for the search node instance sizes.",
							MarkdownDescription: "Hardware specification for the search node instance sizes.",
						},
						"node_count": schema.Int64Attribute{
							Computed:            true,
							Description:         "Number of search nodes in the cluster.",
							MarkdownDescription: "Number of search nodes in the cluster.",
						},
					},
				},
				Description:         "List of settings that configure the search nodes for your cluster.",
				MarkdownDescription: "List of settings that configure the search nodes for your cluster.",
			},
			"state_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Human-readable label that indicates the current operating condition of this search deployment.",
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
