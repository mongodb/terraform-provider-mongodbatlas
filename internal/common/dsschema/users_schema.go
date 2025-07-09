package dsschema

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

var (
	DSOrgUsersSchema = schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"org_membership_status": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"roles": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"org_roles": {
								Type:     schema.TypeSet,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"project_roles_assignments": {
								Type:     schema.TypeSet,
								Computed: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"project_id": {
											Type:     schema.TypeString,
											Computed: true,
										},
										"project_roles": {
											Type:     schema.TypeSet,
											Computed: true,
											Elem:     &schema.Schema{Type: schema.TypeString},
										},
									},
								},
							},
						},
					},
				},
				"team_ids": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"username": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"invitation_created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"invitation_expires_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"inviter_username": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"country": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"created_at": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"first_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"last_auth": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"last_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mobile_number": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
)

func FlattenUsers(users []admin.OrgUserResponse) []map[string]any {
	ret := make([]map[string]any, len(users))
	for i := range users {
		user := &users[i]
		ret[i] = map[string]any{
			"id":                    user.GetId(),
			"org_membership_status": user.GetOrgMembershipStatus(),
			"roles":                 FlattenUserRoles(user.GetRoles()),
			"team_ids":              user.GetTeamIds(),
			"username":              user.GetUsername(),
			"invitation_created_at": user.GetInvitationCreatedAt().Format(time.RFC3339),
			"invitation_expires_at": user.GetInvitationExpiresAt().Format(time.RFC3339),
			"inviter_username":      user.GetInviterUsername(),
			"country":               user.GetCountry(),
			"created_at":            user.GetCreatedAt().Format(time.RFC3339),
			"first_name":            user.GetFirstName(),
			"last_auth":             user.GetLastAuth().Format(time.RFC3339),
			"last_name":             user.GetLastName(),
			"mobile_number":         user.GetMobileNumber(),
		}
	}
	return ret
}

func FlattenUserRoles(roles admin.OrgUserRolesResponse) []map[string]any {
	ret := make([]map[string]any, 0)
	roleMap := map[string]any{
		"org_roles":                 []string{},
		"project_roles_assignments": []map[string]any{},
	}
	if roles.HasOrgRoles() {
		roleMap["org_roles"] = roles.GetOrgRoles()
	}
	if roles.HasGroupRoleAssignments() {
		roleMap["project_roles_assignments"] = FlattenProjectRolesAssignments(roles.GetGroupRoleAssignments())
	}
	ret = append(ret, roleMap)
	return ret
}

func FlattenProjectRolesAssignments(assignments []admin.GroupRoleAssignment) []map[string]any {
	ret := make([]map[string]any, 0, len(assignments))
	for _, assignment := range assignments {
		ret = append(ret, map[string]any{
			"project_id":    assignment.GetGroupId(),
			"project_roles": assignment.GetGroupRoles(),
		})
	}
	return ret
}
