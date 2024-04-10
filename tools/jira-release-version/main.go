package main

import (
	"context"
	"fmt"
	"log"
	"os"

	jira "github.com/andygrunwald/go-jira/v2/onpremise"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
)

const (
	projectID = "CLOUDP"
	jiraURL   = "https://jira.mongodb.org"
)

func main() {
	client := getJiraClient()
	ctx := context.Background()

	_, _, err := client.Project.Get(ctx, projectID)
	if err != nil {
		log.Fatalf("Error getting project info: %v", err)
	}

	versionID := 39020
	_, _, err = client.Version.Get(ctx, versionID)
	if err != nil {
		log.Fatalf("Error getting version: %v", err)
	}

	updateVersion := &jira.Version{ID: "39020", Released: conversion.Pointer(true), ReleaseDate: "2024-04-09"}
	update, _, err := client.Version.Update(ctx, updateVersion)
	fmt.Println(update)
	if err != nil {
		log.Fatalf("Error updating version: %v", err)
	}
}

func getJiraClient() *jira.Client {
	apiToken := os.Getenv("JIRA_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Environment variable JIRA_API_TOKEN is required")
	}

	tp := jira.BearerAuthTransport{Token: apiToken}
	client, err := jira.NewClient(jiraURL, tp.Client())
	if err != nil {
		log.Fatalf("Error getting Jira client: %v", err)
	}
	return client
}
