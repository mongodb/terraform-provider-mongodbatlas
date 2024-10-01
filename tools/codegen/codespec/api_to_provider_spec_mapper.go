//nolint:gocritic
package codespec

import (
	"log"
	// "net/url"

	// "github.com/pb33f/libopenapi"
	// v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
)

func ToProviderSpec(atlasAdminAPISpecFilePath, configPath string, resourceName *string) *CodeSpecification {

	apiDocModel, err := openapi.ParseAtlasAdminAPI(atlasAdminAPISpecFilePath)
	if err != nil {
		panic(err)
	}

	genConfig, _ := config.ParseGenConfigYAML(configPath)

	var resourceSpec config.Resource
	if resourceName != nil {
		resourceSpec = genConfig.Resources[*resourceName]
	}

	log.Print(resourceSpec)
	log.Print(apiDocModel)

	return nil
}

// func ToProviderSpec(atlasAdminAPISpecFilePath, configPath string, resourceName *string) *CodeSpecification {
// 	apiDocModel, err := openapi.ParseAtlasAdminAPI(atlasAdminAPISpecFilePath)
// 	if err != nil {
// 		panic(err)
// 	}

// 	genConfig, _ := config.ParseGenConfigYAML(configPath)

// 	var resourceSpec config.Resource
// 	if resourceName != nil {
// 		resourceSpec = genConfig.Resources[*resourceName]
// 	}

// 	log.Print(resourceSpec)

// 	// TODO: convert openAPIModel and config to  CodeSpecification

// 	// jsonBytes, err := json.MarshalIndent(apiDocModel, "", "  ")
// 	// if err != nil {
// 	// 	log.Fatalf("Error marshaling OpenAPI model to JSON: %v", err)
// 	// }
// 	// fmt.Println(string(jsonBytes))

// 	printRelevantOpenAPIParts(apiDocModel)

// 	// var createPath = apiDocModel.Paths.Find(genConfig.Resources["push_based_log_export"].Create.Path)
// 	// log.Println("path............................" + fmt.Sprint(createPath.Post.MarshalJSON()))

// 	return &CodeSpecification{}
// }

// func printRelevantOpenAPIParts(openAPIModel *openapi3.T) {
// 	if openAPIModel == nil {
// 		log.Println("The OpenAPI model is nil")
// 		return
// 	}

// 	relevantParts := &openapi3.T{
// 		OpenAPI:    openAPIModel.OpenAPI,
// 		Info:       openAPIModel.Info,
// 		Servers:    openAPIModel.Servers,
// 		Components: &openapi3.Components{Schemas: map[string]*openapi3.SchemaRef{}},
// 	}

// 	// Initialize Paths with capacity based on the potential filtered results
// 	relevantParts.Paths = openapi3.NewPaths()

// 	// Filter paths
// 	if openAPIModel.Paths != nil {
// 		for path, item := range openAPIModel.Paths.Map() {
// 			if relevantToPushBasedLogExport(path, item) {
// 				relevantParts.Paths.Set(path, item)
// 			}
// 		}
// 	}

// 	// Filter schemas
// 	if openAPIModel.Components != nil && openAPIModel.Components.Schemas != nil {
// 		for key, schema := range openAPIModel.Components.Schemas {
// 			if schemaContainsKeyword(schema, "PushBasedLogConfiguration") {
// 				relevantParts.Components.Schemas[key] = schema
// 			}
// 		}
// 	}

// 	// Serialize the filtered parts as JSON
// 	jsonBytes, err := json.MarshalIndent(relevantParts, "", "  ")
// 	if err != nil {
// 		log.Fatalf("Error marshaling relevant OpenAPI parts to JSON: %v", err)
// 	}
// 	fmt.Println(string(jsonBytes))
// }

// func relevantToPushBasedLogExport(path string, item *openapi3.PathItem) bool {
// 	return item != nil && (containsKeyword(&item.Summary, "PushBasedLogConfiguration") ||
// 		containsKeyword(&item.Description, "PushBasedLogConfiguration") ||
// 		pathOperationsContainKeyword(item, "PushBasedLogConfiguration"))
// }

// func pathOperationsContainKeyword(item *openapi3.PathItem, keyword string) bool {
// 	for _, operation := range item.Operations() {
// 		if containsKeyword(&operation.Summary, keyword) ||
// 			containsKeyword(&operation.Description, keyword) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func schemaContainsKeyword(schema *openapi3.SchemaRef, keyword string) bool {
// 	if schema == nil || schema.Value == nil {
// 		return false
// 	}
// 	return containsKeyword(&schema.Value.Description, keyword) ||
// 		containsKeyword(&schema.Value.Title, keyword)
// }

// func containsKeyword(text *string, keyword string) bool {
// 	if text == nil {
// 		return false
// 	}
// 	return strings.Contains(strings.ToLower(*text), strings.ToLower(keyword))
// }
