package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
)

func TestIsIndexValue(t *testing.T) {
	assert.True(t, conversion.IsIndexValue(path.Root("replication_specs").AtListIndex(0)))
	assert.True(t, conversion.IsIndexValue(path.Root("replication_specs").AtMapKey("myKey")))
	assert.True(t, conversion.IsIndexValue(path.Root("replication_specs").AtSetValue(types.StringValue("myKey"))))
	assert.False(t, conversion.IsIndexValue(path.Root("replication_specs")))
	assert.False(t, conversion.IsIndexValue(path.Root("replication_specs").AtName("id")))
}

func TestAttributeNameEquals(t *testing.T) {
	var (
		repSpecPath       = path.Root("replication_specs")
		regionConfigsPath = repSpecPath.AtListIndex(0).AtName("region_configs")
	)
	for expectedAttribute, paths := range map[string][]path.Path{
		"replication_specs": {
			repSpecPath,
			repSpecPath.AtListIndex(0),
			repSpecPath.AtMapKey("myKey"),
		},
		"region_configs": {
			regionConfigsPath,
			regionConfigsPath.AtListIndex(0),
			regionConfigsPath.AtMapKey("myKey"),
		},
	} {
		for _, p := range paths {
			assert.True(t, conversion.AttributeNameEquals(p, expectedAttribute))
		}
	}
}

func TestStripSquareBrackets(t *testing.T) {
	assert.Equal(t, "replication_specs", conversion.TrimLastIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.Equal(t, "replication_specs", conversion.TrimLastIndex(path.Root("replication_specs").AtMapKey("myKey")))
	assert.Equal(t, "replication_specs", conversion.TrimLastIndex(path.Root("replication_specs")))
}

func TestIndexMethods(t *testing.T) {
	assert.True(t, conversion.IsListIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.False(t, conversion.IsListIndex(path.Root("replication_specs").AtName("region_configs")))
	assert.False(t, conversion.IsListIndex(path.Root("replication_specs").AtMapKey("region_configs")))
	assert.Equal(t, "replication_specs[+0]", conversion.AsAddedIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.Equal(t, "replication_specs[0].region_configs[+1]", conversion.AsAddedIndex(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1)))
	assert.Equal(t, "replication_specs[-1]", conversion.AsRemovedIndex(path.Root("replication_specs").AtListIndex(1)))
	assert.Equal(t, "replication_specs[0].region_configs[-1]", conversion.AsRemovedIndex(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1)))
}

func TestPathMatches(t *testing.T) {
	prefix := path.Root("replication_specs").AtListIndex(0)
	assert.True(t, conversion.HasPrefix(path.Root("replication_specs").AtListIndex(0), prefix))
	assert.True(t, conversion.HasPrefix(path.Root("replication_specs").AtListIndex(0).AtName("region_configs"), prefix))
	assert.False(t, conversion.HasPrefix(path.Root("replication_specs").AtListIndex(1), prefix))
	assert.True(t, conversion.HasPrefix(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1), path.Empty()))
}
func TestParentPathWithIndex_Found(t *testing.T) {
	diags := new(diag.Diagnostics)
	// Build a nested path: resource -> parent -> child
	basePath := path.Root("resource")
	parentPath := basePath.AtName("parent")
	childPath := parentPath.AtName("child")

	assert.Equal(t, parentPath.String(), conversion.ParentPathWithIndex(childPath, "parent", diags).String())
	assert.Equal(t, basePath.String(), conversion.ParentPathWithIndex(childPath, "resource", diags).String())
	assert.Empty(t, diags, "Diagnostics should not have errors")
}

func TestParentPathWithIndex_FoundIncludesIndex(t *testing.T) {
	diags := new(diag.Diagnostics)
	// Build a nested path: resource[0] -> parent[0] -> child
	basePath := path.Root("resource")
	parentPath := basePath.AtListIndex(0).AtName("parent")
	childPath := parentPath.AtListIndex(0).AtName("child")
	assert.Equal(t, "resource[0].parent[0].child", childPath.String())

	assert.Equal(t, parentPath.AtListIndex(0).String(), conversion.ParentPathWithIndex(childPath, "parent", diags).String())
	assert.Equal(t, basePath.AtListIndex(0).String(), conversion.ParentPathWithIndex(childPath, "resource", diags).String())
	assert.Empty(t, diags, "Diagnostics should not have errors")
}

func TestParentPathNoIndex_RemovesIndex(t *testing.T) {
	diags := new(diag.Diagnostics)
	// Build a nested path: resource[0] -> parent[0] -> child
	basePath := path.Root("resource")
	parentPath := basePath.AtListIndex(0).AtName("parent")
	childPath := parentPath.AtListIndex(0).AtName("child")
	assert.Equal(t, "resource[0].parent[0].child", childPath.String())

	assert.Equal(t, parentPath.String(), conversion.ParentPathNoIndex(childPath, "parent", diags).String())
	assert.Equal(t, basePath.String(), conversion.ParentPathNoIndex(childPath, "resource", diags).String())
	assert.Empty(t, diags, "Diagnostics should not have errors")
}

func TestParentPathWithIndex_NotFound(t *testing.T) {
	diags := new(diag.Diagnostics)
	// Build a path: resource -> child
	basePath := path.Root("resource")
	childPath := basePath.AtName("child")

	result := conversion.ParentPathWithIndex(childPath, "nonexistent", diags)
	// The function should traverse to path.Empty() and add an error.
	assert.True(t, result.Equal(path.Empty()), "Expected result to be empty if parent not found")
	assert.True(t, diags.HasError(), "Diagnostics should have an error when parent attribute is missing")
}

func TestParentPathWithIndex_EmptyPath(t *testing.T) {
	diags := new(diag.Diagnostics)
	emptyPath := path.Empty()
	result := conversion.ParentPathWithIndex(emptyPath, "any", diags)
	// Since the path is empty, it should immediately return empty and add error.
	assert.True(t, result.Equal(path.Empty()), "Expected empty path as result from an empty input path")
	assert.True(t, diags.HasError(), "Diagnostics should have an error for empty input path")
}
