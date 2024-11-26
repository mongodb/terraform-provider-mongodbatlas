package resourcepolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourcePluralSchema(ctx context.Context) schema.Schema {
	dsAttributes1 := dataSourceSchema(true)
	dsAttributes2 := dataSourceSchema(true)
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"org_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
			},
			"resource_policies": schema.ListNestedAttribute{
				DeprecationMessage: fmt.Sprintf(constant.DeprecationParamWithReplacement, "`results`"),
				NestedObject: schema.NestedAttributeObject{
					Attributes: dsAttributes1,
				},
				Computed: true,
			},
			"results": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: dsAttributes2,
				},
				Computed: true,
			},
		},
	}
}

type TFModelDSP struct {
	OrgID            types.String `tfsdk:"org_id"`
	ResourcePolicies []TFModel    `tfsdk:"resource_policies"`
	Results          []TFModel    `tfsdk:"results"`
}
