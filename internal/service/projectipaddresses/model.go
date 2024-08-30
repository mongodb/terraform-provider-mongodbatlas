package projectipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20240805002/admin"
)

func NewTFProjectIPAddresses(ctx context.Context, ipAddresses *admin.GroupIPAddresses) (*TFProjectIpAddressesModel, diag.Diagnostics) {
	clusterObjs := make([]attr.Value, len(ipAddresses.Services.GetClusters()))

	for i, cluster := range ipAddresses.Services.GetClusters() {
		inbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetInbound())
		outbound, _ := types.ListValueFrom(ctx, types.StringType, cluster.GetOutbound())

		clusterObj, _ := types.ObjectValue(ClusterIPsObjectType.AttrTypes, map[string]attr.Value{
			"cluster_name": types.StringPointerValue(cluster.ClusterName),
			"inbound":      inbound,
			"outbound":     outbound,
		})

		clusterObjs[i] = clusterObj
	}

	clustersList := types.ListValueMust(types.ObjectType{AttrTypes: ClusterIPsObjectType.AttrTypes}, clusterObjs)

	servicesObj, diags := types.ObjectValue(ServicesObjectType.AttrTypes, map[string]attr.Value{
		"clusters": clustersList,
	})
	if diags.HasError() {
		return nil, diags
	}

	return &TFProjectIpAddressesModel{
		ProjectId: types.StringPointerValue(ipAddresses.GroupId),
		Services:  servicesObj,
	}, nil
}
