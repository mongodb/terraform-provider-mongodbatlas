package project

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sort"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115002/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	projectResourceName            = "project"
	errorProjectCreate             = "error creating Project: %s"
	errorProjectDelete             = "error deleting project (%s): %s"
	errorProjectUpdate             = "error updating project (%s): %s"
	projectDependentsStateIdle     = "IDLE"
	projectDependentsStateDeleting = "DELETING"
	projectDependentsStateRetry    = "RETRY"
)

var _ resource.ResourceWithConfigure = &projectRS{}
var _ resource.ResourceWithImportState = &projectRS{}

func Resource() resource.Resource {
	return &projectRS{
		RSCommon: config.RSCommon{
			ResourceName: projectResourceName,
		},
	}
}

type projectRS struct {
	config.RSCommon
}

type TfProjectRSModel struct {
	Limits                                      types.Set    `tfsdk:"limits"`
	Teams                                       types.Set    `tfsdk:"teams"`
	RegionUsageRestrictions                     types.String `tfsdk:"region_usage_restrictions"`
	Name                                        types.String `tfsdk:"name"`
	OrgID                                       types.String `tfsdk:"org_id"`
	Created                                     types.String `tfsdk:"created"`
	ProjectOwnerID                              types.String `tfsdk:"project_owner_id"`
	ID                                          types.String `tfsdk:"id"`
	ClusterCount                                types.Int64  `tfsdk:"cluster_count"`
	IsDataExplorerEnabled                       types.Bool   `tfsdk:"is_data_explorer_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool   `tfsdk:"is_performance_advisor_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool   `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool   `tfsdk:"is_schema_advisor_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool   `tfsdk:"is_extended_storage_sizes_enabled"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool   `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	WithDefaultAlertsSettings                   types.Bool   `tfsdk:"with_default_alerts_settings"`
}

type TfTeamModel struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.Set    `tfsdk:"role_names"`
}

type TfLimitModel struct {
	Name         types.String `tfsdk:"name"`
	Value        types.Int64  `tfsdk:"value"`
	CurrentUsage types.Int64  `tfsdk:"current_usage"`
	DefaultLimit types.Int64  `tfsdk:"default_limit"`
	MaximumLimit types.Int64  `tfsdk:"maximum_limit"`
}

var TfTeamObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"team_id":    types.StringType,
	"role_names": types.SetType{ElemType: types.StringType},
}}
var TfLimitObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name":          types.StringType,
	"value":         types.Int64Type,
	"current_usage": types.Int64Type,
	"default_limit": types.Int64Type,
	"maximum_limit": types.Int64Type,
}}

// Resources that need to be cleaned up before a project can be deleted
type AtlasProjectDependants struct {
	AdvancedClusters *admin.PaginatedAdvancedClusterDescription
}

