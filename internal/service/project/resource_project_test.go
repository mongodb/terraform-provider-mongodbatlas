package project_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mocksvc"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

var (
	name             = types.StringValue("sameName")
	diffName         = types.StringValue("diffName")
	projectStateName = project.TFProjectRSModel{
		Name: name,
	}
	projectStateNameDiff = project.TFProjectRSModel{
		Name: diffName,
	}
	dummyProjectID = "6575af27f93c7a6a4b50b239"
)

func TestGetProjectPropsFromAPI(t *testing.T) {
	successfulTeamRoleResponse := TeamRoleResponse{
		TeamRole: &admin.PaginatedTeamRole{},
		Err:      nil,
	}
	successfulLimitsResponse := LimitsResponse{
		Limits: []admin.DataFederationLimit{},
		Err:    nil,
	}
	successfulGroupSettingsResponse := GroupSettingsResponse{
		GroupSettings: &admin.GroupSettings{},
		Err:           nil,
	}
	testCases := []struct {
		teamRoleReponse     TeamRoleResponse
		groupResponse       GroupSettingsResponse
		ipAddressesResponse IPAddressesResponse
		name                string
		limitResponse       LimitsResponse
		expectedError       bool
	}{
		{
			name:            "Successful",
			teamRoleReponse: successfulTeamRoleResponse,
			limitResponse:   successfulLimitsResponse,
			groupResponse:   successfulGroupSettingsResponse,
			expectedError:   false,
		},
		{
			name: "Fail to get project's teams assigned ",
			teamRoleReponse: TeamRoleResponse{
				TeamRole:     nil,
				HTTPResponse: &http.Response{StatusCode: 503},
				Err:          errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
		{
			name:            "Fail to get project's limits",
			teamRoleReponse: successfulTeamRoleResponse,
			limitResponse: LimitsResponse{
				Limits:       nil,
				HTTPResponse: &http.Response{StatusCode: 503},
				Err:          errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
		{
			name:            "Fail to get project's settings",
			teamRoleReponse: successfulTeamRoleResponse,
			limitResponse:   successfulLimitsResponse,
			groupResponse: GroupSettingsResponse{
				GroupSettings: nil,
				HTTPResponse:  &http.Response{StatusCode: 503},
				Err:           errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
		{
			name:            "Fail to get project's ip addresses",
			teamRoleReponse: successfulTeamRoleResponse,
			limitResponse:   successfulLimitsResponse,
			groupResponse:   successfulGroupSettingsResponse,
			ipAddressesResponse: IPAddressesResponse{
				IPAddresses:  nil,
				HTTPResponse: &http.Response{StatusCode: 503},
				Err:          errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teamsMock := mocksvc.NewTeamsApi(t)
			projectsMock := mocksvc.NewProjectsApi(t)

			teamsMock.EXPECT().ListProjectTeams(mock.Anything, mock.Anything).Return(admin.ListProjectTeamsApiRequest{ApiService: teamsMock})
			teamsMock.EXPECT().ListProjectTeamsExecute(mock.Anything).Return(tc.teamRoleReponse.TeamRole, tc.teamRoleReponse.HTTPResponse, tc.teamRoleReponse.Err)

			projectsMock.EXPECT().ListProjectLimits(mock.Anything, mock.Anything).Return(admin.ListProjectLimitsApiRequest{ApiService: projectsMock}).Maybe()
			projectsMock.EXPECT().ListProjectLimitsExecute(mock.Anything).Return(tc.limitResponse.Limits, tc.limitResponse.HTTPResponse, tc.limitResponse.Err).Maybe()

			projectsMock.EXPECT().GetProjectSettings(mock.Anything, mock.Anything).Return(admin.GetProjectSettingsApiRequest{ApiService: projectsMock}).Maybe()
			projectsMock.EXPECT().GetProjectSettingsExecute(mock.Anything).Return(tc.groupResponse.GroupSettings, tc.groupResponse.HTTPResponse, tc.groupResponse.Err).Maybe()

			projectsMock.EXPECT().ReturnAllIPAddresses(mock.Anything, mock.Anything).Return(admin.ReturnAllIPAddressesApiRequest{ApiService: projectsMock}).Maybe()
			projectsMock.EXPECT().ReturnAllIPAddressesExecute(mock.Anything).Return(tc.ipAddressesResponse.IPAddresses, tc.ipAddressesResponse.HTTPResponse, tc.ipAddressesResponse.Err).Maybe()

			_, err := project.GetProjectPropsFromAPI(context.Background(), projectsMock, teamsMock, dummyProjectID)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestFilterUserDefinedLimits(t *testing.T) {
	testCases := []struct {
		name           string
		allAtlasLimits []admin.DataFederationLimit
		tfLimits       []project.TFLimitModel
		expectedResult []admin.DataFederationLimit
	}{
		{
			name: "FilterUserDefinedLimits",
			allAtlasLimits: []admin.DataFederationLimit{
				createDataFederationLimit("1"),
				createDataFederationLimit("2"),
				createDataFederationLimit("3"),
			},
			tfLimits: []project.TFLimitModel{
				{
					Name: types.StringValue("1"),
				},
				{
					Name: types.StringValue("2"),
				},
			},
			expectedResult: []admin.DataFederationLimit{
				createDataFederationLimit("1"),
				createDataFederationLimit("2"),
			},
		},
		{
			name: "FilterUserDefinedLimits",
			allAtlasLimits: []admin.DataFederationLimit{
				createDataFederationLimit("1"),
			},
			tfLimits:       []project.TFLimitModel{},
			expectedResult: []admin.DataFederationLimit{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultModel := project.FilterUserDefinedLimits(tc.allAtlasLimits, tc.tfLimits)
			if !reflect.DeepEqual(resultModel, tc.expectedResult) {
				t.Errorf("Filtered DataFederationlimit did not match expected output")
			}
		})
	}
}

func TestUpdateProject(t *testing.T) {
	testCases := []struct {
		name          string
		mockResponses ProjectResponse
		projectState  project.TFProjectRSModel
		projectPlan   project.TFProjectRSModel
		expectedError bool
	}{
		{
			name:         "Successful update",
			projectState: projectStateName,
			projectPlan:  projectStateNameDiff,
			mockResponses: ProjectResponse{
				Err: nil,
			},
			expectedError: false,
		},
		{
			name:         "Same project names; No update",
			projectState: projectStateName,
			projectPlan:  projectStateName,
			mockResponses: ProjectResponse{
				Err: nil,
			},
			expectedError: false,
		},
		{
			name:         "Failed API call; Failed update",
			projectState: projectStateName,
			projectPlan:  projectStateNameDiff,
			mockResponses: ProjectResponse{
				Project:      nil,
				HTTPResponse: &http.Response{StatusCode: 503},
				Err:          errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocksvc.NewProjectsApi(t)
			svc.EXPECT().UpdateProject(mock.Anything, mock.Anything, mock.Anything).Return(admin.UpdateProjectApiRequest{ApiService: svc}).Maybe()

			svc.EXPECT().UpdateProjectExecute(mock.Anything).Return(tc.mockResponses.Project, tc.mockResponses.HTTPResponse, tc.mockResponses.Err).Maybe()

			err := project.UpdateProject(context.Background(), svc, &testCases[i].projectState, &testCases[i].projectPlan)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestUpdateProjectLimits(t *testing.T) {
	twoLimits := []project.TFLimitModel{
		{
			Name: types.StringValue("limit1"),
		},
		{
			Name: types.StringValue("limit2"),
		},
	}
	oneLimit := []project.TFLimitModel{
		{
			Name: types.StringValue("limit1"),
		},
	}
	updatedLimit := []project.TFLimitModel{
		{
			Name:  types.StringValue("limit1"),
			Value: types.Int64Value(6),
		},
	}
	singleLimitSet, _ := types.SetValueFrom(context.Background(), project.TfLimitObjectType, oneLimit)
	updatedLimitSet, _ := types.SetValueFrom(context.Background(), project.TfLimitObjectType, updatedLimit)
	twoLimitSet, _ := types.SetValueFrom(context.Background(), project.TfLimitObjectType, twoLimits)
	testCases := []struct {
		name          string
		mockResponses DeleteProjectLimitResponse
		projectState  project.TFProjectRSModel
		projectPlan   project.TFProjectRSModel
		expectedError bool
	}{
		{
			name: "Limits has not changed",
			projectState: project.TFProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			mockResponses: DeleteProjectLimitResponse{},
			expectedError: false,
		},
		{
			name: "Adding limits",
			projectState: project.TFProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:   name,
				Limits: twoLimitSet,
			},
			mockResponses: DeleteProjectLimitResponse{
				Err: nil,
			},
			expectedError: false,
		},
		{
			name: "Removing limits",
			projectState: project.TFProjectRSModel{
				Name:   name,
				Limits: twoLimitSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			mockResponses: DeleteProjectLimitResponse{
				Err: nil,
			},
			expectedError: false,
		},
		{
			name: "Updating limits",
			projectState: project.TFProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:   name,
				Limits: updatedLimitSet,
			},
			mockResponses: DeleteProjectLimitResponse{
				Err: nil,
			},
			expectedError: false,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocksvc.NewProjectsApi(t)

			svc.EXPECT().DeleteProjectLimit(mock.Anything, mock.Anything, mock.Anything).Return(admin.DeleteProjectLimitApiRequest{ApiService: svc}).Maybe()
			svc.EXPECT().DeleteProjectLimitExecute(mock.Anything).Return(tc.mockResponses.DeleteProjectLimit, tc.mockResponses.HTTPResponse, tc.mockResponses.Err).Maybe()

			svc.EXPECT().SetProjectLimit(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(admin.SetProjectLimitApiRequest{ApiService: svc}).Maybe()
			svc.EXPECT().SetProjectLimitExecute(mock.Anything).Return(nil, nil, nil).Maybe()

			err := project.UpdateProjectLimits(context.Background(), svc, &testCases[i].projectState, &testCases[i].projectPlan)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestUpdateProjectTeams(t *testing.T) {
	teamRoles, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"BASIC_PERMISSION"})
	teamOne := project.TFTeamModel{
		TeamID:    types.StringValue("team1"),
		RoleNames: teamRoles,
	}
	teamTwo := project.TFTeamModel{
		TeamID: types.StringValue("team2"),
	}
	teamRolesUpdated, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"ADMIN_PERMISSION"})
	updatedTeam := project.TFTeamModel{
		TeamID:    types.StringValue("team1"),
		RoleNames: teamRolesUpdated,
	}
	singleTeamSet, _ := types.SetValueFrom(context.Background(), project.TfTeamObjectType, []project.TFTeamModel{teamOne})
	twoTeamSet, _ := types.SetValueFrom(context.Background(), project.TfTeamObjectType, []project.TFTeamModel{teamOne, teamTwo})
	updatedTeamSet, _ := types.SetValueFrom(context.Background(), project.TfTeamObjectType, []project.TFTeamModel{updatedTeam})

	testCases := []struct {
		name          string
		projectState  project.TFProjectRSModel
		projectPlan   project.TFProjectRSModel
		expectedError bool
	}{
		{
			name: "Teams has not changed",
			projectState: project.TFProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			expectedError: false,
		},
		{
			name: "Add teams",
			projectState: project.TFProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:  name,
				Teams: twoTeamSet,
			},
			expectedError: false,
		},
		{
			name: "Remove teams",
			projectState: project.TFProjectRSModel{
				Name:  name,
				Teams: twoTeamSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			expectedError: false,
		},
		{
			name: "Update teams",
			projectState: project.TFProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			projectPlan: project.TFProjectRSModel{
				Name:  name,
				Teams: updatedTeamSet,
			},
			expectedError: false,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocksvc.NewTeamsApi(t)

			svc.EXPECT().AddAllTeamsToProject(mock.Anything, mock.Anything, mock.Anything).Return(admin.AddAllTeamsToProjectApiRequest{ApiService: svc}).Maybe()
			svc.EXPECT().AddAllTeamsToProjectExecute(mock.Anything).Return(nil, nil, nil).Maybe()

			svc.EXPECT().RemoveProjectTeam(mock.Anything, mock.Anything, mock.Anything).Return(admin.RemoveProjectTeamApiRequest{ApiService: svc}).Maybe()
			svc.EXPECT().RemoveProjectTeamExecute(mock.Anything).Return(nil, nil).Maybe()

			svc.EXPECT().UpdateTeamRoles(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(admin.UpdateTeamRolesApiRequest{ApiService: svc}).Maybe()
			svc.EXPECT().UpdateTeamRolesExecute(mock.Anything).Return(nil, nil, nil).Maybe()

			err := project.UpdateProjectTeams(context.Background(), svc, &testCases[i].projectState, &testCases[i].projectPlan)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestResourceProjectDependentsDeletingRefreshFunc(t *testing.T) {
	testCases := []struct {
		mockResponses AdvancedClusterDescriptionResponse
		name          string
		expectedError bool
	}{
		{
			name: "Error not from the API",
			mockResponses: AdvancedClusterDescriptionResponse{
				AdvancedClusterDescription: &admin.PaginatedAdvancedClusterDescription{},
				Err:                        errors.New("Non-API error"),
			},
			expectedError: true,
		},
		{
			name: "Error from the API",
			mockResponses: AdvancedClusterDescriptionResponse{
				AdvancedClusterDescription: &admin.PaginatedAdvancedClusterDescription{},
				Err:                        &admin.GenericOpenAPIError{},
			},
			expectedError: true,
		},
		{
			name: "Successful API call",
			mockResponses: AdvancedClusterDescriptionResponse{
				AdvancedClusterDescription: &admin.PaginatedAdvancedClusterDescription{
					TotalCount: conversion.IntPtr(2),
					Results: &[]admin.AdvancedClusterDescription{
						{StateName: conversion.StringPtr("IDLE")},
						{StateName: conversion.StringPtr("DELETING")},
					},
				},
				Err: nil,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocksvc.NewClustersApi(t)

			svc.EXPECT().ListClusters(mock.Anything, dummyProjectID).Return(admin.ListClustersApiRequest{ApiService: svc})
			svc.EXPECT().ListClustersExecute(mock.Anything).Return(tc.mockResponses.AdvancedClusterDescription, tc.mockResponses.HTTPResponse, tc.mockResponses.Err)

			_, _, err := project.ResourceProjectDependentsDeletingRefreshFunc(context.Background(), dummyProjectID, svc)()

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

const (
	resourceName         = "mongodbatlas_project.test"
	dataSourceNameByID   = "data.mongodbatlas_project.test"
	dataSourceNameByName = "data.mongodbatlas_project.test2"
	dataSourcePluralName = "data.mongodbatlas_projects.test"
)

func TestAccProjectRSProject_basic(t *testing.T) {
	var (
		group        admin.Group
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount = "0"
		projectName  = acc.RandomProjectName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 3) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(0)),
							RoleNames: &[]string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(1)),
							RoleNames: &[]string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_addresses.services.clusters.#"),

					resource.TestCheckResourceAttr(dataSourceNameByID, "name", projectName),
					resource.TestCheckResourceAttr(dataSourceNameByID, "org_id", orgID),
					resource.TestCheckResourceAttr(dataSourceNameByID, "cluster_count", clusterCount),
					resource.TestCheckResourceAttr(dataSourceNameByID, "teams.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "ip_addresses.services.clusters.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "is_collect_database_specifics_statistics_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "is_data_explorer_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "is_extended_storage_sizes_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "is_performance_advisor_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "is_realtime_performance_panel_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "is_schema_advisor_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByID, "limits.0.name"),

					resource.TestCheckResourceAttr(dataSourceNameByName, "name", projectName),
					resource.TestCheckResourceAttr(dataSourceNameByName, "org_id", orgID),
					resource.TestCheckResourceAttr(dataSourceNameByName, "cluster_count", clusterCount),
					resource.TestCheckResourceAttr(dataSourceNameByName, "teams.#", "2"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "ip_addresses.services.clusters.#"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "is_collect_database_specifics_statistics_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "is_data_explorer_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "is_extended_storage_sizes_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "is_performance_advisor_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "is_realtime_performance_panel_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "is_schema_advisor_enabled"),
					resource.TestCheckResourceAttrSet(dataSourceNameByName, "limits.0.name"),

					resource.TestCheckResourceAttrSet(dataSourcePluralName, "total_count"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
				),
			},
			{
				Config: configBasic(projectName, orgID,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(0)),
							RoleNames: &[]string{"GROUP_OWNER"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(1)),
							RoleNames: &[]string{"GROUP_DATA_ACCESS_READ_WRITE"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(2)),
							RoleNames: &[]string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "3"),
				),
			},
			{
				Config: configBasic(projectName, orgID,

					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(0)),
							RoleNames: &[]string{"GROUP_READ_ONLY", "GROUP_READ_ONLY"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(1)),
							RoleNames: &[]string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "2"),
				),
			},
		},
	})
}

func TestAccProjectRSProject_withProjectOwner(t *testing.T) {
	var (
		group          admin.Group
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithOwner(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccProjectRSGovProject_withProjectOwner(t *testing.T) {
	acc.SkipTestForCI(t) // Gov test config not set
	var (
		group          admin.Group
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID_GOV")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID_GOV")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckGov(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectGovWithOwner(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}
func TestAccProjectRSProject_withFalseDefaultSettings(t *testing.T) {
	var (
		group          admin.Group
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithFalseDefaultSettings(projectName, orgID, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccProjectRSProject_withUpdatedSettings(t *testing.T) {
	var (
		group          admin.Group
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "project_owner_id", projectOwnerID),
					resource.TestCheckResourceAttr(resourceName, "with_default_alerts_settings", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_collect_database_specifics_statistics_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_data_explorer_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_extended_storage_sizes_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_performance_advisor_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_realtime_performance_panel_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_schema_advisor_enabled", "false"),
				),
			},
			{
				Config: acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, true),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "with_default_alerts_settings", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_collect_database_specifics_statistics_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_data_explorer_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_extended_storage_sizes_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_performance_advisor_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_realtime_performance_panel_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_schema_advisor_enabled", "true"),
				),
			},
			{
				Config: acc.ConfigProjectWithSettings(projectName, orgID, projectOwnerID, false),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "with_default_alerts_settings", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_collect_database_specifics_statistics_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_data_explorer_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_extended_storage_sizes_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_performance_advisor_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_realtime_performance_panel_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_schema_advisor_enabled", "false"),
				),
			},
		},
	})
}

func TestAccProjectRSProject_withUpdatedRole(t *testing.T) {
	var (
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName     = acc.RandomProjectName()
		roleName        = "GROUP_DATA_ACCESS_ADMIN"
		roleNameUpdated = "GROUP_READ_ONLY"
		clusterCount    = "0"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 1) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithUpdatedRole(projectName, orgID, acc.GetProjectTeamsIDsWithPos(0), roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: acc.ConfigProjectWithUpdatedRole(projectName, orgID, acc.GetProjectTeamsIDsWithPos(0), roleNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
		},
	})
}

func TestAccProjectRSProject_updatedToEmptyRoles(t *testing.T) {
	var (
		group       admin.Group
		projectName = acc.RandomProjectName()
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 1) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(0)),
							RoleNames: &[]string{"GROUP_OWNER", "GROUP_READ_ONLY"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "teams.0.team_id", acc.GetProjectTeamsIDsWithPos(0)),
					resource.TestCheckResourceAttr(resourceName, "teams.0.role_names.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "teams.0.role_names.*", "GROUP_OWNER"),
					resource.TestCheckTypeSetElemAttr(resourceName, "teams.0.role_names.*", "GROUP_READ_ONLY"),
				),
			},
			{
				Config: configBasic(projectName, orgID, nil),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "0"),
				),
			},
		},
	})
}

func TestAccProjectRSProject_importBasic(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectName, orgID,
					[]*admin.TeamRole{},
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"with_default_alerts_settings"},
			},
		},
	})
}

