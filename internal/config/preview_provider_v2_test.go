package config_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAdvancedClusterV2Schema_notEnabled(t *testing.T) {
	t.Setenv(config.PreviewProviderV2EnvVar, "true")
	assert.False(t, config.PreviewProviderV2(), "PreviewProviderV2 can't be enabled yet")
}
