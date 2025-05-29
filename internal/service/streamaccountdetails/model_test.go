package streamaccountdetails_test

import (
	"testing"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"

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
	linksInput := []admin.Link{
		{
			Href: admin.PtrString("foo_href"),
			Rel:  admin.PtrString("bar_rel"),
		},
		{
			Href: admin.PtrString("foo2_href"),
			Rel:  admin.PtrString("bar2_rel"),
		},
	}

	linksTF := []streamaccountdetails.TFLinkModel{
		{
			Href: types.StringValue("foo_href"),
			Rel:  types.StringValue("bar_rel"),
		},
		{
			Href: types.StringValue("foo2_href"),
			Rel:  types.StringValue("bar2_rel"),
		},
	}

	linksOutput, _ := types.ListValueFrom(t.Context(), streamaccountdetails.LinkModel, linksTF)

	testCases := map[string]sdkToTFModelTestCase{
		"Complete AWS Account SDK response": {
			cloudProvider: "aws",
			region:        "US_EAST_1",
			SDKResp: &admin.AccountDetails{
				AwsAccountId:  admin.PtrString("123456789"),
				CidrBlock:     admin.PtrString("192.168.0.0/20"),
				CloudProvider: admin.PtrString("aws"),
				Links:         &linksInput,
				VpcId:         admin.PtrString("0987654"),
			},
			expectedTFModel: &streamaccountdetails.TFStreamAccountDetailsModel{
				ProjectId:     types.StringValue(dummyProjectID),
				AwsAccountId:  types.StringValue("123456789"),
				CidrBlock:     types.StringValue("192.168.0.0/20"),
				CloudProvider: types.StringValue("aws"),
				Links:         linksOutput,
				RegionName:    types.StringValue("US_EAST_1"),
				VpcId:         types.StringValue("0987654"),
			},
		},

		"AWS Account SDK response with no links": {
			cloudProvider: "aws",
			region:        "US_EAST_1",
			SDKResp: &admin.AccountDetails{
				AwsAccountId:  admin.PtrString("123456789"),
				CidrBlock:     admin.PtrString("192.168.0.0/20"),
				CloudProvider: admin.PtrString("aws"),
				VpcId:         admin.PtrString("0987654"),
			},
			expectedTFModel: &streamaccountdetails.TFStreamAccountDetailsModel{
				ProjectId:     types.StringValue(dummyProjectID),
				AwsAccountId:  types.StringValue("123456789"),
				CidrBlock:     types.StringValue("192.168.0.0/20"),
				CloudProvider: types.StringValue("aws"),
				Links:         types.ListNull(streamaccountdetails.LinkModel),
				RegionName:    types.StringValue("US_EAST_1"),
				VpcId:         types.StringValue("0987654"),
			},
		},

		"Complete Azure account SDK response": {
			cloudProvider: "azure",
			region:        "EASTUS",
			SDKResp: &admin.AccountDetails{
				CidrBlock:           admin.PtrString("192.168.0.0/20"),
				CloudProvider:       admin.PtrString("azure"),
				Links:               &linksInput,
				AzureSubscriptionId: admin.PtrString("234567890"),
				VirtualNetworkName:  admin.PtrString("876543"),
			},
			expectedTFModel: &streamaccountdetails.TFStreamAccountDetailsModel{
				ProjectId:           types.StringValue(dummyProjectID),
				CidrBlock:           types.StringValue("192.168.0.0/20"),
				CloudProvider:       types.StringValue("azure"),
				Links:               linksOutput,
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
