package mongodbatlas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-test/deep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func resourceMongoDBAtlasSearchIndex() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMongoDBAtlasSearchIndexCreate,
		ReadContext:   resourceMongoDBAtlasSearchIndexRead,
		UpdateContext: resourceMongoDBAtlasSearchIndexUpdate,
		DeleteContext: resourceMongoDBAtlasSearchIndexDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMongoDBAtlasSearchIndexImportState,
		},
		Schema: returnSearchIndexSchema(),
	}
}

func returnSearchIndexSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"cluster_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"index_id": {
			Type:     schema.TypeString,
			Computed: true,
			Required: false,
		},
		"analyzer": {
			Type:     schema.TypeString,
			Required: true,
		},
		"analyzers": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     customAnalyzersSchema(),
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
		"status": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func customAnalyzersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"char_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ignore_tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"mappings": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: validateSearchIndexMappingDiff,
						},
					},
				},
			},
			"tokenizer": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"max_token_length": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"min_gram": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_gram": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"group": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"token_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"original_tokens": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"normalization_form": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min_gram": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_gram": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"terms_not_in_bounds": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"min_shingle_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_shingle_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pattern": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"replacement": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"matches": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"stemmer_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tokens": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ignore_case": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMongoDBAtlasSearchIndexImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := meta.(*MongoDBClient).Atlas

	parts := strings.SplitN(d.Id(), "--", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a search index, use the format {project_id}-{cluster_name}-{index_id}")
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

func resourceMongoDBAtlasSearchIndexDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	_, err := conn.Search.DeleteIndex(ctx, projectID, clusterName, indexID)
	if err != nil {
		return diag.Errorf("error deleting search index (%s): %s", d.Get("name").(string), err)
	}

	return nil
}

func resourceMongoDBAtlasSearchIndexUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		searchIndex.Analyzers = expandCustomAnalyzers(d.Get("analyzers").([]interface{}))
	}

	if d.HasChange("collection_name") {
		searchIndex.CollectionName = d.Get("collectionName").(string)
	}

	if d.HasChange("database") {
		searchIndex.Database = d.Get("database").(string)
	}

	if d.HasChange("name") {
		searchIndex.Name = d.Get("name").(string)
	}

	if d.HasChange("search_analyzer") {
		searchIndex.SearchAnalyzer = d.Get("searchAnalyzer").(string)
	}

	if d.HasChange("mappings_dynamic") {
		searchIndex.Mappings.Dynamic = d.Get("mappings_dynamic").(bool)
	}

	if d.HasChange("mappings_fields") {
		var indexMapping *map[string]matlas.IndexField
		indexMappingResult := unmarshalSearchIndexMappingFields(d.Get("mappings_fields").(string))
		indexMapping = &indexMappingResult
		searchIndex.Mappings.Fields = indexMapping
	}

	searchIndex.IndexID = ""
	_, _, err = conn.Search.UpdateIndex(context.Background(), projectID, clusterName, indexID, searchIndex)
	if err != nil {
		return diag.Errorf("error updating search index (%s): %s", searchIndex.Name, err)
	}

	return resourceMongoDBAtlasSearchIndexRead(ctx, d, meta)
}

func resourceMongoDBAtlasSearchIndexRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	ids := decodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	searchIndex, _, err := conn.Search.GetIndex(ctx, projectID, clusterName, indexID)
	if err != nil {
		// case 404
		// deleted in the backend case
		reset := strings.Contains(err.Error(), "404") && !d.IsNewResource()

		if reset {
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

	searchIndexCustomAnalyzers, err := flattenSearchIndexCustomAnalyzers(searchIndex.Analyzers)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("analyzers", searchIndexCustomAnalyzers); err != nil {
		return diag.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
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

	if searchIndex.Mappings.Fields != nil {
		searchIndexMappingFields, err := marshallSearchIndexMappingFields(*searchIndex.Mappings.Fields)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("mappings_fields", searchIndexMappingFields); err != nil {
			return diag.Errorf("error setting `mappings_fields` for for search index (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func marshallSearchIndexMappingFields(fields map[string]matlas.IndexField) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}

	mappingFieldJSON, err := json.Marshal(fields)
	return string(mappingFieldJSON), err
}

func marshallSearchIndexCharFilterMappingFields(fields map[string]string) (interface{}, error) {
	if len(fields) == 0 {
		return "", nil
	}

	mappingFieldJSON, err := json.Marshal(fields)

	return string(mappingFieldJSON), err
}

func flattenSearchIndexCustomAnalyzers(analyzers []*matlas.CustomAnalyzer) ([]map[string]interface{}, error) {
	if len(analyzers) == 0 {
		return nil, nil
	}

	mapAnalyzers := make([]map[string]interface{}, len(analyzers))

	for i, analyzer := range analyzers {
		tokenizer, err := flattenSearchIndexTokenizer(analyzer.Tokenizer)
		if err != nil {
			return nil, err
		}

		mapAnalyzers[i] = map[string]interface{}{
			"name":      analyzer.Name,
			"tokenizer": tokenizer,
		}

		if len(analyzer.CharFilters) > 0 {
			searchIndexCharFilters, err := flattenSearchIndexCharFilters(analyzer.CharFilters)
			if err != nil {
				return nil, err
			}
			mapAnalyzers[i]["char_filters"] = searchIndexCharFilters
		}

		if len(analyzer.TokenFilters) > 0 {
			mapAnalyzers[i]["token_filters"] = flattenSearchIndexTokenFilters(analyzer.TokenFilters)
		}
	}
	return mapAnalyzers, nil
}

func flattenSearchIndexTokenizer(tokenizer *matlas.AnalyzerTokenizer) ([]map[string]interface{}, error) {
	tokenList := make([]map[string]interface{}, 0)

	mapTokenizer := map[string]interface{}{}

	if tokenizer.Type != "" {
		mapTokenizer["type"] = tokenizer.Type
	}

	if tokenizer.MaxTokenLength != nil {
		mapTokenizer["max_token_length"] = *tokenizer.MaxTokenLength
	}

	if tokenizer.MinGram != nil {
		mapTokenizer["min_gram"] = *tokenizer.MinGram
	}
	if tokenizer.MaxGram != nil {
		mapTokenizer["max_gram"] = *tokenizer.MaxGram
	}
	if tokenizer.Pattern != "" {
		mapTokenizer["pattern"] = tokenizer.Pattern
	}
	if tokenizer.Group != nil {
		mapTokenizer["group"] = *tokenizer.Group
	}

	tokenList = append(tokenList, mapTokenizer)

	return tokenList, nil
}

func flattenSearchIndexTokenFilters(filters []*matlas.AnalyzerTokenFilters) []map[string]interface{} {
	if len(filters) == 0 {
		return nil
	}

	mapCharFilters := make([]map[string]interface{}, len(filters))

	for i, filter := range filters {
		mapCharFilters[i] = map[string]interface{}{
			"type": filter.Type,
		}

		if filter.OriginalTokens != "" {
			mapCharFilters[i]["original_tokens"] = filter.OriginalTokens
		}
		//

		if filter.Min != nil {
			mapCharFilters[i]["min"] = *filter.Min
		}

		if filter.Max != nil {
			mapCharFilters[i]["max"] = *filter.Max
		}

		if filter.NormalizationForm != "" {
			mapCharFilters[i]["normalization_form"] = filter.NormalizationForm
		}

		if filter.MinGram != nil {
			mapCharFilters[i]["min_gram"] = *filter.MinGram
		}

		if filter.MaxGram != nil {
			mapCharFilters[i]["max_gram"] = *filter.MaxGram
		}

		if filter.TermsNotInBounds != "" {
			mapCharFilters[i]["terms_not_in_bounds"] = filter.TermsNotInBounds
		}

		if filter.MinShingleSize != nil {
			mapCharFilters[i]["min_shingle_size"] = *filter.MinShingleSize
		}

		if filter.MaxShingleSize != nil {
			mapCharFilters[i]["max_shingle_size"] = *filter.MaxShingleSize
		}

		if filter.Pattern != "" {
			mapCharFilters[i]["pattern"] = filter.Pattern
		}

		if filter.Replacement != "" {
			mapCharFilters[i]["replacement"] = filter.Replacement
		}

		if filter.Matches != "" {
			mapCharFilters[i]["matches"] = filter.Matches
		}

		if filter.StemmerName != "" {
			mapCharFilters[i]["stemmer_name"] = filter.StemmerName
		}

		if len(filter.Tokens) > 0 {
			mapCharFilters[i]["tokens"] = filter.Tokens
		}

		if filter.IgnoreCase != nil {
			mapCharFilters[i]["ignore_case"] = *filter.IgnoreCase
		}
	}
	return mapCharFilters
}

func flattenSearchIndexCharFilters(filters []*matlas.AnalyzerCharFilter) ([]map[string]interface{}, error) {
	if len(filters) == 0 {
		return nil, nil
	}

	mapCharFilters := make([]map[string]interface{}, len(filters))

	for i, filter := range filters {
		mapCharFilters[i] = map[string]interface{}{
			"type": filter.Type,
		}

		if len(filter.IgnoreTags) > 0 {
			mapCharFilters[i]["ignore_tags"] = filter.IgnoreTags
		}

		if filter.Mappings != nil {
			searchIndexCharFilterMappingFields, err := marshallSearchIndexCharFilterMappingFields(*filter.Mappings)
			if err != nil {
				return nil, err
			}

			mapCharFilters[i]["mappings"] = searchIndexCharFilterMappingFields
		}
	}
	return mapCharFilters, nil
}

func resourceMongoDBAtlasSearchIndexCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas
	projectID := d.Get("project_id").(string)

	clusterName := d.Get("cluster_name").(string)

	var indexMapping *map[string]matlas.IndexField
	indexMappingResult := unmarshalSearchIndexMappingFields(d.Get("mappings_fields").(string))
	indexMapping = &indexMappingResult

	searchIndexRequest := &matlas.SearchIndex{
		Analyzer:       d.Get("analyzer").(string),
		Analyzers:      expandCustomAnalyzers(d.Get("analyzers").([]interface{})),
		CollectionName: d.Get("collection_name").(string),
		Database:       d.Get("database").(string),
		Mappings: &matlas.IndexMapping{
			Dynamic: d.Get("mappings_dynamic").(bool),
			Fields:  indexMapping,
		},
		Name:           d.Get("name").(string),
		SearchAnalyzer: d.Get("search_analyzer").(string),
		Status:         d.Get("status").(string),
	}

	dbSearchIndexRes, _, err := conn.Search.CreateIndex(ctx, projectID, clusterName, searchIndexRequest)
	if err != nil {
		return diag.Errorf("error creating database user: %s", err)
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"index_id":     dbSearchIndexRes.IndexID,
	}))

	return resourceMongoDBAtlasSearchIndexRead(ctx, d, meta)
}

func expandCustomAnalyzers(analyzers []interface{}) []*matlas.CustomAnalyzer {
	if len(analyzers) == 0 {
		return nil
	}

	var analyzersList []*matlas.CustomAnalyzer

	for _, analyzerObj := range analyzers {
		analyzerInterface := analyzerObj.(map[string]interface{})

		analyzer := &matlas.CustomAnalyzer{
			Name: analyzerInterface["name"].(string),
		}

		charFiltersMap, ok := analyzerInterface["char_filters"]
		if ok {
			analyzer.CharFilters = expandIndexCharFilters(charFiltersMap.([]interface{}))
		}

		tokenizer, ok := analyzerInterface["tokenizer"]
		if ok {
			analyzer.Tokenizer = expandIndexTokenizer(tokenizer.([]interface{}))
		}

		tokenFiltersMap, ok := analyzerInterface["token_filters"]
		if ok {
			analyzer.TokenFilters = expandIndexTokenFilters(tokenFiltersMap.([]interface{}))
		}

		analyzersList = append(analyzersList, analyzer)
	}

	return analyzersList
}

func expandIndexTokenFilters(tokenFilters []interface{}) []*matlas.AnalyzerTokenFilters {
	var analyzerTokenFilters []*matlas.AnalyzerTokenFilters
	//TODO: here occurrs the change
	if len(tokenFilters) == 0 {
		return nil
	}

	for _, tf := range tokenFilters {
		tokenFilterMap := tf.(map[string]interface{})

		tokenFilter := &matlas.AnalyzerTokenFilters{}

		tokenFilter.Type = tokenFilterMap["type"].(string)

		if originalToken, ok := tokenFilterMap["original_tokens"]; ok {
			tokenFilter.OriginalTokens = originalToken.(string)
		}

		if min, ok := tokenFilterMap["min"].(int); ok && min > 0 {
			tokenFilter.Min = pointy.Int(min)
		}

		if max, ok := tokenFilterMap["max"].(int); ok && max > 0 {
			tokenFilter.Max = pointy.Int(max)
		}

		if normalizationForm, ok := tokenFilterMap["normalization_form"]; ok {
			tokenFilter.NormalizationForm = normalizationForm.(string)
		}

		if minGram, ok := tokenFilterMap["min_gram"].(int); ok && minGram > 0 {
			tokenFilter.MinGram = pointy.Int(minGram)
		}

		if maxGram, ok := tokenFilterMap["max_gram"].(int); ok && maxGram > 0 {
			tokenFilter.MinGram = pointy.Int(maxGram)
		}

		if termsNotInBounds, ok := tokenFilterMap["terms_not_in_bounds"]; ok {
			tokenFilter.TermsNotInBounds = termsNotInBounds.(string)
		}

		if minShingleSize, ok := tokenFilterMap["min_shingle_size"].(int); ok && minShingleSize > 0 {
			tokenFilter.MinShingleSize = pointy.Int(minShingleSize)
		}

		if maxShingleSize, ok := tokenFilterMap["max_shingle_size"].(int); ok && maxShingleSize > 0 {
			tokenFilter.MaxShingleSize = pointy.Int(maxShingleSize)
		}

		if pattern, ok := tokenFilterMap["pattern"]; ok {
			tokenFilter.Pattern = pattern.(string)
		}

		if replacement, ok := tokenFilterMap["replacement"]; ok {
			tokenFilter.Replacement = replacement.(string)
		}

		if matches, ok := tokenFilterMap["matches"]; ok {
			tokenFilter.Matches = matches.(string)
		}

		if stemmerName, ok := tokenFilterMap["stemmer_name"]; ok {
			tokenFilter.StemmerName = stemmerName.(string)
		}

		if tokens, ok := tokenFilterMap["tokens"]; ok {
			tokenFilter.Tokens = expandIndexTokens(tokens)
		}

		if ignoreCase, ok := tokenFilterMap["ignore_case"].(bool); ok {
			tokenFilter.IgnoreCase = pointy.Bool(ignoreCase)
		}

		analyzerTokenFilters = append(analyzerTokenFilters, tokenFilter)
	}

	return analyzerTokenFilters
}

func expandIndexTokens(tokens interface{}) []string {
	tokensInterfaces := tokens.([]interface{})
	tokensList := make([]string, len(tokensInterfaces))

	for i, token := range tokensInterfaces {
		tokensList[i] = token.(string)
	}
	return tokensList
}

func expandIndexTokenizer(tokenizers []interface{}) *matlas.AnalyzerTokenizer {
	if len(tokenizers) == 0 {
		return nil
	}

	analyzerTokenizer := &matlas.AnalyzerTokenizer{}

	tokenizer := tokenizers[0].(map[string]interface{})

	analyzerTokenizer.Type = tokenizer["type"].(string)

	if maxTokenLength, ok := tokenizer["max_token_length"].(int); ok && maxTokenLength > 0 {
		analyzerTokenizer.MaxTokenLength = pointy.Int(maxTokenLength)
	}

	if minGram, ok := tokenizer["min_gram"].(int); ok && minGram > 0 {
		analyzerTokenizer.MinGram = pointy.Int(minGram)
	}

	if maxGram, ok := tokenizer["max_gram"].(int); ok && maxGram > 0 {
		analyzerTokenizer.MaxGram = pointy.Int(maxGram)
	}

	if pattern, ok := tokenizer["pattern"]; ok {
		analyzerTokenizer.Pattern = pattern.(string)
	}

	if group, ok := tokenizer["group"].(int); ok {
		analyzerTokenizer.Group = pointy.Int(group)
	}

	return analyzerTokenizer
}

func expandIndexCharFilters(charFilters []interface{}) []*matlas.AnalyzerCharFilter {
	var analyzerCharFilters []*matlas.AnalyzerCharFilter

	if len(charFilters) == 0 {
		return nil
	}

	for _, tf := range charFilters {
		charFilterMap := tf.(map[string]interface{})

		charFilter := &matlas.AnalyzerCharFilter{
			Type: charFilterMap["type"].(string),
		}

		if ignoreTags, ok := charFilterMap["ignoreTags"].([]string); ok {
			charFilter.IgnoreTags = ignoreTags
		}

		if mappings, ok := charFilterMap["mappings"]; ok {
			var charFilterMapping *map[string]string
			charFilters := unmarshalSearchIndexCharFilterMapping(mappings.(string))
			charFilterMapping = &charFilters
			charFilter.Mappings = charFilterMapping
		}

		analyzerCharFilters = append(analyzerCharFilters, charFilter)
	}

	return analyzerCharFilters
}

func validateSearchIndexMappingDiff(k, old, newStr string, d *schema.ResourceData) bool {
	var j, j2 interface{}

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

func unmarshalSearchIndexMappingFields(mappingString string) map[string]matlas.IndexField {
	if mappingString == "" {
		return nil
	}

	var fields map[string]matlas.IndexField

	if err := json.Unmarshal([]byte(mappingString), &fields); err != nil {
		log.Printf("[ERROR] cannot unmarshal search index mapping fields: %v", err)
		return nil
	}

	return fields
}

func unmarshalSearchIndexCharFilterMapping(mappingString string) map[string]string {
	var fields map[string]string
	if err := json.Unmarshal([]byte(mappingString), &fields); err != nil {
		log.Printf("[ERROR] json.Unmarshal %v", err)
		return nil
	}
	return fields
}
