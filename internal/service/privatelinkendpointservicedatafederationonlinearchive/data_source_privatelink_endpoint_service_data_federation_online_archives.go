package privatelinkendpointservicedatafederationonlinearchive

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

const (
	errorPrivateEndpointServiceDataFederationOnlineArchiveList = "error reading Private Endpoints for projectId %s: %s"
	errorDataFederationSetting                                 = "error setting `%s` for MongoDB Atlas Data Federation (%s): %s"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluralRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"customer_endpoint_dns_name": {
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

	privateEndpoints, _, err := connV2.DataFederationApi.ListPrivateEndpointIds(ctx, projectID).Execute()
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveList, projectID, err)
	}

	if err := d.Set("results", flattenPrivateLinkEndpointDataFederationResponse(privateEndpoints.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataFederationSetting, "results", projectID, err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenPrivateLinkEndpointDataFederationResponse(entries []admin.PrivateNetworkEndpointIdEntry) []map[string]any {
	results := make([]map[string]any, len(entries))
	for i, entry := range entries {
		results[i] = map[string]any{
			"endpoint_id":                entry.GetEndpointId(),
			"provider_name":              entry.GetProvider(),
			"comment":                    entry.GetComment(),
			"type":                       entry.GetType(),
			"region":                     entry.GetRegion(),
			"customer_endpoint_dns_name": entry.GetCustomerEndpointDNSName(),
		}
	}
	return results
}
