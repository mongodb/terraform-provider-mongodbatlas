package codespec

import (
	"errors"
	"fmt"
	"log"
	"strings"

	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	low "github.com/pb33f/libopenapi/datamodel/low/v3"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

func ToCodeSpecModel(atlasAdminAPISpecFilePath, configPath string, resourceName *string) (*Model, error) {
	apiSpec, err := openapi.ParseAtlasAdminAPI(atlasAdminAPISpecFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Atlas Admin API: %v", err)
	}

	configModel, err := config.ParseGenConfigYAML(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}

	resourceConfigsToIterate := configModel.Resources
	if resourceName != nil { // only generate a specific resource
		resourceConfigsToIterate = map[string]config.Resource{
			*resourceName: configModel.Resources[*resourceName],
		}
	}

	var results []Resource
	for name, resourceConfig := range resourceConfigsToIterate {
		log.Printf("Generating resource: %s", name)
		// find resource operations, schemas, etc from OAS
		oasResource, err := getAPISpecResource(&apiSpec.Model, &resourceConfig, SnakeCaseString(name))
		if err != nil {
			return nil, fmt.Errorf("unable to get APISpecResource schema: %v", err)
		}
		// map OAS resource model to CodeSpecModel
		results = append(results, *apiSpecResourceToCodeSpecModel(oasResource, &resourceConfig, SnakeCaseString(name)))
	}

	return &Model{Resources: results}, nil
}

func apiSpecResourceToCodeSpecModel(oasResource APISpecResource, resourceConfig *config.Resource, name SnakeCaseString) *Resource {
	createOp := oasResource.CreateOp
	readOp := oasResource.ReadOp

	pathParamAttributes := pathParamsToAttributes(createOp)
	createRequestAttributes := opRequestToAttributes(createOp)
	createResponseAttributes := opResponseToAttributes(createOp)
	readResponseAttributes := opResponseToAttributes(readOp)

	attributes := mergeAttributes(pathParamAttributes, createRequestAttributes, createResponseAttributes, readResponseAttributes)

	schema := &Schema{
		Description:        oasResource.Description,
		DeprecationMessage: oasResource.DeprecationMessage,
		Attributes:         attributes,
	}

	operations := obtainOperations(resourceConfig) // paths and version header is grabbed from config
	resource := &Resource{
		Name:       name,
		Schema:     schema,
		Operations: operations,
	}

	applyConfigSchemaOptions(resourceConfig, resource)

	return resource
}

func obtainOperations(resourceConfig *config.Resource) APIOperations {
	return APIOperations{
		CreatePath:    resourceConfig.Create.Path,
		ReadPath:      resourceConfig.Read.Path,
		UpdatePath:    resourceConfig.Update.Path,
		DeletePath:    resourceConfig.Delete.Path,
		VersionHeader: resourceConfig.VersionHeader,
	}
}

func pathParamsToAttributes(createOp *high.Operation) Attributes {
	pathParams := createOp.Parameters

	pathAttributes := Attributes{}
	for _, param := range pathParams {
		if param.In != OASPathParam {
			continue
		}

		s, err := BuildSchema(param.Schema)
		if err != nil {
			continue
		}

		paramName := param.Name
		s.Schema.Description = param.Description
		parameterAttribute, err := s.buildResourceAttr(paramName, Required)
		if err != nil {
			log.Printf("[WARN] Path param %s could not be mapped: %s", paramName, err)
			continue
		}
		pathAttributes = append(pathAttributes, *parameterAttribute)
	}
	return pathAttributes
}

func opRequestToAttributes(op *high.Operation) Attributes {
	var requestAttributes Attributes
	requestSchema, err := buildSchemaFromRequest(op)
	if err != nil {
		log.Printf("[WARN] Request schema could not be mapped (OperationId: %s): %s", op.OperationId, err)
		return nil
	}

	requestAttributes, err = buildResourceAttrs(requestSchema)
	if err != nil {
		log.Printf("[WARN] Request attributes could not be mapped (OperationId: %s): %s", op.OperationId, err)
		return nil
	}

	return requestAttributes
}

func opResponseToAttributes(op *high.Operation) Attributes {
	var responseAttributes Attributes
	responseSchema, err := buildSchemaFromResponse(op)
	if err != nil {
		if errors.Is(err, errSchemaNotFound) {
			log.Printf("[INFO] Operation response body schema not found (OperationId: %s)", op.OperationId)
		} else {
			log.Printf("[WARN] Operation response body schema could not be mapped (OperationId: %s): %s", op.OperationId, err)
		}
	} else {
		responseAttributes, err = buildResourceAttrs(responseSchema)
		if err != nil {
			log.Printf("[WARN] Operation response body schema could not be mapped (OperationId: %s): %s", op.OperationId, err)
		}
	}
	return responseAttributes
}

func getAPISpecResource(spec *high.Document, resourceConfig *config.Resource, name SnakeCaseString) (APISpecResource, error) {
	var errResult error
	var resourceDeprecationMsg *string

	createOp, err := extractOp(spec.Paths, resourceConfig.Create)
	if err != nil {
		errResult = errors.Join(errResult, fmt.Errorf("unable to extract '%s.create' operation: %w", name, err))
	}
	readOp, err := extractOp(spec.Paths, resourceConfig.Read)
	if err != nil {
		errResult = errors.Join(errResult, fmt.Errorf("unable to extract '%s.read' operation: %w", name, err))
	}
	updateOp, err := extractOp(spec.Paths, resourceConfig.Update)
	if err != nil {
		errResult = errors.Join(errResult, fmt.Errorf("unable to extract '%s.update' operation: %w", name, err))
	}
	deleteOp, err := extractOp(spec.Paths, resourceConfig.Delete)
	if err != nil {
		errResult = errors.Join(errResult, fmt.Errorf("unable to extract '%s.delete' operation: %w", name, err))
	}

	commonParameters, err := extractCommonParameters(spec.Paths, resourceConfig.Read.Path)
	if err != nil {
		errResult = errors.Join(errResult, fmt.Errorf("unable to extract '%s' common parameters: %w", name, err))
	}

	if readOp.Deprecated != nil && *readOp.Deprecated {
		resourceDeprecationMsg = conversion.StringPtr(DefaultDeprecationMsg)
	}

	return APISpecResource{
		Description:        &createOp.Description,
		DeprecationMessage: resourceDeprecationMsg,
		CreateOp:           createOp,
		ReadOp:             readOp,
		UpdateOp:           updateOp,
		DeleteOp:           deleteOp,
		CommonParameters:   commonParameters,
	}, errResult
}

func extractOp(paths *high.Paths, apiOp *config.APIOperation) (*high.Operation, error) {
	if apiOp == nil {
		return nil, nil
	}

	if paths == nil || paths.PathItems == nil || paths.PathItems.GetOrZero(apiOp.Path) == nil {
		return nil, fmt.Errorf("path '%s' not found in OpenAPI spec", apiOp.Path)
	}

	pathItem, _ := paths.PathItems.Get(apiOp.Path)

	return extractOpFromPathItem(pathItem, apiOp)
}

func extractOpFromPathItem(pathItem *high.PathItem, apiOp *config.APIOperation) (*high.Operation, error) {
	if pathItem == nil || apiOp == nil {
		return nil, fmt.Errorf("pathItem or apiOp cannot be nil")
	}

	switch strings.ToLower(apiOp.Method) {
	case low.PostLabel:
		return pathItem.Post, nil
	case low.GetLabel:
		return pathItem.Get, nil
	case low.PutLabel:
		return pathItem.Put, nil
	case low.DeleteLabel:
		return pathItem.Delete, nil
	case low.PatchLabel:
		return pathItem.Patch, nil
	default:
		return nil, fmt.Errorf("method '%s' not found at OpenAPI path '%s'", apiOp.Method, apiOp.Path)
	}
}

func extractCommonParameters(paths *high.Paths, path string) ([]*high.Parameter, error) {
	if paths.PathItems.GetOrZero(path) == nil {
		return nil, fmt.Errorf("path '%s' not found in OpenAPI spec", path)
	}

	pathItem, _ := paths.PathItems.Get(path)

	return pathItem.Parameters, nil
}
