package acc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
)

func TestConvertToTPFAttrsMap(t *testing.T) {
	t.Setenv(config.AdvancedClusterV2SchemaEnvVar, "true")
	if !config.AdvancedClusterV2Schema() {
		t.Skip("Skipping test as not in AdvancedClusterV2Schema")
	}
	actual := map[string]string{
		"attr":                            "val1",
		"electable_specs.0":               "val2",
		"prefixbi_connector_config.0":     "val3",
		"advanced_configuration.0postfix": "val4",
		"electable_specs.0advanced_configuration.0bi_connector_config.0": "val5",
	}
	expected := map[string]string{
		"attr":                          "val1",
		"electable_specs":               "val2",
		"prefixbi_connector_config":     "val3",
		"advanced_configurationpostfix": "val4",
		"electable_specsadvanced_configurationbi_connector_config": "val5",
	}
	acc.ConvertToTPFAttrsMap(actual)
	assert.Equal(t, expected, actual)
}
