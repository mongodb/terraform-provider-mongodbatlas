package project

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"sort"
	"time"

	"go.mongodb.org/atlas-sdk/v20240805004/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	ErrorProjectRead               = "error getting project (%s): %s"
	errorProjectDelete             = "error deleting project (%s): %s"
	errorProjectUpdate             = "error updating project (%s): %s"
	errorProjectCreate             = "error creating project: %s"
	projectDependentsStateIdle     = "IDLE"
	projectDependentsStateDeleting = "DELETING"
	projectDependentsStateRetry    = "RETRY"
	projectResourceName            = "project"
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

type TFProjectRSModel struct {
	Limits                                      types.Set    `tfsdk:"limits"`
	Teams                                       types.Set    `tfsdk:"teams"`
	Tags                                        types.Map    `tfsdk:"tags"`
	IPAddresses                                 types.Object `tfsdk:"ip_addresses"`
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

type TFTeamModel struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.Set    `tfsdk:"role_names"`
}

type TFLimitModel struct {
	Name         types.String `tfsdk:"name"`
	Value        types.Int64  `tfsdk:"value"`
	CurrentUsage types.Int64  `tfsdk:"current_usage"`
	DefaultLimit types.Int64  `tfsdk:"default_limit"`
	MaximumLimit types.Int64  `tfsdk:"maximum_limit"`
}

type TFIPAddressesModel struct {
	Services TFServicesModel `tfsdk:"services"`
}

type TFServicesModel struct {
	Clusters []TFClusterIPsModel `tfsdk:"clusters"`
}

type TFClusterIPsModel struct {
	ClusterName types.String `tfsdk:"cluster_name"`
	Inbound     types.List   `tfsdk:"inbound"`
	Outbound    types.List   `tfsdk:"outbound"`
}

var IPAddressesObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"services": ServicesObjectType,
}}

var ServicesObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"clusters": types.ListType{ElemType: ClusterIPsObjectType},
}}

var ClusterIPsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"cluster_name": types.StringType,
	"inbound":      types.ListType{ElemType: types.StringType},
	"outbound":     types.ListType{ElemType: types.StringType},
}}

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
	AdvancedClusters *admin.PaginatedClusterDescription20240805
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
				Computed: true,
				Optional: true,
			},
			"ip_addresses": schema.SingleNestedAttribute{
				Computed:           true,
				DeprecationMessage: fmt.Sprintf(constant.DeprecationParamByVersionWithReplacement, "1.21.0", "mongodbatlas_project_ip_addresses data source"),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"services": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"clusters": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"cluster_name": schema.StringAttribute{
											Computed: true,
										},
										"inbound": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
										},
										"outbound": schema.ListAttribute{
											ElementType: types.StringType,
											Computed:    true,
										},
									},
								},
							},
						},
					},
				},
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
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
	var projectPlan TFProjectRSModel
	var teams []TFTeamModel
	var limits []TFLimitModel

	connV2 := r.Client.AtlasV2

	diags := req.Plan.Get(ctx, &projectPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tags := NewResourceTags(ctx, projectPlan.Tags)
	projectGroup := &admin.Group{
		OrgId:                     projectPlan.OrgID.ValueString(),
		Name:                      projectPlan.Name.ValueString(),
		WithDefaultAlertsSettings: projectPlan.WithDefaultAlertsSettings.ValueBoolPointer(),
		RegionUsageRestrictions:   conversion.StringNullIfEmpty(projectPlan.RegionUsageRestrictions.ValueString()).ValueStringPointer(),
		Tags:                      &tags,
	}

	projectAPIParams := &admin.CreateProjectApiParams{
		Group:          projectGroup,
		ProjectOwnerId: conversion.StringNullIfEmpty(projectPlan.ProjectOwnerID.ValueString()).ValueStringPointer(),
	}

	// create project
	project, _, err := connV2.ProjectsApi.CreateProjectWithParams(ctx, projectAPIParams).Execute()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(errorProjectCreate, projectPlan.Name.ValueString()), err.Error())
		return
	}

	// add teams
	if len(projectPlan.Teams.Elements()) > 0 {
		_ = projectPlan.Teams.ElementsAs(ctx, &teams, false)

		_, _, err := connV2.TeamsApi.AddAllTeamsToProject(ctx, project.GetId(), NewTeamRoleList(ctx, teams)).Execute()
		if err != nil {
			errd := deleteProject(ctx, connV2.ClustersApi, connV2.ProjectsApi, project.GetId())
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
				errd := deleteProject(ctx, connV2.ClustersApi, connV2.ProjectsApi, project.GetId())
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
		errd := deleteProject(ctx, connV2.ClustersApi, connV2.ProjectsApi, project.GetId())
		if errd != nil {
			resp.Diagnostics.AddError("error during project deletion when getting project settings", fmt.Sprintf(errorProjectDelete, project.GetId(), err.Error()))
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project's settings assigned (%s):", project.GetId()), err.Error())
		return
	}

	SetProjectBool(projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled, &projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled)
	SetProjectBool(projectPlan.IsDataExplorerEnabled, &projectSettings.IsDataExplorerEnabled)
	SetProjectBool(projectPlan.IsExtendedStorageSizesEnabled, &projectSettings.IsExtendedStorageSizesEnabled)
	SetProjectBool(projectPlan.IsPerformanceAdvisorEnabled, &projectSettings.IsPerformanceAdvisorEnabled)
	SetProjectBool(projectPlan.IsRealtimePerformancePanelEnabled, &projectSettings.IsRealtimePerformancePanelEnabled)
	SetProjectBool(projectPlan.IsSchemaAdvisorEnabled, &projectSettings.IsSchemaAdvisorEnabled)

	if _, _, err = connV2.ProjectsApi.UpdateProjectSettings(ctx, project.GetId(), projectSettings).Execute(); err != nil {
		errd := deleteProject(ctx, connV2.ClustersApi, connV2.ProjectsApi, project.GetId())
		if errd != nil {
			resp.Diagnostics.AddError("error during project deletion when updating project settings", fmt.Sprintf(errorProjectDelete, project.GetId(), err.Error()))
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error updating project's settings assigned (%s):", project.GetId()), err.Error())
		return
	}

	projectID := project.GetId()
	projectRes, _, err := connV2.ProjectsApi.GetProject(ctx, projectID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error when getting project after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	// get project props
	projectProps, err := GetProjectPropsFromAPI(ctx, connV2.ProjectsApi, connV2.TeamsApi, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	filteredLimits := FilterUserDefinedLimits(projectProps.Limits, limits)
	projectProps.Limits = filteredLimits

	projectPlanNew, diags := NewTFProjectResourceModel(ctx, projectRes, *projectProps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updatePlanFromConfig(projectPlanNew, &projectPlan)

	// set state to fully populated data
	diags = resp.State.Set(ctx, projectPlanNew)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectRS) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var projectState TFProjectRSModel
	var limits []TFLimitModel
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
		if atlasResp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	// get project props
	projectProps, err := GetProjectPropsFromAPI(ctx, connV2.ProjectsApi, connV2.TeamsApi, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	filteredLimits := FilterUserDefinedLimits(projectProps.Limits, limits)
	projectProps.Limits = filteredLimits

	projectStateNew, diags := NewTFProjectResourceModel(ctx, projectRes, *projectProps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updatePlanFromConfig(projectStateNew, &projectState)

	// save read data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectStateNew)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *projectRS) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var projectState TFProjectRSModel
	var projectPlan TFProjectRSModel
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

	err := UpdateProject(ctx, connV2.ProjectsApi, &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	err = UpdateProjectTeams(ctx, connV2.TeamsApi, &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project teams update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	err = UpdateProjectLimits(ctx, connV2.ProjectsApi, &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project limits update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	err = updateProjectSettings(ctx, connV2.ProjectsApi, &projectState, &projectPlan)
	if err != nil {
		resp.Diagnostics.AddError("error in project settings update", fmt.Sprintf(errorProjectUpdate, projectID, err.Error()))
		return
	}

	projectRes, _, err := connV2.ProjectsApi.GetProject(ctx, projectID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error when getting project after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}

	// get project props
	projectProps, err := GetProjectPropsFromAPI(ctx, connV2.ProjectsApi, connV2.TeamsApi, projectID)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties after create", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
		return
	}
	var planLimits []TFLimitModel
	_ = projectPlan.Limits.ElementsAs(ctx, &planLimits, false)

	filteredLimits := FilterUserDefinedLimits(projectProps.Limits, planLimits)
	projectProps.Limits = filteredLimits

	projectPlanNew, diags := NewTFProjectResourceModel(ctx, projectRes, *projectProps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	updatePlanFromConfig(projectPlanNew, &projectPlan)

	// save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectPlanNew)...)
}

func (r *projectRS) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var project *TFProjectRSModel

	// read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &project)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := project.ID.ValueString()
	err := deleteProject(ctx, r.Client.AtlasV2.ClustersApi, r.Client.AtlasV2.ProjectsApi, projectID)

	if err != nil {
		resp.Diagnostics.AddError("error when destroying resource", fmt.Sprintf(errorProjectDelete, projectID, err.Error()))
		return
	}
}

func (r *projectRS) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func updatePlanFromConfig(projectPlanNewPtr, projectPlan *TFProjectRSModel) {
	// we need to reset defaults from what was previously in the state:
	// https://discuss.hashicorp.com/t/boolean-optional-default-value-migration-to-framework/55932
	projectPlanNewPtr.WithDefaultAlertsSettings = projectPlan.WithDefaultAlertsSettings
	projectPlanNewPtr.ProjectOwnerID = projectPlan.ProjectOwnerID
	if projectPlan.Tags.IsNull() && len(projectPlanNewPtr.Tags.Elements()) == 0 {
		projectPlanNewPtr.Tags = types.MapNull(types.StringType)
	}
}

func FilterUserDefinedLimits(allAtlasLimits []admin.DataFederationLimit, tflimits []TFLimitModel) []admin.DataFederationLimit {
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

type AdditionalProperties struct {
	Teams       *admin.PaginatedTeamRole
	Settings    *admin.GroupSettings
	IPAddresses *admin.GroupIPAddresses
	Limits      []admin.DataFederationLimit
}

// GetProjectPropsFromAPI fetches properties obtained from complementary endpoints associated with a project.
func GetProjectPropsFromAPI(ctx context.Context, projectsAPI admin.ProjectsApi, teamsAPI admin.TeamsApi, projectID string) (*AdditionalProperties, error) {
	teams, _, err := teamsAPI.ListProjectTeams(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting project's teams assigned (%s): %v", projectID, err.Error())
	}

	limits, _, err := projectsAPI.ListProjectLimits(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting project's limits (%s): %s", projectID, err.Error())
	}

	projectSettings, _, err := projectsAPI.GetProjectSettings(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting project's settings assigned (%s): %v", projectID, err.Error())
	}

	ipAddresses, _, err := projectsAPI.ReturnAllIPAddresses(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting project's IP addresses (%s): %v", projectID, err.Error())
	}

	return &AdditionalProperties{
		Teams:       teams,
		Limits:      limits,
		Settings:    projectSettings,
		IPAddresses: ipAddresses,
	}, nil
}

func updateProjectSettings(ctx context.Context, projectsAPI admin.ProjectsApi, state, plan *TFProjectRSModel) error {
	projectID := state.ID.ValueString()
	settings, _, err := projectsAPI.GetProjectSettings(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("error getting project's settings assigned: %v", err.Error())
	}

	hasChanged := UpdateProjectBool(plan.IsCollectDatabaseSpecificsStatisticsEnabled, state.IsCollectDatabaseSpecificsStatisticsEnabled, &settings.IsCollectDatabaseSpecificsStatisticsEnabled)
	hasChanged = UpdateProjectBool(plan.IsDataExplorerEnabled, state.IsDataExplorerEnabled, &settings.IsDataExplorerEnabled) || hasChanged
	hasChanged = UpdateProjectBool(plan.IsExtendedStorageSizesEnabled, state.IsExtendedStorageSizesEnabled, &settings.IsExtendedStorageSizesEnabled) || hasChanged
	hasChanged = UpdateProjectBool(plan.IsPerformanceAdvisorEnabled, state.IsPerformanceAdvisorEnabled, &settings.IsPerformanceAdvisorEnabled) || hasChanged
	hasChanged = UpdateProjectBool(plan.IsRealtimePerformancePanelEnabled, state.IsRealtimePerformancePanelEnabled, &settings.IsRealtimePerformancePanelEnabled) || hasChanged
	hasChanged = UpdateProjectBool(plan.IsSchemaAdvisorEnabled, state.IsSchemaAdvisorEnabled, &settings.IsSchemaAdvisorEnabled) || hasChanged

	if hasChanged {
		_, _, err = projectsAPI.UpdateProjectSettings(ctx, projectID, settings).Execute()
		if err != nil {
			return fmt.Errorf("error updating project's settings assigned: %v", err.Error())
		}
	}
	return nil
}

func UpdateProjectLimits(ctx context.Context, projectsAPI admin.ProjectsApi, projectState, projectPlan *TFProjectRSModel) error {
	var planLimits []TFLimitModel
	var stateLimits []TFLimitModel
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
		if _, _, err := projectsAPI.DeleteProjectLimit(ctx, limitName, projectID).Execute(); err != nil {
			return fmt.Errorf("error removing limit %s from the project(%s) during update: %s", limitName, projectID, err)
		}
	}

	// updating values for changed limits
	if len(changedLimits) > 0 {
		if err := setProjectLimits(ctx, projectsAPI, projectID, changedLimits); err != nil {
			return fmt.Errorf("error adding modified limits into the project during update: %v", err.Error())
		}
	}

	// adding new limits into the project
	if len(newLimits) > 0 {
		if err := setProjectLimits(ctx, projectsAPI, projectID, newLimits); err != nil {
			return fmt.Errorf("error adding limits into the project during update: %v", err.Error())
		}
	}

	return nil
}

func setProjectLimits(ctx context.Context, projectsAPI admin.ProjectsApi, projectID string, tfLimits []TFLimitModel) error {
	for _, limit := range tfLimits {
		dataFederationLimit := &admin.DataFederationLimit{
			Name:  limit.Name.ValueString(),
			Value: limit.Value.ValueInt64(),
		}
		_, _, err := projectsAPI.SetProjectLimit(ctx, limit.Name.ValueString(), projectID, dataFederationLimit).Execute()
		if err != nil {
			return fmt.Errorf("error adding limits into the project: %v", err.Error())
		}
	}
	return nil
}

func UpdateProjectTeams(ctx context.Context, teamsAPI admin.TeamsApi, projectState, projectPlan *TFProjectRSModel) error {
	var planTeams []TFTeamModel
	var stateTeams []TFTeamModel
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
		_, err := teamsAPI.RemoveProjectTeam(ctx, projectID, teamID).Execute()
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
		roleNames := conversion.TypesSetToString(ctx, team.RoleNames)
		_, _, err := teamsAPI.UpdateTeamRoles(ctx, projectID, teamID,
			&admin.TeamRole{
				RoleNames: &roleNames,
			},
		).Execute()
		if err != nil {
			return fmt.Errorf("error updating role names for the team(%s): %s", teamID, err.Error())
		}
	}

	// adding new teams into the project
	if _, _, err := teamsAPI.AddAllTeamsToProject(ctx, projectID, NewTeamRoleList(ctx, newTeams)).Execute(); err != nil {
		return fmt.Errorf("error adding teams to the project: %v", err.Error())
	}

	return nil
}

func hasTeamsChanged(planTeams, stateTeams []TFTeamModel) bool {
	sort.Slice(planTeams, func(i, j int) bool {
		return planTeams[i].TeamID.ValueString() < planTeams[j].TeamID.ValueString()
	})
	sort.Slice(stateTeams, func(i, j int) bool {
		return stateTeams[i].TeamID.ValueString() < stateTeams[j].TeamID.ValueString()
	})
	return !reflect.DeepEqual(planTeams, stateTeams)
}

func hasLimitsChanged(planLimits, stateLimits []TFLimitModel) bool {
	sort.Slice(planLimits, func(i, j int) bool {
		return planLimits[i].Name.ValueString() < planLimits[j].Name.ValueString()
	})
	sort.Slice(stateLimits, func(i, j int) bool {
		return stateLimits[i].Name.ValueString() < stateLimits[j].Name.ValueString()
	})
	return !reflect.DeepEqual(planLimits, stateLimits)
}

func UpdateProject(ctx context.Context, projectsAPI admin.ProjectsApi, projectState, projectPlan *TFProjectRSModel) error {
	tagsBefore := NewResourceTags(ctx, projectState.Tags)
	tagsAfter := NewResourceTags(ctx, projectPlan.Tags)
	if projectPlan.Name.Equal(projectState.Name) && reflect.DeepEqual(tagsBefore, tagsAfter) {
		return nil
	}

	projectID := projectState.ID.ValueString()

	if _, _, err := projectsAPI.UpdateProject(ctx, projectID, NewGroupUpdate(projectPlan, &tagsAfter)).Execute(); err != nil {
		return fmt.Errorf("error updating the project(%s): %s", projectID, err)
	}

	return nil
}

func deleteProject(ctx context.Context, clustersAPI admin.ClustersApi, projectsAPI admin.ProjectsApi, projectID string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{projectDependentsStateDeleting, projectDependentsStateRetry},
		Target:     []string{projectDependentsStateIdle},
		Refresh:    ResourceProjectDependentsDeletingRefreshFunc(ctx, projectID, clustersAPI),
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      0,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("[ERROR] could not determine MongoDB project %s dependents status: %s", projectID, err.Error()))
	}

	_, _, err = projectsAPI.DeleteProject(ctx, projectID).Execute()

	return err
}

/*
resourceProjectDependentsDeletingRefreshFunc assumes the project CRUD outcome will be the same for any non-zero number of dependents

If all dependents are deleting, wait to try and delete
Else consider the aggregate dependents idle.

If we get a defined error response, return that right away
Else retry
*/
func ResourceProjectDependentsDeletingRefreshFunc(ctx context.Context, projectID string, clustersAPI admin.ClustersApi) retry.StateRefreshFunc {
	return func() (any, string, error) {
		clusters, _, listClustersErr := clustersAPI.ListClusters(ctx, projectID).Execute()
		dependents := AtlasProjectDependants{AdvancedClusters: clusters}

		if listClustersErr != nil {
			return nil, "", listClustersErr
		}

		if *dependents.AdvancedClusters.TotalCount == 0 {
			return dependents, projectDependentsStateIdle, nil
		}

		results := dependents.AdvancedClusters.GetResults()
		for i := range results {
			if *results[i].StateName != projectDependentsStateDeleting {
				return dependents, projectDependentsStateIdle, nil
			}
		}

		log.Printf("[DEBUG] status for MongoDB project %s dependents: %s", projectID, projectDependentsStateDeleting)

		return dependents, projectDependentsStateDeleting, nil
	}
}

func getChangesInTeamsSet(planTeams, stateTeams []TFTeamModel) (newElements, changedElements, removedElements []TFTeamModel) {
	var removedTeams, newTeams, changedTeams []TFTeamModel

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

func getChangesInLimitsSet(planLimits, stateLimits []TFLimitModel) (newElements, changedElements, removedElements []TFLimitModel) {
	var removedLimits, newLimits, changedLimits []TFLimitModel

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
