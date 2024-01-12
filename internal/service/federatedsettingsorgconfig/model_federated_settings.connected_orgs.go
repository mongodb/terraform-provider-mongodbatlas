package federatedsettingsorgconfig

import (
	"sort"
	"strings"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

type roleMappingsByGroupName []*matlas.RoleMappings

func (ra roleMappingsByGroupName) Len() int      { return len(ra) }
func (ra roleMappingsByGroupName) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }

func (ra roleMappingsByGroupName) Less(i, j int) bool {
	return ra[i].ExternalGroupName < ra[j].ExternalGroupName
}

func FlattenRoleMappings(roleMappings []*matlas.RoleMappings) []map[string]any {
	sort.Sort(roleMappingsByGroupName(roleMappings))

	var roleMappingsMap []map[string]any

	if len(roleMappings) > 0 {
		roleMappingsMap = make([]map[string]any, len(roleMappings))

		for i := range roleMappings {
			roleMappingsMap[i] = map[string]any{
				"external_group_name": roleMappings[i].ExternalGroupName,
				"id":                  roleMappings[i].ID,
				"role_assignments":    FlattenRoleAssignments(roleMappings[i].RoleAssignments),
			}
		}
	}

	return roleMappingsMap
}

type mRoleAssignment []*matlas.RoleAssignments

func (ra mRoleAssignment) Len() int      { return len(ra) }
func (ra mRoleAssignment) Swap(i, j int) { ra[i], ra[j] = ra[j], ra[i] }
func (ra mRoleAssignment) Less(i, j int) bool {
	compareVal := strings.Compare(ra[i].OrgID, ra[j].OrgID)

	if compareVal != 0 {
		return compareVal < 0
	}

	compareVal = strings.Compare(ra[i].GroupID, ra[j].GroupID)

	if compareVal != 0 {
		return compareVal < 0
	}

	return ra[i].Role < ra[j].Role
}

func FlattenRoleAssignments(roleAssignments []*matlas.RoleAssignments) []map[string]any {
	sort.Sort(mRoleAssignment(roleAssignments))

	var roleAssignmentsMap []map[string]any

	if len(roleAssignments) > 0 {
		roleAssignmentsMap = make([]map[string]any, len(roleAssignments))

		for i := range roleAssignments {
			roleAssignmentsMap[i] = map[string]any{
				"group_id": roleAssignments[i].GroupID,
				"org_id":   roleAssignments[i].OrgID,
				"role":     roleAssignments[i].Role,
			}
		}
	}

	return roleAssignmentsMap
}

func FlattenUserConflicts(userConflicts matlas.UserConflicts) []map[string]any {
	var userConflictsMap []map[string]any

	if len(userConflicts) == 0 {
		return nil
	}
	userConflictsMap = make([]map[string]any, len(userConflicts))

	for i := range userConflicts {
		userConflictsMap[i] = map[string]any{
			"email_address":          userConflicts[i].EmailAddress,
			"federation_settings_id": userConflicts[i].FederationSettingsID,
			"first_name":             userConflicts[i].FirstName,
			"last_name":              userConflicts[i].LastName,
			"user_id":                userConflicts[i].UserID,
		}
	}

	return userConflictsMap
}
