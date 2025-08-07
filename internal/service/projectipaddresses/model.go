package projectipaddresses

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20250312006/admin"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewTFProjectIPAddresses(ctx context.Context, ipAddresses *admin.GroupIPAddresses) (*TFProjectIpAddressesModel, diag.Diagnostics) {
	clusterObjs := make([]TFClusterValueModel, len(ipAddresses.Services.GetClusters()))

	for i, cluster := range ipAddresses.Services.GetClusters() {
		inbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetInbound())
		outbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetOutbound())
		futureInbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetFutureInbound())
		futureOutbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetFutureOutbound())

		clusterObjs[i] = TFClusterValueModel{
			ClusterName:    types.StringPointerValue(cluster.ClusterName),
			Inbound:        inbound,
			Outbound:       outbound,
			FutureInbound:  futureInbound,
			FutureOutbound: futureOutbound,
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
