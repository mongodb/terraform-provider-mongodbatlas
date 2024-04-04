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

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/mock"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
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
			teamsMock := mockadmin.NewTeamsApi(t)
			projectsMock := mockadmin.NewProjectsApi(t)

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
			svc := mockadmin.NewProjectsApi(t)
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
			svc := mockadmin.NewProjectsApi(t)

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
			svc := mockadmin.NewTeamsApi(t)

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
			svc := mockadmin.NewClustersApi(t)

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

func TestAccProject_basic(t *testing.T) {
	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)
	commonChecks := map[string]string{
		"name":          projectName,
		"org_id":        orgID,
		"cluster_count": "0",
		"teams.#":       "2",
	}
	commonSetChecks := []string{
		"ip_addresses.services.clusters.#",
		"is_collect_database_specifics_statistics_enabled",
		"is_data_explorer_enabled",
		"is_extended_storage_sizes_enabled",
		"is_performance_advisor_enabled",
		"is_realtime_performance_panel_enabled",
		"is_schema_advisor_enabled",
	}
	checks := acc.AddAttrChecks(resourceName, nil, commonChecks)
	checks = acc.AddAttrChecks(dataSourceNameByID, checks, commonChecks)
	checks = acc.AddAttrChecks(dataSourceNameByName, checks, commonChecks)
	checks = acc.AddAttrSetChecks(resourceName, checks, commonSetChecks...)
	checks = acc.AddAttrSetChecks(dataSourceNameByID, checks, commonSetChecks...)
	checks = acc.AddAttrSetChecks(dataSourceNameByName, checks, commonSetChecks...)
	checks = append(checks, checkExists(resourceName), checkExists(dataSourceNameByID), checkExists(dataSourceNameByName))
	checks = acc.AddAttrSetChecks(dataSourcePluralName, checks, "total_count", "results.#")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 3) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, projectOwnerID, true,
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
				Check: resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: configBasic(orgID, projectName, projectOwnerID, false,
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
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "3"),
				),
			},
			{
				Config: configBasic(orgID, projectName, projectOwnerID, false,

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
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"with_default_alerts_settings", "project_owner_id"},
			},
		},
	})
}

func TestAccProjectGov_withProjectOwner(t *testing.T) {
	acc.SkipTestForCI(t) // Gov test config not set

	var (
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID_GOV")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID_GOV")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckGovBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configGovWithOwner(orgID, projectName, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "region_usage_restrictions", "GOV_REGIONS_ONLY"),
				),
			},
		},
	})
}
func TestAccProject_withFalseDefaultSettings(t *testing.T) {
	var (
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
				Config: configWithFalseDefaultSettings(orgID, projectName, projectOwnerID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
		},
	})
}

func TestAccProject_withUpdatedSettings(t *testing.T) {
	var (
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
					checkExists(resourceName),
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
					checkExists(resourceName),
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
					checkExists(resourceName),
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

func TestAccProject_withUpdatedRole(t *testing.T) {
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
				Config: configWithUpdatedRole(orgID, projectName, acc.GetProjectTeamsIDsWithPos(0), roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: configWithUpdatedRole(orgID, projectName, acc.GetProjectTeamsIDsWithPos(0), roleNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
		},
	})
}

func TestAccProject_updatedToEmptyRoles(t *testing.T) {
	var (
		projectName = acc.RandomProjectName()
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIDsWithMinCount(t, 1) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configBasic(orgID, projectName, "", false,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIDsWithPos(0)),
							RoleNames: &[]string{"GROUP_OWNER", "GROUP_READ_ONLY"},
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "teams.0.team_id", acc.GetProjectTeamsIDsWithPos(0)),
					resource.TestCheckResourceAttr(resourceName, "teams.0.role_names.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "teams.0.role_names.*", "GROUP_OWNER"),
					resource.TestCheckTypeSetElemAttr(resourceName, "teams.0.role_names.*", "GROUP_READ_ONLY"),
				),
			},
			{
				Config: configBasic(orgID, projectName, "", false, nil),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "teams.#", "0"),
				),
			},
		},
	})
}

