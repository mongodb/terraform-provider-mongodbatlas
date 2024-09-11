package mongodbemployeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"cluster_name": schema.StringAttribute{
				Required:            true,
				Description:         "Human-readable label that identifies this cluster.",
				MarkdownDescription: "Human-readable label that identifies this cluster.",
			},
			"grant_type": schema.StringAttribute{
				Required:            true,
				Description:         "Level of access to grant to MongoDB Employees.",
				MarkdownDescription: "Level of access to grant to MongoDB Employees.",
			},
			"expiration_time": schema.StringAttribute{
				Required:            true,
				Description:         "Expiration date for the employee access grant.",
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
