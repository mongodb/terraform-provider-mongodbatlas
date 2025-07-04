package privatelinkendpointserviceserverless

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourcePluralRead,
		DeprecationMessage: fmt.Sprintf(constant.DeprecationDataSourceByDateWithExternalLink, "March 2025", "https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/serverless-shared-migration-guide"),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_provider_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"endpoint_service_name": {
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
				},
			},
		},
	}
}

func dataSourcePluralRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)

	privateLinkEndpoints, _, err := connV2.ServerlessPrivateEndpointsApi.ListServerlessPrivateEndpoints(ctx, projectID, instanceName).Execute()
	if err != nil {
		return diag.Errorf("error getting Serverless PrivateLink Endpoints Information: %s", err)
	}

	if err := d.Set("results", flattenServerlessPrivateLinkEndpoints(privateLinkEndpoints)); err != nil {
		return diag.Errorf("error setting `results`: %s", err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenServerlessPrivateLinkEndpoints(privateLinks []admin.ServerlessTenantEndpoint) []map[string]any {
	if len(privateLinks) == 0 {
		return nil
	}

	results := make([]map[string]any, len(privateLinks))

	for k := range privateLinks {
		results[k] = map[string]any{
			"endpoint_id":                      privateLinks[k].GetId(),
			"endpoint_service_name":            privateLinks[k].GetEndpointServiceName(),
			"cloud_provider_endpoint_id":       privateLinks[k].GetCloudProviderEndpointId(),
			"private_link_service_resource_id": privateLinks[k].GetPrivateLinkServiceResourceId(),
			"private_endpoint_ip_address":      privateLinks[k].GetPrivateEndpointIpAddress(),
			"comment":                          privateLinks[k].GetComment(),
			"error_message":                    privateLinks[k].GetErrorMessage(),
			"status":                           privateLinks[k].GetStatus(),
		}
	}

	return results
}
