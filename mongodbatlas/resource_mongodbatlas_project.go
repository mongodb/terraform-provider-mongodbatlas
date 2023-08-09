package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	"go.mongodb.org/atlas-sdk/v20230201002/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	errorProjectCreate  = "error creating Project: %s"
	errorProjectRead    = "error getting project(%s): %s"
	errorProjectDelete  = "error deleting project (%s): %s"
	errorProjectSetting = "error setting `%s` for project (%s): %s"
)

func resourceMongoDBAtlasProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasProjectCreate,
		ReadContext:   resourceMongoDBAtlasProjectRead,
		UpdateContext: resourceMongoDBAtlasProjectUpdate,
		DeleteContext: resourceMongoDBAtlasProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"teams": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"team_id": {
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
			"project_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"with_default_alerts_settings": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"api_keys": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: fmt.Sprintf(DeprecationMessageParameterToResource, "v1.12.0", "mongodbatlas_project_api_key"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key_id": {
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
			"is_collect_database_specifics_statistics_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_data_explorer_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_extended_storage_sizes_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_performance_advisor_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_realtime_performance_panel_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"is_schema_advisor_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"region_usage_restrictions": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"limits": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"current_usage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"maximum_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Resources that need to be cleaned up before a project can be deleted
type AtlastProjectDependents struct {
	AdvancedClusters *matlas.AdvancedClustersResponse
}

func resourceMongoDBAtlasProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	connV2 := meta.(*MongoDBClient).AtlasV2

	projectReq := &matlas.Project{
		OrgID:                     d.Get("org_id").(string),
		Name:                      d.Get("name").(string),
		WithDefaultAlertsSettings: pointy.Bool(d.Get("with_default_alerts_settings").(bool)),
		RegionUsageRestrictions:   d.Get("region_usage_restrictions").(string),
	}

	var createProjectOptions *matlas.CreateProjectOptions

	if projectOwnerID, ok := d.GetOk("project_owner_id"); ok {
		createProjectOptions = &matlas.CreateProjectOptions{
			ProjectOwnerID: projectOwnerID.(string),
		}
	}

	project, _, err := conn.Projects.Create(ctx, projectReq, createProjectOptions)
	if err != nil {
		return diag.Errorf(errorProjectCreate, err)
	}

	// Check if teams were set, if so we need to add the teams into the project
	if teams, ok := d.GetOk("teams"); ok {
		// adding the teams into the project
		_, _, err := conn.Projects.AddTeamsToProject(ctx, project.ID, expandTeamsSet(teams.(*schema.Set)))
		if err != nil {
			errd := deleteProject(ctx, meta, project.ID)
			if errd != nil {
				return diag.Errorf(errorProjectDelete, project.ID, err)
			}
			return diag.Errorf("error adding teams into the project: %s", err)
		}
	}

	// Check if api keys were set, if so we need to add keys into the project
	if apiKeys, ok := d.GetOk("api_keys"); ok {
		// assign api keys to the project
		for _, apiKey := range expandAPIKeysSet(apiKeys.(*schema.Set)) {
			_, err := conn.ProjectAPIKeys.Assign(ctx, project.ID, apiKey.id, &matlas.AssignAPIKey{
				Roles: apiKey.roles,
			})
			if err != nil {
				errd := deleteProject(ctx, meta, project.ID)
				if errd != nil {
					return diag.Errorf(errorProjectDelete, project.ID, err)
				}
				return diag.Errorf("error assigning api keys to the project: %s", err)
			}
		}
	}

	// Check if limits were set, if so we need to add into the project
	if limits, ok := d.GetOk("limits"); ok {
		// assign limits to the project
		for _, limit := range expandLimitsSet(limits.(*schema.Set)) {
			limitModel := &admin.DataFederationLimit{
				Name:  limit.name,
				Value: limit.value,
			}
			_, _, err := connV2.ProjectsApi.SetProjectLimit(ctx, limit.name, project.ID, limitModel).Execute()

			if err != nil {
				errd := deleteProject(ctx, meta, project.ID)
				if errd != nil {
					return diag.Errorf(errorProjectDelete, project.ID, err)
				}
				return diag.Errorf("error assigning limits to the project: %s", err)
			}
		}
	}

	projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, project.ID)
	if err != nil {
		errd := deleteProject(ctx, meta, project.ID)
		if errd != nil {
			return diag.Errorf(errorProjectDelete, project.ID, err)
		}
		return diag.Errorf("error getting project's settings assigned (%s): %s", project.ID, err)
	}

	if v, ok := d.GetOkExists("is_collect_database_specifics_statistics_enabled"); ok {
		projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled = pointy.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_data_explorer_enabled"); ok {
		projectSettings.IsDataExplorerEnabled = pointy.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_extended_storage_sizes_enabled"); ok {
		projectSettings.IsExtendedStorageSizesEnabled = pointy.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_performance_advisor_enabled"); ok {
		projectSettings.IsPerformanceAdvisorEnabled = pointy.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_realtime_performance_panel_enabled"); ok {
		projectSettings.IsRealtimePerformancePanelEnabled = pointy.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("is_schema_advisor_enabled"); ok {
		projectSettings.IsSchemaAdvisorEnabled = pointy.Bool(v.(bool))
	}

	_, _, err = conn.Projects.UpdateProjectSettings(ctx, project.ID, projectSettings)
	if err != nil {
		errd := deleteProject(ctx, meta, project.ID)
		if errd != nil {
			return diag.Errorf(errorProjectDelete, project.ID, err)
		}
		return diag.Errorf("error updating project's settings assigned (%s): %s", project.ID, err)
	}

	d.SetId(project.ID)

	return resourceMongoDBAtlasProjectRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	connV2 := meta.(*MongoDBClient).AtlasV2
	projectID := d.Id()

	projectRes, resp, err := conn.Projects.GetOneProject(context.Background(), projectID)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}

		return diag.Errorf(errorProjectRead, projectID, err)
	}

	teams, _, err := conn.Projects.GetProjectTeamsAssigned(ctx, projectID)
	if err != nil {
		return diag.Errorf("error getting project's teams assigned (%s): %s", projectID, err)
	}

	apiKeys, err := getProjectAPIKeys(ctx, conn, projectRes.OrgID, projectRes.ID)
	if err != nil {
		var target *matlas.ErrorResponse
		if errors.As(err, &target) && target.ErrorCode != "USER_UNAUTHORIZED" {
			return diag.Errorf("error getting project's api keys (%s): %s", projectID, err)
		}
		log.Println("[WARN] `api_keys` will be empty because the user has no permissions to read the api keys endpoint")
	}

	if definedLimits, ok := d.GetOk("limits"); ok {
		definedLimitsList := expandLimitsSet(definedLimits.(*schema.Set))

		// terraform state will only save user defined limits
		filteredLimits, err := fetchUserDefinedLimits(ctx, definedLimitsList, projectID, connV2)
		if err != nil {
			return err
		}

		if err := d.Set("limits", flattenLimits(filteredLimits)); err != nil {
			return diag.Errorf(errorProjectSetting, `limits`, projectID, err)
		}
	}

	projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, projectID)
	if err != nil {
		return diag.Errorf("error getting project's settings assigned (%s): %s", projectID, err)
	}

	if err := d.Set("name", projectRes.Name); err != nil {
		return diag.Errorf(errorProjectSetting, `name`, projectID, err)
	}

	if err := d.Set("org_id", projectRes.OrgID); err != nil {
		return diag.Errorf(errorProjectSetting, `org_id`, projectID, err)
	}

	if err := d.Set("cluster_count", projectRes.ClusterCount); err != nil {
		return diag.Errorf(errorProjectSetting, `clusterCount`, projectID, err)
	}

	if err := d.Set("created", projectRes.Created); err != nil {
		return diag.Errorf(errorProjectSetting, `created`, projectID, err)
	}

	if err := d.Set("teams", flattenTeams(teams)); err != nil {
		return diag.Errorf(errorProjectSetting, `teams`, projectID, err)
	}

	if err := d.Set("api_keys", flattenAPIKeys(apiKeys)); err != nil {
		return diag.Errorf(errorProjectSetting, `api_keys`, projectID, err)
	}

	if err := d.Set("is_collect_database_specifics_statistics_enabled", projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled); err != nil {
		return diag.Errorf(errorProjectSetting, `is_collect_database_specifics_statistics_enabled`, projectID, err)
	}
	if err := d.Set("is_data_explorer_enabled", projectSettings.IsDataExplorerEnabled); err != nil {
		return diag.Errorf(errorProjectSetting, `is_data_explorer_enabled`, projectID, err)
	}
	if err := d.Set("is_extended_storage_sizes_enabled", projectSettings.IsExtendedStorageSizesEnabled); err != nil {
		return diag.Errorf(errorProjectSetting, `is_extended_storage_sizes_enabled`, projectID, err)
	}
	if err := d.Set("is_performance_advisor_enabled", projectSettings.IsPerformanceAdvisorEnabled); err != nil {
		return diag.Errorf(errorProjectSetting, `is_performance_advisor_enabled`, projectID, err)
	}
	if err := d.Set("is_realtime_performance_panel_enabled", projectSettings.IsRealtimePerformancePanelEnabled); err != nil {
		return diag.Errorf(errorProjectSetting, `is_realtime_performance_panel_enabled`, projectID, err)
	}
	if err := d.Set("is_schema_advisor_enabled", projectSettings.IsSchemaAdvisorEnabled); err != nil {
		return diag.Errorf(errorProjectSetting, `is_schema_advisor_enabled`, projectID, err)
	}

	return nil
}

