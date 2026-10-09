package maintenancewindow_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312026/admin"
	"go.mongodb.org/atlas-sdk/v20250312026/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/maintenancewindow"
)

func TestResourceReadProtectedHours(t *testing.T) {
	testCases := map[string]struct {
		apiResponse         *admin.GroupMaintenanceWindow
		stateProtectedHours []any
		expectedProtectedHL []any
	}{
		// Starts with protected hours in state (stale), API returns none: Read must clear
		// the attribute so the loss shows up as drift.
		"no protected hours in API response clears the stale attribute": {
			apiResponse:         &admin.GroupMaintenanceWindow{},
			stateProtectedHours: []any{map[string]any{"start_hour_of_day": 9, "end_hour_of_day": 17}},
			expectedProtectedHL: []any{},
		},
		"protected hours in API response are set": {
			apiResponse: &admin.GroupMaintenanceWindow{
				ProtectedHours: &admin.ProtectedHours{
					StartHourOfDay: new(9),
					EndHourOfDay:   new(17),
				},
			},
			expectedProtectedHL: []any{map[string]any{"start_hour_of_day": 9, "end_hour_of_day": 17}},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			apiMock := mockadmin.NewMaintenanceWindowsAPI(t)
			apiMock.On("GetMaintenanceWindow", mock.Anything, "projectID").Return(admin.GetMaintenanceWindowApiRequest{ApiService: apiMock})
			apiMock.On("GetMaintenanceWindowExecute", mock.Anything).Return(tc.apiResponse, &http.Response{StatusCode: http.StatusOK}, nil)

			d := schema.TestResourceDataRaw(t, maintenancewindow.Resource().Schema, map[string]any{"protected_hours": tc.stateProtectedHours})
			d.SetId("projectID")

			diags := maintenancewindow.Resource().ReadContext(context.Background(), d, &config.MongoDBClient{AtlasV2: &admin.APIClient{MaintenanceWindowsAPI: apiMock}})
			assert.False(t, diags.HasError(), "unexpected errors: %v", diags)
			assert.Equal(t, tc.expectedProtectedHL, d.Get("protected_hours").([]any))
		})
	}
}
