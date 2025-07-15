package project

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func UsersProjectSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},
				"org_membership_status": schema.StringAttribute{
					Computed: true,
				},
				"roles": schema.SetAttribute{
					Computed:    true,
					ElementType: types.StringType,
				},
				"username": schema.StringAttribute{
					Computed: true,
				},
				"invitation_created_at": schema.StringAttribute{
					Computed: true,
				},
				"invitation_expires_at": schema.StringAttribute{
					Computed: true,
				},
				"inviter_username": schema.StringAttribute{
					Computed: true,
				},
				"country": schema.StringAttribute{
					Computed: true,
				},
				"created_at": schema.StringAttribute{
					Computed: true,
				},
				"first_name": schema.StringAttribute{
					Computed: true,
				},
				"last_auth": schema.StringAttribute{
					Computed: true,
				},
				"last_name": schema.StringAttribute{
					Computed: true,
				},
				"mobile_number": schema.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func dataSourceOverridenFields() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("project_id")),
			},
		},
		"project_id": schema.StringAttribute{
			Optional: true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("name")),
			},
		},
		"users": UsersProjectSchema(),
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
		"project_owner_id":             nil,
		"with_default_alerts_settings": nil,
	}
}

type tfProjectsDSModel struct {
	ID           types.String        `tfsdk:"id"`
	Results      []*TFProjectDSModel `tfsdk:"results"`
	PageNum      types.Int64         `tfsdk:"page_num"`
	ItemsPerPage types.Int64         `tfsdk:"items_per_page"`
	TotalCount   types.Int64         `tfsdk:"total_count"`
}

type TFProjectDSModel struct {
	Tags                                        types.Map              `tfsdk:"tags"`
	IPAddresses                                 types.Object           `tfsdk:"ip_addresses"`
	Created                                     types.String           `tfsdk:"created"`
	OrgID                                       types.String           `tfsdk:"org_id"`
	RegionUsageRestrictions                     types.String           `tfsdk:"region_usage_restrictions"`
	ID                                          types.String           `tfsdk:"id"`
	Name                                        types.String           `tfsdk:"name"`
	ProjectID                                   types.String           `tfsdk:"project_id"`
	Teams                                       []*TFTeamDSModel       `tfsdk:"teams"`
	Limits                                      []*TFLimitModel        `tfsdk:"limits"`
	Users                                       []*TFCloudUsersDSModel `tfsdk:"users"`
	ClusterCount                                types.Int64            `tfsdk:"cluster_count"`
	IsCollectDatabaseSpecificsStatisticsEnabled types.Bool             `tfsdk:"is_collect_database_specifics_statistics_enabled"`
	IsRealtimePerformancePanelEnabled           types.Bool             `tfsdk:"is_realtime_performance_panel_enabled"`
	IsSchemaAdvisorEnabled                      types.Bool             `tfsdk:"is_schema_advisor_enabled"`
	IsPerformanceAdvisorEnabled                 types.Bool             `tfsdk:"is_performance_advisor_enabled"`
	IsExtendedStorageSizesEnabled               types.Bool             `tfsdk:"is_extended_storage_sizes_enabled"`
	IsDataExplorerEnabled                       types.Bool             `tfsdk:"is_data_explorer_enabled"`
	IsSlowOperationThresholdingEnabled          types.Bool             `tfsdk:"is_slow_operation_thresholding_enabled"`
}

type TFTeamDSModel struct {
	TeamID    types.String `tfsdk:"team_id"`
	RoleNames types.List   `tfsdk:"role_names"`
}

type TFCloudUsersDSModel struct {
	ID                  types.String `tfsdk:"id"`
	OrgMembershipStatus types.String `tfsdk:"org_membership_status"`
	Roles               types.Set    `tfsdk:"roles"`
	Username            types.String `tfsdk:"username"`
	InvitationCreatedAt types.String `tfsdk:"invitation_created_at"`
	InvitationExpiresAt types.String `tfsdk:"invitation_expires_at"`
	InviterUsername     types.String `tfsdk:"inviter_username"`
	Country             types.String `tfsdk:"country"`
	CreatedAt           types.String `tfsdk:"created_at"`
	FirstName           types.String `tfsdk:"first_name"`
	LastAuth            types.String `tfsdk:"last_auth"`
	LastName            types.String `tfsdk:"last_name"`
	MobileNumber        types.String `tfsdk:"mobile_number"`
}

