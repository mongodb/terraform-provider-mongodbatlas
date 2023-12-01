package project_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

const (
	limitName           = "limitName"
	limitValue          = int64(64)
	limitCurrentUsage   = int64(64)
	limitDefaultLimit   = int64(32)
	limitMaximumLimit   = int64(16)
	projectID           = "projectId"
	projectName         = "projectName"
	projectOrgID        = "orgId"
	projectClusterCount = int64(1)
	clusterCount        = 1
)

var (
	roles        = []string{"GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"}
	roleList, _  = types.ListValueFrom(context.Background(), types.StringType, roles)
	teamRolesSDK = []admin.TeamRole{
		{
			TeamId:    conversion.StringPtr("teamId"),
			RoleNames: roles,
		},
	}
	teamsTF = []*project.TfTeamDSModel{
		{
			TeamID:    types.StringValue("teamId"),
			RoleNames: roleList,
		},
	}
	limitsSDK = []admin.DataFederationLimit{
		{
			Name:         limitName,
			Value:        limitValue,
			CurrentUsage: admin.PtrInt64(limitCurrentUsage),
			DefaultLimit: admin.PtrInt64(limitDefaultLimit),
			MaximumLimit: admin.PtrInt64(limitMaximumLimit),
		},
	}
	limitsTF = []*project.TfLimitModel{
		{
			Name:         types.StringValue(limitName),
			Value:        types.Int64Value(limitValue),
			CurrentUsage: types.Int64Value(limitCurrentUsage),
			DefaultLimit: types.Int64Value(limitDefaultLimit),
			MaximumLimit: types.Int64Value(limitMaximumLimit),
		},
	}
	projectSDK = admin.Group{
		Id:           admin.PtrString(projectID),
		Name:         projectName,
		OrgId:        projectOrgID,
		ClusterCount: int64(clusterCount),
	}
	projectSettingsSDK = admin.GroupSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: admin.PtrBool(true),
		IsDataExplorerEnabled:                       admin.PtrBool(true),
		IsExtendedStorageSizesEnabled:               admin.PtrBool(true),
		IsPerformanceAdvisorEnabled:                 admin.PtrBool(true),
		IsRealtimePerformancePanelEnabled:           admin.PtrBool(true),
		IsSchemaAdvisorEnabled:                      admin.PtrBool(true),
	}
)

func TestTeamsDataSourceSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name              string
		paginatedTeamRole *admin.PaginatedTeamRole
		expectedTFModel   []*project.TfTeamDSModel
	}{
		{
			name: "TeamRole",
			paginatedTeamRole: &admin.PaginatedTeamRole{
				TotalCount: conversion.IntPtr(0),
			}, // not setting explicitly expected result because we expect it to be nil
		},
		{
			name: "Complete TeamRole",
			paginatedTeamRole: &admin.PaginatedTeamRole{
				Results:    teamRolesSDK,
				TotalCount: conversion.IntPtr(1),
			},
			expectedTFModel: teamsTF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTFTeamsDataSourceModel(context.Background(), tc.paginatedTeamRole)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

func TestLimitsDataSourceSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                 string
		dataFederationLimits []admin.DataFederationLimit
		expectedTFModel      []*project.TfLimitModel
	}{
		{
			name:                 "Limit",
			dataFederationLimits: limitsSDK,
			expectedTFModel:      limitsTF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTFLimitsDataSourceModel(context.Background(), tc.dataFederationLimits)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}

func TestProjectDataSourceSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                 string
		project              *admin.Group
		teams                *admin.PaginatedTeamRole
		projectSettings      *admin.GroupSettings
		dataFederationLimits []admin.DataFederationLimit
		expectedTFModel      project.TfProjectDSModel
	}{
		{
			name:    "Project",
			project: &projectSDK,
			teams: &admin.PaginatedTeamRole{
				Results:    teamRolesSDK,
				TotalCount: conversion.IntPtr(1),
			},
			projectSettings:      &projectSettingsSDK,
			dataFederationLimits: limitsSDK,
			expectedTFModel: project.TfProjectDSModel{

				ID:           types.StringValue(projectID),
				ProjectID:    types.StringValue(projectID),
				Name:         types.StringValue(projectName),
				OrgID:        types.StringValue(projectOrgID),
				ClusterCount: types.Int64Value(clusterCount),
				IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(true),
				IsDataExplorerEnabled:                       types.BoolValue(true),
				IsExtendedStorageSizesEnabled:               types.BoolValue(true),
				IsPerformanceAdvisorEnabled:                 types.BoolValue(true),
				IsRealtimePerformancePanelEnabled:           types.BoolValue(true),
				IsSchemaAdvisorEnabled:                      types.BoolValue(true),
				Teams:                                       teamsTF,
				Limits:                                      limitsTF,
				Created:                                     types.StringValue("0001-01-01T00:00:00Z"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTFProjectDataSourceModel(context.Background(), tc.project, tc.teams, tc.projectSettings, tc.dataFederationLimits)
			if !reflect.DeepEqual(resultModel, tc.expectedTFModel) {
				t.Errorf("created terraform model did not match expected output")
			}
		})
	}
}
