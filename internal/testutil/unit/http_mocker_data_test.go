package unit_test

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

func TestMockHTTPData_UpdateVariables(t *testing.T) {
	mockData := unit.NewMockHTTPData(t, 1, []string{""})
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
	mockData := unit.NewMockHTTPData(t, 1, []string{""})
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

func TestMockDataExtractVars(t *testing.T) {
	config1 := projectAdvClusterExample
	config2 := strings.ReplaceAll(config1, "test-acc-tf-c-8022584361920682288", "test-acc-tf-c-8022584361920682289")
	mockData := unit.NewMockHTTPData(t, 2, []string{config1, config2})
	expected := map[string]string{
		"clusterName":  "test-acc-tf-c-8022584361920682288",
		"clusterName2": "test-acc-tf-c-8022584361920682289",
		"orgId":        "65def6ce0f722a1507105aa5",
		"projectName":  "test-acc-tf-p-664077766951329406",
	}
	assert.Equal(t, expected, mockData.Variables)
}

var expectedDump = `variables: {}
steps:
  - diff_requests: []
    request_responses: []
  - config: |-
      data "mongodbatlas_advanced_cluster" "test" {
        project_id                     = mongodbatlas_advanced_cluster.test.project_id
        name                           = mongodbatlas_advanced_cluster.test.name
      }
      data "mongodbatlas_advanced_clusters" "test" {
        project_id                     = mongodbatlas_advanced_cluster.test.project_id
      }
    diff_requests: []
    request_responses: []
`

var tfDsString = "\ndata \"mongodbatlas_advanced_cluster\" \"test\" {\n  project_id                     = mongodbatlas_advanced_cluster.test.project_id\n  name                           = mongodbatlas_advanced_cluster.test.name\n}\ndata \"mongodbatlas_advanced_clusters\" \"test\" {\n  project_id                     = mongodbatlas_advanced_cluster.test.project_id\n}\n \n"

func TestDumpingConfigUsesLiteralStyle(t *testing.T) {
	mockData := unit.NewMockHTTPData(t, 2, []string{"", tfDsString})
	initialYaml := strings.Builder{}
	e := yaml.NewEncoder(&initialYaml)
	e.SetIndent(1)
	err := e.Encode(mockData)
	require.NoError(t, err)
	assert.Equal(t, expectedDump, initialYaml.String())
}
