package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

func NewTFProjectDataSourceModel(ctx context.Context, project *admin.Group,
	teams *admin.PaginatedTeamRole, projectSettings *admin.GroupSettings, limits []admin.DataFederationLimit) TfProjectDSModel {
	return TfProjectDSModel{
		ID:           types.StringValue(project.GetId()),
		ProjectID:    types.StringValue(project.GetId()),
		Name:         types.StringValue(project.Name),
		OrgID:        types.StringValue(project.OrgId),
		ClusterCount: types.Int64Value(project.ClusterCount),
		Created:      types.StringValue(conversion.TimeToString(project.Created)),
		IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(*projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled),
		IsDataExplorerEnabled:                       types.BoolValue(*projectSettings.IsDataExplorerEnabled),
		IsExtendedStorageSizesEnabled:               types.BoolValue(*projectSettings.IsExtendedStorageSizesEnabled),
		IsPerformanceAdvisorEnabled:                 types.BoolValue(*projectSettings.IsPerformanceAdvisorEnabled),
		IsRealtimePerformancePanelEnabled:           types.BoolValue(*projectSettings.IsRealtimePerformancePanelEnabled),
		IsSchemaAdvisorEnabled:                      types.BoolValue(*projectSettings.IsSchemaAdvisorEnabled),
		Teams:                                       NewTFTeamsDataSourceModel(ctx, teams),
		Limits:                                      NewTFLimitsDataSourceModel(ctx, limits),
	}
}

func NewTFTeamsDataSourceModel(ctx context.Context, atlasTeams *admin.PaginatedTeamRole) []*TfTeamDSModel {
	if atlasTeams.GetTotalCount() == 0 {
		return nil
	}
	teams := make([]*TfTeamDSModel, len(atlasTeams.Results))

	for i, atlasTeam := range atlasTeams.Results {
		roleNames, _ := types.ListValueFrom(ctx, types.StringType, atlasTeam.RoleNames)

		teams[i] = &TfTeamDSModel{
			TeamID:    types.StringValue(atlasTeam.GetTeamId()),
			RoleNames: roleNames,
		}
	}
	return teams
}

func NewTFLimitsDataSourceModel(ctx context.Context, dataFederationLimits []admin.DataFederationLimit) []*TfLimitModel {
	limits := make([]*TfLimitModel, len(dataFederationLimits))

	for i, dataFederationLimit := range dataFederationLimits {
		limits[i] = &TfLimitModel{
			Name:         types.StringValue(dataFederationLimit.Name),
			Value:        types.Int64Value(dataFederationLimit.Value),
			CurrentUsage: types.Int64PointerValue(dataFederationLimit.CurrentUsage),
			DefaultLimit: types.Int64PointerValue(dataFederationLimit.DefaultLimit),
			MaximumLimit: types.Int64PointerValue(dataFederationLimit.MaximumLimit),
		}
	}

	return limits
}

func NewTFProjectResourceModel(ctx context.Context, projectRes *admin.Group,
	teams *admin.PaginatedTeamRole, projectSettings *admin.GroupSettings, limits []admin.DataFederationLimit) *TfProjectRSModel {
	projectPlan := TfProjectRSModel{
		ID:                        types.StringValue(projectRes.GetId()),
		Name:                      types.StringValue(projectRes.Name),
		OrgID:                     types.StringValue(projectRes.OrgId),
		ClusterCount:              types.Int64Value(projectRes.ClusterCount),
		Created:                   types.StringValue(conversion.TimeToString(projectRes.Created)),
		WithDefaultAlertsSettings: types.BoolPointerValue(projectRes.WithDefaultAlertsSettings),
		Teams:                     newTFTeamsResourceModel(ctx, teams),
		Limits:                    newTFLimitsResourceModel(ctx, limits),
	}

	if projectSettings != nil {
		projectPlan.IsCollectDatabaseSpecificsStatisticsEnabled = types.BoolValue(*projectSettings.IsCollectDatabaseSpecificsStatisticsEnabled)
		projectPlan.IsDataExplorerEnabled = types.BoolValue(*projectSettings.IsDataExplorerEnabled)
		projectPlan.IsExtendedStorageSizesEnabled = types.BoolValue(*projectSettings.IsExtendedStorageSizesEnabled)
		projectPlan.IsPerformanceAdvisorEnabled = types.BoolValue(*projectSettings.IsPerformanceAdvisorEnabled)
		projectPlan.IsRealtimePerformancePanelEnabled = types.BoolValue(*projectSettings.IsRealtimePerformancePanelEnabled)
		projectPlan.IsSchemaAdvisorEnabled = types.BoolValue(*projectSettings.IsSchemaAdvisorEnabled)
	}

	return &projectPlan
}

func newTFLimitsResourceModel(ctx context.Context, dataFederationLimits []admin.DataFederationLimit) types.Set {
	limits := make([]TfLimitModel, len(dataFederationLimits))

	for i, dataFederationLimit := range dataFederationLimits {
		limits[i] = TfLimitModel{
			Name:         types.StringValue(dataFederationLimit.Name),
			Value:        types.Int64Value(dataFederationLimit.Value),
			CurrentUsage: types.Int64PointerValue(dataFederationLimit.CurrentUsage),
			DefaultLimit: types.Int64PointerValue(dataFederationLimit.DefaultLimit),
			MaximumLimit: types.Int64PointerValue(dataFederationLimit.MaximumLimit),
		}
	}

	s, _ := types.SetValueFrom(ctx, tfLimitObjectType, limits)
	return s
}

func newTFTeamsResourceModel(ctx context.Context, atlasTeams *admin.PaginatedTeamRole) types.Set {
	teams := make([]TfTeamModel, len(atlasTeams.Results))

	for i, atlasTeam := range atlasTeams.Results {
		roleNames, _ := types.SetValueFrom(ctx, types.StringType, atlasTeam.RoleNames)

		teams[i] = TfTeamModel{
			TeamID:    types.StringValue(atlasTeam.GetTeamId()),
			RoleNames: roleNames,
		}
	}

	s, _ := types.SetValueFrom(ctx, tfTeamObjectType, teams)
	return s
}

func NewTeamRoleList(ctx context.Context, teams []TfTeamModel) *[]admin.TeamRole {
	res := make([]admin.TeamRole, len(teams))

	for i, team := range teams {
		res[i] = admin.TeamRole{
			TeamId:    team.TeamID.ValueStringPointer(),
			RoleNames: conversion.TypesSetToString(ctx, team.RoleNames),
		}
	}
	return &res
}

func NewGroupName(tfProject *TfProjectRSModel) *admin.GroupName {
	return &admin.GroupName{
		Name: tfProject.Name.ValueStringPointer(),
	}
}

func NewTfTeamModelMap(teams []TfTeamModel) map[types.String]TfTeamModel {
	teamsMap := make(map[types.String]TfTeamModel)
	for _, team := range teams {
		teamsMap[team.TeamID] = team
	}
	return teamsMap
}

func NewTfLimitModelMap(limits []TfLimitModel) map[types.String]TfLimitModel {
	limitsMap := make(map[types.String]TfLimitModel)
	for _, limit := range limits {
		limitsMap[limit.Name] = limit
	}
	return limitsMap
}
