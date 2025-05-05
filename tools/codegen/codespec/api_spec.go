package codespec

import (
	"context"
	"fmt"
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

	if len(schema.Type) == 0 {
		return nil, fmt.Errorf("invalid schema. no values for schema.Type found")
	}

	resp.Type = schema.Type[0]
	resp.Schema = schema

	return resp, nil
}

func GetSchemaFromMediaType(mediaTypes *orderedmap.Map[string, *high.MediaType]) (*APISpecSchema, error) {
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

	return GetSchemaFromMediaType(op.RequestBody.Content)
}

func buildSchemaFromResponse(op *high.Operation) (*APISpecSchema, error) {
	if op == nil || op.Responses == nil || op.Responses.Codes == nil || op.Responses.Codes.Len() == 0 {
		return nil, errSchemaNotFound
	}

	okResponse, ok := op.Responses.Codes.Get(OASResponseCodeOK)
	if ok {
		return GetSchemaFromMediaType(okResponse.Content)
	}

	createdResponse, ok := op.Responses.Codes.Get(OASResponseCodeCreated)
	if ok {
		return GetSchemaFromMediaType(createdResponse.Content)
	}

	sortedCodes := orderedmap.SortAlpha(op.Responses.Codes)
	for pair := range orderedmap.Iterate(context.Background(), sortedCodes) {
		responseCode := pair.Value()
		statusCode, err := strconv.Atoi(pair.Key())
		if err != nil {
			continue
		}

		if statusCode >= 200 && statusCode <= 299 {
			return GetSchemaFromMediaType(responseCode.Content)
		}
	}

	return nil, errSchemaNotFound
}
