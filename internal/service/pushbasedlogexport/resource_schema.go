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
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/customplanmodifier"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the bucket to which the agent sends the logs to.",
			},
			"create_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time that this feature was enabled on.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
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
				Required:            true,
				MarkdownDescription: "ID of the AWS IAM role that is used to write to the S3 bucket.",
			},
			"prefix_path": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				MarkdownDescription: "S3 directory in which vector writes in order to store the logs. An empty string denotes the root directory.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Describes whether or not the feature is enabled and what status it is in.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create:            true,
				Update:            true,
				Delete:            true,
				CreateDescription: "A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as \"30s\" or \"2h45m\". Valid time units are \"s\" (seconds), \"m\" (minutes), \"h\" (hours). Default: `15m`.",
				UpdateDescription: "A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as \"30s\" or \"2h45m\". Valid time units are \"s\" (seconds), \"m\" (minutes), \"h\" (hours). Default: `15m`.",
				DeleteDescription: "A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as \"30s\" or \"2h45m\". Valid time units are \"s\" (seconds), \"m\" (minutes), \"h\" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs. Default: `15m`.",
			}),
			"delete_on_create_timeout": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					customplanmodifier.CreateOnlyBoolWithDefault(true),
				},
				MarkdownDescription: "Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.",
			},
		},
	}
}

type TFPushBasedLogExportCommonModel struct {
	BucketName types.String `tfsdk:"bucket_name"`
	CreateDate types.String `tfsdk:"create_date"`
	ProjectID  types.String `tfsdk:"project_id"`
	IamRoleID  types.String `tfsdk:"iam_role_id"`
	PrefixPath types.String `tfsdk:"prefix_path"`
	State      types.String `tfsdk:"state"`
}

type TFPushBasedLogExportRSModel struct {
	TFPushBasedLogExportCommonModel
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
	DeleteOnCreateTimeout types.Bool     `tfsdk:"delete_on_create_timeout"`
}
