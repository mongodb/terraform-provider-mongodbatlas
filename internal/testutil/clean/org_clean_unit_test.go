package clean_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/clean"
)

func TestSkipUnauthorizedErr(t *testing.T) {
	apiErr := errors.New("boom")
	testCases := map[string]struct {
		resp     *http.Response
		err      error
		expected error
	}{
		"nil error returns nil": {
			resp:     &http.Response{StatusCode: http.StatusUnauthorized},
			err:      nil,
			expected: nil,
		},
		"401 with error returns ErrUnauthorized": {
			resp:     &http.Response{StatusCode: http.StatusUnauthorized},
			err:      apiErr,
			expected: clean.ErrUnauthorized,
		},
		"other status returns original error": {
			resp:     &http.Response{StatusCode: http.StatusConflict},
			err:      apiErr,
			expected: apiErr,
		},
		"nil response returns original error": {
			resp:     nil,
			err:      apiErr,
			expected: apiErr,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, clean.SkipUnauthorizedErr(tc.resp, tc.err))
		})
	}
}
