package streamaccountdetails

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312003/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewTFStreamAccountDetails(
	ctx context.Context,
	projectID string,
	cloudProvider string,
	region string,
	accountDetails *admin.AccountDetails,
) (*TFStreamAccountDetailsModel, diag.Diagnostics) {
	var accountDetailLinks []TFLinkModel

	if accountDetails.Links != nil {
		accountDetailLinks = make([]TFLinkModel, len(*accountDetails.Links))
		for i, link := range *accountDetails.Links {
			accountDetailLinks[i] = TFLinkModel{
				Href: types.StringPointerValue(link.Href),
				Rel:  types.StringPointerValue(link.Rel),
			}
		}
	}

	links, _ := types.ListValueFrom(ctx, LinkModel, accountDetailLinks)
	return &TFStreamAccountDetailsModel{
		ProjectId:           types.StringValue(projectID),
		AwsAccountId:        types.StringPointerValue(accountDetails.AwsAccountId),
		AzureSubscriptionId: types.StringPointerValue(accountDetails.AzureSubscriptionId),
		CidrBlock:           types.StringPointerValue(accountDetails.CidrBlock),
		CloudProvider:       types.StringValue(cloudProvider),
		Links:               links,
		RegionName:          types.StringValue(region),
		VirtualNetworkName:  types.StringPointerValue(accountDetails.VirtualNetworkName),
		VpcId:               types.StringPointerValue(accountDetails.VpcId),
	}, nil
}