func (r *projectRS) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"org_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_count": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_owner_id": schema.StringAttribute{
				Optional: true,
			},
			"with_default_alerts_settings": schema.BoolAttribute{
				// Default values also must be Computed otherwise Terraform throws error:
				// Schema Using Attribute Default For Non-Computed Attribute
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"is_collect_database_specifics_statistics_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_data_explorer_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_extended_storage_sizes_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_performance_advisor_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_realtime_performance_panel_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_schema_advisor_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"region_usage_restrictions": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"teams": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"team_id": schema.StringAttribute{
							Required: true,
						},
						"role_names": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"limits": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"value": schema.Int64Attribute{
							Required: true,
						},
						"current_usage": schema.Int64Attribute{
							Computed: true,
						},
						"default_limit": schema.Int64Attribute{
							Computed: true,
						},
						"maximum_limit": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
				// https://discuss.hashicorp.com/t/computed-attributes-and-plan-modifiers/45830/12
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *projectRS) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var projectPlan TfProjectRSModel
	var teams []TfTeamModel
	var limits []TfLimitModel

	connV2 := r.Client.AtlasV2

	diags := req.Plan.Get(ctx, &projectPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectGroup := &admin.Group{
		OrgId:                     projectPlan.OrgID.ValueString(),
		Name:                      projectPlan.Name.ValueString(),
		WithDefaultAlertsSettings: projectPlan.WithDefaultAlertsSettings.ValueBoolPointer(),
		RegionUsageRestrictions:   projectPlan.RegionUsageRestrictions.ValueStringPointer(),
	}

	projectAPIParams := &admin.CreateProjectApiParams{
		Group:          projectGroup,
		ProjectOwnerId: conversion.StringNullIfEmpty(projectPlan.ProjectOwnerID.ValueString()).ValueStringPointer(),
	}

	// create project
	project, _, err := connV2.ProjectsApi.CreateProjectWithParams(ctx, projectAPIParams).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorProjectCreate, err.Error())
		return
	}

	// add teams
	if len(projectPlan.Teams.Elements()) > 0 {
		_ = projectPlan.Teams.ElementsAs(ctx, &teams, false)

		_, _, err := connV2.TeamsApi.AddAllTeamsToProject(ctx, project.GetId(), NewTeamRoleList(ctx, teams)).Execute()
		if err != nil {
			errd := deleteProject(ctx, r.Client.AtlasV2, project.Id)
			if errd != nil {
				resp.Diagnostics.AddError("error during project deletion when adding teams", fmt.Sprintf(errorProjectDelete, project.GetId(), err.Error()))
				return
			}
			resp.Diagnostics.AddError("error adding teams into the project", err.Error())
			return
		}
	}

	// add limits
	if len(projectPlan.Limits.Elements()) > 0 {
		_ = projectPlan.Limits.ElementsAs(ctx, &limits, false)

		for _, limit := range limits {
			dataFederationLimit := &admin.DataFederationLimit{
				Name:  limit.Name.ValueString(),
				Value: limit.Value.ValueInt64(),
			}
			_, _, err := connV2.ProjectsApi.SetProjectLimit(ctx, limit.Name.ValueString(), project.GetId(), dataFederationLimit).Execute()
			if err != nil {
				errd := deleteProject(ctx, r.Client.AtlasV2, project.Id)
				if errd != nil {
					resp.Diagnostics.AddError("error during project deletion when adding limits", fmt.Sprintf(errorProjectDelete, project.GetId(), err.Error()))
					return
				}
				resp.Diagnostics.AddError("error adding limits into the project", err.Error())
				return
			}
		}
	}

	// add settings
	projectSettings, _, err := connV2.ProjectsApi.GetProjectSettings(ctx, *project.Id).Execute()
	if err != nil {
		errd := deleteProject(ctx, r.Client.AtlasV2, project.Id)
		if errd != nil {
			resp.Diagnostics.AddError("error during project deletion when getting project settings", fmt.Sprintf(errorProjectDelete, project.GetId(), err.Error()))
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project's settings assigned (%s):", project.GetId()), err.Error())
		return
	}

	if !projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled.IsUnknown() {
		projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled = projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled.ValueBoolPointer()
	}
	if !projectPlan.IsDataExplorerEnabled.IsUnknown() {
		projectSettings.IsDataExplorerEnabled = projectPlan.IsDataExplorerEnabled.ValueBoolPointer()
	}
	if !projectPlan.IsExtendedStorageSizesEnabled.IsUnknown() {
		projectSettings.IsExtendedStorageSizesEnabled = projectPlan.IsExtendedStorageSizesEnabled.ValueBoolPointer()
	}
	if !projectPlan.IsPerformanceAdvisorEnabled.IsUnknown() {
		projectSettings.IsPerformanceAdvisorEnabled = projectPlan.IsPerformanceAdvisorEnabled.ValueBoolPointer()
	}
	if !projectPlan.IsRealtimePerformancePanelEnabled.IsUnknown() {
		projectSettings.IsRealtimePerformancePanelEnabled = projectPlan.IsRealtimePerformancePanelEnabled.ValueBoolPointer()
	}
	if !projectPlan.IsSchemaAdvisorEnabled.IsUnknown() {
		projectSettings.IsSchemaAdvisorEnabled = projectPlan.IsSchemaAdvisorEnabled.ValueBoolPointer()
	}

	if _, _, err = connV2.ProjectsApi.UpdateProjectSettings(ctx, project.GetId(), projectSettings).Execute(); err != nil {
		errd := deleteProject(ctx, r.Client.AtlasV2, project.Id)
		if errd != nil {
			resp.Diagnostics.AddError("error during project deletion when updating project settings", fmt.Sprintf(errorProjectDelete, project.GetId(), err.Error()))
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error updating project's settings assigned (%s):", project.GetId()), err.Error())
		return
	}

	projectID := project.GetId()
	projectRes, atlasResp, err := connV2.ProjectsApi.GetProject(ctx, projectID).Execute()
	if err != nil {
		if resp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting project after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	// get project props
	atlasTeams, atlasLimits, atlasProjectSettings, err := GetProjectPropsFromAPI(ctx, ServiceFromClient(connV2), projectID)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	atlasLimits = FilterUserDefinedLimits(atlasLimits, limits)
	projectPlanNew := NewTFProjectResourceModel(ctx, projectRes, atlasTeams, atlasProjectSettings, atlasLimits)
	updatePlanFromConfig(projectPlanNew, &projectPlan)

	// set state to fully populated data
	diags = resp.State.Set(ctx, projectPlanNew)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var projectState TfProjectRSModel
	var limits []TfLimitModel
	connV2 := r.Client.AtlasV2

	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &projectState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := projectState.ID.ValueString()
	if len(projectState.Limits.Elements()) > 0 {
		_ = projectState.Limits.ElementsAs(ctx, &limits, false)
	}

	// get project
	projectRes, atlasResp, err := connV2.ProjectsApi.GetProject(ctx, projectID).Execute()
	if err != nil {
		if resp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	// get project props
	atlasTeams, atlasLimits, atlasProjectSettings, err := GetProjectPropsFromAPI(ctx, ServiceFromClient(connV2), projectID)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	atlasLimits = FilterUserDefinedLimits(atlasLimits, limits)
	projectStateNew := NewTFProjectResourceModel(ctx, projectRes, atlasTeams, atlasProjectSettings, atlasLimits)
	updatePlanFromConfig(projectStateNew, &projectState)

	// save read data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectStateNew)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var projectState TfProjectRSModel
	var projectPlan TfProjectRSModel
	connV2 := r.Client.AtlasV2

	// get current state
	resp.Diagnostics.Append(req.State.Get(ctx, &projectState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// get current plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &projectPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := projectState.ID.ValueString()

	err := UpdateProject(ctx, ServiceFromClient(connV2), &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	err = UpdateProjectTeams(ctx, ServiceFromClient(connV2), &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project teams update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	err = UpdateProjectLimits(ctx, ServiceFromClient(connV2), &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project limits update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	err = updateProjectSettings(ctx, connV2, &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project settings update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	projectRes, atlasResp, err := connV2.ProjectsApi.GetProject(ctx, projectID).Execute()
	if err != nil {
		if resp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting project after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	// get project props
	atlasTeams, atlasLimits, atlasProjectSettings, err := GetProjectPropsFromAPI(ctx, ServiceFromClient(connV2), projectID)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}
	var planLimits []TfLimitModel
	_ = projectPlan.Limits.ElementsAs(ctx, &planLimits, false)
	atlasLimits = FilterUserDefinedLimits(atlasLimits, planLimits)
	projectPlanNew := NewTFProjectResourceModel(ctx, projectRes, atlasTeams, atlasProjectSettings, atlasLimits)
	updatePlanFromConfig(projectPlanNew, &projectPlan)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectPlanNew)...)
}

func (r *projectRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var project *TfProjectRSModel

	// read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := project.ID.ValueString()
	err := deleteProject(ctx, r.Client.AtlasV2, &projectID)

	if err != nil {
		resp.Diagnostics.AddError("error when destroying resource", fmt.Sprintf(errorProjectDelete, projectID, err.Error()))
		return
	}
}

func (r *projectRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func updatePlanFromConfig(projectPlanNewPtr, projectPlan *TfProjectRSModel) {
	// we need to reset defaults from what was previously in the state:
	// https://discuss.hashicorp.com/t/boolean-optional-default-value-migration-to-framework/55932
	projectPlanNewPtr.WithDefaultAlertsSettings = projectPlan.WithDefaultAlertsSettings
	projectPlanNewPtr.ProjectOwnerID = projectPlan.ProjectOwnerID
}

func FilterUserDefinedLimits(allAtlasLimits []admin.DataFederationLimit, tflimits []TfLimitModel) []admin.DataFederationLimit {
	filteredLimits := []admin.DataFederationLimit{}
	allLimitsMap := make(map[string]admin.DataFederationLimit)

	for _, limit := range allAtlasLimits {
		allLimitsMap[limit.Name] = limit
	}

	for _, definedTfLimit := range tflimits {
		if limit, ok := allLimitsMap[definedTfLimit.Name.ValueString()]; ok {
			filteredLimits = append(filteredLimits, limit)
		}
	}

	return filteredLimits
}

func GetProjectPropsFromAPI(ctx context.Context, client GroupProjectService, projectID string) (*admin.PaginatedTeamRole, []admin.DataFederationLimit, *admin.GroupSettings, error) {
	teams, _, err := client.ListProjectTeams(ctx, projectID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting project's teams assigned (%s): %v", projectID, err.Error())
	}

	limits, _, err := client.ListProjectLimits(ctx, projectID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting project's limits (%s): %s", projectID, err.Error())
	}

	projectSettings, _, err := client.GetProjectSettings(ctx, projectID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting project's settings assigned (%s): %v", projectID, err.Error())
	}

	return teams, limits, projectSettings, nil
}

func updateProjectSettings(ctx context.Context, connV2 *admin.APIClient, projectState, projectPlan *TfProjectRSModel) error {
	hasChanged := false
	projectID := projectState.ID.ValueString()
	projectSettings, _, err := connV2.ProjectsApi.GetProjectSettings(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("error getting project's settings assigned: %v", err.Error())
	}

	if projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled != projectState.IsCollectDatabaseSpecificsStatisticsEnabled {
		hasChanged = true
		projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled = projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled.ValueBoolPointer()
	}
	if projectPlan.IsDataExplorerEnabled != projectState.IsDataExplorerEnabled {
		hasChanged = true
		projectSettings.IsDataExplorerEnabled = projectPlan.IsDataExplorerEnabled.ValueBoolPointer()
	}
	if projectPlan.IsExtendedStorageSizesEnabled != projectState.IsExtendedStorageSizesEnabled {
		hasChanged = true
		projectSettings.IsExtendedStorageSizesEnabled = projectPlan.IsExtendedStorageSizesEnabled.ValueBoolPointer()
	}
	if projectPlan.IsPerformanceAdvisorEnabled != projectState.IsPerformanceAdvisorEnabled {
		hasChanged = true
		projectSettings.IsPerformanceAdvisorEnabled = projectPlan.IsPerformanceAdvisorEnabled.ValueBoolPointer()
	}
	if projectPlan.IsRealtimePerformancePanelEnabled != projectState.IsRealtimePerformancePanelEnabled {
		hasChanged = true
		projectSettings.IsRealtimePerformancePanelEnabled = projectPlan.IsRealtimePerformancePanelEnabled.ValueBoolPointer()
	}
	if projectPlan.IsSchemaAdvisorEnabled != projectState.IsSchemaAdvisorEnabled {
		hasChanged = true
		projectSettings.IsSchemaAdvisorEnabled = projectPlan.IsSchemaAdvisorEnabled.ValueBoolPointer()
	}

	if hasChanged {
		_, _, err = connV2.ProjectsApi.UpdateProjectSettings(ctx, projectID, projectSettings).Execute()
		if err != nil {
			return fmt.Errorf("error updating project's settings assigned: %v", err.Error())
		}
	}
	return nil
}

func UpdateProjectLimits(ctx context.Context, client GroupProjectService, projectState, projectPlan *TfProjectRSModel) error {
	var planLimits []TfLimitModel
	var stateLimits []TfLimitModel
	_ = projectPlan.Limits.ElementsAs(ctx, &planLimits, false)
	_ = projectState.Limits.ElementsAs(ctx, &stateLimits, false)

	if !hasLimitsChanged(planLimits, stateLimits) {
		return nil
	}

	projectID := projectState.ID.ValueString()
	newLimits, changedLimits, removedLimits := getChangesInLimitsSet(planLimits, stateLimits)

	// removing limits from the project
	for _, limit := range removedLimits {
		limitName := limit.Name.ValueString()
		if _, _, err := client.DeleteProjectLimit(ctx, limitName, projectID); err != nil {
			return fmt.Errorf("error removing limit %s from the project(%s) during update: %s", limitName, projectID, err)
		}
	}

	// updating values for changed limits
	if len(changedLimits) > 0 {
		if err := setProjectLimits(ctx, client, projectID, changedLimits); err != nil {
			return fmt.Errorf("error adding modified limits into the project during update: %v", err.Error())
		}
	}

	// adding new limits into the project
	if len(newLimits) > 0 {
		if err := setProjectLimits(ctx, client, projectID, newLimits); err != nil {
			return fmt.Errorf("error adding limits into the project during update: %v", err.Error())
		}
	}

	return nil
}

func setProjectLimits(ctx context.Context, client GroupProjectService, projectID string, tfLimits []TfLimitModel) error {
	for _, limit := range tfLimits {
		dataFederationLimit := &admin.DataFederationLimit{
			Name:  limit.Name.ValueString(),
			Value: limit.Value.ValueInt64(),
		}
		_, _, err := client.SetProjectLimit(ctx, limit.Name.ValueString(), projectID, dataFederationLimit)
		if err != nil {
			return fmt.Errorf("error adding limits into the project: %v", err.Error())
		}
	}
	return nil
}

func UpdateProjectTeams(ctx context.Context, client GroupProjectService, projectState, projectPlan *TfProjectRSModel) error {
	var planTeams []TfTeamModel
	var stateTeams []TfTeamModel
	_ = projectPlan.Teams.ElementsAs(ctx, &planTeams, false)
	_ = projectState.Teams.ElementsAs(ctx, &stateTeams, false)

	if !hasTeamsChanged(planTeams, stateTeams) {
		return nil
	}

	projectID := projectState.ID.ValueString()
	newTeams, changedTeams, removedTeams := getChangesInTeamsSet(planTeams, stateTeams)

	// removing teams from the project
	for _, team := range removedTeams {
		teamID := team.TeamID.ValueString()
		_, err := client.RemoveProjectTeam(ctx, projectID, team.TeamID.ValueString())
		if err != nil {
			apiError, ok := admin.AsError(err)
			if ok && *apiError.ErrorCode != "USER_UNAUTHORIZED" {
				return fmt.Errorf("error removing team(%s) from the project(%s): %s", teamID, projectID, err)
			}
			log.Printf("[WARN] error removing team(%s) from the project(%s): %s", teamID, projectID, err)
		}
	}

	// updating the role names for a team
	for _, team := range changedTeams {
		teamID := team.TeamID.ValueString()

		_, _, err := client.UpdateTeamRoles(ctx, projectID, teamID,
			&admin.TeamRole{
				RoleNames: conversion.TypesSetToString(ctx, team.RoleNames),
			},
		)
		if err != nil {
			return fmt.Errorf("error updating role names for the team(%s): %s", teamID, err.Error())
		}
	}

	// adding new teams into the project
	if _, _, err := client.AddAllTeamsToProject(ctx, projectID, NewTeamRoleList(ctx, newTeams)); err != nil {
		return fmt.Errorf("error adding teams to the project: %v", err.Error())
	}

	return nil
}

func hasTeamsChanged(planTeams, stateTeams []TfTeamModel) bool {
	sort.Slice(planTeams, func(i, j int) bool {
		return planTeams[i].TeamID.ValueString() < planTeams[j].TeamID.ValueString()
	})
	sort.Slice(stateTeams, func(i, j int) bool {
		return stateTeams[i].TeamID.ValueString() < stateTeams[j].TeamID.ValueString()
	})
	return !reflect.DeepEqual(planTeams, stateTeams)
}

func hasLimitsChanged(planLimits, stateLimits []TfLimitModel) bool {
	sort.Slice(planLimits, func(i, j int) bool {
		return planLimits[i].Name.ValueString() < planLimits[j].Name.ValueString()
	})
	sort.Slice(stateLimits, func(i, j int) bool {
		return stateLimits[i].Name.ValueString() < stateLimits[j].Name.ValueString()
	})
	return !reflect.DeepEqual(planLimits, stateLimits)
}

func UpdateProject(ctx context.Context, client GroupProjectService, projectState, projectPlan *TfProjectRSModel) error {
	if projectPlan.Name.Equal(projectState.Name) {
		return nil
	}

	projectID := projectState.ID.ValueString()

	if _, _, err := client.UpdateProject(ctx, projectID, NewGroupName(projectPlan)); err != nil {
		return fmt.Errorf("error updating the project(%s): %s", projectID, err)
	}

	return nil
}

func deleteProject(ctx context.Context, connV2 *admin.APIClient, projectID *string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{projectDependentsStateDeleting, projectDependentsStateRetry},
		Target:     []string{projectDependentsStateIdle},
		Refresh:    resourceProjectDependentsDeletingRefreshFunc(ctx, projectID, connV2),
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      0,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("[ERROR] could not determine MongoDB project %s dependents status: %s", *projectID, err.Error()))
	}

	_, _, err = connV2.ProjectsApi.DeleteProject(ctx, *projectID).Execute()

	return err
}

/*
resourceProjectDependentsDeletingRefreshFunc assumes the project CRUD outcome will be the same for any non-zero number of dependents

If all dependents are deleting, wait to try and delete
Else consider the aggregate dependents idle.

If we get a defined error response, return that right away
Else retry
*/
func resourceProjectDependentsDeletingRefreshFunc(ctx context.Context, projectID *string, connV2 *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		nonNullProjectID := conversion.StringPtrNullIfEmpty(projectID)
		clusters, _, err := connV2.ClustersApi.ListClusters(ctx, nonNullProjectID.String()).Execute()
		dependents := AtlasProjectDependants{AdvancedClusters: clusters}

		if _, ok := admin.AsError(err); ok {
			return nil, "", err
		}

		if err != nil {
			return nil, projectDependentsStateRetry, nil
		}

		if *dependents.AdvancedClusters.TotalCount == 0 {
			return dependents, projectDependentsStateIdle, nil
		}

		for i := range dependents.AdvancedClusters.Results {
			if *dependents.AdvancedClusters.Results[i].StateName != projectDependentsStateDeleting {
				return dependents, projectDependentsStateIdle, nil
			}
		}

		log.Printf("[DEBUG] status for MongoDB project %s dependents: %s", nonNullProjectID, projectDependentsStateDeleting)

		return dependents, projectDependentsStateDeleting, nil
	}
}

func getChangesInTeamsSet(planTeams, stateTeams []TfTeamModel) (newElements, changedElements, removedElements []TfTeamModel) {
	var removedTeams, newTeams, changedTeams []TfTeamModel

	planTeamsMap := NewTfTeamModelMap(planTeams)
	stateTeamsMap := NewTfTeamModelMap(stateTeams)

	for teamID, stateTeam := range stateTeamsMap {
		if plannedTeam, exists := planTeamsMap[teamID]; exists {
			if !reflect.DeepEqual(plannedTeam, stateTeam) {
				changedTeams = append(changedTeams, plannedTeam)
			}
		} else {
			removedTeams = append(removedTeams, stateTeam)
		}
	}

	for teamID, team := range planTeamsMap {
		if _, exists := stateTeamsMap[teamID]; !exists {
			newTeams = append(newTeams, team)
		}
	}
	return newTeams, changedTeams, removedTeams
}

func getChangesInLimitsSet(planLimits, stateLimits []TfLimitModel) (newElements, changedElements, removedElements []TfLimitModel) {
	var removedLimits, newLimits, changedLimits []TfLimitModel

	planLimitsMap := NewTfLimitModelMap(planLimits)
	stateTeamsMap := NewTfLimitModelMap(stateLimits)

	for name, stateLimit := range stateTeamsMap {
		if plannedTeam, exists := planLimitsMap[name]; exists {
			if !reflect.DeepEqual(plannedTeam, stateLimit) {
				changedLimits = append(changedLimits, plannedTeam)
			}
		} else {
			removedLimits = append(removedLimits, stateLimit)
		}
	}

	for name, limit := range planLimitsMap {
		if _, exists := stateTeamsMap[name]; !exists {
			newLimits = append(newLimits, limit)
		}
	}
	return newLimits, changedLimits, removedLimits
}
