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
	"github.com/zclconf/go-cty/cty"
)

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
