package clouduserteamassignment

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func resourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"org_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
			},
			"team_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the team to which you want to assign the MongoDB Cloud user. Use the [/teams](#tag/Teams/operation/listTeams) endpoint to retrieve all teams to which the authenticated user has access.",
			},
			"user_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.",
			},
			"username": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Email address that represents the username of the MongoDB Cloud user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_membership_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "String enum that indicates whether the MongoDB Cloud user has a pending invitation to join the organization or they are already active in the organization.",
			},
			"roles": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Organization and project level roles to assign the MongoDB Cloud user within one organization.",
				Attributes: map[string]schema.Attribute{
					"project_role_assignments": schema.SetNestedAttribute{
						Computed:            true,
						MarkdownDescription: "List of project level role assignments to assign the MongoDB Cloud user.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"project_id": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the project to which these roles belong.",
								},
								"project_roles": schema.SetAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									MarkdownDescription: "One or more project-level roles assigned to the MongoDB Cloud user.",
								},
							},
						},
					},
					"org_roles": schema.SetAttribute{
						ElementType:         types.StringType,
						Computed:            true,
						MarkdownDescription: "One or more organization level roles to assign the MongoDB Cloud user.",
					},
				},
			},
			"team_ids": schema.SetAttribute{
				ElementType:         types.StringType,
				Computed:            true,
				MarkdownDescription: "List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Cloud user belongs.",
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
			"country": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Two-character alphabetical string that identifies the MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.",
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
			"last_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last name, family name, or surname that belongs to the MongoDB Cloud user.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud created the current account. This value is in the ISO 8601 timestamp format in UTC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_auth": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.",
			},
			"mobile_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Mobile phone number that belongs to the MongoDB Cloud user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func dataSourceSchema() dsschema.Schema {
	return conversion.DataSourceSchemaFromResource(resourceSchema(), &conversion.DataSourceSchemaRequest{
		RequiredFields:  []string{"org_id", "team_id"},
		OverridenFields: dataSourceOverridenFields(),
	})
}

func dataSourceOverridenFields() map[string]dsschema.Attribute {
	return map[string]dsschema.Attribute{
		"user_id": dsschema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.",
		},
		"username": dsschema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Email address that represents the username of the MongoDB Cloud user.",
		},
	}
}

type TFUserTeamAssignmentModel struct {
	OrgId               types.String `tfsdk:"org_id"`
	TeamId              types.String `tfsdk:"team_id"`
	UserId              types.String `tfsdk:"user_id"`
	Username            types.String `tfsdk:"username"`
	OrgMembershipStatus types.String `tfsdk:"org_membership_status"`
	Roles               types.Object `tfsdk:"roles"`
	TeamIds             types.Set    `tfsdk:"team_ids"`
	InvitationCreatedAt types.String `tfsdk:"invitation_created_at"`
	InvitationExpiresAt types.String `tfsdk:"invitation_expires_at"`
	InviterUsername     types.String `tfsdk:"inviter_username"`
	Country             types.String `tfsdk:"country"`
	FirstName           types.String `tfsdk:"first_name"`
	LastName            types.String `tfsdk:"last_name"`
	CreatedAt           types.String `tfsdk:"created_at"`
	LastAuth            types.String `tfsdk:"last_auth"`
	MobileNumber        types.String `tfsdk:"mobile_number"`
}

type TFRolesModel struct {
	ProjectRoleAssignments types.Set `tfsdk:"project_role_assignments"`
	OrgRoles               types.Set `tfsdk:"org_roles"`
}

type TFProjectRoleAssignmentsModel struct {
	ProjectId    types.String `tfsdk:"project_id"`
	ProjectRoles types.Set    `tfsdk:"project_roles"`
}

var ProjectRoleAssignmentsAttrType = types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
	"project_id":    types.StringType,
	"project_roles": types.SetType{ElemType: types.StringType},
}}}

var RolesObjectAttrTypes = map[string]attr.Type{
	"org_roles":                types.SetType{ElemType: types.StringType},
	"project_role_assignments": ProjectRoleAssignmentsAttrType,
}
