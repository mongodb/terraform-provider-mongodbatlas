package streamprivatelinkendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the Private Link connection.",
			},
			"dns_domain": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The domain hostname. Required for the following provider and vendor combinations:<br>- AWS provider with CONFLUENT vendor.<br>- AZURE provider with EVENTHUB or CONFLUENT vendor.",
			},
			"dns_sub_domain": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Sub-Domain name of Confluent cluster. These are typically your availability zones. Required for AWS Provider and CONFLUENT vendor. If your AWS CONFLUENT cluster doesn't use subdomains, you must set this to the empty array [].",
				ElementType:         types.StringType,
			},
			"error_message": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Error message if the connection is in a failed state.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.<br>**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group or project id remains the same. The resource and corresponding endpoints use the term groups.",
			},
			"interface_endpoint_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Interface endpoint ID that is created from the specified service endpoint ID.",
			},
			"interface_endpoint_name": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of interface endpoint that is created from the specified service endpoint ID.",
			},
			"provider_account_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Account ID from the cloud provider.",
			},
			"provider_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Provider where the Kafka cluster is deployed. Valid values are AWS and AZURE.",
			},
			"region": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The region of the Provider’s cluster. See [AZURE](https://www.mongodb.com/docs/atlas/reference/microsoft-azure/#stream-processing-instances) and [AWS](https://www.mongodb.com/docs/atlas/reference/amazon-aws/#stream-processing-instances) supported regions. When the vendor is `CONFLUENT`, this is the domain name of Confluent cluster. When the vendor is `MSK`, this is computed by the API from the provided `arn`.",
			},
			"service_endpoint_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "For AZURE EVENTHUB, this is the [namespace endpoint ID](https://learn.microsoft.com/en-us/rest/api/eventhub/namespaces/get). For AWS CONFLUENT cluster, this is the [VPC Endpoint service name](https://docs.confluent.io/cloud/current/networking/private-links/aws-privatelink.html).",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the connection.",
			},
			"vendor": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Vendor that manages the Kafka cluster. The following are the vendor values per provider:<br>- MSK and CONFLUENT for the AWS provider.<br>- EVENTHUB and CONFLUENT for the AZURE provider.",
			},
			"arn": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Amazon Resource Name (ARN). Required for AWS Provider and MSK vendor.",
			},
		},
	}
}

type TFModel struct {
	Id                    types.String `tfsdk:"id"`
	DnsDomain             types.String `tfsdk:"dns_domain"`
	DnsSubDomain          types.List   `tfsdk:"dns_sub_domain"`
	ErrorMessage          types.String `tfsdk:"error_message"`
	ProjectId             types.String `tfsdk:"project_id"`
	InterfaceEndpointId   types.String `tfsdk:"interface_endpoint_id"`
	InterfaceEndpointName types.String `tfsdk:"interface_endpoint_name"`
	Provider              types.String `tfsdk:"provider_name"`
	ProviderAccountId     types.String `tfsdk:"provider_account_id"`
	Region                types.String `tfsdk:"region"`
	ServiceEndpointId     types.String `tfsdk:"service_endpoint_id"`
	State                 types.String `tfsdk:"state"`
	Vendor                types.String `tfsdk:"vendor"`
	Arn                   types.String `tfsdk:"arn"`
}

type TFModelDSP struct {
	ProjectId types.String `tfsdk:"project_id"`
	Results   []TFModel    `tfsdk:"results"`
}
