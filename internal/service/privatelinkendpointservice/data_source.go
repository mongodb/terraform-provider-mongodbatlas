package privatelinkendpointservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/spf13/cast"
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
					},
				},
			},
			"gcp_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gcp_endpoint_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the GCP endpoint. Only populated for port-based architecture.",
			},
			"port_mapping_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	// TODO: update before merging to master: connV2 := d.Client.AtlasV2
	connV2 := meta.(*config.MongoDBClient).AtlasPreview

	projectID := d.Get("project_id").(string)
	privateLinkID := conversion.GetEncodedID(d.Get("private_link_id").(string), "private_link_id")
	providerName := d.Get("provider_name").(string)
	endpointServiceID := conversion.GetEncodedID(d.Get("endpoint_service_id").(string), "endpoint_service_id")

	serviceEndpoint, _, err := connV2.PrivateEndpointServicesApi.GetPrivateEndpoint(ctx, projectID, providerName, endpointServiceID, privateLinkID).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorServiceEndpointRead, endpointServiceID, err))
	}

	if err := d.Set("delete_requested", cast.ToBool(serviceEndpoint.GetDeleteRequested())); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "delete_requested", endpointServiceID, err))
	}

	if err := d.Set("error_message", serviceEndpoint.GetErrorMessage()); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "error_message", endpointServiceID, err))
	}

	if err := d.Set("aws_connection_status", serviceEndpoint.GetConnectionStatus()); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "aws_connection_status", endpointServiceID, err))
	}

	if err := d.Set("interface_endpoint_id", serviceEndpoint.GetInterfaceEndpointId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "interface_endpoint_id", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_connection_name", serviceEndpoint.GetPrivateEndpointConnectionName()); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "private_endpoint_connection_name", endpointServiceID, err))
	}

	if err := d.Set("private_endpoint_resource_id", serviceEndpoint.GetPrivateEndpointResourceId()); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "private_endpoint_resource_id", endpointServiceID, err))
	}

	if strings.EqualFold(providerName, "azure") {
		if err := d.Set("azure_status", serviceEndpoint.GetStatus()); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "azure_status", endpointServiceID, err))
		}

		if err := d.Set("private_endpoint_ip_address", serviceEndpoint.GetPrivateEndpointIPAddress()); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "private_endpoint_ip_address", endpointServiceID, err))
		}
	}

	if err := d.Set("endpoint_service_id", endpointServiceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorEndpointSetting, "endpoint_service_id", endpointServiceID, err))
	}

	if strings.EqualFold(providerName, "gcp") {
		if err := d.Set("port_mapping_enabled", serviceEndpoint.GetPortMappingEnabled()); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "port_mapping_enabled", privateLinkID, err))
		}

		if err := d.Set("gcp_status", serviceEndpoint.GetStatus()); err != nil {
			return diag.FromErr(fmt.Errorf(errorEndpointSetting, "gcp_status", endpointServiceID, err))
		}

		if serviceEndpoint.GetPortMappingEnabled() && serviceEndpoint.Endpoints != nil && len(*serviceEndpoint.Endpoints) == 1 {
			firstEndpoint := (*serviceEndpoint.Endpoints)[0]

			if err := d.Set("gcp_endpoint_status", firstEndpoint.GetStatus()); err != nil {
				return diag.FromErr(fmt.Errorf(errorEndpointSetting, "gcp_endpoint_status", endpointServiceID, err))
			}

			if err := d.Set("private_endpoint_ip_address", firstEndpoint.GetIpAddress()); err != nil {
				return diag.FromErr(fmt.Errorf(errorEndpointSetting, "private_endpoint_ip_address", endpointServiceID, err))
			}
		} else {
			if err := d.Set("endpoints", flattenGCPEndpoints(serviceEndpoint.Endpoints)); err != nil {
				return diag.FromErr(fmt.Errorf(errorEndpointSetting, "endpoints", endpointServiceID, err))
			}
		}
	}

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":          projectID,
		"private_link_id":     privateLinkID,
		"endpoint_service_id": endpointServiceID,
		"provider_name":       providerName,
	}))

	return nil
}
