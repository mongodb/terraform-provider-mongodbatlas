package conversion

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func FlattenLinks(links []admin.Link) []map[string]string {
	ret := make([]map[string]string, len(links))
	for i, link := range links {
		ret[i] = map[string]string{
			"href": link.GetHref(),
			"rel":  link.GetRel(),
		}
	}
	return ret
}

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
	ret := []map[string]any{}
	roleMap := map[string]any{
		"org_roles":     []string{},
		"project_roles": []map[string]any{},
	}
	if roles.HasOrgRoles() {
		roleMap["org_roles"] = roles.GetOrgRoles()
	}
	if roles.HasGroupRoleAssignments() {
		roleMap["project_roles"] = FlattenGroupRolesAssignments(roles.GetGroupRoleAssignments())
	}
	ret = append(ret, roleMap)
	return ret
}

func FlattenGroupRolesAssignments(assignments []admin.GroupRoleAssignment) []map[string]any {
	ret := make([]map[string]any, len(assignments))
	for i, assignment := range assignments {
		ret[i] = map[string]any{
			"group_id":    assignment.GetGroupId(),
			"group_roles": assignment.GetGroupRoles(),
		}
	}
	return ret
}

func FlattenTags(tags []admin.ResourceTag) []map[string]string {
	ret := make([]map[string]string, len(tags))
	for i, tag := range tags {
		ret[i] = map[string]string{
			"key":   tag.GetKey(),
			"value": tag.GetValue(),
		}
	}
	return ret
}

func ExpandTagsFromSetSchema(d *schema.ResourceData) *[]admin.ResourceTag {
	list := d.Get("tags").(*schema.Set)
	ret := make([]admin.ResourceTag, list.Len())
	for i, item := range list.List() {
		tag := item.(map[string]any)
		ret[i] = admin.ResourceTag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
	}
	return &ret
}

func ExpandStringList(list []any) (res []string) {
	for _, v := range list {
		res = append(res, v.(string))
	}
	return
}

func ExpandStringListFromSetSchema(set *schema.Set) []string {
	res := make([]string, set.Len())
	for i, v := range set.List() {
		res[i] = v.(string)
	}
	return res
}
