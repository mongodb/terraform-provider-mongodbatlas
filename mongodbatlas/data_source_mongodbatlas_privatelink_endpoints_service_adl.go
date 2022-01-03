package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasPrivateLinkEndpointsServiceADL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivateLinkEndpointsServiceADLRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rel": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provider_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasPrivateLinkEndpointsServiceADLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	privateLinkEndpoints, _, err := conn.DataLakes.ListPrivateLinkEndpoint(ctx, projectID)
	if err != nil {
		return diag.Errorf("error getting ADL PrivateLink Endpoints Information: %s", err)
	}

	if err := d.Set("links", flattenADLPrivateEndpointLinks(privateLinkEndpoints.Links)); err != nil {
		return diag.Errorf("error setting `results`: %s", err)
	}

	if err := d.Set("results", flattenADLPrivateLinkEndpoints(privateLinkEndpoints.Results)); err != nil {
		return diag.Errorf("error setting `results`: %s", err)
	}

	if err := d.Set("total_count", privateLinkEndpoints.TotalCount); err != nil {
		return diag.Errorf("error setting `total_count`: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenADLPrivateEndpointLinks(links []*matlas.Link) []map[string]interface{} {
	linksList := make([]map[string]interface{}, 0)

	for _, link := range links {
		mLink := map[string]interface{}{
			"href": link.Href,
			"rel":  link.Rel,
		}
		linksList = append(linksList, mLink)
	}

	return linksList
}

func flattenADLPrivateLinkEndpoints(privateLinks []*matlas.PrivateLinkEndpointDataLake) []map[string]interface{} {
	var results []map[string]interface{}

	if len(privateLinks) == 0 {
		return results
	}

	results = make([]map[string]interface{}, len(privateLinks))

	for k, privateLink := range privateLinks {
		results[k] = map[string]interface{}{
			"endpoint_id":   privateLink.EndpointID,
			"type":          privateLink.Type,
			"provider_name": privateLink.Provider,
			"comment":       privateLink.Comment,
		}
	}

	return results
}