func fetchUserDefinedLimits(ctx context.Context, definedLimitsList []*projectLimit, projectID string, connV2 *admin.APIClient) ([]admin.DataFederationLimit, diag.Diagnostics) {
	definedLimitsMap := make(map[string]*projectLimit)
	for _, definedLimit := range definedLimitsList {
		definedLimitsMap[definedLimit.name] = definedLimit
	}

	fetchedLimits, _, err := connV2.ProjectsApi.ListProjectLimits(ctx, projectID).Execute()
	if err != nil {
		return nil, diag.Errorf("error getting project's limits (%s): %s", projectID, err)
	}
	filteredLimits := []admin.DataFederationLimit{}
	for i := range fetchedLimits {
		limitRes := fetchedLimits[i]
		if _, ok := definedLimitsMap[limitRes.Name]; ok {
			filteredLimits = append(filteredLimits, limitRes)
		}
	}
	return filteredLimits, nil
}

func resourceMongoDBAtlasProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := updateProject(ctx, d, meta); err != nil {
		return err
	}
	if err := updateTeams(ctx, d, meta); err != nil {
		return err
	}
	if err := updateAPIKeys(ctx, d, meta); err != nil {
		return err
	}
	if err := updateLimits(ctx, d, meta); err != nil {
		return err
	}

	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()
	projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, projectID)
	if err != nil {
		return diag.Errorf("error getting project's settings assigned (%s): %s", projectID, err)
	}

	if d.HasChange("is_collect_database_specifics_statistics_enabled") {
		projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled = pointy.Bool(d.Get("is_collect_database_specifics_statistics_enabled").(bool))
	}

	if d.HasChange("is_data_explorer_enabled") {
		projectSettings.IsDataExplorerEnabled = pointy.Bool(d.Get("is_data_explorer_enabled").(bool))
	}
	if d.HasChange("is_extended_storage_sizes_enabled") {
		projectSettings.IsExtendedStorageSizesEnabled = pointy.Bool(d.Get("is_extended_storage_sizes_enabled").(bool))
	}
	if d.HasChange("is_performance_advisor_enabled") {
		projectSettings.IsPerformanceAdvisorEnabled = pointy.Bool(d.Get("is_performance_advisor_enabled").(bool))
	}
	if d.HasChange("is_realtime_performance_panel_enabled") {
		projectSettings.IsRealtimePerformancePanelEnabled = pointy.Bool(d.Get("is_realtime_performance_panel_enabled").(bool))
	}
	if d.HasChange("is_schema_advisor_enabled") {
		projectSettings.IsSchemaAdvisorEnabled = pointy.Bool(d.Get("is_schema_advisor_enabled").(bool))
	}
	if d.HasChange("is_collect_database_specifics_statistics_enabled") || d.HasChange("is_data_explorer_enabled") ||
		d.HasChange("is_performance_advisor_enabled") || d.HasChange("is_realtime_performance_panel_enabled") ||
		d.HasChange("is_schema_advisor_enabled") || d.HasChange("is_extended_storage_sizes_enabled") {
		_, _, err := conn.Projects.UpdateProjectSettings(ctx, projectID, projectSettings)
		if err != nil {
			return diag.Errorf("error updating project's settings assigned (%s): %s", projectID, err)
		}
	}
	return resourceMongoDBAtlasProjectRead(ctx, d, meta)
}

func resourceMongoDBAtlasProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Id()
	return deleteProject(ctx, meta, projectID)
}

/*
This assumes the project CRUD outcome will be the same for any non-zero number of dependents

If all dependents are deleting, wait to try and delete
Else consider the aggregate dependents idle.

If we get a defined error response, return that right away
Else retry
*/
func resourceProjectDependentsDeletingRefreshFunc(ctx context.Context, projectID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var target *matlas.ErrorResponse
		clusters, _, err := client.AdvancedClusters.List(ctx, projectID, nil)
		dependents := AtlastProjectDependents{AdvancedClusters: clusters}

		if errors.As(err, &target) {
			return nil, "", err
		} else if err != nil {
			return nil, "RETRY", nil
		}

		if dependents.AdvancedClusters.TotalCount == 0 {
			return dependents, "IDLE", nil
		}

		for _, v := range dependents.AdvancedClusters.Results {
			if v.StateName != "DELETING" {
				return dependents, "IDLE", nil
			}
		}

		log.Printf("[DEBUG] status for MongoDB project %s dependents: %s", projectID, "DELETING")

		return dependents, "DELETING", nil
	}
}

func expandTeamsSet(teams *schema.Set) []*matlas.ProjectTeam {
	res := make([]*matlas.ProjectTeam, teams.Len())

	for i, value := range teams.List() {
		v := value.(map[string]interface{})
		res[i] = &matlas.ProjectTeam{
			TeamID:    v["team_id"].(string),
			RoleNames: expandStringList(v["role_names"].(*schema.Set).List()),
		}
	}

	return res
}

