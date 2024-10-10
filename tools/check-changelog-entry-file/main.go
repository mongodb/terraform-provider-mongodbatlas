package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-changelog"
)

var (
	skipLabelName              = "skip-changelog-check"
	skipTitles                 = []string{"chore", "test", "doc", "ci", "refactor"} // Dependabot uses chore.
	allowedTypeValues          = getValidTypes("scripts/changelog/allowed-types.txt")
	typesRequireResourcePrefix = []string{"breaking-change", "enhancement", "bug"}
)

func main() {
	var (
		title      = os.Getenv("PR_TITLE")
		number     = os.Getenv("PR_NUMBER")
		jsonLabels = os.Getenv("PR_LABELS")
	)
	if title == "" || number == "" || jsonLabels == "" {
		log.Fatal("Environment variables PR_TITLE, PR_NUMBER and PR_LABELS are required")
	}
	var labels []string
	if err := json.Unmarshal([]byte(jsonLabels), &labels); err != nil {
		log.Fatalf("PR_LABELS is not a stringified JSON array: %v", err)
	}

	filePath := fmt.Sprintf(".changelog/%s.txt", number)
	content, errFile := os.ReadFile(filePath)
	if errFile == nil { // Always validate changelog file if present, skip logic is not considered in this case
		validateChangelog(filePath, string(content))
		return
	}

	if skipTitle(title) {
		fmt.Println("Skipping changelog entry file check because PR title")
		return
	}

	if skipLabel(labels) {
		fmt.Printf("Skipping changelog entry file check because PR label found: %s\n", skipLabelName)
		return
	}

	log.Fatalf("Have you ran the `make generate-changelog-entry` command?\nIf this PR doesn't need a changelog entry file, consider using label %s.\nRead contributing guides (https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/contributing/changelog-process.md) for more info.\nChangelog entry file %s not found due to the following reason: %v.", skipLabelName, filePath, errFile)
}

func validateChangelog(filePath, body string) {
	entry := changelog.Entry{
		Body: body,
	}
	// grabbing validation logic from https://github.com/hashicorp/go-changelog/blob/main/entry.go#L66, if entry types become configurable we can invoke entry.Validate() directly
	notes := changelog.NotesFromEntry(entry)

	if len(notes) < 1 {
		log.Fatalf("Error validating changelog file: %s, no changelog entry found", filePath)
	}

	var unknownTypes []string
	for _, note := range notes {
		if !containsType(note.Type, allowedTypeValues) {
			unknownTypes = append(unknownTypes, note.Type)
		}
	}
	if len(unknownTypes) > 0 {
		log.Fatalf("Error validating changelog file: %s. Unknown changelog types %v, please use only the configured changelog entry types %v", filePath, unknownTypes, allowedTypeValues)
	}

	validateEntryPrefix(notes, filePath)
	fmt.Printf("Changelog entry file is valid: %s\n", filePath)
}

func validateEntryPrefix(entries []changelog.Note, filePath string) {
	for _, entry := range entries {
		entryContent := entry.Body
		if containsType(entry.Type, typesRequireResourcePrefix) &&
			!strings.HasPrefix(entryContent, "resource/mongodbatlas_") && !strings.HasPrefix(entryContent, "data-source/mongodbatlas_") && !strings.HasPrefix(entryContent, "provider") {
			log.Fatalf("Error validating changelog file: %s. An incorrect prefix was found in the definition of the changelog entry. Please use the prefix `resource/mongodbatlas_`, `data-source/mongodbatlas_`, or `provider` accordingly.", filePath)
		}
	}
}

func containsType(entryType string, allowed []string) bool {
	for _, a := range allowed {
		if a == entryType {
			return true
		}
	}
	return false
}

func skipTitle(title string) bool {
	for _, item := range skipTitles {
		if strings.HasPrefix(title, item+":") {
			return true
		}
	}
	return false
}

func skipLabel(labels []string) bool {
	for _, label := range labels {
		if label == skipLabelName {
			return true
		}
	}
	return false
}

func getValidTypes(path string) []string {
	content, errFile := os.ReadFile(path)
	if errFile != nil {
		log.Fatalf("Error getting allowed entry types from %s", path)
	}
	lines := strings.Split(string(content), "\n")
	allowedTypes := lines[:len(lines)-1] // remove last element as it is an empty string
	return allowedTypes
}
