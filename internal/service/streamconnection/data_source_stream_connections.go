package streamconnection

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
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
		RequiredFields:  []string{"project_id", "instance_name"},
		HasLegacyFields: true,
	})
}

func (d *streamConnectionsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamConnectionsConfig TFStreamConnectionsDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamConnectionsConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamConnectionsConfig.ProjectID.ValueString()
	instanceName := streamConnectionsConfig.InstanceName.ValueString()
	itemsPerPage := streamConnectionsConfig.ItemsPerPage.ValueInt64Pointer()
	pageNum := streamConnectionsConfig.PageNum.ValueInt64Pointer()

	apiResp, _, err := connV2.StreamsApi.ListStreamConnectionsWithParams(ctx, &admin.ListStreamConnectionsApiParams{
		GroupId:      projectID,
		TenantName:   instanceName,
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
	ID           types.String              `tfsdk:"id"`
	ProjectID    types.String              `tfsdk:"project_id"`
	InstanceName types.String              `tfsdk:"instance_name"`
	Results      []TFStreamConnectionModel `tfsdk:"results"`
	PageNum      types.Int64               `tfsdk:"page_num"`
	ItemsPerPage types.Int64               `tfsdk:"items_per_page"`
	TotalCount   types.Int64               `tfsdk:"total_count"`
}
