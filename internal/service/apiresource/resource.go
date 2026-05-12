package apiresource

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicjson"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicreshape"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const ResourceName = "api_resource"

var (
	_ resource.ResourceWithConfigure      = &rs{}
	_ resource.ResourceWithImportState    = &rs{}
	_ resource.ResourceWithModifyPlan     = &rs{}
	_ resource.ResourceWithValidateConfig = &rs{}
)

func Resource() resource.Resource {
	return &rs{
		RSCommon: config.RSCommon{ResourceName: ResourceName},
	}
}

type rs struct {
	config.RSCommon
}

func (r *rs) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ResourceSchema(ctx)
}

func (r *rs) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !cfg.VersionHeader.IsNull() && !cfg.VersionHeader.IsUnknown() &&
		!cfg.Preview.IsNull() && !cfg.Preview.IsUnknown() && cfg.Preview.ValueBool() {
		resp.Diagnostics.AddAttributeError(path.Root("preview"),
			"version_header and preview are mutually exclusive",
			"Set either `version_header` or `preview = true`, not both.")
	}
}

func (r *rs) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip on destroy (Plan is null).
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if plan.Path.IsNull() || plan.Path.IsUnknown() {
		return
	}
	emitTypedCounterpartWarning(ctx, plan.Path.ValueString(), plan.Preview.ValueBool(), &resp.Diagnostics)
}

