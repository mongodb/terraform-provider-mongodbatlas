package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
)

func main() {
	baseURL := os.Getenv("MONGODB_ATLAS_BASE_URL")
	clientID := os.Getenv("MONGODB_ATLAS_CLIENT_ID")
	clientSecret := os.Getenv("MONGODB_ATLAS_CLIENT_SECRET")
	if baseURL == "" || clientID == "" || clientSecret == "" {
		fmt.Fprintln(os.Stderr, "Error: MONGODB_ATLAS_BASE_URL, MONGODB_ATLAS_CLIENT_ID, and MONGODB_ATLAS_CLIENT_SECRET environment variables are required")
		os.Exit(1)
	}
	baseURL = strings.TrimRight(baseURL, "/")
	conf := clientcredentials.NewConfig(clientID, clientSecret)
	conf.TokenURL = baseURL + clientcredentials.TokenAPIPath
	token, err := conf.Token(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate OAuth2 token: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(token.AccessToken)
}
