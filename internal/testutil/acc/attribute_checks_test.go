package acc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/stretchr/testify/assert"
)

func TestIntGreaterThan(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		errorMsg string
		value    int
		wantErr  bool
	}{
		{"ValidGreater", "10", "", 5, false},
		{"ValidEqual", "5", "5 is not greater than 5", 5, true},
		{"ValidLess", "3", "3 is not greater than 5", 5, true},
		{"InvalidInput", "abc", "strconv.Atoi: parsing \"abc\": invalid syntax", 5, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkFunc := acc.IntGreatThan(tc.value)
			err := checkFunc(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("IntGreatThan() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil && err.Error() != tc.errorMsg {
				t.Errorf("IntGreatThan() error message = %v, want %v", err, tc.errorMsg)
			}
		})
	}
}

func TestJSONEquals(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		value   string
		wantErr bool
	}{
		{"same", "{\"a\": 1}", "{\"a\": 1}", false},
		{"same with blanks", "{\"a\": 1, \"b\": 2}", "{\"a\": \t1,   \n\"b\": 2}", false},
		{"differenct objects", "{\"a\": 1}", "{\"a\": false}", true},
		{"different types", "{\"a\": 1}", "[{\"a\": 1}]", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			checkFunc := acc.JSONEquals(tc.value)
			err := checkFunc(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("JSONEquals() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestAddNoAttrSetChecks(t *testing.T) {
	asserter := assert.New(t)
	targetName := "unitTestTarget"
	asserter.Empty(acc.AddNoAttrSetChecks(targetName, nil), "an empty and no attributes should have length 0")
	checks1 := acc.AddNoAttrSetChecks(targetName, nil, "attr1")
	asserter.Len(checks1, 1, "empty list 1 extra should have 1 element")
	asserter.Len(acc.AddNoAttrSetChecks(targetName, checks1, "attr2"), 2, "existing list + 1 extra should have 2 elements")
	asserter.Len(checks1, 1, "existing list should not be modified")
	asserter.Len(acc.AddNoAttrSetChecks(targetName, checks1, "attr2", "attr3"), 3, "existing list + 2 extra should have 3 elements")
}
