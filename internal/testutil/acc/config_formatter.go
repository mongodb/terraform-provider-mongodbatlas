package acc

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	localHcl "github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/hcl"
)

// ConfigAddResourceStr is useful when you need to add one or more attributes to a resource block.
func ConfigAddResourceStr(t *testing.T, hclConfig, resourceID, extraResourceStr string) string {
	t.Helper()
	resourceParts := strings.Split(resourceID, ".")
	if len(resourceParts) != 2 {
		t.Fatalf("resourceID must be in the format <type>.<name>, got %s", resourceID)
	}
	resourceType := resourceParts[0]
	resourceName := resourceParts[1]
	resourceBlockDef := fmt.Sprintf("resource %q %q {", resourceType, resourceName)
	if !strings.Contains(hclConfig, resourceBlockDef) {
		t.Fatalf("resource block %q not found in config: %s", resourceBlockDef, hclConfig)
	}
	resourceBlockDefWithExtraResourceStr := fmt.Sprintf("%s\n%s\n", resourceBlockDef, extraResourceStr)
	hclConfigModified := strings.Replace(hclConfig, resourceBlockDef, resourceBlockDefWithExtraResourceStr, 1)
	return localHcl.PrettyHCL(t, hclConfigModified)
}

func FormatToHCLMap(m map[string]string, indent, varName string) string {
	if m == nil {
		return ""
	}
	lines := []string{
		fmt.Sprintf("%s%s = {", indent, varName),
	}
	indentKeyValues := indent + "\t"

	for _, k := range SortStringMapKeys(m) {
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

// setAttributeHcl inserts a raw HCL assignment into the body by parsing a snippet like:
//
//	project_id = mongodbatlas_project.test.id or depends_on = [mongodbatlas_project.test.id] etc
//
// and copying its tokens directly. Used for expressions or references that can’t emit simple Go literals.
func setAttributeHcl(body *hclwrite.Body, tfExpression string) error {
	src := []byte(tfExpression)

	f, diags := hclwrite.ParseConfig(src, "", hcl.InitialPos)
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

// setAttributes iterates over attrs, snake-cases each key, converts the value
// with toCtyValue, and calls body.SetAttributeValue.
func setAttributes(body *hclwrite.Body, attrs map[string]any) {
	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, camel := range keys {
		key := ToSnakeCase(camel)
		if cv, ok := toCtyValue(attrs[camel]); ok {
			body.SetAttributeValue(key, cv)
		}
	}
}

// toCtyValue handles:
//   - bool, string, int, float64
//   - map[string]any (recursively)
func toCtyValue(v any) (cty.Value, bool) {
	switch v := v.(type) {
	case bool:
		return cty.BoolVal(v), true
	case string:
		if v == "" {
			return cty.NullVal(cty.String), false
		}
		return cty.StringVal(v), true
	case int:
		return cty.NumberIntVal(int64(v)), true
	case float64:
		return cty.NumberIntVal(int64(v)), true
	case map[string]any:
		if len(v) == 0 {
			return cty.NullVal(cty.EmptyObject), false
		}
		obj := make(map[string]cty.Value, len(v))
		// sort keys for deterministic output
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if cv, ok := toCtyValue(v[k]); ok {
				obj[ToSnakeCase(k)] = cv
			}
		}
		return cty.ObjectVal(obj), true
	default:
		return cty.NilVal, false
	}
}
