package provider

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

const atlasCLIServiceGov = "cloudgov"

// getAtlasCLICredentials returns credentials from a logged-in AtlasCLI session by
// invoking `atlas auth token`. base_url and is_mongodbgov_cloud are read from the same
// CLI profile so that the resolved configuration always matches the environment the token
// was minted for. It returns nil when AtlasCLI is not installed or not authenticated.
func getAtlasCLICredentials(ctx context.Context, profile string) (*config.Credentials, error) {
	// AtlasCLI lowercases profile names, but `config get` matches the name case-sensitively,
	// so normalize before using it as an argument to both `auth token` and `config get`.
	profile = strings.ToLower(profile)
	if profile == "" {
		profile = defaultAtlasCLIProfile(ctx)
	}
	out, err := atlasCLI(ctx, "auth", "token", "--profile", profile)
	if err != nil {
		return nil, nil //nolint:nilerr // AtlasCLI not installed or not logged in is not an error
	}
	token := strings.TrimSpace(out)
	if token == "" {
		return nil, nil
	}
	creds := &config.Credentials{AccessToken: token}
	if cliConfig, err := atlasCLIConfig(ctx, profile); err == nil {
		creds.BaseURL = cliConfig.OpsManagerURL
		creds.IsMongodbGovCloud = cliConfig.Service == atlasCLIServiceGov
	}
	return creds, nil
}

// defaultAtlasCLIProfile resolves the profile the same way the CLI does when none is
// configured explicitly: the only profile when a single one exists, otherwise "default".
func defaultAtlasCLIProfile(ctx context.Context) string {
	out, err := atlasCLI(ctx, "config", "list", "-o", "json")
	if err != nil {
		return "default"
	}
	var profiles []string
	if err := json.Unmarshal([]byte(out), &profiles); err != nil || len(profiles) != 1 {
		return "default"
	}
	return profiles[0]
}

type atlasCLIProfileConfig struct {
	OpsManagerURL string `json:"ops_manager_url"`
	Service       string `json:"service"`
}

func atlasCLIConfig(ctx context.Context, profile string) (*atlasCLIProfileConfig, error) {
	out, err := atlasCLI(ctx, "config", "get", profile, "-o", "json")
	if err != nil {
		return nil, err
	}
	cliConfig := &atlasCLIProfileConfig{}
	if err := json.Unmarshal([]byte(out), cliConfig); err != nil {
		return nil, err
	}
	return cliConfig, nil
}

func atlasCLI(ctx context.Context, args ...string) (string, error) {
	out, err := exec.CommandContext(ctx, "atlas", args...).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
