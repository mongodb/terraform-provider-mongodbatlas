package acc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
)

func TestPathFromString(t *testing.T) {
	assert.Equal(t, "some_key", tfjsonpath.New("some_key").String())
	assert.Equal(t, "some_key.nested_key", tfjsonpath.New("some_key").AtMapKey("nested_key").String())
	assert.Equal(t, "some_key.0", tfjsonpath.New("some_key").AtSliceIndex(0).String())
	for name, test := range map[string]struct {
		pathStr string
		want    tfjsonpath.Path
	}{
		"simple": {
			pathStr: "some_key",
			want:    tfjsonpath.New("some_key"),
		},
		"nested": {
			pathStr: "some_key.some_key2",
			want:    tfjsonpath.New("some_key").AtMapKey("some_key2"),
		},
		"nested_with_index": {
			pathStr: "some_key.some_key2.0",
			want:    tfjsonpath.New("some_key").AtMapKey("some_key2").AtSliceIndex(0),
		},
		"nested_with_index_and_key": {
			pathStr: "some_key.some_key2.0.nested_key",
			want:    tfjsonpath.New("some_key").AtMapKey("some_key2").AtSliceIndex(0).AtMapKey("nested_key"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := acc.PathFromString(test.pathStr)
			assert.Equal(t, test.want.String(), got.String())
		})
	}
}
