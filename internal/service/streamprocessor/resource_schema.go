package streamprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"change_stream_token": schema.StringAttribute{
				Computed:            true,
				Description:         "The resume token for the change stream. Only used when the pipeline source is Cluster.",
				MarkdownDescription: "The resume token for the change stream. Only used when the pipeline source is Cluster.",
			},
			"id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Unique 24-hexadecimal character string that identifies the stream processor.",
				MarkdownDescription: "Unique 24-hexadecimal character string that identifies the stream processor.",
			},
			"instance_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Human-readable label that identifies the stream instance.",
				MarkdownDescription: "Human-readable label that identifies the stream instance.",
			},
			"pipeline": schema.StringAttribute{
				Validators: []validator.String{
					validate.StringIsJSON(),
				},
				PlanModifiers: []planmodifier.String{
					schemafunc.DiffSuppressJSON(),
				},
				Required:            true,
				Description:         "Stream aggregation pipeline you want to apply to your streaming data.",
				MarkdownDescription: "Stream aggregation pipeline you want to apply to your streaming data.",
			},
			"processor_name": schema.StringAttribute{
				Required:            true,
				Description:         "Human-readable label that identifies the stream processor.",
				MarkdownDescription: "Human-readable label that identifies the stream processor.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"state": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "The state of the stream processor.",
				MarkdownDescription: "The state of the stream processor.",
			},
			"options": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"dlq": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"coll": schema.StringAttribute{
								Required:            true,
								Description:         "Name of the collection that will be used for the DLQ.",
								MarkdownDescription: "Name of the collection that will be used for the DLQ.",
							},
							"connection_name": schema.StringAttribute{
								Required:            true,
								Description:         "Connection name that will be used to write DLQ messages to. Has to be an Atlas connection.",
								MarkdownDescription: "Connection name that will be used to write DLQ messages to. Has to be an Atlas connection.",
							},
							"db": schema.StringAttribute{
								Required:            true,
								Description:         "Name of the database that will be used for the DLQ.",
								MarkdownDescription: "Name of the database that will be used for the DLQ.",
							},
						},
						Optional:            true,
						Computed:            true,
						Description:         "Dead letter queue for the stream processor.",
						MarkdownDescription: "Dead letter queue for the stream processor.",
					},
				},
				Optional:            true,
				Computed:            true,
				Description:         "Optional configuration for the stream processor.",
				MarkdownDescription: "Optional configuration for the stream processor.",
			},
		},
	}
}

type TFStreamProcessorRSModel struct {
	InstanceName  types.String `tfsdk:"instance_name"`
	Options       types.Object `tfsdk:"options"`
	Pipeline      types.String `tfsdk:"pipeline"`
	ProcessorID   types.String `tfsdk:"processor_id"`
	ProcessorName types.String `tfsdk:"processor_name"`
	ProjectID     types.String `tfsdk:"project_id"`
	State         types.String `tfsdk:"state"`
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
