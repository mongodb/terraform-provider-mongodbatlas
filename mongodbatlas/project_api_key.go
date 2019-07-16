package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const projectAPIKeysPath = "groups/%s/apiKeys/%s"

//ProjectAPIKeysService is an interface for interfacing with the APIKeys
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/apiKeys/#organization-api-keys-on-projects-endpoints
type ProjectAPIKeysService interface {
	Assign(context.Context, string, string) (*Response, error)
	Unassign(context.Context, string, string) (*Response, error)
}

//ProjectAPIKeysOp handles communication with the APIKey related methods
// of the MongoDB Atlas API
type ProjectAPIKeysOp struct {
	client *Client
}

var _ ProjectAPIKeysService = &ProjectAPIKeysOp{}

//Assign an API-KEY related to {ORG-ID} to a the project with {PROJECT-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-all/
func (s *ProjectAPIKeysOp) Assign(ctx context.Context, orgID string, projectID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("apiKeyID", "must be set")
	}

	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}

	basePath := fmt.Sprintf(projectAPIKeysPath, orgID, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

//Unassign an API-KEY related to {ORG-ID} to a the project with {PROJECT-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/apiKeys-orgs-get-all/
func (s *ProjectAPIKeysOp) Unassign(ctx context.Context, orgID string, projectID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("apiKeyID", "must be set")
	}

	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}

	basePath := fmt.Sprintf(projectAPIKeysPath, orgID, projectID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
