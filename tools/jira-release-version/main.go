package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	jira "github.com/andygrunwald/go-jira/v2/onpremise"
)

const (
	envJira       = "JIRA_API_TOKEN"
	envVersion    = "VERSION_NUMBER"
	projectID     = "CLOUDP"
	jiraURL       = "https://jira.mongodb.org"
	versionPrefix = "terraform-provider-"
)

func main() {
	client := getJiraClient()
	versionName := versionPrefix + getVersion()
	versionID := getVersionID(client, versionName)
	fmt.Println(versionID)
}

func getJiraClient() *jira.Client {
	apiToken := os.Getenv(envJira)
	if apiToken == "" {
		log.Fatalf("Environment variable %s is required", envJira)
	}
	tp := jira.BearerAuthTransport{Token: apiToken}
	client, err := jira.NewClient(jiraURL, tp.Client())
	if err != nil {
		log.Fatalf("Error getting Jira client: %v", err)
	}
	return client
}

func getVersion() string {
	version := os.Getenv(envVersion)
	if version == "" {
		log.Fatalf("Environment variable %s is required", envVersion)
	}
	if strings.Contains(version, "pre") {
		fmt.Printf("Skipping release version for pre-release: %s\n", version)
		os.Exit(0)
	}
	return strings.TrimPrefix(version, "v")
}

func getVersionID(client *jira.Client, versionName string) string {
	ctx := context.Background()
	projects, _, err := client.Project.Get(ctx, projectID)
	if err != nil {
		log.Fatalf("Error getting project info: %v", err)
	}
	for i := range projects.Versions {
		v := &projects.Versions[i]
		if v.Name == versionName {
			return v.ID
		}
	}
	log.Fatalf("Version not found: %s", versionName)
	return ""
}
