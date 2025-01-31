package advancedcluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/flexcluster"
)

func isValidUpgradeToFlex(d *schema.ResourceData) bool {
	if d.HasChange("replication_specs.0.region_configs.0") {
		oldProviderName, newProviderName := d.GetChange("replication_specs.0.region_configs.0.provider_name")
		oldInstanceSize, newInstanceSize := d.GetChange("replication_specs.0.region_configs.0.electable_specs.instance_size")
		if oldProviderName == constant.TENANT && newProviderName == flexcluster.FlexClusterType && oldInstanceSize != nil && newInstanceSize == nil {
			return true
		}
	}
	return false
}

func isValidUpdateOfFlex(d *schema.ResourceData) bool {
	updatableAttrHaveBeenUpdated := d.HasChange("tags") || d.HasChange("termination_protection_enabled")
	nonUpdatableAttrHaveNotBeenUpdated := !d.HasChange("cluster_type") && !d.HasChange("replication_specs") && !d.HasChange("project_id") && !d.HasChange("name")
	if updatableAttrHaveBeenUpdated && nonUpdatableAttrHaveNotBeenUpdated {
		return true
	}
	return false
}
