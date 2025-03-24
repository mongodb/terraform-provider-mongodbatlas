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
	pkgAdvancedCluster    = "advancedcluster"
	pkgAdvancedClusterTPF = "advancedclustertpf"
	pkgRelPath            = "internal/service"
)

var (
	clusterVariableReplacements = map[string]string{
		"clusterName": unit.MockedClusterName,
		"groupId":     unit.MockedProjectID,
	}
)

type importNameConfig struct {
	VariableReplacments map[string]string
	TestName            string
	SrcPackage          string
	DestPackage         string
	Step                int
}

// Manual test meant for creating the data needed for a MockPlanChecks Test.
func TestConvertMockableTests(t *testing.T) {
	if os.Getenv("CONVERT_MOCKABLE_TESTS") == "" {
		t.Skip("CONVERT_MOCKABLE_TESTS is not set, avoid running this in CI and by accident")
	}
	for importName, config := range map[string]importNameConfig{
		unit.ImportNameClusterTwoRepSpecsWithAutoScalingAndSpecs: {
			TestName:            "TestAccMockableAdvancedCluster_removeBlocksFromConfig",
			Step:                1,
			VariableReplacments: clusterVariableReplacements,
			SrcPackage:          pkgAdvancedCluster,
			DestPackage:         pkgAdvancedClusterTPF,
		},
		unit.ImportNameClusterReplicasetOneRegion: {
			TestName:            "TestAccMockableAdvancedCluster_replicasetAdvConfigUpdate",
			Step:                1,
			VariableReplacments: clusterVariableReplacements,
			SrcPackage:          pkgAdvancedCluster,
			DestPackage:         pkgAdvancedClusterTPF,
		},
	} {
		srcTestdata := unit.RepoPath(path.Join(pkgRelPath, config.SrcPackage, "testdata"))
		destTestdata := unit.RepoPath(path.Join(pkgRelPath, config.DestPackage, "testdata"))
		ensureDir(t, destTestdata)
		srcTestdataPath := path.Join(srcTestdata, config.TestName+".yaml")
		destTestdataPath := path.Join(destTestdata, importName+".tmpl.yaml")
		t.Logf("Converting %s step %d to %s", srcTestdataPath, config.Step, destTestdataPath)
		createImportData(t, srcTestdataPath, destTestdataPath, importName, config.Step, config.VariableReplacments)
	}
}

func createImportData(t *testing.T, srcMockFile, destMockFile, importName string, stepNr int, newVars map[string]string) {
	t.Helper()
	destTestData, _ := path.Split(destMockFile)
	require.True(t, strings.HasSuffix(destTestData, "/testdata/"))
	ensureDir(t, destTestData)
	destOutputDir := ensureDir(t, path.Join(destTestData, importName))

	templateMockHTTPData := createImportMockData(t, srcMockFile, destOutputDir, stepNr, newVars)
	require.NoError(t, templateMockHTTPData.UpdateVariablesIgnoreChanges(t, newVars))

	templateYaml, err := unit.ConfigYaml(templateMockHTTPData)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(destMockFile, []byte(templateYaml), 0o600))
}

func createImportMockData(t *testing.T, srcMockFile, destOutputDir string, stepNr int, newVars map[string]string) *unit.MockHTTPData {
	t.Helper()
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
		jsonFileName := strings.ReplaceAll(fmt.Sprintf("import_%s.json", req.IDShort()), "/", "_")
		jsonFilePath := path.Join(destOutputDir, jsonFileName)
		err = os.WriteFile(jsonFilePath, []byte(lastResponse.Text), 0o600)
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

func ensureDir(t *testing.T, dir string) string {
	t.Helper()
	if !unit.FileExist(dir) {
		err := os.Mkdir(dir, 0o755)
		require.NoError(t, err)
	}
	return dir
}
