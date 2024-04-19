package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	envJira    = "JIRA_API_TOKEN"
	envVersion = "VERSION_NUMBER"
)

func main() {
	apiToken := os.Getenv(envJira)
	if apiToken == "" {
		log.Fatalf("Environment variable %s is required", envJira)
	}
	version := os.Getenv(envVersion)
	if version == "" {
		log.Fatalf("Environment variable %s is required", envVersion)
	}
	fmt.Println("VERSION FOUND " + strings.TrimPrefix(version, "v"))
}
