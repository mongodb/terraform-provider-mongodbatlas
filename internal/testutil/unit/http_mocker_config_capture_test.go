package unit_test

import (
	_ "embed"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed testdata/CaptureTest/createCluster.json
	createClusterReqBody string
	//go:embed testdata/CaptureTest/createClusterResponse.json
	createClusterRespBody string
	//go:embed testdata/CaptureTest/getClusterIdle.json
	getClusterIdleRespBody string
	//go:embed testdata/CaptureTest/getContainersAws.json
	getContainersAws string
	//go:embed testdata/CaptureTest/getContainersAzure.json
	getContainersAzure string
)

func TestFailedFilename(t *testing.T) {
	normalName := unit.MockConfigFilePath(t)
	assert.Equal(t, "testdata/TestFailedFilename.yaml", normalName)
	failedName := unit.FailedFilename(normalName)
	assert.Contains(t, failedName, "testdata/TestFailedFilename_failed")
}

func TestCaptureMockConfigClientModifier_clusterExample(t *testing.T) {
	capturedData := unit.NewMockHTTPData(t, 3, []string{"resource \"dummy\" \"test\"{\n  step = 1\n  someString = \"my-string\"\n}", "", "resource \"dummy\" \"test\"{\n  step = 3\n}"})
	clientModifier := unit.NewCaptureMockConfigClientModifier(t, &unit.MockHTTPDataConfig{QueryVars: []string{"providerName"}}, capturedData)
	transport := httpmock.NewMockTransport()
	client := http.Client{Transport: transport}
	err := clientModifier.ModifyHTTPClient(&client)
	require.NoError(t, err)

	// Step 1: Create a cluster
	clientModifier.IncreaseStepNumber()
	responder1 := httpmock.NewStringResponder(201, createClusterRespBody)
	transport.RegisterRegexpResponder("POST", regexp.MustCompile(".*"), responder1)
	createRequest := request("POST", "/api/atlas/v2/groups/g1/clusters", createClusterReqBody)
	resp, err := client.Do(createRequest)
	require.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	createResponse := parseMapStringAny(t, resp)
	assert.Equal(t, "test-acc-tf-c-7871793563057636102", createResponse["name"])

	// Step 2: Read cluster
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

	// Step 3: Read containers, capture query args
	clientModifier.IncreaseStepNumber()
	containersGetResponses := []string{getContainersAws, getContainersAzure}
	containersExpectedIDs := []string{"6746ceedaef48d1cb265896b", "6746cefbaef48d1cb2658bbb"}
	containersGetPaths := []string{
		"/api/atlas/v2/groups/6746cee66f62fc3c122a3b82/containers?includeCount=true&itemsPerPage=100&pageNum=1&providerName=AWS",
		"/api/atlas/v2/groups/6746cee66f62fc3c122a3b82/containers?includeCount=true&itemsPerPage=100&pageNum=1&providerName=AZURE",
	}
	for i := range containersGetResponses {
		transport.Reset()
		responder3 := httpmock.NewStringResponder(200, containersGetResponses[i])
		transport.RegisterRegexpResponder("GET", regexp.MustCompile(".*"), responder3)
		getRequest := request("GET", containersGetPaths[i], "")
		resp, err = client.Do(getRequest)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		getResponse := parseMapStringAny(t, resp)
		assert.Equal(t, containersExpectedIDs[i], getResponse["results"].([]any)[0].(map[string]any)["id"])
	}

	g := goldie.New(t, goldie.WithTestNameForDir(true), goldie.WithNameSuffix(".yaml"))
	configYaml, err := clientModifier.ConfigYaml()
	require.NoError(t, err)
	g.Assert(t, t.Name(), []byte(configYaml))
}

func Reformat(t *testing.T, filePath string) string {
	t.Helper()
	data, err := unit.ParseTestDataConfigYAML(filePath)
	require.NoError(t, err)
	initialYaml := strings.Builder{}
	e := yaml.NewEncoder(&initialYaml)
	e.SetIndent(1)
	err = e.Encode(data)
	require.NoError(t, err)
	return initialYaml.String()
}

// Manual test used to test reformatting of all yaml files in the testdata directory
func TestReformatConfigs(t *testing.T) {
	testDataPath := os.Getenv("TEST_DATA_PATH")
	if testDataPath == "" {
		t.Skip("TEST_DATA_PATH is not set")
	}
	matches, err := filepath.Glob(path.Join(testDataPath, "*.yaml"))
	require.NoError(t, err)
	for _, p := range matches {
		formatted := Reformat(t, p)
		err = os.WriteFile(p, []byte(formatted), 0o600)
		require.NoError(t, err)
	}
}
