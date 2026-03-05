package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	resourceFlag := flag.String("resource", "", "Name of a single resource model to process (e.g. log_integration)")
	modelsDirFlag := flag.String("models-dir", "tools/codegen/models", "Path to the YAML model files directory")
	docsDirFlag := flag.String("docs-dir", "docs", "Base docs directory containing resources/ and data-sources/")
	flag.Parse()

	var targets []docTarget
	if *resourceFlag != "" {
		resource, err := loadModel(*modelsDirFlag, *resourceFlag)
		if err != nil {
			log.Fatalf("failed to load model for %s: %v", *resourceFlag, err)
		}
		targets = docTargetsForModel(resource, *docsDirFlag)
	} else {
		resources, err := loadAllModels(*modelsDirFlag)
		if err != nil {
			log.Fatalf("failed to load models: %v", err)
		}
		for _, resource := range resources {
			targets = append(targets, docTargetsForModel(resource, *docsDirFlag)...)
		}
	}

	for _, t := range targets {
		changed, err := processFile(t)
		if err != nil {
			log.Fatalf("failed to process %s: %v", t.filePath, err)
		}
		if changed {
			fmt.Printf("Post-processed: %s\n", t.filePath)
		}
	}
}

func processFile(t docTarget) (bool, error) {
	content, err := os.ReadFile(t.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Skipped (file not found): %s\n", t.filePath)
			return false, nil
		}
		return false, fmt.Errorf("reading file: %w", err)
	}

	original := string(content)

	result := RestructurePolymorphicDocs(original, t.discriminators, t.isResource)

	if result == original {
		return false, nil
	}

	if err := os.WriteFile(t.filePath, []byte(result), 0o600); err != nil {
		return false, fmt.Errorf("writing file: %w", err)
	}
	return true, nil
}
