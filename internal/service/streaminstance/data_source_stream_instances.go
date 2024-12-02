package streaminstance

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20241113001/admin"
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

type TFStreamInstancesModel struct {
	ID           types.String            `tfsdk:"id"`
	ProjectID    types.String            `tfsdk:"project_id"`
	Results      []TFStreamInstanceModel `tfsdk:"results"`
	PageNum      types.Int64             `tfsdk:"page_num"`
	ItemsPerPage types.Int64             `tfsdk:"items_per_page"`
	TotalCount   types.Int64             `tfsdk:"total_count"`
}

func (d *streamInstancesDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: THIS WILL BE REMOVED BEFORE MERGING, check old data source schema and new auto-generated schema are the same
	ds1 := dsschema.PaginatedDSSchema(
		map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required: true,
			},
		},
		DSAttributes(false))
	conversion.UpdateSchemaDescription(&ds1)

	requiredFields := []string{"project_id"}
	ds2 := conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil, nil, true)
	if diff := cmp.Diff(ds1, ds2); diff != "" {
		log.Fatal(diff)
	}
	resp.Schema = ds2
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
