package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasPrivateEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasPrivateEndpointrRead,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		DeprecationMessage: "this data source is deprecated, please transition as soon as possible to mongodbatlas_privatelink_endpoint",
	}
}

func dataSourceMongoDBAtlasPrivateEndpointrRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	privateLinkID := d.Get("private_link_id").(string)

	privateEndpoint, _, err := conn.PrivateEndpointsDeprecated.Get(context.Background(), projectID, privateLinkID)
	if err != nil {
		return fmt.Errorf(errorPrivateEndpointsRead, privateLinkID, err)
	}

	if err := d.Set("private_link_id", privateEndpoint.ID); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "private_link_id", privateLinkID, err)
	}

	if err := d.Set("endpoint_service_name", privateEndpoint.EndpointServiceName); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "endpoint_service_name", privateLinkID, err)
	}

	if err := d.Set("error_message", privateEndpoint.ErrorMessage); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "error_message", privateLinkID, err)
	}

	if err := d.Set("interface_endpoints", privateEndpoint.InterfaceEndpoints); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "interface_endpoints", privateLinkID, err)
	}

	if err := d.Set("status", privateEndpoint.Status); err != nil {
		return fmt.Errorf(errorPrivateEndpointsSetting, "status", privateLinkID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"private_link_id": privateEndpoint.ID,
		"project_id":      projectID,
	}))

	return nil
}
