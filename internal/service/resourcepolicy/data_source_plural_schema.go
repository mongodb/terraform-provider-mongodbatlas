package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourcePluralSchema(ctx context.Context) schema.Schema {
	dsAttributes := dataSourceSchema(true)
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"org_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
			},
			"resource_policies": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: dsAttributes,
				},
				Computed: true,
			},
		},
	}
}

type TFModelDSP struct {
	OrgID            types.String `tfsdk:"org_id"`
	ResourcePolicies []TFModel    `tfsdk:"resource_policies"`
}
