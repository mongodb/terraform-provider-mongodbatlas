package streamconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

var _ datasource.DataSource = &streamConnectionsDS{}
var _ datasource.DataSourceWithConfigure = &streamConnectionsDS{}

func PluralDataSource() datasource.DataSource {
	return &streamConnectionsDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", streamConnectionName),
		},
	}
}

type streamConnectionsDS struct {
	config.DSCommon
}

func (d *streamConnectionsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields:  []string{"project_id"},
		HasLegacyFields: true,
		OverridenRootFields: map[string]dsschema.Attribute{
			"instance_name": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Label that identifies the stream processing workspace. Conflicts with `workspace_name`.",
				DeprecationMessage:  fmt.Sprintf(constant.DeprecationParamWithReplacement, "workspace_name"),
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("workspace_name")),
				},
			},
			"workspace_name": dsschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Label that identifies the stream processing workspace. This is an alias for `instance_name`. Conflicts with `instance_name`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("instance_name")),
				},
			},
		},
	})
}

// getWorkspaceOrInstanceNameForDS returns the workspace name from either instance_name or workspace_name field for datasource model
func getWorkspaceOrInstanceNameForDS(model *TFStreamConnectionsDSModel) string {
	if !model.WorkspaceName.IsNull() && !model.WorkspaceName.IsUnknown() {
		return model.WorkspaceName.ValueString()
	}
	if !model.InstanceName.IsNull() && !model.InstanceName.IsUnknown() {
		return model.InstanceName.ValueString()
	}
	return ""
}

func (d *streamConnectionsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamConnectionsConfig TFStreamConnectionsDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamConnectionsConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamConnectionsConfig.ProjectID.ValueString()
	workspaceOrInstanceName := getWorkspaceOrInstanceNameForDS(&streamConnectionsConfig)
	if workspaceOrInstanceName == "" {
		resp.Diagnostics.AddError("validation error", "workspace_name must be provided")
		return
	}
	itemsPerPage := streamConnectionsConfig.ItemsPerPage.ValueInt64Pointer()
	pageNum := streamConnectionsConfig.PageNum.ValueInt64Pointer()

	apiResp, _, err := connV2.StreamsApi.ListStreamConnectionsWithParams(ctx, &admin.ListStreamConnectionsApiParams{
		GroupId:      projectID,
		TenantName:   workspaceOrInstanceName,
		ItemsPerPage: conversion.Int64PtrToIntPtr(itemsPerPage),
		PageNum:      conversion.Int64PtrToIntPtr(pageNum),
	}).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error fetching results", err.Error())
		return
	}

	newStreamConnectionsModel, diags := NewTFStreamConnections(ctx, &streamConnectionsConfig, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionsModel)...)
}

type TFStreamConnectionsDSModel struct {
	ID            types.String              `tfsdk:"id"`
	ProjectID     types.String              `tfsdk:"project_id"`
	InstanceName  types.String              `tfsdk:"instance_name"`
	WorkspaceName types.String              `tfsdk:"workspace_name"`
	Results       []TFStreamConnectionModel `tfsdk:"results"`
	PageNum       types.Int64               `tfsdk:"page_num"`
	ItemsPerPage  types.Int64               `tfsdk:"items_per_page"`
	TotalCount    types.Int64               `tfsdk:"total_count"`
}
