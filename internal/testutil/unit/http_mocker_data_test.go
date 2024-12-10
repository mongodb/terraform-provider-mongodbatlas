package unit_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockHTTPData_UpdateVariables(t *testing.T) {
	mockData := unit.NewMockHTTPData(1)
	err := mockData.UpdateVariables(t, map[string]string{"groupId": "g1", "clusterName": "c1"})
	require.NoError(t, err)
	require.Equal(t, "g1", mockData.Variables["groupId"])
	require.Equal(t, "c1", mockData.Variables["clusterName"])
	err = mockData.UpdateVariables(t, map[string]string{"groupId": "g2", "clusterName": "c2"})
	changeError, ok := err.(*unit.VariablesChangedError)
	require.True(t, ok)
	require.Len(t, changeError.Changes, 2)
	assert.Equal(t, map[string]string{"groupId": "groupId2", "clusterName": "clusterName2"}, changeError.ChangedNamesMap())
	assert.Equal(t, map[string]string{"g1": "g2", "c1": "c2"}, changeError.ChangedValuesMap())
	assert.Equal(t, map[string]string{"clusterName": "c1", "clusterName2": "c2", "groupId": "g1", "groupId2": "g2"}, mockData.Variables)
}

func TestMockHTTPData_AddRoundtrip(t *testing.T) {
	mockData := unit.NewMockHTTPData(1)
	rt := &unit.RoundTrip{
		Variables:  map[string]string{"groupId": "g1", "clusterName": "c1"},
		StepNumber: 1,
		Request:    unit.RequestInfo{},
	}
	err := mockData.AddRoundtrip(t, rt, false)
	require.NoError(t, err)
	require.Equal(t, "g1", mockData.Variables["groupId"])
	require.Equal(t, "c1", mockData.Variables["clusterName"])
	rt2 := &unit.RoundTrip{
		Variables:  map[string]string{"groupId": "g2", "clusterName": "c2"},
		StepNumber: 1,
		Request:    unit.RequestInfo{},
	}
	err = mockData.AddRoundtrip(t, rt2, false)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"clusterName": "c1", "clusterName2": "c2", "groupId": "g1", "groupId2": "g2"}, mockData.Variables)
}
