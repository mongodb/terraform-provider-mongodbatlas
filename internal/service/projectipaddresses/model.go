package projectipaddresses

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.mongodb.org/atlas-sdk/v20240805001/admin"
)

func NewTFProjectIPAddresses(ctx context.Context, ipAddresses *admin.GroupIPAddresses) (types.Object, diag.Diagnostics) {
	clusterIPs := []ClustersValue{}
	if ipAddresses != nil && ipAddresses.Services != nil {
		clusterIPAddresses := ipAddresses.Services.GetClusters()
		clusterIPs = make([]ClustersValue, len(clusterIPAddresses))
		for i := range clusterIPAddresses {
			inbound, _ := types.ListValueFrom(ctx, types.StringType, clusterIPAddresses[i].GetInbound())
			outbound, _ := types.ListValueFrom(ctx, types.StringType, clusterIPAddresses[i].GetOutbound())
			clusterIPs[i] = ClustersValue{
				ClusterName: types.StringPointerValue(clusterIPAddresses[i].ClusterName),
				Inbound:     inbound,
				Outbound:    outbound,
			}
		}
	}
	obj, diags := types.ObjectValueFrom(ctx, IPAddressesObjectType.AttrTypes, ProjectIpAddressesModel{
		Services: ServicesValue{
			Clusters: clusterIPs,
		},
	})
	return obj, diags
}
