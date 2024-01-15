package searchindex

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"go.mongodb.org/atlas-sdk/v20231115003/admin"
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
			"page_num": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.15.0"),
			},
			"items_per_page": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: fmt.Sprintf(constant.DeprecationParamByVersion, "1.15.0"),
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

func dataSourceMongoDBAtlasSearchIndexesRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	projectID, projectIDOK := d.GetOk("project_id")
	clusterName, clusterNameOk := d.GetOk("cluster_name")
	databaseName, databaseNameOK := d.GetOk("database")
	collectionName, collectionNameOK := d.GetOk("collection_name")

	if !(projectIDOK && clusterNameOk && databaseNameOK && collectionNameOK) {
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

func flattenSearchIndexes(searchIndexes []admin.ClusterSearchIndex, projectID, clusterName string) ([]map[string]any, error) {
	var searchIndexesMap []map[string]any

	if len(searchIndexes) == 0 {
		return nil, nil
	}
	searchIndexesMap = make([]map[string]any, len(searchIndexes))

	for i := range searchIndexes {
		searchIndexesMap[i] = map[string]any{
			"project_id":      projectID,
			"cluster_name":    clusterName,
			"analyzer":        searchIndexes[i].Analyzer,
			"collection_name": searchIndexes[i].CollectionName,
			"database":        searchIndexes[i].Database,
			"index_id":        searchIndexes[i].IndexID,
			"name":            searchIndexes[i].Name,
			"search_analyzer": searchIndexes[i].SearchAnalyzer,
			"status":          searchIndexes[i].Status,
			"synonyms":        flattenSearchIndexSynonyms(searchIndexes[i].GetSynonyms()),
			"type":            searchIndexes[i].Type,
		}

		if searchIndexes[i].Mappings != nil {
			searchIndexesMap[i]["mappings_dynamic"] = searchIndexes[i].Mappings.Dynamic

			if len(searchIndexes[i].Mappings.Fields) > 0 {
				searchIndexMappingFields, err := marshalSearchIndex(searchIndexes[i].Mappings.Fields)
				if err != nil {
					return nil, err
				}
				searchIndexesMap[i]["mappings_fields"] = searchIndexMappingFields
			}
		}

		if analyzers := searchIndexes[i].GetAnalyzers(); len(analyzers) > 0 {
			searchIndexAnalyzers, err := marshalSearchIndex(analyzers)
			if err != nil {
				return nil, err
			}
			searchIndexesMap[i]["analyzers"] = searchIndexAnalyzers
		}

		if fields := searchIndexes[i].GetFields(); len(fields) > 0 {
			fieldsMarshaled, err := marshalSearchIndex(fields)
			if err != nil {
				return nil, err
			}
			searchIndexesMap[i]["fields"] = fieldsMarshaled
		}
	}

	return searchIndexesMap, nil
}
