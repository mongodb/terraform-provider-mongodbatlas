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

const resourceID = "mongodbatlas_example.this"

var expectedError = errors.New("avoid full apply by raising an expected error")

type requestHandlerSwitch struct {
	useManualHandler *bool
}

func (r *requestHandlerSwitch) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	*r.useManualHandler = true
	resp.Error = expectedError
}

func MockPlanChecksAndRun(t *testing.T, mockConfig MockHTTPDataConfig, importResourceName, importID, planConfig string, checks []plancheck.PlanCheck) {
	t.Helper()
	importConfig := fillMockDataTemplate(t, planConfig)
	useManualHandler := false
	checks = append(checks, &requestHandlerSwitch{useManualHandler: &useManualHandler})
	testCase := &resource.TestCase{
		IsUnitTest:                true,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config:             importConfig,
				ResourceName:       importResourceName,
				ImportStateId:      importID, // static ID to import
				ImportState:        true,
				ImportStatePersist: true, // save the state to use it in the next plan
			},
			{
				Config: planConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: checks,
				},
				ExpectError: regexp.MustCompile(fmt.Sprintf("^Pre-apply plan check\\(s\\) failed:\n%s$", expectedError.Error())), // To avoid doing a full apply
			},
		},
	}
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

func fillMockDataTemplate(t *testing.T, planCheckConfig string) string{
	t.Helper()
	templatePath := fmt.Sprintf("testdata/%s.tmpl.yaml", t.Name())
	templateContent, err := os.ReadFile(templatePath)
	require.NoError(t, err)
	responseDir := fmt.Sprintf("testdata/%s", t.Name())
	responsePaths, err := filepath.Glob(path.Join(responseDir, "*.json"))
	require.NoError(t, err)
	for _, testFile := range responsePaths {
		testFileContent, err := os.ReadFile(testFile)
		require.NoError(t, err)
		testFileContent = bytes.ReplaceAll(testFileContent, []byte(`"`), []byte(`\"`))
		testFileContent = bytes.ReplaceAll(testFileContent, []byte("\n"), []byte(`\n`))
		templateContent = bytes.ReplaceAll(templateContent, []byte(filepath.Base(testFile)), testFileContent)
	}
	mockDataPath := fmt.Sprintf("testdata/%s.yaml", t.Name())
	err = os.WriteFile(mockDataPath, templateContent, 0644)
	require.NoError(t, err)
	fullImportConfigBytes, err := os.ReadFile(path.Join(responseDir, "main.tf"))
	require.NoError(t, err)
	fullImportConfig := string(fullImportConfigBytes)
	addPlanCheckStepAndReadImportConfig(t, fullImportConfig, planCheckConfig, mockDataPath)
	return fullImportConfig
}

func addPlanCheckStepAndReadImportConfig(t *testing.T, fullImportConfig, planCheckConfig string, mockDataPath string) {
	parseData := ReadMockData(t, []string{fullImportConfig})
	parseData.Steps = append(parseData.Steps, StepRequests{
		Config:           Literal(planCheckConfig),
		RequestResponses: parseData.Steps[0].RequestResponses,
	})
	finalYaml, err := ConfigYaml(parseData)
	require.NoError(t, err)
	err = os.WriteFile(mockDataPath, []byte(finalYaml), 0644)
	require.NoError(t, err)
}
