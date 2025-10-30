package codespec

import (
	"context"
	"fmt"
	"log"
	"path"
	"strconv"

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

func getSchemaFromMediaType(mediaTypes *orderedmap.Map[string, *high.MediaType]) (*APISpecSchema, error) {
	if mediaTypes == nil {
		return nil, errSchemaNotFound
	}

	sortedMediaTypes := orderedmap.SortAlpha(mediaTypes)
	for pair := range orderedmap.Iterate(context.Background(), sortedMediaTypes) {
		mediaType := pair.Value()
		if mediaType.Schema != nil {
			s, err := BuildSchema(mediaType.Schema)
			if err != nil {
				return nil, err
			}
			return s, nil
		}
	}

	return nil, errSchemaNotFound
}

func buildSchemaFromRequest(op *high.Operation) (*APISpecSchema, error) {
	if op == nil || op.RequestBody == nil || op.RequestBody.Content == nil || op.RequestBody.Content.Len() == 0 {
		return nil, errSchemaNotFound
	}

	return getSchemaFromMediaType(op.RequestBody.Content)
}

func buildSchemaFromResponse(op *high.Operation) (*APISpecSchema, error) {
	if op == nil || op.Responses == nil || op.Responses.Codes == nil || op.Responses.Codes.Len() == 0 {
		return nil, errSchemaNotFound
	}

	okResponse, ok := op.Responses.Codes.Get(OASResponseCodeOK)
	if ok {
		return getSchemaFromMediaType(okResponse.Content)
	}

	createdResponse, ok := op.Responses.Codes.Get(OASResponseCodeCreated)
	if ok {
		return getSchemaFromMediaType(createdResponse.Content)
	}

	sortedCodes := orderedmap.SortAlpha(op.Responses.Codes)
	for pair := range orderedmap.Iterate(context.Background(), sortedCodes) {
		responseCode := pair.Value()
		statusCode, err := strconv.Atoi(pair.Key())
		if err != nil {
			continue
		}

		if statusCode >= 200 && statusCode <= 299 {
			return getSchemaFromMediaType(responseCode.Content)
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
