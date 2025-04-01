package conversion_test

import (
	"testing"

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
	setIndex := path.Root("advanced_configuration").AtName("custom_openssl_cipher_config_tls12").AtSetValue(types.StringValue("TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"))
	assert.Equal(t, "advanced_configuration.custom_openssl_cipher_config_tls12[-Value(\"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384\")]", conversion.AsRemovedIndex(setIndex))
	setNoIndex, diags := conversion.AncestorPathNoIndex(setIndex, "custom_openssl_cipher_config_tls12")
	assert.Empty(t, diags)
	assert.Equal(t, "advanced_configuration.custom_openssl_cipher_config_tls12", setNoIndex.String())
	assert.Empty(t, conversion.AsRemovedIndex(path.Root("replication_specs")))
}

func TestHasAncestor(t *testing.T) {
	prefix := path.Root("replication_specs").AtListIndex(0)
	assert.True(t, conversion.HasAncestor(path.Root("replication_specs").AtListIndex(0), prefix))
	assert.True(t, conversion.HasAncestor(path.Root("replication_specs").AtListIndex(0).AtName("region_configs"), prefix))
	assert.False(t, conversion.HasAncestor(path.Root("replication_specs").AtListIndex(1), prefix))
	assert.True(t, conversion.HasAncestor(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1), path.Empty()))
}

func TestParentPathWithIndex_Found(t *testing.T) {
	// Build a nested path: resource -> parent -> child
	basePath := path.Root("resource")
	parentPath := basePath.AtName("parent")
	childPath := parentPath.AtName("child")

	parentPathActual, diags := conversion.AncestorPathWithIndex(childPath, "parent")
	assert.Empty(t, diags)
	assert.Equal(t, parentPath.String(), parentPathActual.String())

	basePathActual, diags := conversion.AncestorPathWithIndex(childPath, "resource")
	assert.Equal(t, basePath.String(), basePathActual.String())
	assert.Empty(t, diags, "Diagnostics should not have errors")
}

func TestParentPathWithIndex_FoundIncludesIndex(t *testing.T) {
	// Build a nested path: resource[0] -> parent[0] -> child
	basePath := path.Root("resource")
	parentPath := basePath.AtListIndex(0).AtName("parent")
	childPath := parentPath.AtListIndex(0).AtName("child")
	assert.Equal(t, "resource[0].parent[0].child", childPath.String())

	parentPathActual, diags := conversion.AncestorPathWithIndex(childPath, "parent")
	assert.Empty(t, diags)
	assert.Equal(t, parentPath.AtListIndex(0).String(), parentPathActual.String())

	basePathActual, diags := conversion.AncestorPathWithIndex(childPath, "resource")
	assert.Empty(t, diags)
	assert.Equal(t, basePath.AtListIndex(0).String(), basePathActual.String())
	assert.Empty(t, diags, "Diagnostics should not have errors")
}

func TestParentPathNoIndex_RemovesIndex(t *testing.T) {
	// Build a nested path: resource[0] -> parent[0] -> child
	basePath := path.Root("resource")
	parentPath := basePath.AtListIndex(0).AtName("parent")
	childPath := parentPath.AtListIndex(0).AtName("child")
	assert.Equal(t, "resource[0].parent[0].child", childPath.String())

	parentPathActual, diags := conversion.AncestorPathNoIndex(childPath, "parent")
	assert.Empty(t, diags)
	assert.Equal(t, parentPath.String(), parentPathActual.String())

	// Get base path without index
	basePathActual, diags := conversion.AncestorPathNoIndex(childPath, "resource")
	assert.Empty(t, diags)
	assert.Equal(t, basePath.String(), basePathActual.String())
}

func TestParentPathWithIndex_NotFound(t *testing.T) {
	// Build a path: resource -> child
	basePath := path.Root("resource")
	childPath := basePath.AtName("child")

	result, diags := conversion.AncestorPathWithIndex(childPath, "nonexistent")
	// The function should traverse to path.Empty() and add an error.
	assert.True(t, result.Equal(path.Empty()), "Expected result to be empty if parent not found")
	assert.True(t, diags.HasError(), "Diagnostics should have an error when parent attribute is missing")
}
