package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
)

func TestHasElementsSliceOrMap(t *testing.T) {
	testCasesTrue := map[string]any{
		"slice":       []string{"hi"},
		"map":         map[string]string{"hi": "there"},
		"int int map": map[int]int{1: 2},
		"double map": map[string]map[string]string{
			"hi": {"there": "bye"},
		},
	}
	testCasesFalse := map[string]any{
		"nil":                           nil,
		"empty slice":                   []string{},
		"empty map":                     map[string]string{},
		"empty int int map":             map[int]int{},
		"not a collection but with len": "hello",
		"random object":                 123,
	}
	for name, value := range testCasesTrue {
		t.Run(name, func(t *testing.T) {
			assert.True(t, conversion.HasElementsSliceOrMap(value))
		})
	}
	for name, value := range testCasesFalse {
		t.Run(name, func(t *testing.T) {
			assert.False(t, conversion.HasElementsSliceOrMap(value))
		})
	}
}
