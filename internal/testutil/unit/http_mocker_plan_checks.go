package unit

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

const resourceID = "mongodbatlas_example.this"

type requestHandlerSwitch struct {
	useManualHandler *bool
}

func (r *requestHandlerSwitch) CheckPlan(_ context.Context, _ plancheck.CheckPlanRequest, _ *plancheck.CheckPlanResponse) {
	*r.useManualHandler = true
}

// TODO: Extract all import parameters instead from the template file
func MockPlanChecksAndRun(t *testing.T, mockConfig MockHTTPDataConfig, importInput, importConfig, importResourceName string, testStep *resource.TestStep) {
	t.Helper()
	exampleResourceConfig := CreateExampleConfig(importInput)
	fullImportConfig := exampleResourceConfig + importConfig
	testStep.Config = exampleResourceConfig + testStep.Config
	testCase := &resource.TestCase{
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: exampleResourceConfig,
				Check:  resource.TestCheckResourceAttr("mongodbatlas_example.this", "import_id", importInput),
			},
			{
				Config:             fullImportConfig,
				ResourceName:       importResourceName,
				ImportStateIdFunc:  importStateImportID(),
				ImportState:        true,
				ImportStatePersist: true, // save the state to use it in the next plan
			},
			*testStep,
		},
	}
	fillMockDataTemplate(t, exampleResourceConfig, fullImportConfig, testStep.Config)
	useManualHandler := false
	testCase.Steps[2].ConfigPlanChecks.PreApply = append(testCase.Steps[2].ConfigPlanChecks.PreApply, &requestHandlerSwitch{useManualHandler: &useManualHandler})
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
				useManualHandler = true
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

func CreateExampleConfig(importID string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_example" "this" {
  import_id = %[1]q
}
`, importID)
}
func importStateImportID() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceID]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return rs.Primary.Attributes["import_id"], nil
	}
}

func fillMockDataTemplate(t *testing.T, exampleConfig, fullImportConfig, planCheckConfig string) {
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
	addPlanCheckStep(t, exampleConfig, fullImportConfig, planCheckConfig, mockDataPath)
}

func addPlanCheckStep(t *testing.T, exampleConfig, fullImportConfig, planCheckConfig string, mockDataPath string) {
	parseData := ReadMockData(t, []string{exampleConfig, fullImportConfig})
	parseData.Steps = append(parseData.Steps, StepRequests{
		Config:           Literal(planCheckConfig),
		RequestResponses: parseData.Steps[1].RequestResponses,
	})
	finalYaml, err := ConfigYaml(parseData)
	require.NoError(t, err)
	err = os.WriteFile(mockDataPath, []byte(finalYaml), 0644)
	require.NoError(t, err)
}
