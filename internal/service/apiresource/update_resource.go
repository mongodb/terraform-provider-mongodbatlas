package apiresource

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dynamicreshape"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/responseproject"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const ResourceNameUpdate = "api_update"

var (
	_ resource.ResourceWithConfigure      = &urs{}
	_ resource.ResourceWithImportState    = &urs{}
	_ resource.ResourceWithValidateConfig = &urs{}
)

func UpdateResource() resource.Resource {
	return &urs{
		RSCommon: config.RSCommon{ResourceName: ResourceNameUpdate},
	}
}

type urs struct {
	config.RSCommon
}

func (r *urs) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = UpdateResourceSchema(ctx)
}

func (r *urs) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg TFModelUpdate
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
	if overlap := responseproject.PathsOverlap(
		exportPaths(cfg.ResponseExportValues), exportPaths(cfg.ResponseExportValuesSensitive),
	); len(overlap) > 0 {
		resp.Diagnostics.AddAttributeError(path.Root("response_export_values_sensitive"),
			"path declared in both response_export_values and response_export_values_sensitive",
			fmt.Sprintf("each path must appear in only one list. Overlapping: %v", overlap))
	}
}

// Create issues the PATCH against `path` and uses the response body to populate
// state. The design spec mentions a follow-up GET, but Atlas PATCH endpoints
// reliably return the updated document on these surfaces, so we reuse the sibling
// api_resource pattern. If a target endpoint ever returns 204 (no body), revisit
// per open question #1 in the design spec.
func (r *urs) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TFModelUpdate
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bodyMap, sensitiveMap, diags := buildRequestMaps(plan.Body, plan.SensitiveBody)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := plan.Path.ValueString()
	versionHeader := resolveVersionHeader(plan.VersionHeader, plan.Preview)

	mergedBody := mergeMaps(bodyMap, sensitiveMap)
	bodyBytes, err := json.Marshal(mergedBody)
	if err != nil {
		resp.Diagnostics.AddError("encoding request body", err.Error())
		return
	}

	result := callAPI(ctx, r.Client, plan.UpdateMethod.ValueString(), url, versionHeader, bodyBytes)
	if result.Err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("API %s %s failed", plan.UpdateMethod.ValueString(), url),
			responseError(result),
		)
		return
	}

	state := plan
	state.VersionHeader = types.StringValue(versionHeader)
	state.ID = types.StringValue(url)
	resp.Diagnostics.Append(populateAfterWriteUpdate(ctx, &state, bodyMap, sensitiveMap, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// populateAfterWriteUpdate mirrors populateAfterWrite from resource.go but
// skips deriveReadURL (api_update never derives a URL — path IS the URL).
func populateAfterWriteUpdate(ctx context.Context, state *TFModelUpdate, bodyMap, sensitiveMap map[string]any, result callResult) diag.Diagnostics {
	var diags diag.Diagnostics

	respMap := result.Parsed
	if respMap == nil {
		respMap = map[string]any{}
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

	outputDyn, outputSensitiveDyn, err := projectToDynamics(ctx, respMap,
		exportPaths(state.ResponseExportValues), exportPaths(state.ResponseExportValuesSensitive))
	if err != nil {
		diags.AddError("encoding output", err.Error())
		return diags
	}
	state.Output = outputDyn
	state.OutputSensitive = outputSensitiveDyn
	return diags
}

func (r *urs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TFModelUpdate
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bodyMap, _, diags := buildRequestMaps(state.Body, types.DynamicNull())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// No deriveReadURL call (unlike rs.Read in resource.go): api_update has no
	// id_attribute, so path IS the read URL.
	readURL := state.Path.ValueString()
	versionHeader := resolveVersionHeader(state.VersionHeader, state.Preview)

	result := callAPI(ctx, r.Client, defaultReadMethod, readURL, versionHeader, nil)
	// Some endpoints accept a preview content-type on write but not on read
	// (e.g. PATCH /streams/{tenantName} has a preview variant; GET does not).
	// When we configured `preview = true` for writes, fall back to today's GA
	// version for the read. Atlas serializes the same entity regardless of
	// content-type, so hidden fields like failoverRegions still come back.
	if result.Status == 406 && state.Preview.ValueBool() {
		result = callAPI(ctx, r.Client, defaultReadMethod, readURL, todayVersionHeader(), nil)
	}
	if result.NotFound {
		// Entity is gone — likely the typed resource was destroyed. Drop our
		// state. A subsequent apply will fail at Create until the entity is
		// recreated (or the user removes this resource).
		resp.State.RemoveResource(ctx)
		return
	}
	if result.Err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API GET %s failed", readURL), responseError(result))
		return
	}

	resp.Diagnostics.Append(populateAfterReadUpdate(ctx, &state, bodyMap, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// populateAfterReadUpdate is structurally identical to populateAfterRead in
// resource.go, retyped for TFModelUpdate. The reshape engine already gives us
// the filter-to-patched-keys semantics by treating bodyMap as the template.
// Keep in sync with populateAfterRead.
func populateAfterReadUpdate(ctx context.Context, state *TFModelUpdate, bodyMap map[string]any, result callResult) diag.Diagnostics {
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

	outputDyn, outputSensitiveDyn, err := projectToDynamics(ctx, respMap,
		exportPaths(state.ResponseExportValues), exportPaths(state.ResponseExportValuesSensitive))
	if err != nil {
		diags.AddError("encoding output", err.Error())
		return diags
	}
	state.Output = outputDyn
	state.OutputSensitive = outputSensitiveDyn
	return diags
}

func (r *urs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TFModelUpdate
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bodyMap, sensitiveMap, diags := buildRequestMaps(plan.Body, plan.SensitiveBody)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := plan.Path.ValueString()
	versionHeader := resolveVersionHeader(plan.VersionHeader, plan.Preview)

	mergedBody := mergeMaps(bodyMap, sensitiveMap)
	bodyBytes, err := json.Marshal(mergedBody)
	if err != nil {
		resp.Diagnostics.AddError("encoding request body", err.Error())
		return
	}

	result := callAPI(ctx, r.Client, plan.UpdateMethod.ValueString(), url, versionHeader, bodyBytes)
	if result.Err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("API %s %s failed", plan.UpdateMethod.ValueString(), url),
			responseError(result),
		)
		return
	}

	newState := plan
	newState.VersionHeader = types.StringValue(versionHeader)
	newState.ID = types.StringValue(url)
	resp.Diagnostics.Append(populateAfterWriteUpdate(ctx, &newState, bodyMap, sensitiveMap, result)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *urs) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) { //nolint:gocritic // framework interface dictates the signature
	// No-op. The entity belongs to another (typed) resource. Removing this
	// block from config leaves the patched field at its last-applied value.
}

func (r *urs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	// path doubles as id for this resource — set both so subsequent Read can execute.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), req.ID)...)
	resp.Diagnostics.AddWarning(
		"Import is best-effort — re-declare config and rotate secrets",
		"Terraform import recovers only the resource URL into `id` and `path`. To finish: "+
			"(1) re-declare `body` and other config in HCL; "+
			"(2) run `terraform plan` and adjust HCL until the diff is clean; "+
			"(3) re-supply or rotate `sensitive_body` — Atlas does not return secrets on GET, so the previous value cannot be recovered.",
	)
}
