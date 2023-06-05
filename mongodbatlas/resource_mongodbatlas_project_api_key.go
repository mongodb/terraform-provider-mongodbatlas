package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasProjectAPIKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasProjectAPIKeyCreate,
		ReadContext:   resourceMongoDBAtlasProjectAPIKeyRead,
		UpdateContext: resourceMongoDBAtlasProjectAPIKeyUpdate,
		DeleteContext: resourceMongoDBAtlasProjectAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasProjectAPIKeyImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
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
			"role_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"project_assignment"},
			},
			"project_assignment": {
				Type:     schema.TypeSet,
				Optional: true,
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
				ConflictsWith: []string{"role_names"},
			},
		},
	}
}

type APIProjectAssignmentKeyInput struct {
	ProjectID string   `json:"desc,omitempty"`
	RoleNames []string `json:"roles,omitempty"`
}

func resourceMongoDBAtlasProjectAPIKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	createRequest := new(matlas.APIKeyInput)

	var apiKey *matlas.APIKey
	var err error
	var resp *matlas.Response

	createRequest.Desc = d.Get("description").(string)
	if projectAssignments, ok := d.GetOk("project_assignment"); ok {
		projectAssignmentList := ExpandProjectAssignmentSet(projectAssignments.(*schema.Set))
		for _, apiKeyList := range projectAssignmentList {
			if apiKeyList.ProjectID == projectID {
				createRequest.Roles = apiKeyList.RoleNames
				apiKey, resp, err = conn.ProjectAPIKeys.Create(ctx, projectID, createRequest)
				if err != nil {
					if resp != nil && resp.StatusCode == http.StatusNotFound {
						d.SetId("")
						return nil
					}
				}
			}
		}

		for _, apiKeyList := range projectAssignmentList {
			if apiKeyList.ProjectID != projectID {
				createRequest.Roles = apiKeyList.RoleNames
				_, err := conn.ProjectAPIKeys.Assign(ctx, apiKeyList.ProjectID, apiKey.ID, &matlas.AssignAPIKey{
					Roles: createRequest.Roles,
				})
				if err != nil {
					if resp != nil && resp.StatusCode == http.StatusNotFound {
						d.SetId("")
						return nil
					}
				}
			}
		}
	} else {
		createRequest.Roles = expandStringList(d.Get("role_names").(*schema.Set).List())

		apiKey, resp, err = conn.ProjectAPIKeys.Create(ctx, projectID, createRequest)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				d.SetId("")
				return nil
			}

			return diag.FromErr(fmt.Errorf("error create API key: %s", err))
		}
	}

	if err := d.Set("public_key", apiKey.PublicKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
	}

	if err := d.Set("private_key", apiKey.PrivateKey); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `private_key`: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"api_key_id": apiKey.ID,
	}))

	return resourceMongoDBAtlasProjectAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectAPIKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	options := &matlas.ListOptions{}

	apiKeyOrgList, _, err := conn.Root.List(ctx, options)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	projectAssignments := []APIProjectAssignmentKeyInput{}

	for idx, role := range apiKeyOrgList.APIKey.Roles {
		if strings.HasPrefix(role.RoleName, "ORG_") {
			orgKeys, _, err := conn.APIKeys.List(ctx, apiKeyOrgList.APIKey.Roles[idx].OrgID, options)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
			}
			for _, val := range orgKeys {
				if val.ID == apiKeyID {
					for _, r := range val.Roles {
						temp := new(APIProjectAssignmentKeyInput)
						if strings.HasPrefix(r.RoleName, "GROUP_") {
							temp.ProjectID = r.GroupID
							for _, l := range val.Roles {
								if l.GroupID == temp.ProjectID {
									temp.RoleNames = append(temp.RoleNames, l.RoleName)
								}
							}
							projectAssignments = append(projectAssignments, *temp)
						}
					}
				}
			}
			break
		}
	}

	_, roleOk := d.GetOk("role_names")
	if !roleOk {
		if err := d.Set("role_names", nil); err != nil {
			return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
		}
		if err := d.Set("project_assignment", FlattenProjectAssignment(projectAssignments)); err != nil {
			return diag.Errorf(errorProjectSetting, `created`, projectID, err)
		}
	}

	projectAPIKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}
	for _, val := range projectAPIKeys {
		if val.ID == apiKeyID {
			if err := d.Set("api_key_id", val.ID); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `api_key_id`: %s", err))
			}

			if err := d.Set("description", val.Desc); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `description`: %s", err))
			}

			if err := d.Set("public_key", val.PublicKey); err != nil {
				return diag.FromErr(fmt.Errorf("error setting `public_key`: %s", err))
			}
			if roleOk {
				if err := d.Set("role_names", flattenProjectAPIKeyRoles(projectID, val.Roles)); err != nil {
					return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
				}
			}
		}
	}

	if err := d.Set("project_id", projectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `project_id`: %s", err))
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id": projectID,
		"api_key_id": apiKeyID,
	}))

	return nil
}

func resourceMongoDBAtlasProjectAPIKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]

	updateRequest := new(matlas.AssignAPIKey)

	if d.HasChange("role_names") {
		updateRequest.Roles = expandStringList(d.Get("role_names").(*schema.Set).List())
		if updateRequest.Roles != nil {
			_, err := conn.ProjectAPIKeys.Assign(ctx, projectID, apiKeyID, updateRequest)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error updating API key: %s", err))
			}
		}
	}

	if d.HasChange("project_assignment") {
		// get the current api_keys and the new api_keys with changes
		newAPIKeys, changedAPIKeys, removedAPIKeys := getStateProjectAssignmentAPIKeys(d)

		// adding new api_keys into the project
		if len(newAPIKeys) > 0 {
			for _, apiKey := range newAPIKeys {
				projectID := apiKey.(map[string]interface{})["project_id"].(string)
				roles := expandStringList(apiKey.(map[string]interface{})["role_names"].(*schema.Set).List())
				_, err := conn.ProjectAPIKeys.Assign(ctx, projectID, apiKeyID, &matlas.AssignAPIKey{
					Roles: roles,
				})
				if err != nil {
					return diag.Errorf("error assigning api_keys into the project(%s): %s", projectID, err)
				}
			}
		}

		// Removing api_keys from the project
		for _, apiKey := range removedAPIKeys {
			projectID := apiKey.(map[string]interface{})["project_id"].(string)
			_, err := conn.ProjectAPIKeys.Unassign(ctx, projectID, apiKeyID)
			if err != nil {
				return diag.Errorf("error removing api_key(%s) from the project(%s): %s", apiKeyID, projectID, err)
			}
		}

		// Updating the role names for the api_key
		for _, apiKey := range changedAPIKeys {
			projectID := apiKey.(map[string]interface{})["project_id"].(string)
			roles := expandStringList(apiKey.(map[string]interface{})["role_names"].(*schema.Set).List())
			_, err := conn.ProjectAPIKeys.Assign(ctx, projectID, apiKeyID, &matlas.AssignAPIKey{
				Roles: roles,
			})
			if err != nil {
				return diag.Errorf("error updating role names for the api_key(%s): %s", apiKey, err)
			}
		}
	}
	return resourceMongoDBAtlasProjectAPIKeyRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectAPIKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	apiKeyID := ids["api_key_id"]
	var orgID string

	projectAPIKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
	}

	for _, val := range projectAPIKeys {
		if val.ID == apiKeyID {
			for i, role := range val.Roles {
				if strings.HasPrefix(role.RoleName, "ORG_") {
					orgID = val.Roles[i].OrgID
				}
			}
		}
	}

	_, roleOk := d.GetOk("role_names")
	if !roleOk {
		options := &matlas.ListOptions{}

		apiKeyOrgList, _, err := conn.Root.List(ctx, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
		}

		projectAssignments := []APIProjectAssignmentKeyInput{}

		for idx, role := range apiKeyOrgList.APIKey.Roles {
			if strings.HasPrefix(role.RoleName, "ORG_") {
				orgKeys, _, err := conn.APIKeys.List(ctx, apiKeyOrgList.APIKey.Roles[idx].OrgID, options)
				if err != nil {
					return diag.FromErr(fmt.Errorf("error getting api key information: %s", err))
				}

				for _, val := range orgKeys {
					if val.ID == apiKeyID {
						for _, r := range val.Roles {
							temp := new(APIProjectAssignmentKeyInput)
							if strings.HasPrefix(r.RoleName, "GROUP_") {
								temp.ProjectID = r.GroupID
								for _, l := range val.Roles {
									if l.GroupID == temp.ProjectID {
										temp.RoleNames = append(temp.RoleNames, l.RoleName)
									}
								}
								projectAssignments = append(projectAssignments, *temp)
							}
						}
					}
				}
				break
			}
		}
		for _, apiKey := range projectAssignments {
			_, err = conn.ProjectAPIKeys.Unassign(ctx, apiKey.ProjectID, apiKeyID)
			if err != nil {
				return diag.FromErr(fmt.Errorf("error deleting project api key: %s", err))
			}
		}
	} else {
		_, err = conn.ProjectAPIKeys.Unassign(ctx, projectID, apiKeyID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting project api key: %s", err))
		}
	}

	_, err = conn.APIKeys.Delete(ctx, orgID, apiKeyID)
	if err != nil {
		log.Printf("[WARN] unable to delete Key (%s): %s\n", apiKeyID, err)
	}

	d.SetId("")
	return nil
}

