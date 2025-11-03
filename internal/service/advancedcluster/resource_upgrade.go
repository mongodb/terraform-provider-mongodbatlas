package advancedcluster

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
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
	req := admin.LegacyAtlasTenantClusterUpgradeRequest{
		Name: state.GetName(),
		ProviderSettings: &admin.ClusterProviderSettings{
			ProviderName:     newProviderName,
			RegionName:       newRegion.RegionName,
			InstanceSizeName: newRegion.GetElectableSpecs().InstanceSize,
		},
	}
	if patch.GetBackupEnabled() {
		// ProviderBackupEnabled must be used instead of BackupEnabled for tenant upgrade request, details in CLOUDP-327109
		req.ProviderBackupEnabled = conversion.Pointer(true)
	}
	return &req
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
	req := admin.AtlasTenantClusterUpgradeRequest20240805{
		Name:             state.GetName(),
		ClusterType:      state.ClusterType,
		ReplicationSpecs: patch.ReplicationSpecs,
	}

	// checking for state value as a flex cluster can already have backup enabled
	if state.GetBackupEnabled() || patch.GetBackupEnabled() {
		req.BackupEnabled = conversion.Pointer(true)
	}
	return &req
}