func TestAccProjectRSProject_withUpdatedLimits(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{
					{
						Name:  "atlas.project.deployment.clusters",
						Value: 1,
					},
					{
						Name:  "atlas.project.deployment.nodesPerPrivateLinkRegion",
						Value: 1,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "limits.0.name", "atlas.project.deployment.clusters"),
					resource.TestCheckResourceAttr(resourceName, "limits.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "limits.1.name", "atlas.project.deployment.nodesPerPrivateLinkRegion"),
					resource.TestCheckResourceAttr(resourceName, "limits.1.value", "1"),
				),
			},
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{
					{
						Name:  "atlas.project.deployment.nodesPerPrivateLinkRegion",
						Value: 2,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "limits.0.name", "atlas.project.deployment.nodesPerPrivateLinkRegion"),
					resource.TestCheckResourceAttr(resourceName, "limits.0.value", "2"),
				),
			},
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{
					{
						Name:  "atlas.project.deployment.nodesPerPrivateLinkRegion",
						Value: 3,
					},
					{
						Name:  "atlas.project.security.databaseAccess.customRoles",
						Value: 110,
					},
					{
						Name:  "atlas.project.security.databaseAccess.users",
						Value: 30,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"limits.*",
						map[string]string{
							"name":  "atlas.project.deployment.nodesPerPrivateLinkRegion",
							"value": "3",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"limits.*",
						map[string]string{
							"name":  "atlas.project.security.databaseAccess.customRoles",
							"value": "110",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"limits.*",
						map[string]string{
							"name":  "atlas.project.security.databaseAccess.users",
							"value": "30",
						},
					),
				),
			},
		},
	})
}

