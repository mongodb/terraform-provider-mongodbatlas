package searchindex

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20250312001/admin"
)

func PluralDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMongoDBAtlasSearchIndexesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"collection_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: returnSearchIndexDSSchema(),
				},
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceMongoDBAtlasSearchIndexesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	projectID, projectIDOK := d.GetOk("project_id")
	clusterName, clusterNameOk := d.GetOk("cluster_name")
	databaseName, databaseNameOK := d.GetOk("database")
	collectionName, collectionNameOK := d.GetOk("collection_name")

	if !projectIDOK || !clusterNameOk || !databaseNameOK || !collectionNameOK {
		return diag.Errorf("project_id, cluster_name, database and collection_name must be configured")
	}

	connV2 := meta.(*config.MongoDBClient).AtlasV2
	searchIndexes, _, err := connV2.AtlasSearchApi.ListAtlasSearchIndexes(ctx, projectID.(string), clusterName.(string), collectionName.(string), databaseName.(string)).Execute()

	if err != nil {
		return diag.Errorf("error getting search indexes information: %s", err)
	}

	flattedSearchIndexes, err := flattenSearchIndexes(searchIndexes, projectID.(string), clusterName.(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("results", flattedSearchIndexes); err != nil {
		return diag.Errorf("error setting `result` for search indexes: %s", err)
	}

	if err := d.Set("total_count", len(searchIndexes)); err != nil {
		return diag.Errorf("error setting `name`: %s", err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func flattenSearchIndexes(searchIndexes []admin.SearchIndexResponse, projectID, clusterName string) ([]map[string]any, error) {
	var searchIndexesMap []map[string]any

	if len(searchIndexes) == 0 {
		return nil, nil
	}
	searchIndexesMap = make([]map[string]any, len(searchIndexes))

	for i := range searchIndexes {
		searchIndexesMap[i] = map[string]any{
			"project_id":      projectID,
			"cluster_name":    clusterName,
			"analyzer":        searchIndexes[i].LatestDefinition.Analyzer,
			"collection_name": searchIndexes[i].CollectionName,
			"database":        searchIndexes[i].Database,
			"index_id":        searchIndexes[i].IndexID,
			"name":            searchIndexes[i].Name,
			"search_analyzer": searchIndexes[i].LatestDefinition.SearchAnalyzer,
			"status":          searchIndexes[i].Status,
			"synonyms":        flattenSearchIndexSynonyms(searchIndexes[i].LatestDefinition.GetSynonyms()),
			"type":            searchIndexes[i].Type,
		}

		if searchIndexes[i].LatestDefinition.Mappings != nil {
			searchIndexesMap[i]["mappings_dynamic"] = searchIndexes[i].LatestDefinition.Mappings.Dynamic

			if conversion.HasElementsSliceOrMap(searchIndexes[i].LatestDefinition.Mappings.Fields) {
				searchIndexMappingFields, err := marshalSearchIndex(searchIndexes[i].LatestDefinition.Mappings.Fields)
				if err != nil {
					return nil, err
				}
				searchIndexesMap[i]["mappings_fields"] = searchIndexMappingFields
			}
		}

		if analyzers := searchIndexes[i].LatestDefinition.GetAnalyzers(); len(analyzers) > 0 {
			searchIndexAnalyzers, err := marshalSearchIndex(analyzers)
			if err != nil {
				return nil, err
			}
			searchIndexesMap[i]["analyzers"] = searchIndexAnalyzers
		}

		if fields := searchIndexes[i].LatestDefinition.GetFields(); len(fields) > 0 {
			fieldsMarshaled, err := marshalSearchIndex(fields)
			if err != nil {
				return nil, err
			}
			searchIndexesMap[i]["fields"] = fieldsMarshaled
		}

		storedSource := searchIndexes[i].LatestDefinition.GetStoredSource()
		strStoredSource, errStoredSource := MarshalStoredSource(storedSource)
		if errStoredSource != nil {
			return nil, errStoredSource
		}
		searchIndexesMap[i]["stored_source"] = strStoredSource
	}
	return searchIndexesMap, nil
}
