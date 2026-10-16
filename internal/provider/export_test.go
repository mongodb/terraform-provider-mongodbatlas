package provider

import (
	"context"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

// GetAtlasCLICredentialsForTest exposes getAtlasCLICredentials for package provider_test.

func GetAtlasCLICredentialsForTest(ctx context.Context) (*config.Credentials, error) {
	return getAtlasCLICredentials(ctx)
}
