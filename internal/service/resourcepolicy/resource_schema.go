package resourcepolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
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
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				Required:            true,
			},
			"org_id": schema.StringAttribute{
				Description:         "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
				Required:            true,
			},
			"policies": schema.ListNestedAttribute{
				Description:         "List of policies that make up the Atlas resource policy.",
				MarkdownDescription: "List of policies that make up the Atlas resource policy.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"body": schema.StringAttribute{
							Description:         "A string that defines the permissions for the policy. The syntax used is the Cedar Policy language.",
							MarkdownDescription: "A string that defines the permissions for the policy. The syntax used is the Cedar Policy language.",
							Required:            true,
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type TFModel struct {
	CreatedByUser     types.Object    `tfsdk:"created_by_user"`
	CreatedDate       types.String    `tfsdk:"created_date"`
	ID                types.String    `tfsdk:"id"`
	LastUpdatedByUser types.Object    `tfsdk:"last_updated_by_user"`
	LastUpdatedDate   types.String    `tfsdk:"last_updated_date"`
	Name              types.String    `tfsdk:"name"`
	OrgID             types.String    `tfsdk:"org_id"`
	Version           types.String    `tfsdk:"version"`
	Policies          []TFPolicyModel `tfsdk:"policies"`
}

type TFUserMetadataModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

var UserMetadataObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"id":   types.StringType,
	"name": types.StringType,
}}

type TFPolicyModel struct {
	Body types.String `tfsdk:"body"`
	ID   types.String `tfsdk:"id"`
}
