package provider

import (
	"context"
	"fmt"

	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/utils"
)

var _ datasource.DataSource = &ProjectDataSource{}
var _ datasource.DataSourceWithConfigure = &ProjectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

type ProjectDataSource struct {
	client *MongoDBClient
}

type projectDataSourceModel struct {
	ID                                          types.String              `tfsdk:"id"`
	ProjectID                                   types.String              `tfsdk:"project_id"`
	Name                                        types.String              `tfsdk:"name"`
	OrgID                                       types.String              `tfsdk:"org_id"`
	ClusterCount                                types.Int64               `tfsdk:"cluster_count"`
	Created                                     types.String              `tfsdk:"created"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool                `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	IsDataExplorerEnabled                       types.Bool                `tfsdk:"is_data_explorer_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool                `tfsdk:"is_extended_storage_sizes_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool                `tfsdk:"is_performance_advisor_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool                `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool                `tfsdk:"is_schema_advisor_enabled"`
	RegionUsageRestrictions                     types.String              `tfsdk:"region_usage_restrictions"`
	Teams                                       []projectDataSourceTeam   `tfsdk:"teams"`
	ApiKeys                                     []projectDataSourceApiKey `tfsdk:"api_keys"`
}

type projectDataSourceTeam struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.List   `tfsdk:"role_names"`
}

type projectDataSourceApiKey struct {
	ApiKeyID  types.String `tfsdk:"api_key_id"`
	RoleNames types.List   `tfsdk:"role_names"`
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("name")),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("project_id")),
				},
			},
			"org_id": schema.StringAttribute{
				Computed: true,
			},
			"cluster_count": schema.Int64Attribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
			"is_collect_database_specifics_statistics_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"is_data_explorer_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"is_extended_storage_sizes_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"is_performance_advisor_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"is_realtime_performance_panel_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"is_schema_advisor_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"region_usage_restrictions": schema.StringAttribute{
				Computed: true,
			},
			"api_keys": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_key_id": schema.StringAttribute{
							Computed: true,
						},
						"role_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
			"teams": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"team_id": schema.StringAttribute{
							Computed: true,
						},
						"role_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
		// Blocks: map[string]schema.Block{
		// 	"teams": schema.ListNestedBlock{
		// 		NestedObject: schema.NestedBlockObject{
		// 			Attributes: map[string]schema.Attribute{
		// 				"team_id": schema.StringAttribute{
		// 					Computed: true,
		// 				},
		// 				"role_names": schema.ListAttribute{
		// 					Computed:    true,
		// 					ElementType: types.StringType,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	"api_keys": schema.ListNestedBlock{
		// 		NestedObject: schema.NestedBlockObject{
		// 			Attributes: map[string]schema.Attribute{
		// 				"api_key_id": schema.StringAttribute{
		// 					Computed: true,
		// 				},
		// 				"role_names": schema.ListAttribute{
		// 					Computed:    true,
		// 					ElementType: types.StringType,
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var stateModel projectDataSourceModel
	conn := d.client.Atlas

	resp.Diagnostics.Append(req.Config.Get(ctx, &stateModel)...)

	if stateModel.ProjectID.IsNull() && stateModel.Name.IsNull() {
		resp.Diagnostics.AddError("missing required attributes for data source", "either project_id or name must be configured")
		return
	}

	var (
		err     error
		project *matlas.Project
	)

	if !stateModel.ProjectID.IsNull() {
		project, _, err = conn.Projects.GetOneProject(ctx, stateModel.ProjectID.ValueString())
	} else {
		project, _, err = conn.Projects.GetOneProjectByName(ctx, stateModel.Name.ValueString())
	}
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf(errorProjectRead, stateModel.ProjectID.ValueString()), err.Error())
		return

	}

	teams, apiKeys, projectSettings, err := getProjectPropsFromAtlas(ctx, conn, project)
	if err != nil {
		resp.Diagnostics.AddError("error in monogbatlas_projects data source", err.Error())
		return
	}

	stateModel = toProjectDataSourceModel(ctx, project, teams, apiKeys, projectSettings)

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getProjectPropsFromAtlas(ctx context.Context, conn *matlas.Client, projectRes *matlas.Project) (*matlas.TeamsAssigned, []matlas.APIKey, *matlas.ProjectSettings, error) {
	projectID := projectRes.ID
	teams, _, err := conn.Projects.GetProjectTeamsAssigned(ctx, projectID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting project's teams assigned (%s): %v", projectID, err.Error())
	}

	apiKeys, err := getProjectAPIKeys(ctx, conn, projectID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting project's teams assigned (%s): %v", projectID, err.Error())
	}

	projectSettings, _, err := conn.Projects.GetProjectSettings(ctx, projectID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting project's settings assigned (%s): %v", projectID, err.Error())
	}

	return teams, apiKeys, projectSettings, nil
}

func toProjectDataSourceModel(ctx context.Context, project *matlas.Project, teams *matlas.TeamsAssigned, apiKeys []matlas.APIKey, projectSettings *matlas.ProjectSettings) projectDataSourceModel {
	projectStateModel := projectDataSourceModel{
		ID:           types.StringValue(project.ID),
		ProjectID:    types.StringValue(project.ID),
		Name:         types.StringValue(project.Name),
		OrgID:        types.StringValue(project.OrgID),
		ClusterCount: types.Int64Value(int64(project.ClusterCount)),
		Created:      types.StringValue(project.Created),
		IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(*projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled),
		IsDataExplorerEnabled:                       types.BoolValue(*projectSettings.IsDataExplorerEnabled),
		IsExtendedStorageSizesEnabled:               types.BoolValue(*projectSettings.IsExtendedStorageSizesEnabled),
		IsPerformanceAdvisorEnabled:                 types.BoolValue(*projectSettings.IsPerformanceAdvisorEnabled),
		IsRealtimePerformancePanelEnabled:           types.BoolValue(*projectSettings.IsRealtimePerformancePanelEnabled),
		IsSchemaAdvisorEnabled:                      types.BoolValue(*projectSettings.IsSchemaAdvisorEnabled),
		Teams:                                       toTeamsDataSourceModel(ctx, teams),
		ApiKeys:                                     toApiKeysDataSourceModel(ctx, apiKeys),
	}

	return projectStateModel
}

func toApiKeysDataSourceModel(ctx context.Context, atlasApiKeys []matlas.APIKey) []projectDataSourceApiKey {
	res := []projectDataSourceApiKey{}

	for _, atlasKey := range atlasApiKeys {
		id := atlasKey.ID

		var atlasRoles []attr.Value
		for _, role := range atlasKey.Roles {
			atlasRoles = append(atlasRoles, types.StringValue(role.RoleName))

		}

		res = append(res, projectDataSourceApiKey{
			ApiKeyID:  types.StringValue(id),
			RoleNames: utils.ArrToListValue(atlasRoles),
		})
	}
	return res
}

func toTeamsDataSourceModel(ctx context.Context, atlasTeams *matlas.TeamsAssigned) []projectDataSourceTeam {
	if atlasTeams.TotalCount == 0 {
		return nil
	}
	teams := make([]projectDataSourceTeam, atlasTeams.TotalCount)

	for i, atlasTeam := range atlasTeams.Results {
		roleNames, _ := types.ListValueFrom(ctx, types.StringType, atlasTeam.RoleNames)

		teams[i] = projectDataSourceTeam{
			TeamID:    types.StringValue(atlasTeam.TeamID),
			RoleNames: roleNames,
		}
	}
	return teams
}
