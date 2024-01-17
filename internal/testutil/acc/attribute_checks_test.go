package acc_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
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
	objMap := map[string]interface{}{
		"str":    "my_string",
		"number": float64(1234),
		"bool1":  true,
		"bool2":  false,
		"nilvar": nil,
	}
	strMap := `
	{
		"str": "my_string",
		"number": 1234,
		"bool1": true,
		"bool2": false,
		"nilvar": null
	}
`
	if err := acc.JSONEquals(objMap)(strMap); err != nil {
		t.Errorf("JSONEquals() error = %v", err)
	}
}

func TestJSONStringEquals(t *testing.T) {
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
			checkFunc := acc.JSONStringEquals(tc.value)
			err := checkFunc(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("JSONStringEquals() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
