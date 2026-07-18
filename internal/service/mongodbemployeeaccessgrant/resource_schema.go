package mongodbemployeeaccessgrant

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ResourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project, also known as `groupId` in the official documentation.",
			},
			"cluster_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Human-readable label that identifies this cluster.",
			},
			"grant_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Level of access to grant to MongoDB Employees. Possible values are CLUSTER_DATABASE_LOGS, CLUSTER_INFRASTRUCTURE or CLUSTER_INFRASTRUCTURE_AND_APP_SERVICES_SYNC_DATA.",
			},
			"expiration_time": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Expiration date for the employee access grant.",
			},
		},
	}
}

type TFModel struct {
	ProjectID      types.String `tfsdk:"project_id"`
	ClusterName    types.String `tfsdk:"cluster_name"`
	GrantType      types.String `tfsdk:"grant_type"`
	ExpirationTime types.String `tfsdk:"expiration_time"`
}