func NewTFProjectDataSourceModel(ctx context.Context, project *admin.Group, projectProps *AdditionalProperties) (*TFProjectDSModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if project == nil {
		diags.AddError("Invalid Project Data", "Project data is nil and cannot be processed")
		return nil, diags
	}
	if projectProps == nil {
		diags.AddError("Invalid Project Properties", "Project properties data is nil and cannot be processed")
		return nil, diags
	}
	ipAddressesModel, ipDiags := NewTFIPAddressesModel(ctx, projectProps.IPAddresses)
	diags.Append(ipDiags...)
	if diags.HasError() {
		return nil, diags
	}
	projectSettings := projectProps.Settings
	return &TFProjectDSModel{
		ID:                      types.StringValue(project.GetId()),
		ProjectID:               types.StringValue(project.GetId()),
		Name:                    types.StringValue(project.Name),
		OrgID:                   types.StringValue(project.OrgId),
		ClusterCount:            types.Int64Value(project.ClusterCount),
		Created:                 types.StringValue(conversion.TimeToString(project.Created)),
		RegionUsageRestrictions: types.StringPointerValue(project.RegionUsageRestrictions),
		IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(*projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled),
		IsDataExplorerEnabled:                       types.BoolValue(*projectSettings.IsDataExplorerEnabled),
		IsExtendedStorageSizesEnabled:               types.BoolValue(*projectSettings.IsExtendedStorageSizesEnabled),
		IsPerformanceAdvisorEnabled:                 types.BoolValue(*projectSettings.IsPerformanceAdvisorEnabled),
		IsRealtimePerformancePanelEnabled:           types.BoolValue(*projectSettings.IsRealtimePerformancePanelEnabled),
		IsSchemaAdvisorEnabled:                      types.BoolValue(*projectSettings.IsSchemaAdvisorEnabled),
		Teams:                                       NewTFTeamsDataSourceModel(ctx, projectProps.Teams),
		Limits:                                      NewTFLimitsDataSourceModel(ctx, projectProps.Limits),
		IPAddresses:                                 ipAddressesModel,
		Tags:                                        conversion.NewTFTags(project.GetTags()),
		IsSlowOperationThresholdingEnabled:          types.BoolValue(projectProps.IsSlowOperationThresholdingEnabled),
		Users:                                       NewTFCloudUsersDataSourceModel(ctx, projectProps.Users),
	}, nil
}

func NewTFTeamsDataSourceModel(ctx context.Context, atlasTeams *admin.PaginatedTeamRole) []*TFTeamDSModel {
	if atlasTeams == nil || atlasTeams.GetTotalCount() == 0 {
		return nil
	}
	results := atlasTeams.GetResults()
	teams := make([]*TFTeamDSModel, len(results))
	for i, atlasTeam := range results {
		roleNames, _ := types.ListValueFrom(ctx, types.StringType, atlasTeam.RoleNames)
		teams[i] = &TFTeamDSModel{
			TeamID:    types.StringValue(atlasTeam.GetTeamId()),
			RoleNames: roleNames,
		}
	}
	return teams
}

func NewTFLimitsDataSourceModel(ctx context.Context, dataFederationLimits []admin.DataFederationLimit) []*TFLimitModel {
	limits := make([]*TFLimitModel, len(dataFederationLimits))

	for i, dataFederationLimit := range dataFederationLimits {
		limits[i] = &TFLimitModel{
			Name:         types.StringValue(dataFederationLimit.Name),
			Value:        types.Int64Value(dataFederationLimit.Value),
			CurrentUsage: types.Int64PointerValue(dataFederationLimit.CurrentUsage),
			DefaultLimit: types.Int64PointerValue(dataFederationLimit.DefaultLimit),
			MaximumLimit: types.Int64PointerValue(dataFederationLimit.MaximumLimit),
		}
	}

	return limits
}

