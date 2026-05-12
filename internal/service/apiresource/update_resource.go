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
}

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

	outputDyn, err := mapToDynamic(ctx, respMap, types.DynamicNull())
	if err != nil {
		diags.AddError("encoding output", err.Error())
		return diags
	}
	state.Output = outputDyn
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

	outputDyn, err := mapToDynamic(ctx, respMap, types.DynamicNull())
	if err != nil {
		diags.AddError("encoding output", err.Error())
		return diags
	}
	state.Output = outputDyn
	return diags
}

func (r *urs) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "api_update Update is not yet implemented")
}

func (r *urs) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op. The entity belongs to another (typed) resource. Removing this
	// block from config leaves the patched field at its last-applied value.
}

func (r *urs) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	// path doubles as id for this resource — set both so subsequent Read can execute.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("path"), req.ID)...)
	resp.Diagnostics.AddWarning(
		"Import does not recover body / sensitive_body",
		"After import, re-declare body and sensitive_body in HCL. The next plan will surface drift until the configured body matches what Atlas returns.",
	)
}
