package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// OrgOwner is an Organization
	OrgOwner = "ORG_OWNER"
	// OrgMember is an Organization
	OrgMember = "ORG_MEMBER"
	// OrgProjectCreator is an Organization
	OrgProjectCreator = "ORG_PROJECT_CREATOR"
	// OrgBillingAdmin is an Organization
	OrgBillingAdmin = "ORG_BILLING_ADMIN"
	// OrgReadOnly is an Organization
	OrgReadOnly = "ORG_READ_ONLY"

	organizationBasePath = "orgs"
)

// OrganizationsService is an interface for interfacing with the Organizations
// endpoints of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/organizations/
type OrganizationsService interface {
	GetAllOrganizations(context.Context) (*Organizations, *Response, error)
	GetOneOrganization(context.Context, string) (*Organization, *Response, error)
	GetAllOrganizationUsers(context.Context, string) (*OrganizationUsers, *Response, error)
	Create(context.Context, string) (*Organization, *Response, error)
	UpdateOrganizationName(context.Context, *Organization) (*Organization, *Response, error)
	Delete(context.Context, string) (*Response, error)
	GetAllOrganizationProjects(context.Context, string) (*Projects, *Response, error)
}

//OrganizationsServiceOp handles communication with the DatabaseUsers related methos of the
//MongoDB Atlas API
type OrganizationsServiceOp struct {
	client *Client
}

var _ OrganizationsService = &OrganizationsServiceOp{}

// Organization is the structure of Organization request.
type Organization struct {
	ID    string  `json:"id,omitempty"`    //The unique identifier for the organization.
	Name  string  `json:"name,omitempty"`  //The name of the organization you want to create.
	Links []*Link `json:"links,omitempty"` //One or more links to sub-resources and/or related resources
}

// Organizations represents an array of organizartions
type Organizations struct {
	Links      []*Link         `json:"links"`      //One or more links to sub-resources and/or related resources
	Results    []*Organization `json:"results"`    //Results is an array of organizations
	TotalCount int             `json:"totalCount"` //It is total about organization array
}

// OrganizationUser represent a organization user
type OrganizationUser struct {
	ID           string              `json:"id,omitempty"`           //The user’s id.
	TeamIDS      []string            `json:"teamIds,omitempty"`      //An array of the team ids for the organization.
	Username     string              `json:"username,omitempty"`     //The username for authenticating to MongoDB.
	Country      string              `json:"country,omitempty"`      //The country where the user lives.
	EmailAddress string              `json:"emailAddress,omitempty"` //The user’s email address.
	FirstName    string              `json:"firstName,omitempty"`    //The user’s first name.
	LastName     string              `json:"lastName,omitempty"`     //ID of the Atlas project the user belongs to.
	MobileNumber string              `json:"mobileNumber,omitempty"` //The user’s mobile phone number.
	Roles        []*OrganizationRole `json:"roles,omitempty"`        //An array of the user’s roles within the Organization and for each Project to which the user belongs.
	Links        []*Link             `json:"links,omitempty"`        //One or more links to sub-resources and/or related resources.
}

// OrganizationRole represents the roles that exists in a organization for the organization user
type OrganizationRole struct {
	GroupID  string `json:"groupId,omitempty"`  //The {groupId} represents the Organization or Project to which this role applies. Possible values are: orgId or groupId.
	OrgID    string `json:"orgId,omitempty"`    //The orgId} represents the Organization or Project to which this role applies. Possible values are: orgId or groupId.
	RoleName string `json:"roleName,omitempty"` //The name of the role. The users resource returns all the roles the user has in either Atlas or
}

// OrganizationUsers represents an array og all user in a organization
type OrganizationUsers struct {
	Links      []*Link             `json:"links,omitempty"` //One or more links to sub-resources and/or related resources.
	Results    []*OrganizationUser `json:"results"`         //Results is an array of organizations user
	TotalCount int                 `json:"totalCount"`      //It is total about organization users array
}

// GetAllOrganizations gets all organizations.
// See more: https://docs.atlas.mongodb.com/reference/api/organization-get-all/
func (s *OrganizationsServiceOp) GetAllOrganizations(ctx context.Context) (*Organizations, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, organizationBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organizations)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// GetOneOrganization gets a single organization.
// See more: https://docs.atlas.mongodb.com/reference/api/organization-get-one/
func (s *OrganizationsServiceOp) GetOneOrganization(ctx context.Context, organizationID string) (*Organization, *Response, error) {
	if organizationID == "" {
		return nil, nil, NewArgError("organizationID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", organizationBasePath, organizationID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetAllOrganizationUsers gets all organization oser.
// See more: https://docs.atlas.mongodb.com/reference/api/organization-users-get-all-users/
func (s *OrganizationsServiceOp) GetAllOrganizationUsers(ctx context.Context, organizationID string) (*OrganizationUsers, *Response, error) {
	if organizationID == "" {
		return nil, nil, NewArgError("organizationID", "must be set")
	}

	path := fmt.Sprintf("%s/%s/users", organizationBasePath, organizationID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(OrganizationUsers)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create one organization.
// See more: https://docs.atlas.mongodb.com/reference/api/organization-create-one/
func (s *OrganizationsServiceOp) Create(ctx context.Context, organizationName string) (*Organization, *Response, error) {
	if organizationName == "" {
		return nil, nil, NewArgError("organizationName", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, organizationBasePath, &Organization{
		Name: organizationName,
	})
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// UpdateOrganizationName update only the organization name
//See more: https://docs.atlas.mongodb.com/reference/api/organization-rename/
func (s *OrganizationsServiceOp) UpdateOrganizationName(ctx context.Context, createRequest *Organization) (*Organization, *Response, error) {
	if createRequest.ID == "" {
		return nil, nil, NewArgError("organizationID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s", organizationBasePath, createRequest.ID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Delete an organization by organization id.
// See more: https://docs.atlas.mongodb.com/reference/api/organization-delete-one/
func (s *OrganizationsServiceOp) Delete(ctx context.Context, organizationID string) (*Response, error) {
	if organizationID == "" {
		return nil, NewArgError("organizationID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", organizationBasePath, organizationID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

//GetAllOrganizationProjects gets all projects in an organization.
//See more: https://docs.atlas.mongodb.com/reference/api/organization-get-all-projects/
func (s *OrganizationsServiceOp) GetAllOrganizationProjects(ctx context.Context, organizationID string) (*Projects, *Response, error) {
	if organizationID == "" {
		return nil, nil, NewArgError("organizationID", "must be set")
	}

	path := fmt.Sprintf("%s/%s/groups", organizationBasePath, organizationID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Projects)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}
