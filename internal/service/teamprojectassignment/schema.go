package teamprojectassignment

import (
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func resourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project, also known as `groupId` in the official documentation.",
			},
			"role_names": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "One or more project-level roles assigned to the team.",
			},
			"team_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the team.",
			},
		},
	}
}

func dataSourceSchema() dsschema.Schema {
	return conversion.DataSourceSchemaFromResource(resourceSchema(), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "team_id"},
	})
}

type TFModel struct {
	ProjectId types.String `tfsdk:"project_id"`
	RoleNames types.Set    `tfsdk:"role_names"`
	TeamId    types.String `tfsdk:"team_id"`
}
