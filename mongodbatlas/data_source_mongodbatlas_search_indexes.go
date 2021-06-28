package mongodbatlas

import (
	"context"
	"errors"
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
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"collection_name": {
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
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: returnSearchIndexSchema(),
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasSearchIndexesRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)

	projectID, projectIDOK := d.GetOk("project_id")
	clusterName, clusterNameOk := d.GetOk("cluster_name")
	databaseName, databaseNameOK := d.GetOk("database_name")
	collectionName, collectionNameOK := d.GetOk("collection_name")

	if !projectIDOK || !clusterNameOk || !databaseNameOK || !collectionNameOK {
		return errors.New("project_id, cluster_name, database_name and collection_name must be configured")
	}

	options := &matlas.ListOptions{
		PageNum:      d.Get("page_num").(int),
		ItemsPerPage: d.Get("items_per_page").(int),
	}

	searchIndexes, _, err := conn.Search.ListIndexes(context.Background(), projectID.(string), clusterName.(string), databaseName.(string), collectionName.(string), options)
	if err != nil {
		return fmt.Errorf("error getting search indexes information: %s", err)
	}

	if err := d.Set("results", flattenSearchIndexes(searchIndexes)); err != nil {
		return fmt.Errorf("error setting `result` for search indexes: %s", err)
	}

	if err := d.Set("total_count", len(searchIndexes)); err != nil {
		return fmt.Errorf("error setting `name`: %s", err)
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
				"analyzer":        searchIndexes[i].Analyzer,
				"analyzers":       flattenSearchIndexCustomAnalyzers(searchIndexes[i].Analyzers),
				"collectionName":  searchIndexes[i].CollectionName,
				"database":        searchIndexes[i].Database,
				"indexID":         searchIndexes[i].IndexID,
				"mapping_dynamic": searchIndexes[i].Mappings.Dynamic,
				"mappings_fields": marshallSearchIndexMappingFields(searchIndexes[i].Mappings.Fields),
				"name":            searchIndexes[i].Name,
				"searchAnalyzer":  searchIndexes[i].SearchAnalyzer,
				"status":          searchIndexes[i].Status,
			}
		}
	}

	return searchIndexesMap
}
