package pushbasedlogexport

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the bucket to which the agent sends the logs to.",
				MarkdownDescription: "The name of the bucket to which the agent sends the logs to.",
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
				Validators: []validator.String{
					stringvalidator.LengthBetween(24, 24),
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
				},
			},
			"iam_role_id": schema.StringAttribute{
				Computed:            true,
				Description:         "ID of the AWS IAM role that is used to write to the S3 bucket.",
				MarkdownDescription: "ID of the AWS IAM role that is used to write to the S3 bucket.",
			},
			"prefix_path": schema.StringAttribute{
				Computed:            true,
				Description:         "S3 directory in which vector writes in order to store the logs. An empty string denotes the root directory.",
				MarkdownDescription: "S3 directory in which vector writes in order to store the logs. An empty string denotes the root directory.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				Description:         "Describes whether or not the feature is enabled and what status it is in.",
				MarkdownDescription: "Describes whether or not the feature is enabled and what status it is in.",
			},
		},
	}
}

type TFPushBasedLogExportDSModel struct {
	BucketName types.String `tfsdk:"bucket_name"`
	CreateDate types.String `tfsdk:"create_date"`
	ProjectID  types.String `tfsdk:"project_id"`
	IamRoleID  types.String `tfsdk:"iam_role_id"`
	PrefixPath types.String `tfsdk:"prefix_path"`
	State      types.String `tfsdk:"state"`
}
