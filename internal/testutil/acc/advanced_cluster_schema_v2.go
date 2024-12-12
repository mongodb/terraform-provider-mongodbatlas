package acc

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/zclconf/go-cty/cty"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckResourceAttrSchemaV2(isAcc bool, name, key, value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(name, key, value)
}

func TestCheckResourceAttrSetSchemaV2(isAcc bool, name, key string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(name, key)
}

func TestCheckResourceAttrWithSchemaV2(isAcc bool, name, key string, checkValueFunc resource.CheckResourceAttrWithFunc) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(name, key, checkValueFunc)
}

func TestCheckTypeSetElemNestedAttrsSchemaV2(isAcc bool, name, attr string, values map[string]string) resource.TestCheckFunc {
	return resource.TestCheckTypeSetElemNestedAttrs(name, attr, values)
}

// AddAttrChecksSchemaV2 is like AddAttrChecks but adding V2 schema support
func AddAttrChecksSchemaV2(isAcc bool, name string, checks []resource.TestCheckFunc, mapChecks map[string]string) []resource.TestCheckFunc {
	return AddAttrChecks(name, checks, ConvertToTPFAttrsMap(mapChecks))
}

// AddAttrChecksSchemaV2 is like AddAttrSetChecks but adding V2 schema support
func AddAttrSetChecksSchemaV2(isAcc bool, name string, checks []resource.TestCheckFunc, attrNames ...string) []resource.TestCheckFunc {
	return AddAttrSetChecks(name, checks, ConvertToTPFAttrsSet(attrNames)...)
}

// AddAttrChecksPrefixSchemaV2 is like AddAttrChecksPrefix but adding V2 schema support
func AddAttrChecksPrefixSchemaV2(isAcc bool, name string, checks []resource.TestCheckFunc, mapChecks map[string]string, prefix string, skipNames ...string) []resource.TestCheckFunc {
	return AddAttrChecksPrefix(name, checks, ConvertToTPFAttrsMap(mapChecks), prefix, skipNames...)
}

func ConvertToTPFAttrsMap(attrsMap map[string]string) map[string]string {
	if !config.AdvancedClusterV2Schema() {
		return attrsMap
	}
	ret := make(map[string]string, len(attrsMap))
	for name, value := range attrsMap {
		ret[AttrNameToSchemaV2(name)] = value
	}
	return ret
}

func ConvertToTPFAttrsSet(attrsSet []string) []string {
	if !config.AdvancedClusterV2Schema() {
		return attrsSet
	}
	ret := make([]string, 0, len(attrsSet))
	for _, name := range attrsSet {
		ret = append(ret, AttrNameToSchemaV2(name))
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

func AttrNameToSchemaV2(name string) string {
	if !config.AdvancedClusterV2Schema() {
		return name
	}
	for _, singleAttrName := range tpfSingleNestedAttrs {
		name = strings.ReplaceAll(name, singleAttrName+".0", singleAttrName)
	}
	return name
}

func ConvertAdvancedClusterToTPF(t *testing.T, def string) string {
	t.Helper()
	if !config.AdvancedClusterV2Schema() {
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

func ConvertAdvancedClusterToTPFIfEnabled(t *testing.T, enabled bool, def string) string {
	t.Helper()
	if enabled {
		return ConvertAdvancedClusterToTPF(t, def)
	}
	return def
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
