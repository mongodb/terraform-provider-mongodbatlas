package mongodbatlas

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMongoDBAtlasSearchIndex() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceMongoDBAtlasSearchIndexRead,
		Schema: returnSearchIndexSchema(),
	}
}

func dataSourceMongoDBAtlasSearchIndexRead(d *schema.ResourceData, meta interface{}) error {
	// Get client connection.
	conn := meta.(*MongoDBClient).Atlas

	projectID, projectIDOk := d.GetOk("project_id")
	clusterName, clusterNameOK := d.GetOk("cluster_name")
	indexID, indexIDOk := d.GetOk("index_id")

	if !(projectIDOk && clusterNameOK && indexIDOk) {
		return errors.New("project_id, cluster_name and index_id must be configured")
	}

	searchIndex, _, err := conn.Search.GetIndex(context.Background(), projectID.(string), clusterName.(string), indexID.(string))
	if err != nil {
		return fmt.Errorf("error getting search index information: %s", err)
	}

	if err := d.Set("index_id", indexID); err != nil {
		return fmt.Errorf("error setting `index_id` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("analyzer", searchIndex.Analyzer); err != nil {
		return fmt.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
	}

	searchIndexCustomAnalyzers, err := flattenSearchIndexCustomAnalyzers(searchIndex.Analyzers)
	if err != nil {
		return nil
	}

	if err := d.Set("analyzers", searchIndexCustomAnalyzers); err != nil {
		return fmt.Errorf("error setting `analyzer` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("collection_name", searchIndex.CollectionName); err != nil {
		return fmt.Errorf("error setting `collectionName` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("database", searchIndex.Database); err != nil {
		return fmt.Errorf("error setting `database` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("name", searchIndex.Name); err != nil {
		return fmt.Errorf("error setting `name` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("search_analyzer", searchIndex.SearchAnalyzer); err != nil {
		return fmt.Errorf("error setting `searchAnalyzer` for search index (%s): %s", d.Id(), err)
	}

	if err := d.Set("mapping_dynamic", searchIndex.Mappings.Dynamic); err != nil {
		return fmt.Errorf("error setting `mapping_dynamic` for search index (%s): %s", d.Id(), err)
	}

	searchIndexMappingFields, err := marshallSearchIndexMappingFields(*searchIndex.Mappings.Fields)
	if err != nil {
		return err
	}
	if err := d.Set("mapping_fields", searchIndexMappingFields); err != nil {
		return fmt.Errorf("error setting `mapping_fields` for for search index (%s): %s", d.Id(), err)
	}

	return nil
}
