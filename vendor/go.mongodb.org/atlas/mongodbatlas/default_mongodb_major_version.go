package mongodbatlas

import (
	"bytes"
	"context"
	"net/http"
)

const defaultMongoDBMajorVersionPath = "api/private/unauth/nds/defaultMongoDBMajorVersion"

// DefaultMongoDBMajorVersionService this service is to be used by other MongoDB tools
// to determine the current default major version of MongoDB Server in Atlas.
//
// We currently make no promise to support or document this service or endpoint
// beyond what can be seen here.
type DefaultMongoDBMajorVersionService interface {
	Get(context.Context) (string, *Response, error)
}

// DefaultMongoDBMajorVersionServiceOp is an implementation of DefaultMongoDBMajorVersionService
type DefaultMongoDBMajorVersionServiceOp struct {
	Client PlainRequestDoer
}

var _ DefaultMongoDBMajorVersionService = &DefaultMongoDBMajorVersionServiceOp{}

// Get gets the current major MongoDB version in Atlas
func (s *DefaultMongoDBMajorVersionServiceOp) Get(ctx context.Context) (string, *Response, error) {
	req, err := s.Client.NewPlainRequest(ctx, http.MethodGet, defaultMongoDBMajorVersionPath)
	if err != nil {
		return "", nil, err
	}
	root := new(bytes.Buffer)

	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return "", resp, err
	}

	return root.String(), resp, err
}
