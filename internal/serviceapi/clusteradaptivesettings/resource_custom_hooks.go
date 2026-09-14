package clusteradaptivesettings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

var (
	_ autogen.PostReadAPICallHook = &rs{}
	_ autogen.ResourceSchemaHook  = &rs{}
)

// PostReadAPICall normalizes the response when Atlas omits cleared overrides.
// After a whole-map reset (null), Atlas omits adaptiveSettingsOverrides from
// the response; the shared decoder preserves absent and null fields.
func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	if result.Err == nil {
		req.State.(*TFModel).AdaptiveSettingsOverrides = jsontypes.NewNormalizedNull()
	}
	return result
}

// ResourceSchema adds a validator to adaptive_settings_overrides rejecting
// null-valued entries. Under replacement PATCH semantics, removing a key from
// the configured object resets that override; a null entry is invalid.
func (r *rs) ResourceSchema(_ context.Context, s schema.Schema) schema.Schema {
	overrides := s.Attributes["adaptive_settings_overrides"].(schema.StringAttribute)
	overrides.Validators = append(overrides.Validators, adaptiveSettingsOverridesValidator{})
	s.Attributes["adaptive_settings_overrides"] = overrides
	return s
}

type adaptiveSettingsOverridesValidator struct{}

func (v adaptiveSettingsOverridesValidator) Description(context.Context) string {
	return "A JSON object without null-valued entries. Remove a key to reset its override."
}

func (v adaptiveSettingsOverridesValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v adaptiveSettingsOverridesValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var settings map[string]json.RawMessage
	if err := json.Unmarshal([]byte(req.ConfigValue.ValueString()), &settings); err != nil || settings == nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid Adaptive Settings overrides", "Use jsonencode with an object. To reset all overrides, omit adaptive_settings_overrides or configure an empty object.")
		return
	}
	for key, value := range settings {
		if bytes.Equal(bytes.TrimSpace(value), []byte("null")) {
			resp.Diagnostics.AddAttributeError(req.Path, "Invalid Adaptive Settings override", fmt.Sprintf("Override %q is null. To reset an override, remove its key from the configured object.", key))
		}
	}
}
