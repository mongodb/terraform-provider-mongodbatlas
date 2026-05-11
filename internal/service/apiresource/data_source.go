package apiresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
		MarkdownDescription: "Reads any Atlas Admin API GET endpoint. The full response is exposed via `output`.",
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
			"output": schema.DynamicAttribute{
				Computed:    true,
				Description: "Full API response from the most recent successful operation.",
			},
		},
	}
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
	outputDyn, err := mapToDynamic(ctx, respMap, types.DynamicNull())
	if err != nil {
		resp.Diagnostics.AddError("encoding output", err.Error())
		return
	}
	state.Output = outputDyn
	state.ID = types.StringValue(state.Path.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
