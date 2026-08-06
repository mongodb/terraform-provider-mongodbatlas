package clusteradaptivesettings

import (
	"bytes"
	"encoding/json"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const adaptiveSettingsOverridesAPIField = "adaptiveSettingsOverrides"

var adaptiveSettingKeys = []string{
	"OVERLOAD_PROTECTION",
	"SEARCH_OVERLOAD_PROTECTION",
}

var _ autogen.PreCreateAPICallHook = &rs{}
var _ autogen.PreUpdateAPICallHook = &rs{}

func (r *rs) PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (updatedParams config.APICallParams, updatedBody []byte) {
	return callParams, replaceAdaptiveSettingsOverrides(bodyReq)
}

func (r *rs) PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (updatedParams config.APICallParams, updatedBody []byte) {
	return callParams, replaceAdaptiveSettingsOverrides(bodyReq)
}

// replaceAdaptiveSettingsOverrides converts Terraform's desired object into an RFC 7396 patch by resetting omitted known keys.
func replaceAdaptiveSettingsOverrides(bodyReq []byte) []byte {
	var request map[string]json.RawMessage
	if err := json.Unmarshal(bodyReq, &request); err != nil {
		return bodyReq
	}

	rawOverrides, ok := request[adaptiveSettingsOverridesAPIField]
	if !ok || bytes.Equal(bytes.TrimSpace(rawOverrides), []byte("null")) {
		request[adaptiveSettingsOverridesAPIField] = json.RawMessage("null")
		return marshalRequestOrOriginal(request, bodyReq)
	}

	var overrides map[string]json.RawMessage
	if err := json.Unmarshal(rawOverrides, &overrides); err != nil {
		return bodyReq
	}
	for _, key := range adaptiveSettingKeys {
		if _, exists := overrides[key]; !exists {
			overrides[key] = json.RawMessage("null")
		}
	}

	updatedOverrides, err := json.Marshal(overrides)
	if err != nil {
		return bodyReq
	}
	request[adaptiveSettingsOverridesAPIField] = updatedOverrides
	return marshalRequestOrOriginal(request, bodyReq)
}

func marshalRequestOrOriginal(request map[string]json.RawMessage, original []byte) []byte {
	updatedRequest, err := json.Marshal(request)
	if err != nil {
		return original
	}
	return updatedRequest
}
