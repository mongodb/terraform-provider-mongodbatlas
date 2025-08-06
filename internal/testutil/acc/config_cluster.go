package acc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312005/admin"

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

	if req.RetainBackupsEnabled {
		clusterRootAttributes["retain_backups_enabled"] = req.RetainBackupsEnabled
	}
	addPrimitiveAttributes(cluster, clusterRootAttributes)
	cluster.AppendNewline()
	if len(req.AdvancedConfiguration) > 0 {
		if err := writeAdvancedConfiguration(cluster, req.AdvancedConfiguration); err != nil {
			return "", "", "", err
		}

	}
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

func recursiveJSONToCty(raw any) (cty.Value, error) {
	switch v := raw.(type) {
	case map[string]any:
		obj := make(map[string]cty.Value, len(v))
		for key, val := range v {
			snake := ToSnakeCase(key)
			conv, err := recursiveJSONToCty(val)
			if err != nil {
				return cty.NilVal, err
			}
			obj[snake] = conv
		}
		return cty.ObjectVal(obj), nil

	case []any:
		list := make([]cty.Value, 0, len(v))
		for _, elem := range v {
			conv, err := recursiveJSONToCty(elem)
			if err != nil {
				return cty.NilVal, err
			}
			list = append(list, conv)
		}
		return cty.ListVal(list), nil

	case bool:
		return cty.BoolVal(v), nil
	case string:
		return cty.StringVal(v), nil
	case float64:
		if float64(int64(v)) == v {
			return cty.NumberIntVal(int64(v)), nil
		}
		return cty.NumberFloatVal(v), nil

	case nil:
		return cty.NullVal(cty.DynamicPseudoType), nil

	default:
		return cty.NilVal, fmt.Errorf("unsupported JSON value type %T", v)
	}
}

func structToCtyObject(obj any) (map[string]cty.Value, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	top, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected JSON object, got %T", raw)
	}

	result := make(map[string]cty.Value, len(top))
	for key, val := range top {
		snake := ToSnakeCase(key)
		conv, err := recursiveJSONToCty(val)
		if err != nil {
			return nil, err
		}
		if conv.IsNull() {
			continue
		}
		result[snake] = conv
	}
	return result, nil
}

func writeReplicationSpec(cluster *hclwrite.Body, specs []admin.ReplicationSpec20240805) error {
	var allSpecs []cty.Value

	for _, spec := range specs {
		specMap, err := structToCtyObject(spec)
		if err != nil {
			return err
		}
		delete(specMap, "region_configs") // Handle region_configs separately below

		var rcList []cty.Value
		for _, rc := range spec.GetRegionConfigs() {
			if rc.Priority == nil {
				rc.SetPriority(7)
			}

			rcMap, err := structToCtyObject(rc)
			if err != nil {
				return err
			}

			delete(rcMap, "auto_scaling")

			if rc.AutoScaling == nil {
				rcMap["auto_scaling"] = cty.ObjectVal(map[string]cty.Value{
					"disk_gb_enabled": cty.BoolVal(false),
				})
			} else {
				autoScaling := rc.GetAutoScaling()
				asDisk := autoScaling.GetDiskGB()
				if autoScaling.Compute != nil {
					return fmt.Errorf("auto_scaling.compute is not supportd yet %v", autoScaling)
				}
				rcMap["auto_scaling"] = cty.ObjectVal(map[string]cty.Value{
					"disk_gb_enabled": cty.BoolVal(asDisk.GetEnabled()),
				})
			}

			nodeSpec := rc.GetElectableSpecs()
			esMap, err := structToCtyObject(nodeSpec)
			if err != nil {
				return err
			}
			rcMap["electable_specs"] = cty.ObjectVal(esMap)

			readOnlySpecs := rc.GetReadOnlySpecs()
			if readOnlySpecs.GetNodeCount() != 0 {
				roMap, err := structToCtyObject(readOnlySpecs)
				if err != nil {
					return err
				}
				rcMap["read_only_specs"] = cty.ObjectVal(roMap)
			} else {
				delete(rcMap, "read_only_specs")
			}

			rcList = append(rcList, cty.ObjectVal(rcMap))
		}

		specMap["region_configs"] = cty.ListVal(rcList)

		allSpecs = append(allSpecs, cty.ObjectVal(specMap))
	}

	cluster.SetAttributeValue("replication_specs", cty.ListVal(allSpecs))

	return nil
}

func writeAdvancedConfiguration(
	cluster *hclwrite.Body,
	advConfig map[string]any,
) error {
	if len(advConfig) == 0 {
		return nil
	}

	for _, key := range sortStringMapKeysAny(advConfig) {
		if !knownAdvancedConfig[key] {
			return fmt.Errorf("unknown key in advanced configuration: %s", key)
		}
	}

	confMap, err := structToCtyObject(advConfig)
	if err != nil {
		return fmt.Errorf("failed to convert advanced configuration: %w", err)
	}

	cluster.SetAttributeValue("advanced_configuration", cty.ObjectVal(confMap))
	cluster.AppendNewline()

	return nil
}
