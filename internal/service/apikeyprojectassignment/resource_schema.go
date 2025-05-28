package apikeyprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies this organization API key that you want to assign to one project.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"role_names": schema.SetAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable label that identifies the collection of privileges that MongoDB Cloud grants a specific API key, MongoDB Cloud user, or MongoDB Cloud team. These roles include only the specific project-level roles.",
				ElementType:         types.StringType,
			},
		},
	}
}

type TFModel struct {
	ApiKeyId  types.String `tfsdk:"api_key_id"`
	ProjectId types.String `tfsdk:"project_id"`
	RoleNames types.Set    `tfsdk:"role_names"`
}
