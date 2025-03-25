package conversion_test

import (
	"testing"

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
