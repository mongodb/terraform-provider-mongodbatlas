package projectipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"
)

func NewTFProjectIPAddresses(ctx context.Context, ipAddresses *admin.GroupIPAddresses) (*TFProjectIpAddressesModel, diag.Diagnostics) {
	clusterObjs := make([]TFClusterValueModel, len(ipAddresses.Services.GetClusters()))

	for i, cluster := range ipAddresses.Services.GetClusters() {
		inbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetInbound())
		outbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetOutbound())

		clusterObjs[i] = TFClusterValueModel{
			ClusterName: types.StringPointerValue(cluster.ClusterName),
			Inbound:     inbound,
			Outbound:    outbound,
		}
	}

	servicesObj, diags := types.ObjectValueFrom(ctx, ServicesObjectType.AttrTypes, TFServicesModel{
		Clusters: clusterObjs,
	})
	if diags.HasError() {
		return nil, diags
	}

	return &TFProjectIpAddressesModel{
		ProjectId: types.StringPointerValue(ipAddresses.GroupId),
		Services:  servicesObj,
	}, nil
}
