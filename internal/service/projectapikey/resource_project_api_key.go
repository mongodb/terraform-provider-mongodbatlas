package projectapikey

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

const errorNoProjectAssignmentDefined = "could not obtain a project id as no assignments are defined"

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

	if d.HasChange("project_assignment") {
		// Getting the changes to api key project assignments
		newAssignments, changedAssignments, removedAssignments := getStateProjectAssignmentAPIKeys(d)

		// Adding new projects assignments
		if len(newAssignments) > 0 {
			for _, apiKey := range newAssignments {
				projectID := apiKey.(map[string]any)["project_id"].(string)
				roles := conversion.ExpandStringList(apiKey.(map[string]any)["role_names"].(*schema.Set).List())
				assignment := []admin.UserAccessRoleAssignment{{Roles: &roles}}
				_, _, err := connV2.ProgrammaticAPIKeysApi.AddProjectApiKey(ctx, projectID, apiKeyID, &assignment).Execute()
				if err != nil {
					return diag.Errorf("error assigning api_keys into the project(%s): %s", projectID, err)
				}
			}
		}

		// Removing projects assignments
		for _, apiKey := range removedAssignments {
			projectID := apiKey.(map[string]any)["project_id"].(string)
			_, _, err := connV2.ProgrammaticAPIKeysApi.RemoveProjectApiKey(ctx, projectID, apiKeyID).Execute()
			if err != nil && strings.Contains(err.Error(), "GROUP_NOT_FOUND") {
				continue // allows removing assignment for a project that has been deleted
			}
			if err != nil {
				return diag.Errorf("error removing api_key(%s) from the project(%s): %s", apiKeyID, projectID, err)
			}
		}

		// Updating the role names for the project assignments
		for _, apiKey := range changedAssignments {
			projectID := apiKey.(map[string]any)["project_id"].(string)
			roles := conversion.ExpandStringList(apiKey.(map[string]any)["role_names"].(*schema.Set).List())
			assignment := admin.UpdateAtlasProjectApiKey{Roles: &roles}
			_, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKeyRoles(ctx, projectID, apiKeyID, &assignment).Execute()
			if err != nil {
				return diag.Errorf("error updating role names for the api_key(%s): %s", apiKey, err)
			}
		}
	}

	firstProjectID, err := getFirstProjectIDFromAssignments(d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not obtain a project id from state: %s", err))
	}

	if d.HasChange("description") {
		newDescription := d.Get("description").(string)
		if _, _, err := connV2.ProgrammaticAPIKeysApi.UpdateApiKeyRoles(ctx, *firstProjectID, apiKeyID, &admin.UpdateAtlasProjectApiKey{
			Desc: &newDescription,
		}).Execute(); err != nil {
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
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a api key use the format {project_id}-{api_key_id}")
	}

	projectID := parts[0]
	apiKeyID := parts[1]

	projectAPIKeys, _, err := connV2.ProgrammaticAPIKeysApi.ListProjectApiKeys(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key %s in project %s, error: %s", projectID, apiKeyID, err)
	}
	for _, val := range projectAPIKeys.GetResults() {
		if val.GetId() != apiKeyID {
			continue
		}
		if err := d.Set("description", val.GetDesc()); err != nil {
			return nil, fmt.Errorf("error setting `description`: %s", err)
		}

		if err := d.Set("public_key", val.GetPublicKey()); err != nil {
			return nil, fmt.Errorf("error setting `public_key`: %s", err)
		}

		apiAssigments, err := getAPIProjectAssignments(ctx, connV2, apiKeyID)
		if err != nil {
			return nil, fmt.Errorf("error getting api key information: %s", err)
		}
		if err := d.Set("project_assignment", flattenProjectAssignments(apiAssigments)); err != nil {
			return nil, fmt.Errorf("error setting  `project_assignment`: %s", err)
		}

		d.SetId(conversion.EncodeStateID(map[string]string{
			"api_key_id": val.GetId(),
		}))
	}
	return []*schema.ResourceData{d}, nil
}

