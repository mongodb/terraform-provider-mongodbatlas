package advancedclustertpf_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
	"github.com/stretchr/testify/assert"
)

func HasPrefix(p path.Path, prefix path.Path) bool {
	prefixString := prefix.String()
	pString := p.String()
	return strings.HasPrefix(pString, prefixString)
}

func LastPart(p path.Path) string {
	parts := strings.Split(p.String(), ".")
	return parts[len(parts)-1]
}

func IsListIndex(p path.Path) bool {
	lastPart := LastPart(p)
	if IsMapIndex(p) {
		return false
	}
	return strings.HasSuffix(lastPart, "]")
}

func IsMapIndex(p path.Path) bool {
	lastPart := LastPart(p)
	return strings.HasSuffix(lastPart, "\"]")
}

func AsAddedIndex(p path.Path) string {
	parentString := p.ParentPath().ParentPath().String()
	lastPart := LastPart(p)
	indexWithSign := strings.Replace(lastPart, "[", "[+", 1)
	if parentString == "" {
		return indexWithSign
	}
	return parentString + "." + indexWithSign
}

func AsRemovedIndex(p path.Path) string {
	parentString := p.ParentPath().ParentPath().String()
	lastPart := LastPart(p)
	indexWithSign := strings.Replace(lastPart, "[", "[-", 1)
	if parentString == "" {
		return indexWithSign
	}
	return parentString + "." + indexWithSign
}

func StripSquareBrackets(p path.Path) string {
	if IsListIndex(p) {
		return p.ParentPath().String()
	}
	if IsMapIndex(p) {
		return p.ParentPath().String()
	}
	return p.String()
}

func AttributeNameEquals(p path.Path, name string) bool {
	noBrackets := StripSquareBrackets(p)
	return noBrackets == name || strings.HasSuffix(noBrackets, fmt.Sprintf(".%s", name))
}

func TestIsAttributeValueOnly(t *testing.T) {
	assert.True(t, advancedclustertpf.IsAttributeValueOnly(path.Root("replication_specs").AtListIndex(0)))
	assert.True(t, advancedclustertpf.IsAttributeValueOnly(path.Root("replication_specs").AtMapKey("myKey")))
	assert.True(t, advancedclustertpf.IsAttributeValueOnly(path.Root("replication_specs").AtSetValue(types.StringValue("myKey"))))
}


func TestAttributeNameEquals(t *testing.T) {
	assert.True(t, AttributeNameEquals(path.Root("replication_specs").AtListIndex(0), "replication_specs"))
	assert.True(t, AttributeNameEquals(path.Root("replication_specs").AtMapKey("myKey"), "replication_specs"))
	assert.True(t, AttributeNameEquals(path.Root("replication_specs"), "replication_specs"))
	assert.True(t, AttributeNameEquals(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1), "region_configs"))
	assert.False(t, AttributeNameEquals(path.Root("replication_specs").AtListIndex(0), "region_configs"))
	assert.False(t, AttributeNameEquals(path.Root("replication_specs").AtMapKey("myKey"), "region_configs"))
	assert.False(t, AttributeNameEquals(path.Root("replication_specs"), "region_configs"))
}

func TestStripSquareBrackets(t *testing.T) {
	assert.Equal(t, "replication_specs", StripSquareBrackets(path.Root("replication_specs").AtListIndex(0)))
	assert.Equal(t, "replication_specs", StripSquareBrackets(path.Root("replication_specs").AtMapKey("myKey")))
	assert.Equal(t, "replication_specs", StripSquareBrackets(path.Root("replication_specs")))
}

func TestIndexMethods(t *testing.T) {
	assert.True(t, IsListIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.False(t, IsListIndex(path.Root("replication_specs").AtName("region_configs")))
	assert.False(t, IsListIndex(path.Root("replication_specs").AtMapKey("region_configs")))
	assert.Equal(t, "replication_specs[+0]", AsAddedIndex(path.Root("replication_specs").AtListIndex(0)))
	assert.Equal(t, "replication_specs[0].region_configs[+1]", AsAddedIndex(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1)))
	assert.Equal(t, "replication_specs[-1]", AsRemovedIndex(path.Root("replication_specs").AtListIndex(1)))
	assert.Equal(t, "replication_specs[0].region_configs[-1]", AsRemovedIndex(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1)))
}

func TestPathMatches(t *testing.T) {
	prefix := path.Root("replication_specs").AtListIndex(0)
	assert.True(t, HasPrefix(path.Root("replication_specs").AtListIndex(0), prefix))
	assert.True(t, HasPrefix(path.Root("replication_specs").AtListIndex(0).AtName("region_configs"), prefix))
	assert.False(t, HasPrefix(path.Root("replication_specs").AtListIndex(1), prefix))
	assert.True(t, HasPrefix(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1), path.Empty()))
}
