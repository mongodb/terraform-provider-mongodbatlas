package privatelinkendpoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_link_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"interface_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"private_link_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_link_service_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_group_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_attachment_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	connV2 := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	privateLinkID := conversion.GetEncodedID(d.Get("private_link_id").(string), "private_link_id")
	providerName := d.Get("provider_name").(string)

	privateEndpoint, _, err := connV2.PrivateEndpointServicesApi.GetPrivateEndpointService(ctx, projectID, providerName, privateLinkID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsRead, privateLinkID, err))
	}

	if err := d.Set("private_link_id", privateEndpoint.GetId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "private_link_id", privateLinkID, err))
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.GetEndpointServiceName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "endpoint_service_name", privateLinkID, err))
	}

	if err := d.Set("error_message", privateEndpoint.GetErrorMessage()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "error_message", privateLinkID, err))
	}

	if err := d.Set("interface_endpoints", privateEndpoint.GetInterfaceEndpoints()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_endpoints", privateEndpoint.GetPrivateEndpoints()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "private_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_link_service_name", privateEndpoint.GetPrivateLinkServiceName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "private_link_service_name", privateLinkID, err))
	}

	if err := d.Set("private_link_service_resource_id", privateEndpoint.GetPrivateLinkServiceResourceId()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "private_link_service_resource_id", privateLinkID, err))
	}

	if err := d.Set("status", privateEndpoint.GetStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "status", privateLinkID, err))
	}

	if err := d.Set("endpoint_group_names", privateEndpoint.GetEndpointGroupNames()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "endpoint_group_names", privateLinkID, err))
	}

	if err := d.Set("region_name", privateEndpoint.GetRegionName()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "region_name", privateLinkID, err))
	}

	if err := d.Set("service_attachment_names", privateEndpoint.GetServiceAttachmentNames()); err != nil {
		return diag.FromErr(fmt.Errorf(ErrorPrivateLinkEndpointsSetting, "service_attachment_names", privateLinkID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"private_link_id": privateEndpoint.GetId(),
		"project_id":      projectID,
		"provider_name":   providerName,
	}))

	return nil
}
