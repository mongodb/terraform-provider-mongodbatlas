package streamaccountdetails

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"aws_account_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The AWS Account ID.",
				MarkdownDescription: "The AWS Account ID.",
			},
			"azure_subscription_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The Azure Subscription ID.",
				MarkdownDescription: "The Azure Subscription ID.",
			},
			"cidr_block": schema.StringAttribute{
				Computed:            true,
				Description:         "The VPC CIDR Block.",
				MarkdownDescription: "The VPC CIDR Block.",
			},
			"cloud_provider": schema.StringAttribute{
				Required:            true,
				Description:         "One of \"aws\", \"azure\" or \"gcp\".",
				MarkdownDescription: "One of \"aws\", \"azure\" or \"gcp\".",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				Description:         "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^([a-f0-9]{24})$"), ""),
				},
			},
			"region_name": schema.StringAttribute{
				Required:            true,
				Description:         "The cloud provider specific region name, i.e. \"US_EAST_1\" for cloud provider \"aws\".",
				MarkdownDescription: "The cloud provider specific region name, i.e. \"US_EAST_1\" for cloud provider \"aws\".",
			},
			"virtual_network_name": schema.StringAttribute{
				Computed:            true,
				Description:         "The name of the virtual network.",
				MarkdownDescription: "The name of the virtual network.",
			},
			"vpc_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The VPC ID.",
				MarkdownDescription: "The VPC ID.",
			},
		},
	}
}

type TFStreamAccountDetailsModel struct {
	AwsAccountId        types.String `tfsdk:"aws_account_id"`
	AzureSubscriptionId types.String `tfsdk:"azure_subscription_id"`
	CidrBlock           types.String `tfsdk:"cidr_block"`
	CloudProvider       types.String `tfsdk:"cloud_provider"`
	ProjectId           types.String `tfsdk:"project_id"`
	RegionName          types.String `tfsdk:"region_name"`
	VirtualNetworkName  types.String `tfsdk:"virtual_network_name"`
	VpcId               types.String `tfsdk:"vpc_id"`
}
