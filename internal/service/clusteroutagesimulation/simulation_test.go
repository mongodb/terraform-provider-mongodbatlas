package clusteroutagesimulation_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
	"go.mongodb.org/atlas-sdk/v20250312018/mockadmin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/retrystrategy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clusteroutagesimulation"
)

const (
	testProjectID   = "test-project-id"
	testClusterName = "test-cluster"
)

var fastTimeConfig = retrystrategy.TimeConfig{
	Timeout:    5 * time.Second,
	MinTimeout: 10 * time.Millisecond,
	Delay:      0,
}

func TestParseActionTimeout(t *testing.T) {
	tests := []struct {
		name        string
		input       types.String
		expected    time.Duration
		expectError bool
	}{
		{
			name:     "null returns default",
			input:    types.StringNull(),
			expected: 25 * time.Minute,
		},
		{
			name:     "empty string returns default",
			input:    types.StringValue(""),
			expected: 25 * time.Minute,
		},
		{
			name:     "valid minutes",
			input:    types.StringValue("10m"),
			expected: 10 * time.Minute,
		},
		{
			name:     "valid hours and minutes",
			input:    types.StringValue("1h30m"),
			expected: 90 * time.Minute,
		},
		{
			name:        "invalid string returns error",
			input:       types.StringValue("not-a-duration"),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := clusteroutagesimulation.ParseActionTimeout(tc.input)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, got)
			}
		})
	}
}

func TestSimulateOutage_Success(t *testing.T) {
	m := mockadmin.NewClusterOutageSimulationApi(t)

	simulatingState := "SIMULATING"
	m.EXPECT().StartOutageSimulation(mock.Anything, testProjectID, testClusterName, mock.Anything).
		Return(admin.StartOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().StartOutageSimulationExecute(mock.Anything).
		Return(nil, nil, nil)

	m.EXPECT().GetOutageSimulation(mock.Anything, testProjectID, testClusterName).
		Return(admin.GetOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().GetOutageSimulationExecute(mock.Anything).
		Return(&admin.ClusterOutageSimulation{State: &simulatingState}, nil, nil)

	filters := []admin.AtlasClusterOutageSimulationOutageFilter{
		{CloudProvider: strPtr("AWS"), RegionName: strPtr("US_EAST_1")},
	}

	err := clusteroutagesimulation.SimulateOutage(t.Context(), m, testProjectID, testClusterName, filters, true, fastTimeConfig)
	require.NoError(t, err)
}

func TestSimulateOutage_APIError(t *testing.T) {
	m := mockadmin.NewClusterOutageSimulationApi(t)

	m.EXPECT().StartOutageSimulation(mock.Anything, testProjectID, testClusterName, mock.Anything).
		Return(admin.StartOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().StartOutageSimulationExecute(mock.Anything).
		Return(nil, &http.Response{StatusCode: http.StatusInternalServerError}, errors.New("internal server error"))

	err := clusteroutagesimulation.SimulateOutage(t.Context(), m, testProjectID, testClusterName, nil, true, fastTimeConfig)
	require.Error(t, err)
	assert.Contains(t, err.Error(), testProjectID)
	assert.Contains(t, err.Error(), testClusterName)
}

func TestStopSimulation_Success(t *testing.T) {
	m := mockadmin.NewClusterOutageSimulationApi(t)

	deletedState := "DELETED"
	m.EXPECT().EndOutageSimulation(mock.Anything, testProjectID, testClusterName).
		Return(admin.EndOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().EndOutageSimulationExecute(mock.Anything).
		Return(nil, nil, nil)

	m.EXPECT().GetOutageSimulation(mock.Anything, testProjectID, testClusterName).
		Return(admin.GetOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().GetOutageSimulationExecute(mock.Anything).
		Return(&admin.ClusterOutageSimulation{State: &deletedState}, nil, nil)

	err := clusteroutagesimulation.StopSimulation(t.Context(), m, testProjectID, testClusterName, fastTimeConfig)
	require.NoError(t, err)
}

func TestStopSimulation_APIError(t *testing.T) {
	m := mockadmin.NewClusterOutageSimulationApi(t)

	m.EXPECT().EndOutageSimulation(mock.Anything, testProjectID, testClusterName).
		Return(admin.EndOutageSimulationApiRequest{ApiService: m})
	m.EXPECT().EndOutageSimulationExecute(mock.Anything).
		Return(nil, &http.Response{StatusCode: http.StatusInternalServerError}, errors.New("internal server error"))

	err := clusteroutagesimulation.StopSimulation(t.Context(), m, testProjectID, testClusterName, fastTimeConfig)
	require.Error(t, err)
	assert.Contains(t, err.Error(), testProjectID)
	assert.Contains(t, err.Error(), testClusterName)
}

func strPtr(s string) *string {
	return &s
}
