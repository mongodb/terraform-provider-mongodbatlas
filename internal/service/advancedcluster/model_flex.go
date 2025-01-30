package advancedcluster

import (
	"go.mongodb.org/atlas-sdk/v20241113004/admin"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
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

// TODO: refactor this. 1. Move to tpf. Instead of d as parameter, pass tags and bool
func getFlexClusterUpdateRequest(d *schema.ResourceData) *admin.FlexClusterDescriptionUpdate20241113 {
	return &admin.FlexClusterDescriptionUpdate20241113{
		Tags:                         conversion.ExpandTagsFromSetSchema(d),
		TerminationProtectionEnabled: conversion.Pointer(d.Get("termination_protection_enabled").(bool)),
	}
}

func getUpgradeToFlexClusterRequest() *admin.LegacyAtlasTenantClusterUpgradeRequest {
	// WIP: will be finished as part of CLOUDP-296220
	return &admin.LegacyAtlasTenantClusterUpgradeRequest{
		ProviderSettings: &admin.ClusterProviderSettings{
			ProviderName: flexcluster.FlexClusterType,
		},
	}
}
