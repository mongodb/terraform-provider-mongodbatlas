package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	switch command := os.Args[1]; command {
	case "create":
		if err := createToken(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "revoke":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: revoke command requires an access token as second argument")
			printUsage()
			os.Exit(1)
		}
		accessToken := os.Args[2]
		if err := revokeToken(accessToken); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  access-token create             # Generate a new OAuth2 access token")
	fmt.Fprintln(os.Stderr, "  access-token revoke <token>     # Revoke an existing OAuth2 access token")
}

func createToken() error {
	conf, err := getConfig()
	if err != nil {
		return err
	}
	token, err := conf.Token(context.Background())
	if err != nil {
		return fmt.Errorf("failed to generate OAuth2 token: %w", err)
	}
	accessToken := token.AccessToken
	if accessToken == "" {
		return fmt.Errorf("generated access token is empty")
	}
	return outputToken(accessToken)
}

func revokeToken(accessToken string) error {
	if accessToken == "" {
		return fmt.Errorf("access token cannot be empty")
	}
	conf, err := getConfig()
	if err != nil {
		return err
	}
	// OAuth2 revocation is always successful as per RFC 7009 for security and idempotency, even for invalid tokens.
	_ = conf.RevokeToken(context.Background(), &oauth2.Token{AccessToken: accessToken})
	return nil
}

func getConfig() (*clientcredentials.Config, error) {
	baseURL := strings.TrimRight(os.Getenv("MONGODB_ATLAS_BASE_URL"), "/")
	clientID := os.Getenv("MONGODB_ATLAS_CLIENT_ID")
	clientSecret := os.Getenv("MONGODB_ATLAS_CLIENT_SECRET")
	if baseURL == "" || clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("MONGODB_ATLAS_BASE_URL, MONGODB_ATLAS_CLIENT_ID, and MONGODB_ATLAS_CLIENT_SECRET environment variables are required")
	}
	conf := clientcredentials.NewConfig(clientID, clientSecret)
	conf.TokenURL = baseURL + clientcredentials.TokenAPIPath
	conf.RevokeURL = baseURL + clientcredentials.RevokeAPIPath
	return conf, nil
}

func outputToken(accessToken string) error {
	// Check if running in GitHub Actions.
	if githubOutput := os.Getenv("GITHUB_OUTPUT"); githubOutput != "" {
		return writeGitHubOutput(githubOutput, accessToken)
	}
	// Local usage: just print the token.
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
