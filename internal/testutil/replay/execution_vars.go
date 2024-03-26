package replay

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

const envVarFilePostfix = "-execution-vars"

type ExecutionVariables struct {
	ProjectID string `json:"projectId"`
}

// ManageProjectID is a function that will store the generated project id during capture mode,
// and in case of simulate mode will simple fetch the project id defined in execution variables file.
func ManageProjectID(t *testing.T, projectIDProvider func(t testing.TB) string) string {
	t.Helper()
	if IsInCaptureMode() {
		id := projectIDProvider(t)
		serializeValueToEnvFile(t, id)
		return id
	}

	if IsInSimulateMode() {
		vars, err := ObtainExecutionVariables(t)
		if err != nil {
			log.Fatal("failed to obtain env file during simulation mode")
		}
		return vars.ProjectID // returns project stored in execution vars file
	}

	// case where no replay mode is configured
	return projectIDProvider(t)
}

func serializeValueToEnvFile(t *testing.T, projectID string) {
	t.Helper()
	jsonData, err := json.MarshalIndent(ExecutionVariables{ProjectID: projectID}, "", "    ")
	if err != nil {
		log.Fatal("failed to serialize json with env variables")
	}
	err = createFileInSimulationDir(jsonData, fmt.Sprintf("%s%s.json", t.Name(), envVarFilePostfix))
	if err != nil {
		log.Fatal("failed to write json with env variables")
	}
}

func ObtainExecutionVariables(t *testing.T) (*ExecutionVariables, error) {
	t.Helper()
	filePath := filePathInSimulationDir(fmt.Sprintf("%s%s.json", t.Name(), envVarFilePostfix))
	pairsFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("got error when opening execution variables file %s", err.Error())
	}
	var vars ExecutionVariables
	body, err := io.ReadAll(pairsFile)
	if err != nil {
		return nil, fmt.Errorf("got error when opening execution variables file %s", err.Error())
	}
	err = json.Unmarshal(body, &vars)
	if err != nil {
		return nil, fmt.Errorf("got error while parsing execution variables %s", err.Error())
	}
	return &vars, nil
}
