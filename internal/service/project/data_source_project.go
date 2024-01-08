package project

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115002/admin"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &projectDS{}
var _ datasource.DataSourceWithConfigure = &projectDS{}

func DataSource() datasource.DataSource {
	return &projectDS{
		DSCommon: config.DSCommon{
			DataSourceName: projectResourceName,
		},
	}
}

type projectDS struct {
	config.DSCommon
}

type TfProjectDSModel struct {
	RegionUsageRestrictions                     types.String     `tfsdk:"region_usage_restrictions"`
	ProjectID                                   types.String     `tfsdk:"project_id"`
	Name                                        types.String     `tfsdk:"name"`
	OrgID                                       types.String     `tfsdk:"org_id"`
	Created                                     types.String     `tfsdk:"created"`
	ID                                          types.String     `tfsdk:"id"`
	Limits                                      []*TfLimitModel  `tfsdk:"limits"`
	Teams                                       []*TfTeamDSModel `tfsdk:"teams"`
	ClusterCount                                types.Int64      `tfsdk:"cluster_count"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool       `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool       `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool       `tfsdk:"is_schema_advisor_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool       `tfsdk:"is_performance_advisor_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool       `tfsdk:"is_extended_storage_sizes_enabled"`
	IsDataExplorerEnabled                       types.Bool       `tfsdk:"is_data_explorer_enabled"`
}

type TfTeamDSModel struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.List   `tfsdk:"role_names"`
}

func (d *projectDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *projectDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var projectState TfProjectDSModel
	connV2 := d.Client.AtlasV2

	resp.Diagnostics.Append(req.Config.Get(ctx, &projectState)...)

	if projectState.ProjectID.IsNull() && projectState.Name.IsNull() {
		resp.Diagnostics.AddError("missing required attributes for data source", "either project_id or name must be configured")
		return
	}

	var (
		err     error
		project *admin.Group
	)

	if !projectState.ProjectID.IsNull() {
		projectID := projectState.ProjectID.ValueString()
		project, _, err = connV2.ProjectsApi.GetProject(ctx, projectID).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(ErrorProjectRead, projectID, err.Error()))
			return
		}
	} else {
		name := projectState.Name.ValueString()
		project, _, err = connV2.ProjectsApi.GetProjectByName(ctx, name).Execute()
		if err != nil {
			resp.Diagnostics.AddError("error when getting project from Atlas", fmt.Sprintf(ErrorProjectRead, name, err.Error()))
			return
		}
	}

	atlasTeams, atlasLimits, atlasProjectSettings, err := GetProjectPropsFromAtlas(ctx, ServiceFromClient(connV2), project.GetId())
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties", fmt.Sprintf(ErrorProjectRead, project.GetId(), err.Error()))
		return
	}

	projectState = NewTFProjectDataSourceModel(ctx, project, atlasTeams, atlasProjectSettings, atlasLimits)

	resp.Diagnostics.Append(resp.State.Set(ctx, &projectState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
