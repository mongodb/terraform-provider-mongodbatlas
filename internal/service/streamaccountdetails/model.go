package streamaccountdetails

import (
	"go.mongodb.org/atlas-sdk/v20250312022/admin"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewTFStreamAccountDetails(
	projectID string,
	cloudProvider string,
	region string,
	accountDetails *admin.AccountDetails,
) *TFStreamAccountDetailsModel {
	return &TFStreamAccountDetailsModel{
		ProjectId:           types.StringValue(projectID),
		AwsAccountId:        types.StringPointerValue(accountDetails.AwsAccountId),
		AzureSubscriptionId: types.StringPointerValue(accountDetails.AzureSubscriptionId),
		CidrBlock:           types.StringPointerValue(accountDetails.CidrBlock),
		CloudProvider:       types.StringValue(cloudProvider),
		RegionName:          types.StringValue(region),
		VirtualNetworkName:  types.StringPointerValue(accountDetails.VirtualNetworkName),
		VpcId:               types.StringPointerValue(accountDetails.VpcId),
	}
}
