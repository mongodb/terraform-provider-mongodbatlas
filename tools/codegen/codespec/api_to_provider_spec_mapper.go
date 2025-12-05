package codespec

import (
	"errors"
	"fmt"
	"log"
	"strings"

	high "github.com/pb33f/libopenapi/datamodel/high/v3"
	low "github.com/pb33f/libopenapi/datamodel/low/v3"
	"github.com/pb33f/libopenapi/orderedmap"

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
		resource, ok := configModel.Resources[*resourceName]
		if !ok {
			log.Printf("[INFO] Resource %s not found in config file, skipping model generation", *resourceName)
			return &Model{}, nil
		}
		resourceConfigsToIterate = map[string]config.Resource{
			*resourceName: resource,
		}
	}

	if err := validateRequiredOperations(resourceConfigsToIterate); err != nil {
		return nil, err
	}

	var resources []Resource
	for name, resourceConfig := range resourceConfigsToIterate {
		log.Printf("[INFO] Generating resource model: %s", name)
		// find resource operations, schemas, etc from OAS
		oasResource, err := getAPISpecResource(&apiSpec.Model, &resourceConfig, name)
		if err != nil {
			return nil, fmt.Errorf("unable to get APISpecResource schema: %v", err)
		}
		// map OAS resource model to CodeSpecModel
		resource, err := apiSpecResourceToCodeSpecModel(oasResource, &resourceConfig, name)
		if err != nil {
			return nil, fmt.Errorf("unable to map to code spec model for %s: %w", name, err)
		}

		// Generate DataSources only when datasources block is defined in config
		if resourceConfig.DataSources != nil {
			dataSources, err := apiSpecToDataSourcesModel(&apiSpec.Model, &resourceConfig)
			if err != nil {
				return nil, fmt.Errorf("unable to map to data sources model for %s: %w", name, err)
			}
			resource.DataSources = dataSources
			log.Printf("[INFO] Generated data sources model for: %s", name)
		}

		resources = append(resources, *resource)
	}

	return &Model{Resources: resources}, nil
}

func validateRequiredOperations(resourceConfigs map[string]config.Resource) error {
	var validationErrors []error
	for name, resourceConfig := range resourceConfigs {
		if resourceConfig.Create == nil {
			validationErrors = append(validationErrors, fmt.Errorf("resource %s missing Create operation in config file", name))
		}
		if resourceConfig.Read == nil {
			validationErrors = append(validationErrors, fmt.Errorf("resource %s missing Read operation in config file", name))
		}
	}
	if len(validationErrors) > 0 {
		return errors.Join(validationErrors...)
	}
	return nil
}

