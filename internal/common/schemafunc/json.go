package schemafunc

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func DiffSuppressJSON() planmodifier.String {
	return diffSuppressJSONModifier{}
}

type diffSuppressJSONModifier struct {
}

func (m diffSuppressJSONModifier) Description(_ context.Context) string {
	return "If the parsed jsons are the same, the value of this attribute in state will not change."
}

func (m diffSuppressJSONModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m diffSuppressJSONModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if req.PlanValue.IsUnknown() || req.ConfigValue.IsUnknown() {
		return
	}
	old, newStr := req.StateValue.String(), req.PlanValue.String()
	if EqualJSON(old, newStr, req.Path.String()) {
		resp.PlanValue = req.StateValue
	}
}

func EqualJSON(old, newStr, errContext string) bool {
	var j, j2 any

	if old == "" {
		old = "{}"
	}

	if newStr == "" {
		newStr = "{}"
	}

	if err := json.Unmarshal([]byte(old), &j); err != nil {
		log.Printf("[ERROR] cannot unmarshal old %s json %v", errContext, err)
		return false
	}
	if err := json.Unmarshal([]byte(newStr), &j2); err != nil {
		log.Printf("[ERROR] cannot unmarshal new %s json %v", errContext, err)
		return false
	}
	return reflect.DeepEqual(&j, &j2)
}
