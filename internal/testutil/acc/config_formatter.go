package acc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
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

func AddNonZeroAttributes(b *hclwrite.Body, obj any) error {
	var objMap map[string]interface{}
	inrec, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(inrec, &objMap)
	if err != nil {
		return err
	}
	addNonEmptyAttributes(b, objMap)
	return nil
}

func addNonEmptyAttributes(b *hclwrite.Body, values map[string]any) {
	for _, keyCamel := range sortStringMapKeysAny(values) {
		key := ToSnakeCase(keyCamel)
		value := values[keyCamel]
		switch value := value.(type) {
		case bool:
			b.SetAttributeValue(key, cty.BoolVal(value))
		case string:
			if value == "" {
				continue
			}
			b.SetAttributeValue(key, cty.StringVal(value))
		case int:
			if value == 0 {
				continue
			}
			b.SetAttributeValue(key, cty.NumberIntVal(int64(value)))
		case float64:
			b.SetAttributeValue(key, cty.NumberIntVal(int64(value)))
		default:
			continue
		}
	}
}

func ClusterResourceHcl(projectID string, req *ClusterRequest, specs []admin.ReplicationSpec) (configStr, clusterName string, err error) {
	if len(specs) == 0 {
		specs = append(specs, ReplicationSpec(nil))
	}
	if req == nil {
		req = new(ClusterRequest)
	}
	if req.ProviderName == "" {
		req.ProviderName = constant.AWS
	}
	clusterName = req.ClusterNameExplicit
	if clusterName == "" {
		clusterName = RandomClusterName()
	}
	clusterTypeStr := "REPLICASET"
	if req.Geosharded {
		clusterTypeStr = "GEOSHARDED"
	}

	f := hclwrite.NewEmptyFile()
	root := f.Body()
	cluster := root.AppendNewBlock("resource", []string{"mongodbatlas_advanced_cluster", "cluster_info"}).Body()
	addNonEmptyAttributes(cluster, map[string]any{
		"project_id":     projectID,
		"clusterType":    clusterTypeStr,
		"name":           clusterName,
		"backup_enabled": req.CloudBackup,
	})
	cluster.AppendNewline()
	for i, spec := range specs {
		err = writeReplicationSpec(cluster, spec, req.ProviderName)
		if err != nil {
			return "", "", fmt.Errorf("error writing hcl for replication spec %d: %w", i, err)
		}
	}
	cluster.AppendNewline()
	if req.ResourceDependencyName != "" {
		if !strings.Contains(req.ResourceDependencyName, ".") {
			return "", "", fmt.Errorf("req.ResourceDependencyName must have a '.'")
		}
		dependsOnAttr, err := extractValueTokens(fmt.Sprintf("dummy = [%s]", req.ResourceDependencyName))
		if err != nil {
			return "", "", err
		}
		cluster.SetAttributeRaw("depends_on", dependsOnAttr)
	}
	return "\n" + string(f.Bytes()), clusterName, err
}

func writeReplicationSpec(cluster *hclwrite.Body, spec admin.ReplicationSpec, providerName string) error {
	replicationBlock := cluster.AppendNewBlock("replication_specs", nil).Body()
	err := AddNonZeroAttributes(replicationBlock, spec)
	if err != nil {
		return err
	}
	for _, rc := range spec.GetRegionConfigs() {
		if rc.ProviderName == nil {
			rc.SetProviderName(providerName)
		}
		if rc.Priority == nil {
			rc.SetPriority(7)
		}
		replicationBlock.AppendNewline()
		rcBlock := replicationBlock.AppendNewBlock("region_configs", nil).Body()
		err = AddNonZeroAttributes(rcBlock, rc)
		if err != nil {
			return err
		}
		autoScalingBlock := rcBlock.AppendNewBlock("auto_scaling", nil).Body()
		if rc.AutoScaling == nil {
			autoScalingBlock.SetAttributeValue("disk_gb_enabled", cty.BoolVal(false))
		} else {
			autoScaling := rc.GetAutoScaling()
			return fmt.Errorf("auto_scaling on replication spec is not supportd yet %v", autoScaling)
		}
		nodeSpec := rc.GetElectableSpecs()
		nodeSpecBlock := rcBlock.AppendNewBlock("electable_specs", nil).Body()
		err = AddNonZeroAttributes(nodeSpecBlock, nodeSpec)
	}
	return err
}

func extractValueTokens(tfExpression string) ([]*hclwrite.Token, error) {
	src := []byte(tfExpression)

	f, diags := hclwrite.ParseConfig(src, "", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("extract attribute error %s\nparsing %s", diags, tfExpression)
	}
	for _, attr := range f.Body().Attributes() {
		tokens := hclwrite.Tokens{}
		tokens = attr.BuildTokens(tokens)
		equalFound := false
		returnTokens := []*hclwrite.Token{}
		for _, token := range tokens {
			if equalFound {
				returnTokens = append(returnTokens, token)
			}
			if token.Type == hclsyntax.TokenEqual {
				equalFound = true
			}
		}
		return returnTokens, nil
	}
	return nil, fmt.Errorf("extract attribute error parsing %s", tfExpression)
}
