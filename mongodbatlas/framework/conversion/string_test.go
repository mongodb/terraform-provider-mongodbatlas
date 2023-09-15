package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/mongodbatlas/framework/conversion"
)

func TestIsStringPresent(t *testing.T) {
	var (
		empty     = ""
		one_blank = " "
		str       = "text"
	)
	tests := []struct {
		strPtr   *string
		expected bool
	}{
		{nil, false},
		{&empty, false},
		{&one_blank, true},
		{&str, true},
	}
	for _, test := range tests {
		if resp := conversion.IsStringPresent(test.strPtr); resp != test.expected {
			t.Errorf("IsStringPresent(%v) = %v; want %v", test.strPtr, resp, test.expected)
		}
	}
}
