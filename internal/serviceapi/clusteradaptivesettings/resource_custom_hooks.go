package clusteradaptivesettings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	adaptiveSettingsOverridesAPIField = "adaptiveSettingsOverrides"
	effectiveAdaptiveSettingsAPIField = "effectiveAdaptiveSettings"
)

var _ autogen.PreCreateAPICallHook = &rs{}
var _ autogen.PreUpdateAPICallHook = &rs{}
var _ autogen.PostReadAPICallHook = &rs{}

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err == nil {
		// Atlas omits cleared overrides; the shared decoder preserves absent and null fields.
		req.State.(*TFModel).AdaptiveSettingsOverrides = jsontypes.NewNormalizedNull()
	}
	return result
}

func (r *rs) PreCreateAPICall(ctx context.Context, callParams config.APICallParams, bodyReq []byte) (updatedParams config.APICallParams, updatedBody []byte, err error) {
	updatedBody, err = r.adaptiveSettingsPatchBody(ctx, callParams, bodyReq)
	return callParams, updatedBody, err
}

func (r *rs) PreUpdateAPICall(ctx context.Context, callParams config.APICallParams, bodyReq []byte) (updatedParams config.APICallParams, updatedBody []byte, err error) {
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
		rawOverrides = json.RawMessage("null")
	}

	var plannedOverrides map[string]json.RawMessage
	if err := json.Unmarshal(rawOverrides, &plannedOverrides); err != nil {
		return nil, fmt.Errorf("unmarshal planned %s: %w", adaptiveSettingsOverridesAPIField, err)
	}
	if plannedOverrides == nil {
		request[adaptiveSettingsOverridesAPIField] = json.RawMessage("null")
		return json.Marshal(request)
	}

	// The effective map provides the complete key set without coupling the provider to current setting names.
	effectiveSettings, err := r.getEffectiveAdaptiveSettings(ctx, callParams)
	if err != nil {
		return nil, err
	}
	for key := range effectiveSettings {
		effectiveSettings[key] = json.RawMessage("null")
	}
	maps.Copy(effectiveSettings, plannedOverrides)
	request[adaptiveSettingsOverridesAPIField], err = json.Marshal(effectiveSettings)
	if err != nil {
		return nil, err
	}
	return json.Marshal(request)
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
	var response struct {
		EffectiveSettings map[string]json.RawMessage `json:"effectiveAdaptiveSettings"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("unmarshal Adaptive Settings before PATCH: %w", err)
	}
	if response.EffectiveSettings == nil {
		return nil, fmt.Errorf("adaptive settings response contains missing or null %s", effectiveAdaptiveSettingsAPIField)
	}
	return response.EffectiveSettings, nil
}
