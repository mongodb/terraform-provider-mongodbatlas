package federatedsettingsorgrolemapping

import (
	"sort"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115013/admin"
)

type mRoleAssignment []admin.RoleAssignment

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

func FlattenRoleAssignments(roleAssignments []admin.RoleAssignment) []map[string]any {
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
// func FlattenRoleAssignmentsLegacy(roleAssignments []*matlas.RoleAssignments) []map[string]any {
// 	sort.Sort(mRoleAssignment(roleAssignments))

// 	var roleAssignmentsMap []map[string]any

// 	if len(roleAssignments) > 0 {
// 		roleAssignmentsMap = make([]map[string]any, len(roleAssignments))

// 		for i := range roleAssignments {
// 			roleAssignmentsMap[i] = map[string]any{
// 				"group_id": roleAssignments[i].GroupID,
// 				"org_id":   roleAssignments[i].OrgID,
// 				"role":     roleAssignments[i].Role,
// 			}
// 		}
// 	}

// 	return roleAssignmentsMap
// }