func resourceMongoDBAtlasProjectAPIKeyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a api key use the format {org_id}-{api_key_id}")
	}

	projectID := parts[0]
	apiKeyID := parts[1]

	projectAPIKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't import api key %s in project %s, error: %s", projectID, apiKeyID, err)
	}
	for _, val := range projectAPIKeys {
		if val.ID == apiKeyID {
			if err := d.Set("description", val.Desc); err != nil {
				return nil, fmt.Errorf("error setting `description`: %s", err)
			}

			if err := d.Set("public_key", val.PublicKey); err != nil {
				return nil, fmt.Errorf("error setting `public_key`: %s", err)
			}

			d.SetId(encodeStateID(map[string]string{
				"project_id": projectID,
				"api_key_id": val.ID,
			}))
		}
	}
	return []*schema.ResourceData{d}, nil
}

func flattenProjectAPIKeys(ctx context.Context, conn *matlas.Client, projectID string, apiKeys []matlas.APIKey) []map[string]interface{} {
	var results []map[string]interface{}

	if len(apiKeys) > 0 {
		results = make([]map[string]interface{}, len(apiKeys))
		for k, apiKey := range apiKeys {
			results[k] = map[string]interface{}{
				"api_key_id":  apiKey.ID,
				"description": apiKey.Desc,
				"public_key":  apiKey.PublicKey,
				"private_key": apiKey.PrivateKey,
				"role_names":  flattenProjectAPIKeyRoles(projectID, apiKey.Roles),
			}
		}
	}
	return results
}

func flattenProjectAPIKeyRoles(projectID string, apiKeyRoles []matlas.AtlasRole) []string {
	if len(apiKeyRoles) == 0 {
		return nil
	}

	flattenedOrgRoles := []string{}

	for _, role := range apiKeyRoles {
		if strings.HasPrefix(role.RoleName, "GROUP_") && role.GroupID == projectID {
			flattenedOrgRoles = append(flattenedOrgRoles, role.RoleName)
		}
	}

	return flattenedOrgRoles
}

func ExpandProjectAssignmentSet(projectAssignments *schema.Set) []*APIProjectAssignmentKeyInput {
	res := make([]*APIProjectAssignmentKeyInput, projectAssignments.Len())

	for i, value := range projectAssignments.List() {
		v := value.(map[string]interface{})
		res[i] = &APIProjectAssignmentKeyInput{
			ProjectID: v["project_id"].(string),
			RoleNames: expandStringList(v["role_names"].(*schema.Set).List()),
		}
	}

	return res
}

func FlattenProjectAssignment(apiKeyAssignmentSet []APIProjectAssignmentKeyInput) []map[string]interface{} {
	var results []map[string]interface{}
	var atlasRoles []matlas.AtlasRole
	var atlasRole matlas.AtlasRole
	if len(apiKeyAssignmentSet) > 0 {
		results = make([]map[string]interface{}, len(apiKeyAssignmentSet))
		for k, apiKey := range apiKeyAssignmentSet {
			for _, roleName := range apiKey.RoleNames {
				atlasRole.GroupID = apiKey.ProjectID
				atlasRole.RoleName = roleName
				atlasRoles = append(atlasRoles, atlasRole)
			}
			results[k] = map[string]interface{}{
				"project_id": apiKey.ProjectID,
				"role_names": flattenProjectAPIKeyRoles(apiKey.ProjectID, atlasRoles),
			}
		}
	}
	return results
}

func getStateProjectAssignmentAPIKeys(d *schema.ResourceData) (newAPIKeys, changedAPIKeys, removedAPIKeys []interface{}) {
	currentAPIKeys, changes := d.GetChange("project_assignment")

	rAPIKeys := currentAPIKeys.(*schema.Set).Difference(changes.(*schema.Set))
	nAPIKeys := changes.(*schema.Set).Difference(currentAPIKeys.(*schema.Set))
	changedAPIKeys = make([]interface{}, 0)

	for _, changed := range nAPIKeys.List() {
		for _, removed := range rAPIKeys.List() {
			if changed.(map[string]interface{})["project_id"] == removed.(map[string]interface{})["project_id"] {
				rAPIKeys.Remove(removed)
			}
		}

		for _, current := range currentAPIKeys.(*schema.Set).List() {
			if changed.(map[string]interface{})["project_id"] == current.(map[string]interface{})["project_id"] {
				changedAPIKeys = append(changedAPIKeys, changed.(map[string]interface{}))
				nAPIKeys.Remove(changed)
			}
		}
	}

	newAPIKeys = nAPIKeys.List()
	removedAPIKeys = rAPIKeys.List()

	return
}
