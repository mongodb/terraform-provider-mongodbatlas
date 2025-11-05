package project_test

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
)

const (
	limitName               = "limitName"
	limitValue              = int64(64)
	limitCurrentUsage       = int64(64)
	limitDefaultLimit       = int64(32)
	limitMaximumLimit       = int64(16)
	projectID               = "projectId"
	projectName             = "projectName"
	projectOrgID            = "orgId"
	projectClusterCount     = int64(1)
	clusterCount            = 1
	regionUsageRestrictions = "GOV_REGIONS_ONLY"
	userOrgMembershipStatus = "ACTIVE"
	country                 = "US"
	inviterUsername         = ""
	mobileNumber            = ""
)

var (
	roles              = []string{"GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"}
	roleList, _        = types.ListValueFrom(context.Background(), types.StringType, roles)
	roleSet, _         = types.SetValueFrom(context.Background(), types.StringType, roles)
	ipAddresses        = []string{"13.13.13.13"}
	ipAddressesList, _ = types.ListValueFrom(context.Background(), types.StringType, ipAddresses)
	empptyTFList, _    = types.ListValueFrom(context.Background(), types.StringType, []string{})
	teamRolesSDK       = []admin.TeamRole{
		{
			TeamId:    conversion.StringPtr("teamId"),
			RoleNames: &roles,
		},
	}
	teamsDSTF = []*project.TFTeamDSModel{
		{
			TeamID:    types.StringValue("teamId"),
			RoleNames: roleList,
		},
	}
	teamsTFSet, _ = types.SetValueFrom(context.Background(), project.TfTeamObjectType, []project.TFTeamModel{
		{
			TeamID:    types.StringValue("teamId"),
			RoleNames: roleSet,
		},
	})
	limitsSDK = []admin.DataFederationLimit{
		{
			Name:         limitName,
			Value:        limitValue,
			CurrentUsage: admin.PtrInt64(limitCurrentUsage),
			DefaultLimit: admin.PtrInt64(limitDefaultLimit),
			MaximumLimit: admin.PtrInt64(limitMaximumLimit),
		},
	}
	limitsTF = []*project.TFLimitModel{
		{
			Name:         types.StringValue(limitName),
			Value:        types.Int64Value(limitValue),
			CurrentUsage: types.Int64Value(limitCurrentUsage),
			DefaultLimit: types.Int64Value(limitDefaultLimit),
			MaximumLimit: types.Int64Value(limitMaximumLimit),
		},
	}
	limitsTFSet, _ = types.SetValueFrom(context.Background(), project.TfLimitObjectType, []project.TFLimitModel{
		*limitsTF[0],
	})

	usersSDK = []admin.GroupUserResponse{
		{
			Id:                  "user-id-1",
			Username:            "user1@example.com",
			FirstName:           admin.PtrString("FirstName1"),
			LastName:            admin.PtrString("LastName1"),
			Roles:               roles,
			InvitationCreatedAt: nil,
			InvitationExpiresAt: nil,
			InviterUsername:     admin.PtrString(inviterUsername),
			OrgMembershipStatus: userOrgMembershipStatus,
			Country:             admin.PtrString("US"),
			CreatedAt:           admin.PtrTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			LastAuth:            admin.PtrTime(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
			MobileNumber:        admin.PtrString(mobileNumber),
		},
		{
			Id:                  "user-id-2",
			Username:            "user2@example.com",
			FirstName:           admin.PtrString("FirstName2"),
			LastName:            admin.PtrString("LastName2"),
			Roles:               roles,
			InvitationCreatedAt: nil,
			InvitationExpiresAt: nil,
			InviterUsername:     admin.PtrString(inviterUsername),
			OrgMembershipStatus: userOrgMembershipStatus,
			Country:             admin.PtrString(country),
			CreatedAt:           admin.PtrTime(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
			LastAuth:            admin.PtrTime(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
			MobileNumber:        admin.PtrString(mobileNumber),
		},
	}
	usersTF = []*project.TFCloudUsersDSModel{
		{
			ID:                  types.StringValue("user-id-1"),
			Username:            types.StringValue("user1@example.com"),
			FirstName:           types.StringValue("FirstName1"),
			LastName:            types.StringValue("LastName1"),
			Roles:               roleSet,
			InvitationCreatedAt: types.StringNull(),
			InvitationExpiresAt: types.StringNull(),
			InviterUsername:     types.StringValue(inviterUsername),
			OrgMembershipStatus: types.StringValue(userOrgMembershipStatus),
			Country:             types.StringValue(country),
			CreatedAt:           types.StringValue("2025-01-01T00:00:00Z"),
			LastAuth:            types.StringValue("2025-01-02T00:00:00Z"),
			MobileNumber:        types.StringValue(mobileNumber),
		},
		{
			ID:                  types.StringValue("user-id-2"),
			Username:            types.StringValue("user2@example.com"),
			FirstName:           types.StringValue("FirstName2"),
			LastName:            types.StringValue("LastName2"),
			Roles:               roleSet,
			InvitationCreatedAt: types.StringNull(),
			InvitationExpiresAt: types.StringNull(),
			InviterUsername:     types.StringValue(inviterUsername),
			OrgMembershipStatus: types.StringValue(userOrgMembershipStatus),
			Country:             types.StringValue(country),
			CreatedAt:           types.StringValue("2025-01-01T00:00:00Z"),
			LastAuth:            types.StringValue("2025-01-02T00:00:00Z"),
			MobileNumber:        types.StringValue(mobileNumber),
		},
	}

	ipAddressesTF, _ = types.ObjectValueFrom(context.Background(), project.IPAddressesObjectType.AttrTypes, project.TFIPAddressesModel{
		Services: project.TFServicesModel{
			Clusters: []project.TFClusterIPsModel{
				{
					Inbound:     ipAddressesList,
					Outbound:    ipAddressesList,
					ClusterName: types.StringValue("Cluster0"),
				},
			},
		},
	})
	IPAddressesNoClusterTF, _ = types.ObjectValueFrom(context.Background(), project.IPAddressesObjectType.AttrTypes, project.TFIPAddressesModel{
		Services: project.TFServicesModel{
			Clusters: []project.TFClusterIPsModel{},
		},
	})
	IPAddressesWithClusterNoIPsTF, _ = types.ObjectValueFrom(context.Background(), project.IPAddressesObjectType.AttrTypes, project.TFIPAddressesModel{
		Services: project.TFServicesModel{
			Clusters: []project.TFClusterIPsModel{
				{
					Inbound:     empptyTFList,
					Outbound:    empptyTFList,
					ClusterName: types.StringValue("Cluster0"),
				},
			},
		},
	})
	projectSDK = admin.Group{
		Id:           admin.PtrString(projectID),
		Name:         projectName,
		OrgId:        projectOrgID,
		ClusterCount: int64(clusterCount),
	}
	projectGovSDK = admin.Group{
		Id:                      admin.PtrString(projectID),
		Name:                    projectName,
		OrgId:                   projectOrgID,
		ClusterCount:            int64(clusterCount),
		RegionUsageRestrictions: admin.PtrString(regionUsageRestrictions),
	}
	projectSettingsSDK = admin.GroupSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: admin.PtrBool(true),
		IsDataExplorerEnabled:                       admin.PtrBool(true),
		IsExtendedStorageSizesEnabled:               admin.PtrBool(true),
		IsPerformanceAdvisorEnabled:                 admin.PtrBool(true),
		IsRealtimePerformancePanelEnabled:           admin.PtrBool(true),
		IsSchemaAdvisorEnabled:                      admin.PtrBool(true),
	}
	IPAddressesSDK = admin.GroupIPAddresses{
		GroupId: admin.PtrString(projectID),
		Services: &admin.GroupService{
			Clusters: &[]admin.ClusterIPAddresses{
				{
					Inbound:     &[]string{"13.13.13.13"},
					Outbound:    &[]string{"13.13.13.13"},
					ClusterName: admin.PtrString("Cluster0"),
				},
			},
		},
	}
	IPAddressesWithClusterNoIPsSDK = admin.GroupIPAddresses{
		GroupId: admin.PtrString(projectID),
		Services: &admin.GroupService{
			Clusters: &[]admin.ClusterIPAddresses{
				{
					Inbound:     &[]string{},
					Outbound:    &[]string{},
					ClusterName: admin.PtrString("Cluster0"),
				},
			},
		},
	}
)

func TestTeamsDataSourceSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name              string
		paginatedTeamRole *admin.PaginatedTeamRole
		expectedTFModel   []*project.TFTeamDSModel
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
				Results:    &teamRolesSDK,
				TotalCount: conversion.IntPtr(1),
			},
			expectedTFModel: teamsDSTF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTFTeamsDataSourceModel(t.Context(), tc.paginatedTeamRole)
			assert.Equal(t, tc.expectedTFModel, resultModel)
		})
	}
}

