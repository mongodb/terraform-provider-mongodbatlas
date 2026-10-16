package provider

import (
	"context"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// Helpers exposed for package provider_test.

func GetAtlasCLICredentialsForTest(ctx context.Context, profile string) (*config.Credentials, error) {
	return getAtlasCLICredentials(ctx, profile)
}

func DefaultAtlasCLIProfileForTest(ctx context.Context) string {
	return defaultAtlasCLIProfile(ctx)
}
