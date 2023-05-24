package mongodbatlas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const errorPrivateEndpointServiceDataFederationOnlineArchiveList = "error reading Private Endpoings for projectId %s: %s"

func dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchives() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchivesRead,
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
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasPrivatelinkEndpointServiceDataFederationOnlineArchivesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	privateEndpoints, _, err := conn.DataLakes.ListPrivateLinkEndpoint(context.Background(), projectID)
	if err != nil {
		return diag.Errorf(errorPrivateEndpointServiceDataFederationOnlineArchiveList, projectID, err)
	}

	if err := d.Set("results", flattenPrivateLinkEndpointDataLakeResponse(privateEndpoints.Results)); err != nil {
		return diag.FromErr(fmt.Errorf(errorDataLakeSetting, "results", projectID, err))
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenPrivateLinkEndpointDataLakeResponse(atlasPrivateLinkEndpointDataLakes []*matlas.PrivateLinkEndpointDataLake) []map[string]interface{} {

	if len(atlasPrivateLinkEndpointDataLakes) == 0 {
		return []map[string]interface{}{}
	}

	results := make([]map[string]interface{}, len(atlasPrivateLinkEndpointDataLakes))

	for i, atlasPrivateLinkEndpointDataLake := range atlasPrivateLinkEndpointDataLakes {
		results[i] = map[string]interface{}{
			"endpoint_id":   atlasPrivateLinkEndpointDataLake.EndpointID,
			"provider_name": atlasPrivateLinkEndpointDataLake.Provider,
			"comment":       atlasPrivateLinkEndpointDataLake.Comment,
			"type":          atlasPrivateLinkEndpointDataLake.Type,
		}
	}

	return results
}
