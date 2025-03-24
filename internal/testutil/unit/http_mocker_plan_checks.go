package unit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const (
	ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs = "ClusterTwoRepSpecsWithAutoScalingAndSpecs"
	ImportNameClusterReplicasetOneRegion                = "ClusterReplicasetOneRegion"
	MockedClusterName                                   = "mocked-cluster"
	MockedProjectID                                     = "111111111111111111111111"
)

var (
	errToSkipApply  = errors.New("avoid full apply by raising an expected error")
	clusterImportID = fmt.Sprintf("%s-%s", MockedProjectID, MockedClusterName)

	importIDMapping = map[string]string{
		ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs: fmt.Sprintf("%s-%s", MockedProjectID, MockedClusterName),
		ImportNameClusterReplicasetOneRegion:                clusterImportID,
	}
	// later this could be inferred when reading the src main.tf
	importResourceNameMapping = map[string]string{
		ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs: "mongodbatlas_advanced_cluster.test",
		ImportNameClusterReplicasetOneRegion:                "mongodbatlas_advanced_cluster.test",
	}
)

func NewMockPlanChecksConfig(t *testing.T, mockConfig *MockHTTPDataConfig, importName string) *MockPlanChecksConfig {
	t.Helper()
	importID := importIDMapping[importName]
	require.NotEmpty(t, importID, "import ID not found for import name: %s", importName)
	resourceName := importResourceNameMapping[importName]
	require.NotEmpty(t, resourceName, "resource name not found for import name: %s", importName)
	return &MockPlanChecksConfig{
		ImportName:   importName,
		MockConfig:   *mockConfig,
		ImportID:     importID,
		ResourceName: resourceName,
	}
}

type MockPlanChecksConfig struct {
	ImportID       string
	ResourceName   string
	ImportName     string
	ConfigFilename string
	Checks         []plancheck.PlanCheck
	MockConfig     MockHTTPDataConfig
}

func (m *MockPlanChecksConfig) WithPlanCheckTest(testConfig PlanCheckTest) *MockPlanChecksConfig {
	return &MockPlanChecksConfig{
		Checks:         testConfig.Checks,
		ConfigFilename: testConfig.ConfigFilename,
		ImportName:     m.ImportName,
		ImportID:       m.ImportID,
		ResourceName:   m.ResourceName,
		MockConfig:     m.MockConfig,
	}
}

type PlanCheckTest struct {
	ConfigFilename string // .tf filename
	Checks         []plancheck.PlanCheck
}

func RunPlanCheckTests(t *testing.T, baseConfig *MockPlanChecksConfig, tests []PlanCheckTest) {
	t.Helper()
	for _, testCase := range tests {
		t.Run(testCase.ConfigFilename, func(t *testing.T) {
			MockPlanChecksAndRun(t, baseConfig.WithPlanCheckTest(testCase))
		})
	}
}

