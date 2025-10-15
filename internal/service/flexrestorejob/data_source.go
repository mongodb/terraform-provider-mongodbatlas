package flexrestorejob

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const (
	resourceName = "flex_restore_job"
	errorRead    = "error reading flex restore job"
)

var _ datasource.DataSource = &ds{}
var _ datasource.DataSourceWithConfigure = &ds{}

func DataSource() datasource.DataSource {
	return &ds{
		DSCommon: config.DSCommon{
			DataSourceName: resourceName,
		},
	}
}

type ds struct {
	config.DSCommon
}

func (d *ds) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(DataSourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id", "name", "restore_job_id"},
	})
}

func (d *ds) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfModel TFModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
	apiResp, _, err := d.Client.AtlasV2.FlexRestoreJobsApi.GetFlexRestoreJob(ctx, tfModel.ProjectID.ValueString(), tfModel.Name.ValueString(), tfModel.RestoreJobID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(errorRead, err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, NewTFModel(apiResp))...)
}
