package resourcepolicy_test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/resourcepolicy"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/unit"
)

var (
	//go:embed testdata/policy_clusterForbidCloudProvider.json
	clusterForbidCloudProviderJSON string
	//go:embed testdata/policy_multipleEntries.json
	policyMultipleEntriesJSON string
)

type tfModelTestCase struct {
	description     *string
	name            string
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

func (tc *tfModelTestCase) addDefaults() {
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

func createTFModel(t *testing.T, testCase *tfModelTestCase) *resourcepolicy.TFModel {
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
	return &resourcepolicy.TFModel{
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
		Name:            types.StringValue(testCase.name),
		OrgID:           types.StringValue(testCase.orgID),
		Version:         types.StringValue(testCase.version),
		Description:     types.StringPointerValue(testCase.description),
	}
}

func TestNewTFModel(t *testing.T) {
	testCases := map[string]tfModelTestCase{
		"clusterForbidCloudProvider": {
			name:           "clusterForbidCloudProvider",
			SDKRespJSON:    clusterForbidCloudProviderJSON,
			userIDUpdate:   "65def6f00f722a1507105ad9",
			userNameUpdate: "updateUser",
		},
		"policyMultipleEntriesJSON": {
			SDKRespJSON:     policyMultipleEntriesJSON,
			name:            "multipleEntries",
			createdDate:     "2024-09-11T13:36:18Z",
			lastUpdatedDate: "2024-09-11T13:36:18Z",
			policyID:        "66e19cd2fdc0332d1fa5e877",
			description:     conversion.StringPtr("test description"),
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			SDKModel := parseSDKModel(t, tc.SDKRespJSON)
			expectedModel := createTFModel(t, &tc)
			resultModel, diags := resourcepolicy.NewTFModel(t.Context(), &SDKModel)
			unit.AssertDiagsOK(t, diags)
			assert.Equal(t, expectedModel, resultModel)
		})
	}
}

func TestNewUserMetadataObjectTypeWithNilArg(t *testing.T) {
	var metadataNil *admin.ApiAtlasUserMetadata
	diags := diag.Diagnostics{}
	obj := resourcepolicy.NewUserMetadataObjectType(t.Context(), metadataNil, &diags)
	unit.AssertDiagsOK(t, diags)
	assert.Equal(t, types.ObjectNull(resourcepolicy.UserMetadataObjectType.AttrTypes), obj)
}

func TestNewAdminPolicies(t *testing.T) {
	policies := []resourcepolicy.TFPolicyModel{
		{
			Body: types.StringValue("policy1"),
			ID:   types.StringValue("id1"),
		},
		{
			Body: types.StringValue("policy2"),
		},
	}
	apiModels := resourcepolicy.NewAdminPolicies(t.Context(), policies)
	assert.Len(t, apiModels, 2)
	assert.Equal(t, "policy1", apiModels[0].GetBody())
	assert.Equal(t, "policy2", apiModels[1].GetBody())
}

func TestNewTFModelDSP(t *testing.T) {
	orgID := "65def6ce0f722a1507105aa5"
	input := []admin.ApiAtlasResourcePolicy{
		parseSDKModel(t, clusterForbidCloudProviderJSON),
		parseSDKModel(t, policyMultipleEntriesJSON),
	}
	resultModel, diags := resourcepolicy.NewTFModelDSP(t.Context(), orgID, input)
	unit.AssertDiagsOK(t, diags)
	assert.Len(t, resultModel.ResourcePolicies, 2)

	assert.Equal(t, orgID, resultModel.OrgID.ValueString())
}

func TestNewTFModelDSPEmptyModel(t *testing.T) {
	orgID := "65def6ce0f722a1507105aa5"
	resultModel, diags := resourcepolicy.NewTFModelDSP(t.Context(), orgID, []admin.ApiAtlasResourcePolicy{})
	unit.AssertDiagsOK(t, diags)
	assert.Empty(t, resultModel.ResourcePolicies)
	assert.Equal(t, orgID, resultModel.OrgID.ValueString())
}

func TestModifyPlan_SkipsValidationWhenValuesUnknownOrNull(t *testing.T) {
	testCases := map[string]struct {
		orgIDValue any
		nameValue  any
	}{
		"unknown org_id": {
			orgIDValue: tftypes.UnknownValue,
			nameValue:  "test-name",
		},
		"unknown name": {
			orgIDValue: "65def6ce0f722a1507105aa5",
			nameValue:  tftypes.UnknownValue,
		},
		"null org_id": {
			orgIDValue: nil,
			nameValue:  "test-name",
		},
		"null name": {
			orgIDValue: "65def6ce0f722a1507105aa5",
			nameValue:  nil,
		},
		"both unknown": {
			orgIDValue: tftypes.UnknownValue,
			nameValue:  tftypes.UnknownValue,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			plan := buildResourcePolicyPlan(t, tc.orgIDValue, tc.nameValue)
			req := resource.ModifyPlanRequest{Plan: plan}
			resp := &resource.ModifyPlanResponse{Plan: plan}

			rs := resourcepolicy.Resource()
			rsMp, ok := rs.(resource.ResourceWithModifyPlan)
			if !ok {
				t.Fatal("resource does not implement ResourceWithModifyPlan")
			}
			rsMp.ModifyPlan(t.Context(), req, resp)

			assert.False(t, resp.Diagnostics.HasError())
		})
	}
}

func buildResourcePolicyPlan(t *testing.T, orgIDValue, nameValue any) tfsdk.Plan {
	t.Helper()
	ctx := t.Context()
	schema := resourcepolicy.ResourceSchema(ctx)
	userMetadataType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":   tftypes.String,
			"name": tftypes.String,
		},
	}
	policyType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"body": tftypes.String,
			"id":   tftypes.String,
		},
	}
	raw := tftypes.NewValue(schema.Type().TerraformType(ctx), map[string]tftypes.Value{
		"org_id":               tftypes.NewValue(tftypes.String, orgIDValue),
		"name":                 tftypes.NewValue(tftypes.String, nameValue),
		"description":          tftypes.NewValue(tftypes.String, nil),
		"id":                   tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"version":              tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"created_date":         tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"last_updated_date":    tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"created_by_user":      tftypes.NewValue(userMetadataType, tftypes.UnknownValue),
		"last_updated_by_user": tftypes.NewValue(userMetadataType, tftypes.UnknownValue),
		"policies": tftypes.NewValue(tftypes.List{ElementType: policyType}, []tftypes.Value{
			tftypes.NewValue(policyType, map[string]tftypes.Value{
				"body": tftypes.NewValue(tftypes.String, "forbid(principal, action, resource);"),
				"id":   tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			}),
		}),
	})
	return tfsdk.Plan{
		Schema: schema,
		Raw:    raw,
	}
}
