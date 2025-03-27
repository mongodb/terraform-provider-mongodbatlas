package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
)

func TestIsIndexTypes(t *testing.T) {
	listIndexPath := path.Root("replication_specs").AtListIndex(0)
	mapIndexPath := path.Root("replication_specs").AtMapKey("myKey")
	setIndexPath := path.Root("replication_specs").AtSetValue(types.StringValue("myKey"))
	assert.True(t, conversion.IsListIndex(listIndexPath))
	assert.False(t, conversion.IsListIndex(setIndexPath))
	assert.False(t, conversion.IsListIndex(mapIndexPath))

	assert.True(t, conversion.IsSetIndex(setIndexPath))
	assert.False(t, conversion.IsSetIndex(mapIndexPath))
	assert.False(t, conversion.IsSetIndex(listIndexPath))

	assert.True(t, conversion.IsMapIndex(mapIndexPath))
	assert.False(t, conversion.IsMapIndex(setIndexPath))
	assert.False(t, conversion.IsMapIndex(listIndexPath))
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

func TestTrimLastIndex(t *testing.T) {
	assert.Equal(t, "replication_specs", conversion.TrimLastIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.Equal(t, "replication_specs", conversion.TrimLastIndex(path.Root("replication_specs").AtMapKey("myKey")))
	assert.Equal(t, "replication_specs[Value(\"myKey\")]", path.Root("replication_specs").AtSetValue(types.StringValue("myKey")).String())
	assert.Equal(t, "replication_specs", conversion.TrimLastIndex(path.Root("replication_specs").AtSetValue(types.StringValue("myKey"))))
}

func TestIndexMethods(t *testing.T) {
	assert.True(t, conversion.IsListIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.False(t, conversion.IsListIndex(path.Root("replication_specs").AtName("region_configs")))
	assert.False(t, conversion.IsListIndex(path.Root("replication_specs").AtMapKey("region_configs")))
	assert.Equal(t, "replication_specs[+0]", conversion.AsAddedIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.Equal(t, "replication_specs[0].region_configs[+1]", conversion.AsAddedIndex(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1)))
	assert.Equal(t, "replication_specs[+\"myKey\"]", conversion.AsAddedIndex(path.Root("replication_specs").AtMapKey("myKey")))
	assert.Equal(t, "replication_specs[+Value(\"myKey\")]", conversion.AsAddedIndex(path.Root("replication_specs").AtSetValue(types.StringValue("myKey"))))
	assert.Equal(t, "replication_specs[-1]", conversion.AsRemovedIndex(path.Root("replication_specs").AtListIndex(1)))
	assert.Equal(t, "replication_specs[0].region_configs[-1]", conversion.AsRemovedIndex(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1)))
	assert.Equal(t, "replication_specs[-\"myKey\"]", conversion.AsRemovedIndex(path.Root("replication_specs").AtMapKey("myKey")))
	assert.Equal(t, "replication_specs[-Value(\"myKey\")]", conversion.AsRemovedIndex(path.Root("replication_specs").AtSetValue(types.StringValue("myKey"))))
	assert.Equal(t, "advanced_configuration.custom_openssl_cipher_config_tls12[-Value(\"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384\")]", conversion.AsRemovedIndex(path.Root("advanced_configuration").AtName("custom_openssl_cipher_config_tls12").AtSetValue(types.StringValue("TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"))))
	assert.Equal(t, "", conversion.AsRemovedIndex(path.Root("replication_specs")))
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
