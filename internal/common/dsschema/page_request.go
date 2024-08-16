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

func AllPages[T any](ctx context.Context, call func(ctx context.Context, pageNum int) (PaginateResponse[T], *http.Response, error)) ([]T, error) {
	var results []T
	for i := 1; ; i++ {
		resp, _, err := call(ctx, i)
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

type PaginateRequest[T any] interface {
	Execute() (PaginateResponse[T], *http.Response, error)
	PageNum(int) PaginateRequest[T]
}

func AllPagesFromRequest[T any](ctx context.Context, req PaginateRequest[T]) ([]T, error) {
	return AllPages(ctx, func(ctx context.Context, pageNum int) (PaginateResponse[T], *http.Response, error) {
		request := req.PageNum(pageNum)
		return request.Execute()
	})
}
