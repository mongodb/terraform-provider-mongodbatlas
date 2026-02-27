package privatelinkendpoint

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"private_link_id": {
							Type:     schema.TypeString,
							Computed: true,
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
						"port_mapping_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Flag that indicates whether this resource uses GCP port-mapping. When `true`, it uses the port-mapped architecture. When `false` or unset, it uses the GCP legacy private endpoint architecture. Only applicable for GCP provider.",
						},
					},
				},
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	conn := meta.(*config.MongoDBClient).AtlasV2

	projectID := d.Get("project_id").(string)
	providerName := d.Get("provider_name").(string)
	privateEndpoints, _, err := conn.PrivateEndpointServicesApi.ListPrivateEndpointService(ctx, projectID, providerName).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting Private Endpoints: %s", err))
	}

	if err := d.Set("results", flattenPrivateEndpoints(privateEndpoints)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `results`: %s", err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenPrivateEndpoints(privateEndpoints []admin.EndpointService) []map[string]any {
	var results []map[string]any

	if len(privateEndpoints) > 0 {
		results = make([]map[string]any, len(privateEndpoints))

		for k, privateEndpoint := range privateEndpoints {
			results[k] = map[string]any{
				"private_link_id":                  privateEndpoint.GetId(),
				"endpoint_service_name":            privateEndpoint.GetEndpointServiceName(),
				"error_message":                    privateEndpoint.GetErrorMessage(),
				"interface_endpoints":              privateEndpoint.GetInterfaceEndpoints(),
				"private_endpoints":                privateEndpoint.GetPrivateEndpoints(),
				"private_link_service_name":        privateEndpoint.GetPrivateLinkServiceName(),
				"private_link_service_resource_id": privateEndpoint.GetPrivateLinkServiceResourceId(),
				"status":                           privateEndpoint.GetStatus(),
				"endpoint_group_names":             privateEndpoint.GetEndpointGroupNames(),
				"region_name":                      privateEndpoint.GetRegionName(),
				"service_attachment_names":         privateEndpoint.GetServiceAttachmentNames(),
				"port_mapping_enabled":             privateEndpoint.GetPortMappingEnabled(),
			}
		}
	}

	return results
}
