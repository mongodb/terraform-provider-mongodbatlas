package apikeyprojectassignment_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/apikeyprojectassignment"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20250219001/admin"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.PaginatedApiApiUser
	expectedTFModel *apikeyprojectassignment.TFModel
}

func TestApiKeyProjectAssignmentSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		"Complete SDK response": {
			SDKResp: &admin.PaginatedApiApiUser{
				Results: &[]admin.ApiKeyUserDetails{
					{
						Id: admin.PtrString("TargetAPIKeyID"),
						Roles: &[]admin.CloudAccessRoleAssignment{
							{
								GroupId:  admin.PtrString("TargetProjectID"),
								RoleName: admin.PtrString("MY_ROLE"),
							},
							{
								GroupId:  admin.PtrString("TargetProjectID"),
								RoleName: admin.PtrString("MY_ROLE_2"),
							},
						},
					},
				},
			},
			expectedTFModel: &apikeyprojectassignment.TFModel{
				RoleNames: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("MY_ROLE"),
					types.StringValue("MY_ROLE_2"),
				}),
			},
		},
		"Complete SDK response - No assigned roles": {
			SDKResp: &admin.PaginatedApiApiUser{
				Results: &[]admin.ApiKeyUserDetails{
					{
						Id: admin.PtrString("NotMyTargetAPIKeyID"),
						Roles: &[]admin.CloudAccessRoleAssignment{
							{
								GroupId:  admin.PtrString("TargetProjectID"),
								RoleName: admin.PtrString("MY_ROLE"),
							},
							{
								GroupId:  admin.PtrString("TargetProjectID"),
								RoleName: admin.PtrString("MY_ROLE_2"),
							},
						},
					},
				},
			},
			expectedTFModel: &apikeyprojectassignment.TFModel{
				RoleNames: types.SetNull(nil),
			},
		},
		"Complete SDK response - Wrong project": {
			SDKResp: &admin.PaginatedApiApiUser{
				Results: &[]admin.ApiKeyUserDetails{
					{
						Id: admin.PtrString("TargetAPIKeyID"),
						Roles: &[]admin.CloudAccessRoleAssignment{
							{
								GroupId:  admin.PtrString("NotMyTargetProjectID"),
								RoleName: admin.PtrString("MY_ROLE"),
							},
							{
								GroupId:  admin.PtrString("NotMyTargetProjectID"),
								RoleName: admin.PtrString("MY_ROLE_2"),
							},
						},
					},
				},
			},
			expectedTFModel: &apikeyprojectassignment.TFModel{
				RoleNames: types.SetValueMust(types.StringType, []attr.Value{}),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := apikeyprojectassignment.NewTFModel(t.Context(), tc.SDKResp, "TargetAPIKeyID", "TargetProjectID")
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}

type tfToSDKModelUpdateTestCase struct {
	tfModel        *apikeyprojectassignment.TFModel
	expectedSDKReq *admin.UpdateAtlasProjectApiKey
}

func TestApiKeyProjectAssignmentTFModelToSDKPatch(t *testing.T) {
	testCases := map[string]tfToSDKModelUpdateTestCase{
		"Complete TF state": {
			tfModel: &apikeyprojectassignment.TFModel{
				ApiUserId: types.StringValue("TargetAPIKeyID"),
				ProjectId: types.StringValue("TargetProject"),
				RoleNames: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("MY_ROLE"),
					types.StringValue("MY_ROLE_2"),
				}),
			},
			expectedSDKReq: &admin.UpdateAtlasProjectApiKey{
				Desc: nil,
				Roles: &[]string{
					"MY_ROLE",
					"MY_ROLE_2",
				},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := apikeyprojectassignment.NewAtlasUpdateReq(t.Context(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}

type tfToSDKModelCreateTestCase struct {
	tfModel        *apikeyprojectassignment.TFModel
	expectedSDKReq *[]admin.UserAccessRoleAssignment
}

func TestApiKeyProjectAssignmentTFModelToSDKCreate(t *testing.T) {
	testCases := map[string]tfToSDKModelCreateTestCase{
		"Complete TF state": {
			tfModel: &apikeyprojectassignment.TFModel{
				ApiUserId: types.StringValue("TargetAPIKeyID"),
				ProjectId: types.StringValue("TargetProject"),
				RoleNames: types.SetValueMust(types.StringType, []attr.Value{
					types.StringValue("MY_ROLE"),
					types.StringValue("MY_ROLE_2"),
				}),
			},
			expectedSDKReq: &[]admin.UserAccessRoleAssignment{
				{
					UserId: admin.PtrString("TargetAPIKeyID"),
					Roles: &[]string{
						"MY_ROLE",
						"MY_ROLE_2",
					},
				},
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			apiReqResult, diags := apikeyprojectassignment.NewAtlasCreateReq(t.Context(), tc.tfModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
		})
	}
}