func TestLimitsDataSourceSDKToTFModel(t *testing.T) {
	testCases := []struct {
		name                 string
		dataFederationLimits []admin.DataFederationLimit
		expectedTFModel      []*project.TFLimitModel
	}{
		{
			name:                 "Limit",
			dataFederationLimits: limitsSDK,
			expectedTFModel:      limitsTF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTFLimitsDataSourceModel(t.Context(), tc.dataFederationLimits)
			assert.Equal(t, tc.expectedTFModel, resultModel)
		})
	}
}

func TestUsersDataSourceSDKToDataSourceTFModel(t *testing.T) {
	testCases := []struct {
		name            string
		users           []admin.GroupUserResponse
		expectedTFModel []*project.TFCloudUsersDSModel
	}{
		{
			name:            "Users",
			users:           usersSDK,
			expectedTFModel: usersTF,
		},
		{
			name:            "Empty Users",
			users:           []admin.GroupUserResponse{},
			expectedTFModel: []*project.TFCloudUsersDSModel{},
		},
		{
			name:            "Nil Users",
			users:           nil,
			expectedTFModel: []*project.TFCloudUsersDSModel{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTFCloudUsersDataSourceModel(t.Context(), tc.users)
			assert.Equal(t, tc.expectedTFModel, resultModel)
		})
	}
}

