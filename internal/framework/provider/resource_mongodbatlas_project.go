package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/framework/utils"
)

const (
	errorProjectCreate  = "error creating atlas project"
	errorProjectRead    = "error getting atlas project(%s)"
	errorProjectDelete  = "error deleting atlas project (%s)"
	errorProjectSetting = "error setting `%s` for atlas project (%s)"
)

var _ resource.Resource = &MongoDBAtlasProjectResource{}
var _ resource.ResourceWithImportState = &MongoDBAtlasProjectResource{}

func NewMongoDBAtlasProjectResource() resource.Resource {
	return &MongoDBAtlasProjectResource{}
}

type MongoDBAtlasProjectResource struct {
	client *MongoDBClient
}

type mongoDBAtlasProjectResourceModel struct {
	ID                                          types.String `tfsdk:"id"`
	Name                                        types.String `tfsdk:"name"`
	OrgID                                       types.String `tfsdk:"org_id"`
	ClusterCount                                types.Int64  `tfsdk:"cluster_count"`
	Created                                     types.String `tfsdk:"created"`
	ProjectOwnerID                              types.String `tfsdk:"project_owner_id"`
	WithDefaultAlertsSettings                   types.Bool   `tfsdk:"with_default_alerts_settings"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool   `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	IsDataExplorerEnabled                       types.Bool   `tfsdk:"is_data_explorer_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool   `tfsdk:"is_extended_storage_sizes_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool   `tfsdk:"is_performance_advisor_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool   `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool   `tfsdk:"is_schema_advisor_enabled"`
	RegionUsageRestrictions                     types.String `tfsdk:"region_usage_restrictions"`
	Teams                                       []team       `tfsdk:"teams"`
	ApiKeys                                     []apiKey     `tfsdk:"api_keys"`
}

type team struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.Set    `tfsdk:"role_names"`
}

type apiKey struct {
	ApiKeyID  types.String `tfsdk:"api_key_id"`
	RoleNames types.Set    `tfsdk:"role_names"`
}

func (r *MongoDBAtlasProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *MongoDBAtlasProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_count": schema.NumberAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"project_owner_id": schema.StringAttribute{
				Optional: true,
			},
			"with_default_alerts_settings": schema.BoolAttribute{
				Optional: true,
				Default:  booldefault.StaticBool(true),
			},
			// Since api_keys is a Computed attribute it will not be added as a Block:
			// https://developer.hashicorp.com/terraform/plugin/framework/migrating/attributes-blocks/blocks-computed
			"api_keys": schema.SetNestedAttribute{
				Optional:           true,
				Computed:           true,
				DeprecationMessage: fmt.Sprintf(DeprecationMessageParameterToResource, "v1.12.0", "mongodbatlas_project_api_key"),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_key_id": schema.StringAttribute{
							Required: true,
						},
						"role_names": schema.SetAttribute{
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"is_collect_database_specifics_statistics_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"is_data_explorer_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"is_extended_storage_sizes_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"is_performance_advisor_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"is_realtime_performance_panel_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"is_schema_advisor_enabled": schema.BoolAttribute{
				Computed: true,
				Optional: true,
			},
			"region_usage_restrictions": schema.StringAttribute{
				Computed: true,
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
		},
	}
}

func (r *MongoDBAtlasProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*MongoDBClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *MongoDBClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *MongoDBAtlasProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var projectPlan mongoDBAtlasProjectResourceModel
	conn := r.client.Atlas

	diags := req.Plan.Get(ctx, &projectPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectReq := &matlas.Project{
		OrgID:                     projectPlan.OrgID.ValueString(),
		Name:                      projectPlan.Name.ValueString(),
		WithDefaultAlertsSettings: projectPlan.WithDefaultAlertsSettings.ValueBoolPointer(),
		RegionUsageRestrictions:   projectPlan.RegionUsageRestrictions.ValueString(),
	}

	var createProjectOptions *matlas.CreateProjectOptions

	createProjectOptions = &matlas.CreateProjectOptions{
		ProjectOwnerID: projectPlan.ProjectOwnerID.ValueString(),
	}

	if !projectPlan.ProjectOwnerID.IsNull() {
		createProjectOptions = &matlas.CreateProjectOptions{
			ProjectOwnerID: projectPlan.ProjectOwnerID.ValueString(),
		}
	}

	project, _, err := conn.Projects.Create(ctx, projectReq, createProjectOptions)
	if err != nil {
		resp.Diagnostics.AddError(errorProjectCreate, err.Error())
		return
	}

	// Check if teams were set, if so we need to add the teams into the project
	if len(projectPlan.Teams) > 0 {
		// adding the teams into the project
		_, _, err := conn.Projects.AddTeamsToProject(ctx, project.ID, expandTeamsSet(ctx, projectPlan.Teams))
		if err != nil {
			errd := deleteProject(ctx, conn, project.ID)
			if errd != nil {
				resp.Diagnostics.AddError(fmt.Sprintf(errorProjectDelete, project.ID), err.Error())
				return
			}
			resp.Diagnostics.AddError("error adding teams into the project", err.Error())
			return
		}
	}

	// Check if api keys were set, if so we need to add keys into the project
	if len(projectPlan.ApiKeys) > 0 {
		// assign api keys to the project
		for _, apiKey := range projectPlan.ApiKeys {
			_, err := conn.ProjectAPIKeys.Assign(ctx, project.ID, apiKey.ApiKeyID.ValueString(), &matlas.AssignAPIKey{
				Roles: utils.StringSet(ctx, apiKey.RoleNames),
			})
			if err != nil {
				errd := deleteProject(ctx, conn, project.ID)
				if errd != nil {
					resp.Diagnostics.AddError(fmt.Sprintf(errorProjectDelete, project.ID), err.Error())
					return
				}
				resp.Diagnostics.AddError("error assigning api keys to the project", err.Error())
				return
			}
		}
	}

	projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, project.ID)
	if err != nil {
		errd := deleteProject(ctx, conn, project.ID)
		if errd != nil {
			resp.Diagnostics.AddError(fmt.Sprintf(errorProjectDelete, project.ID), err.Error())
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error getting project's settings assigned (%s):", project.ID), err.Error())
		return
	}
	projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled = projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled.ValueBoolPointer()
	projectSettings.IsDataExplorerEnabled = projectPlan.IsDataExplorerEnabled.ValueBoolPointer()
	projectSettings.IsExtendedStorageSizesEnabled = projectPlan.IsExtendedStorageSizesEnabled.ValueBoolPointer()
	projectSettings.IsPerformanceAdvisorEnabled = projectPlan.IsPerformanceAdvisorEnabled.ValueBoolPointer()
	projectSettings.IsRealtimePerformancePanelEnabled = projectPlan.IsRealtimePerformancePanelEnabled.ValueBoolPointer()
	projectSettings.IsSchemaAdvisorEnabled = projectPlan.IsSchemaAdvisorEnabled.ValueBoolPointer()

	_, _, err = conn.Projects.UpdateProjectSettings(ctx, project.ID, projectSettings)
	if err != nil {
		errd := deleteProject(ctx, conn, project.ID)
		if errd != nil {
			resp.Diagnostics.AddError(fmt.Sprintf(errorProjectDelete, project.ID), err.Error())
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("error updating project's settings assigned (%s):", project.ID), err.Error())
		return
	}

	// do a Read GET request
	projectID := project.ID
	projectRes, atlasResp, err := conn.Projects.GetOneProject(ctx, projectID)
	if err != nil {
		if resp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf(errorProjectRead, projectID), err.Error())
		return
	}

	// ---return resourceMongoDBAtlasProjectRead(ctx, d, meta)
	projectPlan, errString := getProjectPropsFromAtlas(ctx, conn, projectRes, projectPlan)

	if errString != "" {
		resp.Diagnostics.AddError(fmt.Sprintf(errorProjectRead, projectID), errString)
		return
	}

	// // Save updated data into Terraform state
	// resp.Diagnostics.Append(resp.State.Set(ctx, &projectPlan)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// Set state to fully populated data
	diags = resp.State.Set(ctx, projectPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func expandTeamsSet(ctx context.Context, teams []team) []*matlas.ProjectTeam {
	res := make([]*matlas.ProjectTeam, len(teams))

	for i, team := range teams {
		res[i] = &matlas.ProjectTeam{
			TeamID:    team.TeamID.ValueString(),
			RoleNames: utils.StringSet(ctx, team.RoleNames),
		}
	}
	return res
}

func deleteProject(ctx context.Context, conn *matlas.Client, projectID string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING", "RETRY"},
		Target:  []string{"IDLE"},
		// Refresh:    resourceProjectDependentsDeletingRefreshFunc(ctx, projectID, conn),
		Timeout:    30 * time.Minute,
		MinTimeout: 30 * time.Second,
		Delay:      0,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("[ERROR] could not determine MongoDB project %s dependents status: %s", projectID, err.Error()))
	}

	_, err = conn.Projects.Delete(ctx, projectID)

	// if err != nil {

	// 	return diag.Errorf(errorProjectDelete, projectID, err)
	// }

	return err
}

func (r *MongoDBAtlasProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var projectPlan mongoDBAtlasProjectResourceModel
	conn := r.client.Atlas

	// Get current state
	diags := req.State.Get(ctx, &projectPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectID := projectPlan.ID.ValueString()

	projectRes, atlasResp, err := conn.Projects.GetOneProject(ctx, projectID)
	if err != nil {
		if resp != nil && atlasResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf(errorProjectRead, projectID), err.Error())
		return
	}

	// teams, _, err := conn.Projects.GetProjectTeamsAssigned(ctx, projectID)
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("error getting project's teams assigned (%s)", projectID), err.Error())
	// 	return
	// }

	// apiKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, &matlas.ListOptions{})
	// if err != nil {
	// 	var target *matlas.ErrorResponse
	// 	if errors.As(err, &target) && target.ErrorCode != "USER_UNAUTHORIZED" {
	// 		resp.Diagnostics.AddError(fmt.Sprintf("error getting project's api keys (%s)", projectID), err.Error())
	// 		return
	// 	}
	// 	tflog.Info(ctx, "[WARN] `api_keys` will be empty because the user has no permissions to read the api keys endpoint")
	// }

	// projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, projectID)
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("error getting project's settings assigned (%s)", projectID), err.Error())
	// 	return
	// }

	// projectPlan = convertProjectToModel(ctx, projectPlan, projectRes, teams, apiKeys, projectSettings)

	projectPlan, errString := getProjectPropsFromAtlas(ctx, conn, projectRes, projectPlan)

	if errString != "" {
		resp.Diagnostics.AddError(fmt.Sprintf(errorProjectRead, projectID), errString)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectPlan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getProjectPropsFromAtlas(ctx context.Context, conn *matlas.Client, projectRes *matlas.Project, projectPlan mongoDBAtlasProjectResourceModel) (mongoDBAtlasProjectResourceModel, string) {
	projectID := projectRes.ID
	teams, _, err := conn.Projects.GetProjectTeamsAssigned(ctx, projectID)
	if err != nil {

		return projectPlan, fmt.Sprintf("error getting project's teams assigned (%s): %v", projectID, err.Error())
	}

	apiKeys, _, err := conn.ProjectAPIKeys.List(ctx, projectID, &matlas.ListOptions{})
	if err != nil {
		var target *matlas.ErrorResponse
		if errors.As(err, &target) && target.ErrorCode != "USER_UNAUTHORIZED" {
			return projectPlan, fmt.Sprintf("error getting project's api keys (%s): %v", projectID, err.Error())
		}
		tflog.Info(ctx, "[WARN] `api_keys` will be empty because the user has no permissions to read the api keys endpoint")
	}

	projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, projectID)
	if err != nil {
		return projectPlan, fmt.Sprintf("error getting project's settings assigned (%s): %v", projectID, err.Error())
	}

	return convertProjectToModel(ctx, projectID, projectRes, teams, apiKeys, projectSettings), ""
}

func convertProjectToModel(ctx context.Context, projectID string, projectRes *matlas.Project, teams *matlas.TeamsAssigned, apiKeys []matlas.APIKey, projectSettings *matlas.ProjectSettings) mongoDBAtlasProjectResourceModel {
	projectPlan := mongoDBAtlasProjectResourceModel{
		ID:           types.StringValue(projectID),
		Name:         types.StringValue(projectRes.Name),
		OrgID:        types.StringValue(projectRes.OrgID),
		ClusterCount: types.Int64Value(int64(projectRes.ClusterCount)),
		Created:      types.StringValue(projectRes.Created),
		IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(*projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled),
		IsDataExplorerEnabled:                       types.BoolValue(*projectSettings.IsDataExplorerEnabled),
		IsExtendedStorageSizesEnabled:               types.BoolValue(*projectSettings.IsExtendedStorageSizesEnabled),
		IsPerformanceAdvisorEnabled:                 types.BoolValue(*projectSettings.IsPerformanceAdvisorEnabled),
		IsRealtimePerformancePanelEnabled:           types.BoolValue(*projectSettings.IsRealtimePerformancePanelEnabled),
		IsSchemaAdvisorEnabled:                      types.BoolValue(*projectSettings.IsSchemaAdvisorEnabled),
		Teams:                                       convertTeamsToModel(ctx, teams),
		ApiKeys:                                     convertApiKeysToModel(ctx, apiKeys, projectID),
	}
	// projectPlan.Name = types.StringValue(projectRes.Name)
	// projectPlan.OrgID = types.StringValue(projectRes.OrgID)
	// projectPlan.ClusterCount = types.Int64Value(int64(projectRes.ClusterCount))
	// projectPlan.Created = types.StringValue(projectRes.Created)
	// projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled = types.BoolValue(*projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled)
	// projectPlan.IsDataExplorerEnabled = types.BoolValue(*projectSettings.IsDataExplorerEnabled)
	// projectPlan.IsExtendedStorageSizesEnabled = types.BoolValue(*projectSettings.IsExtendedStorageSizesEnabled)
	// projectPlan.IsPerformanceAdvisorEnabled = types.BoolValue(*projectSettings.IsPerformanceAdvisorEnabled)
	// projectPlan.IsRealtimePerformancePanelEnabled = types.BoolValue(*projectSettings.IsRealtimePerformancePanelEnabled)
	// projectPlan.IsSchemaAdvisorEnabled = types.BoolValue(*projectSettings.IsSchemaAdvisorEnabled)
	// projectPlan.Teams = convertTeamsToModel(ctx, teams)
	// projectPlan.ApiKeys = convertApiKeysToModel(ctx, apiKeys, projectID)

	return projectPlan
}

func convertApiKeysToModel(ctx context.Context, atlasApiKeys []matlas.APIKey, projectID string) []apiKey {
	res := make([]apiKey, len(atlasApiKeys))

	for _, atlasKey := range atlasApiKeys {
		id := atlasKey.ID

		var atlasRoles []string
		for _, role := range atlasKey.Roles {

			// ProjectAPIKeys.List returns all API keys of the Project, including the org and project roles
			// For more details: https://docs.atlas.mongodb.com/reference/api/projectApiKeys/get-all-apiKeys-in-one-project/
			if !strings.HasPrefix(role.RoleName, "ORG_") && role.GroupID == projectID {
				atlasRoles = append(atlasRoles, role.RoleName)
			}
		}
		roleNames, _ := types.SetValueFrom(ctx, types.StringType, atlasRoles)

		res = append(res, apiKey{
			ApiKeyID:  types.StringValue(id),
			RoleNames: roleNames,
		})
	}
	return res
}

func convertTeamsToModel(ctx context.Context, atlasTeams *matlas.TeamsAssigned) []team {
	teams := make([]team, atlasTeams.TotalCount)

	for i, atlasTeam := range atlasTeams.Results {
		roleNames, _ := types.SetValueFrom(ctx, types.StringType, atlasTeam.RoleNames)
		teams[i] = team{
			TeamID:    types.StringValue(atlasTeam.TeamID),
			RoleNames: roleNames,
		}
	}
	return teams
}

func (r *MongoDBAtlasProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var projectPlan *mongoDBAtlasProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &projectPlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &projectPlan)...)
}

func (r *MongoDBAtlasProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ExampleResourceModel

	// Read Terraform prior state data into the model.
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *MongoDBAtlasProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
