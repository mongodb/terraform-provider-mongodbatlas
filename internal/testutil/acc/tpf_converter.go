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
	parse, diags := hclwrite.ParseConfig([]byte(def), "", hcl.Pos{Line: 1, Column: 1})
	require.False(t, diags.HasErrors(), "failed to parse HCL: %s", diags.Error())
	body := parse.Body()
	for _, resource := range body.Blocks() {
		isResource := resource.Type() == "resource"
		resourceType := resource.Labels()[0]
		if !isResource || resourceType != "mongodbatlas_advanced_cluster" {
			continue
		}
		writeBody := resource.Body()
		generateReplicationSpecs(t, writeBody)
	}
	return string(parse.Bytes())
}

func AssertEqualHCL(t *testing.T, expected, actual string, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, canonicalHCL(t, expected), canonicalHCL(t, actual), msgAndArgs...)
}

func canonicalHCL(t *testing.T, def string) string {
	t.Helper()
	parse, diags := hclwrite.ParseConfig([]byte(def), "", hcl.Pos{Line: 1, Column: 1})
	require.False(t, diags.HasErrors(), "failed to parse HCL: %s", diags.Error())
	return string(parse.Bytes())
}

func generateReplicationSpecs(t *testing.T, writeBody *hclwrite.Body) {
	t.Helper()
	const name = "replication_specs"
	var vals []cty.Value
	for {
		match := writeBody.FirstMatchingBlock(name, nil)
		if match == nil {
			break
		}
		parse, diags := hclparse.NewParser().ParseHCL(match.Body().BuildTokens(nil).Bytes(), "")
		require.False(t, diags.HasErrors(), "failed to parse %s: %s", name, diags.Error())
		body, ok := parse.Body.(*hclsyntax.Body)
		require.True(t, ok, "unexpected hclsyntax.Body type: %T", parse.Body)
		vals = append(vals, getReplicationSpecsAttribute(t, body))
		writeBody.RemoveBlock(match)
	}
	require.NotEmpty(t, vals, "there must be at least one %s block", name)
	writeBody.SetAttributeValue(name, cty.ListVal(vals))
}

func getReplicationSpecsAttribute(t *testing.T, body *hclsyntax.Body) cty.Value {
	t.Helper()
	const name = "region_configs"
	var vals []cty.Value

	for _, block := range body.Blocks {
		vals = append(vals, getRegionConfigsAttribute(t, block))
	}
	return cty.ObjectVal(map[string]cty.Value{
		name: cty.ListVal(vals),
	})
}

func getRegionConfigsAttribute(t *testing.T, block *hclsyntax.Block) cty.Value {
	t.Helper()
	valMap := make(map[string]cty.Value)
	for name, attr := range block.Body.Attributes {
		val, diags := attr.Expr.Value(nil)
		require.False(t, diags.HasErrors(), "failed to parse attribute %s: %s", name, diags.Error())
		valMap[name] = val
	}
	return cty.ObjectVal(valMap)
}
