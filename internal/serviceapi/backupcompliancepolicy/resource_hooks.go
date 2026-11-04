package backupcompliancepolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// Manual-parity custom code for mongodbatlas_backup_compliance_policy (lift-and-shift from the handwritten
// SDKv2 resource).

type TFExpandedModel struct {
	// Crafted id of the manual resource, kept so migrated state and references to .id keep working.
	// Injected into the response body by the hooks below under a provider-only key.
	ID types.String `tfsdk:"id" apiname:"tfStateId" autogen:"omitjson"`
}

const (
	stateIDKey = "tfStateId"

	// The manual resource always sends overwriteBackupPolicies=false: the API default (true) would silently
	// overwrite the non-complying backup policies of every cluster in the project.
	overwriteBackupPoliciesParam = "overwriteBackupPolicies"
)

// Create and update responses may not include the project id, it is taken from the request path.
var projectIDPathRegex = regexp.MustCompile(`/groups/([0-9a-f]{24})/backupCompliancePolicy`)

var (
	_ autogen.PreCreateAPICallHook  = (*rs)(nil)
	_ autogen.PreUpdateAPICallHook  = (*rs)(nil)
	_ autogen.PostCreateAPICallHook = (*rs)(nil)
	_ autogen.PostReadAPICallHook   = (*rs)(nil)
	_ autogen.PostUpdateAPICallHook = (*rs)(nil)
	_ autogen.ResourceSchemaHook    = (*rs)(nil)
	// Enables the state upgrade below for state written by the manual SDKv2 resource (schema version 0).
	_ resource.ResourceWithUpgradeState = (*rs)(nil)
)

func (r *rs) ResourceSchema(ctx context.Context, baseSchema schema.Schema) schema.Schema {
	// Version 1: triggers the SDKv2 state upgrade (schema version 0) below.
	baseSchema.Version = 1

	// Crafted id of the manual resource, see TFExpandedModel.
	baseSchema.Attributes["id"] = schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}

	// RequiresReplace instead of the generated CreateOnly: the generated modifier hard-errors at plan time,
	// while replacing (as a project-scoped singleton) is the closest match to the manual behavior.
	projectID := baseSchema.Attributes["project_id"].(schema.StringAttribute)
	projectID.PlanModifiers = []planmodifier.String{stringplanmodifier.RequiresReplace()}
	baseSchema.Attributes["project_id"] = projectID
	return baseSchema
}

func (r *rs) PreCreateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	return withOverwriteBackupPoliciesFalse(callParams), bodyReq
}

func (r *rs) PreUpdateAPICall(callParams config.APICallParams, bodyReq []byte) (modifiedParams config.APICallParams, modifiedBody []byte) {
	return withOverwriteBackupPoliciesFalse(callParams), bodyReq
}

func (r *rs) PostCreateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return adaptResponse(result, projectIDFromRequest(result))
}

func (r *rs) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	projectID := ""
	if model, ok := req.State.(*TFModel); ok {
		projectID = model.ProjectId.ValueString()
	}
	return adaptResponse(result, projectID)
}

func (r *rs) PostUpdateAPICall(result autogen.APICallResult) autogen.APICallResult {
	return adaptResponse(result, projectIDFromRequest(result))
}

func withOverwriteBackupPoliciesFalse(callParams config.APICallParams) config.APICallParams {
	if callParams.QueryParams == nil {
		callParams.QueryParams = map[string]string{}
	}
	callParams.QueryParams[overwriteBackupPoliciesParam] = "false"
	return callParams
}

func projectIDFromRequest(result autogen.APICallResult) string {
	if result.Resp == nil || result.Resp.Request == nil || result.Resp.Request.URL == nil {
		return ""
	}
	if m := projectIDPathRegex.FindStringSubmatch(result.Resp.Request.URL.Path); m != nil {
		return m[1]
	}
	return ""
}