func TestProjectDataSourceSDKToDataSourceTFModel(t *testing.T) {
	testCases := []struct {
		name            string
		project         *admin.Group
		projectProps    project.AdditionalProperties
		expectedTFModel project.TFProjectDSModel
	}{
		{
			name:    "Project",
			project: &projectSDK,
			projectProps: project.AdditionalProperties{
				Teams: &admin.PaginatedTeamRole{
					Results:    &teamRolesSDK,
					TotalCount: conversion.IntPtr(1),
				},
				Settings:    &projectSettingsSDK,
				IPAddresses: &IPAddressesSDK,
				Limits:      limitsSDK,
				Users:       usersSDK,
			},
			expectedTFModel: project.TFProjectDSModel{

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
				IsSlowOperationThresholdingEnabled:          types.BoolValue(false),
				Teams:                                       teamsDSTF,
				Limits:                                      limitsTF,
				Users:                                       usersTF,
				IPAddresses:                                 ipAddressesTF,
				Created:                                     types.StringValue("0001-01-01T00:00:00Z"),
				Tags:                                        types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
		},
		{
			name:    "ProjectGov",
			project: &projectGovSDK,
			projectProps: project.AdditionalProperties{
				Teams: &admin.PaginatedTeamRole{
					Results:    &teamRolesSDK,
					TotalCount: conversion.IntPtr(1),
				},
				Settings:                           &projectSettingsSDK,
				IPAddresses:                        &IPAddressesSDK,
				Limits:                             limitsSDK,
				Users:                              usersSDK,
				IsSlowOperationThresholdingEnabled: true,
			},
			expectedTFModel: project.TFProjectDSModel{

				ID:                      types.StringValue(projectID),
				ProjectID:               types.StringValue(projectID),
				Name:                    types.StringValue(projectName),
				OrgID:                   types.StringValue(projectOrgID),
				ClusterCount:            types.Int64Value(clusterCount),
				RegionUsageRestrictions: types.StringValue(regionUsageRestrictions),
				IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(true),
				IsDataExplorerEnabled:                       types.BoolValue(true),
				IsExtendedStorageSizesEnabled:               types.BoolValue(true),
				IsPerformanceAdvisorEnabled:                 types.BoolValue(true),
				IsRealtimePerformancePanelEnabled:           types.BoolValue(true),
				IsSchemaAdvisorEnabled:                      types.BoolValue(true),
				IsSlowOperationThresholdingEnabled:          types.BoolValue(true),
				Teams:                                       teamsDSTF,
				Limits:                                      limitsTF,
				Users:                                       usersTF,
				IPAddresses:                                 ipAddressesTF,
				Created:                                     types.StringValue("0001-01-01T00:00:00Z"),
				Tags:                                        types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := project.NewTFProjectDataSourceModel(t.Context(), tc.project, &tc.projectProps)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, *resultModel)
		})
	}
}

