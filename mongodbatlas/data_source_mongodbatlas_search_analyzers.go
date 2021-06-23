package mongodbatlas

import (
	"context"
	"fmt"
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
			"results": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"baseAnalyzer": {
							Type:     schema.TypeString,
							Computed: true,
							Required: false,
						},
						"ignoreCase": {
							Type:     schema.TypeString,
							Required: true,
						},
						"maxTokenLength": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: false,
						},
						"stemExclusionSet": {
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
		},
	}
}

func dataSourceMongoDBAtlasSearchAnalyzersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)

	analyzers, resp, err := conn.Search.ListAnalyzers(context.Background(), projectID, clusterName, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("error reading analyzers list for project(%s): %s", projectID, err)
	}

	d.Set("results", flattenSearchAnalyzers(analyzers))

	return nil
}
