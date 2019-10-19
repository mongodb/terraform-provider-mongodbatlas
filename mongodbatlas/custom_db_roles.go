package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const dbCustomDBRolesBasePath = "groups/%s/customDBRoles/roles"

// CustomDBRolesService is an interface for working wit the Custom MongoDB Roles
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/custom-roles/
type CustomDBRolesService interface {
	List(context.Context, string, *ListOptions) (*[]CustomDbRole, *Response, error)
	Get(context.Context, string, string) (*CustomDbRole, *Response, error)
	Create(context.Context, string, *CustomDbRole) (*CustomDbRole, *Response, error)
	Update(context.Context, string, string, *CustomDbRole) (*CustomDbRole, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

//CustomDBRolesServiceOp handles communication with the CustomDBRoles related methods of the
//MongoDB Atlas API
type CustomDBRolesServiceOp struct {
	client *Client
}

var _ CustomDBRolesService = &CustomDBRolesServiceOp{}

// A Resource describes a specific resource the Role will allow operating on.
type Resource struct {
	Collection string `json:"collection,omitempty"`
	Db         string `json:"db,omitempty"`
	Cluster    bool   `json:"cluster,omitempty"`
}

// An Action describes the operation the role will include, for a specific set of Resources.
type Action struct {
	Action    string     `json:"action,omitempty"`
	Resources []Resource `json:"resources,omitempty"`
}

// An InheritedRole describes the role that this Role inherits from.
type InheritedRole struct {
	Db   string `json:"db,omitempty"`
	Role string `json:"role,omitempty"`
}

// CustomDBRole represents a Custom MongoDB Role in your cluster.
type CustomDbRole struct {
	Actions        []Action        `json:"actions,omitempty"`
	InheritedRoles []InheritedRole `json:"inheritedRoles,omitempty"`
	RoleName       string          `json:"roleName,omitempty"`
}

//List gets all custom db roles in the project.
//See more: https://docs.atlas.mongodb.com/reference/api/custom-roles-get-all-roles/
func (s *CustomDBRolesServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) (*[]CustomDbRole, *Response, error) {
	path := fmt.Sprintf(dbCustomDBRolesBasePath, groupID)

	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]CustomDbRole)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// Gets a single Custom MongoDB Role in the project.
// See more: https://docs.atlas.mongodb.com/reference/api/custom-roles-get-single-role/
func (s *CustomDBRolesServiceOp) Get(ctx context.Context, groupID string, roleName string) (*CustomDbRole, *Response, error) {
	if roleName == "" {
		return nil, nil, NewArgError("roleName", "must be set")
	}

	basePath := fmt.Sprintf(dbCustomDBRolesBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, roleName)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(CustomDbRole)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Creates a new Custom MongoDB Role in the project.
// See more: https://docs.atlas.mongodb.com/reference/api/custom-roles-create-a-role/
func (s *CustomDBRolesServiceOp) Create(ctx context.Context, groupID string, createRequest *CustomDbRole) (*CustomDbRole, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(dbCustomDBRolesBasePath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(CustomDbRole)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Updates a single Custom MongoDB Role.
// See more: https://docs.atlas.mongodb.com/reference/api/custom-roles-update-a-role/
func (s *CustomDBRolesServiceOp) Update(ctx context.Context, groupID string, roleName string, updateRequest *CustomDbRole) (*CustomDbRole, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(dbCustomDBRolesBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, roleName)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(CustomDbRole)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Deletes a single Custom MongoDB Role.
// See more: https://docs.atlas.mongodb.com/reference/api/custom-roles-delete-a-role/
func (s *CustomDBRolesServiceOp) Delete(ctx context.Context, groupID string, roleName string) (*Response, error) {
	if roleName == "" {
		return nil, NewArgError("roleName", "must be set")
	}

	basePath := fmt.Sprintf(dbCustomDBRolesBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, roleName)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
