package pushbasedlogexport

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

var _ datasource.DataSource = &pushBasedLogExportDS{}
var _ datasource.DataSourceWithConfigure = &pushBasedLogExportDS{}

func DataSource() datasource.DataSource {
	return &pushBasedLogExportDS{
		DSCommon: config.DSCommon{
			DataSourceName: pushBasedLogExportName,
		},
	}
}

type pushBasedLogExportDS struct {
	config.DSCommon
}

func (d *pushBasedLogExportDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.DataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.DataSourceSchemaRequest{
		RequiredFields: []string{"project_id"},
	})
}

func (d *pushBasedLogExportDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfConfig TFPushBasedLogExportDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := tfConfig.ProjectID.ValueString()
	logConfig, _, err := connV2.PushBasedLogExportApi.GetPushBasedLogConfiguration(ctx, projectID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error when getting push-based log export configuration", err.Error())
		return
	}

	newTFModel, diags := NewTFPushBasedLogExport(ctx, projectID, logConfig, nil, nil)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	dsModel := convertToDSModel(newTFModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, dsModel)...)
}

func convertToDSModel(inputModel *TFPushBasedLogExportRSModel) TFPushBasedLogExportDSModel {
	return TFPushBasedLogExportDSModel{
		BucketName: inputModel.BucketName,
		CreateDate: inputModel.CreateDate,
		ProjectID:  inputModel.ProjectID,
		IamRoleID:  inputModel.IamRoleID,
		PrefixPath: inputModel.PrefixPath,
		State:      inputModel.State,
	}
}
