package advancedclustertpf

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20241023001/admin"
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
)

func ReadResponse(number int) (*admin.ClusterDescription20240805, error) {
	responseJSON, ok := responsesCreate[number]
	if !ok {
		return nil, fmt.Errorf("unknown response number %d", number)
	}
	return parseReadResponse(responseJSON)
}

func parseReadResponse(data string) (*admin.ClusterDescription20240805, error) {
	var SDKModel admin.ClusterDescription20240805
	err := json.Unmarshal([]byte(data), &SDKModel)
	return &SDKModel, err
}
