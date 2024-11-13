package advancedclustertpf

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20241023002/admin"
)

var (
	//go:embed testdata/create_1.json
	create1 string
	//go:embed testdata/create_2.json
	create2 string
	//go:embed testdata/create_3.json
	create3         string
	responsesCreate = map[int]string{
		1: create1,
		2: create2,
		3: create3,
	}
	//go:embed testdata/process_args_1.json
	processArgs1         string
	responsesProcessArgs = map[int]string{
		1: processArgs1,
	}
	currentClusterResponse     = 1
	currentProcessArgsResponse = 1
	actualCreate               string
	actualUpdate               string
)

func SetCurrentClusterResponse(responseNumber int) {
	currentClusterResponse = responseNumber
}

func ReadClusterResponse() (*admin.ClusterDescription20240805, error) {
	responseJSON, ok := responsesCreate[currentClusterResponse]
	if !ok {
		return nil, fmt.Errorf("unknown cluster response number %d", currentClusterResponse)
	}
	var SDKModel admin.ClusterDescription20240805
	err := json.Unmarshal([]byte(responseJSON), &SDKModel)
	return &SDKModel, err
}

func ReadClusterProcessArgsResponse() (*admin.ClusterDescriptionProcessArgs20240805, error) {
	responseJSON, ok := responsesProcessArgs[currentProcessArgsResponse]
	if !ok {
		return nil, fmt.Errorf("unknown process args response number %d", currentProcessArgsResponse)
	}
	var SDKModel admin.ClusterDescriptionProcessArgs20240805
	err := json.Unmarshal([]byte(responseJSON), &SDKModel)
	return &SDKModel, err
}

func StoreCreatePayload(payload *admin.ClusterDescription20240805) error {
	localPayload, err := dumpJSON(payload)
	if err != nil {
		return err
	}
	actualCreate = localPayload
	return nil
}

func dumpJSON(payload any) (string, error) {
	jsonPayload := strings.Builder{}
	encoder := json.NewEncoder(&jsonPayload)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(payload)
	if err != nil {
		return "", err
	}
	return jsonPayload.String(), nil
}

func StoreUpdatePayload(payload *admin.ClusterDescription20240805) error {
	localPayload, err := dumpJSON(payload)
	if err != nil {
		return err
	}
	actualUpdate = localPayload
	return nil
}

func ReadLastCreatePayload() (string, error) {
	if actualCreate == "" {
		return "", fmt.Errorf("no create payload has been stored")
	}
	return actualCreate, nil
}

func ReadLastUpdatePayload() (string, error) {
	if actualUpdate == "" {
		return "", fmt.Errorf("no update payload has been stored")
	}
	return actualUpdate, nil
}
