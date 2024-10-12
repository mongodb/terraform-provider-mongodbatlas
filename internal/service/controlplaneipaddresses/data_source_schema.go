package controlplaneipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"inbound": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"aws": schema.MapAttribute{
						ElementType: types.ListType{
							ElemType: types.StringType,
						},
						Computed:            true,
						MarkdownDescription: "Control plane IP addresses in AWS. Each key identifies an Amazon Web Services (AWS) region. Each value identifies control plane IP addresses in the AWS region.",
					},
					"azure": schema.MapAttribute{
						ElementType: types.ListType{
							ElemType: types.StringType,
						},
						Computed:            true,
						MarkdownDescription: "Control plane IP addresses in Azure. Each key identifies an Azure region. Each value identifies control plane IP addresses in the Azure region.",
					},
					"gcp": schema.MapAttribute{
						ElementType: types.ListType{
							ElemType: types.StringType,
						},
						Computed:            true,
						MarkdownDescription: "Control plane IP addresses in GCP. Each key identifies a Google Cloud (GCP) region. Each value identifies control plane IP addresses in the GCP region.",
					},
				},
				Computed:            true,
				MarkdownDescription: "List of inbound IP addresses to the Atlas control plane, categorized by cloud provider. If your application allows outbound HTTP requests only to specific IP addresses, you must allow access to the following IP addresses so that your API requests can reach the Atlas control plane.",
			},
			"outbound": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"aws": schema.MapAttribute{
						ElementType: types.ListType{
							ElemType: types.StringType,
						},
						Computed:            true,
						MarkdownDescription: "Control plane IP addresses in AWS. Each key identifies an Amazon Web Services (AWS) region. Each value identifies control plane IP addresses in the AWS region.",
					},
					"azure": schema.MapAttribute{
						ElementType: types.ListType{
							ElemType: types.StringType,
						},
						Computed:            true,
						MarkdownDescription: "Control plane IP addresses in Azure. Each key identifies an Azure region. Each value identifies control plane IP addresses in the Azure region.",
					},
					"gcp": schema.MapAttribute{
						ElementType: types.ListType{
							ElemType: types.StringType,
						},
						Computed:            true,
						MarkdownDescription: "Control plane IP addresses in GCP. Each key identifies a Google Cloud (GCP) region. Each value identifies control plane IP addresses in the GCP region.",
					},
				},
				Computed:            true,
				MarkdownDescription: "List of outbound IP addresses from the Atlas control plane, categorized by cloud provider. If your network allows inbound HTTP requests only from specific IP addresses, you must allow access from the following IP addresses so that Atlas can communicate with your webhooks and KMS.",
			},
		},
	}
}

type TFControlPlaneIpAddressesModel struct {
	Inbound  InboundValue  `tfsdk:"inbound"`
	Outbound OutboundValue `tfsdk:"outbound"`
}

type InboundValue struct {
	Aws   basetypes.MapValue `tfsdk:"aws"`
	Azure basetypes.MapValue `tfsdk:"azure"`
	Gcp   basetypes.MapValue `tfsdk:"gcp"`
}
type OutboundValue struct {
	Aws   basetypes.MapValue `tfsdk:"aws"`
	Azure basetypes.MapValue `tfsdk:"azure"`
	Gcp   basetypes.MapValue `tfsdk:"gcp"`
}
