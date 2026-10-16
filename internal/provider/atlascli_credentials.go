package provider

import (
	"context"
	"os/exec"
	"strings"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// getAtlasCLICredentials returns credentials from a logged-in AtlasCLI session by
// invoking `atlas auth token` and using the returned token as an access token.
// It returns nil when AtlasCLI is not installed, the command fails (non-zero exit),
// or the token is empty, so that callers can fall back to "no credentials".
func getAtlasCLICredentials(ctx context.Context) (*config.Credentials, error) {
	out, err := exec.CommandContext(ctx, "atlas", "auth", "token").Output()
	if err != nil {
		return nil, nil //nolint:nilerr // AtlasCLI not installed or not logged in is not an error
	}
	token := strings.TrimSpace(string(out))
	if token == "" {
		return nil, nil
	}
	return &config.Credentials{AccessToken: token}, nil
}
