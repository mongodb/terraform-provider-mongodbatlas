package streamconnection

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

type TFStreamConnectionsDSModel struct {
	ID           types.String              `tfsdk:"id"`
	ProjectID    types.String              `tfsdk:"project_id"`
	InstanceName types.String              `tfsdk:"instance_name"`
	Results      []TFStreamConnectionModel `tfsdk:"results"`
	PageNum      types.Int64               `tfsdk:"page_num"`
	ItemsPerPage types.Int64               `tfsdk:"items_per_page"`
	TotalCount   types.Int64               `tfsdk:"total_count"`
}

func (d *streamConnectionsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// TODO: THIS WILL BE REMOVED BEFORE MERGING, check old data source schema and new auto-generated schema are the same
	ds1 := dsschema.PaginatedDSSchema(
		map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"instance_name": schema.StringAttribute{
				Required: true,
			},
		},
		DSAttributes(false))
	conversion.UpdateSchemaDescription(&ds1)

	requiredFields := []string{"project_id", "instance_name"}
	ds2 := conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), requiredFields, nil, nil, "", true)
	if diff := cmp.Diff(ds1, ds2); diff != "" {
		log.Fatal(diff)
	}
	resp.Schema = ds2
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
