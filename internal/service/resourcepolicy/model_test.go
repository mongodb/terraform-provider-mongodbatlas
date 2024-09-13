package resourcepolicy_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/resourcepolicy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas-sdk/v20240805003/admin"
)

var (
	//go:embed testdata/policy_clusterForbidCloudProvider.json
	clusterForbidCloudProviderJSON string
	//go:embed testdata/policy_clusterForbidCloudProvider_no_name.json
	clusterForbidCloudProviderNoNameJSON string
	//go:embed testdata/policy_multipleEntries.json
	policyMultipleEntriesJSON string
)

type sdkToTFModelTestCase struct {
	name            *string
	SDKRespJSON     string
	userIDCreate    string
	userIDUpdate    string
	userNameCreate  string
	userNameUpdate  string
	createdDate     string
	lastUpdatedDate string
	orgID           string
	policyID        string
	version         string
}

func (tc *sdkToTFModelTestCase) addDefaults() {
	if tc.userIDCreate == "" {
		tc.userIDCreate = "65def6f00f722a1507105ad8"
	}
	if tc.userIDUpdate == "" {
		tc.userIDUpdate = "65def6f00f722a1507105ad8"
	}
	if tc.userNameCreate == "" {
		tc.userNameCreate = "mvccpeou"
	}
	if tc.userNameUpdate == "" {
		tc.userNameUpdate = "mvccpeou"
	}
	if tc.createdDate == "" {
		tc.createdDate = "2024-09-10T14:59:34Z"
	}
	if tc.lastUpdatedDate == "" {
		tc.lastUpdatedDate = "2024-09-10T14:59:35Z"
	}
	if tc.orgID == "" {
		tc.orgID = "65def6ce0f722a1507105aa5"
	}
	if tc.policyID == "" {
		tc.policyID = "66e05ed6680f032312b6b22b"
	}
	if tc.version == "" {
		tc.version = "v1"
	}
}

func parseSDKModel(t *testing.T, sdkRespJSON string) admin.ApiAtlasResourcePolicy {
	t.Helper()
	var SDKModel admin.ApiAtlasResourcePolicy
	err := json.Unmarshal([]byte(sdkRespJSON), &SDKModel)
	if err != nil {
		t.Fatalf("failed to unmarshal sdk response: %s", err)
	}
	return SDKModel
}

func createTFModel(t *testing.T, testCase *sdkToTFModelTestCase) *resourcepolicy.TFResourcePolicyModel {
	t.Helper()
	testCase.addDefaults()
	adminModel := parseSDKModel(t, testCase.SDKRespJSON)
	policies := make([]resourcepolicy.TFPolicyModel, len(adminModel.GetPolicies()))
	for i, policy := range adminModel.GetPolicies() {
		policies[i] = resourcepolicy.TFPolicyModel{
			Body: types.StringPointerValue(policy.Body),
			ID:   types.StringPointerValue(policy.Id),
		}
	}
	return &resourcepolicy.TFResourcePolicyModel{
		CreatedByUser: unit.TFObjectValue(t, resourcepolicy.UserMetadataObjectType, resourcepolicy.TFUserMetadataModel{
			ID:   types.StringValue(testCase.userIDCreate),
			Name: types.StringValue(testCase.userNameCreate),
		}),
		LastUpdatedByUser: unit.TFObjectValue(t, resourcepolicy.UserMetadataObjectType, resourcepolicy.TFUserMetadataModel{
			ID:   types.StringValue(testCase.userIDUpdate),
			Name: types.StringValue(testCase.userNameUpdate),
		}),
		Policies:        policies,
		CreatedDate:     types.StringValue(testCase.createdDate),
		ID:              types.StringValue(testCase.policyID),
		LastUpdatedDate: types.StringValue(testCase.lastUpdatedDate),
		Name:            types.StringPointerValue(testCase.name),
		OrgID:           types.StringValue(testCase.orgID),
		Version:         types.StringValue(testCase.version),
	}
}

func TestResourcePolicySDKToTFModelFull(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		// try no name
		"clusterForbidCloudProvider": {
			name:           conversion.StringPtr("clusterForbidCloudProvider"),
			SDKRespJSON:    clusterForbidCloudProviderJSON,
			userIDUpdate:   "65def6f00f722a1507105ad9",
			userNameUpdate: "updateUser",
		},
		"policyNoName": {
			SDKRespJSON:    clusterForbidCloudProviderNoNameJSON,
			userIDUpdate:   "65def6f00f722a1507105ad9",
			userNameUpdate: "updateUser",
		},
		"policyMultipleEntriesJSON": {
			SDKRespJSON:     policyMultipleEntriesJSON,
			name:            conversion.StringPtr("multipleEntries"),
			createdDate:     "2024-09-11T13:36:18Z",
			lastUpdatedDate: "2024-09-11T13:36:18Z",
			policyID:        "66e19cd2fdc0332d1fa5e877",
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			SDKModel := parseSDKModel(t, tc.SDKRespJSON)
			ctx := context.Background()
			expectedModel := createTFModel(t, &tc)
			resultModel, diags := resourcepolicy.NewTFResourcePolicyModel(ctx, &SDKModel)
			unit.AssertDiagsOK(t, diags)
			assert.Equal(t, expectedModel, resultModel)
		})
	}
}

func Test_NewUserMetadataObjectTypeWithNilArg(t *testing.T) {
	ctx := context.Background()
	var metadataNil *admin.ApiAtlasUserMetadata
	diags := diag.Diagnostics{}
	obj := resourcepolicy.NewUserMetadataObjectType(ctx, metadataNil, &diags)
	unit.AssertDiagsOK(t, diags)
	assert.Equal(t, types.ObjectNull(resourcepolicy.UserMetadataObjectType.AttrTypes), obj)
}
