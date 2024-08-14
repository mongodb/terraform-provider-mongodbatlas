package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func (d *streamProcessorsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamConnectionsConfig TFStreamProcessorsDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamConnectionsConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamConnectionsConfig.ProjectID.ValueString()
	instanceName := streamConnectionsConfig.InstanceName.ValueString()
	itemsPerPage := streamConnectionsConfig.ItemsPerPage.ValueInt64Pointer()
	pageNum := streamConnectionsConfig.PageNum.ValueInt64Pointer()

	apiResp, _, err := connV2.StreamsApi.ListStreamProcessorsWithParams(ctx, &admin.ListStreamProcessorsApiParams{
		GroupId:      projectID,
		TenantName:   instanceName,
		ItemsPerPage: conversion.Int64PtrToIntPtr(itemsPerPage),
		PageNum:      conversion.Int64PtrToIntPtr(pageNum),
	}).Execute()

	if err != nil {
		resp.Diagnostics.AddError("error fetching results", err.Error())
		return
	}

	newStreamConnectionsModel, diags := NewTFStreamProcessors(ctx, &streamConnectionsConfig, apiResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionsModel)...)
}
