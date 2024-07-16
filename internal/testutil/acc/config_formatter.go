package acc

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go.mongodb.org/atlas-sdk/v20240530002/admin"
)

func FormatToHCLMap(m map[string]string, indent, varName string) string {
	if m == nil {
		return ""
	}
	lines := []string{
		fmt.Sprintf("%s%s = {", indent, varName),
	}
	indentKeyValues := indent + "\t"

	for _, k := range sortStringMapKeys(m) {
		v := m[k]
		lines = append(lines, fmt.Sprintf("%s%s = %[3]q", indentKeyValues, k, v))
	}
	lines = append(lines, fmt.Sprintf("%s}", indent))
	return strings.Join(lines, "\n")
}

func FormatToHCLLifecycleIgnore(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	ignoreParts := []string{}
	for _, ignoreKey := range keys {
		ignoreParts = append(ignoreParts, fmt.Sprintf("\t\t\t%s,", ignoreKey))
	}
	lines := []string{
		"\tlifecycle {",
		"\t\tignore_changes = [",
		strings.Join(ignoreParts, "\n"),
		"\t\t]",
		"\t}",
	}
	return strings.Join(lines, "\n")
}

func sortStringMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
func sortStringMapKeysAny(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var (
	ClusterAdvConfigOplogMinRetentionHours = "oplog_min_retention_hours"
	knownAdvancedConfig                    = map[string]bool{
		ClusterAdvConfigOplogMinRetentionHours: true,
	}
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
	specs := make([]admin.ReplicationSpec, len(specRequests))
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
		for _, key := range sortStringMapKeys(req.Tags) {
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

func writeReplicationSpec(cluster *hclwrite.Body, spec admin.ReplicationSpec) error {
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

// addPrimitiveAttributesViaJSON  adds "primitive" bool/string/int/float attributes of a struct.
func addPrimitiveAttributesViaJSON(b *hclwrite.Body, obj any) error {
	var objMap map[string]any
	inrec, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, &objMap)
	if err != nil {
		return err
	}
	addPrimitiveAttributes(b, objMap)
	return nil
}

func addPrimitiveAttributes(b *hclwrite.Body, values map[string]any) {
	for _, keyCamel := range sortStringMapKeysAny(values) {
		key := ToSnakeCase(keyCamel)
		value := values[keyCamel]
		switch value := value.(type) {
		case bool:
			b.SetAttributeValue(key, cty.BoolVal(value))
		case string:
			if value != "" {
				b.SetAttributeValue(key, cty.StringVal(value))
			}
		case int:
			b.SetAttributeValue(key, cty.NumberIntVal(int64(value)))
		// int gets parsed as float64 for json
		case float64:
			b.SetAttributeValue(key, cty.NumberIntVal(int64(value)))
		default:
			continue
		}
	}
}

// Sometimes it is easier to set a value using hcl/tf syntax instead of creating complex values like list hcl.Traversal.
func setAttributeHcl(body *hclwrite.Body, tfExpression string) error {
	src := []byte(tfExpression)

	f, diags := hclwrite.ParseConfig(src, "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("extract attribute error %s\nparsing %s", diags, tfExpression)
	}
	expressionAttributes := f.Body().Attributes()
	if len(expressionAttributes) != 1 {
		return fmt.Errorf("must be a single attribute in expression: %s", tfExpression)
	}
	tokens := hclwrite.Tokens{}
	for _, attr := range expressionAttributes {
		tokens = attr.BuildTokens(tokens)
	}
	if len(tokens) == 0 {
		return fmt.Errorf("no tokens found for expression %s", tfExpression)
	}
	var attributeName string
	valueTokens := []*hclwrite.Token{}
	equalFound := false
	for _, token := range tokens {
		if attributeName == "" && token.Type == hclsyntax.TokenIdent {
			attributeName = string(token.Bytes)
		}
		if equalFound {
			valueTokens = append(valueTokens, token)
		}
		if token.Type == hclsyntax.TokenEqual {
			equalFound = true
		}
	}
	if attributeName == "" {
		return fmt.Errorf("unable to find the attribute name set for expr=%s", tfExpression)
	}
	if len(valueTokens) == 0 {
		return fmt.Errorf("unable to find the attribute value set for expr=%s", tfExpression)
	}
	body.SetAttributeRaw(attributeName, valueTokens)
	return nil
}
