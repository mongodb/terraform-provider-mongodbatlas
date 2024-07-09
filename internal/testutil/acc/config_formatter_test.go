package acc_test

import (
	"fmt"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func projectTemplateWithExtra(extra string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"
%s
}`, extra)
}

var projectWithTags = `
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"
	tags = {
		Environment = "test"
		Name = "my-name"
	}
}`
var projectWithEmptyTags = `
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"
	tags = {
	}
}`
var projectWithoutTags = `
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"

}`

func TestFormatToHCLMap(t *testing.T) {
	testCases := map[string]struct {
		values   map[string]string
		expected string
	}{
		"normal map": {map[string]string{
			"Name":        "my-name",
			"Environment": "test",
		}, projectWithTags},
		"empty map": {map[string]string{}, projectWithEmptyTags},
		"nil map":   {nil, projectWithoutTags},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tags := acc.FormatToHCLMap(tc.values, "\t", "tags")
			assert.Equal(t, tc.expected, projectTemplateWithExtra(tags))
		})
	}
}

var projectWithEmptyLifecycleIgnore = `
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"

}`
var projectWithLifecycleIgnoreSingle = `
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"
	lifecycle {
		ignore_changes = [
			tags["Name"],
		]
	}
}`
var projectWithLifecycleIgnoreMultiple = `
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"
	lifecycle {
		ignore_changes = [
			tags["Name"],
			tags["Env"],
		]
	}
}`

func TestFormatToHCLLifecycleIgnore(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
		keys     []string
	}{
		{"empty", projectWithEmptyLifecycleIgnore, []string{}},
		{"single", projectWithLifecycleIgnoreSingle, []string{`tags["Name"]`}},
		{"plural", projectWithLifecycleIgnoreMultiple, []string{`tags["Name"]`, `tags["Env"]`}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, projectTemplateWithExtra(acc.FormatToHCLLifecycleIgnore(tc.keys...)))
		})
	}
}

var standardClusterResource = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  project_id     = "project"

  replication_specs {
    num_shards = 1
    zone_name  = "zone1"

    region_configs {
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_1"
      auto_scaling {
        disk_gb_enabled = false
      }
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }

}
`

var dependsOnClusterResource = `
resource "mongodbatlas_advanced_cluster" "cluster_info" {
  backup_enabled = false
  cluster_type   = "REPLICASET"
  name           = "my-name"
  project_id     = "project"

  replication_specs {
    num_shards = 1
    zone_name  = "zone1"

    region_configs {
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_WEST_1"
      auto_scaling {
        disk_gb_enabled = false
      }
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }

  depends_on = [mongodbatlas_project.project_execution]
}
`

func Test_ClusterResourceHcl(t *testing.T) {
	var (
		clusterName = "my-name"
		testCases   = map[string]struct {
			expected string
			req      acc.ClusterRequest
		}{
			"defaults": {
				standardClusterResource,
				acc.ClusterRequest{ClusterNameExplicit: clusterName},
			},
			"dependsOn": {
				dependsOnClusterResource,
				acc.ClusterRequest{ClusterNameExplicit: clusterName, ResourceDependencyName: "mongodbatlas_project.project_execution"},
			},
		}
	)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			config, actualClusterName, err := acc.ClusterResourceHcl("project", &tc.req, nil)
			require.NoError(t, err)
			assert.Equal(t, clusterName, actualClusterName)
			assert.Equal(t, tc.expected, config)
		})
	}
}
