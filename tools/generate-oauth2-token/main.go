package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
)

func main() {
	baseURL := strings.TrimRight(os.Getenv("MONGODB_ATLAS_BASE_URL"), "/")
	clientID := os.Getenv("MONGODB_ATLAS_CLIENT_ID")
	clientSecret := os.Getenv("MONGODB_ATLAS_CLIENT_SECRET")
	if baseURL == "" || clientID == "" || clientSecret == "" {
		fmt.Fprintln(os.Stderr, "Error: MONGODB_ATLAS_BASE_URL, MONGODB_ATLAS_CLIENT_ID, and MONGODB_ATLAS_CLIENT_SECRET environment variables are required")
		os.Exit(1)
	}
	conf := clientcredentials.NewConfig(clientID, clientSecret)
	conf.TokenURL = baseURL + clientcredentials.TokenAPIPath
	token, err := conf.Token(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate OAuth2 token: %v\n", err)
		os.Exit(1)
	}

	accessToken := token.AccessToken
	if accessToken == "" {
		fmt.Fprintln(os.Stderr, "Error: Generated access token is empty")
		os.Exit(1)
	}

	if err := outputToken(accessToken); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func outputToken(accessToken string) error {
	// Check if running in GitHub Actions
	if githubOutput := os.Getenv("GITHUB_OUTPUT"); githubOutput != "" {
		return writeGitHubOutput(githubOutput, accessToken)
	}
	// Local usage: just print the token
	fmt.Print(accessToken)
	return nil
}

func writeGitHubOutput(githubOutput, accessToken string) error {
	file, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open GITHUB_OUTPUT file: %w", err)
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "access_token<<EOF\n%s\nEOF\n", accessToken); err != nil {
		return fmt.Errorf("failed to write to GITHUB_OUTPUT file: %w", err)
	}
	return nil
}
