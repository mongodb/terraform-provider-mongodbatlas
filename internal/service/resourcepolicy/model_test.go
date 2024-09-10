package resourcepolicy_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/resourcepolicy"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

var (
	//go:embed testdata/resource_policy_clusterForbidCloudProvider.json
	clusterForbidCloudProviderJSON string
	userIDCreate                   = "65def6f00f722a1507105ad8"
	userNameCreate                 = "mvccpeou"
	userIDUpdate                   = "65def6f00f722a1507105ad9"
	userNameUpdate                 = "updateUser"
	createdDate                    = "2024-09-10T14:59:34Z"
	lastUpdatedDate                = "2024-09-10T14:59:35Z"
	orgID                          = "65def6ce0f722a1507105aa5"
	policyID                       = "66e05ed6680f032312b6b22b"
	policyBody                     = "\n\tforbid (\n\t\tprincipal,\n\t\taction == cloud::Action::\"cluster.createEdit\",\n\t\tresource\n\t) when {\n\t\tcontext.cluster.cloudProviders.containsAny([cloud::cloudProvider::\"aws\"])\n\t};"
	policies0ID                    = "66e05ed6680f032312b6b22a"
	version                        = "v1"
)

type sdkToTFModelTestCase struct {
	SDKRespJSON     string
	userIDCreate    string
	userIDUpdate    string
	userNameCreate  string
	userNameUpdate  string
	createdDate     string
	lastUpdatedDate string
	orgID           string
	policyID        string
	policyBody      string
	policies0ID     string
	version         string
}

func TestResourcePolicySDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		// try no name
		"clusterForbidCloudProvider": {
			SDKRespJSON:     clusterForbidCloudProviderJSON,
			userIDCreate:    userIDCreate,
			userIDUpdate:    userIDUpdate,
			userNameCreate:  userNameCreate,
			userNameUpdate:  userNameUpdate,
			createdDate:     createdDate,
			lastUpdatedDate: lastUpdatedDate,
			orgID:           orgID,
			policyID:        policyID,
			policyBody:      policyBody,
			policies0ID:     policies0ID,
			version:         version,
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			var SDKModel admin.ApiAtlasResourcePolicy
			err := json.Unmarshal([]byte(tc.SDKRespJSON), &SDKModel)
			if err != nil {
				t.Fatalf("failed to unmarshal sdk response: %s", err)
			}
			ctx := context.Background()
			resultModel, diags := resourcepolicy.NewTFResourcePolicy(ctx, &SDKModel)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			asserter := assert.New(t)
			asserter.Equal(types.ObjectValueMust(resourcepolicy.UserMetadataObjectType.AttrTypes, map[string]attr.Value{
				"id":   types.StringValue(tc.userIDCreate),
				"name": types.StringValue(tc.userNameCreate),
			}), resultModel.CreatedByUser)
			asserter.Equal(types.ObjectValueMust(resourcepolicy.UserMetadataObjectType.AttrTypes, map[string]attr.Value{
				"id":   types.StringValue(tc.userIDUpdate),
				"name": types.StringValue(tc.userNameUpdate),
			}), resultModel.LastUpdatedByUser)
			asserter.Equal(tc.createdDate, resultModel.CreatedDate.ValueString())
			asserter.Equal(tc.policyID, resultModel.ID.ValueString())
			asserter.Equal(tc.lastUpdatedDate, resultModel.LastUpdatedDate.ValueString())
			asserter.Equal(tc.orgID, resultModel.OrgID.ValueString())
			asserter.Equal(tc.version, resultModel.Version.ValueString())
			var tfPolicies []resourcepolicy.TFPolicyModel
			diags = resultModel.Policies.ElementsAs(ctx, &tfPolicies, false)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			asserter.Len(tfPolicies, 1)
			asserter.Equal(resourcepolicy.TFPolicyModel{
				Body: types.StringValue(tc.policyBody),
				ID:   types.StringValue(tc.policies0ID),
			}, tfPolicies[0])
		})
	}
}

// type tfToSDKModelTestCase struct {
// 	tfModel        *resourcepolicy.TFResourcePolicyModel
// 	expectedSDKReq *admin.ResourcePolicy
// }

// func TestResourcePolicyTFModelToSDK(t *testing.T) {
// 	testCases := map[string]tfToSDKModelTestCase{
// 		"Complete TF state": {
// 			tfModel: &resourcepolicy.TFResourcePolicyModel{
// 			},
// 			expectedSDKReq: &admin.ResourcePolicy{
// 			},
// 		},
// 	}

// 	for testName, tc := range testCases {
// 		t.Run(testName, func(t *testing.T) {
// 			apiReqResult, diags := resourcepolicy.NewResourcePolicyReq(context.Background(), tc.tfModel)
// 			if diags.HasError() {
// 				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
// 			}
// 			assert.Equal(t, tc.expectedSDKReq, apiReqResult, "created sdk model did not match expected output")
// 		})
// 	}
// }
