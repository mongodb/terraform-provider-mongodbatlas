package searchindex

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

const (
	vectorSearch = "vectorSearch"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCreate,
		ReadContext:          resourceRead,
		UpdateWithoutTimeout: resourceUpdate,
		DeleteContext:        resourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportState,
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
		},
		"analyzer": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"analyzers": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: diffSuppressJSON,
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
			DiffSuppressFunc: diffSuppressJSON,
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
			Computed: true,
		},
		"wait_for_index_build_completion": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"fields": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: diffSuppressJSON,
		},
		"stored_source": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: diffSuppressJSON,
		},
	}
}

func resourceImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "--", 3)
	if len(parts) != 3 {
		return nil, errors.New("import format error: to import a search index, use the format {project_id}--{cluster_name}--{index_id}")
	}

	projectID := parts[0]
	clusterName := parts[1]
	indexID := parts[2]

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	_, _, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
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

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"index_id":     indexID,
	}))

	return []*schema.ResourceData{d}, nil
}

func resourceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	_, err := connV2.AtlasSearchApi.DeleteAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	if err != nil {
		return diag.Errorf("error deleting search index (%s): %s", d.Get("name").(string), err)
	}
	return nil
}

func resourceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]
	indexName := d.Get("name").(string)

	if d.HasChange("name") || d.HasChange("type") || d.HasChange("database") || d.HasChange("collection_name") {
		return diag.Errorf("error updating search index (%s): attributes name, type, database and collection_name can't be updated", indexName)
	}

	searchRead, _, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	if err != nil {
		return diag.Errorf("error getting search index information: %s", err)
	}
	searchIndex := &admin.SearchIndexUpdateRequest{
		Definition: admin.SearchIndexUpdateRequestDefinition{
			Analyzer:       searchRead.LatestDefinition.Analyzer,
			Analyzers:      searchRead.LatestDefinition.Analyzers,
			Mappings:       searchRead.LatestDefinition.Mappings,
			SearchAnalyzer: searchRead.LatestDefinition.SearchAnalyzer,
			StoredSource:   searchRead.LatestDefinition.StoredSource,
			Synonyms:       searchRead.LatestDefinition.Synonyms,
			Fields:         searchRead.LatestDefinition.Fields,
		},
	}

	if d.HasChange("analyzer") {
		searchIndex.Definition.Analyzer = conversion.StringPtr(d.Get("analyzer").(string))
	}

	if d.HasChange("search_analyzer") {
		searchIndex.Definition.SearchAnalyzer = conversion.StringPtr(d.Get("search_analyzer").(string))
	}

	if d.HasChange("analyzers") {
		analyzers, err := UnmarshalSearchIndexAnalyzersFields(d.Get("analyzers").(string))
		if err != nil {
			return err
		}
		searchIndex.Definition.Analyzers = &analyzers
	}

	if d.HasChange("mappings_dynamic") {
		dynamic := d.Get("mappings_dynamic").(bool)
		if searchIndex.Definition.Mappings == nil {
			searchIndex.Definition.Mappings = &admin.SearchMappings{}
		}
		searchIndex.Definition.Mappings.Dynamic = &dynamic
	}

	if d.HasChange("mappings_fields") {
		mappingsFields, err := unmarshalSearchIndexMappingFields(d.Get("mappings_fields").(string))
		if err != nil {
			return err
		}
		if searchIndex.Definition.Mappings == nil {
			searchIndex.Definition.Mappings = &admin.SearchMappings{}
		}
		searchIndex.Definition.Mappings.Fields = &mappingsFields
	}

	if d.HasChange("fields") {
		fields, err := unmarshalSearchIndexFields(d.Get("fields").(string))
		if err != nil {
			return err
		}
		searchIndex.Definition.Fields = conversion.ToAnySlicePointer(&fields)
	}

	if d.HasChange("synonyms") {
		synonyms := expandSearchIndexSynonyms(d)
		searchIndex.Definition.Synonyms = &synonyms
	}

	if d.HasChange("stored_source") {
		obj, err := UnmarshalStoredSource(d.Get("stored_source").(string))
		if err != nil {
			return err
		}
		searchIndex.Definition.StoredSource = obj
	}

	if _, _, err := connV2.AtlasSearchApi.UpdateAtlasSearchIndex(ctx, projectID, clusterName, indexID, searchIndex).Execute(); err != nil {
		return diag.Errorf("error updating search index (%s): %s", indexName, err)
	}

	if d.Get("wait_for_index_build_completion").(bool) {
		timeout := d.Timeout(schema.TimeoutUpdate)
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING", "BUILDING", "IN_PROGRESS", "MIGRATING"},
			Target:     []string{"READY", "STEADY"},
			Refresh:    resourceSearchIndexRefreshFunc(ctx, clusterName, projectID, indexID, connV2),
			Timeout:    timeout,
			MinTimeout: 1 * time.Minute,
			Delay:      1 * time.Minute,
		}

		// Wait, catching any errors
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			d.SetId(conversion.EncodeStateID(map[string]string{
				"project_id":   projectID,
				"cluster_name": clusterName,
				"index_id":     indexID,
			}))
			return diag.FromErr(fmt.Errorf("error updating index in cluster (%s). mongodbatlas_search_index resource was not deleted : %s", clusterName, err))
		}
	}

	return resourceRead(ctx, d, meta)
}

func resourceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ids := conversion.DecodeStateID(d.Id())
	projectID := ids["project_id"]
	clusterName := ids["cluster_name"]
	indexID := ids["index_id"]

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	searchIndex, resp, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
	if err != nil {
		// deleted in the backend case
		if validate.StatusNotFound(resp) && !d.IsNewResource() {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error getting search index information: %s", err)
	}

	if err := d.Set("index_id", indexID); err != nil {
		return diag.Errorf("error setting `index_id` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("type", searchIndex.Type); err != nil {
		return diag.Errorf("error setting `type` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("analyzer", searchIndex.LatestDefinition.Analyzer); err != nil {
		return diag.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
	}

	if analyzers := searchIndex.LatestDefinition.GetAnalyzers(); len(analyzers) > 0 {
		searchIndexMappingFields, err := marshalSearchIndex(analyzers)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("analyzers", searchIndexMappingFields); err != nil {
			return diag.Errorf("error setting `analyzers` for search index (%s): %s", d.Id(), err)
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

	if err := d.Set("search_analyzer", searchIndex.LatestDefinition.SearchAnalyzer); err != nil {
		return diag.Errorf("error setting `searchAnalyzer` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("synonyms", flattenSearchIndexSynonyms(searchIndex.LatestDefinition.GetSynonyms())); err != nil {
		return diag.Errorf("error setting `synonyms` for search index (%s): %s", d.Id(), err)
	}

	if searchIndex.LatestDefinition.Mappings != nil {
		if err := d.Set("mappings_dynamic", searchIndex.LatestDefinition.Mappings.Dynamic); err != nil {
			return diag.Errorf("error setting `mappings_dynamic` for search index (%s): %s", d.Id(), err)
		}

		if fields := searchIndex.LatestDefinition.Mappings.Fields; fields != nil && conversion.HasElementsSliceOrMap(*fields) {
			searchIndexMappingFields, err := marshalSearchIndex(*fields)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("mappings_fields", searchIndexMappingFields); err != nil {
				return diag.Errorf("error setting `mappings_fields` for for search index (%s): %s", d.Id(), err)
			}
		}
	}

	if fields := searchIndex.LatestDefinition.GetFields(); len(fields) > 0 {
		fieldsMarshaled, err := marshalSearchIndex(fields)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("fields", fieldsMarshaled); err != nil {
			return diag.Errorf("error setting `fields` for for search index (%s): %s", d.Id(), err)
		}
	}

	storedSource := searchIndex.LatestDefinition.GetStoredSource()
	strStoredSource, errStoredSource := MarshalStoredSource(storedSource)
	if errStoredSource != nil {
		return diag.FromErr(errStoredSource)
	}
	if err := d.Set("stored_source", strStoredSource); err != nil {
		return diag.Errorf("error setting `stored_source` for search index (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	connV2 := meta.(*config.MongoDBClient).AtlasV2
	projectID := d.Get("project_id").(string)
	clusterName := d.Get("cluster_name").(string)
	indexType := d.Get("type").(string)
	searchIndexRequest := &admin.SearchIndexCreateRequest{
		Type:           conversion.StringPtr(indexType),
		CollectionName: d.Get("collection_name").(string),
		Database:       d.Get("database").(string),
		Name:           d.Get("name").(string),
		Definition: &admin.BaseSearchIndexCreateRequestDefinition{
			Analyzer:       conversion.StringPtr(d.Get("analyzer").(string)),
			SearchAnalyzer: conversion.StringPtr(d.Get("search_analyzer").(string)),
		},
	}

	if indexType == vectorSearch {
		fields, err := unmarshalSearchIndexFields(d.Get("fields").(string))
		if err != nil {
			return err
		}
		searchIndexRequest.Definition.Fields = conversion.ToAnySlicePointer(&fields)
	} else {
		analyzers, err := UnmarshalSearchIndexAnalyzersFields(d.Get("analyzers").(string))
		if err != nil {
			return err
		}
		searchIndexRequest.Definition.Analyzers = &analyzers
		mappingsFields, err := unmarshalSearchIndexMappingFields(d.Get("mappings_fields").(string))
		if err != nil {
			return err
		}
		dynamic := d.Get("mappings_dynamic").(bool)
		searchIndexRequest.Definition.Mappings = &admin.SearchMappings{
			Dynamic: &dynamic,
			Fields:  &mappingsFields,
		}
		synonyms := expandSearchIndexSynonyms(d)
		searchIndexRequest.Definition.Synonyms = &synonyms
	}

	objStoredSource, errStoredSource := UnmarshalStoredSource(d.Get("stored_source").(string))
	if errStoredSource != nil {
		return errStoredSource
	}
	searchIndexRequest.Definition.StoredSource = objStoredSource

	dbSearchIndexRes, _, err := connV2.AtlasSearchApi.CreateAtlasSearchIndex(ctx, projectID, clusterName, searchIndexRequest).Execute()
	if err != nil {
		return diag.Errorf("error creating index: %s", err)
	}
	indexID := conversion.SafeString(dbSearchIndexRes.IndexID)
	if d.Get("wait_for_index_build_completion").(bool) {
		timeout := d.Timeout(schema.TimeoutCreate)
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING", "BUILDING", "IN_PROGRESS", "MIGRATING"},
			Target:     []string{"READY", "STEADY"},
			Refresh:    resourceSearchIndexRefreshFunc(ctx, clusterName, projectID, indexID, connV2),
			Timeout:    timeout,
			MinTimeout: 1 * time.Minute,
			Delay:      1 * time.Minute,
		}

		// Wait, catching any errors
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			d.SetId(conversion.EncodeStateID(map[string]string{
				"project_id":   projectID,
				"cluster_name": clusterName,
				"index_id":     indexID,
			}))
			resourceDelete(ctx, d, meta)
			d.SetId("")
			return diag.FromErr(fmt.Errorf("error creating index in cluster (%s): %s", clusterName, err))
		}
	}
	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID,
		"cluster_name": clusterName,
		"index_id":     indexID,
	}))

	return resourceRead(ctx, d, meta)
}
