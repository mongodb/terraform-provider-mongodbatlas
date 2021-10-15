package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivateLinkEndpointRead,
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

func dataSourceMongoDBAtlasPrivateLinkEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	privateLinkID := getEncodedID(d.Get("private_link_id").(string), "private_link_id")
	providerName := d.Get("provider_name").(string)

	privateEndpoint, _, err := conn.PrivateEndpoints.Get(ctx, projectID, providerName, privateLinkID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsRead, privateLinkID, err))
	}

	if err := d.Set("private_link_id", privateEndpoint.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_id", privateLinkID, err))
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.EndpointServiceName); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "endpoint_service_name", privateLinkID, err))
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "error_message", privateLinkID, err))
	}

	if err := d.Set("interface_endpoints", privateEndpoint.InterfaceEndpoints); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_endpoints", privateEndpoint.PrivateEndpoints); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_endpoints", privateLinkID, err))
	}

	if err := d.Set("private_link_service_name", privateEndpoint.PrivateLinkServiceName); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_service_name", privateLinkID, err))
	}

	if err := d.Set("private_link_service_resource_id", privateEndpoint.PrivateLinkServiceResourceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_service_resource_id", privateLinkID, err))
	}

	if err := d.Set("status", privateEndpoint.Status); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "status", privateLinkID, err))
	}

	if err := d.Set("endpoint_group_names", privateEndpoint.EndpointGroupNames); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "endpoint_group_names", privateLinkID, err))
	}

	if err := d.Set("region_name", privateEndpoint.RegionName); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "region_name", privateLinkID, err))
	}

	if err := d.Set("service_attachment_names", privateEndpoint.ServiceAttachmentNames); err != nil {
		return diag.FromErr(fmt.Errorf(errorPrivateLinkEndpointsSetting, "service_attachment_names", privateLinkID, err))
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpoint.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
	}))

	return nil
}
