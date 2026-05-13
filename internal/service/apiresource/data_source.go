package apiresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/responseproject"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const DataSourceName = "api_resource"

var (
	_ datasource.DataSourceWithConfigure      = &ds{}
	_ datasource.DataSourceWithValidateConfig = &ds{}
)

func DataSource() datasource.DataSource {
	return &ds{
		DSCommon: config.DSCommon{DataSourceName: DataSourceName},
	}
}

type ds struct {
	config.DSCommon
}

func (d *ds) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Reads any Atlas Admin API GET endpoint. Declare paths in `response_export_values` " +
			"and/or `response_export_values_sensitive` to opt fields into state. By default both outputs are null.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Synthetic identifier (the path).",
			},
			"path": schema.StringAttribute{
				Required:    true,
				Description: "Atlas API path to GET.",
			},
			"version_header": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Description: "Atlas API version media type. When unset, defaults to today's UTC date in the form " +
					"`application/vnd.atlas.<YYYY-MM-DD>+json` (Atlas snaps the date down to the latest published " +
					"version on or before it). Mutually exclusive with `preview`.",
			},
			"preview": schema.BoolAttribute{
				Optional:    true,
				Description: "Shorthand for `version_header = \"" + previewVersionHeader + "\"`. Mutually exclusive with `version_header`.",
			},
			"response_export_values": schema.ListAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
				Description: "Dotted paths into the API response to retain in `output`.",
			},
			"response_export_values_sensitive": schema.ListAttribute{
				Optional:    true,
				ElementType: basetypes.StringType{},
				Description: "Dotted paths whose values land in `output_sensitive` (Sensitive). A path must not appear in both lists.",
			},
			"output": schema.DynamicAttribute{
				Computed:    true,
				Description: "Projected response containing paths listed in `response_export_values`. Null when none declared.",
			},
			"output_sensitive": schema.DynamicAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Projected response containing paths listed in `response_export_values_sensitive`. Null when none declared.",
			},
		},
	}
	conversion.UpdateSchemaDescription(&resp.Schema)
}

func (d *ds) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var cfg TFModelDS
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

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModelDS
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	versionHeader := resolveVersionHeader(state.VersionHeader, state.Preview)
	state.VersionHeader = types.StringValue(versionHeader)

	result := callAPI(ctx, d.Client, defaultReadMethod, state.Path.ValueString(), versionHeader, nil)
	if result.NotFound {
		resp.Diagnostics.AddError("resource not found", fmt.Sprintf("GET %s returned 404", state.Path.ValueString()))
		return
	}
	if result.Err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("API GET %s failed", state.Path.ValueString()), responseError(result))
		return
	}
	respMap := result.Parsed
	if respMap == nil {
		respMap = map[string]any{}
	}
	outputDyn, outputSensitiveDyn, err := projectToDynamics(ctx, respMap,
		exportPaths(state.ResponseExportValues), exportPaths(state.ResponseExportValuesSensitive))
	if err != nil {
		resp.Diagnostics.AddError("encoding output", err.Error())
		return
	}
	state.Output = outputDyn
	state.OutputSensitive = outputSensitiveDyn
	state.ID = types.StringValue(state.Path.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
