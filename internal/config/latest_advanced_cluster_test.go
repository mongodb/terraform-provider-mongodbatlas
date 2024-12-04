package config_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAdvancedClusterV2Schema_notEnabled(t *testing.T) {
	t.Setenv(config.AdvancedClusterV2SchemaEnvVar, "true")
	assert.False(t, config.AdvancedClusterV2Schema(), "AdvancedClusterV2Schema can't be enabled yet")
}
