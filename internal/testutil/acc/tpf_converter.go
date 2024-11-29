package acc

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ConvertAdvancedClusterToTPF(t *testing.T, def string) string {
	t.Helper()
	if !IsTPFAdvancedCluster() {
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
		generateAllAttributeSpecs(t, "labels", writeBody, func(body *hclsyntax.Body) cty.Value {
			return getVal(t, body)
		})
		generateAllAttributeSpecs(t, "tags", writeBody, func(body *hclsyntax.Body) cty.Value {
			return getVal(t, body)
		})
		generateAllAttributeSpecs(t, "replication_specs", writeBody, func(body *hclsyntax.Body) cty.Value {
			return getOneReplicationSpecs(t, body)
		})
	}
	content := parse.Bytes()
	return string(content)
}

func AssertEqualHCL(t *testing.T, expected, actual string, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, canonicalHCL(t, expected), canonicalHCL(t, actual), msgAndArgs...)
}

func generateAllAttributeSpecs(t *testing.T, name string, writeBody *hclwrite.Body, getOneAttr func(*hclsyntax.Body) cty.Value) {
	t.Helper()
	var vals []cty.Value
	for {
		match := writeBody.FirstMatchingBlock(name, nil)
		if match == nil {
			break
		}
		vals = append(vals, getOneAttr(getBlockBody(t, match)))
		writeBody.RemoveBlock(match) // TODO: RemoveBlock doesn't remove newline just after the block so an extra line is added
	}
	if len(vals) > 0 {
		writeBody.SetAttributeValue(name, cty.TupleVal(vals))
	}
}

func getOneReplicationSpecs(t *testing.T, body *hclsyntax.Body) cty.Value {
	t.Helper()
	const name = "region_configs"
	var vals []cty.Value
	for _, block := range body.Blocks {
		assert.Equal(t, name, block.Type, "unexpected block type: %s", block.Type)
		oneRegionConfigs := getVal(t, block.Body)
		vals = append(vals, oneRegionConfigs)
	}
	return cty.ObjectVal(map[string]cty.Value{
		name: cty.TupleVal(vals),
	})
}

func getVal(t *testing.T, body *hclsyntax.Body) cty.Value {
	t.Helper()
	ret := make(map[string]cty.Value)
	for name, attr := range body.Attributes {
		val, diags := attr.Expr.Value(nil)
		require.False(t, diags.HasErrors(), "failed to parse attribute %s: %s", name, diags.Error())
		ret[name] = val
	}
	for _, block := range body.Blocks {
		ret[block.Type] = getVal(t, block.Body)
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