func (r *rs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bodyMap, sensitiveMap, diags := buildRequestMaps(plan.Body, plan.SensitiveBody)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createURL := plan.Path.ValueString()
	versionHeader := resolveVersionHeader(plan.VersionHeader, plan.Preview)

	mergedBody := mergeMaps(bodyMap, sensitiveMap)
	bodyBytes, err := json.Marshal(mergedBody)
	if err != nil {
		resp.Diagnostics.AddError("encoding request body", err.Error())
		return
	}

	result := callAPI(ctx, r.Client, plan.CreateMethod.ValueString(), createURL, versionHeader, bodyBytes)
	if result.Err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API %s %s failed", plan.CreateMethod.ValueString(), createURL), responseError(result))
		return
	}

	state := plan
	state.VersionHeader = types.StringValue(versionHeader)
	resp.Diagnostics.Append(populateAfterWrite(ctx, &state, bodyMap, sensitiveMap, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bodyMap, _, diags := buildRequestMaps(state.Body, types.DynamicNull())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	prevOutput := dynamicToMap(state.Output)
	readURL, d := deriveReadURL(state.Path.ValueString(), state.IDAttribute, prevOutput)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	versionHeader := resolveVersionHeader(state.VersionHeader, state.Preview)

	result := callAPI(ctx, r.Client, defaultReadMethod, readURL, versionHeader, nil)
	if result.NotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	if result.Err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API GET %s failed", readURL), responseError(result))
		return
	}

	resp.Diagnostics.Append(populateAfterRead(ctx, &state, bodyMap, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *rs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state TFModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bodyMap, sensitiveMap, diags := buildRequestMaps(plan.Body, plan.SensitiveBody)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	prevOutput := dynamicToMap(state.Output)
	updateURL, d := deriveReadURL(plan.Path.ValueString(), plan.IDAttribute, prevOutput)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	versionHeader := resolveVersionHeader(plan.VersionHeader, plan.Preview)

	mergedBody := mergeMaps(bodyMap, sensitiveMap)
	for _, k := range createOnlyKeys(plan.CreateOnlyBodyKeys) {
		delete(mergedBody, k)
	}
	bodyBytes, err := json.Marshal(mergedBody)
	if err != nil {
		resp.Diagnostics.AddError("encoding request body", err.Error())
		return
	}

	result := callAPI(ctx, r.Client, plan.UpdateMethod.ValueString(), updateURL, versionHeader, bodyBytes)
	if result.Err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API %s %s failed", plan.UpdateMethod.ValueString(), updateURL), responseError(result))
		return
	}

	newState := plan
	newState.VersionHeader = types.StringValue(versionHeader)
	resp.Diagnostics.Append(populateAfterWrite(ctx, &newState, bodyMap, sensitiveMap, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *rs) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TFModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	prevOutput := dynamicToMap(state.Output)
	deleteURL, d := deriveReadURL(state.Path.ValueString(), state.IDAttribute, prevOutput)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	versionHeader := resolveVersionHeader(state.VersionHeader, state.Preview)

	result := callAPI(ctx, r.Client, defaultDeleteMethod, deleteURL, versionHeader, nil)
	if result.NotFound {
		return
	}
	if result.Err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API DELETE %s failed", deleteURL), responseError(result))
	}
}

func (r *rs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.AddWarning(
		"Import does not recover body / sensitive_body",
		"After import, re-declare body and sensitive_body in HCL. The next plan will surface drift until the configured body matches what Atlas returns.",
	)
}

// deriveReadURL appends each id_attribute value (looked up in output) to path.
// Returns path unchanged for singletons (no id_attribute configured).
func deriveReadURL(basePath string, idAttr types.List, output map[string]any) (string, diag.Diagnostics) {
	var diags diag.Diagnostics
	if idAttr.IsNull() || idAttr.IsUnknown() {
		return basePath, diags
	}
	elems := idAttr.Elements()
	if len(elems) == 0 {
		return basePath, diags
	}
	var b strings.Builder
	b.WriteString(strings.TrimRight(basePath, "/"))
	for _, e := range elems {
		sv, ok := e.(types.String)
		if !ok || sv.IsNull() || sv.IsUnknown() {
			diags.AddAttributeError(path.Root("id_attribute"),
				"invalid id_attribute entry", "elements must be non-null strings")
			return "", diags
		}
		key := sv.ValueString()
		raw, present := output[key]
		if !present {
			diags.AddAttributeError(path.Root("id_attribute"),
				"id_attribute not found in output",
				fmt.Sprintf("key %q not present in API response. Available output keys: %v", key, outputKeys(output)))
			return "", diags
		}
		s, err := stringifyScalar(raw)
		if err != nil {
			diags.AddAttributeError(path.Root("id_attribute"),
				"id_attribute is not a scalar",
				fmt.Sprintf("output[%q] = %v: %s", key, raw, err))
			return "", diags
		}
		b.WriteByte('/')
		b.WriteString(s)
	}
	return b.String(), diags
}

func outputKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func stringifyScalar(v any) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case bool:
		if t {
			return "true", nil
		}
		return "false", nil
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", t), "0"), "."), nil
	case int, int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%d", t), nil
	}
	if v == nil {
		return "", fmt.Errorf("null value")
	}
	return "", fmt.Errorf("unsupported type %T", v)
}

func resolveVersionHeader(versionHeader types.String, preview types.Bool) string {
	if !preview.IsNull() && !preview.IsUnknown() && preview.ValueBool() {
		return previewVersionHeader
	}
	if !versionHeader.IsNull() && !versionHeader.IsUnknown() && versionHeader.ValueString() != "" {
		return versionHeader.ValueString()
	}
	return todayVersionHeader()
}

// buildRequestMaps unmarshals body and sensitive_body Dynamics into Go maps.
// Nil/Null Dynamics return empty maps.
func buildRequestMaps(body, sensitiveBody types.Dynamic) (bodyMap, sensitiveMap map[string]any, diags diag.Diagnostics) {
	var err error
	bodyMap, err = dynamicToMapStrict(body)
	if err != nil {
		diags.AddAttributeError(path.Root("body"), "invalid body", err.Error())
		return nil, nil, diags
	}
	sensitiveMap, err = dynamicToMapStrict(sensitiveBody)
	if err != nil {
		diags.AddAttributeError(path.Root("sensitive_body"), "invalid sensitive_body", err.Error())
		return nil, nil, diags
	}
	return bodyMap, sensitiveMap, diags
}

