package advancedclustertpf

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"go.mongodb.org/atlas-sdk/v20241113002/admin"
)

func getTenantUpgradeRequest(state, patch *admin.ClusterDescription20240805) *admin.LegacyAtlasTenantClusterUpgradeRequest {
	if patch.ReplicationSpecs == nil {
		return nil
	}
	oldRegion := state.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	oldProviderName := oldRegion.GetProviderName()
	newRegion := patch.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	newProviderName := newRegion.GetProviderName()
	if oldProviderName != constant.TENANT || newProviderName == constant.TENANT {
		return nil
	}
	return &admin.LegacyAtlasTenantClusterUpgradeRequest{
		Name: state.GetName(),
		ProviderSettings: &admin.ClusterProviderSettings{
			ProviderName:     newProviderName,
			RegionName:       newRegion.RegionName,
			InstanceSizeName: newRegion.GetElectableSpecs().InstanceSize,
		},
	}
}
