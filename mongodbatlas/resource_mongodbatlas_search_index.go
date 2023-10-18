package mongodbatlas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20231001001/admin"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasSearchIndex() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMongoDBAtlasSearchIndexCreate,
		ReadContext:          resourceMongoDBAtlasSearchIndexRead,
		UpdateWithoutTimeout: resourceMongoDBAtlasSearchIndexUpdate,
		DeleteContext:        resourceMongoDBAtlasSearchIndexDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasSearchIndexImportState,
		},
		Schema: returnSearchIndexSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Hour),
			Update: schema.DefaultTimeout(3 * time.Hour),
			Delete: schema.DefaultTimeout(3 * time.Hour),
		},
	}
}

func returnSearchIndexSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"index_id": {
			Type:     schema.TypeString,
			Computed: true,
			Required: false,
		},
		"analyzer": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"analyzers": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: validateSearchAnalyzersDiff,
		},
		"collection_name": {
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
		"search_analyzer": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mappings_dynamic": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"mappings_fields": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: validateSearchIndexMappingDiff,
		},
		"synonyms": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"analyzer": {
						Type:     schema.TypeString,
						Required: true,
					},
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"source_collection": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"status": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"wait_for_index_build_completion": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func resourceMongoDBAtlasSearchIndexImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "--", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a search index, use the format {project_id}--{cluster_name}--{index_id}")
	}

	projectID := parts[0]
	clusterName := parts[1]
	indexID := parts[2]

	_, _, err := conn.Search.GetIndex(ctx, projectID, clusterName, indexID)
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

func resourceMongoDBAtlasSearchIndexDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	connV2 := meta.(*MongoDBClient).AtlasV2
	_, _, err := connV2.AtlasSearchApi.DeleteAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	if err != nil {
		return diag.Errorf("error deleting search index (%s): %s", d.Get("name").(string), err)
	}
	return nil
}

func resourceMongoDBAtlasSearchIndexUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	searchIndex, _, err := conn.Search.GetIndex(ctx, projectID, clusterName, indexID)
	if err != nil {
		return diag.Errorf("error getting search index information: %s", err)
	}

	if d.HasChange("analyzer") {
		searchIndex.Analyzer = d.Get("analyzer").(string)
	}

	if d.HasChange("analyzers") {
		searchIndex.Analyzers = unmarshalSearchIndexAnalyzersFields(d.Get("analyzers").(string))
	}

	if d.HasChange("collection_name") {
		searchIndex.CollectionName = d.Get("collection_name").(string)
	}

	if d.HasChange("database") {
		searchIndex.Database = d.Get("database").(string)
	}

	if d.HasChange("name") {
		searchIndex.Name = d.Get("name").(string)
	}

	if d.HasChange("search_analyzer") {
		searchIndex.SearchAnalyzer = d.Get("search_analyzer").(string)
	}

	if d.HasChange("mappings_dynamic") {
		searchIndex.Mappings.Dynamic = d.Get("mappings_dynamic").(bool)
	}

	if d.HasChange("mappings_fields") {
		mappingFields := unmarshalSearchIndexMappingFields(d.Get("mappings_fields").(string))
		searchIndex.Mappings.Fields = &mappingFields
	}

	if d.HasChange("synonyms") {
		synonyms := expandSearchIndexSynonyms(d)
		searchIndex.Synonyms = synonyms
	}

	searchIndex.IndexID = ""
	dbSearchIndexRes, _, err := conn.Search.UpdateIndex(context.Background(), projectID, clusterName, indexID, searchIndex)
	if err != nil {
		return diag.Errorf("error updating search index (%s): %s", searchIndex.Name, err)
	}

	if d.Get("wait_for_index_build_completion").(bool) {
		timeout := d.Timeout(schema.TimeoutCreate)
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"IN_PROGRESS", "MIGRATING"},
			Target:     []string{"STEADY"},
			Refresh:    resourceSearchIndexRefreshFunc(ctx, clusterName, projectID, dbSearchIndexRes.IndexID, conn),
			Timeout:    timeout,
			MinTimeout: 1 * time.Minute,
			Delay:      1 * time.Minute,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			d.SetId(encodeStateID(map[string]string{
				"project_id":   projectID,
				"cluster_name": clusterName,
				"index_id":     dbSearchIndexRes.IndexID,
			}))
			resourceMongoDBAtlasSearchIndexDelete(ctx, d, meta)
			d.SetId("")
			return diag.FromErr(fmt.Errorf("error creating index in cluster (%s): %s", clusterName, err))
		}
	}

	return resourceMongoDBAtlasSearchIndexRead(ctx, d, meta)
}

func resourceMongoDBAtlasSearchIndexRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	connV2 := meta.(*MongoDBClient).AtlasV2
	searchIndex, resp, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	if err != nil {
		// deleted in the backend case
		if resp.StatusCode == 404 && !d.IsNewResource() {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting search index information: %s", err)
	}

	if err := d.Set("index_id", indexID); err != nil {
		return diag.Errorf("error setting `index_id` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("analyzer", searchIndex.Analyzer); err != nil {
		return diag.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
	}

	if len(searchIndex.Analyzers) > 0 {
		searchIndexMappingFields, err := marshallSearchIndexAnalyzers(searchIndex.Analyzers)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("analyzers", searchIndexMappingFields); err != nil {
			return diag.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
		}
	}

	if err := d.Set("collection_name", searchIndex.CollectionName); err != nil {
		return diag.Errorf("error setting `collectionName` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("database", searchIndex.Database); err != nil {
		return diag.Errorf("error setting `database` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", searchIndex.Name); err != nil {
		return diag.Errorf("error setting `name` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("search_analyzer", searchIndex.SearchAnalyzer); err != nil {
		return diag.Errorf("error setting `searchAnalyzer` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("mappings_dynamic", searchIndex.Mappings.Dynamic); err != nil {
		return diag.Errorf("error setting `mappings_dynamic` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("synonyms", flattenSearchIndexSynonyms(searchIndex.Synonyms)); err != nil {
		return diag.Errorf("error setting `synonyms` for search index (%s): %s", d.Id(), err)
	}

	if searchIndex.Mappings.Fields != nil {
		searchIndexMappingFields, err := marshallSearchIndexMappingsField(searchIndex.Mappings.Fields)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("mappings_fields", searchIndexMappingFields); err != nil {
			return diag.Errorf("error setting `mappings_fields` for for search index (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func flattenSearchIndexSynonyms(synonyms []admin.SearchSynonymMappingDefinition) []map[string]any {
	synonymsMap := make([]map[string]any, len(synonyms))
	for i, s := range synonyms {
		synonymsMap[i] = map[string]any{
			"name":              s.Name,
			"analyzer":          s.Analyzer,
			"source_collection": s.Source.Collection,
		}
	}
	return synonymsMap
}

func flattenSearchIndexSynonyms2(synonyms []map[string]any) []map[string]any {
	synonymsMap := make([]map[string]any, 0)

	for _, s := range synonyms {
		sourceCollection := s["source"].(map[string]any)
		synonym := map[string]any{
			"name":              s["name"],
			"analyzer":          s["analyzer"],
			"source_collection": sourceCollection["collection"],
		}
		synonymsMap = append(synonymsMap, synonym)
	}
	return synonymsMap
}

func marshallSearchIndexAnalyzers(fields []admin.ApiAtlasFTSAnalyzers) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	mappingFieldJSON, err := json.Marshal(fields)
	return string(mappingFieldJSON), err
}

func marshallSearchIndexAnalyzers2(fields []map[string]any) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	mappingFieldJSON, err := json.Marshal(fields)
	return string(mappingFieldJSON), err
}

func marshallSearchIndexMappingsField(fields map[string]any) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	mappingFieldJSON, err := json.Marshal(fields)
	return string(mappingFieldJSON), err
}

func resourceMongoDBAtlasSearchIndexCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	clusterName := d.Get("cluster_name").(string)

	indexMapping := unmarshalSearchIndexMappingFields(d.Get("mappings_fields").(string))

	searchIndexRequest := &matlas.SearchIndex{
		Analyzer:       d.Get("analyzer").(string),
		Analyzers:      unmarshalSearchIndexAnalyzersFields(d.Get("analyzers").(string)),
		CollectionName: d.Get("collection_name").(string),
		Database:       d.Get("database").(string),
		Mappings: &matlas.IndexMapping{
			Dynamic: d.Get("mappings_dynamic").(bool),
			Fields:  &indexMapping,
		},
		Name:           d.Get("name").(string),
		SearchAnalyzer: d.Get("search_analyzer").(string),
		Status:         d.Get("status").(string),
		Synonyms:       expandSearchIndexSynonyms(d),
	}

	dbSearchIndexRes, _, err := conn.Search.CreateIndex(ctx, projectID, clusterName, searchIndexRequest)
	if err != nil {
		return diag.Errorf("error creating index: %s", err)
	}
	if d.Get("wait_for_index_build_completion").(bool) {
		timeout := d.Timeout(schema.TimeoutCreate)
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"IN_PROGRESS", "MIGRATING"},
			Target:     []string{"STEADY"},
			Refresh:    resourceSearchIndexRefreshFunc(ctx, clusterName, projectID, dbSearchIndexRes.IndexID, conn),
			Timeout:    timeout,
			MinTimeout: 1 * time.Minute,
			Delay:      1 * time.Minute,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			d.SetId(encodeStateID(map[string]string{
				"project_id":   projectID,
				"cluster_name": clusterName,
				"index_id":     dbSearchIndexRes.IndexID,
			}))
			resourceMongoDBAtlasSearchIndexDelete(ctx, d, meta)
			d.SetId("")
			return diag.FromErr(fmt.Errorf("error creating index in cluster (%s): %s", clusterName, err))
		}
	}
	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"index_id":     dbSearchIndexRes.IndexID,
	}))

	return resourceMongoDBAtlasSearchIndexRead(ctx, d, meta)
}

