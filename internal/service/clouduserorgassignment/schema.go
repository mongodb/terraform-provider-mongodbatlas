package clouduserorgassignment

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func resourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"country": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Two-character alphabetical string that identifies the MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud created the current account. This value is in the ISO 8601 timestamp format in UTC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "First or given name that belongs to the MongoDB Cloud user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invitation_created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud sent the invitation. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invitation_expires_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the invitation from MongoDB Cloud expires. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inviter_username": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Username of the MongoDB Cloud user who sent the invitation to join the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_auth": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last name, family name, or surname that belongs to the MongoDB Cloud user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mobile_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Mobile phone number that belongs to the MongoDB Cloud user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
			},
			"org_membership_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "String enum that indicates whether the MongoDB Cloud user has a pending invitation to join the organization or they are already active in the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"roles": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Organization and project level roles to assign the MongoDB Cloud user within one organization.",
				Attributes: map[string]schema.Attribute{
					"project_role_assignments": schema.ListNestedAttribute{
						Computed:            true,
						MarkdownDescription: "List of project level role assignments to assign the MongoDB Cloud user.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"project_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the project to which these roles belong.",
								},
								"project_roles": schema.SetAttribute{
									Computed:            true,
									MarkdownDescription: "One or more project-level roles assigned to the MongoDB Cloud user.",
									ElementType:         types.StringType,
								},
							},
						},
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"org_roles": schema.SetAttribute{
						Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
						Optional:            true,
						MarkdownDescription: "One or more organization level roles to assign the MongoDB Cloud user.",
						ElementType:         types.StringType,
					},
				},
			},
			"team_ids": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Cloud user belongs.",
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Email address that represents the username of the MongoDB Cloud user.",
			},
		},
	}
}

// func dataSourceSchema(ctx context.Context) dsschema.Schema {
// 	return conversion.DataSourceSchemaFromResource(resourceSchema(ctx), &conversion.DataSourceSchemaRequest{
// 		RequiredFields: []string{"org_id"},

// 		OverridenFields: dataSourceOverridenFields(),
// 	})
// }

// func dataSourceOverridenFields() map[string]dsschema.Attribute {
// 	return map[string]dsschema.Attribute{
// 		"user_id": dsschema.BoolAttribute{
// 			Optional:            true,
// 			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.",
// 		},
// 		"username": dsschema.BoolAttribute{
// 			Optional:            true,
// 			MarkdownDescription: "Email address that represents the username of the MongoDB Cloud user.",
// 		},
// 	}
// }

type TFModel struct {
	Country             types.String `tfsdk:"country"`
	CreatedAt           types.String `tfsdk:"created_at"`
	FirstName           types.String `tfsdk:"first_name"`
	UserId              types.String `tfsdk:"user_id"`
	InvitationCreatedAt types.String `tfsdk:"invitation_created_at"`
	InvitationExpiresAt types.String `tfsdk:"invitation_expires_at"`
	InviterUsername     types.String `tfsdk:"inviter_username"`
	LastAuth            types.String `tfsdk:"last_auth"`
	LastName            types.String `tfsdk:"last_name"`
	MobileNumber        types.String `tfsdk:"mobile_number"`
	OrgId               types.String `tfsdk:"org_id"`
	OrgMembershipStatus types.String `tfsdk:"org_membership_status"`
	Roles               types.Object `tfsdk:"roles"`
	TeamIds             types.Set    `tfsdk:"team_ids"`
	Username            types.String `tfsdk:"username" autogen:"omitjsonupdate"`
}
type TFRolesModel struct {
	ProjectRoleAssignments types.List `tfsdk:"project_role_assignments"`
	OrgRoles               types.Set  `tfsdk:"org_roles"`
}
type TFRolesProjectRoleAssignmentsModel struct {
	ProjectId    types.String `tfsdk:"project_id"`
	ProjectRoles types.Set    `tfsdk:"project_roles"`
}

var ProjectRoleAssignmentsAttrType = types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
	"project_id":    types.StringType,
	"project_roles": types.SetType{ElemType: types.StringType},
}}}

var RolesObjectAttrTypes = map[string]attr.Type{
	"org_roles":                types.SetType{ElemType: types.StringType},
	"project_role_assignments": ProjectRoleAssignmentsAttrType,
}
