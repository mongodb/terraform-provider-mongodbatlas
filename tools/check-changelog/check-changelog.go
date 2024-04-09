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
	skipLabelName = "skip-changelog-check"
	skipTitles    = []string{"chore", "test", "doc", "ci"} // Dependabot uses chore.
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
		entry := changelog.Entry{
			Body: string(content),
		}
		if err := entry.Validate(); err != nil {
			log.Fatalf("Error validating changelog file: %s, err: %v", filePath, err)
		}
		fmt.Printf("Changelog file is valid: %s\n", filePath)
		return
	}

	if skipTitle(title) {
		fmt.Println("Skipping changelog check because PR title")
		return
	}

	if skipLabel(labels) {
		fmt.Printf("Skipping changelog check because PR label found: %s\n", skipLabelName)
		return
	}

	log.Fatalf("Consider using label %s if this PR doesn't need a changelog.\nChangelog file not found: %s, err: %v", skipLabelName, filePath, errFile)
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
