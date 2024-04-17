package pushbasedlogexport

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket_name": schema.StringAttribute{
				Optional:            true,
				Description:         "The name of the bucket to which the agent will send the logs to.",
				MarkdownDescription: "The name of the bucket to which the agent will send the logs to.",
			},
			"create_date": schema.StringAttribute{
				Computed:            true,
				Description:         "Date and time that this feature was enabled on.",
				MarkdownDescription: "Date and time that this feature was enabled on.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(24, 24),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
				},
			},
			"iam_role_id": schema.StringAttribute{
				Optional:            true,
				Description:         "ID of the AWS IAM role that will be used to write to the S3 bucket.",
				MarkdownDescription: "ID of the AWS IAM role that will be used to write to the S3 bucket.",
			},
			"links": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{
							Computed:            true,
							Description:         "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
							MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
						"rel": schema.StringAttribute{
							Computed:            true,
							Description:         "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
							MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
					},
					// CustomType: LinksType{
					// 	ObjectType: types.ObjectType{
					// 		AttrTypes: LinksValue{}.AttributeTypes(ctx),
					// 	},
					// },
				},
				Computed:            true,
				Description:         "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
				MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
			},
			"prefix_path": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				Description:         "S3 directory in which vector will write to in order to store the logs.",
				MarkdownDescription: "S3 directory in which vector will write to in order to store the logs.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "Describes whether or not the feature is enabled and what status it is in.",
				MarkdownDescription: "Describes whether or not the feature is enabled and what status it is in.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
			}),
		},
	}
}

type TFPushBasedLogExportRSModel struct {
	BucketName types.String   `tfsdk:"bucket_name"`
	CreateDate types.String   `tfsdk:"create_date"`
	ProjectID  types.String   `tfsdk:"project_id"`
	IamRoleID  types.String   `tfsdk:"iam_role_id"`
	Links      types.List     `tfsdk:"links"`
	PrefixPath types.String   `tfsdk:"prefix_path"`
	State      types.String   `tfsdk:"state"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
}
