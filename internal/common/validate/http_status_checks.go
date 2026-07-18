package validate

import "net/http"

func StatusNotFound(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

func StatusUnauthorized(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusUnauthorized
}

func StatusServiceUnavailable(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusServiceUnavailable
}

func StatusBadRequest(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusBadRequest
}

func StatusInternalServerError(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusInternalServerError
}
