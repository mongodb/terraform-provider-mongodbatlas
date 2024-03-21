package replay

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const fileName = "execution-variables"

type ExecutionVariables struct {
	ProjectID string `json:"projectId"`
}

func CaptureExecutionVariables(projectID string) error {
	jsonData, err := json.MarshalIndent(ExecutionVariables{ProjectID: projectID}, "", "    ")
	if err != nil {
		return err
	}
	return createFileInSimulationDir(jsonData, fileName)
}

func ObtainExecutionVariables() (*ExecutionVariables, error) {
	filePath := fmt.Sprintf("%s%s.json", resultFilePrefix, fileName)
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