// MockPlanChecksAndRun creates and runs a UnitTest enabled TestCase for Read to State checks and PlanModifier logic.
// The 1st step is always Import
// The 2nd step is always Plan with runConfig.Checks run. Note: No Update logic is executed as we exit after the PlanModifier has run.
// Instead of having to store full mock data files, we re-use the same GET requests from the directory testdata/{runConfig.ImportName}/import_*.json
// Together with the extra step in `testdata/{ImportName}/main_{runConfig.Name}.tf` we fill the template: testdata/{runConfig.ImportName}.tmpl.yaml
func MockPlanChecksAndRun(t *testing.T, runConfig *MockPlanChecksConfig) {
	t.Helper()
	importConfig, planConfig, mockDataPath := fillMockDataTemplate(t, runConfig.ImportName, runConfig.ConfigFilename)
	t.Cleanup(func() {
		require.NoError(t, os.Remove(mockDataPath))
	})
	useManualHandler := false
	runConfig.Checks = append(runConfig.Checks, &requestHandlerSwitch{useManualHandler: &useManualHandler})
	testCase := &resource.TestCase{
		IsUnitTest:                true,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config:             importConfig,
				ResourceName:       runConfig.ResourceName,
				ImportStateId:      runConfig.ImportID, // static ID to import
				ImportState:        true,
				ImportStatePersist: true, // save the state to use it in the next plan
			},
			{
				// Specifying Config AND using PlanChecks are only available when running test in `Config` mode (see https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/teststep#test-modes for the different modes)
				// To avoid doing a full apply we return a known error in the `requestHandlerSwitch` and Expect it in ExpectError
				Config: planConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: runConfig.Checks,
				},
				ExpectError: regexp.MustCompile(fmt.Sprintf("^Pre-apply plan check\\(s\\) failed:\n%s$", errToSkipApply.Error())), // Notice full match using ^ and $ in case some checks also fails
			},
		},
	}
	mockConfig := runConfig.MockConfig
	// Use FilePathOverride to avoid having to follow the standard `MockConfigFilePath` which depends on the test name
	mockConfig.FilePathOverride = mockDataPath
	// Custom RequestHandler that runs after PlanModifier is done, to avoid `mock response not found`` errors in test cleanup functions
	mockConfig.RequestHandler = func(defaulthHandler RequestHandler, req *http.Request, method string) (*http.Response, error) {
		customHandler := func(req *http.Request, method string) (*http.Response, error) {
			switch method {
			case "GET":
				notFoundResponder, err := httpmock.NewJsonResponder(404, map[string]any{
					"errorCode": "RESOURCE_NOT_FOUND",
				})
				require.NoError(t, err)
				return notFoundResponder(req)
			case "DELETE":
				return httpmock.NewStringResponder(202, "")(req)
			}
			return nil, fmt.Errorf("plan check responder doesn't have logic to handle, method: %s, url: %s", method, req.URL)
		}
		if useManualHandler {
			return customHandler(req, method)
		}
		return defaulthHandler(req, method)
	}
	err := enableReplayForTestCase(
		t,
		&mockConfig,
		testCase,
	)
	require.NoError(t, err)
	resource.ParallelTest(t, *testCase)
}

type requestHandlerSwitch struct {
	useManualHandler *bool
}

func (r *requestHandlerSwitch) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	*r.useManualHandler = true
	resp.Error = errToSkipApply
}

func fillMockDataTemplate(t *testing.T, importName, planConfigFilename string) (importConfig, planCheckConfig, mockDataFilePath string) {
	t.Helper()
	templatePath := fmt.Sprintf("testdata/%s.tmpl.yaml", importName)
	templateContent, err := os.ReadFile(templatePath)
	require.NoError(t, err)
	responseDir := fmt.Sprintf("testdata/%s", importName)
	responsePaths, err := filepath.Glob(path.Join(responseDir, "*.json"))
	require.NoError(t, err)
	for _, testFile := range responsePaths {
		testFileContent, err := os.ReadFile(testFile)
		require.NoError(t, err)
		testFileContent = bytes.ReplaceAll(testFileContent, []byte(`"`), []byte(`\"`))
		testFileContent = bytes.ReplaceAll(testFileContent, []byte("\n"), []byte(`\n`))
		templateContent = bytes.ReplaceAll(templateContent, []byte(filepath.Base(testFile)), testFileContent)
	}
	mockDataPath := fmt.Sprintf("testdata/%s_%s.yaml", importName, planConfigFilename)
	err = os.WriteFile(mockDataPath, templateContent, 0o600)
	require.NoError(t, err)
	fullImportConfigBytes, err := os.ReadFile(path.Join(responseDir, "main.tf"))
	require.NoError(t, err)
	fullImportConfig := string(fullImportConfigBytes)
	planCheckConfigBytes, err := os.ReadFile(path.Join(responseDir, planConfigFilename))
	require.NoError(t, err)
	planCheckConfig = string(planCheckConfigBytes)
	addPlanCheckStepAndReadImportConfig(t, fullImportConfig, planCheckConfig, mockDataPath)
	return fullImportConfig, planCheckConfig, mockDataPath
}

func addPlanCheckStepAndReadImportConfig(t *testing.T, fullImportConfig, planCheckConfig, mockDataPath string) {
	t.Helper()
	parseData := ReadMockDataFile(t, mockDataPath, []string{fullImportConfig})
	parseData.Steps = append(parseData.Steps, StepRequests{
		Config:           Literal(planCheckConfig),
		RequestResponses: parseData.Steps[0].RequestResponses,
	})
	finalYaml, err := ConfigYaml(parseData)
	require.NoError(t, err)
	err = os.WriteFile(mockDataPath, []byte(finalYaml), 0o600)
	require.NoError(t, err)
}
