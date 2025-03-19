package unit_test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/require"
)

const (
	advancedClusterRelPath = "internal/service/advancedcluster"
	prefixName             = "TestMockPlanChecks_"
)

func ConvertFileNameToPlanCheckDir(t *testing.T, fileName string) string {
	t.Helper()
	dir, name := filepath.Split(fileName)
	nameParts := strings.SplitN(name, "_", 2)
	require.Len(t, nameParts, 2)
	nameWithSuffix := nameParts[1]
	specficName := strings.TrimSuffix(nameWithSuffix, ".yaml")
	newName := prefixName + specficName
	return path.Join(dir, newName)
}

func CreateImportData(t *testing.T, httpMockFile, outputDir string) {
	t.Helper()
	data, err := unit.ParseTestDataConfigYAML(httpMockFile)
	require.NoError(t, err)
	if !unit.FileExist(outputDir) {
		err = os.Mkdir(outputDir, 0755)
		require.NoError(t, err)
	}
	firstStep := data.Steps[0]
	getRequests := []unit.RequestInfo{}
	for _, req := range firstStep.RequestResponses {
		if req.Method == "GET" {
			getRequests = append(getRequests, req)
		}
	}
	templateMockHTTPData := unit.MockHTTPData{
		Steps: []unit.StepRequests{
			{
				Config: firstStep.Config,
			},
		},
		Variables: data.Variables,
	}
	for _, req := range getRequests {
		lastResponse := req.Responses[len(req.Responses)-1]
		jsonFileName := strings.ReplaceAll(fmt.Sprintf("import_%s.json", req.IdShort()), "/", "_")
		jsonFilePath := path.Join(outputDir, jsonFileName)
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
	templateYaml, err := unit.ConfigYaml(&templateMockHTTPData)
	require.NoError(t, err)
	testDataDir := filepath.Dir(outputDir)
	templateYamlPath := path.Join(testDataDir, filepath.Base(outputDir)+".tmpl.yaml")
	err = os.WriteFile(templateYamlPath, []byte(templateYaml), 0644)
	require.NoError(t, err)
}

func TestConvertMockableTests(t *testing.T) {
	for relPath, usedNames := range map[string][]string{
		advancedClusterRelPath: {"TestAccMockableAdvancedCluster_removeBlocksFromConfig"},
	} {
		testDataPath := unit.RepoPath(relPath + "/testdata")
		gitIgnorePath := path.Join(testDataPath, ".gitignore")
		if !unit.FileExist(gitIgnorePath) {
			err := os.WriteFile(gitIgnorePath, []byte("TestMockPlan*.yaml\n!*.tmpl.yaml\n"), 0644)
			require.NoError(t, err)
		}
		for _, testFileName := range usedNames {
			testFilePath := path.Join(testDataPath, testFileName+".yaml")
			destDir := ConvertFileNameToPlanCheckDir(t, testFilePath)
			t.Logf("Converting %s to %s", testFilePath, destDir)
			CreateImportData(t, testFilePath, destDir)
		}
	}
}