func TestAccProject_withUpdatedLimits(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		limitChecks = []map[string]string{
			{
				"name":  "atlas.project.deployment.clusters",
				"value": "1",
			},
			{
				"name":  "atlas.project.deployment.nodesPerPrivateLinkRegion",
				"value": "1",
			},
		}
	)
	checks := []resource.TestCheckFunc{checkExists(resourceName), checkExists(dataSourceNameByID)}
	for _, check := range limitChecks {
		checks = append(checks,
			resource.TestCheckTypeSetElemNestedAttrs(resourceName, "limits.*", check),
			resource.TestCheckTypeSetElemNestedAttrs(dataSourceNameByID, "limits.*", check),
		)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: configWithLimits(orgID, projectName, []*admin.DataFederationLimit{
					{
						Name:  "atlas.project.deployment.clusters",
						Value: 1,
					},
					{
						Name:  "atlas.project.deployment.nodesPerPrivateLinkRegion",
						Value: 1,
					},
				}),
				Check: resource.ComposeTestCheckFunc(checks...),
			},
			{
				Config: configWithLimits(orgID, projectName, []*admin.DataFederationLimit{
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
				Config: configWithLimits(orgID, projectName, []*admin.DataFederationLimit{
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

func TestAccProject_updatedToEmptyLimits(t *testing.T) {
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
				Config: configWithLimits(orgID, projectName, []*admin.DataFederationLimit{
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
				Config: configWithLimits(orgID, projectName, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "limits.#", "0"),
				),
			},
		},
	})
}

func TestAccProject_withInvalidLimitName(t *testing.T) {
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
				Config: configWithLimits(orgID, projectName, []*admin.DataFederationLimit{
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

func TestAccProject_withInvalidLimitNameOnUpdate(t *testing.T) {
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
				Config: configWithLimits(orgID, projectName, []*admin.DataFederationLimit{}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
				),
			},
			{
				Config: configWithLimits(orgID, projectName, []*admin.DataFederationLimit{
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

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		if _, _, err := acc.ConnV2().ProjectsApi.GetProjectByName(context.Background(), rs.Primary.Attributes["name"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("project (%s) does not exist", rs.Primary.ID)
	}
}

func configBasic(orgID, projectName, projectOwnerID string, includeDataSource bool, teams []*admin.TeamRole) string {
	var dataSourceStr string
	if includeDataSource {
		dataSourceStr = `
			data "mongodbatlas_project" "test" {
				project_id = mongodbatlas_project.test.id
			}

			data "mongodbatlas_project" "test2" {
				name = mongodbatlas_project.test.name
			}

			data "mongodbatlas_projects" "test" {
			}
		`
	}

	var additionalStr string
	if projectOwnerID != "" {
		additionalStr = fmt.Sprintf("project_owner_id = %q\n", projectOwnerID)
	}

	for _, t := range teams {
		additionalStr += fmt.Sprintf(`
		teams {
			team_id = %q
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

		%[4]s
	`, orgID, projectName, additionalStr, dataSourceStr)
}

func configGovWithOwner(orgID, projectName, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name   			 = %[2]q
			project_owner_id = %[3]q
			region_usage_restrictions = "GOV_REGIONS_ONLY"
		}
	`, orgID, projectName, projectOwnerID)
}

func configWithFalseDefaultSettings(orgID, projectName, projectOwnerID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name   			 = %[2]q
			project_owner_id = %[3]q
			with_default_alerts_settings = false
		}
	`, orgID, projectName, projectOwnerID)
}

func configWithLimits(orgID, projectName string, limits []*admin.DataFederationLimit) string {
	var limitsString string

	for _, limit := range limits {
		limitsString += fmt.Sprintf(`
		limits {
			name = %[1]q
			value = %[2]d
		}
		`, limit.Name, limit.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id 			 = %[1]q
			name   			 = %[2]q

			%[3]s
		}

		data "mongodbatlas_project" "test" {
			project_id = mongodbatlas_project.test.id
		}
	`, orgID, projectName, limitsString)
}

func configWithUpdatedRole(orgID, projectName, teamID, roleName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q

			teams {
				team_id = %[3]q
				role_names = [ %[4]q ]
			}
		}
	`, orgID, projectName, teamID, roleName)
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
