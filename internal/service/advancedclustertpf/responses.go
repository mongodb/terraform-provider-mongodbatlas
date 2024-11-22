package advancedclustertpf

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
)

var (
	//go:embed testdata/replicaset_create_resp1.json
	createReplicasetResp1 string
	//go:embed testdata/replicaset_create_resp1_final.json
	createReplicasetResp1Final string
	//go:embed testdata/replicaset_update1_resp.json
	updateReplicasetResp1 string

	//go:embed testdata/replicaset_update2_resp.json
	updateReplicasetResp2 string

	//go:embed testdata/sharded_create_resp1.json
	createSharded1 string

	//go:embed testdata/sharded_create_resp1_final.json
	createSharded2 string

	//go:embed testdata/sharded_update_resp1.json
	updateSharded1 string

	//go:embed testdata/sharded_update_resp1_final.json
	updateSharded2 string

	//go:embed testdata/process_args_1.json
	processArgs1 string

	//go:embed testdata/process_args_2.json
	processArgs2 string

	//go:embed testdata/process_args_1_legacy.json
	processArgsLegacy1 string

	//go:embed testdata/process_args_2_legacy.json
	processArgsLegacy2 string

	responsesCreate = map[string][]string{
		"replicaset": {createReplicasetResp1, createReplicasetResp1Final},
		"sharded":    {createSharded1, createSharded2},
	}
	responsesUpdate = map[string][]string{
		"replicaset": {updateReplicasetResp1, updateReplicasetResp2},
		"sharded":    {updateSharded1, updateSharded2},
	}

	responsesProcessArgs = map[int]string{
		0: processArgs1,
		1: processArgs2,
	}

	responsesProcessArgsLegacy = map[int]string{
		0: processArgsLegacy1,
		1: processArgsLegacy2,
	}
)

type MockData struct {
	ClusterResponse        string
	ResponseIndex          int
	IsUpdate               bool
	ProcessArgsIndex       int
	ProcessArgsIndexLegacy int
}

func (m *MockData) NextResponse(isUpdate, isProcessResponse bool) error {
	switch {
	case isProcessResponse:
		m.ProcessArgsIndex++
		m.ProcessArgsIndexLegacy++
	case isUpdate && !m.IsUpdate:
		m.IsUpdate = true
		m.ResponseIndex = 0
	default:
		m.ResponseIndex++
	}
	mockCallData = MockCallData{} // reset call data
	return validateResponsesResetRequests()
}

func (m *MockData) GetResponse() string {
	var responses map[string][]string
	if m.IsUpdate {
		responses = responsesUpdate
	} else {
		responses = responsesCreate
	}
	responseJSON, ok := responses[m.ClusterResponse]
	if !ok {
		return ""
	}
	return responseJSON[m.ResponseIndex]
}
func (m *MockData) GetProcessArgsResponse() string {
	responseJSON, ok := responsesProcessArgs[m.ProcessArgsIndex]
	if !ok {
		return ""
	}
	return responseJSON
}

func (m *MockData) GetProcessArgsResponseLegacy() string {
	responseJSON, ok := responsesProcessArgsLegacy[m.ProcessArgsIndexLegacy]
	if !ok {
		return ""
	}
	return responseJSON
}

var mockData = &MockData{
	ClusterResponse: "replicaset",
}

type MockCallData struct {
	ReqCreate                  string
	ReqUpdate                  string
	ReqProcessArgs             string
	ReqProcessArgsUpdate       string
	ReqProcessArgsUpdateLegacy string
}

var mockCallData = MockCallData{}

func SetMockDataResetResponses(data *MockData) error {
	mockData = data
	err := validateResponsesResetRequests()
	if err != nil {
		return err
	}
	return err
}

func validateResponsesResetRequests() error {
	_, err := ReadClusterResponse()
	if err != nil {
		return err
	}
	_, err = ReadClusterProcessArgsResponse()
	mockCallData = MockCallData{} // Reset the call data
	return err
}

func ReadClusterResponse() (*admin.ClusterDescription20240805, error) {
	response := mockData.GetResponse()
	if response == "" {
		return nil, fmt.Errorf("unknown cluster response for %s[%d]", mockData.ClusterResponse, mockData.ResponseIndex)
	}
	var SDKModel admin.ClusterDescription20240805
	err := json.Unmarshal([]byte(response), &SDKModel)
	return &SDKModel, err
}

func ReadClusterProcessArgsResponse() (*admin.ClusterDescriptionProcessArgs20240805, error) {
	response := mockData.GetProcessArgsResponse()
	if response == "" {
		return nil, fmt.Errorf("unknown process args response number %d", mockData.ProcessArgsIndex)
	}
	var SDKModel admin.ClusterDescriptionProcessArgs20240805
	err := json.Unmarshal([]byte(response), &SDKModel)
	return &SDKModel, err
}

func ReadClusterProcessArgsResponseLegacy() (*admin20240530.ClusterDescriptionProcessArgs, error) {
	response := mockData.GetProcessArgsResponseLegacy()
	if response == "" {
		return nil, fmt.Errorf("unknown process args response number %d", mockData.ProcessArgsIndex)
	}
	var SDKModel admin20240530.ClusterDescriptionProcessArgs
	err := json.Unmarshal([]byte(response), &SDKModel)
	return &SDKModel, err
}

func StoreCreatePayload(payload *admin.ClusterDescription20240805) error {
	localPayload, err := dumpJSON(payload)
	if err != nil {
		return err
	}
	mockCallData.ReqCreate = localPayload
	return nil
}

func StoreUpdatePayload(payload *admin.ClusterDescription20240805) error {
	localPayload, err := dumpJSON(payload)
	if err != nil {
		return err
	}
	mockCallData.ReqUpdate = localPayload
	return nil
}

func StoreUpdatePayloadProcessArgs(payload *admin.ClusterDescriptionProcessArgs20240805) error {
	localPayload, err := dumpJSON(payload)
	if err != nil {
		return err
	}
	mockCallData.ReqProcessArgsUpdate = localPayload
	return nil
}

func StoreUpdatePayloadProcessArgsLegacy(payload *admin20240530.ClusterDescriptionProcessArgs) error {
	localPayload, err := dumpJSON(payload)
	if err != nil {
		return err
	}
	mockCallData.ReqProcessArgsUpdateLegacy = localPayload
	return nil
}

func ReadLastCreatePayload() (string, error) {
	if mockCallData.ReqCreate == "" {
		return "", fmt.Errorf("no create payload has been stored")
	}
	return mockCallData.ReqCreate, nil
}

func ReadLastUpdatePayload() (string, error) {
	if mockCallData.ReqUpdate == "" {
		return "", fmt.Errorf("no update payload has been stored")
	}
	return mockCallData.ReqUpdate, nil
}

func ReadLastUpdatePayloadProcessArgs() (string, error) {
	if mockCallData.ReqProcessArgsUpdate == "" {
		return "", fmt.Errorf("no update payload has been stored")
	}
	return mockCallData.ReqProcessArgsUpdate, nil
}

func ReadLastUpdatePayloadProcessArgsLegacy() (string, error) {
	if mockCallData.ReqProcessArgsUpdate == "" {
		return "", fmt.Errorf("no update payload has been stored")
	}
	return mockCallData.ReqProcessArgsUpdateLegacy, nil
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
