package acc

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/zclconf/go-cty/cty"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckResourceAttrSchemaV2(isAcc bool, name, key, value string) resource.TestCheckFunc {
	if skipChecks(isAcc, name) {
		return testCheckFuncAlwaysPass
	}
	return resource.TestCheckResourceAttr(name, AttrNameToSchemaV2(isAcc, key), value)
}

func TestCheckResourceAttrSetSchemaV2(isAcc bool, name, key string) resource.TestCheckFunc {
	if skipChecks(isAcc, name) {
		return testCheckFuncAlwaysPass
	}
	return resource.TestCheckResourceAttrSet(name, AttrNameToSchemaV2(isAcc, key))
}

func TestCheckResourceAttrWithSchemaV2(isAcc bool, name, key string, checkValueFunc resource.CheckResourceAttrWithFunc) resource.TestCheckFunc {
	if skipChecks(isAcc, name) {
		return testCheckFuncAlwaysPass
	}
	return resource.TestCheckResourceAttrWith(name, AttrNameToSchemaV2(isAcc, key), checkValueFunc)
}

func TestCheckTypeSetElemNestedAttrsSchemaV2(isAcc bool, name, key string, values map[string]string) resource.TestCheckFunc {
	if skipChecks(isAcc, name) {
		return testCheckFuncAlwaysPass
	}
	return resource.TestCheckTypeSetElemNestedAttrs(name, AttrNameToSchemaV2(isAcc, key), values)
}

func testCheckFuncAlwaysPass(*terraform.State) error {
	return nil
}

func AddAttrChecksSchemaV2(isAcc bool, name string, checks []resource.TestCheckFunc, mapChecks map[string]string) []resource.TestCheckFunc {
	if skipChecks(isAcc, name) {
		return []resource.TestCheckFunc{}
	}
	return AddAttrChecks(name, checks, ConvertToSchemaV2AttrsMap(isAcc, mapChecks))
}

func AddAttrSetChecksSchemaV2(isAcc bool, name string, checks []resource.TestCheckFunc, attrNames ...string) []resource.TestCheckFunc {
	if skipChecks(isAcc, name) {
		return []resource.TestCheckFunc{}
	}
	return AddAttrSetChecks(name, checks, ConvertToSchemaV2AttrsSet(isAcc, attrNames)...)
}

func AddAttrChecksPrefixSchemaV2(isAcc bool, name string, checks []resource.TestCheckFunc, mapChecks map[string]string, prefix string, skipNames ...string) []resource.TestCheckFunc {
	if skipChecks(isAcc, name) {
		return []resource.TestCheckFunc{}
	}
	return AddAttrChecksPrefix(name, checks, ConvertToSchemaV2AttrsMap(isAcc, mapChecks), prefix, skipNames...)
}

func skipChecks(isAcc bool, name string) bool {
	if !config.AdvancedClusterV2Schema() || !isAcc {
		return false
	}
	return strings.HasPrefix(name, "data.mongodbatlas_advanced_cluster")
}

func ConvertToSchemaV2AttrsMap(isAcc bool, attrsMap map[string]string) map[string]string {
	if !config.AdvancedClusterV2Schema() || !isAcc {
		return attrsMap
	}
	ret := make(map[string]string, len(attrsMap))
	for name, value := range attrsMap {
		ret[AttrNameToSchemaV2(isAcc, name)] = value
	}
	return ret
}

func ConvertToSchemaV2AttrsSet(isAcc bool, attrsSet []string) []string {
	if !config.AdvancedClusterV2Schema() || !isAcc {
		return attrsSet
	}
	ret := make([]string, 0, len(attrsSet))
	for _, name := range attrsSet {
		ret = append(ret, AttrNameToSchemaV2(isAcc, name))
	}
	return ret
}

var tpfSingleNestedAttrs = []string{
	"analytics_specs",
	"electable_specs",
	"read_only_specs",
	"auto_scaling", // includes analytics_auto_scaling
	"advanced_configuration",
	"bi_connector_config",
}

