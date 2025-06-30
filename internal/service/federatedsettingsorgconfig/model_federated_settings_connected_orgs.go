package federatedsettingsorgconfig

import (
	"sort"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312004/admin"
)

type roleMappingsByGroupName []admin.AuthFederationRoleMapping

func (ra roleMappingsByGroupName) Len() int      { return len(ra) }
func (ra roleMappingsByGroupName) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }

func (ra roleMappingsByGroupName) Less(i, j int) bool {
	return ra[i].ExternalGroupName < ra[j].ExternalGroupName
}

func FlattenRoleMappings(roleMappings []admin.AuthFederationRoleMapping) []map[string]any {
	sort.Sort(roleMappingsByGroupName(roleMappings))

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

type mRoleAssignment []admin.ConnectedOrgConfigRoleAssignment

func (ra mRoleAssignment) Len() int      { return len(ra) }
func (ra mRoleAssignment) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }
func (ra mRoleAssignment) Less(i, j int) bool {
	compareVal := strings.Compare(ra[i].GetOrgId(), ra[j].GetOrgId())

	if compareVal != 0 {
		return compareVal < 0
	}

	compareVal = strings.Compare(ra[i].GetGroupId(), ra[j].GetGroupId())

	if compareVal != 0 {
		return compareVal < 0
	}

	return ra[i].GetRole() < ra[j].GetRole()
}

func FlattenRoleAssignments(roleAssignments []admin.ConnectedOrgConfigRoleAssignment) []map[string]any {
	sort.Sort(mRoleAssignment(roleAssignments))

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
