package streamaccountdetails_test

import (
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/streamaccountdetails"
	"github.com/stretchr/testify/assert"
)

type sdkToTFModelTestCase struct {
	SDKResp         *admin.AccountDetails
	expectedTFModel *streamaccountdetails.TFStreamAccountDetailsModel
	cloudProvider   string
	region          string
}

const (
	dummyProjectID = "111111111111111111111111"
)

func TestStreamAccountDetailsSDKToTFModel(t *testing.T) {
	testCases := map[string]sdkToTFModelTestCase{
		"Complete AWS Account SDK response": {
			cloudProvider: "aws",
			region:        "US_EAST_1",
			SDKResp: &admin.AccountDetails{
				AwsAccountId:  new("123456789"),
				CidrBlock:     new("192.168.0.0/20"),
				CloudProvider: new("aws"),
				VpcId:         new("0987654"),
			},
			expectedTFModel: &streamaccountdetails.TFStreamAccountDetailsModel{
				ProjectId:     types.StringValue(dummyProjectID),
				AwsAccountId:  types.StringValue("123456789"),
				CidrBlock:     types.StringValue("192.168.0.0/20"),
				CloudProvider: types.StringValue("aws"),
				RegionName:    types.StringValue("US_EAST_1"),
				VpcId:         types.StringValue("0987654"),
			},
		},

		"Complete Azure account SDK response": {
			cloudProvider: "azure",
			region:        "EASTUS",
			SDKResp: &admin.AccountDetails{
				CidrBlock:           new("192.168.0.0/20"),
				CloudProvider:       new("azure"),
				AzureSubscriptionId: new("234567890"),
				VirtualNetworkName:  new("876543"),
			},
			expectedTFModel: &streamaccountdetails.TFStreamAccountDetailsModel{
				ProjectId:           types.StringValue(dummyProjectID),
				CidrBlock:           types.StringValue("192.168.0.0/20"),
				CloudProvider:       types.StringValue("azure"),
				AzureSubscriptionId: types.StringValue("234567890"),
				RegionName:          types.StringValue("EASTUS"),
				VirtualNetworkName:  types.StringValue("876543"),
			},
		},
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			resultModel, diags := streamaccountdetails.NewTFStreamAccountDetails(t.Context(), dummyProjectID, tc.cloudProvider, tc.region, tc.SDKResp)
			if diags.HasError() {
				t.Errorf("unexpected errors found: %s", diags.Errors()[0].Summary())
			}
			assert.Equal(t, tc.expectedTFModel, resultModel, "created terraform model did not match expected output")
		})
	}
}
