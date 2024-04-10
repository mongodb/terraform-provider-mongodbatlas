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
		Name = "my-name"
		Environment = "test"
	}
}`

func TestMapToHcl(t *testing.T) {
	asserter := assert.New(t)

	tags := acc.HclMap(map[string]string{
		"Name":        "my-name",
		"Environment": "test",
	}, "\t", "tags")
	asserter.Equal(projectWithTags, projectTemplateWithExtra(tags))
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

func TestHclLifecycleIgnore(t *testing.T) {
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
			assert.Equal(t, tc.expected, projectTemplateWithExtra(acc.HclLifecycleIgnore(tc.keys...)))
		})
	}
}
