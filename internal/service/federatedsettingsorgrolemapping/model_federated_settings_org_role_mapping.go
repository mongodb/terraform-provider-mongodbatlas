package federatedsettingsorgrolemapping

import (
	"sort"
	"strings"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

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
