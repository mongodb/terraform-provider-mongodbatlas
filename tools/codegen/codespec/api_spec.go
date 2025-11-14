package codespec

import (
	"context"
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

var (
	errSchemaNotFound = fmt.Errorf("schema not found")
)

// This function only builds the schema from a proxy and returns the basic type and format without handling oneOf, anyOf, allOf, or nullable types.
func BuildSchema(proxy *base.SchemaProxy) (*APISpecSchema, error) {
	resp := &APISpecSchema{}
	schema, err := proxy.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to build schema from proxy: %w", err)
	}
	switch {
	case len(schema.Type) > 0:
		resp.Type = schema.Type[0]
	case schema.Properties != nil && schema.Properties.Len() > 0:
		// Infer object type when type is not explicitly defined but properties exist.
		// This handles cases like BaseSearchIndexCreateRequestDefinition and BaseSearchIndexResponseLatestDefinition which have properties but no explicit type.
		// This case can be removed after CLOUDP-355777 is done.
		schemaName := getSchemaName(proxy, schema)
		log.Printf("[WARN] Schema missing explicit type, inferring 'object' type from properties (schema: %s, properties: %d)", schemaName, schema.Properties.Len())
		resp.Type = OASTypeObject
	default:
		schemaName := getSchemaName(proxy, schema)
		return nil, fmt.Errorf("invalid schema. no values for schema.Type found and type cannot be inferred (schema: %s)", schemaName)
	}
	resp.Schema = schema
	return resp, nil
}

func getSchemaFromMediaType(mediaTypes *orderedmap.Map[string, *high.MediaType], configuredVersion *string) (*APISpecSchema, error) {
	if mediaTypes == nil {
		return nil, errSchemaNotFound
	}

	// If a configured version is provided, attempt to select the exact match or the closest older available version.
	if configuredVersion != nil {
		target := *configuredVersion

		// Find the closest version by comparing the Atlas date segment (exact match or closest older version)
		targetDate := extractAtlasDateFromMediaType(target)
		var bestKey string
		var bestDate string

		for pair := range orderedmap.Iterate(context.Background(), mediaTypes) {
			key := pair.Key()
			mt := pair.Value()
			if mt == nil || mt.Schema == nil {
				continue
			}
			keyDate := extractAtlasDateFromMediaType(key)
			if keyDate == "" || targetDate == "" {
				continue
			}
			// We want the greatest keyDate such that keyDate <= targetDate.
			if keyDate <= targetDate && (bestDate == "" || keyDate > bestDate) {
				bestKey = key
				bestDate = keyDate
			}
		}

		if bestKey != "" {
			if mt, ok := mediaTypes.Get(bestKey); ok {
				s, err := BuildSchema(mt.Schema)
				if err != nil {
					return nil, err
				}
				return s, nil
			}
		}
		// If no suitable version was found using the configuredVersion, fall through to default behavior.
	}

	sortedMediaTypes := orderedmap.SortAlpha(mediaTypes)
	if newest := sortedMediaTypes.Newest(); newest != nil {
		mt := newest.Value
		if mt != nil && mt.Schema != nil {
			s, err := BuildSchema(mt.Schema)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}

	return nil, errSchemaNotFound
}

// extractAtlasDateFromMediaType extracts the YYYY-MM-DD date segment from Atlas media types like:
// "application/vnd.atlas.2023-01-01+json". Returns empty string if not found.
func extractAtlasDateFromMediaType(mediaType string) string {
	const marker = "vnd.atlas."
	idx := strings.Index(mediaType, marker)
	if idx == -1 {
		return ""
	}
	start := idx + len(marker)
	rest := mediaType[start:]
	end := strings.IndexByte(rest, '+')
	if end == -1 {
		end = len(rest)
	}
	date := rest[:end]
	return date
}

func buildSchemaFromRequest(op *high.Operation) (*APISpecSchema, error) {
	if op == nil || op.RequestBody == nil || op.RequestBody.Content == nil || op.RequestBody.Content.Len() == 0 {
		return nil, errSchemaNotFound
	}

	return getSchemaFromMediaType(op.RequestBody.Content, nil)
}

func buildSchemaFromResponse(op *high.Operation) (*APISpecSchema, error) {
	if op == nil || op.Responses == nil || op.Responses.Codes == nil || op.Responses.Codes.Len() == 0 {
		return nil, errSchemaNotFound
	}

	okResponse, ok := op.Responses.Codes.Get(OASResponseCodeOK)
	if ok {
		return getSchemaFromMediaType(okResponse.Content, nil)
	}

	createdResponse, ok := op.Responses.Codes.Get(OASResponseCodeCreated)
	if ok {
		return getSchemaFromMediaType(createdResponse.Content, nil)
	}

	sortedCodes := orderedmap.SortAlpha(op.Responses.Codes)
	for pair := range orderedmap.Iterate(context.Background(), sortedCodes) {
		responseCode := pair.Value()
		statusCode, err := strconv.Atoi(pair.Key())
		if err != nil {
			continue
		}

		if statusCode >= 200 && statusCode <= 299 {
			return getSchemaFromMediaType(responseCode.Content, nil)
		}
	}

	return nil, errSchemaNotFound
}

// getSchemaName extracts a human-readable schema name from the proxy reference or schema title.
func getSchemaName(proxy *base.SchemaProxy, schema *base.Schema) string {
	if ref := proxy.GetReference(); ref != "" { // Try to get the name from the $ref.
		return path.Base(ref)
	}
	if schema.Title != "" { // Fall back to the schema title if available.
		return schema.Title
	}
	return "unknown"
}