func expandSearchIndexSynonyms(d *schema.ResourceData) []map[string]any {
	var synonymsList []map[string]any

	synonymsDoc := map[string]any{}

	if vSynonyms, vSynonymsOK := d.GetOk("synonyms"); vSynonymsOK {
		for _, s := range vSynonyms.(*schema.Set).List() {
			synonym := s.(map[string]any)

			synonymsDoc["name"] = synonym["name"]
			synonymsDoc["analyzer"] = synonym["analyzer"]
			synonymsDoc["source"] = map[string]any{
				"collection": synonym["source_collection"],
			}
			synonymsList = append(synonymsList, synonymsDoc)
		}
	}
	return synonymsList
}

func validateSearchIndexMappingDiff(k, old, newStr string, d *schema.ResourceData) bool {
	var j, j2 any

	if old == "" {
		old = "{}"
	}

	if newStr == "" {
		newStr = "{}"
	}

	if err := json.Unmarshal([]byte(old), &j); err != nil {
		log.Printf("[ERROR] cannot unmarshal old search index mapping json %v", err)
	}
	if err := json.Unmarshal([]byte(newStr), &j2); err != nil {
		log.Printf("[ERROR] cannot unmarshal new search index mapping json %v", err)
	}
	if diff := deep.Equal(&j, &j2); diff != nil {
		log.Printf("[DEBUG] deep equal not passed: %v", diff)
		return false
	}

	return true
}

func validateSearchAnalyzersDiff(k, old, newStr string, d *schema.ResourceData) bool {
	var j, j2 any

	if old == "" {
		old = "{}"
	}

	if newStr == "" {
		newStr = "{}"
	}

	if err := json.Unmarshal([]byte(old), &j); err != nil {
		log.Printf("[ERROR] cannot unmarshal old search index analyzer json %v", err)
	}
	if err := json.Unmarshal([]byte(newStr), &j2); err != nil {
		log.Printf("[ERROR] cannot unmarshal new search index analyzer json %v", err)
	}
	if diff := deep.Equal(&j, &j2); diff != nil {
		log.Printf("[DEBUG] deep equal not passed: %v", diff)
		return false
	}

	return true
}

func unmarshalSearchIndexMappingFields(mappingString string) map[string]any {
	if mappingString == "" {
		return nil
	}

	var fields map[string]any

	if err := json.Unmarshal([]byte(mappingString), &fields); err != nil {
		log.Printf("[ERROR] cannot unmarshal search index mapping fields: %v", err)
		return nil
	}

	return fields
}

func unmarshalSearchIndexAnalyzersFields(mappingString string) []map[string]any {
	if mappingString == "" {
		return nil
	}

	var fields []map[string]any

	if err := json.Unmarshal([]byte(mappingString), &fields); err != nil {
		log.Printf("[ERROR] cannot unmarshal search index mapping fields: %v", err)
		return nil
	}

	return fields
}

func resourceSearchIndexRefreshFunc(ctx context.Context, clusterName, projectID, indexID string, client *matlas.Client) retry.StateRefreshFunc {
	return func() (any, string, error) {
		searchIndex, resp, err := client.Search.GetIndex(ctx, projectID, clusterName, indexID)
		if err != nil {
			return nil, "ERROR", err
		}

		if err != nil && searchIndex == nil && resp == nil {
			return nil, "", err
		} else if err != nil {
			if resp.StatusCode == 404 {
				return "", "DELETED", nil
			}
			if resp.StatusCode == 503 {
				return "", "PENDING", nil
			}
			return nil, "", err
		}

		if searchIndex.Status != "" {
			log.Printf("[DEBUG] status for Search Index : %s: %s", clusterName, searchIndex.Status)
		}

		return searchIndex, searchIndex.Status, nil
	}
}
