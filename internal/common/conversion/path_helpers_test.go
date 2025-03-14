package conversion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
)

func TestIsAttributeValueOnly(t *testing.T) {
	assert.True(t, conversion.IsAttributeValueOnly(path.Root("replication_specs").AtListIndex(0)))
	assert.True(t, conversion.IsAttributeValueOnly(path.Root("replication_specs").AtMapKey("myKey")))
	assert.True(t, conversion.IsAttributeValueOnly(path.Root("replication_specs").AtSetValue(types.StringValue("myKey"))))
}

func TestAttributeNameEquals(t *testing.T) {
	assert.True(t, conversion.AttributeNameEquals(path.Root("replication_specs").AtListIndex(0), "replication_specs"))
	assert.True(t, conversion.AttributeNameEquals(path.Root("replication_specs").AtMapKey("myKey"), "replication_specs"))
	assert.True(t, conversion.AttributeNameEquals(path.Root("replication_specs"), "replication_specs"))
	assert.True(t, conversion.AttributeNameEquals(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1), "region_configs"))
	assert.False(t, conversion.AttributeNameEquals(path.Root("replication_specs").AtListIndex(0), "region_configs"))
	assert.False(t, conversion.AttributeNameEquals(path.Root("replication_specs").AtMapKey("myKey"), "region_configs"))
	assert.False(t, conversion.AttributeNameEquals(path.Root("replication_specs"), "region_configs"))
}

func TestStripSquareBrackets(t *testing.T) {
	assert.Equal(t, "replication_specs", conversion.StripSquareBrackets(path.Root("replication_specs").AtListIndex(0)))
	assert.Equal(t, "replication_specs", conversion.StripSquareBrackets(path.Root("replication_specs").AtMapKey("myKey")))
	assert.Equal(t, "replication_specs", conversion.StripSquareBrackets(path.Root("replication_specs")))
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
