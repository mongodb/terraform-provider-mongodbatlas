package advancedclustertpf

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"go.mongodb.org/atlas-sdk/v20241113005/admin"
)

func getUpgradeTenantRequest(state, patch *admin.ClusterDescription20240805) *admin.LegacyAtlasTenantClusterUpgradeRequest {
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

func getUpgradeFlexToDedicatedRequest(state, patch *admin.ClusterDescription20240805) *admin.AtlasTenantClusterUpgradeRequest20240805 {
	if patch.ReplicationSpecs == nil {
		return nil
	}
	(*patch.ReplicationSpecs)[0].Id = nil
	(*patch.ReplicationSpecs)[0].ZoneId = nil
	oldRegion := state.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	oldProviderName := oldRegion.GetProviderName()
	newRegion := patch.GetReplicationSpecs()[0].GetRegionConfigs()[0]
	newProviderName := newRegion.GetProviderName()
	if oldProviderName != flexcluster.FlexClusterType || newProviderName == flexcluster.FlexClusterType {
		return nil
	}
	return &admin.AtlasTenantClusterUpgradeRequest20240805{
		Name:             state.GetName(),
		ClusterType:      state.ClusterType,
		ReplicationSpecs: patch.ReplicationSpecs,
	}
}
