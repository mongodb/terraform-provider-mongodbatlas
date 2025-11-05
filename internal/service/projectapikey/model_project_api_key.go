package projectapikey

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/rolesorgid"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

func expandProjectAssignments(projectAssignments *schema.Set) map[string][]string {
	results := make(map[string][]string)
	for _, val := range projectAssignments.List() {
		results[val.(map[string]any)["project_id"].(string)] = conversion.ExpandStringList(val.(map[string]any)["role_names"].(*schema.Set).List())
	}
	return results
}

func flattenProjectAssignments(roles []admin.CloudAccessRoleAssignment) []map[string]any {
	assignments := make(map[string][]string)
	for _, role := range roles {
		if groupID := role.GetGroupId(); groupID != "" {
			assignments[groupID] = append(assignments[groupID], role.GetRoleName())
		}
	}
	var results []map[string]any
	for projectID, roles := range assignments {
		results = append(results, map[string]any{
			"project_id": projectID,
			"role_names": roles,
		})
	}
	return results
}

func getAssignmentChanges(d *schema.ResourceData) (add, remove, update map[string][]string) {
	before, after := d.GetChange("project_assignment")
	add = expandProjectAssignments(after.(*schema.Set))
	remove = expandProjectAssignments(before.(*schema.Set))
	update = make(map[string][]string)

	for projectID, rolesAfter := range add {
		if rolesBefore, ok := remove[projectID]; ok {
			if !sameRoles(rolesBefore, rolesAfter) {
				update[projectID] = rolesAfter
			}
			delete(remove, projectID)
			delete(add, projectID)
		}
	}
	return
}

func sameRoles(roles1, roles2 []string) bool {
	set1 := make(map[string]struct{})
	for _, role := range roles1 {
		set1[role] = struct{}{}
	}
	set2 := make(map[string]struct{})
	for _, role := range roles2 {
		set2[role] = struct{}{}
	}
	return reflect.DeepEqual(set1, set2)
}

// getKeyDetails returns nil error and nil details if not found as it's not considered an error
func getKeyDetails(ctx context.Context, connV2 *admin.APIClient, apiKeyID string) (*admin.ApiKeyUserDetails, string, error) {
	orgID, err := rolesorgid.GetCurrentOrgID(ctx, connV2)
	if err != nil {
		return nil, "", err
	}
	key, _, err := connV2.ProgrammaticAPIKeysApi.GetOrgApiKey(ctx, orgID, apiKeyID).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "API_KEY_NOT_FOUND") {
			return nil, orgID, nil
		}
		return nil, orgID, fmt.Errorf("error getting api key information: %s", err)
	}
	return key, orgID, nil
}

func validateUniqueProjectIDs(d *schema.ResourceData) error {
	if projectAssignments, ok := d.GetOk("project_assignment"); ok {
		uniqueIDs := make(map[string]bool)
		for _, val := range projectAssignments.(*schema.Set).List() {
			projectID := val.(map[string]any)["project_id"].(string)
			if uniqueIDs[projectID] {
				return fmt.Errorf("duplicated projectID in assignments: %s", projectID)
			}
			uniqueIDs[projectID] = true
		}
	}
	return nil
}
