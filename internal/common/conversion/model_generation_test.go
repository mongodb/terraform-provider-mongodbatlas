package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyModel(t *testing.T) {
	type destType struct {
		AttrStr        string
		attrUnexported string
		AttrInt        int
	}

	testCases := map[string]struct {
		input            any
		expected         any
		expectedErrorStr string
	}{
		"basic": {
			input: &struct {
				AttrStr string
				AttrInt int
			}{
				AttrStr: "val",
				AttrInt: 1,
			},
			expected: &destType{
				AttrStr:        "val",
				AttrInt:        1,
				attrUnexported: "",
			},
		},
		"missing field": {
			input: &struct {
				AttrStr string
			}{
				AttrStr: "val",
			},
			expected: &destType{
				AttrStr: "val",
			},
		},
		"extra field": {
			input: &struct {
				AttrStr   string
				AttrExtra string
				AttrInt   int
			}{
				AttrStr:   "val",
				AttrExtra: "extra",
				AttrInt:   1,
			},
			expected: &destType{
				AttrStr: "val",
				AttrInt: 1,
			},
		},
		"different type": {
			input: &struct {
				AttrStr bool
			}{
				AttrStr: true,
			},
			expectedErrorStr: "field has different type: AttrStr",
		},
		"unexported": {
			input: &struct {
				attrUnexported string
			}{
				attrUnexported: "val",
			},
			expectedErrorStr: "field can't be set, probably unexported: attrUnexported",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			dest, err := conversion.CopyModel[destType](tc.input)
			if err == nil {
				assert.Equal(t, tc.expected, dest)
				assert.Equal(t, "", tc.expectedErrorStr)
			} else {
				require.ErrorContains(t, err, tc.expectedErrorStr)
				assert.Nil(t, dest)
				assert.Nil(t, tc.expected)
			}
		})
	}
}
