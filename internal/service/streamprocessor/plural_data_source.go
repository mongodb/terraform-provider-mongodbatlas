package streamprocessor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
)

var _ datasource.DataSource = &StreamProccesorDS{}
var _ datasource.DataSourceWithConfigure = &StreamProccesorDS{}

func PluralDataSource() datasource.DataSource {
	return &streamProcessorsDS{
		DSCommon: config.DSCommon{
			DataSourceName: fmt.Sprintf("%ss", StreamProcessorName),
		},
	}
}

type streamProcessorsDS struct {
	config.DSCommon
}

func (d *streamProcessorsDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = conversion.PluralDataSourceSchemaFromResource(ResourceSchema(ctx), &conversion.PluralDataSourceSchemaRequest{
		RequiredFields:     []string{"project_id", "instance_name"},
		OverrideResultsDoc: "Returns all Stream Processors within the specified stream instance.\n\nTo use this resource, the requesting API Key must have the Project Owner\n\nrole or Project Stream Processing Owner role.",
	})
}

func (d *streamProcessorsDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var streamConnectionsConfig TFStreamProcessorsDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &streamConnectionsConfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connV2 := d.Client.AtlasV2
	projectID := streamConnectionsConfig.ProjectID.ValueString()
	instanceName := streamConnectionsConfig.InstanceName.ValueString()

	params := admin.GetStreamProcessorsApiParams{
		GroupId:    projectID,
		TenantName: instanceName,
	}
	sdkProcessors, err := dsschema.AllPages(ctx, func(ctx context.Context, pageNum int) (dsschema.PaginateResponse[admin.StreamsProcessorWithStats], *http.Response, error) {
		request := connV2.StreamsApi.GetStreamProcessorsWithParams(ctx, &params)
		request = request.PageNum(pageNum)
		return request.Execute()
	})
	if err != nil {
		resp.Diagnostics.AddError("error fetching results", err.Error())
		return
	}

	newStreamConnectionsModel, diags := NewTFStreamProcessors(ctx, &streamConnectionsConfig, sdkProcessors)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, newStreamConnectionsModel)...)
}
