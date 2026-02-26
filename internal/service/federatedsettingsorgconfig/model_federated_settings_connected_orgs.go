package federatedsettingsorgconfig

import (
	"cmp"
	"slices"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

func compareRoleMappingsByGroupName(a, b admin.AuthFederationRoleMapping) int {
	return cmp.Compare(a.ExternalGroupName, b.ExternalGroupName)
}

func compareRoleAssignment(a, b admin.ConnectedOrgConfigRoleAssignment) int {
	if c := cmp.Compare(a.GetOrgId(), b.GetOrgId()); c != 0 {
		return c
	}
	if c := cmp.Compare(a.GetGroupId(), b.GetGroupId()); c != 0 {
		return c
	}
	return cmp.Compare(a.GetRole(), b.GetRole())
}

func FlattenRoleMappings(roleMappings []admin.AuthFederationRoleMapping) []map[string]any {
	slices.SortFunc(roleMappings, compareRoleMappingsByGroupName)

	var roleMappingsMap []map[string]any

	if len(roleMappings) > 0 {
		roleMappingsMap = make([]map[string]any, len(roleMappings))

		for i := range roleMappings {
			roleMappingsMap[i] = map[string]any{
				"external_group_name": roleMappings[i].GetExternalGroupName(),
				"id":                  roleMappings[i].GetId(),
				"role_assignments":    FlattenRoleAssignments(roleMappings[i].GetRoleAssignments()),
			}
		}
	}

	return roleMappingsMap
}

func FlattenRoleAssignments(roleAssignments []admin.ConnectedOrgConfigRoleAssignment) []map[string]any {
	slices.SortFunc(roleAssignments, compareRoleAssignment)

	var roleAssignmentsMap []map[string]any

	if len(roleAssignments) > 0 {
		roleAssignmentsMap = make([]map[string]any, len(roleAssignments))

		for i := range roleAssignments {
			roleAssignmentsMap[i] = map[string]any{
				"group_id": roleAssignments[i].GetGroupId(),
				"org_id":   roleAssignments[i].GetOrgId(),
				"role":     roleAssignments[i].GetRole(),
			}
		}
	}

	return roleAssignmentsMap
}

func FlattenUserConflicts(userConflicts []admin.FederatedUser) []map[string]any {
	var userConflictsMap []map[string]any

	if len(userConflicts) == 0 {
		return nil
	}
	userConflictsMap = make([]map[string]any, len(userConflicts))

	for i := range userConflicts {
		userConflictsMap[i] = map[string]any{
			"email_address":          userConflicts[i].GetEmailAddress(),
			"federation_settings_id": userConflicts[i].GetFederationSettingsId(),
			"first_name":             userConflicts[i].GetFirstName(),
			"last_name":              userConflicts[i].GetLastName(),
			"user_id":                userConflicts[i].GetUserId(),
		}
	}

	return userConflictsMap
}
