package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: dataSourceSchema(false),
	}
}

func dataSourceSchema(isPlural bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"created_by_user": schema.SingleNestedAttribute{
			Description:         "The user that last updated the Atlas resource policy.",
			MarkdownDescription: "The user that last updated the Atlas resource policy.",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description:         "Unique 24-hexadecimal character string that identifies a user.",
					MarkdownDescription: "Unique 24-hexadecimal character string that identifies a user.",
					Computed:            true,
				},
				"name": schema.StringAttribute{
					Description:         "Human-readable label that describes a user.",
					MarkdownDescription: "Human-readable label that describes a user.",
					Computed:            true,
				},
			},
		},
		"created_date": schema.StringAttribute{
			Description:         "Date and time in UTC when the Atlas resource policy was created.",
			MarkdownDescription: "Date and time in UTC when the Atlas resource policy was created.",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			Description:         "Unique 24-hexadecimal digit string that identifies an Atlas resource policy.",
			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies an Atlas resource policy.",
			Required:            !isPlural,
			Computed:            isPlural,
		},
		"last_updated_by_user": schema.SingleNestedAttribute{
			Description:         "The user that last updated the Atlas resource policy.",
			MarkdownDescription: "The user that last updated the Atlas resource policy.",
			Computed:            true,
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description:         "Unique 24-hexadecimal character string that identifies a user.",
					MarkdownDescription: "Unique 24-hexadecimal character string that identifies a user.",
					Computed:            true,
				},
				"name": schema.StringAttribute{
					Description:         "Human-readable label that describes a user.",
					MarkdownDescription: "Human-readable label that describes a user.",
					Computed:            true,
				},
			},
		},
		"last_updated_date": schema.StringAttribute{
			Description:         "Date and time in UTC when the Atlas resource policy was last updated.",
			MarkdownDescription: "Date and time in UTC when the Atlas resource policy was last updated.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Human-readable label that describes the Atlas resource policy.",
			MarkdownDescription: "Human-readable label that describes the Atlas resource policy.",
			Computed:            true,
		},
		"org_id": schema.StringAttribute{
			Required:            !isPlural,
			Computed:            isPlural,
			Description:         "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
		},
		"policies": schema.ListNestedAttribute{
			Description: "List of policies that make up the Atlas resource policy.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"body": schema.StringAttribute{
						Description:         "A string that defines the permissions for the policy. The syntax used is the Cedar Policy language.",
						MarkdownDescription: "A string that defines the permissions for the policy. The syntax used is the Cedar Policy language.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						Description:         "Unique 24-hexadecimal character string that identifies the policy.",
						MarkdownDescription: "Unique 24-hexadecimal character string that identifies the policy.",
						Computed:            true,
					},
				},
			},
		},
		"version": schema.StringAttribute{
			Description:         "A string that identifies the version of the Atlas resource policy.",
			MarkdownDescription: "A string that identifies the version of the Atlas resource policy.",
			Computed:            true,
		},
	}
}
