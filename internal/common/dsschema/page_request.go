package dsschema

import (
	"context"
	"errors"
	"net/http"
)

type PaginateResponse[T any] interface {
	GetResults() []T
	GetTotalCount() int
}

func AllPages[T any](ctx context.Context, listOnPage func(ctx context.Context, pageNum int) (PaginateResponse[T], *http.Response, error)) ([]T, error) {
	var results []T
	for currentPage := 1; ; currentPage++ {
		resp, _, err := listOnPage(ctx, currentPage)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, errors.New("no response")
		}
		currentResults := resp.GetResults()
		results = append(results, currentResults...)
		if len(currentResults) == 0 || len(results) >= resp.GetTotalCount() {
			break
		}
	}
	return results, nil
}
