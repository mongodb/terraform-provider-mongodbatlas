package streamprocessor

import (
	"context"
	"errors"
	"net/http"

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

	params := admin.ListStreamProcessorsApiParams{
		GroupId:      projectID,
		TenantName:   instanceName,
		ItemsPerPage: conversion.Int64PtrToIntPtr(itemsPerPage),
		PageNum:      conversion.Int64PtrToIntPtr(pageNum),
	}
	sdkProcessors, err := AllPages[admin.StreamsProcessorWithStats](ctx, func(ctx context.Context, pageNum int) (PaginateResponse[admin.StreamsProcessorWithStats], *http.Response, error) {
		request := connV2.StreamsApi.ListStreamProcessorsWithParams(ctx, &params)
		request.PageNum(pageNum)
		return request.Execute()
	})
	
	apiResp, _, err := connV2.StreamsApi.ListStreamProcessorsWithParams(ctx, &params).Execute()
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

type PaginateResponse[T any] interface {
	GetResults() []T
}

// type PaginateRequest[T any] interface {
// 	Execute() (*PaginateResponse[T], *http.Response, error)
// }

func AllPages[T any](ctx context.Context, call func(ctx context.Context, pageNum int) (PaginateResponse[T], *http.Response, error)) ([]T, error) {
	var results []T
	for i := 0; ; i++ {
		resp, _, err := call(ctx, i)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, errors.New("no response")
		}
		currentResults := resp.GetResults()
		if len(currentResults) == 0 {
			break
		}
		results = append(results, currentResults...)
	}
	return results, nil
}

// type Pager[T any] struct {

// }

// var PagerDone = errors.New("no more pages to iterate")

// func AllPages()
