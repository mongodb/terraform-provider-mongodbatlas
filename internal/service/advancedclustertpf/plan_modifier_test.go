package advancedclustertpf_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/stretchr/testify/assert"
)

func HasPrefix(p path.Path, prefix path.Path) bool {
	prefixString := prefix.String()
	pString := p.String()
	return strings.HasPrefix(pString, prefixString)
}

func TestPathMatches(t *testing.T) {
	prefix := path.Root("replication_specs").AtListIndex(0)
	assert.True(t, HasPrefix(path.Root("replication_specs").AtListIndex(0), prefix))
	assert.True(t, HasPrefix(path.Root("replication_specs").AtListIndex(0).AtName("region_configs"), prefix))
	assert.False(t, HasPrefix(path.Root("replication_specs").AtListIndex(1), prefix))
	assert.True(t, HasPrefix(path.Root("replication_specs").AtListIndex(0).AtName("region_configs").AtListIndex(1), path.Empty()))
}
