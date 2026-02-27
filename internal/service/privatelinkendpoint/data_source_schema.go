package privatelinkendpoint

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

func PluralDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"provider_name": schema.StringAttribute{
				Required: true,
			},
			"results": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"private_link_id": schema.StringAttribute{
							Computed: true,
						},
						"endpoint_service_name": schema.StringAttribute{
							Computed: true,
						},
						"error_message": schema.StringAttribute{
							Computed: true,
						},
						"interface_endpoints": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"private_endpoints": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"private_link_service_name": schema.StringAttribute{
							Computed: true,
						},
						"private_link_service_resource_id": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"endpoint_group_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"region_name": schema.StringAttribute{
							Computed: true,
						},
						"service_attachment_names": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
						"port_mapping_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Flag that indicates whether this resource uses GCP port-mapping. When `true`, it uses the port-mapped architecture. When `false` or unset, it uses the GCP legacy private endpoint architecture. Only applicable for GCP provider.",
						},
					},
				},
			},
		},
	}
}

type TFPrivateLinkEndpointsModel struct {
	ProjectID    types.String                   `tfsdk:"project_id"`
	ProviderName types.String                   `tfsdk:"provider_name"`
	Results      []TFPrivateLinkEndpointDSModel `tfsdk:"results"`
}

type TFPrivateLinkEndpointDSModel struct {
	PrivateLinkID                types.String `tfsdk:"private_link_id"`
	EndpointServiceName          types.String `tfsdk:"endpoint_service_name"`
	ErrorMessage                 types.String `tfsdk:"error_message"`
	InterfaceEndpoints           types.List   `tfsdk:"interface_endpoints"`
	PrivateEndpoints             types.List   `tfsdk:"private_endpoints"`
	PrivateLinkServiceName       types.String `tfsdk:"private_link_service_name"`
	PrivateLinkServiceResourceID types.String `tfsdk:"private_link_service_resource_id"`
	Status                       types.String `tfsdk:"status"`
	EndpointGroupNames           types.List   `tfsdk:"endpoint_group_names"`
	RegionName                   types.String `tfsdk:"region_name"`
	ServiceAttachmentNames       types.List   `tfsdk:"service_attachment_names"`
	PortMappingEnabled           types.Bool   `tfsdk:"port_mapping_enabled"`
}

func newTFPrivateLinkEndpointResults(ctx context.Context, endpoints []admin.EndpointService) ([]TFPrivateLinkEndpointDSModel, diag.Diagnostics) {
	results := make([]TFPrivateLinkEndpointDSModel, len(endpoints))
	var diags diag.Diagnostics
	for i := range endpoints {
		model, d := newTFPrivateLinkEndpointDSModel(ctx, &endpoints[i])
		diags.Append(d...)
		results[i] = model
	}
	return results, diags
}

func newTFPrivateLinkEndpointDSModel(ctx context.Context, endpoint *admin.EndpointService) (TFPrivateLinkEndpointDSModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	interfaceEndpoints, d := types.ListValueFrom(ctx, types.StringType, endpoint.GetInterfaceEndpoints())
	diags.Append(d...)
	privateEndpoints, d := types.ListValueFrom(ctx, types.StringType, endpoint.GetPrivateEndpoints())
	diags.Append(d...)
	endpointGroupNames, d := types.ListValueFrom(ctx, types.StringType, endpoint.GetEndpointGroupNames())
	diags.Append(d...)
	serviceAttachmentNames, d := types.ListValueFrom(ctx, types.StringType, endpoint.GetServiceAttachmentNames())
	diags.Append(d...)

	return TFPrivateLinkEndpointDSModel{
		PrivateLinkID:                types.StringValue(endpoint.GetId()),
		EndpointServiceName:          types.StringValue(endpoint.GetEndpointServiceName()),
		ErrorMessage:                 types.StringValue(endpoint.GetErrorMessage()),
		InterfaceEndpoints:           interfaceEndpoints,
		PrivateEndpoints:             privateEndpoints,
		PrivateLinkServiceName:       types.StringValue(endpoint.GetPrivateLinkServiceName()),
		PrivateLinkServiceResourceID: types.StringValue(endpoint.GetPrivateLinkServiceResourceId()),
		Status:                       types.StringValue(endpoint.GetStatus()),
		EndpointGroupNames:           endpointGroupNames,
		RegionName:                   types.StringValue(endpoint.GetRegionName()),
		ServiceAttachmentNames:       serviceAttachmentNames,
		PortMappingEnabled:           types.BoolValue(endpoint.GetPortMappingEnabled()),
	}, diags
}