func NewTFCloudUsersDataSourceModel(ctx context.Context, cloudUsers []admin.GroupUserResponse) []*TFCloudUsersDSModel {
	if len(cloudUsers) == 0 {
		return []*TFCloudUsersDSModel{}
	}
	users := make([]*TFCloudUsersDSModel, len(cloudUsers))
	for i := range cloudUsers {
		cloudUser := &cloudUsers[i]
		roles, _ := types.SetValueFrom(ctx, types.StringType, cloudUser.Roles)
		users[i] = &TFCloudUsersDSModel{
			ID:                  types.StringValue(cloudUser.Id),
			OrgMembershipStatus: types.StringValue(cloudUser.OrgMembershipStatus),
			Roles:               roles,
			Username:            types.StringValue(cloudUser.Username),
			InvitationCreatedAt: types.StringPointerValue(conversion.TimePtrToStringPtr(cloudUser.InvitationCreatedAt)),
			InvitationExpiresAt: types.StringPointerValue(conversion.TimePtrToStringPtr(cloudUser.InvitationExpiresAt)),
			InviterUsername:     types.StringPointerValue(cloudUser.InviterUsername),
			Country:             types.StringPointerValue(cloudUser.Country),
			CreatedAt:           types.StringPointerValue(conversion.TimePtrToStringPtr(cloudUser.CreatedAt)),
			FirstName:           types.StringPointerValue(cloudUser.FirstName),
			LastAuth:            types.StringPointerValue(conversion.TimePtrToStringPtr(cloudUser.LastAuth)),
			LastName:            types.StringPointerValue(cloudUser.LastName),
			MobileNumber:        types.StringPointerValue(cloudUser.MobileNumber),
		}
	}
	return users
}

func NewTFIPAddressesModel(ctx context.Context, ipAddresses *admin.GroupIPAddresses) (types.Object, diag.Diagnostics) {
	clusterIPs := []TFClusterIPsModel{}
	if ipAddresses != nil && ipAddresses.Services != nil {
		clusterIPAddresses := ipAddresses.Services.GetClusters()
		clusterIPs = make([]TFClusterIPsModel, len(clusterIPAddresses))
		for i := range clusterIPAddresses {
			inbound, _ := types.ListValueFrom(ctx, types.StringType, clusterIPAddresses[i].GetInbound())
			outbound, _ := types.ListValueFrom(ctx, types.StringType, clusterIPAddresses[i].GetOutbound())
			clusterIPs[i] = TFClusterIPsModel{
				ClusterName: types.StringPointerValue(clusterIPAddresses[i].ClusterName),
				Inbound:     inbound,
				Outbound:    outbound,
			}
		}
	}
	obj, diags := types.ObjectValueFrom(ctx, IPAddressesObjectType.AttrTypes, TFIPAddressesModel{
		Services: TFServicesModel{
			Clusters: clusterIPs,
		},
	})
	return obj, diags
}

func NewTFProjectResourceModel(ctx context.Context, projectRes *admin.Group, projectProps *AdditionalProperties) (*TFProjectRSModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	if projectRes == nil {
		diags.AddError("Invalid Project Data", "Project data is nil and cannot be processed")
		return nil, diags
	}
	if projectProps == nil {
		diags.AddError("Invalid Project Properties", "Project properties data is nil and cannot be processed")
		return nil, diags
	}
	ipAddressesModel, ipDiags := NewTFIPAddressesModel(ctx, projectProps.IPAddresses)
	diags.Append(ipDiags...)
	if diags.HasError() {
		return nil, diags
	}
	projectPlan := TFProjectRSModel{
		ID:                                 types.StringValue(projectRes.GetId()),
		Name:                               types.StringValue(projectRes.Name),
		OrgID:                              types.StringValue(projectRes.OrgId),
		ClusterCount:                       types.Int64Value(projectRes.ClusterCount),
		RegionUsageRestrictions:            types.StringPointerValue(projectRes.RegionUsageRestrictions),
		Created:                            types.StringValue(conversion.TimeToString(projectRes.Created)),
		WithDefaultAlertsSettings:          types.BoolPointerValue(projectRes.WithDefaultAlertsSettings),
		Teams:                              newTFTeamsResourceModel(ctx, projectProps.Teams),
		Limits:                             newTFLimitsResourceModel(ctx, projectProps.Limits),
		IPAddresses:                        ipAddressesModel,
		Tags:                               conversion.NewTFTags(projectRes.GetTags()),
		IsSlowOperationThresholdingEnabled: types.BoolValue(projectProps.IsSlowOperationThresholdingEnabled),
	}

	projectSettings := projectProps.Settings
	if projectSettings != nil {
		projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled = types.BoolValue(*projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled)
		projectPlan.IsDataExplorerEnabled = types.BoolValue(*projectSettings.IsDataExplorerEnabled)
		projectPlan.IsExtendedStorageSizesEnabled = types.BoolValue(*projectSettings.IsExtendedStorageSizesEnabled)
		projectPlan.IsPerformanceAdvisorEnabled = types.BoolValue(*projectSettings.IsPerformanceAdvisorEnabled)
		projectPlan.IsRealtimePerformancePanelEnabled = types.BoolValue(*projectSettings.IsRealtimePerformancePanelEnabled)
		projectPlan.IsSchemaAdvisorEnabled = types.BoolValue(*projectSettings.IsSchemaAdvisorEnabled)
	}

	return &projectPlan, nil
}

