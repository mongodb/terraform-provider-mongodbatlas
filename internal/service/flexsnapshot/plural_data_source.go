package flexsnapshot

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

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
	var tfModel TFFlexSnapshotsDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	projectID := tfModel.ProjectId.ValueString()
	name := tfModel.Name.ValueString()
	connV2 := d.Client.AtlasV2
	flexSnapshots, err := ListFlexSnapshots(ctx, projectID, name, connV2.FlexSnapshotsApi)
	if err != nil {
		resp.Diagnostics.AddError(errorRead, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModelPluralDS(projectID, name, flexSnapshots))...)
}

func ListFlexSnapshots(ctx context.Context, projectID, name string, client admin.FlexSnapshotsApi) (*[]admin.FlexBackupSnapshot20241113, error) {
	params := admin.ListFlexBackupsApiParams{
		GroupId: projectID,
		Name:    name,
	}
	flexClusters, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.FlexBackupSnapshot20241113], *http.Response, error) {
		request := client.ListFlexBackupsWithParams(ctx, &params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		return nil, err
	}
	return &flexClusters, nil
}
