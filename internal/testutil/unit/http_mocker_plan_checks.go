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
	MockedClusterName = "mocked-cluster"
	MockedProjectID = "111111111111111111111111"
)
var (
	expectedError = errors.New("avoid full apply by raising an expected error")

	importIDMapping = map[string]string{
		ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs: fmt.Sprintf("%s-%s", MockedProjectID, MockedClusterName),
	}
	// later this could be inferred
	importResourceNameMapping = map[string]string{
		ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs: "mongodbatlas_advanced_cluster.test",
	}
)


type requestHandlerSwitch struct {
	useManualHandler *bool
}

func (r *requestHandlerSwitch) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	*r.useManualHandler = true
	resp.Error = expectedError
}

func NewMockPlanChecksConfig(t *testing.T, mockConfig MockHTTPDataConfig, importName string) MockPlanChecksConfig {
	t.Helper()
	importID := importIDMapping[importName]
	require.NotEmpty(t, importID, "import ID not found for import name: %s", importName)
	resourceName := importResourceNameMapping[importName]
	require.NotEmpty(t, resourceName, "resource name not found for import name: %s", importName)
	config := MockPlanChecksConfig{
		ImportName: importName,
		MockConfig: mockConfig,
		ImportID: importID,
		ResourceName: resourceName,
	}
	return config
}

type MockPlanChecksConfig struct {
	Checks       []plancheck.PlanCheck
	ImportID     string
	ResourceName string
	MockConfig   MockHTTPDataConfig
	ImportName   string
	Name         string
}

func (m *MockPlanChecksConfig) WithNameAndChecks(name string, checks []plancheck.PlanCheck) MockPlanChecksConfig {
	return MockPlanChecksConfig{
		Checks:       checks,
		ImportName:   m.ImportName,
		ImportID:     m.ImportID,
		ResourceName: m.ResourceName,
		MockConfig:   m.MockConfig,
		Name:         name,
	}
}

func MockPlanChecksAndRun(t *testing.T, runConfig MockPlanChecksConfig) {
	t.Helper()
	importConfig, planConfig, mockDataPath := fillMockDataTemplate(t, runConfig.ImportName, runConfig.Name)
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
				Config: planConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: runConfig.Checks,
				},
				ExpectError: regexp.MustCompile(fmt.Sprintf("^Pre-apply plan check\\(s\\) failed:\n%s$", expectedError.Error())), // To avoid doing a full apply
			},
		},
	}
	mockConfig := runConfig.MockConfig
	mockConfig.FilePathOverride = mockDataPath
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

func fillMockDataTemplate(t *testing.T, importName string, planName string) (importConfig, planConfig, mockDataFilePath string) {
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
	mockDataPath := fmt.Sprintf("testdata/%s_%s.yaml", importName, planName)
	err = os.WriteFile(mockDataPath, templateContent, 0644)
	require.NoError(t, err)
	fullImportConfigBytes, err := os.ReadFile(path.Join(responseDir, "main.tf"))
	require.NoError(t, err)
	fullImportConfig := string(fullImportConfigBytes)
	planCheckConfigBytes, err := os.ReadFile(path.Join(responseDir, fmt.Sprintf("main_%s.tf", planName)))
	require.NoError(t, err)
	planCheckConfig := string(planCheckConfigBytes)
	addPlanCheckStepAndReadImportConfig(t, fullImportConfig, planCheckConfig, mockDataPath)
	return fullImportConfig, planCheckConfig, mockDataPath
}

func addPlanCheckStepAndReadImportConfig(t *testing.T, fullImportConfig, planCheckConfig string, mockDataPath string) {
	parseData := ReadMockDataFile(t, mockDataPath, []string{fullImportConfig})
	parseData.Steps = append(parseData.Steps, StepRequests{
		Config:           Literal(planCheckConfig),
		RequestResponses: parseData.Steps[0].RequestResponses,
	})
	finalYaml, err := ConfigYaml(parseData)
	require.NoError(t, err)
	err = os.WriteFile(mockDataPath, []byte(finalYaml), 0644)
	require.NoError(t, err)
}
