package projectipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Unique 24-hexadecimal digit string that identifies your project.",
			},
			"services": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"clusters": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"cluster_name": schema.StringAttribute{
									Computed:            true,
									MarkdownDescription: "Human-readable label that identifies the cluster.",
								},
								"inbound": schema.ListAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									MarkdownDescription: "List of inbound IP addresses associated with the cluster. If your network allows outbound HTTP requests only to specific IP addresses, you must allow access to the following IP addresses so that your application can connect to your Atlas cluster.",
								},
								"outbound": schema.ListAttribute{
									ElementType:         types.StringType,
									Computed:            true,
									MarkdownDescription: "List of outbound IP addresses associated with the cluster. If your network allows inbound HTTP requests only from specific IP addresses, you must allow access from the following IP addresses so that your Atlas cluster can communicate with your webhooks and KMS.",
								},
							},
						},
						Computed:            true,
						MarkdownDescription: "IP addresses of clusters.",
					},
				},
				Computed:            true,
				MarkdownDescription: "List of IP addresses in a project categorized by services.",
			},
		},
	}
}

type TFProjectIpAddressesModel struct {
	ProjectId types.String `tfsdk:"project_id"`
	Services  types.Object `tfsdk:"services"`
}

type TFServicesModel struct {
	Clusters []TFClusterValueModel `tfsdk:"clusters"`
}

type TFClusterValueModel struct {
	ClusterName types.String `tfsdk:"cluster_name"`
	Inbound     types.List   `tfsdk:"inbound"`
	Outbound    types.List   `tfsdk:"outbound"`
}

var IPAddressesObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"project_id": types.StringType,
	"services":   ServicesObjectType,
}}

var ServicesObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"clusters": types.ListType{ElemType: ClusterIPsObjectType},
}}

var ClusterIPsObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"cluster_name": types.StringType,
	"inbound":      types.ListType{ElemType: types.StringType},
	"outbound":     types.ListType{ElemType: types.StringType},
}}
