package advancedcluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

func isValidUpgradeToFlex(d *schema.ResourceData) bool {
	if d.HasChange("replication_specs.0.region_configs.0") {
		oldProviderName, newProviderName := d.GetChange("replication_specs.0.region_configs.0.provider_name")
		if oldProviderName == constant.TENANT && newProviderName == flexcluster.FlexClusterType {
			return true
		}
	}
	return false
}

func isValidUpdateOfFlex(d *schema.ResourceData) bool {
	updatableAttrHaveBeenUpdated := d.HasChange("tags") || d.HasChange("termination_protection_enabled")
	nonUpdatableAttrHaveNotBeenUpdated := !d.HasChange("cluster_type") && !d.HasChange("replication_specs") && !d.HasChange("project_id") && !d.HasChange("name")
	return updatableAttrHaveBeenUpdated && nonUpdatableAttrHaveNotBeenUpdated
}

func isUpgradeFromFlex(d *schema.ResourceData) bool {
	if d.HasChange("replication_specs.0.region_configs.0") {
		oldProviderName, _ := d.GetChange("replication_specs.0.region_configs.0.provider_name")
		return oldProviderName == flexcluster.FlexClusterType
	}
	return false
}

func GetUpgradeToDedicatedClusterRequest(d *schema.ResourceData) *admin.AtlasTenantClusterUpgradeRequest20240805 {
	clusterName := d.Get("name").(string)
	var rootDiskSizeGB *float64
	if v, ok := d.GetOk("disk_size_gb"); ok {
		rootDiskSizeGB = conversion.Pointer(v.(float64))
	}
	return &admin.AtlasTenantClusterUpgradeRequest20240805{
		Name:             clusterName,
		ClusterType:      conversion.Pointer(d.Get("cluster_type").(string)),
		ReplicationSpecs: expandAdvancedReplicationSpecs(d.Get("replication_specs").([]any), rootDiskSizeGB),
	}
}
