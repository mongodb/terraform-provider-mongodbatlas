package mongodbatlas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func dataSourceMongoDBAtlasSearchIndexes() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceMongoDBAtlasSearchIndexesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database-name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"collection-name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"indexID": {
							Type:     schema.TypeString,
							Computed: true,
							Required: false,
						},
						"analyzer": {
							Type:     schema.TypeString,
							Required: true,
						},
						"analyzers": {
							Type:     schema.TypeString, //TODO: change type
							Required: true,
						},
						"collectionName": {
							Type:     schema.TypeString,
							Required: true,
						},
						"database": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"searchAnalyzer": {
							Type:     schema.TypeString,
							Required: false,
						},
						"mappings": {
							Type:     schema.TypeSet,
							Optional: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dynamic": {
										Type:     schema.TypeBool,
										Optional: false,
									},
									"fields": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"field": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Required: false,
						},
					},
				},
			},
		},
	}
}

func dataSourceMongoDBAtlasSearchIndexesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	databaseName := d.Get("database_name").(string)
	collectionName := d.Get("collection_name").(string)

	searchIndexes, _, err := conn.Search.ListIndexes(context.Background(), projectID, clusterName, databaseName, collectionName, nil)
	if err != nil {
		return fmt.Errorf("error getting search indexes information: %s", err)
	}

	if err := d.Set("results", flattenSearchIndexes(searchIndexes)); err != nil {
		return fmt.Errorf("error setting `result` for search indexes: %s", err)
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenSearchIndexes(searchIndexes []*matlas.SearchIndex) []map[string]interface{} {
	var searchIndexesMap []map[string]interface{}

	if len(searchIndexes) > 0 {
		searchIndexesMap = make([]map[string]interface{}, len(searchIndexes))

		for i := range searchIndexes {
			searchIndexesMap[i] = map[string]interface{}{
				"analyzer":       searchIndexes[i].Analyzer,
				"analyzers":      searchIndexes[i].Analyzers,
				"collectionName": searchIndexes[i].CollectionName,
				"database":       searchIndexes[i].Database,
				"indexID":        searchIndexes[i].IndexID,
				//"mappings":       searchIndexes[i].Mappings,  //TODO: create Flatten function
				"name":           searchIndexes[i].Name,
				"searchAnalyzer": searchIndexes[i].SearchAnalyzer,
				"status":         searchIndexes[i].Status,
			}
		}
	}

	return searchIndexesMap
}
