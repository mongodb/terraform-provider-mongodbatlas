package mongodbatlas

import (
	"context"
	"fmt"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestMongoDBAtlasDatabaseUserResource_Schema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the resource.Resource and call its Schema method
	NewDatabaseUserRS().Schema(ctx, schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	// Validate the schema
	diagnostics := schemaResponse.Schema.ValidateImplementation(ctx)

	if diagnostics.HasError() {
		t.Fatalf("Schema validation diagnostics: %+v", diagnostics)
	}
}

func TestMongoDBAtlasDatabaseUserResourceSchema_Read(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resource, schemaResponse, err := newDatabaseUserRSMock(ctx)
	if err != nil {
		t.Fatalf("error during the resource creation: %v", err)
	}

	readRequest, readResponse := newReadRequestResponse(schemaResponse)
	requestModel := newRequestModel(ctx)
	readRequest.State.Set(ctx, requestModel)

	resource.Read(ctx, *readRequest, readResponse)

	var responseModel *tfDatabaseUserModel
	diag := readResponse.State.Get(ctx, &responseModel)
	if diag.HasError() {
		t.Fatalf("The READ operation failed: %v", diag)
	}

	validateResponse(t, requestModel, responseModel)
}

func TestMongoDBAtlasDatabaseUserResourceSchema_Create(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resource, schemaResponse, err := newDatabaseUserRSMock(ctx)
	if err != nil {
		t.Fatalf("error during the resource creation: %v", err)
	}

	createRequest, createResponse := newCreateRequestResponse(schemaResponse)
	requestModel := newRequestModel(ctx)
	createRequest.Plan.Set(ctx, requestModel)

	resource.Create(ctx, *createRequest, createResponse)

	var responseModel *tfDatabaseUserModel
	diag := createResponse.State.Get(ctx, &responseModel)
	if diag.HasError() {
		t.Fatalf("The CREATE operation failed: %v", diag)
	}

	validateResponse(t, requestModel, responseModel)
}

func TestMongoDBAtlasDatabaseUserResourceSchema_Update(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resource, schemaResponse, err := newDatabaseUserRSMock(ctx)
	if err != nil {
		t.Fatalf("error during the resource creation: %v", err)
	}

	updateRequest, updateResponse := newUpdateRequestResponse(schemaResponse)
	updateModel := newRequestModel(ctx)
	updateRequest.Plan.Set(ctx, updateModel)

	resource.Update(ctx, *updateRequest, updateResponse)

	var responseModel *tfDatabaseUserModel
	diag := updateResponse.State.Get(ctx, &responseModel)
	if diag.HasError() {
		t.Fatalf("The UPDATE operation failed: %v", diag)
	}

	validateResponse(t, updateModel, responseModel)
}

func TestMongoDBAtlasDatabaseUserResourceSchema_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resource, schemaResponse, err := newDatabaseUserRSMock(ctx)
	if err != nil {
		t.Fatalf("error during the resource creation: %v", err)
	}

	deleteRequest, deleteResponse := newDeleteRequestResponse(schemaResponse)
	deleteModel := newRequestModel(ctx)
	deleteRequest.State.Set(ctx, deleteModel)

	resource.Delete(ctx, *deleteRequest, deleteResponse)

	var responseModel *tfDatabaseUserModel
	diag := deleteResponse.State.Get(ctx, &responseModel)
	if diag.HasError() {
		t.Fatalf("The DELETE operation failed: %v", diag)
	}
}

func validateResponse(t *testing.T, requestModel, responseModel *tfDatabaseUserModel) {
	expectedID := fmt.Sprintf("%s-%s-%s", requestModel.ProjectID.ValueString(), requestModel.Username.ValueString(), requestModel.AuthDatabaseName.ValueString())
	if responseModel.ID.ValueString() != expectedID {
		t.Fatalf("expected %s, got %s", expectedID, responseModel.ID)
	}

	if requestModel.Password != responseModel.Password {
		t.Fatalf("expected %s, got %s", requestModel.Password, responseModel.Password)
	}

	if requestModel.AuthDatabaseName != responseModel.AuthDatabaseName {
		t.Fatalf("expected %s, got %s", requestModel.AuthDatabaseName, responseModel.AuthDatabaseName)
	}

	if requestModel.ProjectID != responseModel.ProjectID {
		t.Fatalf("expected %s, got %s", requestModel.AuthDatabaseName, responseModel.AuthDatabaseName)
	}
}

