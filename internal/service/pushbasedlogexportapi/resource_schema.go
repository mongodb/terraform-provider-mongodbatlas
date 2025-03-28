package pushbasedlogexportapi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
				MarkdownDescription: "Date and time that this feature was enabled on.",
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"iam_role_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the AWS IAM role that will be used to write to the S3 bucket.",
			},
			"links": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of one or more Uniform Resource Locators (URLs) that point to API sub-resources, related API resources, or both. RFC 5988 outlines these relationships.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"href": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Uniform Resource Locator (URL) that points another API resource to which this response has some relationship. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
						"rel": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Uniform Resource Locator (URL) that defines the semantic relationship between this resource and another API resource. This URL often begins with `https://cloud.mongodb.com/api/atlas`.",
						},
					},
				},
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
	CreateDate types.String `tfsdk:"create_date"`
	GroupId    types.String `tfsdk:"group_id"`
	IamRoleId  types.String `tfsdk:"iam_role_id"`
	Links      types.List   `tfsdk:"links"`
	PrefixPath types.String `tfsdk:"prefix_path"`
	State      types.String `tfsdk:"state"`
}
type TFLinksModel struct {
	Href types.String `tfsdk:"href"`
	Rel  types.String `tfsdk:"rel"`
}

var LinksObjType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"href": types.StringType,
	"rel":  types.StringType,
}}
