package searchindex

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func DataSource() *schema.Resource {
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
			Computed: true,
		},
		"type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"fields": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: validateSearchIndexMappingDiff,
		},
	}
}

func dataSourceMongoDBAtlasSearchIndexRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	projectID, projectIDOk := d.GetOk("project_id")
	clusterName, clusterNameOK := d.GetOk("cluster_name")
	indexID, indexIDOk := d.GetOk("index_id")

	if !(projectIDOk && clusterNameOK && indexIDOk) {
		return diag.Errorf("project_id, cluster_name and index_id must be configured")
	}

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	searchIndex, _, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(ctx, projectID.(string), clusterName.(string), indexID.(string)).Execute()
	if err != nil {
		return diag.Errorf("error getting search index information: %s", err)
	}

	if err := d.Set("type", searchIndex.Type); err != nil {
		return diag.Errorf("error setting `type` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("index_id", indexID); err != nil {
		return diag.Errorf("error setting `index_id` for search index (%s): %s", d.Id(), err)
	}

	// if err := d.Set("analyzer", searchIndex.Analyzer); err != nil {
	// 	return diag.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
	// }

	// if analyzers := searchIndex.GetAnalyzers(); len(analyzers) > 0 {
	// 	searchIndexMappingFields, err := marshalSearchIndex(analyzers)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}

	// 	if err := d.Set("analyzers", searchIndexMappingFields); err != nil {
	// 		return diag.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
	// 	}
	// }

	if err := d.Set("collection_name", searchIndex.CollectionName); err != nil {
		return diag.Errorf("error setting `collectionName` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("database", searchIndex.Database); err != nil {
		return diag.Errorf("error setting `database` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", searchIndex.Name); err != nil {
		return diag.Errorf("error setting `name` for search index (%s): %s", d.Id(), err)
	}

	// if err := d.Set("search_analyzer", searchIndex.SearchAnalyzer); err != nil {
	// 	return diag.Errorf("error setting `searchAnalyzer` for search index (%s): %s", d.Id(), err)
	// }

	// if err := d.Set("synonyms", flattenSearchIndexSynonyms(searchIndex.GetSynonyms())); err != nil {
	// 	return diag.Errorf("error setting `synonyms` for search index (%s): %s", d.Id(), err)
	// }

	// if searchIndex.Mappings != nil {
	// 	if err := d.Set("mappings_dynamic", searchIndex.Mappings.Dynamic); err != nil {
	// 		return diag.Errorf("error setting `mappings_dynamic` for search index (%s): %s", d.Id(), err)
	// 	}

	// 	if len(searchIndex.Mappings.Fields) > 0 {
	// 		searchIndexMappingFields, err := marshalSearchIndex(searchIndex.Mappings.Fields)
	// 		if err != nil {
	// 			return diag.FromErr(err)
	// 		}
	// 		if err := d.Set("mappings_fields", searchIndexMappingFields); err != nil {
	// 			return diag.Errorf("error setting `mappings_fields` for for search index (%s): %s", d.Id(), err)
	// 		}
	// 	}
	// }

	// if fields := searchIndex.GetFields(); len(fields) > 0 {
	// 	fieldsMarshaled, err := marshalSearchIndex(fields)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}

	// 	if err := d.Set("fields", fieldsMarshaled); err != nil {
	// 		return diag.Errorf("error setting `fields` for for search index (%s): %s", d.Id(), err)
	// 	}
	// }

	d.SetId(conversion.EncodeStateID(map[string]string{
		"project_id":   projectID.(string),
		"cluster_name": clusterName.(string),
		"index_id":     indexID.(string),
	}))

	return nil
}
