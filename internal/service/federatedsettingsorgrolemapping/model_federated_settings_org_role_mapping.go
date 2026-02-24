package federatedsettingsorgrolemapping

import (
	"cmp"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

func compareRoleAssignment(a, b admin.ConnectedOrgConfigRoleAssignment) int {
	if c := cmp.Compare(a.GetOrgId(), b.GetOrgId()); c != 0 {
		return c
	}
	if c := cmp.Compare(a.GetGroupId(), b.GetGroupId()); c != 0 {
		return c
	}
	return cmp.Compare(a.GetRole(), b.GetRole())
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

func expandRoleAssignments(d *schema.ResourceData) *[]admin.ConnectedOrgConfigRoleAssignment {
	var roleAssignments []admin.ConnectedOrgConfigRoleAssignment

	if v, ok := d.GetOk("role_assignments"); ok {
		if rs := v.(*schema.Set); rs.Len() > 0 {
			for _, r := range rs.List() {
				roleMap := r.(map[string]any)

				for _, role := range roleMap["roles"].(*schema.Set).List() {
					roleAssignment := admin.ConnectedOrgConfigRoleAssignment{
						OrgId:   conversion.StringPtr(roleMap["org_id"].(string)),
						GroupId: conversion.StringPtr(roleMap["group_id"].(string)),
						Role:    conversion.StringPtr(role.(string)),
					}
					roleAssignments = append(roleAssignments, roleAssignment)
				}
			}
		}
	}

	slices.SortFunc(roleAssignments, compareRoleAssignment)
	return &roleAssignments
}

func flattenRoleAssignmentsResource(roleAssignments []admin.ConnectedOrgConfigRoleAssignment) []map[string]any {
	if len(roleAssignments) == 0 {
		return nil
	}
	slices.SortFunc(roleAssignments, compareRoleAssignment)
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
