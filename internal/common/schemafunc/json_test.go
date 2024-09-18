package schemafunc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/schemafunc"
)

func Test_EqualJSON(t *testing.T) {
	testCases := map[string]struct {
		old      string
		new      string
		expected bool
	}{
		"empty strings":                      {"", "", true},
		"different objects":                  {`{"a": 1}`, `{"b": 2}`, false},
		"invalid object":                     {`{{"a": 1}`, `{"b": 2}`, false},
		"double invalid object":              {`{{"a": 1}`, `{"b": 2}}`, false},
		"equal objects with different order": {`{"a": 1, "b": 2}`, `{"b": 2, "a": 1}`, true},
		"equal objects whitespace":           {`{"a": 1, "b": 2}`, `{"a":1,"b":2}`, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := schemafunc.EqualJSON(tc.old, tc.new, "vector search index")
			if actual != tc.expected {
				t.Errorf("Expected: %v, got: %v", tc.expected, actual)
			}
		})
	}
}