// adaptResponse injects the crafted id of the manual resource into the response body.
func adaptResponse(result autogen.APICallResult, projectID string) autogen.APICallResult {
	if result.Err != nil || autogen.IsEmptyJSON(result.Body) || projectID == "" {
		return result
	}
	var obj map[string]any
	if err := autogen.Decode(result.Body, &obj); err != nil {
		result.Err = err
		return result
	}
	obj[stateIDKey] = conversion.EncodeStateID(map[string]string{
		"project_id": projectID,
	})
	body, err := json.Marshal(obj)
	if err != nil {
		result.Err = fmt.Errorf("marshaling backup compliance policy response: %w", err)
		return result
	}
	result.Body = body
	return result
}

// UpgradeState upgrades state written by the manual SDKv2 resource (schema version 0) to the autogenerated
// schema (version 1). The structural change (BL1): the manual `policy_item_hourly`/`daily`/`weekly`/`monthly`/
// `yearly` blocks and the `on_demand_policy_item` block (SDKv2 lists, several with MaxItems: 1) become the
// API's `scheduled_policy_items` list and a single `on_demand_policy_item` object.
func (r *rs) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	return map[int64]resource.StateUpgrader{
		0: {StateUpgrader: upgradeSDKv2StateV0},
	}
}

func upgradeSDKv2StateV0(ctx context.Context, req resource.UpgradeStateRequest, resp *resource.UpgradeStateResponse) {
	rawMap := map[string]any{}
	if err := json.Unmarshal(req.RawState.JSON, &rawMap); err != nil {
		resp.Diagnostics.AddError("Unable to unmarshal previous state", err.Error())
		return
	}
	upgradeOnDemandPolicyItem(rawMap)
	upgradeScheduledPolicyItems(rawMap)

	fixedJSON, err := json.Marshal(rawMap)
	if err != nil {
		resp.Diagnostics.AddError("Unable to marshal upgraded state", err.Error())
		return
	}
	rawStateValue, err := tftypes.ValueFromJSONWithOpts(fixedJSON, resp.State.Schema.Type().TerraformType(ctx), tftypes.ValueFromJSONOpts{
		IgnoreUndefinedAttributes: true, // drops the removed policy_item_* attributes
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Previously Saved State for UpgradeResourceState", err.Error())
		return
	}
	resp.State.Raw = rawStateValue
}

// upgradeOnDemandPolicyItem converts the SDKv2 list (empty or of one item) into the object the autogenerated
// schema stores: an empty or absent list is removed (null object), the item of a single-element list becomes
// the object.
func upgradeOnDemandPolicyItem(rawMap map[string]any) {
	items, ok := rawMap["on_demand_policy_item"].([]any)
	if !ok || len(items) == 0 {
		delete(rawMap, "on_demand_policy_item")
		return
	}
	if item, isMap := items[0].(map[string]any); isMap {
		if item["frequency_type"] == "" || item["frequency_type"] == nil {
			item["frequency_type"] = "ondemand"
		}
		rawMap["on_demand_policy_item"] = item
	}
}

// upgradeScheduledPolicyItems merges the manual per-frequency policy_item_* lists into the API's
// scheduled_policy_items list, taking each item's frequency_type from the block name.
func upgradeScheduledPolicyItems(rawMap map[string]any) {
	scheduled := []any{}
	for _, frequency := range []string{"hourly", "daily", "weekly", "monthly", "yearly"} {
		items, ok := rawMap["policy_item_"+frequency].([]any)
		if !ok {
			continue
		}
		for _, item := range items {
			itemMap, isMap := item.(map[string]any)
			if !isMap {
				continue
			}
			if itemMap["frequency_type"] == "" || itemMap["frequency_type"] == nil {
				itemMap["frequency_type"] = frequency
			}
			scheduled = append(scheduled, itemMap)
		}
	}
	delete(rawMap, "policy_item_hourly")
	delete(rawMap, "policy_item_daily")
	delete(rawMap, "policy_item_weekly")
	delete(rawMap, "policy_item_monthly")
	delete(rawMap, "policy_item_yearly")
	if len(scheduled) > 0 {
		rawMap["scheduled_policy_items"] = scheduled
	}
}
