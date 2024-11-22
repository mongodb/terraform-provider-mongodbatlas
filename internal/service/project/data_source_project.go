package project

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20241113001/admin"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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

type TFProjectDSModel struct {
	IPAddresses                                 types.Object    `tfsdk:"ip_addresses"`
	Created                                     types.String    `tfsdk:"created"`
	OrgID                                       types.String    `tfsdk:"org_id"`
	RegionUsageRestrictions                     types.String    `tfsdk:"region_usage_restrictions"`
	ID                                          types.String    `tfsdk:"id"`
	Name                                        types.String    `tfsdk:"name"`
	ProjectID                                   types.String    `tfsdk:"project_id"`
	Tags                                        types.Map       `tfsdk:"tags"`
	Teams                                       []*TFTeamModel  `tfsdk:"teams"`
	Limits                                      []*TFLimitModel `tfsdk:"limits"`
	ClusterCount                                types.Int64     `tfsdk:"cluster_count"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool      `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool      `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool      `tfsdk:"is_schema_advisor_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool      `tfsdk:"is_performance_advisor_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool      `tfsdk:"is_extended_storage_sizes_enabled"`
	IsDataExplorerEnabled                       types.Bool      `tfsdk:"is_data_explorer_enabled"`
	IsSlowOperationThresholdingEnabled          types.Bool      `tfsdk:"is_slow_operation_thresholding_enabled"`
}

func (d *projectDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	overridenFields := map[string]schema.Attribute{
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
		"project_owner_id":             nil,
		"with_default_alerts_settings": nil,
	}
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), nil, overridenFields)
}

func (d *projectDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var projectState TFProjectDSModel
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

	projectProps, err := GetProjectPropsFromAPI(ctx, connV2.ProjectsApi, connV2.TeamsApi, connV2.PerformanceAdvisorApi, project.GetId(), &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("error when getting project properties", fmt.Sprintf(ErrorProjectRead, project.GetId(), err.Error()))
		return
	}

	newProjectState, diags := NewTFProjectDataSourceModel(ctx, project, *projectProps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newProjectState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