func expandAPIKeysSet(apiKeys *schema.Set) []*apiKey {
	res := make([]*apiKey, apiKeys.Len())

	for i, value := range apiKeys.List() {
		v := value.(map[string]interface{})
		res[i] = &apiKey{
			id:    v["api_key_id"].(string),
			roles: expandStringList(v["role_names"].(*schema.Set).List()),
		}
	}

	return res
}

type projectLimit struct {
	name  string
	value int64
}

func expandLimitsSet(limits *schema.Set) []*projectLimit {
	res := make([]*projectLimit, limits.Len())

	for i, value := range limits.List() {
		v := value.(map[string]interface{})
		res[i] = &projectLimit{
			name:  v["name"].(string),
			value: int64(v["value"].(int)),
		}
	}

	return res
}

func expandTeamsList(teams []interface{}) []*matlas.ProjectTeam {
	res := make([]*matlas.ProjectTeam, len(teams))

	for i, value := range teams {
		v := value.(map[string]interface{})
		res[i] = &matlas.ProjectTeam{
			TeamID:    v["team_id"].(string),
			RoleNames: expandStringList(v["role_names"].(*schema.Set).List()),
		}
	}

	return res
}

func expandAPIKeysList(apiKeys []interface{}) []*apiKey {
	res := make([]*apiKey, len(apiKeys))

	for i, value := range apiKeys {
		v := value.(map[string]interface{})
		res[i] = &apiKey{
			id:    v["api_key_id"].(string),
			roles: expandStringList(v["role_names"].(*schema.Set).List()),
		}
	}

	return res
}

func expandLimitsList(limits []interface{}) []*projectLimit {
	res := make([]*projectLimit, len(limits))

	for i, value := range limits {
		v := value.(map[string]interface{})
		res[i] = &projectLimit{
			name:  v["name"].(string),
			value: int64(v["value"].(int)),
		}
	}

	return res
}

func flattenTeams(ta *matlas.TeamsAssigned) []map[string]interface{} {
	if ta == nil || ta.Results == nil || ta.TotalCount == 0 {
		return nil
	}

	teams := ta.Results
	res := make([]map[string]interface{}, len(teams))

	for i, team := range teams {
		res[i] = map[string]interface{}{
			"team_id":    team.TeamID,
			"role_names": team.RoleNames,
		}
	}

	return res
}

