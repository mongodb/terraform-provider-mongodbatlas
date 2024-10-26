package projectapikey

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20240805005/admin"
)

const (
	ErrorProjectSetting = "error setting `%s` for project (%s): %s"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate,
		ReadContext:   resourceRead,
		UpdateContext: resourceUpdate,
		DeleteContext: resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
		},
		Schema: map[string]*schema.Schema{
			"api_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"project_assignment": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role_names": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

type APIProjectAssignmentKeyInput struct {
	ProjectID string   `json:"desc,omitempty"`
	RoleNames []string `json:"roles,omitempty"`
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := validateUniqueProjectIDs(d); err != nil {
		return diag.FromErr(err)
	}
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	val, ok := d.GetOk("project_assignment")
	if !ok {
		return diag.FromErr(errors.New("project_assignment not found"))
	}
	assignments := expandProjectAssignmentSet(val.(*schema.Set))

	req := &admin.CreateAtlasProjectApiKey{
		Desc:  d.Get("description").(string),
		Roles: assignments[0].RoleNames,
	}
	ret, _, err := connV2.ProgrammaticAPIKeysApi.CreateProjectApiKey(ctx, assignments[0].ProjectID, req).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	apiKeyID := ret.GetId()
	d.SetId(conversion.EncodeStateID(map[string]string{
		"api_key_id": apiKeyID,
	}))
	if err := d.Set("public_key", ret.GetPublicKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}
	if err := d.Set("private_key", ret.GetPrivateKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
	}

	for _, assignment := range assignments[1:] {
		req := &[]admin.UserAccessRoleAssignment{{Roles: &assignment.RoleNames}}
		if _, _, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, assignment.ProjectID, apiKeyID, req).Execute(); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	apiKeyID := ids["api_key_id"]

	details, _, err := getKeyDetails(ctx, connV2, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if details == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("api_key_id", details.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `api_key_id`: %s", err))
	}

	if err := d.Set("description", details.GetDesc()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
	}

	if err := d.Set("public_key", details.GetPublicKey()); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("project_assignment", flattenProjectAssignmentsFromRoles(details.GetRoles())); err != nil {
		return diag.Errorf("error setting `project_assignment` : %s", err)
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	if err := validateUniqueProjectIDs(d); err != nil {
		return diag.FromErr(err)
	}

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	apiKeyID := ids["api_key_id"]

	details, orgID, err := getKeyDetails(ctx, connV2, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if details == nil {
		return diag.Errorf("error updating project api_key (%s): not found", apiKeyID)
	}

	if d.HasChange("project_assignment") {
		add, remove, update := getAssignmentChanges(d)

		for projectID := range remove {
			_, _, err := connV2.ProgrammaticAPIKeysApi.RemoveProjectApiKey(ctx, projectID, apiKeyID).Execute()
			if err != nil {
				if admin.IsErrorCode(err, "GROUP_NOT_FOUND") {
					continue // allows removing assignment for a project that has been deleted
				}
				return diag.Errorf("error removing project_api_key(%s) from project(%s): %s", apiKeyID, projectID, err)
			}
		}

		for projectID, roles := range add {
			req := &[]admin.UserAccessRoleAssignment{{Roles: &roles}}
			_, _, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, projectID, apiKeyID, req).Execute()
			if err != nil {
				return diag.Errorf("error adding project_api_key(%s) to project(%s): %s", apiKeyID, projectID, err)
			}
		}

		for projectID, roles := range update {
			req := &admin.UpdateAtlasProjectApiKey{Roles: &roles}
			_, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKeyRoles(ctx, projectID, apiKeyID, req).Execute()
			if err != nil {
				return diag.Errorf("error changing project_api_key(%s) in project(%s): %s", apiKeyID, projectID, err)
			}
		}
	}

	if d.HasChange("description") {
		req := &admin.UpdateAtlasOrganizationApiKey{Desc: conversion.StringPtr(d.Get("description").(string))}
		if _, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKey(ctx, orgID, apiKeyID, req).Execute(); err != nil {
			return diag.Errorf("error updating description in api key(%s): %s", apiKeyID, err)
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	apiKeyID := ids["api_key_id"]
	details, orgID, err := getKeyDetails(ctx, connV2, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if details != nil && orgID != "" {
		if _, _, err = connV2.ProgrammaticAPIKeysApi.DeleteApiKey(ctx, orgID, apiKeyID).Execute(); err != nil {
			return diag.FromErr(fmt.Errorf("error deleting project key (%s): %s", apiKeyID, err))
		}
	}
	d.SetId("")
	return nil
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	d.SetId(conversion.EncodeStateID(map[string]string{
		"api_key_id": d.Id(),
	}))
	return []*schema.ResourceData{d}, nil
}

func expandProjectAssignmentSet(projectAssignments *schema.Set) []*APIProjectAssignmentKeyInput {
	res := make([]*APIProjectAssignmentKeyInput, projectAssignments.Len())
	for i, value := range projectAssignments.List() {
		v := value.(map[string]any)
		res[i] = &APIProjectAssignmentKeyInput{
			ProjectID: v["project_id"].(string),
			RoleNames: conversion.ExpandStringList(v["role_names"].(*schema.Set).List()),
		}
	}
	return res
}

func flattenProjectAssignmentsFromRoles(roles []admin.CloudAccessRoleAssignment) []map[string]any {
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
	add = make(map[string][]string)
	remove = make(map[string][]string)
	update = make(map[string][]string)
	before, after := d.GetChange("project_assignment")
	for _, val := range after.(*schema.Set).List() {
		add[val.(map[string]any)["project_id"].(string)] = conversion.ExpandStringList(val.(map[string]any)["role_names"].(*schema.Set).List())
	}
	for _, val := range before.(*schema.Set).List() {
		remove[val.(map[string]any)["project_id"].(string)] = conversion.ExpandStringList(val.(map[string]any)["role_names"].(*schema.Set).List())
	}

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
	root, _, err := connV2.RootApi.GetSystemStatus(ctx).Execute()
	if err != nil {
		return nil, "", err
	}
	for _, role := range root.ApiKey.GetRoles() {
		if orgID := role.GetOrgId(); orgID != "" {
			key, _, err := connV2.ProgrammaticAPIKeysApi.GetApiKey(ctx, orgID, apiKeyID).Execute()
			if err != nil {
				if admin.IsErrorCode(err, "API_KEY_NOT_FOUND") {
					return nil, orgID, nil
				}
				return nil, orgID, fmt.Errorf("error getting api key information: %s", err)
			}
			return key, orgID, nil
		}
	}
	return nil, "", nil
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
