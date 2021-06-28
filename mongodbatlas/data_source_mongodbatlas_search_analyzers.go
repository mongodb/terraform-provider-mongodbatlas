package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
	"net/http"
)

func dataSourceMongoDBAtlasSearchAnalyzers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMongoDBAtlasSearchAnalyzersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_analyzer": {
							Type:     schema.TypeString,
							Computed: true,
							Required: false,
						},
						"ignore_case": {
							Type:     schema.TypeString,
							Required: true,
						},
						"max_token_length": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: false,
						},
						"stem_exclusion_set": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"stopwords": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

func dataSourceMongoDBAtlasSearchAnalyzersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID, projectIDOK := d.GetOk("project_id")
	clusterName, clusterNameOK := d.GetOk("cluster_name")

	if !projectIDOK || !clusterNameOK {
		return errors.New("project_id and cluster_name must be configured")
	}

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	analyzers, resp, err := conn.Search.ListAnalyzers(context.Background(), projectID.(string), clusterName.(string), options)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("error reading analyzers list for project(%s): %s", projectID, err)
	}

	if err := d.Set("results", flattenSearchAnalyzers(analyzers)); err != nil {
		return fmt.Errorf("error setting `result` for search analyzers: %s", err)
	}

	if err := d.Set("total_count", len(analyzers)); err != nil {
		return fmt.Errorf("error setting `name`: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}
