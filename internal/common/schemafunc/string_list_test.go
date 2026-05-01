package schemafunc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

func Test_StringSlicesEqualIgnoringOrder(t *testing.T) {
	testCases := map[string]struct {
		a        []string
		b        []string
		expected bool
	}{
		"same order":        {[]string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		"different order":   {[]string{"a", "b", "c"}, []string{"c", "a", "b"}, true},
		"different values":  {[]string{"a", "b"}, []string{"a", "c"}, false},
		"different lengths": {[]string{"a", "b", "c"}, []string{"a", "b"}, false},
		"empty lists":       {[]string{}, []string{}, true},
		"single element":    {[]string{"a"}, []string{"a"}, true},
		"nil slices":        {nil, nil, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := schemafunc.StringSlicesEqualIgnoringOrder(tc.a, tc.b)
			if actual != tc.expected {
				t.Errorf("Expected: %v, got: %v", tc.expected, actual)
			}
		})
	}
}
