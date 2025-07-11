package acc_test

import (
	"fmt"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
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

func TestConfigAddResourceStr(t *testing.T) {
	testCases := map[string]struct {
		hclConfig         string
		resourceID        string
		extraResourceStr  string
		expectedHCLConfig string
	}{
		"add tags": {
			hclConfig:        `resource "mongodbatlas_project" "test" {}`,
			resourceID:       "mongodbatlas_project.test",
			extraResourceStr: `tags = {}`,
			expectedHCLConfig: `resource "mongodbatlas_project" "test" {
  tags = {}
}
`},
		"add timeout on create": {
			hclConfig:  `resource "mongodbatlas_project" "test" {}`,
			resourceID: "mongodbatlas_project.test",
			extraResourceStr: `timeouts = { 
			create = "1m" 
		}`,
			expectedHCLConfig: `resource "mongodbatlas_project" "test" {
  timeouts = {
    create = "1m"
  }
}
`,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedHCLConfig, acc.ConfigAddResourceStr(t, tc.hclConfig, tc.resourceID, tc.extraResourceStr))
		})
	}
}