func TestProjectDataSourceSDKToResourceTFModel(t *testing.T) {
	testCases := []struct {
		name            string
		project         *admin.Group
		projectProps    project.AdditionalProperties
		expectedTFModel project.TFProjectRSModel
	}{
		{
			name:    "Project",
			project: &projectSDK,
			projectProps: project.AdditionalProperties{
				Teams: &admin.PaginatedTeamRole{
					Results:    &teamRolesSDK,
					TotalCount: conversion.IntPtr(1),
				},
				Settings:    &projectSettingsSDK,
				IPAddresses: &IPAddressesSDK,
				Limits:      limitsSDK,
			},
			expectedTFModel: project.TFProjectRSModel{

				ID:           types.StringValue(projectID),
				Name:         types.StringValue(projectName),
				OrgID:        types.StringValue(projectOrgID),
				ClusterCount: types.Int64Value(clusterCount),
				IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(true),
				IsDataExplorerEnabled:                       types.BoolValue(true),
				IsExtendedStorageSizesEnabled:               types.BoolValue(true),
				IsPerformanceAdvisorEnabled:                 types.BoolValue(true),
				IsRealtimePerformancePanelEnabled:           types.BoolValue(true),
				IsSchemaAdvisorEnabled:                      types.BoolValue(true),
				IsSlowOperationThresholdingEnabled:          types.BoolValue(false),
				Teams:                                       teamsTFSet,
				Limits:                                      limitsTFSet,
				IPAddresses:                                 ipAddressesTF,
				Created:                                     types.StringValue("0001-01-01T00:00:00Z"),
				Tags:                                        types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
		},
		{
			name:    "ProjectGov",
			project: &projectGovSDK,
			projectProps: project.AdditionalProperties{
				Teams: &admin.PaginatedTeamRole{
					Results:    &teamRolesSDK,
					TotalCount: conversion.IntPtr(1),
				},
				Settings:    &projectSettingsSDK,
				IPAddresses: &IPAddressesSDK,
				Limits:      limitsSDK,
			},
			expectedTFModel: project.TFProjectRSModel{

				ID:                      types.StringValue(projectID),
				Name:                    types.StringValue(projectName),
				OrgID:                   types.StringValue(projectOrgID),
				ClusterCount:            types.Int64Value(clusterCount),
				RegionUsageRestrictions: types.StringValue(regionUsageRestrictions),
				IsCollectDatabaseSpecificsStatisticsEnabled: types.BoolValue(true),
				IsDataExplorerEnabled:                       types.BoolValue(true),
				IsExtendedStorageSizesEnabled:               types.BoolValue(true),
				IsPerformanceAdvisorEnabled:                 types.BoolValue(true),
				IsRealtimePerformancePanelEnabled:           types.BoolValue(true),
				IsSchemaAdvisorEnabled:                      types.BoolValue(true),
				IsSlowOperationThresholdingEnabled:          types.BoolValue(false),
				Teams:                                       teamsTFSet,
				Limits:                                      limitsTFSet,
				IPAddresses:                                 ipAddressesTF,
				Created:                                     types.StringValue("0001-01-01T00:00:00Z"),
				Tags:                                        types.MapValueMust(types.StringType, map[string]attr.Value{}),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := project.NewTFProjectResourceModel(t.Context(), tc.project, &tc.projectProps)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, *resultModel)
		})
	}
}

