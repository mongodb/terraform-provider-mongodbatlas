package apikeyprojectassignment

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields: []string{"project_id"},
	})
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfModel TFModelDSP
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := tfModel.ProjectId.ValueString()

	params := &admin.ListProjectApiKeysApiParams{
		GroupId: tfModel.ProjectId.ValueString(),
	}

	apiKeys, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.ApiKeyUserDetails], *http.Response, error) {
		request := connV2.ProgrammaticAPIKeysApi.ListProjectApiKeysWithParams(ctx, params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		resp.Diagnostics.AddError("error reading plural data source", err.Error())
		return
	}

	newAPIKeyProjectAssignmentsModel, diags := NewTFModelDSP(ctx, projectID, apiKeys)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newAPIKeyProjectAssignmentsModel)...)
}