func flattenAPIKeys(keys []*apiKey) []map[string]interface{} {
	res := make([]map[string]interface{}, len(keys))

	for i, key := range keys {
		res[i] = map[string]interface{}{
			"api_key_id": key.id,
			"role_names": key.roles,
		}
	}

	return res
}

func flattenLimits(limits []admin.DataFederationLimit) []map[string]interface{} {
	res := make([]map[string]interface{}, len(limits))

	for i, limit := range limits {
		res[i] = map[string]interface{}{
			"name":  limit.Name,
			"value": limit.Value,
		}
		if limit.CurrentUsage != nil {
			res[i]["current_usage"] = *limit.CurrentUsage
		}
		if limit.DefaultLimit != nil {
			res[i]["default_limit"] = *limit.DefaultLimit
		}
		if limit.MaximumLimit != nil {
			res[i]["maximum_limit"] = *limit.MaximumLimit
		}
	}

	return res
}

// getChangesInSet divides modified elements of a Set into 3 distinct groups using an a specific key for comparing elements. This is useful for update in Set values.
// The function receives an attribute key where the Set is stored, and a elementsIdKey to uniquely identify each element.
// - newElements: new elements that where not present in previous state
// - changedElements: elements where some value was modified but is still present
// - removedElements: elements that are no longer present
func getChangesInSet(d *schema.ResourceData, attributeKey, elementsIDKey string) (newElements, changedElements, removedElements []interface{}) {
	oldSet, newSet := d.GetChange(attributeKey)

	removedSchemaSet := oldSet.(*schema.Set).Difference(newSet.(*schema.Set))
	newSchemaSet := newSet.(*schema.Set).Difference(oldSet.(*schema.Set))
	changedElements = make([]interface{}, 0)

	for _, new := range newSchemaSet.List() {
		for _, removed := range removedSchemaSet.List() {
			if new.(map[string]interface{})[elementsIDKey] == removed.(map[string]interface{})[elementsIDKey] {
				removedSchemaSet.Remove(removed)
			}
		}

		for _, old := range oldSet.(*schema.Set).List() {
			if new.(map[string]interface{})[elementsIDKey] == old.(map[string]interface{})[elementsIDKey] {
				changedElements = append(changedElements, new.(map[string]interface{}))
				newSchemaSet.Remove(new)
			}
		}
	}

	newElements = newSchemaSet.List()
	removedElements = removedSchemaSet.List()

	return
}

func deleteProject(ctx context.Context, meta interface{}, projectID string) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING", "RETRY"},
		Target:     []string{"IDLE"},
		Refresh:    resourceProjectDependentsDeletingRefreshFunc(ctx, projectID, conn),
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      0,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		log.Printf("[ERROR] could not determine MongoDB project %s dependents status: %s", projectID, err.Error())
	}

	_, err = conn.Projects.Delete(ctx, projectID)

	if err != nil {
		return diag.Errorf(errorProjectDelete, projectID, err)
	}

	return nil
}

func updateProject(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("name") {
		return nil
	}

	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()

	if _, _, err := conn.Projects.Update(ctx, projectID, newProjectUpdateRequest(d)); err != nil {
		return diag.Errorf("error updating the project(%s): %s", projectID, err)
	}

	return nil
}

func newProjectUpdateRequest(d *schema.ResourceData) *matlas.ProjectUpdateRequest {
	return &matlas.ProjectUpdateRequest{
		Name: d.Get("name").(string),
	}
}

