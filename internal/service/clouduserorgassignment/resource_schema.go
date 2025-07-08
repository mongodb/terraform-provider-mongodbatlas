package clouduserorgassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"country": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Two-character alphabetical string that identifies the MongoDB Cloud user's geographic location. This parameter uses the ISO 3166-1a2 code format.",
			},
			"created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud created the current account. This value is in the ISO 8601 timestamp format in UTC.",
			},
			"first_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "First or given name that belongs to the MongoDB Cloud user.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the MongoDB Cloud user.",
			},
			"invitation_created_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when MongoDB Cloud sent the invitation. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.",
			},
			"invitation_expires_at": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the invitation from MongoDB Cloud expires. MongoDB Cloud represents this timestamp in ISO 8601 format in UTC.",
			},
			"inviter_username": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Username of the MongoDB Cloud user who sent the invitation to join the organization.",
			},
			"last_auth": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time when the current account last authenticated. This value is in the ISO 8601 timestamp format in UTC.",
			},
			"last_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last name, family name, or surname that belongs to the MongoDB Cloud user.",
			},
			"mobile_number": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Mobile phone number that belongs to the MongoDB Cloud user.",
			},
			"org_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [/orgs](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.",
			},
			"org_membership_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "String enum that indicates whether the MongoDB Cloud user has a pending invitation to join the organization or they are already active in the organization.",
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
					},
					"org_roles": schema.SetAttribute{
						Required:            true,
						MarkdownDescription: "One or more organization level roles to assign the MongoDB Cloud user.",
						ElementType:         types.StringType,
					},
				},
			},
			"team_ids": schema.SetAttribute{
				Computed:            true,
				MarkdownDescription: "List of unique 24-hexadecimal digit strings that identifies the teams to which this MongoDB Cloud user belongs.",
				ElementType:         types.StringType,
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email address that represents the username of the MongoDB Cloud user.",
			},
		},
	}
}

type TFModel struct {
	Country             types.String `tfsdk:"country" autogen:"omitjson"`
	CreatedAt           types.String `tfsdk:"created_at" autogen:"omitjson"`
	FirstName           types.String `tfsdk:"first_name" autogen:"omitjson"`
	Id                  types.String `tfsdk:"id" autogen:"omitjson"`
	InvitationCreatedAt types.String `tfsdk:"invitation_created_at" autogen:"omitjson"`
	InvitationExpiresAt types.String `tfsdk:"invitation_expires_at" autogen:"omitjson"`
	InviterUsername     types.String `tfsdk:"inviter_username" autogen:"omitjson"`
	LastAuth            types.String `tfsdk:"last_auth" autogen:"omitjson"`
	LastName            types.String `tfsdk:"last_name" autogen:"omitjson"`
	MobileNumber        types.String `tfsdk:"mobile_number" autogen:"omitjson"`
	OrgId               types.String `tfsdk:"org_id" autogen:"omitjson"`
	OrgMembershipStatus types.String `tfsdk:"org_membership_status" autogen:"omitjson"`
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
