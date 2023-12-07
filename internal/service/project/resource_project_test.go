package project_test

import (
	"context"
	"errors"
	"log"
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
)

func TestGetProjectPropsFromAPI(t *testing.T) {
	testCases := []struct {
		name          string
		projectID     string
		mockResponses []ProjectResponse
		expectedError bool
	}{
		{
			name: "Successful",
			mockResponses: []ProjectResponse{
				successfulPaginatedTeamRole,
				successfulDataFederationLimit,
				successfulGroupSettingsResponse,
			},
			expectedError: false,
			projectID:     "projectId",
		},
		{
			name: "Fail to get project's teams assigned ",
			mockResponses: []ProjectResponse{
				{
					ProjectTeamResp: nil,
					HTTPResponse:    &http.Response{StatusCode: 503},
					Err:             errors.New("Service Unavailable"),
				},
			},
			expectedError: true,
			projectID:     "projectId",
		},
		{
			name: "Fail to get project's limits",
			mockResponses: []ProjectResponse{
				successfulPaginatedTeamRole,
				{
					LimitsResponse: nil,
					HTTPResponse:   &http.Response{StatusCode: 503},
					Err:            errors.New("Service Unavailable"),
				},
			},
			expectedError: true,
			projectID:     "projectId",
		},
		{
			name: "Fail to get project's settings",
			mockResponses: []ProjectResponse{
				successfulPaginatedTeamRole,
				successfulDataFederationLimit,
				{
					GroupSettingsResponse: nil,
					HTTPResponse:          &http.Response{StatusCode: 503},
					Err:                   errors.New("Service Unavailable"),
				},
			},
			expectedError: true,
			projectID:     "projectId",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := MockProjectService{
				MockResponses: tc.mockResponses,
			}
			_, _, _, err := project.GetProjectPropsFromAPI(context.Background(), &mockService, tc.projectID)

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
		mockResponses []ProjectResponse
		projectState  project.TfProjectRSModel
		projectPlan   project.TfProjectRSModel
		expectedError bool
	}{
		{
			name:         "Successful update",
			projectState: projectStateName,
			projectPlan:  projectStateNameDiff,
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
			},
			expectedError: false,
		},
		{
			name:         "Same project names; Failed update",
			projectState: projectStateName,
			projectPlan:  projectStateName,
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
			},
			expectedError: false,
		},
		{
			name:         "Failed API call; Failed update",
			projectState: projectStateName,
			projectPlan:  projectStateNameDiff,
			mockResponses: []ProjectResponse{
				{
					ProjectResp:  nil,
					HTTPResponse: &http.Response{StatusCode: 503},
					Err:          errors.New("Service Unavailable"),
				},
			},
			expectedError: true,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := MockProjectService{
				MockResponses: tc.mockResponses,
			}
			err := project.UpdateProject(context.Background(), &mockService, &testCases[i].projectState, &testCases[i].projectPlan)

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
		mockResponses []ProjectResponse
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
			mockResponses: []ProjectResponse{},
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
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
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
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
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
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
			},
			expectedError: false,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := MockProjectService{
				MockResponses: tc.mockResponses,
			}
			err := project.UpdateProjectLimits(context.Background(), &mockService, &testCases[i].projectState, &testCases[i].projectPlan)

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
			mockResponses: []ProjectResponse{},
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
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
				{
					Err: nil,
				},
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
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
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
			mockResponses: []ProjectResponse{
				{
					Err: nil,
				},
				{
					Err: nil,
				},
			},
			expectedError: false,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := MockProjectService{
				MockResponses: tc.mockResponses,
			}
			err := project.UpdateProjectTeams(context.Background(), &mockService, &testCases[i].projectState, &testCases[i].projectPlan)

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

type MockProjectService struct {
	MockResponses []ProjectResponse
	index         int
}

func (a *MockProjectService) UpdateProject(ctx context.Context, groupID string, groupName *admin.GroupName) (*admin.Group, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.ProjectResp, resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) ListProjectLimits(ctx context.Context, groupID string) ([]admin.DataFederationLimit, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.LimitsResponse, resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) GetProjectSettings(ctx context.Context, groupID string) (*admin.GroupSettings, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.GroupSettingsResponse, resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) ListProjectTeams(ctx context.Context, groupID string) (*admin.PaginatedTeamRole, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.ProjectTeamResp, resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) DeleteProjectLimit(ctx context.Context, limitName, projectID string) (map[string]interface{}, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.DeleteProjectLimitResponse, resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) SetProjectLimit(ctx context.Context, limitName, groupID string, dataFederationLimit *admin.DataFederationLimit) (*admin.DataFederationLimit, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return &resp.LimitResponse, resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) RemoveProjectTeam(ctx context.Context, groupID, teamID string) (*http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) UpdateTeamRoles(ctx context.Context, groupID, teamID string, teamRole *admin.TeamRole) (*admin.PaginatedTeamRole, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.ProjectTeamResp, resp.HTTPResponse, resp.Err
}

func (a *MockProjectService) AddAllTeamsToProject(ctx context.Context, groupID string, teamRole *[]admin.TeamRole) (*admin.PaginatedTeamRole, *http.Response, error) {
	if a.index >= len(a.MockResponses) {
		log.Fatal(errors.New("no more mocked responses available"))
	}
	resp := a.MockResponses[a.index]
	a.index++
	return resp.ProjectTeamResp, resp.HTTPResponse, resp.Err
}

type ProjectResponse struct {
	ProjectResp                *admin.Group
	ProjectTeamResp            *admin.PaginatedTeamRole
	GroupSettingsResponse      *admin.GroupSettings
	HTTPResponse               *http.Response
	Err                        error
	LimitResponse              admin.DataFederationLimit
	DeleteProjectLimitResponse map[string]interface{}
	LimitsResponse             []admin.DataFederationLimit
}
