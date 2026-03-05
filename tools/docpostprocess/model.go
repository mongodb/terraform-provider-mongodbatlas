package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/codespec"
	"github.com/mongodb/terraform-provider-mongodbatlas/tools/codegen/gofilegen/datasource"
	"gopkg.in/yaml.v3"
)

type docContext int

const (
	docResource docContext = iota
	docSingularDS
	docPluralDS
)

type docTarget struct {
	discriminators map[string]*codespec.Discriminator
	filePath       string
	isResource     bool
}

func loadModel(modelsDir, resourceName string) (*codespec.Resource, error) {
	path := filepath.Join(modelsDir, resourceName+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading model file %s: %w", path, err)
	}
	var resource codespec.Resource
	if err := yaml.Unmarshal(data, &resource); err != nil {
		return nil, fmt.Errorf("unmarshaling model %s: %w", path, err)
	}
	return &resource, nil
}

func loadAllModels(modelsDir string) ([]*codespec.Resource, error) {
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, fmt.Errorf("reading models directory %s: %w", modelsDir, err)
	}
	var resources []*codespec.Resource
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".yaml")
		resource, err := loadModel(modelsDir, name)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

// docTargetsForModel derives doc file paths and their discriminators from a loaded model.
func docTargetsForModel(resource *codespec.Resource, docsDir string) []docTarget {
	var targets []docTarget

	if discs := collectDiscriminators(resource, docResource); len(discs) > 0 {
		targets = append(targets, docTarget{
			filePath:       filepath.Join(docsDir, "resources", resource.Name+".md"),
			discriminators: discs,
			isResource:     true,
		})
	}

	if discs := collectDiscriminators(resource, docSingularDS); len(discs) > 0 {
		targets = append(targets, docTarget{
			filePath:       filepath.Join(docsDir, "data-sources", resource.Name+".md"),
			discriminators: discs,
			isResource:     false,
		})
	}

	if discs := collectDiscriminators(resource, docPluralDS); len(discs) > 0 {
		targets = append(targets, docTarget{
			filePath:       filepath.Join(docsDir, "data-sources", datasource.PluralizeName(resource.Name)+".md"),
			discriminators: discs,
			isResource:     false,
		})
	}

	return targets
}

func extractDiscriminatorForDoc(resource *codespec.Resource, ctx docContext) *codespec.Discriminator {
	switch ctx {
	case docResource:
		if resource.Schema != nil {
			return resource.Schema.Discriminator
		}
	case docSingularDS:
		if resource.DataSources != nil && resource.DataSources.Singular != nil {
			return resource.DataSources.Singular.Discriminator
		}
	case docPluralDS:
		if resource.DataSources != nil && resource.DataSources.Plural != nil {
			return resource.DataSources.Plural.Discriminator
		}
	}
	return nil
}

// collectDiscriminators returns a map from nested schema name to its discriminator.
// The empty string key "" holds the root-level discriminator.
func collectDiscriminators(resource *codespec.Resource, ctx docContext) map[string]*codespec.Discriminator {
	result := make(map[string]*codespec.Discriminator)

	rootDisc := extractDiscriminatorForDoc(resource, ctx)
	if rootDisc != nil {
		result[""] = rootDisc
	}

	var attrs codespec.Attributes
	switch ctx {
	case docResource:
		if resource.Schema != nil {
			attrs = resource.Schema.Attributes
		}
	case docSingularDS:
		if resource.DataSources != nil && resource.DataSources.Singular != nil {
			attrs = resource.DataSources.Singular.Attributes
		}
	case docPluralDS:
		if resource.DataSources != nil && resource.DataSources.Plural != nil {
			attrs = resource.DataSources.Plural.Attributes
		}
	}

	walkNestedDiscriminators(attrs, result)
	return result
}

func walkNestedDiscriminators(attrs codespec.Attributes, result map[string]*codespec.Discriminator) {
	for i := range attrs {
		attr := &attrs[i]
		nested := nestedObjectFromAttribute(attr)
		if nested == nil {
			continue
		}
		if nested.Discriminator != nil {
			result[attr.TFSchemaName] = nested.Discriminator
		}
		walkNestedDiscriminators(nested.Attributes, result)
	}
}

func nestedObjectFromAttribute(attr *codespec.Attribute) *codespec.NestedAttributeObject {
	switch {
	case attr.ListNested != nil:
		return &attr.ListNested.NestedObject
	case attr.SetNested != nil:
		return &attr.SetNested.NestedObject
	case attr.SingleNested != nil:
		return &attr.SingleNested.NestedObject
	case attr.MapNested != nil:
		return &attr.MapNested.NestedObject
	default:
		return nil
	}
}
