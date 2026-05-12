package apiresource

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

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
	resp.Diagnostics.AddError("not implemented", "api_update Create is not yet implemented")
}

func (r *urs) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.Diagnostics.AddError("not implemented", "api_update Read is not yet implemented")
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
