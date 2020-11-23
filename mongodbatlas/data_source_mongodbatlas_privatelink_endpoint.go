package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasPrivateLinkEndpointRead,
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
		},
	}
}

func dataSourceMongoDBAtlasPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	privateLinkID := d.Get("private_link_id").(string)
	providerName := d.Get("provider_name").(string)

	privateEndpoint, _, err := conn.PrivateEndpoints.Get(context.Background(), projectID, providerName, privateLinkID)
	if err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsRead, privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateEndpoint.ID); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.EndpointServiceName); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "endpoint_service_name", privateLinkID, err)
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "error_message", privateLinkID, err)
	}

	if err := d.Set("interface_endpoints", privateEndpoint.InterfaceEndpoints); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "interface_endpoints", privateLinkID, err)
	}

	if err := d.Set("private_endpoints", privateEndpoint.PrivateEndpoints); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_endpoints", privateLinkID, err)
	}

	if err := d.Set("private_link_service_name", privateEndpoint.PrivateLinkServiceName); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_service_name", privateLinkID, err)
	}

	if err := d.Set("private_link_service_resource_id", privateEndpoint.PrivateLinkServiceResourceID); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "private_link_service_resource_id", privateLinkID, err)
	}

	if err := d.Set("status", privateEndpoint.Status); err != nil {
		return fmt.Errorf(errorPrivateLinkEndpointsSetting, "status", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpoint.ID,
		"project_id":      projectID,
		"provider_name":   providerName,
	}))

	return nil
}
