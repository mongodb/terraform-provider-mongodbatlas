package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const dbCustomDBRolesBasePath = "groups/%s/customDBRoles/roles"

type CustomDBRolesService interface {
	List(context.Context, string, *ListOptions) (*[]CustomDbRole, *Response, error)
}

//CustomDBRolesServiceOp handles communication with the CustomDBRoles related methods of the
//MongoDB Atlas API
type CustomDBRolesServiceOp struct {
	client *Client
}

var _ CustomDBRolesService = &CustomDBRolesServiceOp{}

type Resource struct {
	Collection string `json:"collection,omitempty"`
	Db         string `json:"db,omitempty"`
	Cluster    bool   `json:"cluster,omitempty"`
}

type Action struct {
	Action    string     `json:"action,omitempty"`
	Resources []Resource `json:"resources,omitempty"`
}

type InheritedRole struct {
	Db   string `json:"db,omitempty"`
	Role string `json:"role,omitempty"`
}

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