func newDatabaseUserRSMock(ctx context.Context) (fwresource.Resource, *fwresource.SchemaResponse, error) {
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	// Instantiate the resource.Resource and call its Schema method
	resource := NewDatabaseUserRSWithClient(&MongoDBClient{
		Atlas: &matlas.Client{
			DatabaseUsers: &MockDatabaseUsersServiceOp{},
		},
	})

	resource.Schema(ctx, schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		return nil, nil, fmt.Errorf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	return resource, schemaResponse, nil
}

func newReadRequestResponse(schemaResponse *fwresource.SchemaResponse) (*fwresource.ReadRequest, *fwresource.ReadResponse) {
	return &fwresource.ReadRequest{
			State: tfsdk.State{
				Schema: schemaResponse.Schema,
			},
		}, &fwresource.ReadResponse{
			State: tfsdk.State{
				Schema: schemaResponse.Schema,
			},
		}
}

func newCreateRequestResponse(schemaResponse *fwresource.SchemaResponse) (*fwresource.CreateRequest, *fwresource.CreateResponse) {
	return &fwresource.CreateRequest{
			Plan: tfsdk.Plan{
				Schema: schemaResponse.Schema,
			},
		}, &fwresource.CreateResponse{
			State: tfsdk.State{
				Schema: schemaResponse.Schema,
			},
		}
}

func newUpdateRequestResponse(schemaResponse *fwresource.SchemaResponse) (*fwresource.UpdateRequest, *fwresource.UpdateResponse) {
	return &fwresource.UpdateRequest{
			Plan: tfsdk.Plan{
				Schema: schemaResponse.Schema,
			},
		}, &fwresource.UpdateResponse{
			State: tfsdk.State{
				Schema: schemaResponse.Schema,
			},
		}
}

func newDeleteRequestResponse(schemaResponse *fwresource.SchemaResponse) (*fwresource.DeleteRequest, *fwresource.DeleteResponse) {
	return &fwresource.DeleteRequest{
			State: tfsdk.State{
				Schema: schemaResponse.Schema,
			},
		}, &fwresource.DeleteResponse{
			State: tfsdk.State{
				Schema: schemaResponse.Schema,
			},
		}
}

func newRequestModel(ctx context.Context) *tfDatabaseUserModel {
	rolesSet, _ := types.SetValueFrom(ctx, RoleObjectType, []tfRoleModel{
		{
			RoleName:       types.StringValue("roleName"),
			CollectionName: types.StringValue("CollectionName"),
			DatabaseName:   types.StringValue("DatabaseName"),
		},
	})

	labelSet, _ := types.SetValueFrom(ctx, LabelObjectType, []tfLabelModel{
		{
			Key:   types.StringValue("key"),
			Value: types.StringValue("value"),
		},
	})
	scopeSet, _ := types.SetValueFrom(ctx, ScopeObjectType, []tfScopeModel{
		{
			Name: types.StringValue("Name"),
			Type: types.StringValue("Type"),
		},
	})

	return &tfDatabaseUserModel{
		Username:         types.StringValue("test"),
		Password:         types.StringValue("test"),
		AuthDatabaseName: types.StringValue("test"),
		ProjectID:        types.StringValue("test"),
		Roles:            rolesSet,
		Labels:           labelSet,
		Scopes:           scopeSet,
	}
}

type MockDatabaseUsersServiceOp struct{}

func (s *MockDatabaseUsersServiceOp) Get(ctx context.Context, databaseName, groupID, username string) (*matlas.DatabaseUser, *matlas.Response, error) {
	return &matlas.DatabaseUser{
		GroupID:      groupID,
		Username:     username,
		DatabaseName: databaseName,
	}, nil, nil
}

func (s *MockDatabaseUsersServiceOp) List(ctx context.Context, groupID string, listOptions *matlas.ListOptions) ([]matlas.DatabaseUser, *matlas.Response, error) {
	return nil, nil, nil
}

func (s *MockDatabaseUsersServiceOp) Create(ctx context.Context, groupID string, createRequest *matlas.DatabaseUser) (*matlas.DatabaseUser, *matlas.Response, error) {
	return &matlas.DatabaseUser{
		GroupID:      groupID,
		Username:     createRequest.Username,
		DatabaseName: createRequest.DatabaseName,
	}, nil, nil
}

func (s *MockDatabaseUsersServiceOp) Update(ctx context.Context, groupID, username string, updateRequest *matlas.DatabaseUser) (*matlas.DatabaseUser, *matlas.Response, error) {
	return &matlas.DatabaseUser{
		GroupID:      groupID,
		Username:     updateRequest.Username,
		DatabaseName: updateRequest.DatabaseName,
	}, nil, nil
}

func (s *MockDatabaseUsersServiceOp) Delete(ctx context.Context, databaseName, groupID, username string) (*matlas.Response, error) {
	return nil, nil
}
