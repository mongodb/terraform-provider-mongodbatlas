package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasPrivateLinkEndpointServerless() *schema.Resource {
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

func dataSourceMongoDBAtlasPrivateEndpointServiceServerlessLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

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

	d.SetId(encodeStateID(map[string]string{
		"project_id":    projectID,
		"instance_name": instanceName,
		"endpoint_id":   endpointID,
	}))

	return nil
}
