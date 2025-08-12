package clouduserprojectassignment

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func resourceSchema(ctx context.Context) schema.Schema {
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
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
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
			"org_membership_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "String enum that indicates whether the MongoDB Cloud user has a pending invitation to join the organization or they are already active in the organization.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"roles": schema.SetAttribute{
				ElementType:         types.StringType,
				Required:            true,
				MarkdownDescription: "One or more project-level roles to assign the MongoDB Cloud user.",
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Email address that represents the username of the MongoDB Cloud user.",
			},
		},
	}
}

type TFModel struct {
	Country             types.String `tfsdk:"country"`
	CreatedAt           types.String `tfsdk:"created_at"`
	FirstName           types.String `tfsdk:"first_name"`
	ProjectId           types.String `tfsdk:"project_id"`
	UserId              types.String `tfsdk:"user_id"`
	InvitationCreatedAt types.String `tfsdk:"invitation_created_at"`
	InvitationExpiresAt types.String `tfsdk:"invitation_expires_at"`
	InviterUsername     types.String `tfsdk:"inviter_username"`
	LastAuth            types.String `tfsdk:"last_auth"`
	LastName            types.String `tfsdk:"last_name"`
	MobileNumber        types.String `tfsdk:"mobile_number"`
	OrgMembershipStatus types.String `tfsdk:"org_membership_status"`
	Roles               types.Set    `tfsdk:"roles"`
	Username            types.String `tfsdk:"username"`
}
