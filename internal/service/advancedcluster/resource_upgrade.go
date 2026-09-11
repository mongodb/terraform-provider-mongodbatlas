package advancedcluster

import (
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
	"go.mongodb.org/atlas-sdk/v20250312024/admin"
)

// The legacy SDK request omits databaseEdition, which the tenant upgrade API accepts.
type tenantUpgradeRequest struct {
	DatabaseEdition       *string                        `json:"databaseEdition,omitempty"`
	ProviderSettings      *admin.ClusterProviderSettings `json:"providerSettings,omitempty"`
	ProviderBackupEnabled *bool                          `json:"providerBackupEnabled,omitempty"`
	ReplicationFactor     *int                           `json:"replicationFactor,omitempty"`
	Name                  string                         `json:"name"`
}

func getUpgradeTenantRequest(state, plan, patch *admin.ClusterDescription20240805) *tenantUpgradeRequest {
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
	req := tenantUpgradeRequest{
		Name:            state.GetName(),
		DatabaseEdition: plan.DatabaseEdition,
		ProviderSettings: &admin.ClusterProviderSettings{
			ProviderName:     newProviderName,
			RegionName:       newRegion.RegionName,
			InstanceSizeName: newRegion.GetElectableSpecs().InstanceSize,
		},
	}
	if patch.GetBackupEnabled() {
		// ProviderBackupEnabled must be used instead of BackupEnabled for tenant upgrade request, details in CLOUDP-327109
		req.ProviderBackupEnabled = new(true)
	}
	// Override the legacy three-node default for configurations such as two-node Infinite clusters.
	if nodeCount := newRegion.GetElectableSpecs().NodeCount; nodeCount != nil && *nodeCount != 3 {
		req.ReplicationFactor = nodeCount
	}
	return &req
}

func getUpgradeFlexToDedicatedRequest(state, plan, patch *admin.ClusterDescription20240805) *admin.AtlasTenantClusterUpgradeRequest20240805 {
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
		DatabaseEdition:  plan.DatabaseEdition,
		ReplicationSpecs: patch.ReplicationSpecs,
	}

	// checking for state value as a flex cluster can already have backup enabled
	if state.GetBackupEnabled() || patch.GetBackupEnabled() {
		req.BackupEnabled = new(true)
	}
	return &req
}
