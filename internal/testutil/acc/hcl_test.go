package acc_test

import (
	"fmt"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
)

func projectTemplateWithTags(tags string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	org_id 			 = "some_org"
	name   			 = "test-hcl"
	%s
}`, tags)
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

	tags := acc.MapToHcl(map[string]string{
		"Name":        "my-name",
		"Environment": "test",
	}, "\t", "tags")
	asserter.Equal(projectWithTags, projectTemplateWithTags(tags))
}