func newTFLimitsResourceModel(ctx context.Context, dataFederationLimits []admin.DataFederationLimit) types.Set {
	limits := make([]TFLimitModel, len(dataFederationLimits))

	for i, dataFederationLimit := range dataFederationLimits {
		limits[i] = TFLimitModel{
			Name:         types.StringValue(dataFederationLimit.Name),
			Value:        types.Int64Value(dataFederationLimit.Value),
			CurrentUsage: types.Int64PointerValue(dataFederationLimit.CurrentUsage),
			DefaultLimit: types.Int64PointerValue(dataFederationLimit.DefaultLimit),
			MaximumLimit: types.Int64PointerValue(dataFederationLimit.MaximumLimit),
		}
	}

	s, _ := types.SetValueFrom(ctx, TfLimitObjectType, limits)
	return s
}

func newTFTeamsResourceModel(ctx context.Context, atlasTeams *admin.PaginatedTeamRole) types.Set {
	if atlasTeams == nil || atlasTeams.GetTotalCount() == 0 {
		return types.SetNull(TfTeamObjectType)
	}
	results := atlasTeams.GetResults()
	teams := make([]TFTeamModel, len(results))
	for i, atlasTeam := range results {
		roleNames, _ := types.SetValueFrom(ctx, types.StringType, atlasTeam.RoleNames)
		teams[i] = TFTeamModel{
			TeamID:    types.StringValue(atlasTeam.GetTeamId()),
			RoleNames: roleNames,
		}
	}

	s, _ := types.SetValueFrom(ctx, TfTeamObjectType, teams)
	return s
}

func NewTeamRoleList(ctx context.Context, teams []TFTeamModel) *[]admin.TeamRole {
	res := make([]admin.TeamRole, len(teams))
	for i, team := range teams {
		roleNames := conversion.TypesSetToString(ctx, team.RoleNames)
		res[i] = admin.TeamRole{
			TeamId:    team.TeamID.ValueStringPointer(),
			RoleNames: &roleNames,
		}
	}
	return &res
}

func NewGroupUpdate(tfProject *TFProjectRSModel, tags *[]admin.ResourceTag) *admin.GroupUpdate {
	return &admin.GroupUpdate{
		Name: tfProject.Name.ValueStringPointer(),
		Tags: tags,
	}
}

func NewTfTeamModelMap(teams []TFTeamModel) map[types.String]TFTeamModel {
	teamsMap := make(map[types.String]TFTeamModel)
	for _, team := range teams {
		teamsMap[team.TeamID] = team
	}
	return teamsMap
}

func NewTfLimitModelMap(limits []TFLimitModel) map[types.String]TFLimitModel {
	limitsMap := make(map[types.String]TFLimitModel)
	for _, limit := range limits {
		limitsMap[limit.Name] = limit
	}
	return limitsMap
}

func SetProjectBool(plan types.Bool, setting **bool) {
	if !plan.IsUnknown() {
		*setting = plan.ValueBoolPointer()
	}
}

func UpdateProjectBool(plan, state types.Bool, setting **bool) bool {
	if plan != state {
		*setting = plan.ValueBoolPointer()
		return true
	}
	return false
}
