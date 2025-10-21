package streamprocessor

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the stream processor.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Label that identifies the stream processing workspace.",
				DeprecationMessage:  fmt.Sprintf(constant.DeprecationParamWithReplacement, "workspace_name"),
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("workspace_name"),
					}...),
				},
			},
			"workspace_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Label that identifies the stream processing workspace.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRelative().AtParent().AtName("instance_name"),
					}...),
				},
			},
			"pipeline": schema.StringAttribute{
				CustomType: jsontypes.NormalizedType{},
				Required:   true,
				MarkdownDescription: "Stream aggregation pipeline you want to apply to your streaming data. [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/atlas-stream-processing/stream-aggregation/#std-label-stream-aggregation)" +
					" contain more information. Using [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) is recommended when setting this attribute. For more details see the [Aggregation Pipelines Documentation](https://www.mongodb.com/docs/atlas/atlas-stream-processing/stream-aggregation/)",
			},
			"processor_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Label that identifies the stream processor.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"state": schema.StringAttribute{
				Optional: true,
				Computed: true,
				MarkdownDescription: "The state of the stream processor. Commonly occurring states are 'CREATED', 'STARTED', 'STOPPED' and 'FAILED'. Used to start or stop the Stream Processor. Valid values are `CREATED`, `STARTED` or `STOPPED`." +
					" When a Stream Processor is created without specifying the state, it will default to `CREATED` state. When a Stream Processor is updated without specifying the state, it will default to the Previous state. \n\n**NOTE** When a Stream Processor is updated without specifying the state, it is stopped and then restored to previous state upon update completion.",
			},
			"options": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Optional configuration for the stream processor.",
				Attributes: map[string]schema.Attribute{
					"dlq": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"coll": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Name of the collection to use for the DLQ.",
							},
							"connection_name": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Name of the connection to write DLQ messages to. Must be an Atlas connection.",
							},
							"db": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Name of the database to use for the DLQ.",
							},
						},
						Required:            true,
						MarkdownDescription: "Dead letter queue for the stream processor. Refer to the [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/reference/glossary/#std-term-dead-letter-queue) for more information.",
					},
				},
			},
			"stats": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The stats associated with the stream processor. Refer to the [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/atlas-stream-processing/manage-stream-processor/#view-statistics-of-a-stream-processor) for more information.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
			}),
			"delete_on_create_timeout": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					customplanmodifier.CreateOnly(),
				},
				MarkdownDescription: "Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.",
			},
		},
	}
}

type TFStreamProcessorRSModel struct {
	InstanceName          types.String         `tfsdk:"instance_name"`
	WorkspaceName         types.String         `tfsdk:"workspace_name"`
	Options               types.Object         `tfsdk:"options"`
	Pipeline              jsontypes.Normalized `tfsdk:"pipeline"`
	ProcessorID           types.String         `tfsdk:"id"`
	ProcessorName         types.String         `tfsdk:"processor_name"`
	ProjectID             types.String         `tfsdk:"project_id"`
	State                 types.String         `tfsdk:"state"`
	Stats                 types.String         `tfsdk:"stats"`
	Timeouts              timeouts.Value       `tfsdk:"timeouts"`
	DeleteOnCreateTimeout types.Bool           `tfsdk:"delete_on_create_timeout"`
}

type TFOptionsModel struct {
	Dlq types.Object `tfsdk:"dlq"`
}

type TFDlqModel struct {
	Coll           types.String `tfsdk:"coll"`
	ConnectionName types.String `tfsdk:"connection_name"`
	DB             types.String `tfsdk:"db"`
}

var OptionsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"dlq": DlqObjectType,
}}

var DlqObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"coll":            types.StringType,
	"connection_name": types.StringType,
	"db":              types.StringType,
},
}

type TFStreamProcessorDSModel struct {
	ID            types.String `tfsdk:"id"`
	InstanceName  types.String `tfsdk:"instance_name"`
	WorkspaceName types.String `tfsdk:"workspace_name"`
	Options       types.Object `tfsdk:"options"`
	Pipeline      types.String `tfsdk:"pipeline"`
	ProcessorName types.String `tfsdk:"processor_name"`
	ProjectID     types.String `tfsdk:"project_id"`
	State         types.String `tfsdk:"state"`
	Stats         types.String `tfsdk:"stats"`
}

type TFStreamProcessorsDSModel struct {
	ProjectID     types.String               `tfsdk:"project_id"`
	InstanceName  types.String               `tfsdk:"instance_name"`
	WorkspaceName types.String               `tfsdk:"workspace_name"`
	Results       []TFStreamProcessorDSModel `tfsdk:"results"`
}
