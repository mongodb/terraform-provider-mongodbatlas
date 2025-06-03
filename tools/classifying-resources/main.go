package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/openapi"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/stringcase"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type Resource struct {
	Name   string
	TypeDescription string
	Read   Operation
	Create *Operation
	Delete *Operation
	Update *Operation
	List   *Operation
}

type Operation struct {
	Path   string
	Method string
	Op     v3.Operation
}

func main() {
	specFile := "open-api-spec.yml"
	docModel, err := openapi.ParseAtlasAdminAPI(specFile)
	if err != nil {
		log.Fatalf("Failed to parse spec: %v", err)
	}

	operationIdToOperation := getOperationIdMap(docModel)
	resources := identifyResources(operationIdToOperation)

	// alphabetically sort resources by name
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].TypeDescription < resources[j].TypeDescription
	})

	// Print declarative-friendly resources
	fmt.Printf("%d resources identified, content is inferred exclusively from API Spec leveraging convention used for operationIds:\n", len(resources))
	for _, resource := range resources {
		fmt.Printf("\n\n- %s - %s:\n", "mongodbatlas_"+stringcase.FromCamelCase(resource.Name).SnakeCase()+"_api", resource.TypeDescription)

		opInfo := func(op *Operation) string {
			if op == nil {
				return "N/A"
			}
			return fmt.Sprintf("%s - %s - %s", op.Op.OperationId, op.Method, op.Path)
		}
		fmt.Printf("    read:   %s\n", opInfo(&resource.Read))
		fmt.Printf("    create: %s\n", opInfo(resource.Create))
		fmt.Printf("    delete: %s\n", opInfo(resource.Delete))
		fmt.Printf("    update: %s\n", opInfo(resource.Update))
		fmt.Printf("    list:   %s\n", opInfo(resource.List))
	}
}

func identifyResources(operationIdToOperation map[string]Operation) []Resource {
	// Map to track potential resources by name (derived from read operations)
	results := []Resource{}

	// First, identify all read operations to get potential resource names
	for opId, op := range operationIdToOperation {
		if len(opId) < 3 || opId[:3] != "get" {
			continue
		}
		resourceName := opId[3:]

		possibleResource := Resource{
			Name: resourceName,
			Read: op,
		}
		createOpId := "create" + resourceName
		deleteOpId := "delete" + resourceName
		updateOpId := "update" + resourceName
		listOpId := fmt.Sprintf("list%s", PluralizeNoun(resourceName))

		if op, exists := operationIdToOperation[createOpId]; exists {
			possibleResource.Create = &op
		} 
		if op, exists := operationIdToOperation[deleteOpId]; exists {
			possibleResource.Delete = &op
		} 
		if op, exists := operationIdToOperation[updateOpId]; exists {
			possibleResource.Update = &op
		}

		// If no associated create, delete, or update operations are found, skip this get operation
		if possibleResource.Create == nil && possibleResource.Delete == nil && possibleResource.Update == nil {
			continue
		}

		if op, exists := operationIdToOperation[listOpId]; exists {
			possibleResource.List = &op
		} else {
			possibleResource.List = nil // No list operation found
		}

		if possibleResource.Create != nil {
			if possibleResource.Delete != nil {
				possibleResource.TypeDescription = "Create and delete (standard)"
			} else {	
				possibleResource.TypeDescription = "Create but missing delete (inconsistent)"
			}	
		} else if possibleResource.Update != nil  {
			if possibleResource.Delete != nil {
				possibleResource.TypeDescription = "Singleton with delete"
			} else {
				possibleResource.TypeDescription = "Singleton"
			}
		} else {
			possibleResource.TypeDescription = "Read only (discarded)"
		}

		results = append(results, possibleResource)
	}
	return results
}

// cluster -> clusters
// index -> indexes
// policy -> policies
func PluralizeNoun(noun string) string {
	lower := strings.ToLower(noun)

	// Check for consonant+y ending
	matched, _ := regexp.MatchString(`[^aeiou]y$`, lower)

	switch {
	case strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "sh") ||
		strings.HasSuffix(lower, "s") ||
		strings.HasSuffix(lower, "x") ||
		strings.HasSuffix(lower, "z"):
		return noun + "es"

	case matched:
		return noun[:len(noun)-1] + "ies"

	case strings.HasSuffix(lower, "is"):
		return noun[:len(noun)-2] + "es"

	default:
		return noun + "s"
	}
}

func getOperationIdMap(doc *libopenapi.DocumentModel[v3.Document]) map[string]Operation {
	operationIdToOperation := map[string]Operation{}
	listOfPaths := doc.Model.Paths.PathItems.FromNewest()

	for path, item := range listOfPaths {
		if item.Get != nil {
			operationIdToOperation[item.Get.OperationId] = Operation{
				Path:   path,
				Method: "GET",
				Op:     *item.Get,
			}
		}
		if item.Post != nil {
			operationIdToOperation[item.Post.OperationId] = Operation{
				Path:   path,
				Method: "POST",
				Op:     *item.Post,
			}
		}
		if item.Put != nil {
			operationIdToOperation[item.Put.OperationId] = Operation{
				Path:   path,
				Method: "PUT",
				Op:     *item.Put,
			}
		}
		if item.Patch != nil {
			operationIdToOperation[item.Patch.OperationId] = Operation{
				Path:   path,
				Method: "PATCH",
				Op:     *item.Patch,
			}
		}
		if item.Delete != nil {
			operationIdToOperation[item.Delete.OperationId] = Operation{
				Path:   path,
				Method: "DELETE",
				Op:     *item.Delete,
			}
		}
	}
	return operationIdToOperation
}
