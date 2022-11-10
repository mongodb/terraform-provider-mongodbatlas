package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasPrivateLinkEndpointsServiceServerless() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasPrivateLinkEndpointsServiceServerlessRead,
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
			"page_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"items_per_page": {
				Type:     schema.TypeInt,
				Optional: true,
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

func dataSourceMongoDBAtlasPrivateLinkEndpointsServiceServerlessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)
	instanceName := d.Get("instance_name").(string)

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	privateLinkEndpoints, _, err := conn.ServerlessPrivateEndpoints.List(ctx, projectID, instanceName, options)
	if err != nil {
		return diag.Errorf("error getting Serverless PrivateLink Endpoints Information: %s", err)
	}

	if err := d.Set("results", flattenServerlessPrivateLinkEndpoints(privateLinkEndpoints)); err != nil {
		return diag.Errorf("error setting `results`: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenServerlessPrivateLinkEndpoints(privateLinks []matlas.ServerlessPrivateEndpointConnection) []map[string]interface{} {
	var results []map[string]interface{}

	if len(privateLinks) == 0 {
		return results
	}

	results = make([]map[string]interface{}, len(privateLinks))

	for k := range privateLinks {
		results[k] = map[string]interface{}{
			"endpoint_id":                privateLinks[k].ID,
			"endpoint_service_name":      privateLinks[k].EndpointServiceName,
			"cloud_provider_endpoint_id": privateLinks[k].CloudProviderEndpointID,
			"comment":                    privateLinks[k].Comment,
			"error_message":              privateLinks[k].ErrorMessage,
			"status":                     privateLinks[k].Status,
		}
	}

	return results
}
