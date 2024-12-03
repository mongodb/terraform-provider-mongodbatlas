package config_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestLatestAdvancedClusterEnabled_notEnabled(t *testing.T) {
	t.Setenv(config.LatestAdvancedClusterEnabledEnvVar, "true")
	assert.False(t, config.LatestAdvancedClusterEnabled(), "LatestAdvancedClusterEnabled can't be enabled yet")
}
