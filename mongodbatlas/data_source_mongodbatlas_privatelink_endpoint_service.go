package mongodbatlas

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
)

func dataSourceMongoDBAtlasPrivateEndpointServiceLink() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivateEndpointServiceLinkRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_link_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface_endpoint_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_connection_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_endpoint_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_requested": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_connection_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_attachment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"gcp_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasPrivateEndpointServiceLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID := d.Get("project_id").(string)
	privateLinkID := getEncodedID(d.Get("private_link_id").(string), "private_link_id")
	endpointServiceID := getEncodedID(d.Get("endpoint_service_id").(string), "endpoint_service_id")
	providerName := d.Get("provider_name").(string)

	serviceEndpoint, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(ctx, projectID, providerName, privateLinkID, endpointServiceID)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServiceEndpointRead, endpointServiceID, err))
	}

	if err := d.Set("delete_requested", cast.ToBool(serviceEndpoint.DeleteRequested)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "delete_requested", endpointServiceID, err))
	}

	if err := d.Set("error_message", serviceEndpoint.ErrorMessage); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "error_message", endpointServiceID, err))
	}

	if err := d.Set("aws_connection_status", serviceEndpoint.AWSConnectionStatus); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "aws_connection_status", endpointServiceID, err))
	}

	if strings.EqualFold(providerName, "azure") {
		if err := d.Set("azure_status", serviceEndpoint.Status); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "azure_status", endpointServiceID, err))
		}
	}

	if err := d.Set("endpoints", flattenGCPEndpoints(serviceEndpoint.Endpoints)); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "endpoints", endpointServiceID, err))
	}

	if strings.EqualFold(providerName, "gcp") {
		if err := d.Set("gcp_status", serviceEndpoint.Status); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "gcp_status", endpointServiceID, err))
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return nil
}
