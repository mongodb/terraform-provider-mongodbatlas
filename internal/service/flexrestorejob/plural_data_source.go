package flexrestorejob

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312008/admin"

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
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(DataSourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "name"},
	})
}

func (d *pluralDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfModel TFFlexRestoreJobsDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := tfModel.ProjectID.ValueString()
	name := tfModel.Name.ValueString()
	flexRestoreJobs, err := ListFlexRestoreJobs(ctx, projectID, name, d.Client.AtlasV2.FlexRestoreJobsApi)
	if err != nil {
		resp.Diagnostics.AddError(errorRead, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModelPluralDS(projectID, name, flexRestoreJobs))...)
}

func ListFlexRestoreJobs(ctx context.Context, projectID, name string, client admin.FlexRestoreJobsApi) (*[]admin.FlexBackupRestoreJob20241113, error) {
	params := admin.ListFlexRestoreJobsApiParams{
		GroupId: projectID,
		Name:    name,
	}
	flexRestoreJobs, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.FlexBackupRestoreJob20241113], *http.Response, error) {
		request := client.ListFlexRestoreJobsWithParams(ctx, &params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		return nil, err
	}
	return &flexRestoreJobs, nil
}
