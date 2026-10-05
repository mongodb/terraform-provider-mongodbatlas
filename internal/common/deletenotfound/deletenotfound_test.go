package deletenotfound_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/deletenotfound"
)

func response(status int) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"detail":"not found","errorCode":"USER_NOT_FOUND"}`)),
		Request:    &http.Request{},
	}
}

func TestIsNotFound(t *testing.T) {
	testCases := map[string]struct {
		resp     *http.Response
		err      error
		expected bool
	}{
		"404 with error is tolerated":        {response(http.StatusNotFound), errors.New("404 USER_NOT_FOUND"), true},
		"204 without error is not tolerated": {response(http.StatusNoContent), nil, false},
		"400 is not tolerated":               {response(http.StatusBadRequest), errors.New("400 BAD_REQUEST"), false},
		"500 is not tolerated":               {response(http.StatusInternalServerError), errors.New("500 INTERNAL"), false},
		"nil response with error":            {nil, errors.New("transport failure"), false},
		"nil response without error":         {nil, nil, false},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, deletenotfound.IsNotFound(context.Background(), tc.resp, tc.err))
		})
	}
}
