package streamconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &streamConnectionDS{}
var _ datasource.DataSourceWithConfigure = &streamConnectionDS{}

func DataSource() datasource.DataSource {
	return &streamConnectionDS{
		DSCommon: config.DSCommon{
			DataSourceName: streamConnectionName,
		},
	}
}

type streamConnectionDS struct {
	config.DSCommon
}

func (d *streamConnectionDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "connection_name"},
		OverridenFields: map[string]dsschema.Attribute{
			"instance_name": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-readable label that identifies the stream instance. Conflicts with `workspace_name`.",
				DeprecationMessage:  fmt.Sprintf(constant.DeprecationParamWithReplacement, "workspace_name"),
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("workspace_name")),
				},
			},
			"workspace_name": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Human-readable label that identifies the stream instance. This is an alias for `instance_name`. Conflicts with `instance_name`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("instance_name")),
				},
			},
		},
	})
}

func (d *streamConnectionDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamConnectionConfig TFStreamConnectionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamConnectionConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamConnectionConfig.ProjectID.ValueString()
	workspaceOrInstanceName := getWorkspaceOrInstanceName(&streamConnectionConfig)
	if workspaceOrInstanceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}
	connectionName := streamConnectionConfig.ConnectionName.ValueString()
	apiResp, _, err := connV2.StreamsApi.GetStreamConnection(ctx, projectID, workspaceOrInstanceName, connectionName).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching resource", err.Error())
		return
	}

	instanceName := streamConnectionConfig.InstanceName.ValueString()
	workspaceName := streamConnectionConfig.WorkspaceName.ValueString()
	newStreamConnectionModel, diags := NewTFStreamConnectionWithInstanceName(ctx, projectID, instanceName, workspaceName, nil, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionModel)...)
}
