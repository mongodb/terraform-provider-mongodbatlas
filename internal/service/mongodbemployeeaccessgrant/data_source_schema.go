package mongodbemployeeaccessgrant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchemaDelete(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"cluster_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable label that identifies this cluster.",
			},
			"grant_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Level of access to grant to MongoDB Employees. Possible values are CLUSTER_DATABASE_LOGS, CLUSTER_INFRASTRUCTURE or CLUSTER_INFRASTRUCTURE_AND_APP_SERVICES_SYNC_DATA.",
			},
			"expiration_time": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Expiration date for the employee access grant.",
			},
		},
	}
}
