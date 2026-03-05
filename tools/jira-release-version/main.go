package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira/v2/onpremise"
)

const (
	envJira         = "JIRA_API_TOKEN"
	envVersion      = "VERSION_NUMBER"
	projectKey      = "CLOUDP"
	jiraURL         = "https://jira.mongodb.org"
	versionPrefix   = "terraform-provider-"
	versionNameNext = "next-terraform-provider-release"
)

func main() {
	client := getJiraClient()
	versionName := versionPrefix + getVersion()
	versionID := getOrCreateVersion(client, versionName)
	moveDoneIssues(client, versionName)
	setReleased(client, versionID)
	url := fmt.Sprintf("%s/projects/%s/versions/%s", jiraURL, projectKey, versionID)
	fmt.Printf("Version released, please check all tickets are marked as done: %s\n", url)
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
	return strings.TrimPrefix(version, "v")
}

func getOrCreateVersion(client *jira.Client, versionName string) string {
	var projectID int
	ctx := context.Background()
	projects, _, err := client.Project.Get(ctx, projectKey)
	if err != nil {
		log.Fatalf("Error getting project info: %v", err)
	}
	for i := range projects.Versions {
		v := &projects.Versions[i]
		if projectID == 0 {
			projectID = v.ProjectID
		}
		if v.Name == versionName {
			return v.ID
		}
	}

	version, _, err := client.Version.Create(ctx, &jira.Version{ProjectID: projectID, Name: versionName})
	if err != nil {
		log.Fatalf("Error creating version %s: %v", versionName, err)
	}
	fmt.Printf("Version not found so it has been created: %s, id: %s\n", versionName, version.ID)
	return version.ID
}

func moveDoneIssues(client *jira.Client, versionName string) {
	jql := fmt.Sprintf("project = %s AND status in (Resolved, Closed) AND fixVersion = %s", projectKey, versionNameNext)
	options := &jira.SearchOptions{MaxResults: 1000, Fields: []string{"NoFieldsNeeded"}}
	list, _, err := client.Issue.Search(context.Background(), jql, options)
	if err != nil {
		log.Fatalf("Error retrieving issues: %v", err)
	}
	keys := make([]string, len(list))
	for i := range list {
		key := list[i].Key
		keys[i] = key
		moveIssue(client, key, versionName)
	}
	if len(keys) > 0 {
		fmt.Println("Done issues moved:", strings.Join(keys, ", "))
	}
}

func moveIssue(client *jira.Client, issueKey, versionName string) {
	issue := &jira.Issue{
		Key: issueKey,
		Fields: &jira.IssueFields{
			FixVersions: []*jira.FixVersion{{Name: versionName}},
		},
	}
	_, _, err := client.Issue.Update(context.Background(), issue, nil)
	if err != nil {
		log.Fatalf("Error moving issue %s: %v", issueKey, err)
	}
}

func setReleased(client *jira.Client, versionID string) {
	version := &jira.Version{
		ID:          versionID,
		Released:    new(true),
		ReleaseDate: time.Now().UTC().Format("2006-01-02"),
	}
	_, _, err := client.Version.Update(context.Background(), version)
	if err != nil {
		log.Fatalf("Error updating version %s: %v", versionID, err)
	}
}