func TestTeamRoleListTFtoSDK(t *testing.T) {
	var rolesSet, _ = types.SetValueFrom(t.Context(), types.StringType, roles)
	teamsTF := []project.TFTeamModel{
		{
			TeamID:    types.StringValue("teamId"),
			RoleNames: rolesSet,
		},
	}
	testCases := []struct {
		name           string
		expectedResult *[]admin.TeamRole
		teamRolesTF    []project.TFTeamModel
	}{
		{
			name:           "Team roles",
			teamRolesTF:    teamsTF,
			expectedResult: &teamRolesSDK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTeamRoleList(t.Context(), tc.teamRolesTF)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestTeamModelMapTF(t *testing.T) {
	teams := []project.TFTeamModel{
		{
			TeamID: types.StringValue("id1"),
		},
		{
			TeamID: types.StringValue("id2"),
		},
	}
	testCases := []struct {
		name           string
		expectedResult map[types.String]project.TFTeamModel
		teamRolesTF    []project.TFTeamModel
	}{
		{
			name:        "Team roles",
			teamRolesTF: teams,
			expectedResult: map[types.String]project.TFTeamModel{
				types.StringValue("id1"): teams[0],
				types.StringValue("id2"): teams[1],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTfTeamModelMap(tc.teamRolesTF)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestLimitModelMapTF(t *testing.T) {
	limits := []project.TFLimitModel{
		{
			Name: types.StringValue("limit1"),
		},
		{
			Name: types.StringValue("limit2"),
		},
	}
	testCases := []struct {
		name           string
		expectedResult map[types.String]project.TFLimitModel
		limitsTF       []project.TFLimitModel
	}{
		{
			name:     "Limits",
			limitsTF: limits,
			expectedResult: map[types.String]project.TFLimitModel{
				types.StringValue("limit1"): limits[0],
				types.StringValue("limit2"): limits[1],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.NewTfLimitModelMap(tc.limitsTF)
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestIPAddressesModelToTF(t *testing.T) {
	testCases := []struct {
		name           string
		sdkModel       *admin.GroupIPAddresses
		expectedResult types.Object
	}{
		{
			name:           "No response",
			sdkModel:       nil,
			expectedResult: IPAddressesNoClusterTF,
		},
		{
			name: "Empty response when no clusters are created",
			sdkModel: &admin.GroupIPAddresses{
				GroupId: admin.PtrString(projectID),
				Services: &admin.GroupService{
					Clusters: &[]admin.ClusterIPAddresses{},
				},
			},
			expectedResult: IPAddressesNoClusterTF,
		},
		{
			name:           "One cluster with empty IP lists",
			sdkModel:       &IPAddressesWithClusterNoIPsSDK,
			expectedResult: IPAddressesWithClusterNoIPsTF,
		},
		{
			name:           "Full response",
			sdkModel:       &IPAddressesSDK,
			expectedResult: ipAddressesTF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel, diags := project.NewTFIPAddressesModel(t.Context(), tc.sdkModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedResult, resultModel)
		})
	}
}

func TestSetProjectBool(t *testing.T) {
	testCases := []struct {
		name     string
		plan     types.Bool
		expected bool
	}{
		{"unknown", types.BoolUnknown(), false},
		{"false", types.BoolValue(false), false},
		{"true", types.BoolValue(true), true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setting := new(bool)
			project.SetProjectBool(tc.plan, &setting)
			assert.Equal(t, tc.expected, *setting)
		})
	}
}

func TestUpdateProjectBool(t *testing.T) {
	testCases := []struct {
		name            string
		plan            types.Bool
		state           types.Bool
		expectedSetting bool
		expected        bool
	}{
		{"same state unknown", types.BoolUnknown(), types.BoolUnknown(), false, false},
		{"same state false", types.BoolValue(false), types.BoolValue(false), false, false},
		{"same state true", types.BoolValue(true), types.BoolValue(true), false, false},
		{"different state unknown", types.BoolUnknown(), types.BoolValue(false), false, true},
		{"different state false", types.BoolValue(false), types.BoolValue(true), false, true},
		{"different state true", types.BoolValue(true), types.BoolValue(false), true, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setting := new(bool)
			assert.Equal(t, tc.expected, project.UpdateProjectBool(tc.plan, tc.state, &setting))
			assert.Equal(t, tc.expectedSetting, *setting)
		})
	}
}
