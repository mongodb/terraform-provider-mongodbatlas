package streaminstance

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

var _ datasource.DataSource = &streamInstancesDS{}
var _ datasource.DataSourceWithConfigure = &streamInstancesDS{}

func PluralDataSource() datasource.DataSource {
	return &streamInstancesDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", streamInstanceName),
		},
	}
}

type streamInstancesDS struct {
	config.DSCommon
}

func (d *streamInstancesDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields:  []string{"project_id"},
		HasLegacyFields: true,
	})
}

func (d *streamInstancesDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamInstancesConfig TFStreamInstancesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamInstancesConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamInstancesConfig.ProjectID.ValueString()
	itemsPerPage := streamInstancesConfig.ItemsPerPage.ValueInt64Pointer()
	pageNum := streamInstancesConfig.PageNum.ValueInt64Pointer()
	apiResp, _, err := connV2.StreamsApi.ListStreamInstancesWithParams(ctx, &admin.ListStreamInstancesApiParams{
		GroupId:      projectID,
		ItemsPerPage: conversion.Int64PtrToIntPtr(itemsPerPage),
		PageNum:      conversion.Int64PtrToIntPtr(pageNum),
	}).Execute()
	if err != nil {
		resp.Diagnostics.AddError("error fetching results", err.Error())
		return
	}

	newStreamInstancesModel, diags := NewTFStreamInstances(ctx, &streamInstancesConfig, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamInstancesModel)...)
}

type TFStreamInstancesModel struct {
	ID           types.String            `tfsdk:"id"`
	ProjectID    types.String            `tfsdk:"project_id"`
	Results      []TFStreamInstanceModel `tfsdk:"results"`
	PageNum      types.Int64             `tfsdk:"page_num"`
	ItemsPerPage types.Int64             `tfsdk:"items_per_page"`
	TotalCount   types.Int64             `tfsdk:"total_count"`
}
