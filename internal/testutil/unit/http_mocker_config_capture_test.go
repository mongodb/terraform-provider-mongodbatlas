package unit_test

import (
	_ "embed"
	"net/http"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/CaptureTest/createCluster.json
	createClusterReqBody string
	//go:embed testdata/CaptureTest/createClusterResponse.json
	createClusterRespBody string
	//go:embed testdata/CaptureTest/getClusterIdle.json
	getClusterIdleRespBody string
)

func TestFailedFilename(t *testing.T) {
	normalName := unit.MockConfigFilePath(t)
	assert.Equal(t, "testdata/TestFailedFilename.yaml", normalName)
	failedName := unit.FailedFilename(normalName)
	assert.Contains(t, failedName, "testdata/TestFailedFilename_failed")
}

func TestCaptureMockConfigClientModifier_clusterExample(t *testing.T) {
	t.Setenv(unit.EnvNameHTTPMockerCapture, "true")
	clientModifier := unit.NewCaptureMockConfigClientModifier(t, 2, nil)
	transport := httpmock.NewMockTransport()
	client := http.Client{Transport: transport}
	err := clientModifier.ModifyHTTPClient(&client)
	require.NoError(t, err)
	clientModifier.IncreaseStepNumber()
	responder1 := httpmock.NewStringResponder(201, createClusterRespBody)
	transport.RegisterRegexpResponder("POST", regexp.MustCompile(".*"), responder1)
	createRequest := request("POST", "/api/atlas/v2/groups/g1/clusters", createClusterReqBody)
	resp, err := client.Do(createRequest)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	createResponse := parseMapStringAny(t, resp)
	assert.Equal(t, "test-acc-tf-c-7871793563057636102", createResponse["name"])

	clientModifier.IncreaseStepNumber()
	getResponses := []string{createClusterRespBody, getClusterIdleRespBody, getClusterIdleRespBody}
	expectedState := []string{"CREATING", "IDLE", "IDLE"}
	for i := range getResponses {
		transport.Reset()
		responder2 := httpmock.NewStringResponder(200, getResponses[i])
		transport.RegisterRegexpResponder("GET", regexp.MustCompile(".*"), responder2)
		getRequest := request("GET", "/api/atlas/v2/groups/6746ceed6f62fc3c122a3e0e/clusters/test-acc-tf-c-7871793563057636102", "")
		resp, err = client.Do(getRequest)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		getResponse := parseMapStringAny(t, resp)
		assert.Equal(t, expectedState[i], getResponse["stateName"])
	}

	g := goldie.New(t, goldie.WithTestNameForDir(true), goldie.WithNameSuffix(".yaml"))
	configYaml, err := clientModifier.ConfigYaml()
	require.NoError(t, err)
	g.Assert(t, t.Name(), []byte(configYaml))
}
