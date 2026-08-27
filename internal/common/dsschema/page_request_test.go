package dsschema_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/dsschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testPageResponse struct {
	results    []string
	totalCount int
}

func (r *testPageResponse) GetResults() []string {
	return r.results
}

func (r *testPageResponse) GetTotalCount() int {
	return r.totalCount
}

func TestAllPagesEmptyFirstPageReturnsNonNilSlice(t *testing.T) {
	results, err := dsschema.AllPages(context.Background(), func(_ context.Context, pageNum int) (dsschema.PaginateResponse[string], *http.Response, error) {
		require.Equal(t, 1, pageNum)
		return &testPageResponse{
			results:    []string{},
			totalCount: 0,
		}, nil, nil
	})
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Empty(t, results)
}