func getFirstProjectIDFromAssignments(d *schema.ResourceData) (*string, error) {
	if projectAssignments, ok := d.GetOk("project_assignment"); ok {
		projectAssignmentList := expandProjectAssignmentSet(projectAssignments.(*schema.Set))
		if len(projectAssignmentList) > 0 {
			return admin.PtrString(projectAssignmentList[0].ProjectID), nil
		}
	}
	return nil, errors.New(errorNoProjectAssignmentDefined)
}

func flattenProjectAPIKeyRoles(projectID string, apiKeyRoles []admin.CloudAccessRoleAssignment) []string {
	if len(apiKeyRoles) == 0 {
		return nil
	}

	flattenedOrgRoles := []string{}

	for _, role := range apiKeyRoles {
		if strings.HasPrefix(role.GetRoleName(), "GROUP_") && role.GetGroupId() == projectID {
			flattenedOrgRoles = append(flattenedOrgRoles, role.GetRoleName())
		}
	}

	return flattenedOrgRoles
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

func flattenProjectAssignments(assignments []APIProjectAssignmentKeyInput) []map[string]any {
	var results []map[string]any
	var atlasRoles []admin.CloudAccessRoleAssignment
	if len(assignments) > 0 {
		results = make([]map[string]any, len(assignments))
		for k, apiKey := range assignments {
			for _, roleName := range apiKey.RoleNames {
				atlasRole := admin.CloudAccessRoleAssignment{
					GroupId:  &apiKey.ProjectID,
					RoleName: &roleName,
				}
				atlasRoles = append(atlasRoles, atlasRole)
			}
			results[k] = map[string]any{
				"project_id": apiKey.ProjectID,
				"role_names": flattenProjectAPIKeyRoles(apiKey.ProjectID, atlasRoles),
			}
		}
	}
	return results
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

func getStateProjectAssignmentAPIKeys(d *schema.ResourceData) (newAssignments, changedAssignments, removedAssignments []any) {
	prevAssignments, currAssignments := d.GetChange("project_assignment")

	rAssignments := prevAssignments.(*schema.Set).Difference(currAssignments.(*schema.Set))
	nAssignments := currAssignments.(*schema.Set).Difference(prevAssignments.(*schema.Set))
	changedAssignments = make([]any, 0)

	for _, changed := range nAssignments.List() {
		for _, removed := range rAssignments.List() {
			if changed.(map[string]any)["project_id"] == removed.(map[string]any)["project_id"] {
				rAssignments.Remove(removed)
			}
		}

		for _, current := range prevAssignments.(*schema.Set).List() {
			if changed.(map[string]any)["project_id"] == current.(map[string]any)["project_id"] {
				changedAssignments = append(changedAssignments, changed.(map[string]any))
				nAssignments.Remove(changed)
			}
		}
	}

	newAssignments = nAssignments.List()
	removedAssignments = rAssignments.List()

	return
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

func getAPIProjectAssignments(ctx context.Context, connV2 *admin.APIClient, apiKeyID string) ([]APIProjectAssignmentKeyInput, error) {
	apiKeyOrgList, _, err := connV2.RootApi.GetSystemStatus(ctx).Execute()
	if err != nil {
		return nil, err
	}

	projectAssignments := []APIProjectAssignmentKeyInput{}
	for idx, role := range apiKeyOrgList.ApiKey.GetRoles() {
		if !strings.HasPrefix(*role.RoleName, "ORG_") {
			continue
		}
		roles := apiKeyOrgList.ApiKey.GetRoles()
		orgKeys, _, err := connV2.ProgrammaticAPIKeysApi.ListApiKeys(ctx, *roles[idx].OrgId).Execute()
		if err != nil {
			return nil, fmt.Errorf("error getting api key information: %s", err)
		}
		for _, val := range orgKeys.GetResults() {
			if val.GetId() == apiKeyID {
				for _, r := range val.GetRoles() {
					temp := new(APIProjectAssignmentKeyInput)
					if strings.HasPrefix(r.GetRoleName(), "GROUP_") {
						temp.ProjectID = r.GetGroupId()
						for _, l := range val.GetRoles() {
							if l.GetGroupId() == temp.ProjectID {
								temp.RoleNames = append(temp.RoleNames, l.GetRoleName())
							}
						}
						projectAssignments = append(projectAssignments, *temp)
					}
				}
			}
		}
	}
	return projectAssignments, nil
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
