package searchdeployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Label that identifies the cluster to create Search Nodes for.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the search deployment.",
			},
			"specs": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of settings that configure the Search Nodes for your cluster.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_size": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Hardware specification for the Search Node instance sizes.",
						},
						"node_count": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "Number of Search Nodes in the cluster.",
						},
					},
				},
			},
			"state_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Human-readable label that indicates the current operating condition of this search deployment.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

type TFModel struct {
	ClusterName types.String   `tfsdk:"cluster_name"`
	ProjectId   types.String   `tfsdk:"project_id"`
	Id          types.String   `tfsdk:"id"`
	Specs       types.List     `tfsdk:"specs"`
	StateName   types.String   `tfsdk:"state_name"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}
type TFSpecsModel struct {
	InstanceSize types.String `tfsdk:"instance_size"`
	NodeCount    types.Int64  `tfsdk:"node_count"`
}

var SpecsObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"instance_size": types.StringType,
	"node_count":    types.Int64Type,
}}
