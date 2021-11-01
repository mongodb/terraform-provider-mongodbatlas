package mongodbatlas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMongoDBAtlasSearchIndex() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasSearchIndexRead,
		Schema:      returnSearchIndexDSSchema(),
	}
}

func returnSearchIndexDSSchema() map[string]*schema.Schema {
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
			Required: true,
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
			Optional: true,
		},
		"database": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
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
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"analyzer": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"source_collection": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"status": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func dataSourceMongoDBAtlasSearchIndexRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID, projectIDOk := d.GetOk("project_id")
	clusterName, clusterNameOK := d.GetOk("cluster_name")
	indexID, indexIDOk := d.GetOk("index_id")

	if !(projectIDOk && clusterNameOK && indexIDOk) {
		return diag.Errorf("project_id, cluster_name and index_id must be configured")
	}

	searchIndex, _, err := conn.Search.GetIndex(ctx, projectID.(string), clusterName.(string), indexID.(string))
	if err != nil {
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
		searchIndexMappingFields, err := marshallSearchIndexMappingsField(*searchIndex.Mappings.Fields)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("mappings_fields", searchIndexMappingFields); err != nil {
			return diag.Errorf("error setting `mappings_fields` for for search index (%s): %s", d.Id(), err)
		}
	}

	d.SetId(encodeStateID(map[string]string{
		"project_id":   projectID.(string),
		"cluster_name": clusterName.(string),
		"index_id":     indexID.(string),
	}))

	return nil
}
