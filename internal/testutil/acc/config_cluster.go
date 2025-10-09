package acc

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312008/admin"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
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
	setAttributes(cluster, clusterRootAttributes)
	return "\n" + string(f.Bytes()), clusterName, clusterResourceName, err
}

// ClusterResourceHcl generates Terraform HCL configuration for MongoDB Atlas advanced clusters.
// It converts a ClusterRequest into valid Terraform HCL string
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

	if req.RetainBackupsEnabled {
		clusterRootAttributes["retain_backups_enabled"] = req.RetainBackupsEnabled
	}
	setAttributes(cluster, clusterRootAttributes)

	if err := validateAdvancedConfig(req.AdvancedConfiguration); err != nil {
		return "", "", "", err
	}
	if len(req.AdvancedConfiguration) > 0 {
		cluster.AppendNewline()
		setAttributes(cluster, map[string]any{
			"advanced_configuration": req.AdvancedConfiguration,
		})
	}

	cluster.AppendNewline()
	err = writeReplicationSpec(cluster, specs)
	if err != nil {
		return "", "", "", fmt.Errorf("error writing hcl for replication specs: %w", err)
	}

	if len(req.Tags) > 0 {
		tagMap := make(map[string]cty.Value, len(req.Tags))
		for _, key := range SortStringMapKeys(req.Tags) {
			tagMap[key] = cty.StringVal(req.Tags[key])
		}
		cluster.SetAttributeValue("tags", cty.ObjectVal(tagMap))
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

func writeReplicationSpec(cluster *hclwrite.Body, specs []admin.ReplicationSpec20240805) error {
	var allSpecs []cty.Value

	for _, spec := range specs {
		specMap := make(map[string]cty.Value)

		if spec.ZoneName != nil {
			specMap["zone_name"] = cty.StringVal(*spec.ZoneName)
		}

		var rcList []cty.Value
		for _, rc := range spec.GetRegionConfigs() {
			if rc.Priority == nil {
				rc.SetPriority(7)
			}

			rcMap := map[string]cty.Value{
				"priority":      cty.NumberIntVal(int64(*rc.Priority)),
				"provider_name": cty.StringVal(*rc.ProviderName),
				"region_name":   cty.StringVal(*rc.RegionName),
			}
			if rc.BackingProviderName != nil {
				rcMap["backing_provider_name"] = cty.StringVal(*rc.BackingProviderName)
			}

			if rc.AutoScaling == nil {
				rcMap["auto_scaling"] = cty.ObjectVal(map[string]cty.Value{
					"disk_gb_enabled": cty.BoolVal(false),
				})
			} else {
				as := rc.GetAutoScaling()
				asDisk := as.GetDiskGB()
				if as.Compute != nil {
					return fmt.Errorf("auto_scaling.compute is not supported yet %v", as)
				}
				rcMap["auto_scaling"] = cty.ObjectVal(map[string]cty.Value{
					"disk_gb_enabled": cty.BoolVal(asDisk.GetEnabled()),
				})
			}

			es := rc.GetElectableSpecs()
			esMap := map[string]cty.Value{}
			if es.InstanceSize != nil {
				esMap["instance_size"] = cty.StringVal(*es.InstanceSize)
			}
			if es.NodeCount != nil {
				esMap["node_count"] = cty.NumberIntVal(int64(*es.NodeCount))
			}
			if es.EbsVolumeType != nil && *es.EbsVolumeType != "" {
				esMap["ebs_volume_type"] = cty.StringVal(*es.EbsVolumeType)
			}
			if es.DiskIOPS != nil {
				esMap["disk_iops"] = cty.NumberIntVal(int64(*es.DiskIOPS))
			}
			if len(esMap) > 0 {
				rcMap["electable_specs"] = cty.ObjectVal(esMap)
			}

			ros := rc.GetReadOnlySpecs()
			roMap := map[string]cty.Value{}
			if ros.InstanceSize != nil {
				roMap["instance_size"] = cty.StringVal(*ros.InstanceSize)
			}
			if ros.NodeCount != nil && *ros.NodeCount != 0 {
				roMap["node_count"] = cty.NumberIntVal(int64(*ros.NodeCount))
			}
			if ros.DiskIOPS != nil {
				roMap["disk_iops"] = cty.NumberIntVal(int64(*ros.DiskIOPS))
			}
			if len(roMap) > 0 {
				rcMap["read_only_specs"] = cty.ObjectVal(roMap)
			}

			rcList = append(rcList, cty.ObjectVal(rcMap))
		}
		// Use TupleVal instead of ListVal so region/spec objects can have different fields without type conflicts.
		specMap["region_configs"] = cty.TupleVal(rcList)
		allSpecs = append(allSpecs, cty.ObjectVal(specMap))
	}

	cluster.SetAttributeValue("replication_specs", cty.TupleVal(allSpecs))
	return nil
}

func validateAdvancedConfig(cfg map[string]any) error {
	for k := range cfg {
		if !knownAdvancedConfig[k] {
			return fmt.Errorf("unknown advanced configuration key: %s", k)
		}
	}
	return nil
}