func apiSpecResourceToCodeSpecModel(oasResource APISpecResource, resourceConfig *config.Resource, name string) (*Resource, error) {
	createOp := oasResource.CreateOp
	updateOp := oasResource.UpdateOp
	readOp := oasResource.ReadOp

	createPathParams := pathParamsToAttributes(createOp)
	var configuredVersion *string
	if resourceConfig.VersionHeader != "" {
		configuredVersion = &resourceConfig.VersionHeader
	}

	createRequestAttributes, err := opRequestToAttributes(createOp, configuredVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to process create request attributes for %s: %w", name, err)
	}
	updateRequestAttributes, err := opRequestToAttributes(updateOp, configuredVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to process update request attributes for %s: %w", name, err)
	}
	createResponseAttributes := opResponseToAttributes(createOp, configuredVersion)
	readResponseAttributes := opResponseToAttributes(readOp, configuredVersion)

	attributes := mergeAttributes(&attributeDefinitionSources{
		createPathParams: createPathParams,
		createRequest:    createRequestAttributes,
		updateRequest:    updateRequestAttributes,
		createResponse:   createResponseAttributes,
		readResponse:     readResponseAttributes,
	})

	schema := &Schema{
		Description:        oasResource.Description,
		DeprecationMessage: resourceConfig.DeprecationMessage,
		Attributes:         attributes,
	}

	operations := getOperationsFromConfig(resourceConfig)
	if operations.VersionHeader == "" { // version was not defined in config file
		operations.VersionHeader = getLatestVersionFromAPISpec(readOp)
	}

	var moveState *MoveState
	if resourceConfig.MoveState != nil {
		if len(resourceConfig.MoveState.SourceResources) == 0 {
			return nil, fmt.Errorf("resource %s missing source_resources for move_state in config file", name)
		}
		moveState = &MoveState{SourceResources: resourceConfig.MoveState.SourceResources}
	}

	resource := &Resource{
		Name:        name,
		PackageName: strings.ReplaceAll(name, "_", ""),
		Schema:      schema,
		MoveState:   moveState,
		Operations:  operations,
	}

	if err := applyTransformationsWithConfigOpts(resourceConfig, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

func getLatestVersionFromAPISpec(readOp *high.Operation) string {
	okResponse, ok := readOp.Responses.Codes.Get(OASResponseCodeOK)
	if !ok {
		return ""
	}
	versionsMap := okResponse.Content
	if versionsMap == nil {
		return ""
	}
	return orderedmap.SortAlpha(versionsMap).First().Key()
}

func getOperationsFromConfig(resourceConfig *config.Resource) APIOperations {
	return APIOperations{
		Create:        operationConfigToModel(resourceConfig.Create),
		Read:          operationConfigToModel(resourceConfig.Read),
		Update:        operationConfigToModel(resourceConfig.Update),
		Delete:        operationConfigToModel(resourceConfig.Delete),
		VersionHeader: resourceConfig.VersionHeader,
	}
}

func operationConfigToModel(opConfig *config.APIOperation) *APIOperation {
	if opConfig == nil {
		return nil
	}
	return &APIOperation{
		HTTPMethod:        opConfig.Method,
		Path:              opConfig.Path,
		Wait:              waitConfigToModel(opConfig.Wait),
		StaticRequestBody: opConfig.StaticRequestBody,
	}
}

func waitConfigToModel(waitConfig *config.Wait) *Wait {
	if waitConfig == nil {
		return nil
	}
	return &Wait{
		StateProperty:     waitConfig.StateProperty,
		PendingStates:     waitConfig.PendingStates,
		TargetStates:      waitConfig.TargetStates,
		TimeoutSeconds:    waitConfig.TimeoutSeconds,
		MinTimeoutSeconds: waitConfig.MinTimeoutSeconds,
		DelaySeconds:      waitConfig.DelaySeconds,
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
		parameterAttribute, err := s.buildResourceAttr(paramName, "", Required, false)
		if err != nil {
			log.Printf("[WARN] Path param %s could not be mapped: %s", paramName, err)
			continue
		}
		pathAttributes = append(pathAttributes, *parameterAttribute)
	}
	return pathAttributes
}

func opRequestToAttributes(op *high.Operation, configuredVersion *string) (Attributes, error) {
	if op == nil {
		return nil, nil
	}
	var requestAttributes Attributes
	requestSchema, err := buildSchemaFromRequest(op, configuredVersion)
	if err != nil {
		return nil, fmt.Errorf("request schema could not be mapped (OperationId: %s): %w", op.OperationId, err)
	}

	requestAttributes, err = buildResourceAttrs(requestSchema, "", true)
	if err != nil {
		return nil, fmt.Errorf("request attributes could not be mapped (OperationId: %s): %w", op.OperationId, err)
	}

	return requestAttributes, nil
}

func opResponseToAttributes(op *high.Operation, configuredVersion *string) Attributes {
	var responseAttributes Attributes
	responseSchema, err := buildSchemaFromResponse(op, configuredVersion)
	if err != nil {
		if errors.Is(err, errSchemaNotFound) {
			log.Printf("[INFO] Operation response body schema not found (OperationId: %s)", op.OperationId)
		} else {
			log.Printf("[WARN] Operation response body schema could not be mapped (OperationId: %s): %s", op.OperationId, err)
		}
	} else {
		responseAttributes, err = buildResourceAttrs(responseSchema, "", false)
		if err != nil {
			log.Printf("[WARN] Operation response body schema could not be mapped (OperationId: %s): %s", op.OperationId, err)
		}
	}
	return responseAttributes
}

func getAPISpecResource(spec *high.Document, resourceConfig *config.Resource, name string) (APISpecResource, error) {
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

// apiSpecToDataSourcesModel creates a DataSources model from the API spec using the datasources config.
// The data source has its own schema options (aliases, overrides, ignores) independent from the resource.
func apiSpecToDataSourcesModel(spec *high.Document, resourceConfig *config.Resource) (*DataSources, error) {
	dsConfig := resourceConfig.DataSources
	if dsConfig == nil {
		return nil, nil // no data source to generate
	}

	// Use resource's version header
	versionHeader := resourceConfig.VersionHeader
	var configuredVersion *string
	if versionHeader != "" {
		configuredVersion = &versionHeader
	}

	var attributes Attributes
	var readOp *APIOperation
	var listOp *APIOperation

	// Process Read operation if defined
	if dsConfig.Read != nil {
		oasReadOp, err := extractOp(spec.Paths, dsConfig.Read)
		if err != nil {
			return nil, fmt.Errorf("unable to extract data source read operation: %w", err)
		}

		// Build attributes from the read response
		readResponseAttributes := opResponseToAttributes(oasReadOp, configuredVersion)

		// Get path parameters as required attributes
		pathParams := pathParamsToAttributes(oasReadOp)

		// Merge path params with response attributes, considering aliases to avoid duplicates
		attributes = mergeDataSourceAttributes(pathParams, readResponseAttributes, dsConfig.SchemaOptions.Aliases)

		// Apply alias to path params in the read path
		readPath := applyAliasToPath(dsConfig.Read.Path, dsConfig.SchemaOptions.Aliases)

		readOp = &APIOperation{
			HTTPMethod: dsConfig.Read.Method,
			Path:       readPath,
		}

		// If version header wasn't explicitly set, get from API spec
		if versionHeader == "" {
			versionHeader = getLatestVersionFromAPISpec(oasReadOp)
		}
	}

	// Process List operation if defined
	if dsConfig.List != nil {
		listPath := applyAliasToPath(dsConfig.List.Path, dsConfig.SchemaOptions.Aliases)
		listOp = &APIOperation{
			HTTPMethod: dsConfig.List.Method,
			Path:       listPath,
		}
	}

	// Apply the data source's own schema transformations (aliases, overrides, ignores)
	if err := applyDataSourceTransformations(dsConfig, &attributes); err != nil {
		return nil, fmt.Errorf("failed to apply data source transformations: %w", err)
	}

	return &DataSources{
		Schema: &DataSourceSchema{
			Attributes: attributes,
		},
		Operations: APIOperations{
			Read:          readOp,
			List:          listOp,
			VersionHeader: versionHeader,
		},
	}, nil
}

// mergeDataSourceAttributes merges path parameters with response attributes.
// Path params are marked as required, response attributes are computed.
// Aliases are considered to avoid duplicates when a path param is aliased to match a response attribute.
func mergeDataSourceAttributes(pathParams, responseAttrs Attributes, aliases map[string]string) Attributes {
	// Create a map to track which attributes already exist from response (by APIName for alias matching)
	responseAttrByAPIName := make(map[string]bool)
	for i := range responseAttrs {
		responseAttrByAPIName[responseAttrs[i].APIName] = true
	}

	var result Attributes

	// Add path params as required (they identify the data source)
	for i := range pathParams {
		// Check if this path param will be aliased to match a response attribute
		// If so, the response attribute should be excluded to avoid duplication
		aliasedAPIName := pathParams[i].APIName
		if alias, ok := aliases[pathParams[i].APIName]; ok {
			aliasedAPIName = alias
		}

		// If the aliased param name matches a response attribute, mark it as handled
		if responseAttrByAPIName[aliasedAPIName] {
			responseAttrByAPIName[aliasedAPIName] = false // mark as handled
		}

		result = append(result, pathParams[i])
	}

	// Add remaining response attributes as computed
	for i := range responseAttrs {
		if responseAttrByAPIName[responseAttrs[i].APIName] {
			// Set as computed for data source
			responseAttrs[i].ComputedOptionalRequired = Computed
			responseAttrs[i].ReqBodyUsage = OmitAlways
			result = append(result, responseAttrs[i])
		}
	}

	return result
}

// applyDataSourceTransformations applies the data source's schema options to attributes.
func applyDataSourceTransformations(dsConfig *config.DataSources, attributes *Attributes) error {
	return applyAttributeTransformations(dsConfig.SchemaOptions, attributes, &attrPaths{schemaPath: "", apiPath: ""})
}

// applyAliasToPath replaces path parameters with their aliased names.
func applyAliasToPath(path string, aliases map[string]string) string {
	for original, alias := range aliases {
		path = strings.ReplaceAll(path, fmt.Sprintf("{%s}", original), fmt.Sprintf("{%s}", alias))
	}
	return path
}
