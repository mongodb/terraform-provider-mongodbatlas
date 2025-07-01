package searchindex

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
	"go.mongodb.org/atlas-sdk/v20250312005/admin"
)

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

func expandSearchIndexSynonyms(d *schema.ResourceData) []admin.SearchSynonymMappingDefinition {
	var synonymsList []admin.SearchSynonymMappingDefinition
	if vSynonyms, ok := d.GetOk("synonyms"); ok {
		for _, s := range vSynonyms.(*schema.Set).List() {
			synonym := s.(map[string]any)
			synonymsDoc := admin.SearchSynonymMappingDefinition{
				Name:     synonym["name"].(string),
				Analyzer: synonym["analyzer"].(string),
				Source: admin.SynonymSource{
					Collection: synonym["source_collection"].(string),
				},
			}
			synonymsList = append(synonymsList, synonymsDoc)
		}
	}
	return synonymsList
}

func marshalSearchIndex(fields any) (string, error) {
	respBytes, err := json.Marshal(fields)
	return string(respBytes), err
}

func unmarshalSearchIndexMappingFields(str string) (map[string]any, diag.Diagnostics) {
	fields := map[string]any{}
	if str == "" {
		return fields, nil
	}
	if err := json.Unmarshal([]byte(str), &fields); err != nil {
		return nil, diag.Errorf("cannot unmarshal search index attribute `mappings_fields` because it has an incorrect format")
	}
	return fields, nil
}

func unmarshalSearchIndexFields(str string) ([]map[string]any, diag.Diagnostics) {
	fields := []map[string]any{}
	if str == "" {
		return fields, nil
	}
	if err := json.Unmarshal([]byte(str), &fields); err != nil {
		return nil, diag.Errorf("cannot unmarshal search index attribute `fields` because it has an incorrect format")
	}

	return fields, nil
}

func UnmarshalSearchIndexAnalyzersFields(str string) ([]admin.AtlasSearchAnalyzer, diag.Diagnostics) {
	fields := []admin.AtlasSearchAnalyzer{}
	if str == "" {
		return nil, nil // don't send analyzers field to Atlas if empty
	}
	dec := json.NewDecoder(bytes.NewReader([]byte(str)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&fields); err != nil {
		return nil, diag.Errorf("cannot unmarshal search index attribute `analyzers` because it has an incorrect format")
	}
	return fields, nil
}

func MarshalStoredSource(obj any) (string, error) {
	if obj == nil {
		return "", nil
	}
	if b, ok := obj.(bool); ok {
		return strconv.FormatBool(b), nil
	}
	respBytes, err := json.Marshal(obj)
	return string(respBytes), err
}

func UnmarshalStoredSource(str string) (any, diag.Diagnostics) {
	switch str {
	case "":
		return any(nil), nil
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		var obj any
		if err := json.Unmarshal([]byte(str), &obj); err != nil {
			return nil, diag.Errorf("cannot unmarshal search index attribute `stored_source` because it has an incorrect format")
		}
		return obj, nil
	}
}

func diffSuppressJSON(k, old, newStr string, d *schema.ResourceData) bool {
	return schemafunc.EqualJSON(old, newStr, "vector search index")
}

func resourceSearchIndexRefreshFunc(ctx context.Context, clusterName, projectID, indexID string, connV2 *admin.APIClient) retry.StateRefreshFunc {
	return func() (any, string, error) {
		searchIndex, _, err := connV2.AtlasSearchApi.GetAtlasSearchIndex(ctx, projectID, clusterName, indexID).Execute()
		if err != nil {
			return nil, "ERROR", err
		}
		status := conversion.SafeString(searchIndex.Status)
		return searchIndex, status, nil
	}
}
