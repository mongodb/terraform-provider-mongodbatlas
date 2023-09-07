package mongodbatlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20230201006/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ProjectDS{}
var _ datasource.DataSourceWithConfigure = &ProjectDS{}

func NewProjectDS() datasource.DataSource {
	return &ProjectDS{
		DSCommon: DSCommon{
			dataSourceName: projectResourceName,
		},
	}
}

type ProjectDS struct {
	DSCommon
}

type tfProjectDSModel struct {
	RegionUsageRestrictions                     types.String     `tfsdk:"region_usage_restrictions"`
	ProjectID                                   types.String     `tfsdk:"project_id"`
	Name                                        types.String     `tfsdk:"name"`
	OrgID                                       types.String     `tfsdk:"org_id"`
	Created                                     types.String     `tfsdk:"created"`
	ID                                          types.String     `tfsdk:"id"`
	Limits                                      []*tfLimitModel  `tfsdk:"limits"`
	Teams                                       []*tfTeamDSModel `tfsdk:"teams"`
	ClusterCount                                types.Int64      `tfsdk:"cluster_count"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool       `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool       `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool       `tfsdk:"is_schema_advisor_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool       `tfsdk:"is_performance_advisor_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool       `tfsdk:"is_extended_storage_sizes_enabled"`
	IsDataExplorerEnabled                       types.Bool       `tfsdk:"is_data_explorer_enabled"`
}

type tfTeamDSModel struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.List   `tfsdk:"role_names"`
}

func (d *ProjectDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_id": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("name")),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("project_id")),
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
			"limits": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"value": schema.Int64Attribute{
							Computed: true,
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
			},
		},
	}
}

func (d *ProjectDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var projectState tfProjectDSModel
	conn := d.client.Atlas
	connV2 := d.client.AtlasV2

	resp.Diagnostics.Append(req.Config.Get(ctx, &projectState)...)

	if projectState.ProjectID.IsNull() && projectState.Name.IsNull() {
		resp.Diagnostics.AddError("missing required attributes for data source", "either project_id or name must be configured")
		return
	}

	var (
		err     error
		project *matlas.Project
	)

	if !projectState.ProjectID.IsNull() {
		projectID := projectState.ProjectID.ValueString()
		project, _, err = conn.Projects.GetOneProject(ctx, projectID)
		if err != nil {
			resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(errorProjectRead, projectID, err.Error()))
			return
		}
	} else {
		name := projectState.Name.ValueString()
		project, _, err = conn.Projects.GetOneProjectByName(ctx, name)
		if err != nil {
			resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(errorProjectRead, name, err.Error()))
			return
		}
	}

	atlasTeams, atlasLimits, atlasProjectSettings, err := getProjectPropsFromAPI(ctx, conn, connV2, project.ID)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties", fmt.Sprintf(errorProjectRead, project.ID, err.Error()))
		return
	}

	projectState = newTFProjectDataSourceModel(ctx, project, atlasTeams, atlasProjectSettings, atlasLimits)

	resp.Diagnostics.Append(resp.State.Set(ctx, &projectState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func newTFProjectDataSourceModel(ctx context.Context, project *matlas.Project,
	teams *matlas.TeamsAssigned, projectSettings *matlas.ProjectSettings, limits []admin.DataFederationLimit) tfProjectDSModel {
	return tfProjectDSModel{
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
		Teams:                                       newTFTeamsDataSourceModel(ctx, teams),
		Limits:                                      newTFLimitsDataSourceModel(ctx, limits),
	}
}

func newTFTeamsDataSourceModel(ctx context.Context, atlasTeams *matlas.TeamsAssigned) []*tfTeamDSModel {
	if atlasTeams.TotalCount == 0 {
		return nil
	}
	teams := make([]*tfTeamDSModel, len(atlasTeams.Results))

	for i, atlasTeam := range atlasTeams.Results {
		roleNames, _ := types.ListValueFrom(ctx, types.StringType, atlasTeam.RoleNames)

		teams[i] = &tfTeamDSModel{
			TeamID:    types.StringValue(atlasTeam.TeamID),
			RoleNames: roleNames,
		}
	}
	return teams
}

func newTFLimitsDataSourceModel(ctx context.Context, dataFederationLimits []admin.DataFederationLimit) []*tfLimitModel {
	limits := make([]*tfLimitModel, len(dataFederationLimits))

	for i, dataFederationLimit := range dataFederationLimits {
		limits[i] = &tfLimitModel{
			Name:         types.StringValue(dataFederationLimit.Name),
			Value:        types.Int64Value(dataFederationLimit.Value),
			CurrentUsage: types.Int64PointerValue(dataFederationLimit.CurrentUsage),
			DefaultLimit: types.Int64PointerValue(dataFederationLimit.DefaultLimit),
			MaximumLimit: types.Int64PointerValue(dataFederationLimit.MaximumLimit),
		}
	}

	return limits
}
