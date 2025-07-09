package clouduserorgassignment_test

import (
	"context"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/clouduserorgassignment"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.OrgUserResponse
	expectedTFModel *clouduserorgassignment.TFModel
}

func TestCloudUserOrgAssignmentSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{ // TODO: consider adding test cases to contemplate all possible API responses
		"Complete SDK response": {
			SDKResp:         &admin.OrgUserResponse{},
			expectedTFModel: &clouduserorgassignment.TFModel{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := clouduserorgassignment.NewTFModel(context.Background(), tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfToSDKModelTestCase struct {
	tfModel        *clouduserorgassignment.TFModel
	expectedSDKReq *admin.OrgUserResponse
}

func TestCloudUserOrgAssignmentTFModelToSDK(t *testing.T) {
	testCases := map[string]tfToSDKModelTestCase{
		"Complete TF state": {
			tfModel:        &clouduserorgassignment.TFModel{},
			expectedSDKReq: &admin.OrgUserResponse{},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := clouduserorgassignment.NewOrgUserReq(context.Background(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}

func TestNewTFRoles(t *testing.T) {
	t.Run("nil roles returns null object", func(t *testing.T) {
		obj, diags := clouduserorgassignment.NewTFRoles(context.Background(), nil)
		assert.False(t, diags.HasError())
		assert.True(t, obj.IsNull())
	})

	t.Run("empty roles returns non-null object", func(t *testing.T) {
		roles := &admin.OrgUserRolesResponse{}
		obj, diags := clouduserorgassignment.NewTFRoles(context.Background(), roles)
		assert.False(t, diags.HasError())
		assert.False(t, obj.IsNull())
	})
}

func TestNewTFProjectRoleAssignments(t *testing.T) {
	t.Run("nil groupRoleAssignments returns empty list", func(t *testing.T) {
		lst := clouduserorgassignment.NewTFProjectRoleAssignments(context.Background(), nil)
		assert.True(t, lst.IsNull() || len(lst.Elements()) == 0)
	})

	t.Run("populated groupRoleAssignments returns non-empty list", func(t *testing.T) {
		groupRoles := []admin.GroupRoleAssignment{
			{
				GroupId:    admin.PtrString("project1"),
				GroupRoles: &[]string{"ROLE1", "ROLE2"},
			},
		}
		lst := clouduserorgassignment.NewTFProjectRoleAssignments(context.Background(), &groupRoles)
		assert.False(t, lst.IsNull())
		assert.False(t, len(lst.Elements()) == 0)
	})
}

func TestNewOrgUserRolesRequest(t *testing.T) {
	t.Run("null object returns empty OrgUserRolesRequest", func(t *testing.T) {
		obj := types.ObjectNull(clouduserorgassignment.RolesObjectAttrTypes)
		req, diags := clouduserorgassignment.NewOrgUserRolesRequest(context.Background(), obj)
		assert.False(t, diags.HasError())
		assert.NotNil(t, req)
		assert.Empty(t, req.OrgRoles)
	})

	t.Run("populated object returns OrgUserRolesRequest with org_roles", func(t *testing.T) {
		rolesObj, _ := types.ObjectValue(
			clouduserorgassignment.RolesObjectAttrTypes,
			map[string]attr.Value{
				"org_roles":                types.SetValueMust(types.StringType, []attr.Value{types.StringValue("ORG_ROLE1")}),
				"project_role_assignments": types.ListNull(types.ObjectType{AttrTypes: clouduserorgassignment.ProjectRoleAssignmentsAttrType.ElemType.(types.ObjectType).AttrTypes}),
			},
		)
		req, diags := clouduserorgassignment.NewOrgUserRolesRequest(context.Background(), rolesObj)
		assert.False(t, diags.HasError())
		assert.NotNil(t, req)
		assert.Equal(t, []string{"ORG_ROLE1"}, req.OrgRoles)
	})
}
