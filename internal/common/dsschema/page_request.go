package dsschema

import (
	"context"
	"errors"
	"net/http"
)

type PaginateResponse[ElementsModel any] interface {
	GetResults() []ElementsModel
	GetTotalCount() int
}

func AllPages[ElementsModel any](ctx context.Context, call func(ctx context.Context, pageNum int) (PaginateResponse[ElementsModel], *http.Response, error)) ([]ElementsModel, error) {
	var results []ElementsModel
	for pageNum := 1; ; pageNum++ {
		resp, _, err := call(ctx, pageNum)
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

type PaginateRequest[ElementsModel any] interface {
	PageNum(int) PaginateRequest[ElementsModel]
	Execute() (*PaginateResponse[ElementsModel], *http.Response, error)
}

func AllPagesFromRequest[ElementsModel any, V PaginateRequest[ElementsModel]](ctx context.Context, req V) ([]ElementsModel, error) {
	return AllPages[ElementsModel](ctx, func(ctx context.Context, pageNum int) (PaginateResponse[ElementsModel], *http.Response, error) {
		request := req.PageNum(pageNum)
		a, b, c := request.Execute()
		return *a, b, c
	})
}
