package unit_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/require"
)

var projectAdvClusterExample = `
resource "mongodbatlas_project" "test" {
	org_id = "65def6ce0f722a1507105aa5"
	name   = "test-acc-tf-p-664077766951329406"
}

resource "mongodbatlas_advanced_cluster" "test" {
	project_id   = mongodbatlas_project.test.id
	name         = "test-acc-tf-c-8022584361920682288"
	cluster_type   = "REPLICASET"
	backup_enabled = false
	
	replication_specs {
		region_configs {
			provider_name = "AWS"
			priority      = 6
			region_name   = "US_WEST_2"
			electable_specs {
				node_count    = 1
				instance_size = "M10"
			}
		}

		region_configs {
			provider_name = "AWS"
			priority      = 7
			region_name   = "US_EAST_1"
			electable_specs {
				node_count    = 2
				instance_size = "M10"
			}
		}
	}
}`

func TestExtractConfigVariables(t *testing.T) {
	tests := map[string]struct {
		expected map[string]string
		config   string
	}{
		"Extract variables from a long example": {
			config: projectAdvClusterExample,
			expected: map[string]string{
				"orgId":       "65def6ce0f722a1507105aa5",
				"projectName": "test-acc-tf-p-664077766951329406",
				"clusterName": "test-acc-tf-c-8022584361920682288",
			},
		},
		"Extract variables from an empty config": {
			config:   "",
			expected: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := unit.ExtractConfigVariables(t, tc.config)
			require.Equal(t, tc.expected, result)
		})
	}
}
