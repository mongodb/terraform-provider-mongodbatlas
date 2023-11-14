package testutils

import "testing"

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
			checkFunc := IntGreatThan(tc.value)
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
