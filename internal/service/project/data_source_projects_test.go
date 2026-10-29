package project_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312027/admin"
	"go.mongodb.org/atlas-sdk/v20250312027/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
)

func TestPopulateProjectsDataSourceModel(t *testing.T) {
	notFoundErr := &admin.GenericOpenAPIError{}
	notFoundErr.SetModel(admin.ApiError{Error: http.StatusNotFound, ErrorCode: "GROUP_NOT_FOUND"})

	testCases := map[string]struct {
		firstTeamsErr   error
		expectedIDs     []string
		expectedError   bool
		secondProjectOK bool
	}{
		"all projects are returned": {
			expectedIDs:     []string{"project-1", "project-2"},
			secondProjectOK: true,
		},
		"project deleted after the list call is skipped": {
			firstTeamsErr:   notFoundErr,
			expectedIDs:     []string{"project-2"},
			secondProjectOK: true,
		},
		"transient error is reported instead of dropping the project": {
			firstTeamsErr: errors.New("EOF"),
			expectedError: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			teamsMock := mockadmin.NewTeamsAPI(t)
			projectsMock := mockadmin.NewProjectsAPI(t)
			perfMock := mockadmin.NewPerformanceAdvisorAPI(t)
			usersMock := mockadmin.NewMongoDBCloudUsersAPI(t)

			teamsMock.EXPECT().ListGroupTeams(mock.Anything, "project-1").Return(admin.ListGroupTeamsApiRequest{ApiService: teamsMock}).Once()
			if tc.firstTeamsErr != nil {
				teamsMock.EXPECT().ListGroupTeamsExecute(mock.Anything).Return(nil, nil, tc.firstTeamsErr).Once()
			}
			if tc.secondProjectOK {
				teamsMock.EXPECT().ListGroupTeams(mock.Anything, "project-2").Return(admin.ListGroupTeamsApiRequest{ApiService: teamsMock}).Once()
				teamsMock.EXPECT().ListGroupTeamsExecute(mock.Anything).Return(&admin.PaginatedTeamRole{}, nil, nil)
				projectsMock.EXPECT().ListGroupLimits(mock.Anything, mock.Anything).Return(admin.ListGroupLimitsApiRequest{ApiService: projectsMock})
				projectsMock.EXPECT().ListGroupLimitsExecute(mock.Anything).Return([]admin.DataFederationLimit{}, nil, nil)
				projectsMock.EXPECT().GetGroupSettings(mock.Anything, mock.Anything).Return(admin.GetGroupSettingsApiRequest{ApiService: projectsMock})
				projectsMock.EXPECT().GetGroupSettingsExecute(mock.Anything).Return(&projectSettingsSDK, nil, nil)
				projectsMock.EXPECT().GetGroupIpAddresses(mock.Anything, mock.Anything).Return(admin.GetGroupIpAddressesApiRequest{ApiService: projectsMock})
				projectsMock.EXPECT().GetGroupIpAddressesExecute(mock.Anything).Return(&IPAddressesSDK, nil, nil)
				perfMock.EXPECT().GetManagedSlowMs(mock.Anything, mock.Anything).Return(admin.GetManagedSlowMsApiRequest{ApiService: perfMock})
				perfMock.EXPECT().GetManagedSlowMsExecute(mock.Anything).Return(true, nil, nil)
				usersMock.EXPECT().ListGroupUsers(mock.Anything, mock.Anything).Return(admin.ListGroupUsersApiRequest{ApiService: usersMock})
				usersMock.EXPECT().ListGroupUsersExecute(mock.Anything).Return(&admin.PaginatedGroupUser{}, nil, nil)
			}

			connV2 := &admin.APIClient{
				ProjectsAPI:           projectsMock,
				TeamsAPI:              teamsMock,
				PerformanceAdvisorAPI: perfMock,
				MongoDBCloudUsersAPI:  usersMock,
			}
			projectsRes := &admin.PaginatedAtlasGroup{
				Results: []admin.Group{
					{Id: new("project-1"), Name: "project-1", OrgId: "org"},
					{Id: new("project-2"), Name: "project-2", OrgId: "org"},
				},
				TotalCount: new(2),
			}

			results, diags := project.PopulateProjectsDataSourceResultsForTest(t.Context(), connV2, projectsRes)

			assert.Equal(t, tc.expectedError, diags.HasError(), "diagnostics: %v", diags)
			if tc.expectedError {
				return
			}
			ids := make([]string, len(results))
			for i, result := range results {
				ids[i] = result.ID.ValueString()
			}
			assert.Equal(t, tc.expectedIDs, ids)
		})
	}
}
