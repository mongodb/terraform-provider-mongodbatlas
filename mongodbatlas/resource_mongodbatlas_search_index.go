package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"strings"

	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasSearchIndex() *schema.Resource {

	return &schema.Resource{
		Create: resourceMongoDBAtlasSearchIndexCreate,
		Read:   resourceMongoDBAtlasSearchIndexRead,
		Update: resourceMongoDBAtlasSearchIndexUpdate,
		Delete: resourceMongoDBAtlasSearchIndexDelete,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasSearchIndexImportState,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
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
	}
}

func resourceMongoDBAtlasSearchIndexImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*matlas.Client)

	parts := strings.SplitN(d.Id(), "-", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a search index, use the format {project_id}-{cluster_name}-{index_id}")
	}

	projectID := parts[0]
	clusterName := parts[1]
	indexID := parts[2]

	_, _, err := conn.Search.GetIndex(context.Background(), projectID, clusterName, indexID)
	if err != nil {
		return nil, fmt.Errorf("couldn't import search index (%s) in projectID (%s) and Cluster (%s), error: %s", indexID, projectID, clusterName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", projectID, err)
	}

	if err := d.Set("cluster_name", clusterName); err != nil {
		log.Printf("[WARN] Error setting cluster_name for (%s): %s", clusterName, err)
	}

	if err := d.Set("index_id", indexID); err != nil {
		log.Printf("[WARN] Error setting index_id for (%s): %s", indexID, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"index_id":     indexID,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceMongoDBAtlasSearchIndexDelete(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	_, err := conn.Search.DeleteIndex(context.Background(), projectID, clusterName, indexID)
	if err != nil {
		return fmt.Errorf("error deleting search index (%s): %s", d.Get("name").(string), err)
	}

	return nil
}

func resourceMongoDBAtlasSearchIndexUpdate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	searchIndex, _, err := conn.Search.GetIndex(context.Background(), projectID, clusterName, indexID)
	if err != nil {
		return fmt.Errorf("error getting search index information: %s", err)
	}

	if d.HasChange("analyzer") {
		searchIndex.Analyzer = d.Get("analyzer").(string)
	}

	if d.HasChange("collectionName") {
		searchIndex.CollectionName = d.Get("collectionName").(string)
	}

	if d.HasChange("database") {
		searchIndex.Database = d.Get("database").(string)
	}

	if d.HasChange("name") {
		searchIndex.Name = d.Get("name").(string)
	}

	if d.HasChange("searchAnalyzer") {
		searchIndex.SearchAnalyzer = d.Get("searchAnalyzer").(string)
	}

	/*
		if d.HasChange("mappings") {
			searchIndex.Mappings = d.Get("mappings")
		}
	*/

	_, _, err = conn.Search.UpdateIndex(context.Background(), projectID, clusterName, indexID, searchIndex)
	if err != nil {
		return fmt.Errorf("error updating search index (%s): %s", searchIndex.Name, err)
	}

	return resourceMongoDBAtlasSearchIndexRead(d, meta)
}

func resourceMongoDBAtlasSearchIndexRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	searchIndex, _, err := conn.Search.GetIndex(context.Background(), projectID, clusterName, indexID)
	if err != nil {
		// case 404
		// deleted in the backend case
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error getting search index information: %s", err)
	}

	if err := d.Set("analyzer", searchIndex.Analyzer); err != nil {
		return fmt.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("collectionName", searchIndex.CollectionName); err != nil {
		return fmt.Errorf("error setting `collectionName` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("database", searchIndex.Database); err != nil {
		return fmt.Errorf("error setting `database` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", searchIndex.Name); err != nil {
		return fmt.Errorf("error setting `name` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("searchAnalyzer", searchIndex.SearchAnalyzer); err != nil {
		return fmt.Errorf("error setting `searchAnalyzer` for search index (%s): %s", d.Id(), err)
	}

	/*if err := d.Set("mapping", flattenSearchIndexFields(*searchIndex.Mappings)); err != nil {
		return fmt.Errorf("error setting `scopes` for database user (%s): %s", d.Id(), err)
	}
	*/

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"index_id":     indexID,
	}))

	return nil
}

func resourceMongoDBAtlasSearchIndexCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*matlas.Client)
	projectID := d.Get("project_id").(string)

	clusterName := d.Get("cluster_name").(string)

	searchIndexRequest := &matlas.SearchIndex{
		Analyzer: d.Get("analyzer").(string),
		//Analyzers:      d.Get("analyzers").(map[string]interface{}), //TODO: check if is correct type
		CollectionName: d.Get("collectionName").(string),
		Database:       d.Get("database").(string),
		Mappings:       d.Get("mappings").(*matlas.IndexMapping),
		Name:           d.Get("name").(string),
		SearchAnalyzer: d.Get("searchAnalyzer").(string),
		Status:         d.Get("status").(string),
	}

	dbSearchIndexRes, _, err := conn.Search.CreateIndex(context.Background(), projectID, clusterName, searchIndexRequest)
	if err != nil {
		return fmt.Errorf("error creating database user: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"index_id":     dbSearchIndexRes.IndexID,
	}))

	return resourceMongoDBAtlasSearchIndexRead(d, meta)
}

/*
func flattenSearchIndexFields(l matlas.IndexMapping) []map[string]interface{} {
	mapping := make([]map[string]interface{}, 1)

	mapping[0]["dynamic"] = l.Dynamic

	if !l.Dynamic {
		mapping[0]["fields"] = make([]map[string]interface{}, len(*l.Fields))

		for i, field := range *l.Fields {
			scopes[i] = map[string]interface{}{
				"name": v.Name,
				"type": v.Type,
			}
		}
	}

	return scopes
}

func returnFlattenFields(fields []map[string]interface{}) []map[string]interface{} {
	mapFields = make([]map[string]interface{}, len(fields))

	for i, field := range fields {
		scopes[i] = map[string]interface{}{
			"name":  f,
			"field": v.Type,
		}
	}

	return mapFields
}
*/
