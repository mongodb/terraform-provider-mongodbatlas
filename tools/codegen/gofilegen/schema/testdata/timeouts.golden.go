package testname

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"string_attr": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "string description",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

type TFModel struct {
	StringAttr types.String   `tfsdk:"string_attr"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}
