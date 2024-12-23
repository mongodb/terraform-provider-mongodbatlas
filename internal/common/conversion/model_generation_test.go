package conversion_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/stretchr/testify/assert"
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
		expectedPanicStr string
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
			expectedPanicStr: "field has different type: AttrStr",
		},
		"unexported": {
			input: &struct {
				attrUnexported string
			}{
				attrUnexported: "val",
			},
			expectedPanicStr: "field can't be set, probably unexported: attrUnexported",
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.expectedPanicStr == "" {
				assert.Equal(t, tc.expected, conversion.CopyModel[destType](tc.input))
			} else {
				assert.Nil(t, tc.expected)
				assert.PanicsWithValue(t, tc.expectedPanicStr, func() {
					conversion.CopyModel[destType](tc.input)
				})
			}
		})
	}
}
