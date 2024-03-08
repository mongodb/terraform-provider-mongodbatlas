package validate_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/validate"
)

func TestStringIsUppercase(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError bool
	}{
		{
			name:          "AWS",
			expectedError: false,
		},
		{
			name:          "aws",
			expectedError: true,
		},
		{
			name:          "",
			expectedError: false,
		},
		{
			name:          "AwS",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diag := validate.StringIsUppercase()(tc.name, nil)
			if diag.HasError() != tc.expectedError {
				t.Errorf("Case %s: Received unexpected error: %v", tc.name, diag[0].Summary)
			}
		})
	}
}
