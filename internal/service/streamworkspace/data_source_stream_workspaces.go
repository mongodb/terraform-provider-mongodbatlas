package streamworkspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

var _ datasource.DataSource = &streamWorkspacesDS{}
var _ datasource.DataSourceWithConfigure = &streamWorkspacesDS{}

func PluralDataSource() datasource.DataSource {
	return &streamWorkspacesDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", streamWorkspaceName),
		},
	}
}

type streamWorkspacesDS struct {
	config.DSCommon
}

func (d *streamWorkspacesDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields:  []string{"project_id"},
		HasLegacyFields: true,
	})
}

func (d *streamWorkspacesDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamWorkspacesConfig TFStreamWorkspacesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamWorkspacesConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamWorkspacesConfig.ProjectID.ValueString()
	itemsPerPage := streamWorkspacesConfig.ItemsPerPage.ValueInt64Pointer()
	pageNum := streamWorkspacesConfig.PageNum.ValueInt64Pointer()
	apiResp, _, err := connV2.StreamsApi.ListStreamInstancesWithParams(ctx, &admin.ListStreamInstancesApiParams{
		GroupId:      projectID,
		ItemsPerPage: conversion.Int64PtrToIntPtr(itemsPerPage),
		PageNum:      conversion.Int64PtrToIntPtr(pageNum),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching results", err.Error())
		return
	}

	newStreamWorkspacesModel, diags := NewTFStreamWorkspaces(ctx, &streamWorkspacesConfig, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamWorkspacesModel)...)
}

type TFStreamWorkspacesModel struct {
	ID           types.String             `tfsdk:"id"`
	ProjectID    types.String             `tfsdk:"project_id"`
	Results      []TFStreamWorkspaceModel `tfsdk:"results"`
	PageNum      types.Int64              `tfsdk:"page_num"`
	ItemsPerPage types.Int64              `tfsdk:"items_per_page"`
	TotalCount   types.Int64              `tfsdk:"total_count"`
}
