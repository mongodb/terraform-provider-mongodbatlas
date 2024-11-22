package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchemaDelete(ctx context.Context) schema.Schema {
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
			Computed: true,
			MarkdownDescription: "Stream aggregation pipeline you want to apply to your streaming data. [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/atlas-stream-processing/stream-aggregation/#std-label-stream-aggregation)" +
				" contain more information. Using [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) is recommended when settig this attribute. For more details see the [Aggregation Pipelines Documentation](https://www.mongodb.com/docs/atlas/atlas-stream-processing/stream-aggregation/)",
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
			Computed: true,
			MarkdownDescription: "The state of the stream processor. Commonly occurring states are 'CREATED', 'STARTED', 'STOPPED' and 'FAILED'. Used to start or stop the Stream Processor. Valid values are `CREATED`, `STARTED` or `STOPPED`." +
				" When a Stream Processor is created without specifying the state, it will default to `CREATED` state.\n\n**NOTE** When creating a stream processor, setting the state to STARTED can automatically start the stream processor.",
		},
		"stats": schema.StringAttribute{
			Computed:            true,
			MarkdownDescription: "The stats associated with the stream processor. Refer to the [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/atlas-stream-processing/manage-stream-processor/#view-statistics-of-a-stream-processor) for more information.",
		},
		"options": schema.SingleNestedAttribute{
			Computed:            true,
			MarkdownDescription: "Optional configuration for the stream processor.",
			Attributes: map[string]schema.Attribute{
				"dlq": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"coll": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the collection to use for the DLQ.",
						},
						"connection_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the connection to write DLQ messages to. Must be an Atlas connection.",
						},
						"db": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Name of the database to use for the DLQ.",
						},
					},
					Computed:            true,
					MarkdownDescription: "Dead letter queue for the stream processor. Refer to the [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/reference/glossary/#std-term-dead-letter-queue) for more information.",
				},
			},
		},
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