func updateTeams(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("teams") {
		return nil
	}

	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()

	// get the current teams and the new teams with changes
	newTeams, changedTeams, removedTeams := getChangesInSet(d, "teams", "team_id")

	// adding new teams into the project
	if len(newTeams) > 0 {
		_, _, err := conn.Projects.AddTeamsToProject(ctx, projectID, expandTeamsList(newTeams))
		if err != nil {
			return diag.Errorf("error adding teams into the project(%s): %s", projectID, err)
		}
	}

	// Removing teams from the project
	for _, team := range removedTeams {
		teamID := team.(map[string]interface{})["team_id"].(string)

		_, err := conn.Teams.RemoveTeamFromProject(ctx, projectID, teamID)
		if err != nil {
			var target *matlas.ErrorResponse
			if errors.As(err, &target) && target.ErrorCode != "USER_UNAUTHORIZED" {
				return diag.Errorf("error removing team(%s) from the project(%s): %s", teamID, projectID, err)
			}
			log.Printf("[WARN] error removing team(%s) from the project(%s): %s", teamID, projectID, err)
		}
	}

	// Updating the role names for a team
	for _, t := range changedTeams {
		team := t.(map[string]interface{})

		_, _, err := conn.Teams.UpdateTeamRoles(ctx, projectID, team["team_id"].(string),
			&matlas.TeamUpdateRoles{
				RoleNames: expandStringList(team["role_names"].(*schema.Set).List()),
			},
		)
		if err != nil {
			return diag.Errorf("error updating role names for the team(%s): %s", team["team_id"], err)
		}
	}

	return nil
}

func updateAPIKeys(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("api_keys") {
		return nil
	}

	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Id()

	// get the current api_keys and the new api_keys with changes
	newAPIKeys, changedAPIKeys, removedAPIKeys := getChangesInSet(d, "api_keys", "api_key_id")

	// adding new api_keys into the project
	if len(newAPIKeys) > 0 {
		for _, apiKey := range expandAPIKeysList(newAPIKeys) {
			_, err := conn.ProjectAPIKeys.Assign(ctx, projectID, apiKey.id, &matlas.AssignAPIKey{
				Roles: apiKey.roles,
			})
			if err != nil {
				return diag.Errorf("error assigning api_keys into the project(%s): %s", projectID, err)
			}
		}
	}

	// Removing api_keys from the project
	for _, apiKey := range removedAPIKeys {
		apiKeyID := apiKey.(map[string]interface{})["api_key_id"].(string)
		_, err := conn.ProjectAPIKeys.Unassign(ctx, projectID, apiKeyID)
		if err != nil {
			return diag.Errorf("error removing api_key(%s) from the project(%s): %s", apiKeyID, projectID, err)
		}
	}

	// Updating the role names for the api_key
	for _, apiKey := range expandAPIKeysList(changedAPIKeys) {
		_, err := conn.ProjectAPIKeys.Assign(ctx, projectID, apiKey.id, &matlas.AssignAPIKey{
			Roles: apiKey.roles,
		})
		if err != nil {
			return diag.Errorf("error updating role names for the api_key(%s): %s", apiKey, err)
		}
	}

	return nil
}

func updateLimits(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if !d.HasChange("limits") {
		return nil
	}

	connV2 := meta.(*MongoDBClient).AtlasV2
	projectID := d.Id()

	// get the current limits and the new limits with changes
	newLimits, changedLimits, removedLimits := getChangesInSet(d, "limits", "name")

	// adding new limits into the project
	if len(newLimits) > 0 {
		for _, limit := range expandLimitsList(newLimits) {
			if err := setLimitToProject(ctx, limit, projectID, connV2); err != nil {
				return err
			}
		}
	}

	// Removing limits from the project
	for _, limit := range removedLimits {
		limitName := limit.(map[string]interface{})["name"].(string)
		_, _, err := connV2.ProjectsApi.DeleteProjectLimit(ctx, limitName, projectID).Execute()
		if err != nil {
			return diag.Errorf("error removing limit %s from the project(%s): %s", limitName, projectID, err)
		}
	}

	// Updating values for changed limits
	for _, limit := range expandLimitsList(changedLimits) {
		if err := setLimitToProject(ctx, limit, projectID, connV2); err != nil {
			return err
		}
	}

	return nil
}

func setLimitToProject(ctx context.Context, limit *projectLimit, projectID string, connV2 *admin.APIClient) diag.Diagnostics {
	limitModel := &admin.DataFederationLimit{
		Name:  limit.name,
		Value: limit.value,
	}
	_, _, err := connV2.ProjectsApi.SetProjectLimit(ctx, limit.name, projectID, limitModel).Execute()
	if err != nil {
		return diag.Errorf("error assigning limit %s to the project: %s", limitModel.Name, err)
	}
	return nil
}
