package searchdeployment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies the search deployment.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the search deployment.",
			},
			"cluster_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Label that identifies the cluster to return the search nodes for.",
				MarkdownDescription: "Label that identifies the cluster to return the search nodes for.",
			},
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description:         "Unique 24-hexadecimal character string that identifies the project.",
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the project.",
			},
			"specs": schema.ListNestedAttribute{
				Validators: []validator.List{
					listvalidator.SizeAtMost(1),
					listvalidator.SizeAtLeast(1),
				},
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_size": schema.StringAttribute{
							Required:            true,
							Description:         "Hardware specification for the search node instance sizes. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Atlas-Search/operation/createAtlasSearchDeployment) describes the valid values. More details can also be found in the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-tier).",
							MarkdownDescription: "Hardware specification for the search node instance sizes. The [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Atlas-Search/operation/createAtlasSearchDeployment) describes the valid values. More details can also be found in the [Search Node Documentation](https://www.mongodb.com/docs/atlas/cluster-config/multi-cloud-distribution/#search-tier).",
						},
						"node_count": schema.Int64Attribute{
							Required:            true,
							Description:         "Number of search nodes in the cluster.",
							MarkdownDescription: "Number of search nodes in the cluster.",
						},
					},
				},
				Description:         "List of settings that configure the search nodes for your cluster. This list is currently limited to defining a single element.",
				MarkdownDescription: "List of settings that configure the search nodes for your cluster. This list is currently limited to defining a single element.",
			},
			"state_name": schema.StringAttribute{
				Computed:            true,
				Description:         "Human-readable label that indicates the current operating condition of this search deployment.",
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

type TFSearchDeploymentRSModel struct {
	ID          types.String   `tfsdk:"id"`
	ClusterName types.String   `tfsdk:"cluster_name"`
	ProjectID   types.String   `tfsdk:"project_id"`
	Specs       types.List     `tfsdk:"specs"`
	StateName   types.String   `tfsdk:"state_name"`
	Timeouts    timeouts.Value `tfsdk:"timeouts"`
}

type TFSearchNodeSpecModel struct {
	InstanceSize types.String `tfsdk:"instance_size"`
	NodeCount    types.Int64  `tfsdk:"node_count"`
}

var SpecObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"instance_size": types.StringType,
	"node_count":    types.Int64Type,
}}
