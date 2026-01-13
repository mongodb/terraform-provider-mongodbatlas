package serviceaccountaccesslistentry

import (
	"context"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"
)

const ItemsPerPage = 500 // Max items per page

type ListPageFunc func(ctx context.Context, pageNum int) (*admin.PaginatedServiceAccountIPAccessEntry, *http.Response, error)

// ReadAccessListEntry Iterates through access list pages looking for the entry.
// The first page can be provided to skip an API call. Useful for Create operation, which returns the first page.
func ReadAccessListEntry(
	ctx context.Context,
	firstPage *admin.PaginatedServiceAccountIPAccessEntry,
	listPageFunc ListPageFunc,
	cidrOrIP string,
) (*admin.ServiceAccountIPAccessListEntry, *http.Response, error) {
	var err error
	var apiResp *http.Response

	count := 0
	page := firstPage
	for currentPage := 1; ; currentPage++ {
		if page == nil {
			page, apiResp, err = listPageFunc(ctx, currentPage)
			if err != nil {
				return nil, apiResp, err
			}
		}

		results := page.GetResults()
		count += len(results)

		for i := range results {
			entry := &results[i]
			if entry.GetIpAddress() == cidrOrIP || entry.GetCidrBlock() == cidrOrIP {
				return entry, nil, nil
			}
		}

		if len(results) == 0 || count >= page.GetTotalCount() {
			break
		}
		page = nil
	}

	return nil, nil, nil
}
