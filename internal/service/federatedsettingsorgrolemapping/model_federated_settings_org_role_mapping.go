package federatedsettingsorgrolemapping

import (
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
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

func expandRoleAssignments(d *schema.ResourceData) *[]admin.RoleAssignment {
	var roleAssignments []admin.RoleAssignment

	if v, ok := d.GetOk("role_assignments"); ok {
		if rs := v.(*schema.Set); rs.Len() > 0 {
			for _, r := range rs.List() {
				roleMap := r.(map[string]any)

				for _, role := range roleMap["roles"].(*schema.Set).List() {
					roleAssignment := admin.RoleAssignment{
						OrgId:   conversion.StringPtr(roleMap["org_id"].(string)),
						GroupId: conversion.StringPtr(roleMap["group_id"].(string)),
						Role:    conversion.StringPtr(role.(string)),
					}
					roleAssignments = append(roleAssignments, roleAssignment)
				}
			}
		}
	}

	sort.Sort(mRoleAssignment(roleAssignments))
	return &roleAssignments
}

func flattenRoleAssignmentsResource(roleAssignments []admin.RoleAssignment) []map[string]any {
	if len(roleAssignments) == 0 {
		return nil
	}
	sort.Sort(mRoleAssignment(roleAssignments))
	var flattenedRoleAssignments []map[string]any
	var roleAssignment = map[string]any{
		"group_id": roleAssignments[0].GetGroupId(),
		"org_id":   roleAssignments[0].GetOrgId(),
		"roles":    []string{},
	}

	for _, row := range roleAssignments {
		if (roleAssignment["org_id"] != "" && roleAssignment["org_id"] != row.GetOrgId()) ||
			(roleAssignment["group_id"] != "" && roleAssignment["group_id"] != row.GetGroupId()) {
			flattenedRoleAssignments = append(flattenedRoleAssignments, roleAssignment)

			roleAssignment = map[string]any{
				"group_id": row.GetGroupId(),
				"org_id":   row.GetOrgId(),
				"roles":    []string{},
			}
		}

		roleAssignment["roles"] = append(roleAssignment["roles"].([]string), row.GetRole())
	}

	flattenedRoleAssignments = append(flattenedRoleAssignments, roleAssignment)
	return flattenedRoleAssignments
}
