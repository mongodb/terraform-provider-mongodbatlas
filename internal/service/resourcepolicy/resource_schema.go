package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"created_by_user": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The user that last updated the atlas resource policy.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Unique 24-hexadecimal character string that identifies a user.",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Human-readable label that describes a user.",
					},
				},
			},
			"created_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time in UTC when the atlas resource policy was created.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the atlas resource policy.",
			},
			"last_updated_by_user": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The user that last updated the atlas resource policy.",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Unique 24-hexadecimal character string that identifies a user.",
					},
					"name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "Human-readable label that describes a user.",
					},
				},
			},
			"last_updated_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time in UTC when the atlas resource policy was last updated.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Human-readable label that describes the atlas resource policy.",
			},
			"org_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
			},
			"policies": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of policies that make up the atlas resource policy.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"body": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "A string that defines the permissions for the policy. The syntax used is the Cedar Policy language.",
						},
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Unique 24-hexadecimal character string that identifies the policy.",
						},
					},
				},
			},
			"version": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "A string that identifies the version of the atlas resource policy.",
			},
		},
	}
}

type TFModel struct {
	CreatedByUser     types.Object `tfsdk:"created_by_user"`
	CreatedDate       types.String `tfsdk:"created_date"`
	Id                types.String `tfsdk:"id"`
	LastUpdatedByUser types.Object `tfsdk:"last_updated_by_user"`
	LastUpdatedDate   types.String `tfsdk:"last_updated_date"`
	Name              types.String `tfsdk:"name"`
	OrgId             types.String `tfsdk:"org_id"`
	Policies          types.List   `tfsdk:"policies"`
	Version           types.String `tfsdk:"version"`
}
type TFCreatedByUserModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

var CreatedByUserObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":   types.StringType,
	"name": types.StringType,
}}

type TFLastUpdatedByUserModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

var LastUpdatedByUserObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":   types.StringType,
	"name": types.StringType,
}}

type TFPoliciesModel struct {
	Body types.String `tfsdk:"body"`
	Id   types.String `tfsdk:"id"`
}

var PoliciesObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"body": types.StringType,
	"id":   types.StringType,
}}
