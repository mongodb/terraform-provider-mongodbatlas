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
				MarkdownDescription: "Domain name of Privatelink connected cluster.",
			},
			"dns_sub_domain": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Sub-Domain name of Confluent cluster. These are typically your availability zones.",
				ElementType:         types.StringType,
			},
			"error_message": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Error message if the connection is in a failed state.",
			},
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.\n\n**NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group or project id remains the same. The resource and corresponding endpoints use the term groups.",
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
				MarkdownDescription: "Provider where the Kafka cluster is deployed.",
			},
			"region": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "When the vendor is `CONFLUENT`, this is the domain name of Confluent cluster. When the vendor is `MSK`, this is computed by the API from the provided `arn`.",
			},
			"service_endpoint_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Service Endpoint ID.",
			},
			"state": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the connection.",
			},
			"vendor": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Vendor who manages the Kafka cluster. Possible values are `CONFLUENT`, `MSK` or `GENERIC`.",
			},
			"arn": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Amazon Resource Name (ARN).",
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
