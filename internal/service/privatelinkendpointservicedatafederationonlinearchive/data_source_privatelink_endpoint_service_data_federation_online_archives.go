package privatelinkendpointservicedatafederationonlinearchive

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/datalakepipeline"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

const errorPrivateEndpointServiceDataFederationOnlineArchiveList = "error reading Private Endpoings for projectId %s: %s"

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

	privateEndpoints, _, err := connV2.DataFederationApi.ListDataFederationPrivateEndpoints(ctx, projectID).Execute()
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveList, projectID, err)
	}

	if err := d.Set("results", flattenPrivateLinkEndpointDataLakeResponse(privateEndpoints.GetResults())); err != nil {
		return diag.FromErr(fmt.Errorf(datalakepipeline.ErrorDataLakeSetting, "results", projectID, err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenPrivateLinkEndpointDataLakeResponse(atlasPrivateLinkEndpointDataLakes []admin.PrivateNetworkEndpointIdEntry) []map[string]any {
	if len(atlasPrivateLinkEndpointDataLakes) == 0 {
		return []map[string]any{}
	}

	results := make([]map[string]any, len(atlasPrivateLinkEndpointDataLakes))

	for i, atlasPrivateLinkEndpointDataLake := range atlasPrivateLinkEndpointDataLakes {
		results[i] = map[string]any{
			"endpoint_id":                atlasPrivateLinkEndpointDataLake.GetEndpointId(),
			"provider_name":              atlasPrivateLinkEndpointDataLake.GetProvider(),
			"comment":                    atlasPrivateLinkEndpointDataLake.GetComment(),
			"type":                       atlasPrivateLinkEndpointDataLake.GetType(),
			"region":                     atlasPrivateLinkEndpointDataLake.GetRegion(),
			"customer_endpoint_dns_name": atlasPrivateLinkEndpointDataLake.GetCustomerEndpointDNSName(),
		}
	}

	return results
}
