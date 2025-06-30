package project

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312004/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

func NewTFProjectDataSourceModel(ctx context.Context, project *admin.Group, projectProps *AdditionalProperties) (*TFProjectDSModel, diag.Diagnostics) {
	ipAddressesModel, diags := NewTFIPAddressesModel(ctx, projectProps.IPAddresses)
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
	if atlasTeams.GetTotalCount() == 0 {
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
		return nil
	}
	users := make([]*TFCloudUsersDSModel, len(cloudUsers))
	for i := range cloudUsers {
		cloudUser := &cloudUsers[i]
		roles, _ := types.ListValueFrom(ctx, types.StringType, cloudUser.Roles)
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
	ipAddressesModel, diags := NewTFIPAddressesModel(ctx, projectProps.IPAddresses)
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
