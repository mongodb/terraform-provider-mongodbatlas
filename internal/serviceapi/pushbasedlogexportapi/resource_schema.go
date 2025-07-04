// Code generated by terraform-provider-mongodbatlas using `make generate-resource`. DO NOT EDIT.

package pushbasedlogexportapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bucket_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the bucket to which the agent will send the logs to.",
			},
			"create_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Date and time that this feature was enabled on. This parameter expresses its value in the ISO 8601 timestamp format in UTC.",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"iam_role_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the AWS IAM role that will be used to write to the S3 bucket.",
			},
			"prefix_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "S3 directory in which vector will write to in order to store the logs. An empty string denotes the root directory.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Describes whether or not the feature is enabled and what status it is in.",
			},
		},
	}
}

type TFModel struct {
	BucketName types.String `tfsdk:"bucket_name"`
	CreateDate types.String `tfsdk:"create_date" autogen:"omitjson"`
	GroupId    types.String `tfsdk:"group_id" autogen:"omitjson"`
	IamRoleId  types.String `tfsdk:"iam_role_id"`
	PrefixPath types.String `tfsdk:"prefix_path"`
	State      types.String `tfsdk:"state" autogen:"omitjson"`
}
