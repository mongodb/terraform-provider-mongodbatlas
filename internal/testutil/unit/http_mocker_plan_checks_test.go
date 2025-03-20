package unit_test

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/require"
)

const (
	pkgAdvancedCluster = "advancedcluster"
	pkgAdvancedClusterTPF = "advancedclustertpf"
	pkgRelPath = "internal/service"
	prefixName             = "TestMockPlanChecks_"
)

func ensureDir(t *testing.T, dir string) string{
	t.Helper()
	if !unit.FileExist(dir) {
		err := os.Mkdir(dir, 0755)
		require.NoError(t, err)
	}
	return dir
}

func CreateImportData(t *testing.T, srcMockFile, destMockFile, importName string, stepNr int, newVars map[string]string) {
	t.Helper()
	destTestData, _ := path.Split(destMockFile)
	require.True(t, strings.HasSuffix(destTestData, "/testdata/"))
	ensureDir(t, destTestData)
	destOutputDir := ensureDir(t, path.Join(destTestData, importName))

	templateMockHTTPData := createImportMockData(t, srcMockFile, destOutputDir, stepNr, newVars)
	require.NoError(t, templateMockHTTPData.UpdateVariablesIgnoreChanges(t, newVars))

	templateYaml, err := unit.ConfigYaml(templateMockHTTPData)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(destMockFile, []byte(templateYaml), 0644))
}

func createImportMockData(t *testing.T, srcMockFile, destOutputDir string, stepNr int, newVars map[string]string) *unit.MockHTTPData {
	data, err := unit.ParseTestDataConfigYAML(srcMockFile)
	require.NoError(t, err)
	relevantStep := data.Steps[stepNr-1]
	getRequestsInStep := []unit.RequestInfo{}
	for _, req := range relevantStep.RequestResponses {
		if req.Method == "GET" {
			getRequestsInStep = append(getRequestsInStep, req)
		}
	}
	replaceVarsInConfig := func(config string) unit.Literal {
		for key, oldValue := range data.Variables {
			newValue, ok := newVars[key]
			require.True(t, ok, "Missing variable %s in newVars", key)
			config = strings.ReplaceAll(config, oldValue, newValue)
		}
		return unit.Literal(config)
	}
	templateMockHTTPData := unit.MockHTTPData{
		Steps: []unit.StepRequests{
			{
				Config: replaceVarsInConfig(string(relevantStep.Config)),
			},
		},
		Variables: newVars,
	}
	for _, req := range getRequestsInStep {
		lastResponse := req.Responses[len(req.Responses)-1]
		jsonFileName := strings.ReplaceAll(fmt.Sprintf("import_%s.json", req.IdShort()), "/", "_")
		jsonFilePath := path.Join(destOutputDir, jsonFileName)
		err = os.WriteFile(jsonFilePath, []byte(lastResponse.Text), 0644)
		require.NoError(t, err)
		templateReqResponse := unit.RequestInfo{
			Path:    req.Path,
			Method:  req.Method,
			Version: req.Version,
			Text:    req.Text,
			Responses: []unit.StatusText{
				{
					Status: lastResponse.Status,
					Text:   jsonFileName,
				},
			},
		}
		templateMockHTTPData.Steps[0].RequestResponses = append(templateMockHTTPData.Steps[0].RequestResponses, templateReqResponse)
	}
	return &templateMockHTTPData
}

type importNameConfig struct {
	TestName string
	Step int
	VariableReplacments map[string]string
	SrcPackage string
	DestPackage string
}

func TestConvertMockableTests(t *testing.T) {
	for importName, config := range map[string]importNameConfig{
		unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs: {
			TestName: "TestAccMockableAdvancedCluster_removeBlocksFromConfig",
			Step: 1,
			VariableReplacments: map[string]string{
				"clusterName": unit.MockedClusterName,
				"groupId": unit.MockedProjectID,
			},
			SrcPackage: pkgAdvancedCluster,
			DestPackage: pkgAdvancedClusterTPF,
		},
	} {
		srcTestdata := unit.RepoPath(path.Join(pkgRelPath, config.SrcPackage, "testdata"))
		destTestdata := unit.RepoPath(path.Join(pkgRelPath, config.DestPackage, "testdata"))
		ensureDir(t, destTestdata)
		gitIgnorePath := path.Join(destTestdata, ".gitignore")
		gitIgnoreExpectedContent := fmt.Sprintf("%s_*.yaml\n", importName)
		if unit.FileExist(gitIgnorePath) {
			gitIgnoreContent, err := os.ReadFile(gitIgnorePath)
			require.NoError(t, err)
			require.Contains(t, string(gitIgnoreContent), gitIgnoreExpectedContent, "Missing:\n%s in %s\nused to avoid pushing dynamic planCheck files", gitIgnoreExpectedContent, gitIgnorePath)
		} else {
			require.NoError(t, os.WriteFile(gitIgnorePath, []byte(gitIgnoreExpectedContent), 0644))
		}
		srcTestdataPath := path.Join(srcTestdata, config.TestName+".yaml")
		destTestdataPath := path.Join(destTestdata, importName+".tmpl.yaml")
		t.Logf("Converting %s step %d to %s", srcTestdataPath, config.Step, destTestdataPath)
		CreateImportData(t, srcTestdataPath, destTestdataPath, importName, config.Step, config.VariableReplacments)
	}
}
