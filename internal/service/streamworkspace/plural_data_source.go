package streamworkspace

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streaminstance"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

var _ datasource.DataSource = &streamsWorkspacesDS{}
var _ datasource.DataSourceWithConfigure = &streamsWorkspacesDS{}

func PluralDataSource() datasource.DataSource {
	return &streamsWorkspacesDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", streamsWorkspaceName),
		},
	}
}

type streamsWorkspacesDS struct {
	config.DSCommon
}

func (d *streamsWorkspacesDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields:  []string{"project_id"},
		HasLegacyFields: true,
	})
}

func (d *streamsWorkspacesDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamsWorkspacesConfig TFStreamsWorkspacesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamsWorkspacesConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamsWorkspacesConfig.ProjectID.ValueString()
	itemsPerPage := streamsWorkspacesConfig.ItemsPerPage.ValueInt64Pointer()
	pageNum := streamsWorkspacesConfig.PageNum.ValueInt64Pointer()
	apiResp, _, err := connV2.StreamsApi.ListStreamWorkspacesWithParams(ctx, &admin.ListStreamWorkspacesApiParams{
		GroupId:      projectID,
		ItemsPerPage: conversion.Int64PtrToIntPtr(itemsPerPage),
		PageNum:      conversion.Int64PtrToIntPtr(pageNum),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching results", err.Error())
		return
	}

	newStreamsWorkspacesModel, diags := NewTFStreamsWorkspaces(ctx, &streamsWorkspacesConfig, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamsWorkspacesModel)...)
}

type TFStreamsWorkspacesModel struct {
	ID           types.String `tfsdk:"id"`
	ProjectID    types.String `tfsdk:"project_id"`
	Results      []TFModel    `tfsdk:"results"`
	PageNum      types.Int64  `tfsdk:"page_num"`
	ItemsPerPage types.Int64  `tfsdk:"items_per_page"`
	TotalCount   types.Int64  `tfsdk:"total_count"`
}

func NewTFStreamsWorkspaces(ctx context.Context, streamsWorkspacesConfig *TFStreamsWorkspacesModel, apiResp *admin.PaginatedApiStreamsTenant) (*TFStreamsWorkspacesModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert the stream instances response to stream instances model first
	instancesModel := &streaminstance.TFStreamInstancesModel{
		ID:           streamsWorkspacesConfig.ID,
		ProjectID:    streamsWorkspacesConfig.ProjectID,
		PageNum:      streamsWorkspacesConfig.PageNum,
		ItemsPerPage: streamsWorkspacesConfig.ItemsPerPage,
	}

	newInstancesModel, instanceDiags := streaminstance.NewTFStreamInstances(ctx, instancesModel, apiResp)
	if instanceDiags.HasError() {
		diags.Append(instanceDiags...)
		return nil, diags
	}

	// Convert each instance result to workspace result
	workspaceResults := make([]TFModel, len(newInstancesModel.Results))
	for i := range newInstancesModel.Results {
		workspaceResults[i].FromInstanceModel(&newInstancesModel.Results[i])
	}

	return &TFStreamsWorkspacesModel{
		ID:           newInstancesModel.ID,
		ProjectID:    newInstancesModel.ProjectID,
		Results:      workspaceResults,
		PageNum:      newInstancesModel.PageNum,
		ItemsPerPage: newInstancesModel.ItemsPerPage,
		TotalCount:   newInstancesModel.TotalCount,
	}, diags
}
