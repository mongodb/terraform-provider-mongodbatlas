package acc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go.mongodb.org/atlas-sdk/v20250312007/admin"
)

func ClusterDatasourceHcl(req *ClusterRequest) (configStr, clusterName, resourceName string, err error) {
	if req == nil || req.ProjectID == "" || req.ClusterName == "" {
		return "", "", "", errors.New("must specify a ClusterRequest with at least ProjectID and ClusterName set")
	}
	req.AddDefaults()
	f := hclwrite.NewEmptyFile()
	root := f.Body()
	resourceType := "mongodbatlas_advanced_cluster"
	resourceSuffix := req.ResourceSuffix
	cluster := root.AppendNewBlock("data", []string{resourceType, resourceSuffix}).Body()
	clusterResourceName := fmt.Sprintf("data.%s.%s", resourceType, resourceSuffix)
	clusterName = req.ClusterName
	clusterRootAttributes := map[string]any{
		"name": clusterName,
	}
	projectID := req.ProjectID
	if strings.Contains(req.ProjectID, ".") {
		err = setAttributeHcl(cluster, fmt.Sprintf("project_id = %s", projectID))
		if err != nil {
			return "", "", "", fmt.Errorf("failed to set project_id = %s", projectID)
		}
	} else {
		clusterRootAttributes["project_id"] = projectID
	}
	addPrimitiveAttributes(cluster, clusterRootAttributes)
	return "\n" + string(f.Bytes()), clusterName, clusterResourceName, err
}

func ClusterResourceHcl(req *ClusterRequest) (configStr, clusterName, resourceName string, err error) {
	if req == nil || req.ProjectID == "" {
		return "", "", "", errors.New("must specify a ClusterRequest with at least ProjectID set")
	}
	projectID := req.ProjectID
	req.AddDefaults()
	specRequests := req.ReplicationSpecs
	specs := make([]admin.ReplicationSpec20240805, len(specRequests))
	for i := range specRequests {
		specRequest := specRequests[i]
		specs[i] = replicationSpec(&specRequest)
	}
	clusterName = req.ClusterName
	resourceSuffix := req.ResourceSuffix
	clusterType := req.ClusterType()

	f := hclwrite.NewEmptyFile()
	root := f.Body()
	resourceType := "mongodbatlas_advanced_cluster"
	cluster := root.AppendNewBlock("resource", []string{resourceType, resourceSuffix}).Body()
	clusterRootAttributes := map[string]any{
		"cluster_type":           clusterType,
		"name":                   clusterName,
		"backup_enabled":         req.CloudBackup,
		"pit_enabled":            req.PitEnabled,
		"mongo_db_major_version": req.MongoDBMajorVersion,
	}
	if strings.Contains(req.ProjectID, ".") {
		err = setAttributeHcl(cluster, fmt.Sprintf("project_id = %s", projectID))
		if err != nil {
			return "", "", "", fmt.Errorf("failed to set project_id = %s", projectID)
		}
	} else {
		clusterRootAttributes["project_id"] = projectID
	}
	if req.DiskSizeGb != 0 {
		clusterRootAttributes["disk_size_gb"] = req.DiskSizeGb
	}
	if req.RetainBackupsEnabled {
		clusterRootAttributes["retain_backups_enabled"] = req.RetainBackupsEnabled
	}
	addPrimitiveAttributes(cluster, clusterRootAttributes)
	cluster.AppendNewline()
	if len(req.AdvancedConfiguration) > 0 {
		for _, key := range sortStringMapKeysAny(req.AdvancedConfiguration) {
			if !knownAdvancedConfig[key] {
				return "", "", "", fmt.Errorf("unknown key in advanced configuration: %s", key)
			}
		}
		advancedClusterBlock := cluster.AppendNewBlock("advanced_configuration", nil).Body()
		addPrimitiveAttributes(advancedClusterBlock, req.AdvancedConfiguration)
		cluster.AppendNewline()
	}
	for i, spec := range specs {
		err = writeReplicationSpec(cluster, spec)
		if err != nil {
			return "", "", "", fmt.Errorf("error writing hcl for replication spec %d: %w", i, err)
		}
	}
	if len(req.Tags) > 0 {
		for _, key := range SortStringMapKeys(req.Tags) {
			value := req.Tags[key]
			tagBlock := cluster.AppendNewBlock("tags", nil).Body()
			tagBlock.SetAttributeValue("key", cty.StringVal(key))
			tagBlock.SetAttributeValue("value", cty.StringVal(value))
		}
	}
	cluster.AppendNewline()
	if req.ResourceDependencyName != "" {
		if !strings.Contains(req.ResourceDependencyName, ".") {
			return "", "", "", fmt.Errorf("req.ResourceDependencyName must have a '.'")
		}
		err = setAttributeHcl(cluster, fmt.Sprintf("depends_on = [%s]", req.ResourceDependencyName))
		if err != nil {
			return "", "", "", err
		}
	}
	clusterResourceName := fmt.Sprintf("%s.%s", resourceType, resourceSuffix)
	return "\n" + string(f.Bytes()), clusterName, clusterResourceName, err
}

func writeReplicationSpec(cluster *hclwrite.Body, spec admin.ReplicationSpec20240805) error {
	replicationBlock := cluster.AppendNewBlock("replication_specs", nil).Body()
	err := addPrimitiveAttributesViaJSON(replicationBlock, spec)
	if err != nil {
		return err
	}
	for _, rc := range spec.GetRegionConfigs() {
		if rc.Priority == nil {
			rc.SetPriority(7)
		}
		replicationBlock.AppendNewline()
		rcBlock := replicationBlock.AppendNewBlock("region_configs", nil).Body()
		err = addPrimitiveAttributesViaJSON(rcBlock, rc)
		if err != nil {
			return err
		}
		autoScalingBlock := rcBlock.AppendNewBlock("auto_scaling", nil).Body()
		if rc.AutoScaling == nil {
			autoScalingBlock.SetAttributeValue("disk_gb_enabled", cty.BoolVal(false))
		} else {
			autoScaling := rc.GetAutoScaling()
			asDisk := autoScaling.GetDiskGB()
			autoScalingBlock.SetAttributeValue("disk_gb_enabled", cty.BoolVal(asDisk.GetEnabled()))
			if autoScaling.Compute != nil {
				return fmt.Errorf("auto_scaling.compute is not supportd yet %v", autoScaling)
			}
		}
		nodeSpec := rc.GetElectableSpecs()
		nodeSpecBlock := rcBlock.AppendNewBlock("electable_specs", nil).Body()
		err = addPrimitiveAttributesViaJSON(nodeSpecBlock, nodeSpec)

		readOnlySpecs := rc.GetReadOnlySpecs()
		if readOnlySpecs.GetNodeCount() != 0 {
			readOnlyBlock := rcBlock.AppendNewBlock("read_only_specs", nil).Body()
			err = addPrimitiveAttributesViaJSON(readOnlyBlock, readOnlySpecs)
		}
	}
	return err
}
