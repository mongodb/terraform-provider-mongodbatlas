package advancedcluster

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

func PluralDataSource() datasource.DataSource {
	return &pluralDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", resourceName),
		},
	}
}

type pluralDS struct {
	config.DSCommon
}

func (d *pluralDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = pluralDataSourceSchema(ctx)
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TFModelPluralDS
	diags := &resp.Diagnostics
	diags.Append(req.Config.Get(ctx, &state)...)
	if diags.HasError() {
		return
	}
	model, diags := d.readClusters(ctx, diags, &state)
	resp.Diagnostics = *diags
	if model != nil {
		model.UseEffectiveFields = state.UseEffectiveFields // Set Optional Terraform-only attribute.
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *pluralDS) readClusters(ctx context.Context, diags *diag.Diagnostics, pluralModel *TFModelPluralDS) (*TFModelPluralDS, *diag.Diagnostics) {
	projectID := pluralModel.ProjectID.ValueString()
	outs := &TFModelPluralDS{
		ProjectID: pluralModel.ProjectID,
	}
	basicClusters := d.getBasicClusters(ctx, diags, projectID, pluralModel.UseEffectiveFields)
	if diags.HasError() {
		return nil, diags
	}
	outs.Results = append(outs.Results, basicClusters...)

	flexClusters := d.getFlexClusters(ctx, diags, projectID, pluralModel.UseEffectiveFields)
	if diags.HasError() {
		return nil, diags
	}
	outs.Results = append(outs.Results, flexClusters...)
	return outs, diags
}

// getBasicClusters gets the dedicated and tenant clusters.
func (d *pluralDS) getBasicClusters(ctx context.Context, diags *diag.Diagnostics, projectID string, useEffectiveFields types.Bool) []*TFModelDS {
	var results []*TFModelDS
	api := d.Client.AtlasV2.ClustersApi
	params := admin.ListClustersApiParams{
		GroupId:                    projectID,
		UseEffectiveInstanceFields: conversion.Pointer(useEffectiveFields.ValueBool()),
	}
	list, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ClusterDescription20240805], *http.Response, error) {
		return api.ListClustersWithParams(ctx, &params).PageNum(pageNum).Execute()
	})
	if err != nil {
		addListError(diags, projectID, err)
		RemoveClusterNotFoundErrors(diags)
		return nil
	}
	for i := range list {
		clusterResp := &list[i]
		modelOutDS := convertBasicClusterToDS(ctx, diags, d.Client, clusterResp)
		if !appendClusterModelIfValid(diags, modelOutDS, useEffectiveFields, &results) {
			return nil
		}
	}
	return results
}

func (d *pluralDS) getFlexClusters(ctx context.Context, diags *diag.Diagnostics, projectID string, useEffectiveFields types.Bool) []*TFModelDS {
	var results []*TFModelDS
	listFlexClusters, err := flexcluster.ListFlexClusters(ctx, projectID, d.Client.AtlasV2.FlexClustersApi)
	if err != nil {
		addListError(diags, projectID, err)
		RemoveClusterNotFoundErrors(diags)
		return nil
	}
	for i := range *listFlexClusters {
		flexClusterResp := (*listFlexClusters)[i]
		modelOutDS := convertFlexClusterToDS(ctx, diags, &flexClusterResp)
		if !appendClusterModelIfValid(diags, modelOutDS, useEffectiveFields, &results) {
			return nil
		}
	}
	return results
}

// addListError adds a standardized error for cluster list operations.
func addListError(diags *diag.Diagnostics, projectID string, err error) {
	diags.AddError(errorList, fmt.Sprintf(errorListDetail, projectID, err.Error()))
}

// appendClusterModelIfValid removes CLUSTER_NOT_FOUND errors from diags and appends the model to results if valid.
// Returns false if processing should stop due to remaining errors after filtering.
func appendClusterModelIfValid(diags *diag.Diagnostics, modelOutDS *TFModelDS, useEffectiveFields types.Bool, results *[]*TFModelDS) bool {
	RemoveClusterNotFoundErrors(diags)
	if diags.HasError() {
		return false
	}
	if modelOutDS != nil { // diags could be empty because of RemoveClusterNotFoundErrors but modelOutDS be nil.
		modelOutDS.UseEffectiveFields = useEffectiveFields
		*results = append(*results, modelOutDS)
	}
	return true
}

// RemoveClusterNotFoundErrors removes CLUSTER_NOT_FOUND errors from diags in-place.
func RemoveClusterNotFoundErrors(diags *diag.Diagnostics) {
	filtered := diag.Diagnostics{}
	for _, d := range *diags {
		if d.Severity() == diag.SeverityError && strings.Contains(d.Detail(), "CLUSTER_NOT_FOUND") {
			continue // Skip CLUSTER_NOT_FOUND errors
		}
		filtered.Append(d)
	}
	*diags = filtered
}
