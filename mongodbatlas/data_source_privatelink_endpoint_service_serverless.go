package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSourcePrivateLinkEndpointServerless() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivateEndpointServiceServerlessLinkRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_provider_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_link_service_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_message": {
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

func dataSourceMongoDBAtlasPrivateEndpointServiceServerlessLinkRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*config.MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)
	endpointID := d.Get("endpoint_id").(string)

	serviceEndpoint, _, err := conn.ServerlessPrivateEndpoints.Get(ctx, projectID, instanceName, endpointID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServiceEndpointRead, endpointID, err))
	}

	if err := d.Set("error_message", serviceEndpoint.ErrorMessage); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "error_message", endpointID, err))
	}

	if err := d.Set("status", serviceEndpoint.Status); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "status", endpointID, err))
	}

	if err := d.Set("comment", serviceEndpoint.Comment); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "comment", endpointID, err))
	}

	if err := d.Set("error_message", serviceEndpoint.ErrorMessage); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "error_message", endpointID, err))
	}

	if err := d.Set("endpoint_service_name", serviceEndpoint.EndpointServiceName); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "endpoint_service_name", endpointID, err))
	}

	if err := d.Set("cloud_provider_endpoint_id", serviceEndpoint.CloudProviderEndpointID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "cloud_provider_endpoint_id", endpointID, err))
	}

	if err := d.Set("private_link_service_resource_id", serviceEndpoint.PrivateLinkServiceResourceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "private_link_service_resource_id", endpointID, err))
	}

	if err := d.Set("private_endpoint_ip_address", serviceEndpoint.PrivateEndpointIPAddress); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "private_endpoint_ip_address", endpointID, err))
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endpointID,
	}))

	return nil
}
