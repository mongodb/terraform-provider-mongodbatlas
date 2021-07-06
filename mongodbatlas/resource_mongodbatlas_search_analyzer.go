package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
	"log"
	"strings"
)

func resourceMongoDBAtlasSearchAnalyzers() *schema.Resource {
	return &schema.Resource{
		Create: resourceMongoDBAtlasSearchAnalyzersCreate,
		Read:   resourceMongoDBAtlasSearchAnalyzersRead,
		Update: resourceMongoDBAtlasSearchAnalyzersUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceMongoDBAtlasSearchAnalyzersImportState,
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
			"search_analyzers": {
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
							Type:     schema.TypeBool,
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
		},
	}
}

func resourceMongoDBAtlasSearchAnalyzersImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a search analyzer, use the format {project_id}-{cluster_name}")
	}

	projectID := parts[0]
	clusterName := parts[1]

	_, _, err := conn.Search.ListAnalyzers(context.Background(), projectID, clusterName, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't import search analyzers in projectID (%s) and Cluster (%s), error: %s", projectID, clusterName, err)
	}

	if err := d.Set("project_id", projectID); err != nil {
		log.Printf("[WARN] Error setting project_id for (%s): %s", projectID, err)
	}

	if err := d.Set("cluster_name", clusterName); err != nil {
		log.Printf("[WARN] Error setting cluster_name for (%s): %s", clusterName, err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceMongoDBAtlasSearchAnalyzersUpdate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	var searchAnalyzers []*matlas.SearchAnalyzer

	if d.HasChange("search_analyzers") {
		searchAnalyzers = expandSearchAnalyzers(d.Get("search_analyzers").([]interface{}))
	}

	_, _, err := conn.Search.UpdateAllAnalyzers(context.Background(), projectID, clusterName, searchAnalyzers)
	if err != nil {
		return fmt.Errorf("error updating search analyzers : %s", err)
	}

	return resourceMongoDBAtlasSearchAnalyzersRead(d, meta)
}

func resourceMongoDBAtlasSearchAnalyzersRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]

	searchAnalyzers, _, err := conn.Search.ListAnalyzers(context.Background(), projectID, clusterName, nil)
	if err != nil {
		// case 404
		// deleted in the backend case
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error getting search analyzers information: %s", err)
	}

	if err := d.Set("search_analyzers", flattenSearchAnalyzers(searchAnalyzers)); err != nil {
		return fmt.Errorf("error setting `search_analyzers` : %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return nil
}

func resourceMongoDBAtlasSearchAnalyzersCreate(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	clusterName := d.Get("cluster_name").(string)

	searchAnalyzers := expandSearchAnalyzers(d.Get("search_analyzers").([]interface{}))

	_, _, err := conn.Search.UpdateAllAnalyzers(context.Background(), projectID, clusterName, searchAnalyzers)
	if err != nil {
		return fmt.Errorf("error updating search analyzers: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
	}))

	return resourceMongoDBAtlasSearchAnalyzersRead(d, meta)
}

func expandSearchAnalyzers(p []interface{}) []*matlas.SearchAnalyzer {
	mappings := make([]*matlas.SearchAnalyzer, len(p))

	for k, v := range p {
		mapping := v.(map[string]interface{})
		mappings[k] = &matlas.SearchAnalyzer{
			BaseAnalyzer:     mapping["base_analyzer"].(string),
			IgnoreCase:       pointy.Bool(mapping["ignore_case"].(bool)),
			MaxTokenLength:   pointy.Int(mapping["max_token_length"].(int)),
			Name:             mapping["name"].(string),
			StemExclusionSet: mapping["stem_exclusion_set"].([]string),
			Stopwords:        mapping["stopwords"].([]string),
		}
	}

	return mappings
}

func flattenSearchAnalyzers(analyzers []*matlas.SearchAnalyzer) []map[string]interface{} {
	analyzersMap := make([]map[string]interface{}, 0)
	for _, v := range analyzers {
		analyzersMap = append(analyzersMap, map[string]interface{}{
			"base_analyzer":      v.BaseAnalyzer,
			"Ignore_case":        v.IgnoreCase,
			"max_token_length":   v.MaxTokenLength,
			"name":               v.Name,
			"stem_exclusion_set": v.StemExclusionSet,
			"stopwords":          v.Stopwords,
		})
	}

	return analyzersMap
}
