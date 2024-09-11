package resourcepolicy_test

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/resourcepolicy"
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
	userIDCreate              = "65def6f00f722a1507105ad8"
	userNameCreate            = "mvccpeou"
	userIDUpdate              = "65def6f00f722a1507105ad9"
	userNameUpdate            = "updateUser"
	createdDate               = "2024-09-10T14:59:34Z"
	lastUpdatedDate           = "2024-09-10T14:59:35Z"
	orgID                     = "65def6ce0f722a1507105aa5"
	policyID                  = "66e05ed6680f032312b6b22b"
	policyBody                = "\n\tforbid (\n\t\tprincipal,\n\t\taction == cloud::Action::\"cluster.createEdit\",\n\t\tresource\n\t) when {\n\t\tcontext.cluster.cloudProviders.containsAny([cloud::cloudProvider::\"aws\"])\n\t};"
	policies0ID               = "66e05ed6680f032312b6b22a"
	version                   = "v1"
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

func TestResourcePolicySDKToTFModelFull(t *testing.T) {
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
			SDKModel := parseSDKModel(t, tc.SDKRespJSON)
			ctx := context.Background()
			resultModel, diags := resourcepolicy.NewTFResourcePolicyModel(ctx, &SDKModel)
			assertDiagsOK(t, diags)

			asserter := assert.New(t)
			createdByUserExpected := resourcepolicy.TFUserMetadataModel{
				ID:   types.StringValue(tc.userIDCreate),
				Name: types.StringValue(tc.userNameCreate),
			}
			expectedTFUserCreate, diags := types.ObjectValueFrom(ctx, resourcepolicy.UserMetadataObjectType.AttrTypes, createdByUserExpected)
			assertDiagsOK(t, diags)
			asserter.Equal(expectedTFUserCreate, resultModel.CreatedByUser)

			updatedByUserExpected := resourcepolicy.TFUserMetadataModel{
				ID:   types.StringValue(tc.userIDUpdate),
				Name: types.StringValue(tc.userNameUpdate),
			}
			expectedTFUserUpdate, diags := types.ObjectValueFrom(ctx, resourcepolicy.UserMetadataObjectType.AttrTypes, updatedByUserExpected)
			assertDiagsOK(t, diags)
			asserter.Equal(expectedTFUserUpdate, resultModel.LastUpdatedByUser)

			asserter.Equal(tc.createdDate, resultModel.CreatedDate.ValueString())
			asserter.Equal(tc.policyID, resultModel.ID.ValueString())
			asserter.Equal(tc.lastUpdatedDate, resultModel.LastUpdatedDate.ValueString())
			asserter.Equal(tc.orgID, resultModel.OrgID.ValueString())
			asserter.Equal(tc.version, resultModel.Version.ValueString())

			tfPolicies := asTfPolicies(ctx, t, resultModel)
			asserter.Len(tfPolicies, 1)
			asserter.Equal(resourcepolicy.TFPolicyModel{
				Body: types.StringValue(tc.policyBody),
				ID:   types.StringValue(tc.policies0ID),
			}, tfPolicies[0])
			asserter.Equal(testName, resultModel.Name.ValueString())
		})
	}
}

func TestResourcePolicyNoName(t *testing.T) {
	SDKModel := parseSDKModel(t, clusterForbidCloudProviderNoNameJSON)
	assert.Nil(t, SDKModel.Name)
	ctx := context.Background()
	resultModel, diags := resourcepolicy.NewTFResourcePolicyModel(ctx, &SDKModel)
	assertDiagsOK(t, diags)
	assert.Equal(t, types.StringNull(), resultModel.Name)
}

func TestResourcePolicyMultipleEntries(t *testing.T) {
	SDKModel := parseSDKModel(t, policyMultipleEntriesJSON)
	ctx := context.Background()
	resultModel, diags := resourcepolicy.NewTFResourcePolicyModel(ctx, &SDKModel)
	assertDiagsOK(t, diags)
	tfPolicies := asTfPolicies(ctx, t, resultModel)
	assert.Len(t, tfPolicies, 3)
	for i, expectedSDK := range SDKModel.GetPolicies() {
		assert.Equal(t, expectedSDK.GetId(), tfPolicies[i].ID.ValueString())
		assert.Equal(t, expectedSDK.GetBody(), tfPolicies[i].Body.ValueString())
	}
}

func TestNewApiAtlasPolicyCreate(t *testing.T) {
	ctx := context.Background()
	SDKModel := parseSDKModel(t, clusterForbidCloudProviderJSON)
	resultModel, diags := resourcepolicy.NewTFResourcePolicyModel(ctx, &SDKModel)
	assertDiagsOK(t, diags)
	actualPolicies, diags := resourcepolicy.NewTFPoliciesModelToSDK(ctx, resultModel.Policies)
	assertDiagsOK(t, diags)
	assert.NotNil(t, actualPolicies)
	indexAble := *actualPolicies
	assert.Len(t, indexAble, 1)
	assert.Equal(t, *indexAble[0].Body, policyBody)
}

func asTfPolicies(ctx context.Context, t *testing.T, resultModel *resourcepolicy.TFResourcePolicyModel) []resourcepolicy.TFPolicyModel {
	t.Helper()
	var tfPolicies []resourcepolicy.TFPolicyModel
	diags := resultModel.Policies.ElementsAs(ctx, &tfPolicies, false)
	if diags.HasError() {
		t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
	}
	return tfPolicies
}

func assertDiagsOK(t *testing.T, diags diag.Diagnostics) {
	t.Helper()
	if diags.HasError() {
		t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
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
