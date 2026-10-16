package advancedcluster

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312026/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

// listRetryTimeout is the total time budget for retrying a cluster list on CLUSTER_NOT_FOUND.
const listRetryTimeout = 1 * time.Minute

var _ datasource.DataSource = &pluralDS{}
var _ datasource.DataSourceWithConfigure = &pluralDS{}

func PluralDataSource() datasource.DataSource {
	return &pluralDS{
		DataSourceName: fmt.Sprintf("%ss", resourceName),
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
	model := d.readClusters(ctx, diags, &state)
	if model != nil {
		model.UseEffectiveFields = state.UseEffectiveFields // Set Optional Terraform-only attribute.
		diags.Append(resp.State.Set(ctx, model)...)
	}
}

func (d *pluralDS) readClusters(ctx context.Context, diags *diag.Diagnostics, pluralModel *TFModelPluralDS) *TFModelPluralDS {
	projectID := pluralModel.ProjectID.ValueString()
	outs := &TFModelPluralDS{
		ProjectID: pluralModel.ProjectID,
		// Non-nil so a project without clusters serializes results as an empty list, not null.
		Results: make([]*TFModelDS, 0),
	}
	basicClusters := d.getBasicClusters(ctx, diags, projectID, pluralModel.UseEffectiveFields)
	if diags.HasError() {
		return nil
	}
	outs.Results = append(outs.Results, basicClusters...)

	flexClusters := d.getFlexClusters(ctx, diags, projectID, pluralModel.UseEffectiveFields)
	if diags.HasError() {
		return nil
	}
	outs.Results = append(outs.Results, flexClusters...)
	return outs
}

// RetryListOnClusterNotFound retries a cluster list while Atlas returns the transient
// CLUSTER_NOT_FOUND error, which the list endpoint can return for the whole request when a
// cluster is deleted concurrently.
func RetryListOnClusterNotFound[T any](ctx context.Context, list func() (T, error)) (T, error) {
	var result T
	err := retry.RetryContext(ctx, listRetryTimeout, func() *retry.RetryError {
		var err error
		result, err = list()
		if err == nil {
			return nil
		}
		if admin.IsErrorCode(err, ErrorCodeClusterNotFound) {
			return retry.RetryableError(err)
		}
		return retry.NonRetryableError(err)
	})
	return result, err
}

// listClusterModels lists every cluster and converts it, skipping clusters deleted between
// the list and their conversion. It returns nil when an error must fail the read.
func listClusterModels[T any](ctx context.Context, diags *diag.Diagnostics, projectID string, useEffectiveFields types.Bool, list func() ([]T, error), convert func(*T) *TFModelDS) []*TFModelDS {
	items, err := RetryListOnClusterNotFound(ctx, list)
	if err != nil {
		addListError(diags, projectID, err)
		return nil
	}
	var results []*TFModelDS
	for i := range items {
		modelOutDS := convert(&items[i])
		RemoveClusterNotFoundErrors(diags)
		if diags.HasError() {
			return nil
		}
		if modelOutDS != nil { // diags can be empty after RemoveClusterNotFoundErrors while the model is nil.
			modelOutDS.UseEffectiveFields = useEffectiveFields
			results = append(results, modelOutDS)
		}
	}
	return results
}

// getBasicClusters gets the dedicated and tenant clusters.
func (d *pluralDS) getBasicClusters(ctx context.Context, diags *diag.Diagnostics, projectID string, useEffectiveFields types.Bool) []*TFModelDS {
	api := d.Client.AtlasV2.ClustersAPI
	params := admin.ListClustersApiParams{
		GroupId:                    projectID,
		UseEffectiveInstanceFields: new(useEffectiveFields.ValueBool()),
	}
	return listClusterModels(ctx, diags, projectID, useEffectiveFields,
		func() ([]admin.ClusterDescription20240805, error) {
			return dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ClusterDescription20240805], *http.Response, error) {
				return api.ListClustersWithParams(ctx, &params).PageNum(pageNum).Execute()
			})
		},
		func(clusterResp *admin.ClusterDescription20240805) *TFModelDS {
			return convertBasicClusterToDS(ctx, diags, d.Client, clusterResp)
		})
}

func (d *pluralDS) getFlexClusters(ctx context.Context, diags *diag.Diagnostics, projectID string, useEffectiveFields types.Bool) []*TFModelDS {
	return listClusterModels(ctx, diags, projectID, useEffectiveFields,
		func() ([]admin.FlexClusterDescription20241113, error) {
			list, err := flexcluster.ListFlexClusters(ctx, projectID, d.Client.AtlasV2.FlexClustersAPI)
			if list == nil {
				return nil, err
			}
			return *list, err
		},
		func(flexClusterResp *admin.FlexClusterDescription20241113) *TFModelDS {
			return convertFlexClusterToDS(ctx, diags, flexClusterResp)
		})
}

// addListError adds a standardized error for cluster list operations.
func addListError(diags *diag.Diagnostics, projectID string, err error) {
	diags.AddError(errorList, fmt.Sprintf(errorListDetail, projectID, err.Error()))
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
