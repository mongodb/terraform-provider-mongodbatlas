package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: DSAttributes(true),
	}
}

func DSAttributes(withArguments bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Unique 24-hexadecimal character string that identifies the stream processor.",
		},
		"instance_name": schema.StringAttribute{
			Required:            withArguments,
			Computed:            !withArguments,
			MarkdownDescription: "Human-readable label that identifies the stream instance.",
		},
		"pipeline": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "Stream aggregation pipeline you want to apply to your streaming data.",
		},
		"processor_name": schema.StringAttribute{
			Required:            withArguments,
			Computed:            !withArguments,
			MarkdownDescription: "Human-readable label that identifies the stream processor.",
		},
		"project_id": schema.StringAttribute{
			Required:            withArguments,
			Computed:            !withArguments,
			MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
		},
		"state": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The state of the stream processor.",
		},
		"stats": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The stats associated with the stream processor.",
		},
		"options": optionsSchema(true),
	}
}

type TFStreamProcessorDSModel struct {
	ID            types.String `tfsdk:"id"`
	InstanceName  types.String `tfsdk:"instance_name"`
	Options       types.Object `tfsdk:"options"`
	Pipeline      types.String `tfsdk:"pipeline"`
	ProcessorName types.String `tfsdk:"processor_name"`
	ProjectID     types.String `tfsdk:"project_id"`
	State         types.String `tfsdk:"state"`
	Stats         types.String `tfsdk:"stats"`
}