func TestAccProjectRSProject_updatedToEmptyLimits(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{
					{
						Name:  "atlas.project.deployment.clusters",
						Value: 1,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "limits.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "limits.0.name", "atlas.project.deployment.clusters"),
					resource.TestCheckResourceAttr(resourceName, "limits.0.value", "1"),
				),
			},
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "limits.#", "0"),
				),
			},
		},
	})
}

func TestAccProjectRSProject_withInvalidLimitName(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{
					{
						Name:  "incorrect.name",
						Value: 1,
					},
				}),
				ExpectError: regexp.MustCompile("Not Found"),
			},
		},
	})
}

func TestAccProjectRSProject_withInvalidLimitNameOnUpdate(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			{
				Config: acc.ConfigProjectWithLimits(projectName, orgID, []*admin.DataFederationLimit{
					{
						Name:  "incorrect.name",
						Value: 1,
					},
				}),
				ExpectError: regexp.MustCompile("Not Found"),
			},
		},
	})
}

func createDataFederationLimit(limitName string) admin.DataFederationLimit {
	return admin.DataFederationLimit{
		Name: limitName,
	}
}

func configBasic(orgID, projectName string, includePluralDataSource bool, teams []*admin.TeamRole) string {
	var ts string
	var plural string

	if includePluralDataSource {
		plural = `
			data "mongodbatlas_projects" "test" {
			}
		`
	}

	for _, t := range teams {
		ts += fmt.Sprintf(`
		teams {
			team_id = "%s"
			role_names = %s
		}
		`, t.GetTeamId(), strings.ReplaceAll(fmt.Sprintf("%+q", *t.RoleNames), " ", ","))
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name  			 = %[2]q

			%[3]s
		}

		data "mongodbatlas_project" "test" {
			project_id = mongodbatlas_project.test.id
		}

		data "mongodbatlas_project" "test2" {
			name = mongodbatlas_project.test.name
		}

		%[4]s
	`, orgID, projectName, ts, plural)
}

type TeamRoleResponse struct {
	TeamRole     *admin.PaginatedTeamRole
	HTTPResponse *http.Response
	Err          error
}
type LimitsResponse struct {
	Err          error
	HTTPResponse *http.Response
	Limits       []admin.DataFederationLimit
}
type GroupSettingsResponse struct {
	GroupSettings *admin.GroupSettings
	HTTPResponse  *http.Response
	Err           error
}

type IPAddressesResponse struct {
	IPAddresses  *admin.GroupIPAddresses
	HTTPResponse *http.Response
	Err          error
}
type ProjectResponse struct {
	Project      *admin.Group
	HTTPResponse *http.Response
	Err          error
}
type DeleteProjectLimitResponse struct {
	DeleteProjectLimit map[string]interface{}
	HTTPResponse       *http.Response
	Err                error
}
type AdvancedClusterDescriptionResponse struct {
	AdvancedClusterDescription *admin.PaginatedAdvancedClusterDescription
	HTTPResponse               *http.Response
	Err                        error
}
