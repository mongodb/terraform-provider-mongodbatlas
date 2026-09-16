package serviceaccount

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var (
	_ autogen.ResourceSchemaHook   = (*rs)(nil)
	_ autogen.PreCreateAPICallHook = (*rs)(nil)
)

// ResourceSchema drops the create-required guardrail on `secret_expires_after_hours`.
// Omitting the attribute is the signal for "create the Service Account without a bootstrap secret",
// so it must be valid on create. `CreateOnly` stays: changing the hours after creation still fails.
func (r *rs) ResourceSchema(ctx context.Context, s schema.Schema) schema.Schema {
	attr, ok := s.Attributes["secret_expires_after_hours"].(schema.Int64Attribute)
	if !ok {
		return s
	}
	modifiers := make([]planmodifier.Int64, 0, len(attr.PlanModifiers))
	for _, modifier := range attr.PlanModifiers {
		if isRequestOnlyRequiredOnCreate(modifier) {
			continue
		}
		modifiers = append(modifiers, modifier)
	}
	attr.PlanModifiers = modifiers
	s.Attributes["secret_expires_after_hours"] = attr
	return s
}

// isRequestOnlyRequiredOnCreate compares the dynamic type of a plan modifier. An interface assertion
// cannot be used: customplanmodifier.RequestOnlyRequiredOnCreateModifier has the same method set as
// CreateOnlyModifier, so every create-only modifier would match it too.
func isRequestOnlyRequiredOnCreate(modifier planmodifier.Int64) bool {
	return reflect.TypeOf(modifier) == reflect.TypeOf(customplanmodifier.RequestOnlyRequiredOnCreate())
}

// PreCreateAPICall adds `withoutInitialSecret: true` when `secretExpiresAfterHours` is absent.
// The generated schema hides `without_initial_secret` (see tools/codegen/config.yml), so this is the
// only place where the API learns the account should be created without an initial secret.
// A decode or marshal error returns the original body so the API fails on it instead of silently
// creating an account with a bootstrap secret.
func (r *rs) PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	var body map[string]any
	if err := autogen.Decode(bodyReq, &body); err != nil {
		return callParams, bodyReq
	}
	if _, hasHours := body["secretExpiresAfterHours"]; hasHours {
		return callParams, bodyReq
	}
	body["withoutInitialSecret"] = true
	updated, err := json.Marshal(body)
	if err != nil {
		return callParams, bodyReq
	}
	return callParams, updated
}
