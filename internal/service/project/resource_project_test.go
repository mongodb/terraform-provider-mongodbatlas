package project_test

import (
	"context"
	"errors"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mocksvc"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115002/admin"
)

var (
	successfulPaginatedTeamRole = ProjectResponse{
		ProjectTeamResp: &admin.PaginatedTeamRole{},
		Err:             nil,
	}
	successfulDataFederationLimit = ProjectResponse{
		LimitsResponse: []admin.DataFederationLimit{},
		Err:            nil,
	}
	successfulGroupSettingsResponse = ProjectResponse{
		GroupSettingsResponse: &admin.GroupSettings{},
		Err:                   nil,
	}
	name             = types.StringValue("sameName")
	diffName         = types.StringValue("diffName")
	projectStateName = project.TfProjectRSModel{
		Name: name,
	}
	projectStateNameDiff = project.TfProjectRSModel{
		Name: diffName,
	}
	dummyProjectID = "6575af27f93c7a6a4b50b239"
)

func TestGetProjectPropsFromAPI(t *testing.T) {
	testCases := []struct {
		name            string
		teamRoleReponse ProjectResponse
		limitResponse   ProjectResponse
		groupResponse   ProjectResponse
		expectedError   bool
	}{
		{
			name:            "Successful",
			teamRoleReponse: successfulPaginatedTeamRole,
			limitResponse:   successfulDataFederationLimit,
			groupResponse:   successfulGroupSettingsResponse,
			expectedError:   false,
		},
		{
			name: "Fail to get project's teams assigned ",
			teamRoleReponse: ProjectResponse{
				ProjectTeamResp: nil,
				HTTPResponse:    &http.Response{StatusCode: 503},
				Err:             errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
		{
			name:            "Fail to get project's limits",
			teamRoleReponse: successfulPaginatedTeamRole,
			limitResponse: ProjectResponse{
				LimitsResponse: nil,
				HTTPResponse:   &http.Response{StatusCode: 503},
				Err:            errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
		{
			name:            "Fail to get project's settings",
			teamRoleReponse: successfulPaginatedTeamRole,
			limitResponse:   successfulDataFederationLimit,
			groupResponse: ProjectResponse{
				GroupSettingsResponse: nil,
				HTTPResponse:          &http.Response{StatusCode: 503},
				Err:                   errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mocksvc.NewGroupProjectService(t)
			testObject.On("ListProjectTeams", mock.Anything, mock.Anything).Return(tc.teamRoleReponse.ProjectTeamResp, tc.teamRoleReponse.HTTPResponse, tc.teamRoleReponse.Err)
			testObject.On("ListProjectLimits", mock.Anything, mock.Anything).Return(tc.limitResponse.LimitsResponse, tc.limitResponse.HTTPResponse, tc.limitResponse.Err).Maybe()
			testObject.On("GetProjectSettings", mock.Anything, mock.Anything).Return(tc.groupResponse.GroupSettingsResponse, tc.groupResponse.HTTPResponse, tc.groupResponse.Err).Maybe()

			_, _, _, err := project.GetProjectPropsFromAPI(context.Background(), testObject, dummyProjectID)

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
		tfLimits       []project.TfLimitModel
		expectedResult []admin.DataFederationLimit
	}{
		{
			name: "FilterUserDefinedLimits",
			allAtlasLimits: []admin.DataFederationLimit{
				createDataFederationLimit("1"),
				createDataFederationLimit("2"),
				createDataFederationLimit("3"),
			},
			tfLimits: []project.TfLimitModel{
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
			tfLimits:       []project.TfLimitModel{},
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
		projectState  project.TfProjectRSModel
		projectPlan   project.TfProjectRSModel
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
				ProjectResp:  nil,
				HTTPResponse: &http.Response{StatusCode: 503},
				Err:          errors.New("Service Unavailable"),
			},
			expectedError: true,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mocksvc.NewGroupProjectService(t)

			testObject.On("UpdateProject", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockResponses.ProjectResp, tc.mockResponses.HTTPResponse, tc.mockResponses.Err).Maybe()

			err := project.UpdateProject(context.Background(), testObject, &testCases[i].projectState, &testCases[i].projectPlan)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestUpdateProjectLimits(t *testing.T) {
	twoLimits := []project.TfLimitModel{
		{
			Name: types.StringValue("limit1"),
		},
		{
			Name: types.StringValue("limit2"),
		},
	}
	oneLimit := []project.TfLimitModel{
		{
			Name: types.StringValue("limit1"),
		},
	}
	updatedLimit := []project.TfLimitModel{
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
		mockResponses ProjectResponse
		projectState  project.TfProjectRSModel
		projectPlan   project.TfProjectRSModel
		expectedError bool
	}{
		{
			name: "Limits has not changed",
			projectState: project.TfProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			mockResponses: ProjectResponse{},
			expectedError: false,
		},
		{
			name: "Adding limits",
			projectState: project.TfProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:   name,
				Limits: twoLimitSet,
			},
			mockResponses: ProjectResponse{
				Err: nil,
			},
			expectedError: false,
		},
		{
			name: "Removing limits",
			projectState: project.TfProjectRSModel{
				Name:   name,
				Limits: twoLimitSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			mockResponses: ProjectResponse{
				Err: nil,
			},
			expectedError: false,
		},
		{
			name: "Updating limits",
			projectState: project.TfProjectRSModel{
				Name:   name,
				Limits: singleLimitSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:   name,
				Limits: updatedLimitSet,
			},
			mockResponses: ProjectResponse{
				Err: nil,
			},
			expectedError: false,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mocksvc.NewGroupProjectService(t)

			testObject.On("DeleteProjectLimit", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockResponses.DeleteProjectLimitResponse, tc.mockResponses.HTTPResponse, tc.mockResponses.Err).Maybe()
			testObject.On("SetProjectLimit", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil).Maybe()

			err := project.UpdateProjectLimits(context.Background(), testObject, &testCases[i].projectState, &testCases[i].projectPlan)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestUpdateProjectTeams(t *testing.T) {
	teamRoles, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"BASIC_PERMISSION"})
	teamOne := project.TfTeamModel{
		TeamID:    types.StringValue("team1"),
		RoleNames: teamRoles,
	}
	teamTwo := project.TfTeamModel{
		TeamID: types.StringValue("team2"),
	}
	teamRolesUpdated, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"ADMIN_PERMISSION"})
	updatedTeam := project.TfTeamModel{
		TeamID:    types.StringValue("team1"),
		RoleNames: teamRolesUpdated,
	}
	singleTeamSet, _ := types.SetValueFrom(context.Background(), project.TfTeamObjectType, []project.TfTeamModel{teamOne})
	twoTeamSet, _ := types.SetValueFrom(context.Background(), project.TfTeamObjectType, []project.TfTeamModel{teamOne, teamTwo})
	updatedTeamSet, _ := types.SetValueFrom(context.Background(), project.TfTeamObjectType, []project.TfTeamModel{updatedTeam})

	testCases := []struct {
		name          string
		mockResponses []ProjectResponse
		projectState  project.TfProjectRSModel
		projectPlan   project.TfProjectRSModel
		expectedError bool
	}{
		{
			name: "Teams has not changed",
			projectState: project.TfProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			expectedError: false,
		},
		{
			name: "Add teams",
			projectState: project.TfProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:  name,
				Teams: twoTeamSet,
			},
			expectedError: false,
		},
		{
			name: "Remove teams",
			projectState: project.TfProjectRSModel{
				Name:  name,
				Teams: twoTeamSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			expectedError: false,
		},
		{
			name: "Update teams",
			projectState: project.TfProjectRSModel{
				Name:  name,
				Teams: singleTeamSet,
			},
			projectPlan: project.TfProjectRSModel{
				Name:  name,
				Teams: updatedTeamSet,
			},
			expectedError: false,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testObject := mocksvc.NewGroupProjectService(t)

			testObject.On("AddAllTeamsToProject", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil).Maybe()
			testObject.On("RemoveProjectTeam", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil).Maybe()
			testObject.On("UpdateTeamRoles", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, nil).Maybe()

			err := project.UpdateProjectTeams(context.Background(), testObject, &testCases[i].projectState, &testCases[i].projectPlan)

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestResourceProjectDependentsDeletingRefreshFunc(t *testing.T) {
	testCases := []struct {
		name          string
		mockResponses ProjectResponse
		expectedError bool
	}{
		{
			name: "Error not from the API",
			mockResponses: ProjectResponse{
				clusterReponse: &admin.PaginatedAdvancedClusterDescription{},
				Err:            errors.New("Non-API error"),
			},
			expectedError: true,
		},
		{
			name: "Error from the API",
			mockResponses: ProjectResponse{
				clusterReponse: &admin.PaginatedAdvancedClusterDescription{},
				Err:            &admin.GenericOpenAPIError{},
			},
			expectedError: true,
		},
		{
			name: "Successful API call",
			mockResponses: ProjectResponse{
				clusterReponse: &admin.PaginatedAdvancedClusterDescription{
					TotalCount: conversion.IntPtr(2),
					Results: []admin.AdvancedClusterDescription{
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
			testObject := mocksvc.NewGroupProjectService(t)

			testObject.On("ListClusters", mock.Anything, dummyProjectID).Return(tc.mockResponses.clusterReponse, tc.mockResponses.HTTPResponse, tc.mockResponses.Err)

			_, _, err := project.ResourceProjectDependentsDeletingRefreshFunc(context.Background(), dummyProjectID, testObject)()

			if (err != nil) != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, err)
			}
		})
	}
}

func TestAccProjectRSProject_basic(t *testing.T) {
	var (
		group        admin.Group
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("test-acc")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		clusterCount = "0"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIdsWithMinCount(t, 3) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProject(projectName, orgID,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIdsWithPos(0)),
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIdsWithPos(1)),
							RoleNames: []string{"GROUP_DATA_ACCESS_ADMIN", "GROUP_OWNER"},
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
			{
				Config: acc.ConfigProject(projectName, orgID,
					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIdsWithPos(0)),
							RoleNames: []string{"GROUP_OWNER"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIdsWithPos(1)),
							RoleNames: []string{"GROUP_DATA_ACCESS_READ_WRITE"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIdsWithPos(2)),
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_DATA_ACCESS_ADMIN"},
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
				Config: acc.ConfigProject(projectName, orgID,

					[]*admin.TeamRole{
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIdsWithPos(0)),
							RoleNames: []string{"GROUP_READ_ONLY", "GROUP_READ_ONLY"},
						},
						{
							TeamId:    conversion.StringPtr(acc.GetProjectTeamsIdsWithPos(1)),
							RoleNames: []string{"GROUP_OWNER", "GROUP_DATA_ACCESS_ADMIN"},
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
			{
				Config: acc.ConfigProject(projectName, orgID, []*admin.TeamRole{}),
				Check: resource.ComposeTestCheckFunc(
					acc.CheckProjectExists(resourceName, &group),
					acc.CheckProjectAttributes(&group, projectName),
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
					resource.TestCheckNoResourceAttr(resourceName, "teams.#"),
				),
			},
		},
	})
}

func TestAccProjectRSProject_CreateWithProjectOwner(t *testing.T) {
	var (
		group          admin.Group
		resourceName   = "mongodbatlas_project.test"
		projectName    = acctest.RandomWithPrefix("test-acc")
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
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

func TestAccProjectRSGovProject_CreateWithProjectOwner(t *testing.T) {
	var (
		group          admin.Group
		resourceName   = "mongodbatlas_project.test"
		projectName    = acctest.RandomWithPrefix("tf-acc-project")
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID_GOV")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID_GOV")
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
func TestAccProjectRSProject_CreateWithFalseDefaultSettings(t *testing.T) {
	var (
		group          admin.Group
		resourceName   = "mongodbatlas_project.test"
		projectName    = acctest.RandomWithPrefix("tf-acc-project")
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
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

func TestAccProjectRSProject_CreateWithFalseDefaultAdvSettings(t *testing.T) {
	var (
		group          admin.Group
		resourceName   = "mongodbatlas_project.test"
		projectName    = acctest.RandomWithPrefix("tf-acc-project")
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectOwnerID = os.Getenv("MONGODB_ATLAS_PROJECT_OWNER_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasicOwnerID(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithFalseDefaultAdvSettings(projectName, orgID, projectOwnerID),
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

func TestAccProjectRSProject_withUpdatedRole(t *testing.T) {
	var (
		resourceName    = "mongodbatlas_project.test"
		projectName     = acctest.RandomWithPrefix("tf-acc-project")
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		roleName        = "GROUP_DATA_ACCESS_ADMIN"
		roleNameUpdated = "GROUP_READ_ONLY"
		clusterCount    = "0"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t); acc.PreCheckProjectTeamsIdsWithMinCount(t, 1) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProjectWithUpdatedRole(projectName, orgID, acc.GetProjectTeamsIdsWithPos(0), roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
			{
				Config: acc.ConfigProjectWithUpdatedRole(projectName, orgID, acc.GetProjectTeamsIdsWithPos(0), roleNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", projectName),
					resource.TestCheckResourceAttr(resourceName, "org_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "cluster_count", clusterCount),
				),
			},
		},
	})
}

func TestAccProjectRSProject_importBasic(t *testing.T) {
	var (
		projectName  = acctest.RandomWithPrefix("tf-acc-project")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		resourceName = "mongodbatlas_project.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyProject,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigProject(projectName, orgID,
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
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("tf-acc-project")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
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

func TestAccProjectRSProject_withInvalidLimitName(t *testing.T) {
	var (
		projectName = acctest.RandomWithPrefix("tf-acc-project")
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
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
		resourceName = "mongodbatlas_project.test"
		projectName  = acctest.RandomWithPrefix("tf-acc-project")
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
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

type ProjectResponse struct {
	ProjectResp                *admin.Group
	ProjectTeamResp            *admin.PaginatedTeamRole
	GroupSettingsResponse      *admin.GroupSettings
	clusterReponse             *admin.PaginatedAdvancedClusterDescription
	HTTPResponse               *http.Response
	Err                        error
	LimitResponse              admin.DataFederationLimit
	DeleteProjectLimitResponse map[string]interface{}
	LimitsResponse             []admin.DataFederationLimit
}
