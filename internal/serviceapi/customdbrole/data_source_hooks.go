package customdbrole

import (
	"context"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/autogen"
)

// Manual-parity custom code for the mongodbatlas_custom_db_role data source.

type TFDSExpandedModel struct {
	// Crafted id of the manual data source, injected into the response body by the hook below.
	ID types.String `tfsdk:"id" apiname:"tfStateId" autogen:"omitjson"`
}

var _ autogen.DataSourceSchemaHook = (*ds)(nil)
var _ autogen.PostReadAPICallHook = (*ds)(nil)

func (d *ds) DataSourceSchema(_ context.Context, baseSchema datasourceschema.Schema) datasourceschema.Schema {
	baseSchema.Attributes["id"] = datasourceschema.StringAttribute{
		Computed: true,
	}
	return baseSchema
}

func (d *ds) PostReadAPICall(req autogen.HandleReadReq, result autogen.APICallResult) autogen.APICallResult {
	projectID := ""
	if model, ok := req.State.(*TFDSModel); ok {
		projectID = model.ProjectId.ValueString()
	}
	return adaptResponse(result, projectID)
}
