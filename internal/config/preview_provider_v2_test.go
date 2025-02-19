package config_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAdvancedClusterV2Schema_notEnabled(t *testing.T) {
	t.Setenv(config.PreviewProviderV2AdvancedClusterEnvVar, "true")
	assert.False(t, config.PreviewProviderV2AdvancedCluster(), "PreviewProviderV2AdvancedCluster can't be enabled yet")
}
