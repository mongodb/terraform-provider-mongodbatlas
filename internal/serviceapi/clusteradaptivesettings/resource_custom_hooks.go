package clusteradaptivesettings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	adaptiveSettingsOverridesAPIField = "adaptiveSettingsOverrides"
	effectiveAdaptiveSettingsAPIField = "effectiveAdaptiveSettings"
)

var _ autogen.PreCreateAPICallWithContextHook = &rs{}
var _ autogen.PreUpdateAPICallWithContextHook = &rs{}

func (r *rs) PreCreateAPICallWithContext(ctx context.Context, callParams config.APICallParams, bodyReq []byte) (updatedParams config.APICallParams, updatedBody []byte, err error) {
	updatedBody, err = r.adaptiveSettingsPatchBody(ctx, callParams, bodyReq)
	return callParams, updatedBody, err
}

func (r *rs) PreUpdateAPICallWithContext(ctx context.Context, callParams config.APICallParams, bodyReq []byte) (updatedParams config.APICallParams, updatedBody []byte, err error) {
	updatedBody, err = r.adaptiveSettingsPatchBody(ctx, callParams, bodyReq)
	return callParams, updatedBody, err
}

func (r *rs) adaptiveSettingsPatchBody(ctx context.Context, callParams config.APICallParams, bodyReq []byte) ([]byte, error) {
	var request map[string]json.RawMessage
	if err := json.Unmarshal(bodyReq, &request); err != nil {
		return nil, fmt.Errorf("unmarshal Adaptive Settings PATCH request: %w", err)
	}
	rawOverrides, ok := request[adaptiveSettingsOverridesAPIField]
	if !ok {
		request[adaptiveSettingsOverridesAPIField] = json.RawMessage("null")
		return marshalAdaptiveSettingsRequest(request)
	}

	var plannedOverrides map[string]json.RawMessage
	if err := json.Unmarshal(rawOverrides, &plannedOverrides); err != nil {
		return nil, fmt.Errorf("unmarshal planned %s: %w", adaptiveSettingsOverridesAPIField, err)
	}
	if plannedOverrides == nil {
		request[adaptiveSettingsOverridesAPIField] = json.RawMessage("null")
		return marshalAdaptiveSettingsRequest(request)
	}

	// The effective map provides the complete key set without coupling the provider to current setting names.
	effectiveSettings, err := r.getEffectiveAdaptiveSettings(ctx, callParams)
	if err != nil {
		return nil, err
	}
	return replaceAdaptiveSettingsOverrides(request, plannedOverrides, effectiveSettings)
}

func (r *rs) getEffectiveAdaptiveSettings(ctx context.Context, callParams config.APICallParams) (map[string]json.RawMessage, error) {
	getParams := callParams
	getParams.Method = http.MethodGet
	resp, err := r.Client.UntypedAPICall(ctx, getParams, nil)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("get Adaptive Settings before PATCH: %w", err)
	}
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("get Adaptive Settings before PATCH: empty response")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read Adaptive Settings before PATCH: %w", err)
	}
	var response map[string]json.RawMessage
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal Adaptive Settings before PATCH: %w", err)
	}
	rawEffectiveSettings, ok := response[effectiveAdaptiveSettingsAPIField]
	if !ok {
		return nil, fmt.Errorf("adaptive settings response is missing %s", effectiveAdaptiveSettingsAPIField)
	}
	var effectiveSettings map[string]json.RawMessage
	if err := json.Unmarshal(rawEffectiveSettings, &effectiveSettings); err != nil {
		return nil, fmt.Errorf("unmarshal %s before PATCH: %w", effectiveAdaptiveSettingsAPIField, err)
	}
	if effectiveSettings == nil {
		return nil, fmt.Errorf("adaptive settings response contains null %s", effectiveAdaptiveSettingsAPIField)
	}
	return effectiveSettings, nil
}

// replaceAdaptiveSettingsOverrides resets every effective setting before overlaying Terraform's desired overrides.
func replaceAdaptiveSettingsOverrides(request, plannedOverrides, effectiveSettings map[string]json.RawMessage) ([]byte, error) {
	patchOverrides := make(map[string]json.RawMessage, len(effectiveSettings))
	for key := range effectiveSettings {
		patchOverrides[key] = json.RawMessage("null")
	}
	maps.Copy(patchOverrides, plannedOverrides)

	updatedOverrides, err := json.Marshal(patchOverrides)
	if err != nil {
		return nil, fmt.Errorf("marshal %s PATCH: %w", adaptiveSettingsOverridesAPIField, err)
	}
	request[adaptiveSettingsOverridesAPIField] = updatedOverrides
	return marshalAdaptiveSettingsRequest(request)
}

func marshalAdaptiveSettingsRequest(request map[string]json.RawMessage) ([]byte, error) {
	updatedRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal Adaptive Settings PATCH request: %w", err)
	}
	return updatedRequest, nil
}