func dynamicToMapStrict(d types.Dynamic) (map[string]any, error) {
	if d.IsNull() || d.IsUnknown() {
		return map[string]any{}, nil
	}
	raw, err := dynamicjson.ToJSON(d)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 || string(raw) == "null" {
		return map[string]any{}, nil
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("body must be a JSON object: %w", err)
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}

func dynamicToMap(d types.Dynamic) map[string]any {
	m, _ := dynamicToMapStrict(d)
	return m
}

// createOnlyKeys returns the list of body keys to strip before Update.
func createOnlyKeys(s types.Set) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	out := make([]string, 0, len(s.Elements()))
	for _, e := range s.Elements() {
		if sv, ok := e.(types.String); ok && !sv.IsNull() && !sv.IsUnknown() {
			out = append(out, sv.ValueString())
		}
	}
	return out
}

// mergeMaps returns a deep merge: values from b override or extend a.
func mergeMaps(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if existing, ok := out[k]; ok {
			if em, isMap := existing.(map[string]any); isMap {
				if nm, alsoMap := v.(map[string]any); alsoMap {
					out[k] = mergeMaps(em, nm)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

func populateAfterWrite(ctx context.Context, state *TFModel, bodyMap, sensitiveMap map[string]any, result callResult) diag.Diagnostics {
	var diags diag.Diagnostics

	respMap := result.Parsed
	if respMap == nil {
		respMap = map[string]any{}
	}

	readURL, d := deriveReadURL(state.Path.ValueString(), state.IDAttribute, respMap)
	diags.Append(d...)
	if d.HasError() {
		return diags
	}

	reshaped := dynamicreshape.Reshape(bodyMap, respMap, dynamicreshape.Options{
		SensitivePaths: dynamicreshape.CollectSensitivePaths(sensitiveMap),
	})
	bodyDyn, err := mapToDynamic(ctx, reshaped, state.Body)
	if err != nil {
		diags.AddError("encoding body", err.Error())
		return diags
	}
	state.Body = bodyDyn

	outputDyn, err := mapToDynamic(ctx, respMap, types.DynamicNull())
	if err != nil {
		diags.AddError("encoding output", err.Error())
		return diags
	}
	state.Output = outputDyn

	state.ID = types.StringValue(readURL)
	return diags
}

func populateAfterRead(ctx context.Context, state *TFModel, bodyMap map[string]any, result callResult) diag.Diagnostics {
	var diags diag.Diagnostics

	respMap := result.Parsed
	if respMap == nil {
		respMap = map[string]any{}
	}

	sensitiveMap, _ := dynamicToMapStrict(state.SensitiveBody)
	reshaped := dynamicreshape.Reshape(bodyMap, respMap, dynamicreshape.Options{
		SensitivePaths: dynamicreshape.CollectSensitivePaths(sensitiveMap),
	})
	bodyDyn, err := mapToDynamic(ctx, reshaped, state.Body)
	if err != nil {
		diags.AddError("encoding body", err.Error())
		return diags
	}
	state.Body = bodyDyn

	outputDyn, err := mapToDynamic(ctx, respMap, types.DynamicNull())
	if err != nil {
		diags.AddError("encoding output", err.Error())
		return diags
	}
	state.Output = outputDyn

	return diags
}

func mapToDynamic(ctx context.Context, m any, prior types.Dynamic) (types.Dynamic, error) {
	raw, err := json.Marshal(m)
	if err != nil {
		return types.DynamicNull(), err
	}
	var priorType attr.Type
	if !prior.IsNull() && !prior.IsUnknown() {
		priorType = prior.UnderlyingValue().Type(ctx)
	}
	return dynamicjson.FromJSON(raw, priorType)
}

func responseError(r callResult) string {
	var b strings.Builder
	fmt.Fprintf(&b, "status=%d", r.Status)
	if r.Err != nil {
		fmt.Fprintf(&b, " err=%s", r.Err.Error())
	}
	if len(r.Raw) > 0 {
		fmt.Fprintf(&b, " body=%s", string(r.Raw))
	}
	return b.String()
}
