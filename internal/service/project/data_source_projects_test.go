package project_test

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312028/admin"
	"go.mongodb.org/atlas-sdk/v20250312028/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/project"
)

func TestPopulateProjectsDataSourceModel(t *testing.T) {
	groupNotFound := &admin.GenericOpenAPIError{}
	groupNotFound.SetModel(admin.ApiError{Error: http.StatusNotFound, ErrorCode: "GROUP_NOT_FOUND"})
	groupNotFound.SetError("404 Not Found")

	serverError := &admin.GenericOpenAPIError{}
	serverError.SetModel(admin.ApiError{Error: http.StatusInternalServerError, ErrorCode: "UNEXPECTED_ERROR"})
	serverError.SetError("500 Internal Server Error")

	testCases := []struct {
		err           error
		httpResponse  *http.Response
		name          string
		expectedError bool
	}{
		{
			name:         "Project deleted after the list is skipped",
			err:          groupNotFound,
			httpResponse: &http.Response{StatusCode: http.StatusNotFound},
		},
		{
			name:          "Server error is returned",
			err:           serverError,
			httpResponse:  &http.Response{StatusCode: http.StatusInternalServerError},
			expectedError: true,
		},
		{
			name:          "Connection error is returned",
			err:           &url.Error{Op: "Get", URL: "https://cloud.mongodb.com/api/atlas/v2/groups/" + dummyProjectID + "/teams", Err: io.EOF},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			teamsMock := mockadmin.NewTeamsAPI(t)
			teamsMock.EXPECT().ListGroupTeams(mock.Anything, dummyProjectID).Return(admin.ListGroupTeamsApiRequest{ApiService: teamsMock})
			teamsMock.EXPECT().ListGroupTeamsExecute(mock.Anything).Return(nil, tc.httpResponse, tc.err)

			connV2 := &admin.APIClient{
				ProjectsAPI:           mockadmin.NewProjectsAPI(t),
				TeamsAPI:              teamsMock,
				PerformanceAdvisorAPI: mockadmin.NewPerformanceAdvisorAPI(t),
				MongoDBCloudUsersAPI:  mockadmin.NewMongoDBCloudUsersAPI(t),
			}
			projectsRes := &admin.PaginatedAtlasGroup{
				Results:    []admin.Group{{Id: new(dummyProjectID)}},
				TotalCount: new(1),
			}

			var model project.TFProjectsDSModel
			diags := project.PopulateProjectsDataSourceModel(t.Context(), connV2, &model, projectsRes)

			if tc.expectedError {
				require.True(t, diags.HasError())
				assert.Contains(t, diags.Errors()[0].Detail(), dummyProjectID)
				return
			}
			assert.False(t, diags.HasError())
			assert.Empty(t, model.Results)
		})
	}
}