func AttrNameToSchemaV2(isAcc bool, name string) string {
	if !config.AdvancedClusterV2Schema() || !isAcc {
		return name
	}
	for _, singleAttrName := range tpfSingleNestedAttrs {
		name = strings.ReplaceAll(name, singleAttrName+".0", singleAttrName)
	}
	return name
}

func ConvertAdvancedClusterToSchemaV2(t *testing.T, isAcc bool, def string) string {
	t.Helper()
	if !config.AdvancedClusterV2Schema() || !isAcc {
		return def
	}
	parse := getDefParser(t, def)
	for _, resource := range parse.Body().Blocks() {
		isResource := resource.Type() == "resource"
		resourceName := resource.Labels()[0]
		if !isResource || resourceName != "mongodbatlas_advanced_cluster" {
			continue
		}
		writeBody := resource.Body()
		convertAttrs(t, "labels", writeBody, true, getAttrVal)
		convertAttrs(t, "tags", writeBody, true, getAttrVal)
		convertAttrs(t, "replication_specs", writeBody, true, getReplicationSpecs)
		convertAttrs(t, "advanced_configuration", writeBody, false, getAttrVal)
		convertAttrs(t, "bi_connector_config", writeBody, false, getAttrVal)
	}
	content := parse.Bytes()
	return string(content)
}

func AssertEqualHCL(t *testing.T, expected, actual string, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, canonicalHCL(t, expected), canonicalHCL(t, actual), msgAndArgs...)
}

func convertAttrs(t *testing.T, name string, writeBody *hclwrite.Body, isList bool, getOneAttr func(*testing.T, *hclsyntax.Body) cty.Value) {
	t.Helper()
	var vals []cty.Value
	for {
		match := writeBody.FirstMatchingBlock(name, nil)
		if match == nil {
			break
		}
		vals = append(vals, getOneAttr(t, getBlockBody(t, match)))
		writeBody.RemoveBlock(match) // TODO: RemoveBlock doesn't remove newline just after the block so an extra line is added
	}
	if len(vals) == 0 {
		return
	}
	if isList {
		writeBody.SetAttributeValue(name, cty.TupleVal(vals))
	} else {
		assert.Len(t, vals, 1, "can be only one of %s", name)
		writeBody.SetAttributeValue(name, vals[0])
	}
}

func getReplicationSpecs(t *testing.T, body *hclsyntax.Body) cty.Value {
	t.Helper()
	const name = "region_configs"
	var vals []cty.Value
	for _, block := range body.Blocks {
		assert.Equal(t, name, block.Type, "unexpected block type: %s", block.Type)
		vals = append(vals, getAttrVal(t, block.Body))
	}
	return cty.ObjectVal(map[string]cty.Value{
		name: cty.TupleVal(vals),
	})
}

func getAttrVal(t *testing.T, body *hclsyntax.Body) cty.Value {
	t.Helper()
	ret := make(map[string]cty.Value)
	for name, attr := range body.Attributes {
		val, diags := attr.Expr.Value(nil)
		require.False(t, diags.HasErrors(), "failed to parse attribute %s: %s", name, diags.Error())
		ret[name] = val
	}
	for _, block := range body.Blocks {
		ret[block.Type] = getAttrVal(t, block.Body)
	}
	return cty.ObjectVal(ret)
}

func canonicalHCL(t *testing.T, def string) string {
	t.Helper()
	return string(getDefParser(t, def).Bytes())
}

func getDefParser(t *testing.T, def string) *hclwrite.File {
	t.Helper()
	parser, diags := hclwrite.ParseConfig([]byte(def), "", hcl.Pos{Line: 1, Column: 1})
	require.False(t, diags.HasErrors(), "failed to parse def: %s", diags.Error())
	return parser
}

func getBlockBody(t *testing.T, block *hclwrite.Block) *hclsyntax.Body {
	t.Helper()
	parser, diags := hclparse.NewParser().ParseHCL(block.Body().BuildTokens(nil).Bytes(), "")
	require.False(t, diags.HasErrors(), "failed to parse block: %s", diags.Error())

	body, ok := parser.Body.(*hclsyntax.Body)
	require.True(t, ok, "unexpected *hclsyntax.Body type: %T", parser.Body)
	return body
}
